package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/0xMinomus/bukalontar/apps/api/internal/kb"
	"github.com/0xMinomus/bukalontar/apps/api/internal/model"
)

func serviceTestResponse() model.AnalysisResponse {
	lempir := 3
	return model.AnalysisResponse{
		ManuscriptID: "manuscript-1",
		SourceReference: model.SourceReference{
			LontarGroup:  "Usadha Taru Pramana",
			LempirNumber: &lempir,
			Language:     "Kawi / Jawa Kuno",
		},
		Paleography: model.Paleography{
			AksaraBaliDetected:    true,
			Confidence:            "high",
			Transliteration:       "Tityang taru sotong, don titiang tis, ketumbar bolong ulig lan tempel ring pungsed.",
			IndonesianTranslation: "Hamba pohon jambu biji",
			EnglishTranslation:    "I am the guava tree",
		},
		Ethnomedicinal: model.Ethnomedicinal{
			PlantLocalName:          "Taru Sotong",
			SourceQuote:             "Tityang taru sotong, don titiang tis, ketumbar bolong ulig lan tempel ring pungsed.",
			PlantScientificName:     "Model guess",
			BotanicalFamily:         "Wrong family",
			KBMatch:                 false,
			PlantPartUsed:           "Don (Daun)",
			EnergyCharacter:         "Tis (Dingin)",
			FormulationType:         "Sasap (Kompres)",
			Adjuvants:               []string{"Ketumbar bolong"},
			TargetAilment:           "Mising (diare)",
			PreparationInstructions: "Pucuk daun ditumbuk lalu ditempelkan pada pusar.",
		},
		Scientific: model.Scientific{
			PhytochemicalCompounds: []string{"Model compound"},
			ModernAction:           "Model action",
			ReferenceStudies:       []string{"Model reference"},
			SafetyDisclaimer:       "Model disclaimer",
		},
	}
}

func TestAnalyzeRepairsInvalidOutputAndAppliesKB(t *testing.T) {
	store, err := kb.Load()
	if err != nil {
		t.Fatal(err)
	}
	valid, err := json.Marshal(serviceTestResponse())
	if err != nil {
		t.Fatal(err)
	}

	var calls atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if serveAcceptedScript(w, r) {
			return
		}
		call := calls.Add(1)
		content := "not json"
		if call > 1 {
			content = string(valid)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []any{
				map[string]any{"message": map[string]any{"content": content}},
			},
		})
	}))
	defer upstream.Close()

	service, err := NewService(Config{
		Primary: ProviderConfig{
			BaseURL: upstream.URL,
			Model:   "test",
			Timeout: 5 * time.Second,
		},
		Timeout: 5 * time.Second,
		KB:      store,
	})
	if err != nil {
		t.Fatal(err)
	}

	response, err := service.Analyze(context.Background(), []byte{0xff, 0xd8, 0xff}, "image/jpeg", "kemenuh")
	if err != nil {
		t.Fatal(err)
	}
	if !response.Ethnomedicinal.KBMatch {
		t.Fatal("KB match should be true")
	}
	if response.Ethnomedicinal.PlantScientificName != "Psidium guajava" {
		t.Fatalf("scientific name not grounded: %q", response.Ethnomedicinal.PlantScientificName)
	}
	if response.Ethnomedicinal.BotanicalFamily != "Myrtaceae" {
		t.Fatalf("family not grounded: %q", response.Ethnomedicinal.BotanicalFamily)
	}
	if response.Ethnomedicinal.EnergyCharacter != "Tis (Dingin)" {
		t.Fatalf("KB must not replace the source's traditional property: %q", response.Ethnomedicinal.EnergyCharacter)
	}
	if response.Scientific.SafetyDisclaimer != model.SafetyDisclaimer {
		t.Fatal("safety disclaimer not canonical")
	}
	if response.ManuscriptID != "manuscript-1" {
		t.Fatalf("model manuscript id should be preserved: %q", response.ManuscriptID)
	}
}

