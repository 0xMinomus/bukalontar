package ai

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/0xMinomus/digifest/apps/api/internal/model"
)

const focusPrompt = `You extract only the first medicinal recipe from a Balinese lontar photograph.

Ignore colophons, dates, place names, scribe names, provenance, and ritual instructions. Read the central/main text first.

Return one raw JSON object:
{
  "found": true,
  "plant_local_name": "string",
  "plant_part_used": "Don (Daun)",
  "energy_character": "Panes (Hangat)",
  "formulation_type": "Loloh (Minum)",
  "adjuvants": ["string"],
  "target_ailment": "string",
  "preparation_instructions": "string",
  "evidence": "short transliterated line that names the plant"
}

Rules:
- If no plant is explicitly named in the main text, set found=false and leave strings empty.
- Use the first named plant and the first complete recipe in reading order.
- Allowed plant_part_used: "Don (Daun)", "Babakan (Kulit Kayu)", "Akah (Akar)", "Buah", "Sekar (Bunga)", "Campuran".
- Allowed energy_character: "Panes (Hangat)", "Tis (Dingin)", "Dumelada (Sedang)".
- Allowed formulation_type: "Loloh (Minum)", "Boreh (Balur)", "Tutuh (Tetes)", "Sasap (Kompres)".
- Do not invent a plant, recipe, or citation. Do not use markdown.`

type plantExtraction struct {
	Found                   bool     `json:"found"`
	PlantLocalName          string   `json:"plant_local_name"`
	PlantPartUsed           string   `json:"plant_part_used"`
	EnergyCharacter         string   `json:"energy_character"`
	FormulationType         string   `json:"formulation_type"`
	Adjuvants               []string `json:"adjuvants"`
	TargetAilment           string   `json:"target_ailment"`
	PreparationInstructions string   `json:"preparation_instructions"`
	Evidence                string   `json:"evidence"`
}

func (s *Service) focusPlant(ctx context.Context, image []byte, mimeType string) (*plantExtraction, error) {
	dataURL := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(image)
	messages := []Message{
		{Role: "system", Content: focusPrompt},
		{
			Role: "user",
			Content: []ContentPart{
				{Type: "text", Text: "Find the first medicinal plant recipe in this leaf. Ignore the colophon."},
				{Type: "image_url", ImageURL: &ImageURL{URL: dataURL}},
			},
		},
	}
	content, err := s.primary.CompleteWithTokens(ctx, messages, 1200)
	if err != nil {
		return nil, err
	}
	raw, err := extractJSON(content)
	if err != nil {
		return nil, err
	}
	var extraction plantExtraction
	if err := json.Unmarshal(raw, &extraction); err != nil {
		return nil, fmt.Errorf("decode plant extraction: %w", err)
	}
	return &extraction, nil
}

func (s *Service) mergePlant(response *model.AnalysisResponse, extraction *plantExtraction, sourceID string) {
	if extraction == nil || !extraction.Found || strings.TrimSpace(extraction.PlantLocalName) == "" {
		return
	}
	entity := &response.Ethnomedicinal
	entity.PlantLocalName = strings.TrimSpace(extraction.PlantLocalName)
	if validPlantPart(extraction.PlantPartUsed) {
		entity.PlantPartUsed = extraction.PlantPartUsed
	}
	if validEnergy(extraction.EnergyCharacter) {
		entity.EnergyCharacter = extraction.EnergyCharacter
	}
	if validFormulation(extraction.FormulationType) {
		entity.FormulationType = extraction.FormulationType
	}
	if len(extraction.Adjuvants) > 0 {
		entity.Adjuvants = extraction.Adjuvants
	}
	if strings.TrimSpace(extraction.TargetAilment) != "" {
		entity.TargetAilment = strings.TrimSpace(extraction.TargetAilment)
	}
	if strings.TrimSpace(extraction.PreparationInstructions) != "" {
		entity.PreparationInstructions = strings.TrimSpace(extraction.PreparationInstructions)
	}
	response.Normalize()
	s.enrich(response, sourceID)
}

func validPlantPart(value string) bool {
	switch value {
	case "Don (Daun)", "Babakan (Kulit Kayu)", "Akah (Akar)", "Buah", "Sekar (Bunga)", "Campuran":
		return true
	default:
		return false
	}
}

func validEnergy(value string) bool {
	switch value {
	case "Panes (Hangat)", "Tis (Dingin)", "Dumelada (Sedang)":
		return true
	default:
		return false
	}
}

func validFormulation(value string) bool {
	switch value {
	case "Loloh (Minum)", "Boreh (Balur)", "Tutuh (Tetes)", "Sasap (Kompres)":
		return true
	default:
		return false
	}
}
