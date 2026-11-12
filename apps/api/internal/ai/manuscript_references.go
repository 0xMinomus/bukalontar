package ai

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"github.com/0xMinomus/bukalontar/apps/api/internal/model"
)

//go:embed manuscript_references.json
var manuscriptReferenceData []byte

type manuscriptReference struct {
	model.ReadingReference
	ID              string `json:"id"`
	SourceText      string `json:"source_text"`
	Transliteration string `json:"transliteration"`
	hash            [sha256.Size]byte
}

func loadManuscriptReferences() ([]manuscriptReference, error) {
	var references []manuscriptReference
	if err := json.Unmarshal(manuscriptReferenceData, &references); err != nil {
		return nil, fmt.Errorf("load manuscript references: %w", err)
	}
	withoutSpaces := func(value string) string {
		return strings.Map(func(char rune) rune {
			if unicode.IsSpace(char) {
				return -1
			}
			return char
		}, value)
	}
	for index := range references {
		reference := &references[index]
		if len(reference.ImageSHA256) != sha256.Size*2 {
			return nil, fmt.Errorf("invalid reference image hash: %s", reference.ID)
		}
		if _, err := hex.Decode(reference.hash[:], []byte(reference.ImageSHA256)); err != nil {
			return nil, fmt.Errorf("decode reference image hash %s: %w", reference.ID, err)
		}
		if strings.TrimSpace(reference.SourceText) == "" || withoutSpaces(reference.SourceText) != withoutSpaces(reference.Transliteration) {
			return nil, fmt.Errorf("editorial spacing changes source characters: %s", reference.ID)
		}
		reference.Method = "source_transcription"
	}
	return references, nil
}

func (s *Service) referenceForImage(image []byte) *manuscriptReference {
	hash := sha256.Sum256(image)
	for index := range s.references {
		if s.references[index].hash == hash {
			return &s.references[index]
		}
	}
	return nil
}
