package ocr

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"image"
	"math"
)

//go:embed hog_model.json
var glyphModelData []byte

type glyphModel struct {
	Version string      `json:"version"`
	Classes []int       `json:"classes"`
	Weights [][]float64 `json:"weights"`
	Bias    []float64   `json:"bias"`
}

type GlyphClass struct {
	ClassID int
	Margin  float64
}

// ClassifyExpertBox classifies an already localized character crop into an
// opaque DeepLontar ID. It does not detect characters, transliterate, or provide
// a calibrated confidence. It must not be applied to page proposals as text.
func (r *Recognizer) ClassifyExpertBox(ctx context.Context, crop image.Image) (GlyphClass, error) {
	var result GlyphClass
	if err := ctx.Err(); err != nil {
		return result, err
	}
	b := crop.Bounds()
	if b.Dx() == 0 || b.Dy() == 0 {
		return result, errors.New("empty character crop")
	}
	var pixels [16][16]float32
	var mean float32
	for y := range 16 {
		for x := range 16 {
			pixels[y][x] = float32(grayAt(crop, b.Min.X+(2*x+1)*b.Dx()/32, b.Min.Y+(2*y+1)*b.Dy()/32))
			mean += pixels[y][x]
		}
	}
	mean /= 256
	var cropNorm float32
	for y := range 16 {
		for x := range 16 {
			pixels[y][x] -= mean
			cropNorm += pixels[y][x] * pixels[y][x]
		}
	}
	if cropNorm == 0 {
		return result, errors.New("character crop has no gradients")
	}
	cropNorm = float32(math.Sqrt(float64(cropNorm)))
	for y := range 16 {
		for x := range 16 {
			pixels[y][x] /= cropNorm
		}
	}
	var features [128]float32
	for y := range 16 {
		for x := range 16 {
			gx, gy := float32(0), float32(0)
			if x > 0 && x < 15 {
				gx = pixels[y][x+1] - pixels[y][x-1]
			}
			if y > 0 && y < 15 {
				gy = pixels[y+1][x] - pixels[y-1][x]
			}
			// Training HOG uses numpy float32 arithmetic; keep the same precision.
			angle := float32(math.Mod(math.Atan2(float64(gy), float64(gx))+math.Pi, math.Pi))
			bin := int(math.Floor(float64(angle)*(8/math.Pi) + 0.0001))
			if bin > 7 {
				bin = 7
			}
			features[((y/4)*4+x/4)*8+bin] += float32(math.Hypot(float64(gx), float64(gy)))
		}
	}
	var norm float32
	for _, v := range features {
		norm += v * v
	}
	if norm == 0 {
		return result, errors.New("character crop has no gradients")
	}
	norm = float32(math.Sqrt(float64(norm)))
	for i := range features {
		features[i] /= norm
	}
	best, second := math.Inf(-1), math.Inf(-1)
	for i, weights := range r.glyph.Weights {
		score := r.glyph.Bias[i]
		for j, w := range weights {
			score += w * float64(features[j])
		}
		if score > best {
			second = best
			best = score
			result.ClassID = r.glyph.Classes[i]
		} else if score > second {
			second = score
		}
	}
	result.Margin = best - second
	return result, nil
}

func loadGlyphModel() (glyphModel, error) {
	var g glyphModel
	if err := json.Unmarshal(glyphModelData, &g); err != nil {
		return g, err
	}
	if len(g.Classes) < 2 || len(g.Weights) != len(g.Classes) || len(g.Bias) != len(g.Classes) {
		return g, errors.New("invalid glyph model")
	}
	for _, w := range g.Weights {
		if len(w) != 128 {
			return g, errors.New("invalid glyph feature dimensions")
		}
	}
	return g, nil
}
