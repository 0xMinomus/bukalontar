package ai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/0xMinomus/bukalontar/apps/api/internal/kb"
	"github.com/0xMinomus/bukalontar/apps/api/internal/model"
)

func TestReadingAssistanceModes(t *testing.T) {
	image, err := os.ReadFile("../../../../data/reading-examples/holdout/aji-kreket-page-4.jpg")
	if err != nil {
		t.Skip("labeled holdout photograph is not present")
	}
	store, err := kb.Load()
	if err != nil {
		t.Fatal(err)
	}
	reading, err := json.Marshal(serviceTestResponse())
	if err != nil {
		t.Fatal(err)
	}
	expectedImages := map[string][2]int{"none": {1, 1}, "examples": {4, 4}, "ocr": {2, 3}, "both": {5, 6}}
	for mode, bounds := range expectedImages {
		var images atomic.Int32
		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			images.Store(int32(strings.Count(string(body), `"type":"image_url"`)))
			_ = json.NewEncoder(w).Encode(map[string]any{
				"choices": []any{map[string]any{"message": map[string]any{"content": string(reading)}}},
			})
		}))
		service, err := NewService(Config{
			Primary:           ProviderConfig{BaseURL: upstream.URL, Model: "test", Timeout: 30 * time.Second},
			KB:                store,
			ReadingAssistance: mode,
		})
		if err != nil {
			t.Fatalf("%s: %v", mode, err)
		}
		result, err := service.Analyze(context.Background(), image, "image/jpeg", "holdout")
		upstream.Close()
		if err != nil {
			t.Fatalf("%s: %v", mode, err)
		}
		if result.ReadingAssistance == nil || result.ReadingAssistance.Mode != mode {
			t.Fatalf("%s: missing assistance metadata", mode)
		}
		got := int(images.Load())
		if got < bounds[0] || got > bounds[1] {
			t.Fatalf("%s: sent %d images, want %d..%d", mode, got, bounds[0], bounds[1])
		}
		if mode == "none" && (result.ReadingAssistance.FocusRegions != 0 || result.ReadingAssistance.ExampleCount != 0) {
			t.Fatalf("none mode still added assistance")
		}
		if mode != "none" && result.ReadingAssistance.ExampleCount == 0 && result.ReadingAssistance.FocusRegions == 0 {
			t.Fatalf("%s mode added no assistance", mode)
		}
		if result.ReadingAssistance.FocusRegions > 2 || result.ReadingAssistance.ExampleCount > 3 {
			t.Fatalf("%s: assistance caps exceeded: %+v", mode, result.ReadingAssistance)
		}
	}
}

func TestSanitizeBalineseDropsLookalikeScripts(t *testing.T) {
	paleography := model.Paleography{BalineseTranscription: "ᬅᬓ᭄ꦲꦥᬤ᭄ᬤ 3 ᬲᬂ"}
	if removed := sanitizeBalinese(&paleography); removed != 2 {
		t.Fatalf("removed %d lookalike characters, want 2", removed)
	}
	if paleography.BalineseTranscription != "ᬅᬓ᭄ᬤ᭄ᬤ 3 ᬲᬂ" {
		t.Fatalf("unexpected cleaned transcription: %q", paleography.BalineseTranscription)
	}
}
