// Package readingexamples provides isolated, expert-labelled manuscript glyphs.
package readingexamples

import (
	"bytes"
	"embed"
	"encoding/json"
	"image/gif"
	"image/png"
	"sort"
	"strings"
	"sync"
)

//go:embed assets/*
var assets embed.FS

//go:embed assets/glyph-chart.jpg
var chartJPEG []byte

// Chart returns the labelled contact sheet of all published glyph exemplars.
// It is a reference chart, not target-manuscript text.
func Chart() ([]byte, string) {
	return chartJPEG, "image/jpeg"
}

// Example is a character illustration, not a transcription of the target image.
type Example struct {
	Image       []byte
	MIME        string
	Label       string
	Source      string
	Granularity string
}

type record struct {
	File  string `json:"file"`
	Label string `json:"label"`
	URL   string `json:"url"`
}

var catalog []Example
var once sync.Once

func load() {
	data, err := assets.ReadFile("assets/catalog.json")
	if err != nil {
		panic(err)
	}
	var records []record
	if err := json.Unmarshal(data, &records); err != nil {
		panic(err)
	}
	for _, r := range records {
		data, err := assets.ReadFile("assets/" + r.File)
		if err != nil {
			panic(err)
		}
		img, err := gif.Decode(bytes.NewReader(data))
		if err != nil {
			panic(err)
		}
		var out bytes.Buffer
		if err := png.Encode(&out, img); err != nil {
			panic(err)
		}
		catalog = append(catalog, Example{
			Image: out.Bytes(), MIME: "image/png", Label: r.Label,
			Source:      "Siahaan et al., DeepLontar, Scientific Data 2022, Table 1 (CC BY 4.0): " + r.URL,
			Granularity: "character",
		})
	}
}

// Select returns at most three isolated glyph illustrations. Matching is based
// only on published class names; it is not visual similarity or OCR. When no
// class name matches, three common base-character examples are returned.
// The source article does not identify the leaves shown, so these must not be
// used as examples for evaluations on DeepLontar images: overlap is unknown.
func Select(query string) []Example {
	once.Do(load)
	q := strings.ToLower(query)
	type scored struct{ index, score int }
	var matches []scored
	for i, e := range catalog {
		label := strings.ToLower(e.Label)
		// Short labels must match a whole token, not an arbitrary substring.
		for _, part := range strings.FieldsFunc(q, func(r rune) bool {
			return !(r >= 'a' && r <= 'z')
		}) {
			if part == label {
				matches = append(matches, scored{i, 2})
				break
			}
		}
		if strings.Contains(label, " ") && strings.Contains(q, label) {
			matches = append(matches, scored{i, 3})
		}
	}
	sort.SliceStable(matches, func(i, j int) bool { return matches[i].score > matches[j].score })
	result := make([]Example, 0, 3)
	seen := make(map[int]bool)
	for _, m := range matches {
		if !seen[m.index] {
			result = append(result, catalog[m.index])
			seen[m.index] = true
		}
		if len(result) == 3 {
			return result
		}
	}
	for _, name := range []string{"Ha", "Na", "Ca"} {
		for i, e := range catalog {
			if e.Label == name && !seen[i] {
				result = append(result, e)
				seen[i] = true
				break
			}
		}
		if len(result) == 3 {
			break
		}
	}
	return result
}
