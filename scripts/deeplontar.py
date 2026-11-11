#!/usr/bin/env python3
"""Download and validate the pinned DeepLontar v2 original scans and labels."""
import argparse
import hashlib
import io
import json
from pathlib import Path
import urllib.request
import zipfile

from PIL import Image

ROOT = Path(__file__).resolve().parent.parent
MANIFEST = ROOT / 'data/deeplontar-manifest.json'
DEST = ROOT / 'data/corpus/deeplontar'
ARCHIVES = {
    'DeepLontar.zip': ('https://ndownloader.figshare.com/files/35962061', 'dbb42b0bc84d7dace57d81c94e0457a1a7176b6bbd57a9dabf5de2ae07461ea8'),
    'DeepLontar_Labels.zip': ('https://ndownloader.figshare.com/files/35962064', '98f580e51f4921cc424d879893d4c8a953e70ae822a263217735f3297ea3be60'),
}


def digest(data):
    return hashlib.sha256(data).hexdigest()


def acquire(samples):
    DEST.mkdir(parents=True, exist_ok=True)
    for name, (url, expected) in ARCHIVES.items():
        path = DEST / name
        if not path.exists() or digest(path.read_bytes()) != expected:
            request = urllib.request.Request(url, headers={'User-Agent': 'BukaLontarCorpus/1.0 (research dataset acquisition)'})
            with urllib.request.urlopen(request, timeout=120) as response:
                data = response.read()
            if digest(data) != expected:
                raise ValueError(f'Archive checksum mismatch: {name}')
            path.write_bytes(data)
    with zipfile.ZipFile(DEST / 'DeepLontar.zip') as images, zipfile.ZipFile(DEST / 'DeepLontar_Labels.zip') as labels:
        for sample in samples:
            for archive, folder, key in ((images, 'images', 'image'), (labels, 'labels', 'label')):
                name = sample[key]
                if Path(name).name != name:
                    raise ValueError(f'Invalid filename: {name}')
                data = archive.read(name)
                if digest(data) != sample[key + '_sha256']:
                    raise ValueError(f'Checksum mismatch in archive: {name}')
                target = DEST / folder / name
                target.parent.mkdir(parents=True, exist_ok=True)
                target.write_bytes(data)


def validate(samples):
    seen = set()
    image_bytes = label_bytes = boxes = quarantined = usable_boxes = 0
    for sample in samples:
        name = sample['image']
        if Path(name).name != name or Path(sample['label']).name != sample['label']:
            raise ValueError('Invalid filename in manifest')
        image = (DEST / 'images' / name).read_bytes()
        label = (DEST / 'labels' / sample['label']).read_bytes()
        if len(image) != sample['image_bytes'] or digest(image) != sample['image_sha256']:
            raise ValueError(f'Invalid image: {name}')
        if len(label) != sample['label_bytes'] or digest(label) != sample['label_sha256']:
            raise ValueError(f'Invalid annotation: {name}')
        if sample['image_sha256'] in seen:
            raise ValueError(f'Duplicate image: {name}')
        seen.add(sample['image_sha256'])
        with Image.open(io.BytesIO(image)) as decoded:
            decoded.load()
            if decoded.format != 'JPEG' or decoded.size != (sample['width'], sample['height']):
                raise ValueError(f'Invalid dimensions or format: {name}')
        lines = label.splitlines()
        if len(lines) != sample['annotations']:
            raise ValueError(f'Invalid annotation count: {name}')
        invalid_classes = 0
        for line in lines:
            fields = line.split()
            if len(fields) != 5 or not all(0 <= float(x) <= 1 for x in fields[1:]):
                raise ValueError(f'Invalid YOLO coordinates: {name}')
            if not 0 <= int(fields[0]) < 55:
                invalid_classes += 1
        if sample.get('label_status') == 'quarantined':
            if not invalid_classes or not sample.get('label_issue'):
                raise ValueError(f'Unexplained annotation quarantine: {name}')
            quarantined += 1
            print(f'QUARANTINED {name}: {invalid_classes} invalid class IDs; {len(lines)} boxes excluded')
        elif invalid_classes:
            raise ValueError(f'Invalid YOLO class: {name}')
        else:
            usable_boxes += len(lines)
        image_bytes += len(image)
        label_bytes += len(label)
        boxes += len(lines)
    if len(samples) != 200:
        raise ValueError('Expected 200 original image faces')
    print(f'{len(samples)} images, {image_bytes} image bytes, {len(samples)} label files, {label_bytes} label bytes, {boxes} character boxes')
    print(f'{len(samples) - quarantined} usable label files, {usable_boxes} usable boxes, {quarantined} quarantined label files (not training-ready)')


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('command', choices=('acquire', 'validate'))
    args = parser.parse_args()
    samples = json.loads(MANIFEST.read_text())['samples']
    if args.command == 'acquire':
        acquire(samples)
    validate(samples)


if __name__ == '__main__':
    main()
