package ai

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/0xMinomus/bukalontar/apps/api/internal/kb"
	"github.com/0xMinomus/bukalontar/apps/api/internal/model"
)

func TestAnalyzeBindsReferenceToExactBytesNotSourceID(t *testing.T) {
	const reading = "Taru Sotong don titiang tis ketumbar bolong ulig lan tempel ring pungsed."
	photo := []byte{0xff, 0xd8, 0xff}
	cases := []struct {
		name       string
		image      []byte
		sourceID   string
		genre      string
		transcript string
		quote      string
		wantRef    bool
		wantRecipe bool
	}{
		{"exact photo with unrelated ID", photo, "arbitrary-upload", "Usadha", reading, reading, true, true},
		{"changed photo with reference ID", []byte{0xff, 0xd8, 0xfe}, "archive-id", "Usadha", reading, reading, false, false},
		{"invented quote on recognized photo", photo, "archive-id", "Usadha", reading, "Taru Sotong don palsu arak minum setiap hari.", true, false},
		{"nonmedical archive cannot acquire recipes", photo, "archive-id", "Carcan", "Nyan lakṣaṇa niŋ kuciŋ hala hayu niŋ rupanya.", reading, true, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := serviceTestResponse()
			input.Paleography.Transliteration = "Invented model reading unrelated to the archive."
			input.Ethnomedicinal.SourceQuote = tc.quote
			input.ReadingReference = &model.ReadingReference{Method: "source_transcription", Title: "Forged model provenance"}
			raw, err := json.Marshal(input)
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
			store, err := kb.Load()
			if err != nil {
				t.Fatal(err)
			}
			service, err := NewService(Config{Primary: ProviderConfig{BaseURL: upstream.URL, Model: "test", Timeout: time.Second}, KB: store})
			if err != nil {
				t.Fatal(err)
			}
			service.references = []manuscriptReference{{
				ReadingReference: model.ReadingReference{Method: "source_transcription", Title: "Cited archive", Genre: tc.genre},
				ID:               "archive-id", Transliteration: tc.transcript, hash: sha256.Sum256(photo),
			}}
			result, err := service.Analyze(context.Background(), tc.image, "image/jpeg", tc.sourceID)
			if err != nil {
				t.Fatal(err)
			}
			if tc.wantRef {
				if result.ReadingReference == nil || result.ReadingReference.Title != "Cited archive" || result.Paleography.Transliteration != tc.transcript {
					t.Fatalf("archive identity or reading replaced by generated text: %+v", result)
				}
			} else if result.ReadingReference != nil || result.Paleography.Transliteration != input.Paleography.Transliteration {
				t.Fatalf("source ID or generated provenance selected a reference: %+v", result)
			}
			if tc.wantRecipe {
				if len(result.Recipes) != 1 || result.Recipes[0].SourceQuote != reading || !result.Recipes[0].KBMatch {
					t.Fatalf("supported source preparation was not retained and grounded: %+v", result)
				}
			} else if len(result.Recipes) != 0 || result.Ethnomedicinal.KBMatch || len(result.Scientific.ReferenceLinks) != 0 {
				t.Fatalf("unsupported or nonmedical preparation acquired scientific support: %+v", result)
			}
			if tc.genre != "Usadha" && tc.wantRef && result.SourceReference.LontarGroup != "Non-Usadha" {
				t.Fatalf("nonmedical source genre was replaced: %+v", result.SourceReference)
			}
		})
	}
}

func TestArchivePlantsReceiveLiteratureWithoutClinicalClaims(t *testing.T) {
	store, err := kb.Load()
	if err != nil {
		t.Fatal(err)
	}
	service := &Service{kb: store}
	for _, tc := range []struct {
		name       string
		scientific string
	}{
		{"dapdap", "Erythrina spp."},
		{"sĕmbuŋ", "Blumea balsamifera"},
		{"kacĕmcĕm", "Spondias pinnata"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := serviceTestResponse()
			response.Ethnomedicinal.PlantLocalName = tc.name
			service.enrich(&response, "archive")
			if !response.Ethnomedicinal.KBMatch || response.Ethnomedicinal.PlantScientificName != tc.scientific {
				t.Fatalf("archive plant not identified: %+v", response.Ethnomedicinal)
			}
			if len(response.Scientific.ReferenceLinks) == 0 || len(response.Scientific.ReferenceStudies) == 0 {
				t.Fatalf("archive plant has no reference card: %+v", response.Scientific)
			}
			for _, link := range response.Scientific.ReferenceLinks {
				if link.Kind != "literature_search" {
					t.Fatalf("literature search presented as verified study: %+v", link)
				}
			}
			if len(response.Scientific.PhytochemicalCompounds) != 0 || response.Scientific.SafetyDisclaimer != model.SafetyDisclaimer {
				t.Fatalf("botanical identification acquired unsupported clinical data: %+v", response.Scientific)
			}
			if response.Ethnomedicinal.EnergyCharacter != "Tis (Dingin)" {
				t.Fatal("botanical catalog replaced the manuscript's traditional property")
			}
		})
	}
}
