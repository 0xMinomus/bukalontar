package ai

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/0xMinomus/digifest/apps/api/internal/model"
)

const (
	promptVersion = "2026-11-08.1"
	cacheTTL      = time.Hour
	cacheLimit    = 200
)

type cacheEntry struct {
	result  []byte
	expires time.Time
}

func (s *Service) cacheKey(image []byte, sourceID string) string {
	stable := sourceID
	if stable == "" || stable == "upload" {
		sum := sha256.Sum256(image)
		stable = hex.EncodeToString(sum[:])
	}
	return fmt.Sprintf("%s|%s|%s|%s", s.primary.Model(), promptVersion, sourceID, stable)
}

func (s *Service) cached(key string) (*model.AnalysisResponse, bool) {
	s.cacheMu.Lock()
	entry, ok := s.cache[key]
	if !ok {
		s.cacheMu.Unlock()
		return nil, false
	}
	if time.Now().After(entry.expires) {
		delete(s.cache, key)
		s.cacheMu.Unlock()
		return nil, false
	}
	raw := append([]byte(nil), entry.result...)
	s.cacheMu.Unlock()

	var response model.AnalysisResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		return nil, false
	}
	return &response, true
}

func (s *Service) storeCache(key string, response *model.AnalysisResponse) {
	raw, err := json.Marshal(response)
	if err != nil {
		return
	}
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	if len(s.cache) >= cacheLimit {
		now := time.Now()
		for cacheKey, entry := range s.cache {
			if now.After(entry.expires) {
				delete(s.cache, cacheKey)
			}
		}
		if len(s.cache) >= cacheLimit {
			s.cache = make(map[string]cacheEntry)
		}
	}
	s.cache[key] = cacheEntry{result: raw, expires: time.Now().Add(cacheTTL)}
}
