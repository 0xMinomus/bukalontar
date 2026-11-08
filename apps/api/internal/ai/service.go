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

func (s *Service) Analyze(ctx context.Context, image []byte, mimeType, sourceID string) (*model.AnalysisResponse, error) {
	dataURL := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(image)
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
			if validationErr := response.Validate(); validationErr == nil {
				if attempt == 0 && needsContentRetry(response) && canRetryContent(ctx) {
					err = fmt.Errorf("previous reading appears to cover only a colophon or omit the central text; read the main text of every leaf and extract any named plant and preparation")
				} else {
					return response, nil
				}
			} else {
				err = validationErr
			}
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

func canRetryContent(ctx context.Context) bool {
	deadline, ok := ctx.Deadline()
	if !ok {
		return true
	}
	return time.Until(deadline) > 90*time.Second
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
