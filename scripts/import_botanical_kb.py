"""Build the botanical lookup catalog from a pinned Wikidata/GBIF facts snapshot.

No medical claims are imported. Refreshing remote facts is a separate research step;
this command reproduces the checked-in catalog without network access.
"""

import collections
import copy
import json
from pathlib import Path
import unicodedata
from urllib.parse import quote

ROOT = Path(__file__).resolve().parents[1]
SOURCE = ROOT / "data/botanical-sources/wikidata-gbif.json"
BASE = ROOT / "apps/api/internal/kb/data/taru-pramana.json"
OUTPUT = ROOT / "apps/api/internal/kb/data/wikidata-plants.json"
REPORT = ROOT / "data/botanical-sources/import-report.json"

# These names already have documented ambiguity in the manuscript research.
AMBIGUOUS_NAMES = {"daun dewa", "leunca", "bayam", "kemangi", "gedang", "kedongdong", "tutup bumi"}
LETTERS = str.maketrans({"ĕ": "e", "ə": "e", "ā": "a", "ī": "i", "ū": "u", "ṭ": "t", "ḍ": "d", "ṇ": "n", "ś": "s", "ḥ": "h", "ŋ": "ng", "ṅ": "ng", "ŕ": "r"})


def normalize(name):
    text = name.strip().lower().translate(LETTERS)
    text = "".join(c for c in text if unicodedata.category(c)[0] in {"L", "N"} or c.isspace())
    text = " ".join(text.split()).removeprefix("taru ")
    return "" if text == "taru" else text


def accepted_record(record):
    match = record.get("match", {})
    accepted = record.get("accepted", match)
    if (
        record.get("error")
        or match.get("matchType") != "EXACT"
        or match.get("rank") != "SPECIES"
        or accepted.get("rank") != "SPECIES"
        or accepted.get("kingdom") != "Plantae"
        or accepted.get("taxonomicStatus", accepted.get("status")) != "ACCEPTED"
        or not accepted.get("canonicalName")
        or not accepted.get("family")
    ):
        return None
    key = match.get("acceptedUsageKey", match.get("usageKey"))
    return (key, accepted) if isinstance(key, int) else None


def build(snapshot, base):
    base = copy.deepcopy(base)
    groups = {}
    rejected = []
    for name, record in sorted(snapshot["records"].items()):
        result = accepted_record(record)
        if result is None:
            rejected.append(name)
            continue
        key, accepted = result
        group = groups.setdefault(key, {"accepted": accepted, "sources": []})
        group["sources"].append((name, record))

    existing = {}
    for entry in base:
        record = snapshot["records"].get(entry["scientific_name"], {})
        result = accepted_record(record)
        if result:
            if result[0] in existing:
                raise ValueError(f"curated entries identify the same GBIF species: {entry['id']}")
            existing[result[0]] = entry

    imported = []
    for key, group in sorted(groups.items()):
        accepted = group["accepted"]
        canonical = accepted["canonicalName"]
        gbif_url = f"https://www.gbif.org/species/{key}"
        aliases = [canonical]
        notes = [f"GBIF Backbone Taxonomy: {canonical}, {accepted['family']} ({gbif_url})."]
        for name, record in group["sources"]:
            aliases.append(name)
            if record.get("wikidata"):
                aliases.extend([record.get("label"), *record.get("aliases", [])])
                url = f"https://www.wikidata.org/w/index.php?title={record['wikidata']}&oldid={record['revision']}"
                notes.append(f"Nama dan alias menurut Wikidata, revisi {record['revision']} ({url}).")
        if key in existing:
            entry = existing[key]
            entry["local_names"].extend(a for a in aliases if a)
            for note in notes:
                if note not in entry["reference_studies"]:
                    entry["reference_studies"].append(note)
            continue
        if not any(record.get("wikidata") for _, record in group["sources"]):
            continue
        names = [a for a in aliases if isinstance(a, str) and a.strip()]
        # A scientific-only entry remains usable by its scientific name; no local name is invented.
        display = next((r.get("label") for _, r in group["sources"] if r.get("label") and normalize(r["label"]) != normalize(canonical)), canonical)
        imported.append({
            "id": f"gbif-{key}",
            "local_names": [display, *names],
            "scientific_name": canonical,
            "botanical_family": accepted["family"],
            "energy_character": "Tidak disebutkan",
            "plant_part_used": "Tidak disebutkan",
            "phytochemical_compounds": [],
            "modern_pharmacological_action": "Identitas takson berasal dari GBIF; nama lokal dari Wikidata belum diverifikasi terhadap lempir. Khasiat dan keamanan racikan belum diverifikasi.",
            "reference_studies": [*dict.fromkeys(notes), "Rujukan identitas botani, bukan bukti klinis racikan naskah."],
            "reference_links": [{"title": f"PubMed: pencarian {canonical}", "url": "https://pubmed.ncbi.nlm.nih.gov/?term=" + quote(canonical, safe=""), "kind": "literature_search"}],
        })

    all_entries = [*base, *imported]
    owners = collections.defaultdict(set)
    for entry in all_entries:
        for name in [entry["scientific_name"], *entry["local_names"]]:
            owners[normalize(name)].add(entry["id"])
    ambiguous = {name for name, ids in owners.items() if len(ids) > 1} | AMBIGUOUS_NAMES
    for entry in all_entries:
        names = []
        seen = set()
        for name in entry["local_names"]:
            key = normalize(name)
            if not key or key in ambiguous or key in seen:
                continue
            seen.add(key)
            names.append(name.strip())
        scientific = normalize(entry["scientific_name"])
        if scientific in ambiguous:
            raise ValueError(f"scientific name collision: {entry['scientific_name']}")
        entry["local_names"] = names or [entry["scientific_name"]]
    report = {
        "source_date": snapshot["retrieved_at"],
        "curated_entries": len(base),
        "imported_entries": len(imported),
        "total_entries": len(all_entries),
        "ambiguous_aliases_excluded": sorted(ambiguous),
        "rejected_scientific_names": rejected,
    }
    return base, imported, report


def main():
    snapshot = json.loads(SOURCE.read_text())
    base = snapshot["curated_entries"]
    base, imported, report = build(snapshot, base)
    if report["total_entries"] < 1001:
        raise ValueError(f"insufficient distinct botanical entries: {report['total_entries']}")
    for path, value in [(BASE, base), (OUTPUT, imported), (REPORT, report)]:
        path.write_text(json.dumps(value, ensure_ascii=False, indent=2) + "\n")
    print(json.dumps({k: v for k, v in report.items() if not isinstance(v, list)}, ensure_ascii=False))


if __name__ == "__main__":
    main()
