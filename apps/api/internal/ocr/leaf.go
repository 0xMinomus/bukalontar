package ocr

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"sort"
)

// leafRegions finds wide foreground strips against an approximately uniform
// corner background. This is not a learned detector and can miss tilted leaves.
func leafRegions(src image.Image) []image.Rectangle {
	b := src.Bounds()
	const width = 320
	height := clamp(b.Dy()*width/b.Dx(), 32, 320)
	pixels := make([]uint8, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixels[y*width+x] = grayAt(src, b.Min.X+(2*x+1)*b.Dx()/(2*width), b.Min.Y+(2*y+1)*b.Dy()/(2*height))
		}
	}
	corners := make([]int, 0, 4*8*8)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if (x < 8 || x >= width-8) && (y < 8 || y >= height-8) {
				corners = append(corners, int(pixels[y*width+x]))
			}
		}
	}
	sort.Ints(corners)
	background := corners[len(corners)/2]
	foreground := make([]bool, len(pixels))
	rows := make([]bool, height)
	for y := 0; y < height; y++ {
		count := 0
		for x := 0; x < width; x++ {
			delta := int(pixels[y*width+x]) - background
			if delta < 0 {
				delta = -delta
			}
			if delta >= 35 {
				foreground[y*width+x] = true
				count++
			}
		}
		rows[y] = count >= width*35/100
	}
	// Close short gaps caused by holes and lines of writing.
	for y := 1; y+2 < height; y++ {
		if rows[y-1] && rows[y+2] {
			rows[y] = true
			rows[y+1] = true
		}
	}
	var regions []image.Rectangle
	for y := 0; y < height; {
		if !rows[y] {
			y++
			continue
		}
		y0 := y
		for y < height && rows[y] {
			y++
		}
		y1 := y
		if y1-y0 < height*3/100 || y1-y0 < 3 {
			continue
		}
		x0, x1 := width, 0
		for x := 0; x < width; x++ {
			n := 0
			for yy := y0; yy < y1; yy++ {
				if foreground[yy*width+x] {
					n++
				}
			}
			if n*2 >= y1-y0 {
				if x < x0 {
					x0 = x
				}
				x1 = x + 1
			}
		}
		if x1-x0 < width*35/100 {
			continue
		}
		left := b.Min.X + x0*b.Dx()/width
		right := b.Min.X + (x1*b.Dx()+width-1)/width
		top := b.Min.Y + y0*b.Dy()/height
		bottom := b.Min.Y + (y1*b.Dy()+height-1)/height
		rw, rh := right-left, bottom-top
		if rw < 3*rh {
			continue
		}
		// DeepLontar scans have 5:1 aspect; keep long leaves undistorted.
		if rw > 7*rh {
			mid := (left + right) / 2
			left = mid - 7*rh/2
			right = left + 7*rh
		}
		regions = append(regions, image.Rect(left, top, right, bottom))
	}
	return regions
}

// LeafCrops returns up to max full-leaf enlargements in top-to-bottom order.
// The crops help the multimodal model read small strokes; they are not OCR text.
func (r *Recognizer) LeafCrops(ctx context.Context, data []byte, max int) ([]FocusImage, error) {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if max <= 0 {
		max = 2
	}
	regions := leafRegions(src)
	bounds := src.Bounds()
	crops := make([]FocusImage, 0, 2)
	for _, rect := range regions {
		if err := ctx.Err(); err != nil {
			return crops, err
		}
		rect = rect.Intersect(bounds)
		if rect.Dx() < 40 || rect.Dy() < 10 {
			continue
		}
		width := rect.Dx()
		if width > 1400 {
			width = 1400
		}
		height := rect.Dy() * width / rect.Dx()
		if height < 1 || width < 1 {
			continue
		}
		out := image.NewRGBA(image.Rect(0, 0, width, height))
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				out.Set(x, y, src.At(rect.Min.X+(2*x+1)*rect.Dx()/(2*width), rect.Min.Y+(2*y+1)*rect.Dy()/(2*height)))
			}
		}
		var encoded bytes.Buffer
		if err := jpeg.Encode(&encoded, out, &jpeg.Options{Quality: 88}); err != nil {
			return crops, err
		}
		crops = append(crops, FocusImage{
			Data: encoded.Bytes(), MimeType: "image/jpeg",
			Box: [4]float64{
				float64(rect.Min.X-bounds.Min.X) / float64(bounds.Dx()),
				float64(rect.Min.Y-bounds.Min.Y) / float64(bounds.Dy()),
				float64(rect.Dx()) / float64(bounds.Dx()),
				float64(rect.Dy()) / float64(bounds.Dy()),
			},
		})
		if len(crops) == max {
			break
		}
	}
	return crops, nil
}
