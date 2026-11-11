#!/usr/bin/env python3
"""Acquire the labelled glyph illustrations from DeepLontar's CC BY 4.0 paper.

These are isolated character images, not page transcriptions. The article does
not identify the original leaf for each illustration, so all are development
examples and none may be treated as leaf-independent evaluation cases.
"""
import concurrent.futures
import hashlib
import json
import re
import urllib.request
import sys
import urllib.parse
import xml.etree.ElementTree as ET
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DEST = ROOT / 'apps/api/internal/readingexamples/assets'
SOURCE = 'https://pmc.ncbi.nlm.nih.gov/articles/PMC9741579/'
XML = 'https://www.ebi.ac.uk/europepmc/webservices/rest/PMC9741579/fullTextXML'


def fetch(url):
    with urllib.request.urlopen(urllib.request.Request(url, headers={'User-Agent': 'BukaLontar research corpus/1.0'}), timeout=45) as response:
        return response.read()


def main():
    html = fetch(SOURCE).decode('utf-8')
    xml = ET.fromstring(fetch(XML))
    table = xml.find(".//table-wrap[@id='Tab1']")
    entries = []
    for cell in table.findall('.//tbody/tr/td'):
        graphic = cell.find('inline-graphic')
        if graphic is None:
            continue
        filename = graphic.attrib['{http://www.w3.org/1999/xlink}href']
        match = re.search(r'src="(https://cdn\.ncbi\.nlm\.nih\.gov/pmc/blobs/[^" ]+/' + re.escape(filename) + r')"', html)
        if not match:
            raise ValueError(f'Missing image URL: {filename}')
        entries.append({'file': filename, 'label': ''.join(cell.itertext()).strip(), 'url': match.group(1)})
    if len(entries) != 55:
        raise ValueError(f'Expected 55 labelled classes, got {len(entries)}')
    DEST.mkdir(parents=True, exist_ok=True)
    with concurrent.futures.ThreadPoolExecutor(max_workers=8) as pool:
        images = list(pool.map(fetch, (entry['url'] for entry in entries)))
    for entry, image in zip(entries, images):
        if not image.startswith((b'GIF87a', b'GIF89a')):
            raise ValueError(f'Expected GIF: {entry["url"]}')
        (DEST / entry['file']).write_bytes(image)
        entry['sha256'] = hashlib.sha256(image).hexdigest()
    (DEST / 'catalog.json').write_text(json.dumps(entries, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
    print(f'Acquired {len(entries)} labelled handwritten character illustrations')


def acquire_holdout():
    dest = ROOT / 'data/reading-examples/holdout'
    dest.mkdir(parents=True, exist_ok=True)
    revisions = {4: 225, 5: 232, 6: 240, 7: 247, 8: 254}
    cases = []
    for page, revision in revisions.items():
        params = urllib.parse.urlencode({
            'action': 'query', 'prop': 'revisions', 'revids': revision,
            'rvprop': 'ids|content|user', 'rvslots': 'main', 'format': 'json',
        })
        payload = json.loads(fetch('https://ban.wikisource.org/w/api.php?' + params))
        page_data = next(iter(payload['query']['pages'].values()))
        rev = page_data['revisions'][0]
        raw = rev['slots']['main']['*']
        if rev['revid'] != revision or '<pagequality level="4" user="Carma citrawati"' not in raw:
            raise ValueError(f'Unexpected proofread revision {revision}')
        text = re.sub(r'<noinclude>.*?</noinclude>', '', raw, flags=re.DOTALL)
        text = re.sub(r'<br\s*/?>', '', text).strip() + '\n'
        if '<' in text or '{{' in text:
            raise ValueError(f'Unexpected markup in revision {revision}')
        text_file = f'aji-kreket-page-{page}.txt'
        image_file = f'aji-kreket-page-{page}.jpg'
        image_url = f'https://thumb.wikimedia.org/wikipedia/commons/thumb/a/a3/Bali-lontar-aji-kreket-550ppi.pdf/page{page}-1280px-Bali-lontar-aji-kreket-550ppi.pdf.jpg'
        image = fetch(image_url)
        if not image.startswith(b'\xff\xd8\xff'):
            raise ValueError(f'Unexpected image format for page {page}')
        (dest / image_file).write_bytes(image)
        (dest / text_file).write_text(text, encoding='utf-8')
        cases.append({
            'id': f'aji-kreket-page-{page}', 'group': 'aji-kreket', 'manuscript_id': 'aji-kreket',
            'split': 'holdout', 'granularity': 'two-leaf-faces-per-scan', 'quality': 'Wikisource validated',
            'image': image_file, 'image_sha256': hashlib.sha256(image).hexdigest(),
            'image_url': image_url, 'image_license': 'Wikimedia Commons asserts PD-ID-exempt',
            'image_source': 'https://commons.wikimedia.org/wiki/File:Bali-lontar-aji-kreket-550ppi.pdf',
            'reference': text_file, 'reference_sha256': hashlib.sha256(text.encode()).hexdigest(),
            'expected_balinese_text': text,
            'reference_script': 'Balinese', 'reference_revision': revision,
            'reference_source': f'https://ban.wikisource.org/w/index.php?oldid={revision}',
            'reference_license': 'CC BY-SA 4.0', 'reference_validator': 'Carma citrawati',
            'reference_history': f'https://ban.wikisource.org/w/index.php?title=Kaca:Bali-lontar-aji-kreket-550ppi.pdf/{page}&action=history',
        })
    (dest / 'manifest.json').write_text(json.dumps(cases, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
    print(f'Acquired {len(cases)} human-validated held-out scan/transcription pairs')


if __name__ == '__main__':
    acquire_holdout() if sys.argv[1:] == ['--heldout'] else main()
