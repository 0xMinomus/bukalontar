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
		"Taru Sotong": "Psidium guajava",
		"jambu biji":  "Psidium guajava",
		"kalyasĕm":    "Syzygium polycephalum",
		"Kayu Pule":   "Alstonia scholaris",
		"Turi Putih":  "Sesbania grandiflora",
		"Tabia Bun":   "Capsicum annuum",
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
	if _, ok := store.Lookup("Taru Khayalan"); ok {
		t.Fatal("unknown plant must not match")
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
