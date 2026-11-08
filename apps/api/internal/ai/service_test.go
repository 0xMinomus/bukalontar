package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/0xMinomus/digifest/apps/api/internal/kb"
	"github.com/0xMinomus/digifest/apps/api/internal/model"
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
			Transliteration:       "Tityang taru sotong",
			IndonesianTranslation: "Hamba pohon jambu biji",
			EnglishTranslation:    "I am the guava tree",
		},
		Ethnomedicinal: model.Ethnomedicinal{
			PlantLocalName:          "Taru Sotong",
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
	if response.Ethnomedicinal.EnergyCharacter != "Panes (Hangat)" {
		t.Fatalf("energy not grounded: %q", response.Ethnomedicinal.EnergyCharacter)
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
