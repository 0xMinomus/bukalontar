package model

import (
	"strings"
	"testing"
)

func validResponse() AnalysisResponse {
	lempir := 3
	return AnalysisResponse{
		ManuscriptID: "taru-pramana-1384-10-lp-03",
		SourceReference: SourceReference{
			LontarGroup:  "Usadha Taru Pramana",
			LempirNumber: &lempir,
			Language:     "Kawi / Jawa Kuno",
		},
		Paleography: Paleography{
			AksaraBaliDetected:    true,
			Confidence:            "high",
			Transliteration:       "Tityang taru sotong",
			IndonesianTranslation: "Hamba pohon jambu biji",
			EnglishTranslation:    "I am the guava tree",
		},
		Ethnomedicinal: Ethnomedicinal{
			PlantLocalName:          "Taru Sotong",
			PlantScientificName:     "Psidium guajava",
			BotanicalFamily:         "Myrtaceae",
			KBMatch:                 true,
			PlantPartUsed:           "Don (Daun)",
			EnergyCharacter:         "Panes (Hangat)",
			FormulationType:         "Sasap (Kompres)",
			Adjuvants:               []string{"Ketumbar bolong"},
			TargetAilment:           "Mising (diare)",
			PreparationInstructions: "Pucuk daun ditumbuk lalu ditempelkan pada pusar.",
		},
		Scientific: Scientific{
			PhytochemicalCompounds: []string{"Tannin"},
			ModernAction:           "Antimikroba.",
			ReferenceStudies:       []string{"Tinjauan fitokimia."},
			SafetyDisclaimer:       SafetyDisclaimer,
		},
	}
}

func TestValidateAcceptsCanonicalResponse(t *testing.T) {
	response := validResponse()
	if err := response.Validate(); err != nil {
		t.Fatalf("valid response rejected: %v", err)
	}
}

func TestValidateAllowsNotDetectedGuard(t *testing.T) {
	response := AnalysisResponse{ManuscriptID: "x"}
	if err := response.Validate(); err != nil {
		t.Fatalf("not-detected guard should pass so handler can return 422: %v", err)
	}
}

func TestValidateRejectsBadEnum(t *testing.T) {
	response := validResponse()
	response.Ethnomedicinal.PlantPartUsed = "Daun"
	err := response.Validate()
	if err == nil || !strings.Contains(err.Error(), "plant_part_used") {
		t.Fatalf("expected plant_part_used validation error, got %v", err)
	}
}

func TestNormalizeFillsListsAndCanonicalDisclaimer(t *testing.T) {
	response := validResponse()
	response.Ethnomedicinal.Adjuvants = nil
	response.Scientific.PhytochemicalCompounds = []string{"  ", ""}
	response.Scientific.ReferenceStudies = nil
	response.Scientific.SafetyDisclaimer = "anything"
	response.Normalize()

	if response.Ethnomedicinal.Adjuvants == nil {
		t.Fatal("adjuvants must be an empty array, not nil")
	}
	if len(response.Scientific.PhytochemicalCompounds) != 0 {
		t.Fatalf("blank compounds should be removed: %#v", response.Scientific.PhytochemicalCompounds)
	}
	if response.Scientific.ReferenceStudies == nil {
		t.Fatal("references must be an empty array, not nil")
	}
	if response.Scientific.SafetyDisclaimer != SafetyDisclaimer {
		t.Fatal("safety disclaimer must be canonical")
	}
}
