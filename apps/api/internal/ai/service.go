package ai

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/0xMinomus/bukalontar/apps/api/internal/kb"
	"github.com/0xMinomus/bukalontar/apps/api/internal/model"
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

	chosen := chooseResponse(responses)
	for _, response := range responses {
		if response != chosen {
			mergeRecipes(chosen, response)
		}
	}
	s.enrich(chosen, sourceID)
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

func chooseResponse(responses []*model.AnalysisResponse) *model.AnalysisResponse {
	best := responses[0]
	for _, response := range responses[1:] {
		if len(response.Recipes) > len(best.Recipes) ||
			(len(response.Recipes) == len(best.Recipes) &&
				len([]rune(response.Paleography.Transliteration)) > len([]rune(best.Paleography.Transliteration))) {
			best = response
		}
	}
	return best
}

func mergeRecipes(base, other *model.AnalysisResponse) {
	seen := make(map[string]struct{}, len(base.Recipes))
	for _, recipe := range base.Recipes {
		seen[plantKey(recipe.PlantLocalName)] = struct{}{}
	}
	for _, recipe := range other.Recipes {
		key := plantKey(recipe.PlantLocalName)
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		if len(base.Recipes) >= 6 {
			break
		}
		seen[key] = struct{}{}
		base.Recipes = append(base.Recipes, recipe)
	}
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
	if len(response.Recipes) == 0 && strings.TrimSpace(response.Ethnomedicinal.PlantLocalName) != "" {
		response.Recipes = []model.Recipe{model.RecipeFromResponse(response)}
	}
	for index := range response.Recipes {
		s.enrichRecipe(&response.Recipes[index])
	}
	if len(response.Recipes) > 0 {
		response.Ethnomedicinal = response.Recipes[0].ToEthnomedicinal()
		response.Scientific = response.Recipes[0].ToScientific()
	}
	response.Scientific.SafetyDisclaimer = model.SafetyDisclaimer
	if strings.TrimSpace(response.ManuscriptID) == "" {
		response.ManuscriptID = fmt.Sprintf("bukalontar-%s-%d", safeSourceID(sourceID), time.Now().Unix())
	}
}

func (s *Service) enrichRecipe(recipe *model.Recipe) {
	entry, matched := s.kb.Lookup(recipe.PlantLocalName)
	if matched {
		recipe.KBMatch = true
		recipe.PlantScientificName = entry.ScientificName
		recipe.BotanicalFamily = entry.BotanicalFamily
		if entry.EnergyCharacter != "" {
			recipe.EnergyCharacter = entry.EnergyCharacter
		}
		if len(entry.PhytochemicalCompounds) > 0 {
			recipe.PhytochemicalCompounds = append([]string(nil), entry.PhytochemicalCompounds...)
		}
		if entry.ModernAction != "" {
			recipe.ModernAction = entry.ModernAction
		}
		if len(entry.ReferenceStudies) > 0 {
			recipe.ReferenceStudies = append([]string(nil), entry.ReferenceStudies...)
		}
		return
	}
	recipe.KBMatch = false
	recipe.PlantScientificName = ""
	recipe.BotanicalFamily = ""
	recipe.PhytochemicalCompounds = []string{}
	recipe.ReferenceStudies = []string{}
	if strings.TrimSpace(recipe.ModernAction) == "" {
		recipe.ModernAction = "Belum ada validasi farmakologis modern yang terverifikasi di Taru Pramana Knowledge Base untuk tanaman ini."
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
