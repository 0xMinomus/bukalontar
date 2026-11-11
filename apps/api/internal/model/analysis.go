package model

import (
	"fmt"
	"strings"
)

const SafetyDisclaimer = "Informasi ini untuk preservasi warisan budaya dan edukasi fitofarmaka. Untuk penanganan medis, konsultasikan dengan dokter atau tenaga kesehatan profesional."

type AnalysisResponse struct {
	ManuscriptID      string             `json:"manuscript_id"`
	SourceReference   SourceReference    `json:"source_reference"`
	Paleography       Paleography        `json:"paleography"`
	Ethnomedicinal    Ethnomedicinal     `json:"ethnomedicinal_entity"`
	Scientific        Scientific         `json:"scientific_validation"`
	Recipes           []Recipe           `json:"recipes"`
	ReadingAssistance *ReadingAssistance `json:"reading_assistance,omitempty"`
}

type ReadingAssistance struct {
	Mode           string   `json:"mode"`
	OCRModel       string   `json:"ocr_model,omitempty"`
	OCRMS          int64    `json:"ocr_ms"`
	FocusRegions   int      `json:"focus_regions"`
	LeafRegions    int      `json:"leaf_regions"`
	ExampleCount   int      `json:"example_count"`
	ExampleSources []string `json:"example_sources"`
	Warnings       []string `json:"warnings"`
}

type SourceReference struct {
	LontarGroup  string `json:"lontar_group"`
	LempirNumber *int   `json:"lempir_number,omitempty"`
	Language     string `json:"language"`
}

type Paleography struct {
	AksaraBaliDetected    bool   `json:"aksara_bali_detected"`
	Confidence            string `json:"confidence"`
	BalineseTranscription string `json:"balinese_transcription,omitempty"`
	Transliteration       string `json:"transliteration"`
	IndonesianTranslation string `json:"indonesian_translation"`
	EnglishTranslation    string `json:"english_translation"`
}

type Ethnomedicinal struct {
	PlantLocalName          string   `json:"plant_local_name"`
	PlantScientificName     string   `json:"plant_scientific_name"`
	BotanicalFamily         string   `json:"botanical_family"`
	KBMatch                 bool     `json:"kb_match"`
	PlantPartUsed           string   `json:"plant_part_used"`
	EnergyCharacter         string   `json:"energy_character"`
	FormulationType         string   `json:"formulation_type"`
	Adjuvants               []string `json:"adjuvants"`
	TargetAilment           string   `json:"target_ailment"`
	PreparationInstructions string   `json:"preparation_instructions"`
}

type ReferenceLink struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Kind  string `json:"kind"`
}

type Scientific struct {
	PhytochemicalCompounds []string        `json:"phytochemical_compounds"`
	ModernAction           string          `json:"modern_pharmacological_action"`
	ReferenceStudies       []string        `json:"reference_studies"`
	ReferenceLinks         []ReferenceLink `json:"reference_links"`
	SafetyDisclaimer       string          `json:"safety_disclaimer"`
}

type ValidationError struct {
	Problems []string
}

func (e *ValidationError) Error() string {
	return "invalid analysis response: " + strings.Join(e.Problems, "; ")
}

func (r *AnalysisResponse) Normalize() {
	r.ManuscriptID = strings.TrimSpace(r.ManuscriptID)
	r.SourceReference.LontarGroup = strings.TrimSpace(r.SourceReference.LontarGroup)
	r.SourceReference.Language = strings.TrimSpace(r.SourceReference.Language)
	r.Paleography.Confidence = strings.TrimSpace(r.Paleography.Confidence)
	r.Paleography.Transliteration = strings.TrimSpace(r.Paleography.Transliteration)
	r.Paleography.BalineseTranscription = strings.TrimSpace(r.Paleography.BalineseTranscription)
	r.Paleography.IndonesianTranslation = strings.TrimSpace(r.Paleography.IndonesianTranslation)
	r.Paleography.EnglishTranslation = strings.TrimSpace(r.Paleography.EnglishTranslation)
	r.Ethnomedicinal.PlantLocalName = strings.TrimSpace(r.Ethnomedicinal.PlantLocalName)
	r.Ethnomedicinal.PlantScientificName = strings.TrimSpace(r.Ethnomedicinal.PlantScientificName)
	r.Ethnomedicinal.BotanicalFamily = strings.TrimSpace(r.Ethnomedicinal.BotanicalFamily)
	r.Ethnomedicinal.PlantPartUsed = strings.TrimSpace(r.Ethnomedicinal.PlantPartUsed)
	r.Ethnomedicinal.EnergyCharacter = strings.TrimSpace(r.Ethnomedicinal.EnergyCharacter)
	r.Ethnomedicinal.FormulationType = strings.TrimSpace(r.Ethnomedicinal.FormulationType)
	r.Ethnomedicinal.TargetAilment = strings.TrimSpace(r.Ethnomedicinal.TargetAilment)
	r.Ethnomedicinal.PreparationInstructions = strings.TrimSpace(r.Ethnomedicinal.PreparationInstructions)
	r.Scientific.ModernAction = strings.TrimSpace(r.Scientific.ModernAction)
	r.Scientific.SafetyDisclaimer = strings.TrimSpace(r.Scientific.SafetyDisclaimer)
	r.Ethnomedicinal.Adjuvants = cleanList(r.Ethnomedicinal.Adjuvants)
	r.Scientific.PhytochemicalCompounds = cleanList(r.Scientific.PhytochemicalCompounds)
	r.Scientific.ReferenceStudies = cleanList(r.Scientific.ReferenceStudies)

	for index := range r.Recipes {
		r.Recipes[index].Normalize()
	}
	filtered := r.Recipes[:0]
	for _, recipe := range r.Recipes {
		if !recipe.Empty() {
			filtered = append(filtered, recipe)
		}
	}
	r.Recipes = filtered
	if len(r.Recipes) == 0 {
		recipe := RecipeFromResponse(r)
		if !recipe.Empty() {
			recipe.Normalize()
			r.Recipes = []Recipe{recipe}
		}
	}
	if r.Recipes == nil {
		r.Recipes = []Recipe{}
	}
	if len(r.Recipes) > 6 {
		r.Recipes = r.Recipes[:6]
	}
	if len(r.Recipes) > 0 {
		r.Ethnomedicinal = r.Recipes[0].ToEthnomedicinal()
		r.Scientific = r.Recipes[0].ToScientific()
	}
	if r.Scientific.ReferenceLinks == nil {
		r.Scientific.ReferenceLinks = []ReferenceLink{}
	}

	if r.SourceReference.LempirNumber != nil && *r.SourceReference.LempirNumber <= 0 {
		r.SourceReference.LempirNumber = nil
	}
	r.Scientific.SafetyDisclaimer = SafetyDisclaimer
}

