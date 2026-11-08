package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port               string
	PublicOrigin       string
	AllowedOrigin      string
	AIBaseURL          string
	AIModel            string
	AIAPIKey           string
	AITimeout          time.Duration
	FallbackEnabled    bool
	FallbackBaseURL    string
	FallbackModel      string
	FallbackAPIKey     string
	RateLimitPerMinute int
}

func Load() Config {
	timeoutSeconds := envInt("AI_TIMEOUT_SECONDS", 120)
	return Config{
		Port:               env("PORT", "8080"),
		PublicOrigin:       env("PUBLIC_ORIGIN", "http://localhost:5173"),
		AllowedOrigin:      env("ALLOWED_ORIGIN", ""),
		AIBaseURL:          env("AI_BASE_URL", "http://127.0.0.1:20128/v1"),
		AIModel:            env("AI_MODEL", "ag/gemini-3.8-flash-low"),
		AIAPIKey:           env("AI_API_KEY", ""),
		AITimeout:          time.Duration(timeoutSeconds) * time.Second,
		FallbackEnabled:    envBool("AI_FALLBACK_ENABLED", false),
		FallbackBaseURL:    env("AI_FALLBACK_BASE_URL", ""),
		FallbackModel:      env("AI_FALLBACK_MODEL", ""),
		FallbackAPIKey:     env("AI_FALLBACK_API_KEY", ""),
		RateLimitPerMinute: envInt("RATE_LIMIT_PER_MINUTE", 10),
	}
}

// LoadDotEnv loads the first existing .env file without overriding real
// environment variables. It supports the local `apps/api` and root layouts.
func LoadDotEnv() {
	for _, path := range []string{".env", filepath.Join("apps", "api", ".env")} {
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			key, value, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}
			key = strings.TrimSpace(key)
			value = strings.Trim(strings.TrimSpace(value), `"'`)
			if key == "" || os.Getenv(key) != "" {
				continue
			}
			_ = os.Setenv(key, value)
		}
		file.Close()
		return
	}
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value, err := strconv.Atoi(strings.TrimSpace(os.Getenv(key)))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

func envBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
