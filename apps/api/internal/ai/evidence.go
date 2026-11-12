package ai

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/0xMinomus/bukalontar/apps/api/internal/model"
)

var leafLabel = regexp.MustCompile(`(?i)\[\s*lembar\s+[0-9]+\s*\]`)

// This proves consistency with the reading, not that the reading matches the image.
func retainSupportedRecipes(response *model.AnalysisResponse) {
	reading := response.Paleography.Transliteration
	if labels := leafLabel.FindAllStringIndex(reading, 2); len(labels) > 1 {
		reading = reading[:labels[1][0]]
	}
	source := evidenceWords(reading)
	hadRecipes := len(response.Recipes) > 0
	kept := response.Recipes[:0]
	for _, recipe := range response.Recipes {
		quote := evidenceWords(recipe.SourceQuote)
		plant := strings.TrimPrefix(evidenceWords(recipe.PlantLocalName), "taru ")
		if len(strings.Fields(quote)) < 4 || plant == "" ||
			!containsWords(source, quote) || !containsWords(quote, plant) {
			continue
		}
		supported := true
		for _, ingredient := range recipe.Adjuvants {
			if !containsWords(quote, evidenceWords(ingredient)) {
				supported = false
				break
			}
		}
		if supported {
			kept = append(kept, recipe)
		}
	}
	response.Recipes = kept
	if len(kept) != 0 {
		return
	}
	summary := response.Ethnomedicinal.PreparationInstructions
	if hadRecipes || strings.TrimSpace(summary) == "" {
		summary = "Belum ada resep dengan kutipan bacaan yang mendukung."
	}
	response.Ethnomedicinal = model.Ethnomedicinal{
		PlantLocalName:          "Tidak disebutkan.",
		PlantPartUsed:           "Tidak disebutkan.",
		EnergyCharacter:         "Tidak disebutkan.",
		FormulationType:         "Tidak disebutkan.",
		Adjuvants:               []string{},
		TargetAilment:           "Tidak disebutkan.",
		PreparationInstructions: summary,
	}
}

func evidenceWords(value string) string {
	value = strings.ReplaceAll(strings.ToLower(value), "[tidak terbaca]", " ")
	value = leafLabel.ReplaceAllString(value, " ")
	value = strings.Map(func(char rune) rune {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			return char
		}
		return ' '
	}, value)
	return strings.Join(strings.Fields(value), " ")
}

func containsWords(text, quote string) bool {
	return quote != "" && strings.Contains(" "+text+" ", " "+quote+" ")
}
