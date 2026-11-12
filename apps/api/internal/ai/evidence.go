package ai

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/0xMinomus/bukalontar/apps/api/internal/model"
)

var leafLabel = regexp.MustCompile(`(?i)\[\s*(?:lembar|leaf)\s+([0-9]+)\s*\]`)
var evidenceBreak = regexp.MustCompile(`\[[^\]]*\]|\.(?:\s*\.)+|…`)

// This proves consistency with the reading, not that the reading matches the image.
func retainSupportedRecipes(response *model.AnalysisResponse) {
	reading := response.Paleography.Transliteration
	if labels := leafLabel.FindAllStringSubmatchIndex(reading, 2); len(labels) > 0 {
		first := labels[0]
		number, _ := strconv.Atoi(reading[first[2]:first[3]])
		if number != 1 {
			reading = ""
		} else {
			end := len(reading)
			if len(labels) > 1 {
				end = labels[1][0]
			}
			reading = reading[first[1]:end]
		}
	}
	source := evidenceWords(evidenceBreak.ReplaceAllString(reading, "\x00"))
	hadRecipes := len(response.Recipes) > 0
	kept := response.Recipes[:0]
	for _, recipe := range response.Recipes {
		if evidenceBreak.MatchString(recipe.SourceQuote) {
			continue
		}
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
	value = strings.ToLower(value)
	value = strings.Map(func(char rune) rune {
		if char == 0 || unicode.IsLetter(char) || unicode.IsDigit(char) {
			return char
		}
		return ' '
	}, value)
	return strings.Join(strings.Fields(value), " ")
}

func containsWords(text, quote string) bool {
	return quote != "" && strings.Contains(" "+text+" ", " "+quote+" ")
}
