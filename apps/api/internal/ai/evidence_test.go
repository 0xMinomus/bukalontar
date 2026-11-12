package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/0xMinomus/bukalontar/apps/api/internal/kb"
)

func TestAnalyzeRequiresTopLeafEvidenceBeforeGrounding(t *testing.T) {
	const supported = "Tityang taru sotong, don titiang tis, ketumbar bolong ulig lan tempel ring pungsed."
	cases := []struct {
		name       string
		reading    string
		quote      string
		plant      string
		extra      string
		wantRecipe bool
	}{
		{name: "supported preparation", reading: supported, quote: supported, plant: "Taru Sotong", wantRecipe: true},
		{name: "missing quotation", reading: supported, plant: "Taru Sotong"},
		{name: "invented quotation", reading: "Iki usada lare, makaten tenger wong rare.", quote: supported, plant: "Taru Sotong"},
		{name: "lower leaf cannot supply preparation", reading: "[Lembar 1]\nIki usada lare.\n[Lembar 2]\n" + supported, quote: supported, plant: "Taru Sotong"},
		{name: "plant not in quoted preparation", reading: supported, quote: supported, plant: "Taru Pule"},
		{name: "mere plant mention", reading: supported, quote: "taru sotong", plant: "Taru Sotong"},
		{name: "unquoted complementary ingredient", reading: supported, quote: supported, plant: "Taru Sotong", extra: "arak"},
		{name: "plant substring is not a plant name", reading: supported, quote: supported, plant: "Otong"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := serviceTestResponse()
			input.Paleography.Transliteration = tc.reading
			input.Ethnomedicinal.SourceQuote = tc.quote
			input.Ethnomedicinal.PlantLocalName = tc.plant
			if tc.extra != "" {
				input.Ethnomedicinal.Adjuvants = append(input.Ethnomedicinal.Adjuvants, tc.extra)
			}
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
			service, err := NewService(Config{
				Primary: ProviderConfig{BaseURL: upstream.URL, Model: "test", Timeout: time.Second},
				KB:      store,
			})
			if err != nil {
				t.Fatal(err)
			}
			result, err := service.Analyze(context.Background(), []byte{0xff, 0xd8, 0xff}, "image/jpeg", "evidence")
			if err != nil {
				t.Fatal(err)
			}
			if tc.wantRecipe {
				if len(result.Recipes) != 1 || !result.Recipes[0].KBMatch || result.Recipes[0].SourceQuote != supported || result.Recipes[0].EnergyCharacter != "Tis (Dingin)" {
					t.Fatalf("supported source was lost or replaced by KB assumptions: %+v", result)
				}
			} else if len(result.Recipes) != 0 || result.Ethnomedicinal.KBMatch || result.Ethnomedicinal.PlantLocalName != "Tidak disebutkan." || len(result.Scientific.ReferenceLinks) != 0 {
				t.Fatalf("unsupported recipe escaped or reappeared during normalization: %+v", result)
			}
		})
	}
}
