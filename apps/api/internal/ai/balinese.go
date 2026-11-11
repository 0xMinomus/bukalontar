package ai

import (
	"strings"

	"github.com/0xMinomus/bukalontar/apps/api/internal/model"
)

// sanitizeBalinese drops characters from lookalike scripts (Javanese, Kawi,
// Devanagari, Sundanese) that a multimodal model can emit when it mistakes the
// writing system. Balinese is U+1B00-U+1B7F; Latin, digits, punctuation, and
// the Indonesian "[tidak terbaca]" marker remain.
func sanitizeBalinese(paleography *model.Paleography) int {
	removed := 0
	var builder strings.Builder
	builder.Grow(len(paleography.BalineseTranscription))
	for _, char := range paleography.BalineseTranscription {
		if !foreignScript(char) {
			builder.WriteRune(char)
			continue
		}
		removed++
	}
	paleography.BalineseTranscription = builder.String()
	return removed
}

func hasBalinese(value string) bool {
	for _, char := range value {
		if char >= 0x1B00 && char <= 0x1B7F {
			return true
		}
	}
	return false
}

func foreignScript(char rune) bool {
	switch {
	case char >= 0x1B00 && char <= 0x1B7F:
		return false // Balinese
	case char >= 0x0900 && char <= 0x1CFF:
		return true // Devanagari through Sundanese, including Javanese
	case char >= 0xA980 && char <= 0xA9DF:
		return true // Javanese extended
	case char >= 0x11F00 && char <= 0x11F5F:
		return true // Kawi
	default:
		return false
	}
}
