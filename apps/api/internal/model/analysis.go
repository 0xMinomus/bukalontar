package model

import (
	"fmt"
	"strings"
)

const SafetyDisclaimer = "Informasi ini untuk preservasi warisan budaya dan edukasi fitofarmaka. Untuk penanganan medis, konsultasikan dengan dokter atau tenaga kesehatan profesional."

type AnalysisResponse struct {
	ManuscriptID    string          `json:"manuscript_id"`
	SourceReference SourceReference `json:"source_reference"`
	Paleography     Paleography     `json:"paleography"`
	Ethnomedicinal  Ethnomedicinal  `json:"ethnomedicinal_entity"`
	Scientific      Scientific      `json:"scientific_validation"`
}

type SourceReference struct {
	LontarGroup  string `json:"lontar_group"`
	LempirNumber *int   `json:"lempir_number,omitempty"`
	Language     string `json:"language"`
}

type Paleography struct {
	AksaraBaliDetected    bool   `json:"aksara_bali_detected"`
	Confidence            string `json:"confidence"`
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

type Scientific struct {
	PhytochemicalCompounds []string `json:"phytochemical_compounds"`
	ModernAction           string   `json:"modern_pharmacological_action"`
	ReferenceStudies       []string `json:"reference_studies"`
	SafetyDisclaimer       string   `json:"safety_disclaimer"`
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
		"Usadha Taru Pramana", "Usadha Tiwang", "Usadha Dalem", "Usadha Yeh",
	})
	enum(r.SourceReference.Language, "source_reference.language", []string{
		"Kawi / Jawa Kuno", "Bali Alus",
	})
	enum(r.Paleography.Confidence, "paleography.confidence", []string{"low", "medium", "high"})
	require(r.Paleography.Transliteration, "paleography.transliteration")
	require(r.Paleography.IndonesianTranslation, "paleography.indonesian_translation")
	require(r.Paleography.EnglishTranslation, "paleography.english_translation")
	enum(r.Ethnomedicinal.PlantPartUsed, "ethnomedicinal_entity.plant_part_used", []string{
		"Don (Daun)", "Babakan (Kulit Kayu)", "Akah (Akar)", "Buah", "Sekar (Bunga)", "Campuran",
	})
	enum(r.Ethnomedicinal.EnergyCharacter, "ethnomedicinal_entity.energy_character", []string{
		"Panes (Hangat)", "Tis (Dingin)", "Dumelada (Sedang)",
	})
	enum(r.Ethnomedicinal.FormulationType, "ethnomedicinal_entity.formulation_type", []string{
		"Loloh (Minum)", "Boreh (Balur)", "Tutuh (Tetes)", "Sasap (Kompres)",
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
