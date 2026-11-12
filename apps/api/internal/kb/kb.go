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

//go:embed data/wikidata-plants.json
var botanicalSeed []byte

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
	for _, data := range [][]byte{seed, botanicalSeed} {
		var batch []Entry
		if err := json.Unmarshal(data, &batch); err != nil {
			return nil, fmt.Errorf("parse botanical KB: %w", err)
		}
		entries = append(entries, batch...)
	}
	store := &Store{entries: entries, aliases: make(map[string]int)}
	ids := make(map[string]bool, len(entries))
	for index, entry := range entries {
		if entry.ID == "" || ids[entry.ID] || len(entry.LocalNames) == 0 || entry.ScientificName == "" || entry.BotanicalFamily == "" {
			return nil, fmt.Errorf("invalid or duplicate botanical entry: %q", entry.ID)
		}
		ids[entry.ID] = true
		for nameIndex := 0; nameIndex <= len(entry.LocalNames); nameIndex++ {
			name := entry.ScientificName
			if nameIndex > 0 {
				name = entry.LocalNames[nameIndex-1]
			}
			normalized := normalize(name)
			if normalized == "" {
				return nil, fmt.Errorf("empty botanical alias in %q", entry.ID)
			}
			if previous, exists := store.aliases[normalized]; exists && previous != index {
				return nil, fmt.Errorf("ambiguous botanical alias %q: %s and %s", name, entries[previous].ID, entry.ID)
			}
			store.aliases[normalized] = index
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
	for _, prefix := range [...]string{"kulit kayu ", "babakan ", "daun ", "don ", "akar ", "akah ", "buah ", "bwah ", "bunga ", "sekar "} {
		if name, ok := strings.CutPrefix(normalized, prefix); ok {
			if index, found := s.aliases[strings.TrimPrefix(name, "taru ")]; found {
				return s.entries[index], true
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

var transliterationLetters = strings.NewReplacer(
	"ĕ", "e", "ə", "e",
	"ā", "a", "ī", "i", "ū", "u",
	"ṭ", "t", "ḍ", "d", "ṇ", "n",
	"ś", "s", "ḥ", "h",
	"ŋ", "ng", "ṅ", "ng", "ŕ", "r",
)

func normalize(value string) string {
	value = transliterationLetters.Replace(strings.ToLower(strings.TrimSpace(value)))
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
