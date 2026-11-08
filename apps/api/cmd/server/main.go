package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xMinomus/digifest/apps/api/internal/ai"
	"github.com/0xMinomus/digifest/apps/api/internal/config"
	"github.com/0xMinomus/digifest/apps/api/internal/handler"
	"github.com/0xMinomus/digifest/apps/api/internal/kb"
	"github.com/0xMinomus/digifest/apps/api/internal/ratelimit"
	"github.com/0xMinomus/digifest/apps/api/internal/webui"
)

func main() {
	config.LoadDotEnv()
	cfg := config.Load()

	store, err := kb.Load()
	if err != nil {
		log.Fatalf("load knowledge base: %v", err)
	}

	primary := ai.ProviderConfig{
		Name:    "9router",
		BaseURL: cfg.AIBaseURL,
		Model:   cfg.AIModel,
		APIKey:  cfg.AIAPIKey,
		Timeout: cfg.AITimeout,
	}
	var fallback *ai.ProviderConfig
	if cfg.FallbackEnabled && cfg.FallbackBaseURL != "" && cfg.FallbackModel != "" && cfg.FallbackAPIKey != "" {
		fallback = &ai.ProviderConfig{
			Name:    "command-code",
			BaseURL: cfg.FallbackBaseURL,
			Model:   cfg.FallbackModel,
			APIKey:  cfg.FallbackAPIKey,
			Timeout: cfg.AITimeout,
		}
	}

	analyzer, err := ai.NewService(ai.Config{
		Primary:  primary,
		Fallback: fallback,
		Timeout:  cfg.AITimeout,
		KB:       store,
	})
	if err != nil {
		log.Fatalf("initialize AI service: %v", err)
	}

	files, hasIndex := webui.Dist()
	app := handler.New(handler.Config{
		Analyzer:      analyzer,
		Limiter:       ratelimit.New(cfg.RateLimitPerMinute),
		KBEntries:     store.Count(),
		Timeout:       cfg.AITimeout,
		AllowedOrigin: cfg.AllowedOrigin,
		Static:        webui.NewHandler(files, hasIndex),
	})

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           app,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       25 * time.Second,
		WriteTimeout:      cfg.AITimeout + 20*time.Second,
		IdleTimeout:       120 * time.Second,
	}

	if cfg.AIAPIKey == "" {
		log.Printf("warning: AI_API_KEY kosong; 9router yang memakai auth akan menolak request")
	}
	log.Printf("UsadhaAI API listening on :%s (model=%s, kb=%d entries, frontend_embedded=%t)", cfg.Port, analyzer.Model(), store.Count(), hasIndex)
	if fallback == nil {
		log.Printf("AI fallback disabled")
	} else {
		log.Printf("AI fallback configured: %s", fallback.Model)
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