func TestEnrichClearsUnknownTaxonomy(t *testing.T) {
	store, err := kb.Load()
	if err != nil {
		t.Fatal(err)
	}
	service := &Service{kb: store}
	response := serviceTestResponse()
	response.Ethnomedicinal.PlantLocalName = "Taru Khayalan"
	response.Ethnomedicinal.KBMatch = true
	response.Scientific.ReferenceLinks = []model.ReferenceLink{{Title: "Invented", URL: "https://example.org/fake", Kind: "literature"}}
	response.Scientific.SafetyDisclaimer = "wrong"
	service.enrich(&response, "upload")

	if response.Ethnomedicinal.KBMatch {
		t.Fatal("unknown plant must not claim KB match")
	}
	if response.Ethnomedicinal.PlantScientificName != "" || response.Ethnomedicinal.BotanicalFamily != "" {
		t.Fatal("unknown plant taxonomy must be cleared")
	}
	if len(response.Scientific.PhytochemicalCompounds) != 0 || len(response.Scientific.ReferenceStudies) != 0 {
		t.Fatal("unknown plant citations must be cleared")
	}
	if len(response.Scientific.ReferenceLinks) != 0 || response.Scientific.ModernAction == "Model action" {
		t.Fatal("unknown plant must not retain model-supplied links or efficacy claims")
	}
	if response.Scientific.SafetyDisclaimer != model.SafetyDisclaimer {
		t.Fatal("safety disclaimer must be canonical")
	}
}

func TestAnalyzeRejectsBadOutputAfterRepair(t *testing.T) {
	store, err := kb.Load()
	if err != nil {
		t.Fatal(err)
	}
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if serveAcceptedScript(w, r) {
			return
		}
		fmt.Fprint(w, `{"choices":[{"message":{"content":"not json"}}]}`)
	}))
	defer upstream.Close()

	service, err := NewService(Config{
		Primary: ProviderConfig{BaseURL: upstream.URL, Model: "test", Timeout: 5 * time.Second},
		Timeout: 5 * time.Second,
		KB:      store,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Analyze(context.Background(), []byte{0xff, 0xd8, 0xff}, "image/jpeg", "upload"); err == nil {
		t.Fatal("expected invalid output error")
	}
}

func TestAnalyzeReplacesModelCitationsWithCuratedSources(t *testing.T) {
	store, err := kb.Load()
	if err != nil {
		t.Fatal(err)
	}
	reading := serviceTestResponse()
	reading.Scientific.ReferenceLinks = []model.ReferenceLink{{Title: "Fake study", URL: "javascript:alert(1)", Kind: "literature"}}
	raw, err := json.Marshal(reading)
	if err != nil {
		t.Fatal(err)
	}
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if serveAcceptedScript(w, r) {
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []any{map[string]any{"message": map[string]any{"content": string(raw)}}},
		})
	}))
	defer upstream.Close()
	service, err := NewService(Config{
		Primary: ProviderConfig{BaseURL: upstream.URL, Model: "test", Timeout: time.Second},
		KB:      store,
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err := service.Analyze(context.Background(), []byte{0xff, 0xd8, 0xff}, "image/jpeg", "upload")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Recipes) != 1 {
		t.Fatalf("expected one grounded medicinal preparation, got %d", len(result.Recipes))
	}
	for _, validation := range []model.Scientific{result.Scientific, result.Recipes[0].ToScientific()} {
		foundReview := false
		for _, link := range validation.ReferenceLinks {
			if link.URL == "https://pubmed.ncbi.nlm.nih.gov/18353572/" {
				foundReview = true
			}
			if link.URL == "javascript:alert(1)" {
				t.Fatal("untrusted model URL escaped grounding")
			}
		}
		if !foundReview {
			t.Fatal("curated guava review missing from result")
		}
	}
}

