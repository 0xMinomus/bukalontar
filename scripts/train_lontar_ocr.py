#!/usr/bin/env python3
"""Fit tile and glyph models; evaluate only on held-out leaf groups.

Requires numpy, scipy, Pillow and scikit-learn. Run from repository root:
OPENBLAS_NUM_THREADS=1 OMP_NUM_THREADS=1 python scripts/train_lontar_ocr.py
"""
import hashlib
import json
from collections import Counter
from pathlib import Path

import numpy as np
from PIL import Image
from scipy import ndimage
from sklearn.svm import LinearSVC

ROOT = Path(__file__).resolve().parents[1]
CORPUS = ROOT / 'data/corpus/deeplontar'
OUT = ROOT / 'apps/api/internal/ocr/model.json'
REPORT = ROOT / 'data/ocr/evaluation.json'
N = 16


def feature(im, box):
    x0, y0, x1, y1 = box
    if x1 <= x0 or y1 <= y0:
        return None
    # Match Go's nearest-neighbor pixel-center sampling exactly.
    xs = np.minimum(im.shape[1]-1, x0 + ((np.arange(N)*2+1)*(x1-x0)//(2*N)))
    ys = np.minimum(im.shape[0]-1, y0 + ((np.arange(N)*2+1)*(y1-y0)//(2*N)))
    f = im[np.ix_(ys, xs)].astype(np.float32).reshape(-1)
    f -= f.mean()
    norm = np.linalg.norm(f)
    return f / norm if norm > 1e-6 else None


def load():
    splits = {'train': [], 'dev': [], 'test': []}
    groups = {'train': [], 'dev': [], 'test': []}
    counts = {k: Counter() for k in splits}
    # Fixed arithmetic split of physical leaf numbers; both faces stay together.
    for leaf in range(1, 101):
        split = 'test' if leaf % 10 == 0 else 'dev' if leaf % 10 == 1 else 'train'
        groups[split].append(leaf)
        for face in 'ab':
            stem = f'{leaf}{face}'
            if stem == '76a':
                continue
            im = np.asarray(Image.open(CORPUS / 'images' / (stem+'.jpg')).convert('RGB'))
            # Integer BT.601 grayscale; same formula as Go.
            im = ((299*im[:,:,0].astype(np.uint32)+587*im[:,:,1].astype(np.uint32)+114*im[:,:,2].astype(np.uint32)+500)//1000).astype(np.uint8)
            h, w = im.shape
            for line in (CORPUS / 'labels' / (stem+'.txt')).read_text().splitlines():
                c, x, y, bw, bh = line.split()
                c = int(c)
                if not 0 <= c < 55:
                    raise ValueError(f'invalid class {c} in {stem}')
                x, y, bw, bh = map(float, (x, y, bw, bh))
                box = (max(0, int((x-bw/2)*w)), max(0, int((y-bh/2)*h)), min(w, int((x+bw/2)*w)), min(h, int((y+bh/2)*h)))
                f = feature(im, box)
                if f is not None:
                    splits[split].append((c, f))
                    counts[split][c] += 1
    return splits, groups, counts


def page_proposals(leaves, means):
    records = []
    for leaf in leaves:
        for face in 'ab':
            stem = f'{leaf}{face}'
            if stem == '76a':
                continue
            rgb = np.asarray(Image.open(CORPUS / 'images' / (stem+'.jpg')).convert('RGB')).astype(np.uint32)
            im = ((299*rgb[:,:,0]+587*rgb[:,:,1]+114*rgb[:,:,2]+500)//1000).astype(np.uint8)
            h, w = im.shape
            truth = []
            for line in (CORPUS / 'labels' / (stem+'.txt')).read_text().splitlines():
                c, x, y, bw, bh = line.split()
                x, y, bw, bh = map(float, (x, y, bw, bh))
                truth.append((int(c), int((x-bw/2)*w), int((y-bh/2)*h), int((x+bw/2)*w), int((y+bh/2)*h)))
            local = ndimage.uniform_filter(im.astype(np.float32), size=31, mode='nearest')
            labels, _ = ndimage.label(im.astype(np.float32) < local - 25, structure=np.ones((3,3)))
            for component, s in enumerate(ndimage.find_objects(labels), 1):
                if s is None:
                    continue
                y0,y1=s[0].start,s[0].stop
                x0,x1=s[1].start,s[1].stop
                if not (8 <= x1-x0 <= 80 and 8 <= y1-y0 <= 80):
                    continue
                f = feature(im, (x0,y0,x1,y1))
                if f is None:
                    continue
                scores = means @ f
                order = np.argsort(scores)
                pred = int(order[-1])
                margin = float(scores[order[-1]]-scores[order[-2]])
                best_iou, matched, match_index = 0.0, -1, -1
                for index, (c,tx0,ty0,tx1,ty1) in enumerate(truth):
                    intersection = max(0,min(x1,tx1)-max(x0,tx0))*max(0,min(y1,ty1)-max(y0,ty0))
                    union = (x1-x0)*(y1-y0)+(tx1-tx0)*(ty1-ty0)-intersection
                    iou = intersection / union if union else 0
                    if iou > best_iou:
                        best_iou, matched, match_index = iou, c, index
                patch = labels[y0:y1,x0:x1]
                occupancy = float((patch == component).mean())
                records.append((margin, best_iou >= 0.3 and pred == matched, best_iou >= 0.3,
                                (x1-x0)/80, (y1-y0)/80, occupancy,
                                float(scores[order[-1]]), (x0+x1)/(2*w), (y0+y1)/(2*h),
                                leaf*10000 + (0 if face == 'a' else 5000) + match_index if best_iou >= .3 else -1))
    return np.array(records, dtype=float)


def row_data(leaves):
    features, targets = [], []
    for leaf in leaves:
        for face in 'ab':
            stem = f'{leaf}{face}'
            if stem == '76a':
                continue
            rgb = np.asarray(Image.open(CORPUS / 'images' / (stem+'.jpg')).convert('RGB')).astype(np.uint32)
            im = ((299*rgb[:,:,0]+587*rgb[:,:,1]+114*rgb[:,:,2]+500)//1000).astype(np.uint8)
            h, w = im.shape
            local = ndimage.uniform_filter(im.astype(np.float32), size=31, mode='nearest')
            mask = (im.astype(np.float32) < local - 25)
            density = mask.reshape(10,30,30,50).mean(axis=(1,3))
            target = np.zeros((10,30), dtype=float)
            for line in (CORPUS / 'labels' / (stem+'.txt')).read_text().splitlines():
                _, x, y, _, _ = line.split()
                col, row = min(29, int(float(x)*30)), min(9, int(float(y)*10))
                target[row,col] = 1
            padded = np.pad(density, 1, mode='edge')
            for row in range(10):
                for col in range(30):
                    features.append([1., *padded[row:row+3,col:col+3].reshape(-1)])
            targets.extend(target.reshape(-1))
    return np.array(features), np.array(targets)

def hog_features(items):
    images = np.stack([f for _, f in items]).reshape(-1,16,16)
    gx = np.zeros_like(images)
    gy = np.zeros_like(images)
    gx[:,:,1:-1] = images[:,:,2:] - images[:,:,:-2]
    gy[:,1:-1,:] = images[:,2:,:] - images[:,:-2,:]
    magnitude = np.sqrt(gx*gx+gy*gy)
    direction = np.floor((np.arctan2(gy,gx) % np.pi) * (8/np.pi) + 0.0001).astype(int).clip(0,7)
    out = np.zeros((len(items),4,4,8), dtype=np.float32)
    for cy in range(4):
        for cx in range(4):
            bins = direction[:,cy*4:cy*4+4,cx*4:cx*4+4].reshape(-1,16)
            mags = magnitude[:,cy*4:cy*4+4,cx*4:cx*4+4].reshape(-1,16)
            for b in range(8):
                out[:,cy,cx,b] = (mags*(bins==b)).sum(axis=1)
    out = out.reshape(len(items),-1)
    out /= np.maximum(np.linalg.norm(out,axis=1,keepdims=True), 1e-8)
    return out



def main():
    splits, groups, counts = load()
    means = np.zeros((55, N*N), dtype=np.float32)
    for c, f in splits['train']:
        means[c] += f
    norms = np.linalg.norm(means, axis=1)
    means /= np.maximum(norms[:,None], 1e-8)
    def score(items):
        y = np.array([c for c, _ in items])
        x = np.stack([f for _, f in items])
        sims = x @ means.T
        best = np.argmax(sims, axis=1)
        top = np.max(sims, axis=1)
        second = np.partition(sims, -2, axis=1)[:,-2]
        return y, best, top-second
    yd, pd, md = score(splits['dev'])
    # A conservative global threshold: at least 50 dev boxes and >=80% precision.
    thresholds = np.linspace(0, .5, 501)
    valid = [float(t) for t in thresholds if (md >= t).sum() >= 50 and (yd[md >= t] == pd[md >= t]).mean() >= .8]
    threshold = min(valid) if valid else None
    def metrics(split):
        y, p, margin = score(splits[split])
        accepted = margin >= threshold if threshold is not None else np.zeros(len(y), dtype=bool)
        return {'boxes': len(y), 'top1_accuracy': float((y==p).mean()), 'accepted': int(accepted.sum()), 'accepted_coverage': float(accepted.mean()), 'accepted_accuracy': float((y[accepted]==p[accepted]).mean()) if accepted.any() else None, 'class_counts': dict(sorted(counts[split].items()))}
    page_dev = page_proposals(groups['dev'], means)
    page_test = page_proposals(groups['test'], means)
    page_valid = [float(t) for t in thresholds if (page_dev[:,0] >= t).sum() >= 50 and page_dev[page_dev[:,0] >= t,1].mean() >= .8]
    page_threshold = min(page_valid) if page_valid else None
    location_valid = [float(t) for t in thresholds if (page_dev[:,0] >= t).sum() >= 50 and page_dev[page_dev[:,0] >= t,2].mean() >= .8]
    location_threshold = min(location_valid) if location_valid else None
    def location_metrics(records):
        accepted = records[:,0] >= location_threshold if location_threshold is not None else np.zeros(len(records), dtype=bool)
        return {'accepted': int(accepted.sum()), 'precision': float(records[accepted,2].mean()) if accepted.any() else None}
    def page_metrics(records):
        accepted = records[:,0] >= page_threshold if page_threshold is not None else np.zeros(len(records), dtype=bool)
        return {'proposals': len(records), 'proposal_localization_precision_iou_0_3': float(records[:,2].mean()), 'proposal_class_and_localization_precision_iou_0_3': float(records[:,1].mean()), 'accepted': int(accepted.sum()), 'accepted_precision': float(records[accepted,1].mean()) if accepted.any() else None}
    x_train, y_train = row_data(groups['train'])
    x_dev, y_dev = row_data(groups['dev'])
    x_test, y_test = row_data(groups['test'])
    weights = np.linalg.solve(x_train.T @ x_train + np.eye(x_train.shape[1])*.001, x_train.T @ y_train)
    def row_metrics(x, y, t):
        predicted = x @ weights >= t
        tp = int((predicted & (y == 1)).sum())
        fp = int((predicted & (y == 0)).sum())
        fn = int((~predicted & (y == 1)).sum())
        return {'tiles': len(y), 'tp': tp, 'fp': fp, 'fn': fn, 'precision': tp/(tp+fp) if tp+fp else 0, 'recall': tp/(tp+fn) if tp+fn else 0, 'f1': 2*tp/(2*tp+fp+fn) if 2*tp+fp+fn else 0}
    row_threshold = max(np.linspace(0,1,201), key=lambda t: row_metrics(x_dev,y_dev,t)['f1'])
    def focus_metrics(x, y):
        scores = (x @ weights).reshape(-1,10,30)
        targets = y.reshape(-1,10,30)
        selected = np.zeros_like(targets, dtype=bool)
        pages = 0
        for page, grid in enumerate(scores):
            windows = []
            for yy in range(8):
                for xx in range(16):
                    patch = grid[yy:yy+3,xx:xx+15]
                    good = patch >= row_threshold
                    if good.sum() >= 5:
                        windows.append((float((patch*good).sum()/45), xx, yy))
            windows.sort(key=lambda item:-item[0])
            chosen = []
            for _, xx, yy in windows:
                if any(xx < px+15 and px < xx+15 and yy < py+3 and py < yy+3 for px,py in chosen):
                    continue
                chosen.append((xx,yy))
                selected[page,max(0,yy-1):min(10,yy+4),xx:xx+15] = True
                if len(chosen) == 2:
                    break
            pages += bool(chosen)
        hits = int((selected & (targets == 1)).sum())
        positive = int(targets.sum())
        area = int(selected.sum())
        return {'pages':len(scores), 'pages_with_focus':pages, 'selected_tiles':area, 'selected_tiles_with_character_center':hits, 'selected_area_precision':hits/area if area else 0, 'character_occupied_tile_recall':hits/positive if positive else 0, 'selected_area_fraction':area/targets.size}
    hx_train = hog_features(splits['train'])
    hx_dev = hog_features(splits['dev'])
    hy_train = np.array([c for c,_ in splits['train']])
    hy_dev = np.array([c for c,_ in splits['dev']])
    trials = []
    for balance in (None, 'balanced'):
        clf = LinearSVC(C=1.0, class_weight=balance, random_state=42, max_iter=5000)
        clf.fit(hx_train, hy_train)
        trials.append((float((clf.predict(hx_dev)==hy_dev).mean()), balance, clf))
    _, balance, clf = max(trials, key=lambda item:item[0])
    hy_test = np.array([c for c,_ in splits['test']])
    hp_test = clf.predict(hog_features(splits['test']))
    hog_model = {'version':'deeplontar-hog-linearsvc-1', 'classes':clf.classes_.tolist(), 'weights':clf.coef_.tolist(), 'bias':clf.intercept_.tolist(), 'class_weight':balance}
    hog_path = ROOT / 'apps/api/internal/ocr/hog_model.json'
    hog_path.write_text(json.dumps(hog_model, separators=(',',':'))+'\n')
    model = {'version': 'deeplontar-tile-localizer-1', 'dimension': N, 'threshold': page_threshold, 'location_threshold': location_threshold, 'centroids': means.tolist(), 'row_weights': weights.tolist(), 'row_threshold': float(row_threshold)}
    OUT.write_text(json.dumps(model, separators=(',', ':'))+'\n')
    report = {'source': 'https://doi.org/10.6084/m9.figshare.20103803.v2', 'source_publication': 'https://doi.org/10.1038/s41597-022-01867-5', 'license': 'https://creativecommons.org/licenses/by/4.0/', 'citation': 'Sutramiani, Ni Putu; Suciati, Nanik; Siahaan, Daniel Oranova; Duija, I Nengah; Darma, I Wayan Agus Surya (2022). DeepLontar. figshare. Dataset. https://doi.org/10.6084/m9.figshare.20103803.v2', 'mapping': 'No authoritative class-ID to character-name mapping supplied in dataset archives or paper. IDs are opaque, not text.', 'method': 'mean-centered 16x16 nearest-pixel crops; normalized per-box vectors; class-average cosine classifier; minimum top-two score margin chosen on dev for >=80% precision on >=50 boxes', 'groups': groups, 'excluded': {'76a': 'contains out-of-range class IDs'}, 'image_counts': {'train': 159, 'dev': 20, 'test': 20}, 'model_sha256': hashlib.sha256(OUT.read_bytes()).hexdigest(), 'margin_threshold': threshold, 'crop_metrics': {s: metrics(s) for s in ('dev','test')}, 'limitation': 'Gold-box crop classification only; this does not measure page detection, transcription, translation, or medical accuracy.'}
    report['page_proposals'] = {'method': '31x31 local mean minus 25 grayscale, 8-connected components with dimensions 8..80 pixels; >=0.3 IoU to any same-class expert box; all other proposals including background are errors', 'threshold': page_threshold, 'dev': page_metrics(page_dev), 'test': page_metrics(page_test)}
    report['page_proposals']['location_only'] = {'threshold': location_threshold, 'dev': location_metrics(page_dev), 'test': location_metrics(page_test)}
    report['tile_localization'] = {'method': 'ridge regression fitted to bbox-center occupancy on 30x10 image tiles using neighboring 3x3 contrast-mask densities; development F1 threshold; not glyph recognition', 'threshold': float(row_threshold), 'train': row_metrics(x_train,y_train,row_threshold), 'dev': row_metrics(x_dev,y_dev,row_threshold), 'test': row_metrics(x_test,y_test,row_threshold)}
    report['hog_linear_svc'] = {'method':'16x16 grayscale crops, 4x4 cells, 8 unsigned gradient bins, L2 normalized; LinearSVC C=1, class balance selected by dev top-1', 'dev_trials': [{'class_weight': b, 'accuracy': a} for a,b,_ in trials], 'selected_class_weight':balance, 'test_crop_top1_accuracy':float((hp_test==hy_test).mean()), 'artifact':'apps/api/internal/ocr/hog_model.json', 'sha256':hashlib.sha256(hog_path.read_bytes()).hexdigest(), 'limitation':'Expert-box crop classification only; opaque class IDs, no page OCR or transliteration.'}
    report['focus_windows'] = {'method':'up to two 15x3 tile windows ranked by positive scores; each returned crop is padded vertically by one tile; requires at least 5 positive tiles; not full-page recognition', 'dev':focus_metrics(x_dev,y_dev), 'test':focus_metrics(x_test,y_test)}
    report['training_class_counts'] = dict(sorted(counts['train'].items()))
    REPORT.write_text(json.dumps(report, indent=2)+'\n')
    print(json.dumps({'threshold': threshold, 'dev': report['crop_metrics']['dev'], 'test': report['crop_metrics']['test'], 'model_sha256': report['model_sha256']}))


if __name__ == '__main__':
    main()
