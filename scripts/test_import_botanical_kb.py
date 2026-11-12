import unittest

from import_botanical_kb import accepted_record, build, normalize


def source(name, key, label):
    return {
        "wikidata": f"Q{key}",
        "revision": 1,
        "label": label,
        "aliases": [],
        "match": {
            "usageKey": key,
            "matchType": "EXACT",
            "rank": "SPECIES",
            "kingdom": "Plantae",
            "status": "ACCEPTED",
            "canonicalName": name,
            "family": "Solanaceae",
        },
    }


class BotanicalImportTests(unittest.TestCase):
    def test_shared_local_name_does_not_identify_either_species(self):
        snapshot = {
            "retrieved_at": "2026-11-12",
            "records": {
                "Solanum nigrum": source("Solanum nigrum", 1, "Leunca"),
                "Solanum americanum": source("Solanum americanum", 2, "Leunca"),
            },
        }
        _, entries, report = build(snapshot, [])
        self.assertEqual(len(entries), 2)
        self.assertIn("leunca", report["ambiguous_aliases_excluded"])
        for entry in entries:
            self.assertNotIn("leunca", [normalize(n) for n in entry["local_names"]])
            self.assertIn(entry["scientific_name"], entry["local_names"])

    def test_synonyms_do_not_inflate_species_count(self):
        record = source("Solanum lycopersicum", 3, "Tomat")
        synonym = source("Lycopersicon esculentum", 4, "Tomat")
        synonym["match"].update(status="SYNONYM", acceptedUsageKey=3)
        synonym["accepted"] = {
            **record["match"], "taxonomicStatus": "ACCEPTED"
        }
        _, entries, _ = build({
            "retrieved_at": "2026-11-12",
            "records": {"Solanum lycopersicum": record, "Lycopersicon esculentum": synonym},
        }, [])
        self.assertEqual(len(entries), 1)
        self.assertEqual(entries[0]["scientific_name"], "Solanum lycopersicum")
        self.assertIn("Lycopersicon esculentum", entries[0]["local_names"])
        self.assertEqual(entries[0]["phytochemical_compounds"], [])

    def test_fuzzy_genus_and_nonplant_matches_are_rejected(self):
        for field, value in [("matchType", "FUZZY"), ("rank", "GENUS"), ("kingdom", "Animalia")]:
            with self.subTest(field=field):
                record = source("Solanum nigrum", 1, "Leunca")
                record["match"][field] = value
                self.assertIsNone(accepted_record(record))

    def test_rebuild_does_not_restore_a_previously_ambiguous_alias(self):
        base = [{
            "id": "curated-nigrum",
            "local_names": ["Leunca", "Ranti"],
            "scientific_name": "Solanum nigrum",
            "botanical_family": "Solanaceae",
            "reference_studies": [],
        }]
        snapshot = {
            "retrieved_at": "2026-11-12",
            "records": {
                "Solanum nigrum": source("Solanum nigrum", 1, "Ranti"),
                "Solanum americanum": source("Solanum americanum", 2, "Leunca"),
            },
        }
        first = build(snapshot, base)
        second = build(snapshot, base)
        self.assertEqual(first, second)
        self.assertIn("Leunca", base[0]["local_names"])
        for entry in [*second[0], *second[1]]:
            self.assertNotIn("leunca", [normalize(n) for n in entry["local_names"]])


if __name__ == "__main__":
    unittest.main()