func TestAcceptedScriptWithholdsUnreadableOrConflictingOutput(t *testing.T) {
	for _, script := range []string{"rejected", "ꦲꦤ", "𑼒𑼀", "ᬅꦲ"} {
		t.Run(script, func(t *testing.T) {
			input := serviceTestResponse()
			input.Paleography.AksaraBaliDetected = script != "rejected"
			input.Paleography.BalineseTranscription = script
			content, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if serveAcceptedScript(w, r) {
					return
				}
				json.NewEncoder(w).Encode(map[string]any{
					"choices": []any{map[string]any{"message": map[string]any{"content": string(content)}}},
				})
			}))
			defer upstream.Close()
			store, err := kb.Load()
			if err != nil {
				t.Fatal(err)
			}
			service, err := NewService(Config{
				Primary: ProviderConfig{BaseURL: upstream.URL, Model: "test", Timeout: time.Second},
				KB:      store,
			})
			if err != nil {
				t.Fatal(err)
			}
			result, err := service.Analyze(context.Background(), []byte{0xff, 0xd8, 0xff}, "image/jpeg", "scope")
			if err != nil {
				t.Fatal(err)
			}
			if !result.Paleography.AksaraBaliDetected || result.Paleography.Confidence != "low" ||
				result.Paleography.BalineseTranscription != "" ||
				result.Paleography.Transliteration != "Tidak terbaca." ||
				result.Paleography.IndonesianTranslation != "Tidak terbaca." ||
				result.Paleography.EnglishTranslation != "Tidak terbaca." {
				t.Fatalf("rejected target leaked a reading: %+v", result.Paleography)
			}
			if len(result.Recipes) != 0 || result.Ethnomedicinal.KBMatch ||
				result.Ethnomedicinal.PlantLocalName != "Tidak disebutkan." ||
				len(result.Scientific.ReferenceLinks) != 0 ||
				result.SourceReference.Language != "Tidak diketahui" || result.SourceReference.LempirNumber != nil {
				t.Fatalf("rejected target leaked findings: %+v", result)
			}
		})
	}
}

func serveAcceptedScript(w http.ResponseWriter, r *http.Request) bool {
	var request struct {
		MaxTokens int `json:"max_tokens"`
	}
	_ = json.NewDecoder(r.Body).Decode(&request)
	if request.MaxTokens != 256 {
		return false
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"choices": []any{map[string]any{"message": map[string]any{
			"content": "```json\n{\"script\":\"balinese\",\"confidence\":\"high\"}\n```",
		}}},
	})
	return true
}

func TestScriptGateNeverReadsRejectedOrInvalidDecisions(t *testing.T) {
	for _, decision := range []string{
		`{"script":"javanese","confidence":"high"}`,
		`{"script":"kawi","confidence":"high"}`,
		"```json\n{\"script\":\"kawi\",\"confidence\":\"high\"}\n```",
		`{"script":"mixed","confidence":"high"}`,
		`{"script":"uncertain","confidence":"low"}`,
		`{"script":"balinese","confidence":"low"}`,
		`{"script":"unknown","confidence":"high"}`,
		`{}`,
		`not json`,
	} {
		t.Run(decision, func(t *testing.T) {
			var calls atomic.Int32
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"choices": []any{map[string]any{"message": map[string]any{"content": decision}}},
				})
			}))
			defer upstream.Close()
			store, err := kb.Load()
			if err != nil {
				t.Fatal(err)
			}
			service, err := NewService(Config{
				Primary: ProviderConfig{BaseURL: upstream.URL, Model: "test", Timeout: time.Second},
				KB:      store,
			})
			if err != nil {
				t.Fatal(err)
			}
			result, err := service.Analyze(context.Background(), []byte{0xff, 0xd8, 0xff}, "image/jpeg", "scope")
			if calls.Load() != 1 {
				t.Fatal("rejected or invalid script decision reached the manuscript reader")
			}
			if err == nil && (result.Paleography.AksaraBaliDetected || len(result.Recipes) != 0 ||
				result.Paleography.Transliteration != "Tidak terdeteksi.") {
				t.Fatalf("rejected script produced a reading: %+v", result)
			}
			valid := strings.Contains(decision, `"script":`) && !strings.Contains(decision, `"unknown"`)
			if valid && err != nil {
				t.Fatalf("valid refusal returned an error: %v", err)
			}
			if !valid && err == nil {
				t.Fatal("invalid classifier output must fail closed")
			}
		})
	}
}
