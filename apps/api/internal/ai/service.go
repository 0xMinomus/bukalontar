package ai

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/0xMinomus/digifest/apps/api/internal/kb"
	"github.com/0xMinomus/digifest/apps/api/internal/model"
)

type Config struct {
	Primary  ProviderConfig
	Fallback *ProviderConfig
	Timeout  time.Duration
	KB       *kb.Store
}

type Service struct {
	primary  *Client
	fallback *Client
	kb       *kb.Store
}

func NewService(cfg Config) (*Service, error) {
	if cfg.KB == nil {
		return nil, fmt.Errorf("knowledge base is required")
	}
	if strings.TrimSpace(cfg.Primary.BaseURL) == "" || strings.TrimSpace(cfg.Primary.Model) == "" {
		return nil, fmt.Errorf("AI_BASE_URL and AI_MODEL are required")
	}
	if cfg.Primary.Timeout <= 0 {
		cfg.Primary.Timeout = cfg.Timeout
	}
	service := &Service{
		primary: NewClient(cfg.Primary),
		kb:      cfg.KB,
	}
	if cfg.Fallback != nil && strings.TrimSpace(cfg.Fallback.BaseURL) != "" && strings.TrimSpace(cfg.Fallback.Model) != "" {
		if cfg.Fallback.Timeout <= 0 {
			cfg.Fallback.Timeout = cfg.Timeout
		}
		service.fallback = NewClient(*cfg.Fallback)
	}
	return service, nil
}

func (s *Service) Model() string {
	return s.primary.Model()
}

func (s *Service) Health(ctx context.Context) error {
	return s.primary.Health(ctx)
}

type analysisResult struct {
	response *model.AnalysisResponse
	err      error
}

func (s *Service) Analyze(ctx context.Context, image []byte, mimeType, sourceID string) (*model.AnalysisResponse, error) {
	dataURL := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(image)

	results := make(chan analysisResult, 2)
	for range 2 {
		go func() {
			response, err := s.analyzeOnce(ctx, dataURL, sourceID)
			results <- analysisResult{response: response, err: err}
		}()
	}

	responses := make([]*model.AnalysisResponse, 0, 2)
	var firstErr error
	for range 2 {
		result := <-results
		if result.err != nil {
			if firstErr == nil {
				firstErr = result.err
			}
			continue
		}
		responses = append(responses, result.response)
	}
	if len(responses) == 0 {
		return nil, firstErr
	}

	chosen, consistent := chooseResponse(responses)
	if !consistent {
		if hasPlantResponse(chosen) || needsContentRetry(chosen) {
			if !s.highReasoning() && canRetryContent(ctx) {
				candidates := candidatePlants(responses)
				if len(candidates) > 0 {
					if focused, focusErr := s.focusPlant(ctx, image, mimeType, candidates); focusErr == nil && focused.Found {
						if match := findResponse(plantKey(focused.PlantLocalName), responses); match != nil {
							chosen = match
							consistent = true
						}
					}
				}
			}
			if !consistent {
				return nil, ErrUnstable
			}
		} else {
			return chosen, nil
		}
	}
	return chosen, nil
}

func (s *Service) analyzeOnce(ctx context.Context, dataURL, sourceID string) (*model.AnalysisResponse, error) {
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{
			Role: "user",
			Content: []ContentPart{
				{Type: "text", Text: userPrompt(sourceID, s.kb)},
				{Type: "image_url", ImageURL: &ImageURL{URL: dataURL}},
			},
		},
	}

	var lastErr error
	for attempt := range 2 {
		content, err := s.complete(ctx, messages)
		if err != nil {
			return nil, err
		}
		response, err := parseAnalysis(content)
		if err == nil {
			response.Normalize()
			s.enrich(response, sourceID)
			err = response.Validate()
		}
		if err == nil {
			return response, nil
		}
		lastErr = err
		if attempt == 0 {
			messages = append(messages,
				Message{Role: "assistant", Content: content},
				Message{Role: "user", Content: correctionPrompt(err)},
			)
		}
	}
	return nil, fmt.Errorf("%w: %v", ErrInvalidOutput, lastErr)
}

func (s *Service) highReasoning() bool {
	return strings.Contains(strings.ToLower(s.primary.Model()), "high")
}

func chooseResponse(responses []*model.AnalysisResponse) (*model.AnalysisResponse, bool) {
	best := responses[0]
	for _, response := range responses[1:] {
		if len([]rune(response.Paleography.Transliteration)) > len([]rune(best.Paleography.Transliteration)) {
			best = response
		}
	}

	if !hasPlantResponse(best) {
		for _, response := range responses {
			if hasPlantResponse(response) {
				best = response
				break
			}
		}
	}

	consistent := true
	firstKey := ""
	for _, response := range responses {
		key := plantKey(response.Ethnomedicinal.PlantLocalName)
		if key == "" {
			consistent = false
			continue
		}
		if firstKey == "" {
			firstKey = key
		} else if key != firstKey {
			consistent = false
		}
	}
	consistent = consistent && firstKey != ""
	return best, consistent
}

