package ai

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/0xMinomus/bukalontar/apps/api/internal/model"
	"github.com/0xMinomus/bukalontar/apps/api/internal/ocr"
	"github.com/0xMinomus/bukalontar/apps/api/internal/readingexamples"
)

func (s *Service) initializeAssistance(mode string) error {
	if mode == "" {
		mode = "both"
	}
	switch mode {
	case "none", "examples", "ocr", "both", "leaves", "all", "chart":
	default:
		return fmt.Errorf("invalid reading assistance mode %q", mode)
	}
	s.assistanceMode = mode
	if mode == "ocr" || mode == "both" || mode == "leaves" || mode == "all" {
		var err error
		s.ocr, err = ocr.Load()
		if err != nil {
			return fmt.Errorf("load trained reading detector: %w", err)
		}
	}
	if mode == "chart" {
		data, mimeType := readingexamples.Chart()
		s.exampleParts = append(s.exampleParts,
			ContentPart{Type: "text", Text: "REFERENCE GLYPH CHART: one labelled chart of the 55 published Balinese glyph classes (DeepLontar Table 1, CC BY 4.0). It is a legend, not text from the target manuscript."},
			ContentPart{Type: "image_url", ImageURL: &ImageURL{URL: imageDataURL(data, mimeType)}},
		)
		s.exampleSources = append(s.exampleSources, "Siahaan et al., DeepLontar, Scientific Data 2022, Table 1 (CC BY 4.0); chart assembled by BukaLontar")
	}
	if mode == "examples" || mode == "both" || mode == "all" {
		for _, example := range readingexamples.Select("") {
			s.exampleParts = append(s.exampleParts,
				ContentPart{Type: "text", Text: fmt.Sprintf("REFERENCE GLYPH EXAMPLE (not target text): %s. Granularity: %s. Source: %s", example.Label, example.Granularity, example.Source)},
				ContentPart{Type: "image_url", ImageURL: &ImageURL{URL: imageDataURL(example.Image, example.MIME)}},
			)
			s.exampleSources = append(s.exampleSources, example.Source)
		}
	}
	return nil
}

func (s *Service) prepareReading(ctx context.Context, image []byte, mimeType string) ([]ContentPart, *model.ReadingAssistance, error) {
	metadata := &model.ReadingAssistance{
		Mode: s.assistanceMode, ExampleCount: len(s.exampleSources),
		ExampleSources: append([]string{}, s.exampleSources...), Warnings: []string{},
	}
	parts := make([]ContentPart, 0, 2+len(s.exampleParts)+4)
	parts = append(parts, ContentPart{Type: "text", Text: userPrompt(s.kb)})
	parts = append(parts, s.exampleParts...)
	parts = append(parts,
		ContentPart{Type: "text", Text: "TARGET MANUSCRIPT: read this photograph. Reference glyph examples above are NOT text from this manuscript."},
		ContentPart{Type: "image_url", ImageURL: &ImageURL{URL: imageDataURL(image, mimeType)}},
	)
	if s.ocr != nil && (s.assistanceMode == "leaves" || s.assistanceMode == "all") {
		leaves, err := s.ocr.LeafCrops(ctx, image, 2)
		if ctx.Err() != nil {
			return nil, nil, ctx.Err()
		}
		if err != nil {
			metadata.Warnings = append(metadata.Warnings, "Pemisahan daun tidak tersedia; model membaca foto asli.")
		}
		for i, leaf := range leaves {
			parts = append(parts,
				ContentPart{Type: "text", Text: fmt.Sprintf("LEAF %d OF THE SAME TARGET: full-leaf enlargement in reading order (normalized x,y,width,height %.3f,%.3f,%.3f,%.3f). Use it for small strokes. It is not an extra leaf; do not duplicate its content or recipes.", i+1, leaf.Box[0], leaf.Box[1], leaf.Box[2], leaf.Box[3])},
				ContentPart{Type: "image_url", ImageURL: &ImageURL{URL: imageDataURL(leaf.Data, leaf.MimeType)}},
			)
			metadata.LeafRegions++
		}
	}
	if s.ocr != nil && (s.assistanceMode == "ocr" || s.assistanceMode == "both" || s.assistanceMode == "all") {
		evidence, err := s.ocr.Recognize(ctx, image)
		if ctx.Err() != nil {
			return nil, nil, ctx.Err()
		}
		metadata.OCRModel = evidence.ModelVersion
		metadata.OCRMS = evidence.ElapsedMS
		metadata.Warnings = append(metadata.Warnings, evidence.Warnings...)
		if err != nil {
			metadata.Warnings = append(metadata.Warnings, "Pemilihan fokus lokal tidak tersedia untuk gambar ini; model membaca foto asli.")
		} else {
			for i, focus := range evidence.FocusImages {
				if i >= 2 {
					break
				}
				parts = append(parts,
					ContentPart{Type: "text", Text: fmt.Sprintf("TARGET FOCUS %d: enlarged region of the SAME target photograph (normalized x,y,width,height %.3f,%.3f,%.3f,%.3f). It is not an additional leaf. Verify against the original; do not repeat the text. The detector does not establish glyph identity or a correct reading.", i+1, focus.Box[0], focus.Box[1], focus.Box[2], focus.Box[3])},
					ContentPart{Type: "image_url", ImageURL: &ImageURL{URL: imageDataURL(focus.Data, focus.MimeType)}},
				)
				metadata.FocusRegions++
			}
		}
	}
	return parts, metadata, nil
}

func imageDataURL(data []byte, mimeType string) string {
	return "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(data)
}
