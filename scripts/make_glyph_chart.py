#!/usr/bin/env python3
"""Build one labelled contact sheet of every expert glyph exemplar.

The chart gives the multimodal model the published Balinese glyph inventory in a
single image. Labels come from DeepLontar Table 1 (CC BY 4.0); the chart is not
a transcription of any manuscript.
"""
import json
from pathlib import Path

from PIL import Image, ImageDraw, ImageFont

ASSETS = Path(__file__).resolve().parents[1] / 'apps/api/internal/readingexamples/assets'
OUT = ASSETS / 'glyph-chart.jpg'
COLUMNS = 8
CELL_W, CELL_H = 160, 148
PAD = 10


def font(size):
    for candidate in (
        '/usr/share/fonts/dejavu-sans-fonts/DejaVuSans.ttf',
        '/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf',
    ):
        if Path(candidate).exists():
            return ImageFont.truetype(candidate, size)
    return ImageFont.load_default(size=size)


def main():
    catalog = json.loads((ASSETS / 'catalog.json').read_text())
    rows = (len(catalog) + COLUMNS - 1) // COLUMNS
    chart = Image.new('RGB', (COLUMNS * CELL_W + PAD * 2, rows * CELL_H + PAD * 2 + 40), 'white')
    draw = ImageDraw.Draw(chart)
    title_font, label_font = font(20), font(13)
    draw.text((PAD, PAD), 'Balinese lontar glyph inventory (DeepLontar Table 1, CC BY 4.0); labels are not target text', fill='black', font=title_font)
    for index, record in enumerate(catalog):
        column, row = index % COLUMNS, index // COLUMNS
        x, y = PAD + column * CELL_W, PAD + 40 + row * CELL_H
        glyph = Image.open(ASSETS / record['file']).convert('L')
        glyph.thumbnail((CELL_W - 20, CELL_H - 50))
        chart.paste(glyph.convert('RGB'), (x + (CELL_W - glyph.width) // 2, y + (CELL_H - 46 - glyph.height) // 2))
        label = record['label']
        width = draw.textlength(label, font=label_font)
        draw.text((x + (CELL_W - width) / 2, y + CELL_H - 36), label, fill='black', font=label_font)
    chart.save(OUT, quality=88, optimize=True)
    print(f'{len(catalog)} glyphs -> {OUT} ({OUT.stat().st_size} bytes)')


if __name__ == '__main__':
    main()