func findResponse(key string, responses []*model.AnalysisResponse) *model.AnalysisResponse {
	if key == "" {
		return nil
	}
	for _, response := range responses {
		if plantKey(response.Ethnomedicinal.PlantLocalName) == key {
			return response
		}
	}
	return nil
}

func hasPlantResponse(response *model.AnalysisResponse) bool {
	return plantKey(response.Ethnomedicinal.PlantLocalName) != ""
}

func plantKey(name string) string {
	value := strings.ToLower(strings.TrimSpace(name))
	value = strings.TrimPrefix(value, "taru ")
	value = strings.Map(func(char rune) rune {
		if unicode.IsLetter(char) || unicode.IsDigit(char) || unicode.IsSpace(char) {
			return char
		}
		return -1
	}, value)
	return strings.Join(strings.Fields(value), " ")
}

func canRetryContent(ctx context.Context) bool {
	deadline, ok := ctx.Deadline()
	if !ok {
		return true
	}
	return time.Until(deadline) > 45*time.Second
}

func needsContentRetry(response *model.AnalysisResponse) bool {
	plant := strings.ToLower(strings.TrimSpace(response.Ethnomedicinal.PlantLocalName))
	if plant != "" && !strings.HasPrefix(plant, "tidak disebutkan") && !strings.HasPrefix(plant, "tidak ada") {
		return false
	}
	summary := strings.ToLower(response.Ethnomedicinal.PreparationInstructions + " " + response.Paleography.IndonesianTranslation)
	if strings.Contains(summary, "kolofon") || strings.Contains(summary, "colophon") {
		return true
	}
	transliteration := []rune(strings.TrimSpace(response.Paleography.Transliteration))
	return len(transliteration) < 260
}

func (s *Service) complete(ctx context.Context, messages []Message) (string, error) {
	content, err := s.primary.Complete(ctx, messages)
	if err == nil {
		return content, nil
	}
	if s.fallback == nil || !shouldFallback(err) {
		return "", err
	}
	fallbackContent, fallbackErr := s.fallback.Complete(ctx, messages)
	if fallbackErr == nil {
		return fallbackContent, nil
	}
	return "", fmt.Errorf("primary provider: %v; fallback provider: %w", err, fallbackErr)
}

func (s *Service) enrich(response *model.AnalysisResponse, sourceID string) {
	response.Scientific.SafetyDisclaimer = model.SafetyDisclaimer
	entry, matched := s.kb.Lookup(response.Ethnomedicinal.PlantLocalName)
	if matched {
		response.Ethnomedicinal.KBMatch = true
		response.Ethnomedicinal.PlantScientificName = entry.ScientificName
		response.Ethnomedicinal.BotanicalFamily = entry.BotanicalFamily
		if entry.EnergyCharacter != "" {
			response.Ethnomedicinal.EnergyCharacter = entry.EnergyCharacter
		}
		if len(entry.PhytochemicalCompounds) > 0 {
			response.Scientific.PhytochemicalCompounds = append([]string(nil), entry.PhytochemicalCompounds...)
		}
		if entry.ModernAction != "" {
			response.Scientific.ModernAction = entry.ModernAction
		}
		if len(entry.ReferenceStudies) > 0 {
			response.Scientific.ReferenceStudies = append([]string(nil), entry.ReferenceStudies...)
		}
	} else {
		response.Ethnomedicinal.KBMatch = false
		response.Ethnomedicinal.PlantScientificName = ""
		response.Ethnomedicinal.BotanicalFamily = ""
		response.Scientific.PhytochemicalCompounds = []string{}
		response.Scientific.ReferenceStudies = []string{}
		if strings.TrimSpace(response.Scientific.ModernAction) == "" {
			response.Scientific.ModernAction = "Belum ada validasi farmakologis modern yang terverifikasi di Taru Pramana Knowledge Base untuk tanaman ini."
		}
	}

	if strings.TrimSpace(response.ManuscriptID) == "" {
		response.ManuscriptID = fmt.Sprintf("usadha-%s-%d", safeSourceID(sourceID), time.Now().Unix())
	}
}

func safeSourceID(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.Map(func(char rune) rune {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			return char
		}
		return -1
	}, value)
	if value == "" {
		return "manuscript"
	}
	if len(value) > 32 {
		value = value[:32]
	}
	return value
}
