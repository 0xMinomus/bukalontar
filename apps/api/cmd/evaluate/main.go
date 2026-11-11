package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/0xMinomus/bukalontar/apps/api/internal/ai"
	"github.com/0xMinomus/bukalontar/apps/api/internal/config"
	"github.com/0xMinomus/bukalontar/apps/api/internal/kb"
	"github.com/0xMinomus/bukalontar/apps/api/internal/model"
)

type result struct {
	Model       string                  `json:"model"`
	Variant     string                  `json:"variant"`
	Image       string                  `json:"image"`
	ImageSHA256 string                  `json:"image_sha256"`
	Reasoning   string                  `json:"reasoning_effort,omitempty"`
	Assistance  string                  `json:"assistance"`
	Repeat      int                     `json:"repeat"`
	StartedAt   time.Time               `json:"started_at"`
	ElapsedMS   int64                   `json:"elapsed_ms"`
	Error       string                  `json:"error,omitempty"`
	Response    *model.AnalysisResponse `json:"response,omitempty"`
}

func main() {
	config.LoadDotEnv()
	cfg := config.Load()
	models := flag.String("models", cfg.AIModel, "Comma-separated exact router model IDs")
	repeats := flag.Int("repeat", 3, "Runs per image and model")
	variant := flag.String("variant", "current", "Pipeline label recorded in the result")
	output := flag.String("out", "", "JSONL result file; empty means stdout")
	reasoning := flag.String("reasoning", cfg.AIReasoningEffort, "Reasoning effort for Codex models")
	assistance := flag.String("assistance", cfg.AIReadingAssistance, "Comma-separated modes: none,examples,ocr,both")
	flag.Parse()
	if flag.NArg() == 0 || *repeats < 1 {
		log.Fatal("usage: evaluate [-models id,id] [-repeat 3] [-out results.jsonl] image.jpg ...")
	}
	store, err := kb.Load()
	if err != nil {
		log.Fatal(err)
	}
	writer := os.Stdout
	if *output != "" {
		writer, err = os.OpenFile(*output, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer writer.Close()
	}
	encoder := json.NewEncoder(writer)
	var modelIDs, modes []string
	var services []*ai.Service
	for _, id := range strings.Split(*models, ",") {
		for _, mode := range strings.Split(*assistance, ",") {
			id, mode = strings.TrimSpace(id), strings.TrimSpace(mode)
			service, err := ai.NewService(ai.Config{
				Primary: ai.ProviderConfig{BaseURL: cfg.AIBaseURL, Model: id, ReasoningEffort: *reasoning, APIKey: cfg.AIAPIKey, Timeout: cfg.AITimeout},
				Timeout: cfg.AITimeout, KB: store, ReadingAssistance: mode,
			})
			if err != nil {
				log.Fatal(err)
			}
			modelIDs, modes = append(modelIDs, id), append(modes, mode)
			services = append(services, service)
		}
	}
	for _, path := range flag.Args() {
		image, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		mimeType := http.DetectContentType(image)
		if mimeType != "image/jpeg" && mimeType != "image/png" && mimeType != "image/webp" {
			log.Fatalf("unsupported image type: %s", path)
		}
		sum := sha256.Sum256(image)
		for repeat := 1; repeat <= *repeats; repeat++ {
			for i, service := range services {
				row := result{Model: modelIDs[i], Variant: *variant + "-" + modes[i], Assistance: modes[i], Image: filepath.ToSlash(path), ImageSHA256: hex.EncodeToString(sum[:]), Repeat: repeat, StartedAt: time.Now().UTC()}
				if strings.HasPrefix(modelIDs[i], "cx/") {
					row.Reasoning = *reasoning
				}
				ctx, cancel := context.WithTimeout(context.Background(), cfg.AITimeout)
				started := time.Now()
				row.Response, err = service.Analyze(ctx, image, mimeType, "evaluation")
				row.ElapsedMS = time.Since(started).Milliseconds()
				cancel()
				if err != nil {
					row.Error = err.Error()
				}
				if err := encoder.Encode(row); err != nil {
					log.Fatal(err)
				}
				fmt.Fprintf(os.Stderr, "%s mode=%s %s run=%d elapsed=%dms success=%t\n", row.Model, row.Assistance, filepath.Base(path), repeat, row.ElapsedMS, row.Error == "")
			}
		}
	}
}
