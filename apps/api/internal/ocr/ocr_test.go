package ocr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestBlankAndInvalidImages(t *testing.T) {
	r, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	im := image.NewRGBA(image.Rect(0, 0, 150, 30))
	for y := 0; y < 30; y++ {
		for x := 0; x < 150; x++ {
			im.Set(x, y, color.White)
		}
	}
	var data bytes.Buffer
	if err := jpeg.Encode(&data, im, nil); err != nil {
		t.Fatal(err)
	}
	ev, err := r.Recognize(context.Background(), data.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(ev.FocusImages) != 0 {
		t.Fatal("blank image produced focus regions")
	}
	if _, err := r.Recognize(context.Background(), []byte("not an image")); err == nil {
		t.Fatal("invalid image accepted")
	}
	if _, err := r.ClassifyExpertBox(context.Background(), im); err == nil {
		t.Fatal("flat glyph crop accepted")
	}
}

func TestGlyphModelPythonParity(t *testing.T) {
	data, err := os.ReadFile("../../../../data/ocr/hog-parity.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Pixels  []uint8 `json:"pixels"`
		ClassID int     `json:"class_id"`
		Margin  float64 `json:"margin"`
	}
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	r, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for i, c := range cases {
		im := image.NewGray(image.Rect(0, 0, 16, 16))
		copy(im.Pix, c.Pixels)
		got, err := r.ClassifyExpertBox(context.Background(), im)
		if err != nil || got.ClassID != c.ClassID || math.Abs(got.Margin-c.Margin) > 0.0001 {
			t.Fatalf("case %d: got %+v %v, expected %d %.6f", i, got, err, c.ClassID, c.Margin)
		}
	}
}

func TestPhotoFocus(t *testing.T) {
	output := os.Getenv("OCR_PHOTO_OUTPUT")
	if output == "" {
		t.Skip("set OCR_PHOTO_OUTPUT for local photo inspection")
	}
	r, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for i, path := range []string{"../../../../data/reading-examples/holdout/aji-kreket-page-4.jpg", "../../../../Lontar Test/wanagiri-usada-3.jpg", "../../../../Lontar Test/sukawati-usada-07.jpg"} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		ev, err := r.Recognize(context.Background(), data)
		if err != nil {
			t.Fatal(err)
		}
		if len(ev.FocusImages) == 0 {
			t.Fatalf("no crops for %s: %v", path, ev.Warnings)
		}
		for n, f := range ev.FocusImages {
			name := filepath.Join(output, fmt.Sprintf("photo-%d-%d.jpg", i, n))
			if err := os.WriteFile(name, f.Data, 0600); err != nil {
				t.Fatal(err)
			}
			t.Logf("%s box=%v", name, f.Box)
		}
	}
}
