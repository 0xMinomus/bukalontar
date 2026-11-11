package ocr

import (
	"image"
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
