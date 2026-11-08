package ai

import (
	"strings"
	"testing"

	"github.com/0xMinomus/digifest/apps/api/internal/model"
)

func TestParseAnalysisStripsCodeFence(t *testing.T) {
	content := "```json\n" + `{"manuscript_id":"x","source_reference":{"lontar_group":"Usadha Taru Pramana","language":"Kawi / Jawa Kuno"},"paleography":{"aksara_bali_detected":true,"confidence":"high","transliteration":"a","indonesian_translation":"b","english_translation":"c"},"ethnomedicinal_entity":{"plant_local_name":"x","plant_scientific_name":"","botanical_family":"","kb_match":false,"plant_part_used":"Don (Daun)","energy_character":"Panes (Hangat)","formulation_type":"Loloh (Minum)","adjuvants":[],"target_ailment":"x","preparation_instructions":"x"},"scientific_validation":{"phytochemical_compounds":[],"modern_pharmacological_action":"x","reference_studies":[],"safety_disclaimer":"` + model.SafetyDisclaimer + `"}}` + "\n```"
	response, err := parseAnalysis(content)
	if err != nil {
		t.Fatal(err)
	}
	response.Normalize()
	if err := response.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestParseAnalysisRejectsMissingJSON(t *testing.T) {
	_, err := parseAnalysis("aneh, bukan JSON")
	if err == nil || !strings.Contains(err.Error(), "JSON") {
		t.Fatalf("expected JSON error, got %v", err)
	}
}
