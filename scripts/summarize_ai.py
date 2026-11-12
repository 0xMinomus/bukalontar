#!/usr/bin/env python3
"""Summarize real pipeline runs. Agreement is not accuracy; abstention is not a correct reading."""
import argparse
from collections import defaultdict
from difflib import SequenceMatcher
from itertools import combinations
import json
from pathlib import Path
import re
import unicodedata
import statistics


def reading_text(response):
    text = response.get('paleography', {}).get('transliteration', '').lower()
    text = re.sub(r'\[lembar\s+\d+\]|\[tidak terbaca\]|tidak terbaca\.?|tidak terdeteksi\.?|tidak disebutkan\.?', '', text)
    return ' '.join(re.findall(r'[^\W\d_]+', text, re.UNICODE))


def plant_names(response):
    names = set()
    for recipe in response.get('recipes', []):
        name = recipe.get('plant_local_name', '').strip().lower()
        if not name or name.startswith(('tidak ', '[tidak ', 'bukan ')):
            continue
        names.add(re.sub(r'^taru\s+', '', name))
    return names

def balinese_text(value):
    value = re.sub(r'\[[^\]]*\]', '', unicodedata.normalize('NFC', value))
    return ''.join(char for char in value if '\u1b00' <= char <= '\u1b7f')


def edit_distance(left, right):
    previous = list(range(len(right) + 1))
    for i, a in enumerate(left, 1):
        current = [i]
        for j, b in enumerate(right, 1):
            current.append(min(current[-1] + 1, previous[j] + 1, previous[j - 1] + (a != b)))
        previous = current
    return previous[-1]


def summarize(rows, references):
    grouped = defaultdict(list)
    for row in rows:
        grouped[(row['variant'], row['model'])].append(row)
    reports = []
    for (variant, model), runs in grouped.items():
        successful = [row for row in runs if not row.get('error') and row.get('response')]
        balinese_errors = balinese_characters = balinese_scored = 0
        balinese_raw_cer = []
        for row in runs:
            reference = references.get(row['image_sha256'], {})
            expected = reference.get('expected_balinese_text')
            if expected:
                actual = (row.get('response') or {}).get('paleography', {}).get('balinese_transcription', '')
                normalized = balinese_text(expected)
                if normalized:
                    balinese_errors += edit_distance(normalized, balinese_text(actual))
                    balinese_characters += len(normalized)
                    balinese_scored += 1
                    balinese_raw_cer.append(edit_distance(expected, actual) / len(expected))
        durations = sorted(row['elapsed_ms'] / 1000 for row in runs)
        positive = [row for row in runs if references.get(row['image_sha256'], {}).get('aksara_bali_detected') is True]
        negative = [row for row in runs if references.get(row['image_sha256'], {}).get('aksara_bali_detected') is False]
        by_image = defaultdict(list)
        detection = []
        cer = []
        for row in successful:
            by_image[row['image_sha256']].append(row['response'])
            reference = references.get(row['image_sha256'], {})
            actual = row['response']['paleography']['aksara_bali_detected']
            if 'aksara_bali_detected' in reference:
                detection.append(actual == reference['aksara_bali_detected'])
            expected_text = reference.get('transliteration')
            if expected_text:
                actual_text = row['response']['paleography']['transliteration']
                cer.append(edit_distance(expected_text, actual_text) / len(expected_text))
        agreements = []
        recipe_agreements = []
        for responses in by_image.values():
            for a, b in combinations(responses, 2):
                left, right = reading_text(a), reading_text(b)
                if left and right:
                    agreements.append(SequenceMatcher(None, left, right, autojunk=False).ratio())
                sa, sb = plant_names(a), plant_names(b)
                if sa or sb:
                    recipe_agreements.append(len(sa & sb) / len(sa | sb))
        detected = [row for row in successful if row['response']['paleography']['aksara_bali_detected']]
        reports.append({
            'variant': variant, 'model': model, 'runs': len(runs), 'successful': len(successful),
            'median_seconds': round(statistics.median(durations), 3),
            'max_seconds': round(max(durations), 3),
            'detection_correct': sum(detection), 'detection_scored': len(detection),
            'expected_balinese_runs': len(positive),
            'accepted_balinese_runs': sum(not row.get('error') and (row.get('response') or {}).get('paleography', {}).get('aksara_bali_detected') is True for row in positive),
            'expected_non_balinese_runs': len(negative),
            'rejected_non_balinese_runs': sum(not row.get('error') and (row.get('response') or {}).get('paleography', {}).get('aksara_bali_detected') is False for row in negative),
            'detected_runs_with_no_readable_text': sum(not reading_text(row['response']) for row in detected),
            'detected_runs': len(detected),
            'transliteration_pair_agreement': round(statistics.mean(agreements), 4) if agreements else None,
            'transliteration_pairs_scored': len(agreements),
            'plant_set_pair_jaccard': round(statistics.mean(recipe_agreements), 4) if recipe_agreements else None,
            'plant_set_pairs_scored': len(recipe_agreements),
            'mean_character_error_rate': round(statistics.mean(cer), 4) if cer else None,
            'transcription_references_scored': len(cer),
            'balinese_character_error_rate': round(balinese_errors / balinese_characters, 4) if balinese_characters else None,
            'balinese_reference_runs': balinese_scored,
            'balinese_raw_mean_character_error_rate': round(statistics.mean(balinese_raw_cer), 4) if balinese_raw_cer else None,
        })
    return reports


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('results', nargs='+', type=Path)
    parser.add_argument('--references', type=Path, help='JSON object keyed by image SHA-256; optional human transliteration and expected aksara_bali_detected')
    args = parser.parse_args()
    rows = [json.loads(line) for path in args.results for line in path.read_text().splitlines() if line.strip()]
    if args.references:
        references = json.loads(args.references.read_text())
        if isinstance(references, list):
            references = {record['image_sha256']: record for record in references}
    else:
        references = {}
    print(json.dumps(summarize(rows, references), ensure_ascii=False, indent=2))


if __name__ == '__main__':
    main()
