// Package ocr finds possible text-bearing horizontal bands; it does not transcribe them.
package ocr

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"image"
	"image/jpeg"
	_ "image/png"
	"math"
	"sort"
	"time"
)

//go:embed model.json
var modelData []byte

type model struct {
	Version      string    `json:"version"`
	RowWeights   []float64 `json:"row_weights"`
	RowThreshold float64   `json:"row_threshold"`
}

type Recognizer struct {
	model model
	glyph glyphModel
}

type Candidate struct {
	Box   [4]float64 `json:"box"`
	Score float64    `json:"score"` // Regression score, not a probability.
}

type FocusImage struct {
	Data     []byte     `json:"-"`
	MimeType string     `json:"-"`
	Box      [4]float64 `json:"-"`
}

type Evidence struct {
	ModelVersion string       `json:"modelVersion"`
	ElapsedMS    int64        `json:"elapsedMs"`
	Candidates   []Candidate  `json:"candidates"`
	Warnings     []string     `json:"warnings"`
	FocusImages  []FocusImage `json:"-"`
}

func Load() (*Recognizer, error) {
	var m model
	if err := json.Unmarshal(modelData, &m); err != nil {
		return nil, err
	}
	if len(m.RowWeights) != 10 || m.Version == "" || math.IsNaN(m.RowThreshold) {
		return nil, errors.New("invalid OCR model")
	}
	glyph, err := loadGlyphModel()
	if err != nil {
		return nil, err
	}
	return &Recognizer{model: m, glyph: glyph}, nil
}

