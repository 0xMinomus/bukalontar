package kb

import (
	"net/url"
	"strings"
	"testing"
)

func TestLookupAliases(t *testing.T) {
	store, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string]string{
		"Taru Sotong":      "Psidium guajava",
		"jambu biji":       "Psidium guajava",
		"kalyasĕm":         "Syzygium polycephalum",
		"Kayu Pule":        "Alstonia scholaris",
		"Turi Putih":       "Sesbania grandiflora",
		"Tabia Bun":        "Capsicum annuum",
		"dapdap":           "Erythrina spp.",
		"sĕmbuŋ":           "Blumea balsamifera",
		"kacĕmcĕm":         "Spondias pinnata",
		"poḥgdhaŋ":         "Mangifera indica",
		"don sĕmbuŋ":       "Blumea balsamifera",
		"Daun Taru Sotong": "Psidium guajava",
		"kulit kayu Pule":  "Alstonia scholaris",
		"Psidium guajava":  "Psidium guajava",
	}
	for name, scientific := range cases {
		entry, ok := store.Lookup(name)
		if !ok {
			t.Fatalf("lookup %q failed", name)
		}
		if entry.ScientificName != scientific {
			t.Fatalf("lookup %q: got %q, want %q", name, entry.ScientificName, scientific)
		}
	}
	for _, name := range []string{"Taru Khayalan", "bwaḥ base", "humbintuñtuŋ", "daun dewa", "leunca", "Sotong dan Kunyit", "Sotong yang belum dikenal", "Sirih varietas tidak dikenal", "pacaŋl̥pas·"} {
		if entry, ok := store.Lookup(name); ok {
			t.Fatalf("unknown or ambiguous plant %q matched %s", name, entry.ScientificName)
		}
	}
}

func TestCatalogListsEntries(t *testing.T) {
	store, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	catalog := store.Catalog()
	if !strings.Contains(catalog, "Psidium guajava") || !strings.Contains(catalog, "Alstonia scholaris") {
		t.Fatalf("catalog missing expected entries:\n%s", catalog)
	}
}

func TestCuratedReferenceLinks(t *testing.T) {
	store, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range store.entries {
		if len(entry.ReferenceLinks) == 0 {
			t.Fatalf("%s has no reference links", entry.ID)
		}
		for _, link := range entry.ReferenceLinks {
			parsed, err := url.Parse(link.URL)
			if err != nil || parsed.Scheme != "https" || parsed.Host != "pubmed.ncbi.nlm.nih.gov" || link.Title == "" || (link.Kind != "literature" && link.Kind != "literature_search") {
				t.Fatalf("%s has invalid reference link: %#v", entry.ID, link)
			}
		}
	}
}

func TestCatalogAliasesDoNotOverwriteOtherPlants(t *testing.T) {
	store, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if store.Count() < 1001 {
		t.Fatalf("botanical coverage below requested 1000+: %d", store.Count())
	}
	ids := make(map[string]bool)
	aliases := make(map[string]string)
	for _, entry := range store.entries {
		if entry.ID == "" || ids[entry.ID] {
			t.Fatalf("empty or duplicate plant ID: %q", entry.ID)
		}
		ids[entry.ID] = true
		for _, name := range entry.LocalNames {
			key := normalize(name)
			if previous, exists := aliases[key]; exists && previous != entry.ID {
				t.Fatalf("alias %q identifies both %s and %s", name, previous, entry.ID)
			}
			aliases[key] = entry.ID
			found, ok := store.Lookup(name)
			if !ok || found.ID != entry.ID {
				t.Fatalf("alias %q resolves to %s instead of %s", name, found.ID, entry.ID)
			}
		}
	}
}
