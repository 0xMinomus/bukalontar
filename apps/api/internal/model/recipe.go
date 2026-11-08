package model

import "strings"

type Recipe struct {
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
	PhytochemicalCompounds  []string `json:"phytochemical_compounds"`
	ModernAction            string   `json:"modern_pharmacological_action"`
	ReferenceStudies        []string `json:"reference_studies"`
}

func (r *Recipe) Normalize() {
	r.PlantLocalName = strings.TrimSpace(r.PlantLocalName)
	r.PlantScientificName = strings.TrimSpace(r.PlantScientificName)
	r.BotanicalFamily = strings.TrimSpace(r.BotanicalFamily)
	r.PlantPartUsed = strings.TrimSpace(r.PlantPartUsed)
	r.EnergyCharacter = strings.TrimSpace(r.EnergyCharacter)
	r.FormulationType = strings.TrimSpace(r.FormulationType)
	r.TargetAilment = strings.TrimSpace(r.TargetAilment)
	r.PreparationInstructions = strings.TrimSpace(r.PreparationInstructions)
	r.ModernAction = strings.TrimSpace(r.ModernAction)
	r.Adjuvants = cleanList(r.Adjuvants)
	r.PhytochemicalCompounds = cleanList(r.PhytochemicalCompounds)
	r.ReferenceStudies = cleanList(r.ReferenceStudies)

	if !validPlantPart(r.PlantPartUsed) {
		r.PlantPartUsed = "Campuran"
	}
	if !validEnergy(r.EnergyCharacter) {
		r.EnergyCharacter = "Dumelada (Sedang)"
	}
	if !validFormulation(r.FormulationType) {
		r.FormulationType = "Loloh (Minum)"
	}
}

func (r Recipe) Empty() bool {
	return strings.TrimSpace(r.PlantLocalName) == "" &&
		strings.TrimSpace(r.TargetAilment) == "" &&
		strings.TrimSpace(r.PreparationInstructions) == ""
}

func (r Recipe) ToEthnomedicinal() Ethnomedicinal {
	return Ethnomedicinal{
		PlantLocalName:          r.PlantLocalName,
		PlantScientificName:     r.PlantScientificName,
		BotanicalFamily:         r.BotanicalFamily,
		KBMatch:                 r.KBMatch,
		PlantPartUsed:           r.PlantPartUsed,
		EnergyCharacter:         r.EnergyCharacter,
		FormulationType:         r.FormulationType,
		Adjuvants:               r.Adjuvants,
		TargetAilment:           r.TargetAilment,
		PreparationInstructions: r.PreparationInstructions,
	}
}

func (r Recipe) ToScientific() Scientific {
	return Scientific{
		PhytochemicalCompounds: r.PhytochemicalCompounds,
		ModernAction:           r.ModernAction,
		ReferenceStudies:       r.ReferenceStudies,
		SafetyDisclaimer:       SafetyDisclaimer,
	}
}

func RecipeFromResponse(response *AnalysisResponse) Recipe {
	return Recipe{
		PlantLocalName:          response.Ethnomedicinal.PlantLocalName,
		PlantScientificName:     response.Ethnomedicinal.PlantScientificName,
		BotanicalFamily:         response.Ethnomedicinal.BotanicalFamily,
		KBMatch:                 response.Ethnomedicinal.KBMatch,
		PlantPartUsed:           response.Ethnomedicinal.PlantPartUsed,
		EnergyCharacter:         response.Ethnomedicinal.EnergyCharacter,
		FormulationType:         response.Ethnomedicinal.FormulationType,
		Adjuvants:               response.Ethnomedicinal.Adjuvants,
		TargetAilment:           response.Ethnomedicinal.TargetAilment,
		PreparationInstructions: response.Ethnomedicinal.PreparationInstructions,
		PhytochemicalCompounds:  response.Scientific.PhytochemicalCompounds,
		ModernAction:            response.Scientific.ModernAction,
		ReferenceStudies:        response.Scientific.ReferenceStudies,
	}
}

func validPlantPart(value string) bool {
	switch value {
	case "Don (Daun)", "Babakan (Kulit Kayu)", "Akah (Akar)", "Buah", "Sekar (Bunga)", "Campuran":
		return true
	default:
		return false
	}
}

func validEnergy(value string) bool {
	switch value {
	case "Panes (Hangat)", "Tis (Dingin)", "Dumelada (Sedang)":
		return true
	default:
		return false
	}
}

func validFormulation(value string) bool {
	switch value {
	case "Loloh (Minum)", "Boreh (Balur)", "Tutuh (Tetes)", "Sasap (Kompres)":
		return true
	default:
		return false
	}
}
