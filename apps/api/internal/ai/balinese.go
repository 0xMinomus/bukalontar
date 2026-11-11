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

// A rejected target must not retain a reading or acquire recipes during enrichment.
// This checks output consistency; image classification still depends on the model.
func enforceBalineseScope(response *model.AnalysisResponse) {
	conflict := false
	for _, value := range []string{
		response.Paleography.BalineseTranscription,
		response.Paleography.Transliteration,
	} {
		for _, char := range value {
			if foreignScript(char) {
				conflict = true
				break
			}
		}
	}
	if response.Paleography.AksaraBaliDetected && !conflict {
		return
	}
	response.SourceReference = model.SourceReference{
		LontarGroup: "Tidak diketahui",
		Language:    "Tidak diketahui",
	}
	response.Paleography = model.Paleography{
		Confidence:            "low",
		Transliteration:       "Tidak terdeteksi.",
		IndonesianTranslation: "Tidak terdeteksi.",
		EnglishTranslation:    "Tidak terdeteksi.",
	}
	response.Recipes = []model.Recipe{}
	response.Ethnomedicinal = model.Ethnomedicinal{
		PlantLocalName:          "Tidak disebutkan.",
		PlantPartUsed:           "Tidak disebutkan.",
		EnergyCharacter:         "Tidak disebutkan.",
		FormulationType:         "Tidak disebutkan.",
		Adjuvants:               []string{},
		TargetAilment:           "Tidak disebutkan.",
		PreparationInstructions: "Tidak terdeteksi.",
	}
	response.Scientific = model.Scientific{
		PhytochemicalCompounds: []string{},
		ReferenceStudies:       []string{},
		ReferenceLinks:         []model.ReferenceLink{},
		SafetyDisclaimer:       model.SafetyDisclaimer,
	}
}