func (r *AnalysisResponse) Validate() error {
	problems := make([]string, 0, 8)
	require := func(value, field string) {
		if strings.TrimSpace(value) == "" {
			problems = append(problems, field+" wajib diisi")
		}
	}
	enum := func(value, field string, allowed []string) {
		for _, candidate := range allowed {
			if value == candidate {
				return
			}
		}
		problems = append(problems, field+" tidak valid: "+value)
	}

	if !r.Paleography.AksaraBaliDetected {
		if strings.TrimSpace(r.ManuscriptID) == "" {
			return &ValidationError{Problems: []string{"manuscript_id wajib diisi"}}
		}
		return nil
	}

	require(r.ManuscriptID, "manuscript_id")
	enum(r.SourceReference.LontarGroup, "source_reference.lontar_group", []string{
		"Usadha Taru Pramana", "Usadha Tiwang", "Usadha Dalem", "Usadha Yeh", "Tidak diketahui", "Non-Usadha",
	})
	enum(r.SourceReference.Language, "source_reference.language", []string{
		"Kawi / Jawa Kuno", "Bali", "Bali Alus", "Tidak diketahui",
	})
	enum(r.Paleography.Confidence, "paleography.confidence", []string{"low", "medium", "high"})
	require(r.Paleography.Transliteration, "paleography.transliteration")
	require(r.Paleography.IndonesianTranslation, "paleography.indonesian_translation")
	require(r.Paleography.EnglishTranslation, "paleography.english_translation")
	enum(r.Ethnomedicinal.PlantPartUsed, "ethnomedicinal_entity.plant_part_used", []string{
		"Don (Daun)", "Babakan (Kulit Kayu)", "Akah (Akar)", "Buah", "Sekar (Bunga)", "Campuran", "Tidak disebutkan.",
	})
	enum(r.Ethnomedicinal.EnergyCharacter, "ethnomedicinal_entity.energy_character", []string{
		"Panes (Hangat)", "Tis (Dingin)", "Dumelada (Sedang)", "Tidak disebutkan.",
	})
	enum(r.Ethnomedicinal.FormulationType, "ethnomedicinal_entity.formulation_type", []string{
		"Loloh (Minum)", "Boreh (Balur)", "Tutuh (Tetes)", "Sasap (Kompres)", "Tidak disebutkan.",
	})
	require(r.Ethnomedicinal.PlantLocalName, "ethnomedicinal_entity.plant_local_name")
	require(r.Ethnomedicinal.TargetAilment, "ethnomedicinal_entity.target_ailment")
	require(r.Ethnomedicinal.PreparationInstructions, "ethnomedicinal_entity.preparation_instructions")
	require(r.Scientific.SafetyDisclaimer, "scientific_validation.safety_disclaimer")

	if len(problems) > 0 {
		return &ValidationError{Problems: problems}
	}
	return nil
}

func cleanList(values []string) []string {
	cleaned := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			cleaned = append(cleaned, value)
		}
	}
	return cleaned
}

func (e *ValidationError) Summary() string {
	if len(e.Problems) == 0 {
		return e.Error()
	}
	return fmt.Sprintf("%s (%d masalah)", e.Error(), len(e.Problems))
}
