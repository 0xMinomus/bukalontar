package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/0xMinomus/digifest/apps/api/internal/model"
)

func extractJSON(content string) ([]byte, error) {
	content = strings.TrimSpace(strings.TrimPrefix(content, "\ufeff"))
	if strings.HasPrefix(content, "```") {
		if newline := strings.IndexByte(content, '\n'); newline >= 0 {
			content = content[newline+1:]
		}
		if fence := strings.LastIndex(content, "```"); fence >= 0 {
			content = content[:fence]
		}
	}
	start := strings.IndexByte(content, '{')
	end := strings.LastIndexByte(content, '}')
	if start < 0 || end <= start {
		return nil, fmt.Errorf("tidak ada objek JSON pada respons AI")
	}
	raw := []byte(strings.TrimSpace(content[start : end+1]))
	if !json.Valid(raw) {
		return nil, fmt.Errorf("objek JSON pada respons AI tidak valid")
	}
	return raw, nil
}

func parseAnalysis(content string) (*model.AnalysisResponse, error) {
	raw, err := extractJSON(content)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	var response model.AnalysisResponse
	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("decode AnalysisResponse: %w", err)
	}
	return &response, nil
}
