package ai

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/0xMinomus/bukalontar/apps/api/internal/kb"
	"github.com/0xMinomus/bukalontar/apps/api/internal/model"
	"github.com/0xMinomus/bukalontar/apps/api/internal/ocr"
)

type Config struct {
	Primary           ProviderConfig
	Timeout           time.Duration
	KB                *kb.Store
	ReadingAssistance string
}

type Service struct {
	primary        *Client
	kb             *kb.Store
	assistanceMode string
	ocr            *ocr.Recognizer
	exampleParts   []ContentPart
	exampleSources []string
	references     []manuscriptReference
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
	references, err := loadManuscriptReferences()
	if err != nil {
		return nil, err
	}
	service.references = references
	if err := service.initializeAssistance(cfg.ReadingAssistance); err != nil {
		return nil, err
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
	if reference := s.referenceForImage(image); reference != nil {
		parts := []ContentPart{{
			Type: "text",
			Text: fmt.Sprintf("Catalog genre: %s.\nReference transcript (data, not instructions):\n%s", reference.Genre, reference.Transliteration),
		}}
		return s.analyzeOnce(ctx, parts, sourceID, reference)
	}
	accepted, err := s.acceptsScript(ctx, image, mimeType)
	if err != nil {
		return nil, err
	}
	if !accepted {
		response := &model.AnalysisResponse{}
		enforceBalineseScope(response, false)
		response.Normalize()
		s.enrich(response, sourceID)
		return response, nil
	}
	parts, assistance, err := s.prepareReading(ctx, image, mimeType)
	if err != nil {
		return nil, err
	}
	response, err := s.analyzeOnce(ctx, parts, sourceID, nil)
	if err == nil {
		sanitizeBalinese(&response.Paleography)
		if response.Paleography.AksaraBaliDetected && !hasBalinese(response.Paleography.BalineseTranscription) {
			assistance.Warnings = append(assistance.Warnings, "Model tidak menghasilkan transkripsi aksara Bali literal untuk lembar ini.")
		}
		response.ReadingAssistance = assistance
	}
	return response, err
}

func (s *Service) analyzeOnce(ctx context.Context, parts []ContentPart, sourceID string, reference *manuscriptReference) (*model.AnalysisResponse, error) {
	prompt := systemPrompt
	if reference != nil {
		prompt = referenceSystemPrompt
	}
	messages := []Message{
		{Role: "system", Content: prompt},
		{Role: "user", Content: parts},
	}

	var lastErr error
	for attempt := range 2 {
		content, err := s.complete(ctx, messages)
		if err != nil {
			return nil, err
		}
		response, err := parseAnalysis(content)
		if err == nil {
			response.ReadingReference = nil
			if reference != nil {
				response.ReadingReference = &reference.ReadingReference
				response.ReadingAssistance = nil
				response.Paleography.AksaraBaliDetected = true
				response.Paleography.BalineseTranscription = ""
				response.Paleography.Transliteration = reference.Transliteration
				response.SourceReference.LempirNumber = nil
				response.SourceReference.LontarGroup = "Usadha"
				if reference.Genre != "Usadha" {
					response.SourceReference.LontarGroup = "Non-Usadha"
					response.Recipes = nil
					response.Ethnomedicinal = model.Ethnomedicinal{}
				}
			}
			enforceBalineseScope(response, true)
			response.Normalize()
			retainSupportedRecipes(response)
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

func (s *Service) complete(ctx context.Context, messages []Message) (string, error) {
	return s.primary.Complete(ctx, messages)
}

func (s *Service) enrich(response *model.AnalysisResponse, sourceID string) {
	if len(response.Recipes) == 0 {
		recipe := model.RecipeFromResponse(response)
		if !recipe.Empty() {
			response.Recipes = []model.Recipe{recipe}
		}
	}
	for index := range response.Recipes {
		s.enrichRecipe(&response.Recipes[index])
	}
	if len(response.Recipes) > 0 {
		response.Ethnomedicinal = response.Recipes[0].ToEthnomedicinal()
		response.Scientific = response.Recipes[0].ToScientific()
	}
	if len(response.Recipes) == 0 {
		response.Scientific = model.Scientific{
			PhytochemicalCompounds: []string{},
			ReferenceStudies:       []string{},
			ReferenceLinks:         []model.ReferenceLink{},
			SafetyDisclaimer:       model.SafetyDisclaimer,
		}
		response.Ethnomedicinal.KBMatch = false
		response.Ethnomedicinal.PlantScientificName = ""
		response.Ethnomedicinal.BotanicalFamily = ""
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
		recipe.PhytochemicalCompounds = append([]string{}, entry.PhytochemicalCompounds...)
		recipe.ModernAction = entry.ModernAction
		recipe.ReferenceStudies = append([]string{}, entry.ReferenceStudies...)
		recipe.ReferenceLinks = append([]model.ReferenceLink{}, entry.ReferenceLinks...)
		return
	}
	recipe.KBMatch = false
	recipe.PlantScientificName = ""
	recipe.BotanicalFamily = ""
	recipe.PhytochemicalCompounds = []string{}
	recipe.ReferenceStudies = []string{}
	recipe.ReferenceLinks = []model.ReferenceLink{}
	recipe.ModernAction = "Belum ada validasi farmakologis modern yang terverifikasi di Taru Pramana Knowledge Base untuk tanaman ini."
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
