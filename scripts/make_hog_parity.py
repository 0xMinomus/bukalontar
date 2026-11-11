#!/usr/bin/env python3
"""Regenerate the Go/Python HOG parity fixture from the training pipeline.

The fixture stores the exact 16x16 grayscale samples, the expected opaque class
ID, and the classifier margin computed by the same numpy code that trains the
artifact. Run from the repository root after dependencies are installed.
"""
import json
from pathlib import Path

import numpy as np
from PIL import Image

from train_lontar_ocr import CORPUS, ROOT, feature, hog_features

N = 16
OUT = ROOT / 'data/ocr/hog-parity.json'
MODEL = ROOT / 'apps/api/internal/ocr/hog_model.json'


def main():
    model = json.loads(MODEL.read_text())
    weights = np.array(model['weights'])
    bias = np.array(model['bias'])
    classes = np.array(model['classes'])
    pixels, labels = [], []
    # Deterministic sample of annotated characters from held-out leaves.
    for leaf in (10, 20, 30, 40, 50, 60, 70, 80):
        stem = f'{leaf}a'
        rgb = np.asarray(Image.open(CORPUS / 'images' / (stem + '.jpg')).convert('RGB')).astype(np.uint32)
        gray = ((299 * rgb[:, :, 0] + 587 * rgb[:, :, 1] + 114 * rgb[:, :, 2] + 500) // 1000).astype(np.uint8)
        h, w = gray.shape
        lines = (CORPUS / 'labels' / (stem + '.txt')).read_text().splitlines()
        for index in (0, len(lines) // 2):
            c, x, y, bw, bh = lines[index].split()
            c, x, y, bw, bh = int(c), float(x), float(y), float(bw), float(bh)
            box = (max(0, int((x - bw / 2) * w)), max(0, int((y - bh / 2) * h)), min(w, int((x + bw / 2) * w)), min(h, int((y + bh / 2) * h)))
            xs = np.minimum(w - 1, box[0] + ((np.arange(N) * 2 + 1) * (box[2] - box[0]) // (2 * N)))
            ys = np.minimum(h - 1, box[1] + ((np.arange(N) * 2 + 1) * (box[3] - box[1]) // (2 * N)))
            crop = gray[np.ix_(ys, xs)].astype(np.uint8)
            if feature(gray, box) is None:
                continue
            pixels.append(crop.reshape(-1).tolist())
            labels.append(c)
    features = hog_features([(c, np.asarray(p).astype(np.float32)) for p, c in zip(pixels, labels)])
    scores = features.astype(np.float64) @ weights.T + bias
    order = np.argsort(scores, axis=1)
    best = order[:, -1]
    margins = scores[np.arange(len(pixels)), best] - scores[np.arange(len(pixels)), order[:, -2]]
    cases = [
        {'pixels': p, 'class_id': int(classes[b]), 'margin': float(m)}
        for p, b, m in zip(pixels, best, margins)
    ]
    OUT.write_text(json.dumps(cases, separators=(',', ':')) + '\n')
    print(f'{len(cases)} parity cases written to {OUT}')


if __name__ == '__main__':
    main()
