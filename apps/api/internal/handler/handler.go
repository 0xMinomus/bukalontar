package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/0xMinomus/bukalontar/apps/api/internal/ai"
	"github.com/0xMinomus/bukalontar/apps/api/internal/model"
	"github.com/0xMinomus/bukalontar/apps/api/internal/ratelimit"
)

const maxUploadBytes = 8 << 20

type Analyzer interface {
	Analyze(ctx context.Context, image []byte, mimeType, sourceID string) (*model.AnalysisResponse, error)
	Model() string
	Health(ctx context.Context) error
}

type Config struct {
	Analyzer      Analyzer
	Limiter       *ratelimit.Limiter
	KBEntries     int
	Timeout       time.Duration
	AllowedOrigin string
	Static        http.Handler
}

type handler struct {
	config Config
}

func New(config Config) http.Handler {
	h := &handler{config: config}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/analyze-manuscript", h.analyze)
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /api/health", h.health)
	static := config.Static
	if static == nil {
		static = http.NotFoundHandler()
	}
	mux.Handle("/", static)
	return h.withCORS(mux)
}

func (h *handler) analyze(w http.ResponseWriter, r *http.Request) {
	if h.config.Limiter != nil && !h.config.Limiter.Allow(clientIP(r)) {
		w.Header().Set("Retry-After", "60")
		writeError(w, http.StatusTooManyRequests, "Batas permintaan tercapai. Tunggu sebentar lalu coba lagi.")
		return
	}

	image, sourceID, err := readMultipart(w, r)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			writeError(w, http.StatusRequestEntityTooLarge, "Ukuran unggahan melebihi 8 MB.")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(image) == 0 {
		writeError(w, http.StatusBadRequest, "Gambar lontar wajib dikirim pada field image.")
		return
	}
	if len(image) > maxUploadBytes {
		writeError(w, http.StatusRequestEntityTooLarge, "Ukuran gambar melebihi 8 MB.")
		return
	}
	mimeType, err := detectImageMime(image)
	if err != nil {
		writeError(w, http.StatusUnsupportedMediaType, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.config.Timeout)
	defer cancel()

	result, err := h.config.Analyzer.Analyze(ctx, image, mimeType, sourceID)
	if err != nil {
		var upstream *ai.UpstreamError
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			log.Printf("analyze: timeout")
			writeError(w, http.StatusGatewayTimeout, "Analisis melebihi batas waktu. Coba lagi dengan foto yang lebih kecil.")
		case errors.Is(err, ai.ErrInvalidOutput):
			log.Printf("analyze: invalid output: %v", err)
			writeError(w, http.StatusBadGateway, "AI mengembalikan format yang tidak valid. Coba foto lain atau ulangi.")
		case errors.As(err, &upstream):
			log.Printf("analyze: upstream %d: %s", upstream.Status, upstream.Body)
			writeError(w, http.StatusBadGateway, "Layanan AI gagal merespons. Coba lagi beberapa saat.")
		default:
			log.Printf("analyze: %v", err)
			writeError(w, http.StatusBadGateway, "Analisis AI gagal. Coba lagi beberapa saat.")
		}
		return
	}
	if !result.Paleography.AksaraBaliDetected {
		writeError(w, http.StatusUnprocessableEntity, "Foto belum dapat dipastikan beraksara Bali. BukaLontar hanya menerima aksara Bali, bukan aksara Jawa/Kawi atau campuran. Gunakan foto yang jelas.")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("analyze: encode response: %v", err)
	}
}

func (h *handler) health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	aiErr := h.config.Analyzer.Health(ctx)
	payload := map[string]any{
		"ok":         aiErr == nil,
		"model":      h.config.Analyzer.Model(),
		"kb_entries": h.config.KBEntries,
		"checked_at": time.Now().UTC().Format(time.RFC3339),
	}
	status := http.StatusOK
	if aiErr != nil {
		status = http.StatusServiceUnavailable
		payload["error"] = truncate(aiErr.Error(), 200)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func readMultipart(w http.ResponseWriter, r *http.Request) ([]byte, string, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes+512*1024)
	reader, err := r.MultipartReader()
	if err != nil {
		return nil, "", fmt.Errorf("form multipart tidak valid: %w", err)
	}

	var image []byte
	var sourceID string
	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, "", fmt.Errorf("baca form multipart: %w", err)
		}
		switch part.FormName() {
		case "image":
			image, err = io.ReadAll(io.LimitReader(part, maxUploadBytes+1))
			if err != nil {
				part.Close()
				return nil, "", fmt.Errorf("baca gambar: %w", err)
			}
		case "source_id":
			raw, readErr := io.ReadAll(io.LimitReader(part, 128))
			if readErr != nil {
				part.Close()
				return nil, "", fmt.Errorf("baca source_id: %w", readErr)
			}
			sourceID = strings.TrimSpace(string(raw))
		}
		part.Close()
	}
	return image, sourceID, nil
}

func detectImageMime(data []byte) (string, error) {
	switch {
	case len(data) >= 3 && bytes.Equal(data[:3], []byte{0xff, 0xd8, 0xff}):
		return "image/jpeg", nil
	case len(data) >= 8 && bytes.Equal(data[:8], []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a}):
		return "image/png", nil
	case len(data) >= 12 && string(data[:4]) == "RIFF" && string(data[8:12]) == "WEBP":
		return "image/webp", nil
	default:
		return "", fmt.Errorf("format gambar harus JPG, PNG, atau WEBP")
	}
}

func clientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		if first := strings.TrimSpace(strings.Split(forwarded, ",")[0]); first != "" {
			return first
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, _ = io.WriteString(w, message)
}

func truncate(value string, limit int) string {
	if len(value) <= limit {
		return value
	}
	return value[:limit] + "…"
}

func (h *handler) withCORS(next http.Handler) http.Handler {
	if h.config.AllowedOrigin == "" {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == h.config.AllowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
