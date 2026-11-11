package kb

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"github.com/0xMinomus/bukalontar/apps/api/internal/model"
)

//go:embed data/taru-pramana.json
var seed []byte

type Entry struct {
	ID                     string                `json:"id"`
	LocalNames             []string              `json:"local_names"`
	ScientificName         string                `json:"scientific_name"`
	BotanicalFamily        string                `json:"botanical_family"`
	EnergyCharacter        string                `json:"energy_character"`
	PlantPartUsed          string                `json:"plant_part_used"`
	PhytochemicalCompounds []string              `json:"phytochemical_compounds"`
	ModernAction           string                `json:"modern_pharmacological_action"`
	ReferenceStudies       []string              `json:"reference_studies"`
	ReferenceLinks         []model.ReferenceLink `json:"reference_links"`
}

type Store struct {
	entries []Entry
	aliases map[string]int
}

func Load() (*Store, error) {
	var entries []Entry
	if err := json.Unmarshal(seed, &entries); err != nil {
		return nil, fmt.Errorf("parse Taru Pramana KB: %w", err)
	}
	store := &Store{entries: entries, aliases: make(map[string]int)}
	for index, entry := range entries {
		for _, name := range entry.LocalNames {
			normalized := normalize(name)
			if normalized != "" {
				store.aliases[normalized] = index
			}
		}
	}
	if len(store.entries) == 0 {
		return nil, fmt.Errorf("Taru Pramana KB is empty")
	}
	return store, nil
}

func (s *Store) Count() int {
	return len(s.entries)
}

func (s *Store) Lookup(localName string) (Entry, bool) {
	normalized := normalize(localName)
	if normalized == "" {
		return Entry{}, false
	}
	if index, ok := s.aliases[normalized]; ok {
		return s.entries[index], true
	}
	padded := " " + normalized + " "
	for _, entry := range s.entries {
		for _, name := range entry.LocalNames {
			alias := normalize(name)
			if alias != "" && strings.Contains(padded, " "+alias+" ") {
				return entry, true
			}
		}
	}
	return Entry{}, false
}

func (s *Store) Catalog() string {
	var builder strings.Builder
	for _, entry := range s.entries {
		builder.WriteString("- ")
		builder.WriteString(entry.LocalNames[0])
		builder.WriteString(" (aliases: ")
		builder.WriteString(strings.Join(entry.LocalNames[1:], ", "))
		builder.WriteString(") -> ")
		builder.WriteString(entry.ScientificName)
		builder.WriteString(" (")
		builder.WriteString(entry.BotanicalFamily)
		builder.WriteString("); energy ")
		builder.WriteString(entry.EnergyCharacter)
		builder.WriteString("; part ")
		builder.WriteString(entry.PlantPartUsed)
		builder.WriteByte('\n')
	}
	return strings.TrimSpace(builder.String())
}

func normalize(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.NewReplacer(
		"ĕ", "e", "ə", "e",
		"ā", "a", "ī", "i", "ū", "u",
		"ṭ", "t", "ḍ", "d", "ṇ", "n",
		"ś", "s", "ḥ", "h",
	).Replace(value)
	var builder strings.Builder
	for _, char := range value {
		if unicode.IsLetter(char) || unicode.IsDigit(char) || unicode.IsSpace(char) {
			builder.WriteRune(char)
		}
	}
	value = strings.Join(strings.Fields(builder.String()), " ")
	value = strings.TrimPrefix(value, "taru ")
	if value == "taru" {
		return ""
	}
	return value
}
