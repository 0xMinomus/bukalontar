package ai

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/0xMinomus/digifest/apps/api/internal/model"
)

const focusPrompt = `You are selecting the primary plant from a Balinese lontar photograph.

Candidates: %s

Read the topmost leaf only. Find the first complete medicinal recipe in the central/main text. If one of the candidates is the main plant of that recipe, return found=true and the exact candidate name. If none of the candidates is named in the first recipe, return found=false. Ignore colophons, dates, place names, scribe names, provenance, and lower leaves.

Return one raw JSON object:
{
  "found": true,
  "plant_local_name": "exact candidate name",
  "evidence": "short transliterated line naming the plant"
}

Do not invent a plant outside the candidate list. Do not use markdown.`

type plantExtraction struct {
	Found          bool   `json:"found"`
	PlantLocalName string `json:"plant_local_name"`
	Evidence       string `json:"evidence"`
}

func (s *Service) focusPlant(ctx context.Context, image []byte, mimeType string, candidates []string) (*plantExtraction, error) {
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no plant candidates")
	}
	dataURL := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(image)
	messages := []Message{
		{Role: "system", Content: fmt.Sprintf(focusPrompt, strings.Join(candidates, ", "))},
		{
			Role: "user",
			Content: []ContentPart{
				{Type: "text", Text: "Which candidate is the primary plant of the first recipe on the topmost leaf?"},
				{Type: "image_url", ImageURL: &ImageURL{URL: dataURL}},
			},
		},
	}
	content, err := s.primary.CompleteWithTokens(ctx, messages, 500)
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

func candidatePlants(responses []*model.AnalysisResponse) []string {
	seen := make(map[string]struct{})
	candidates := make([]string, 0, len(responses))
	for _, response := range responses {
		name := strings.TrimSpace(response.Ethnomedicinal.PlantLocalName)
		key := plantKey(name)
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		candidates = append(candidates, name)
	}
	return candidates
}