func (r *Recognizer) Recognize(ctx context.Context, data []byte) (Evidence, error) {
	start := time.Now()
	ev := Evidence{ModelVersion: r.model.Version, Warnings: []string{"Experimental text-band localization only; scores are not probabilities, glyph readings, or translations."}}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return ev, err
	}
	if cfg.Width < 30 || cfg.Height < 10 || int64(cfg.Width)*int64(cfg.Height) > 24000000 {
		return ev, errors.New("image dimensions outside OCR limit")
	}
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return ev, err
	}
	regions := []image.Rectangle{src.Bounds()}
	ratio := float64(cfg.Width) / float64(cfg.Height)
	if ratio < 3 || ratio > 7 {
		regions = leafRegions(src)
		ev.Warnings = append(ev.Warnings, "Leaf regions are estimated heuristically from background contrast; results may miss text.")
		if len(regions) == 0 {
			ev.Warnings = append(ev.Warnings, "No plausible horizontal leaf was found.")
			ev.ElapsedMS = time.Since(start).Milliseconds()
			return ev, nil
		}
	}
	original := src.Bounds()
	for _, bounds := range regions {
		if len(ev.FocusImages) >= 2 {
			break
		}
		const w, h, pad = 1500, 300, 15
		if bounds.Dx() < 30 || bounds.Dy() < 10 {
			continue
		}
		gray := make([]uint8, w*h)
		for y := 0; y < h; y++ {
			if err := ctx.Err(); err != nil {
				return ev, err
			}
			sy := bounds.Min.Y + (2*y+1)*bounds.Dy()/(2*h)
			for x := 0; x < w; x++ {
				sx := bounds.Min.X + (2*x+1)*bounds.Dx()/(2*w)
				gray[y*w+x] = grayAt(src, sx, sy)
			}
		}
		// Edge-replicated integral image implements the training 31x31 local mean.
		const pw, ph = w + 2*pad, h + 2*pad
		integral := make([]int64, (pw+1)*(ph+1))
		for y := 0; y < ph; y++ {
			sy := clamp(y-pad, 0, h-1)
			var sum int64
			for x := 0; x < pw; x++ {
				sx := clamp(x-pad, 0, w-1)
				sum += int64(gray[sy*w+sx])
				integral[(y+1)*(pw+1)+x+1] = integral[y*(pw+1)+x+1] + sum
			}
		}
		var density [10][30]float64
		for y := 0; y < h; y++ {
			if err := ctx.Err(); err != nil {
				return ev, err
			}
			for x := 0; x < w; x++ {
				a, b, c, d := y*(pw+1)+x, y*(pw+1)+x+31, (y+31)*(pw+1)+x, (y+31)*(pw+1)+x+31
				total := integral[d] - integral[b] - integral[c] + integral[a]
				if int64(gray[y*w+x])*961 < total-25*961 {
					density[y/30][x/50]++
				}
			}
		}
		var scores [10][30]float64
		for y := 0; y < 10; y++ {
			for x := 0; x < 30; x++ {
				s := r.model.RowWeights[0]
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						s += r.model.RowWeights[1+(dy+1)*3+dx+1] * density[clamp(y+dy, 0, 9)][clamp(x+dx, 0, 29)] / 1500
					}
				}
				scores[y][x] = s
			}
		}
		type window struct {
			x, y  int
			score float64
		}
		var windows []window
		for y := 0; y <= 7; y++ {
			for x := 0; x <= 15; x++ {
				var score float64
				positives := 0
				for dy := 0; dy < 3; dy++ {
					for dx := 0; dx < 15; dx++ {
						s := scores[y+dy][x+dx]
						if s >= r.model.RowThreshold {
							positives++
							score += s
						}
					}
				}
				if positives >= 5 {
					windows = append(windows, window{x, y, score / 45})
				}
			}
		}
		sort.SliceStable(windows, func(i, j int) bool { return windows[i].score > windows[j].score })
		var chosen []window
		for _, candidate := range windows {
			overlap := false
			for _, prior := range chosen {
				if candidate.x < prior.x+15 && prior.x < candidate.x+15 && candidate.y < prior.y+3 && prior.y < candidate.y+3 {
					overlap = true
					break
				}
			}
			if !overlap {
				chosen = append(chosen, candidate)
			}
			limit := 2
			if len(regions) > 1 {
				limit = 1
			}
			if len(chosen) == limit {
				break
			}
		}
		for _, b := range chosen {
			low, high := 0, 10
			box := [4]float64{float64(bounds.Min.X-original.Min.X+b.x*bounds.Dx()/30) / float64(original.Dx()), float64(bounds.Min.Y-original.Min.Y+low*bounds.Dy()/10) / float64(original.Dy()), float64(15*bounds.Dx()/30) / float64(original.Dx()), float64((high-low)*bounds.Dy()/10) / float64(original.Dy())}
			ev.Candidates = append(ev.Candidates, Candidate{Box: box, Score: b.score})
			x0 := bounds.Min.X + b.x*bounds.Dx()/30
			x1 := bounds.Min.X + (b.x+15)*bounds.Dx()/30
			y0 := bounds.Min.Y + low*bounds.Dy()/10
			y1 := bounds.Min.Y + high*bounds.Dy()/10
			outw := x1 - x0
			if outw > 1024 {
				outw = 1024
			}
			outh := int(math.Ceil(float64(y1-y0) * float64(outw) / float64(x1-x0)))
			if outh < 1 {
				outh = 1
			}
			out := image.NewRGBA(image.Rect(0, 0, outw, outh))
			for y := 0; y < outh; y++ {
				for x := 0; x < outw; x++ {
					out.Set(x, y, src.At(x0+(2*x+1)*(x1-x0)/(2*outw), y0+(2*y+1)*(y1-y0)/(2*outh)))
				}
			}
			var encoded bytes.Buffer
			if err := jpeg.Encode(&encoded, out, &jpeg.Options{Quality: 85}); err != nil {
				return ev, err
			}
			ev.FocusImages = append(ev.FocusImages, FocusImage{Data: encoded.Bytes(), MimeType: "image/jpeg", Box: box})
		}
		if len(chosen) == 0 {
			ev.Warnings = append(ev.Warnings, "No text regions met the model threshold.")
		}
	}
	ev.ElapsedMS = time.Since(start).Milliseconds()
	return ev, nil
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func grayAt(src image.Image, x, y int) uint8 {
	var r, g, b uint32
	switch im := src.(type) {
	case *image.YCbCr:
		c := im.YCbCrAt(x, y)
		r, g, b, _ = c.RGBA()
	case *image.RGBA:
		c := im.RGBAAt(x, y)
		r, g, b, _ = c.RGBA()
	case *image.NRGBA:
		c := im.NRGBAAt(x, y)
		r, g, b = uint32(c.R)*257, uint32(c.G)*257, uint32(c.B)*257
	case *image.Gray:
		return im.GrayAt(x, y).Y
	default:
		r, g, b, _ = src.At(x, y).RGBA()
	}
	return uint8((299*int(r>>8) + 587*int(g>>8) + 114*int(b>>8) + 500) / 1000)
}
