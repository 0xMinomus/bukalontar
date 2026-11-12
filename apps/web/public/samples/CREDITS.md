# Sample image and transcription credits

## Four selectable examples on current main

These are photographs of Balinese-script palm leaves: two Usadha, one Carcan, and one Kakawin. A Kawi language label does not mean Kawi script. Source identity comes from the cited archive, not Gemini's classification. Community transcriptions are not independent expert certification.

The API matches exact image bytes using SHA-256 to the references embedded in `apps/api/internal/ai/manuscript_references.json`. For a match, Latin text comes from a pinned Wikisource rendering, with editorial whitespace changes only. Gemini translates and extracts from that supplied text on every request. This is explicitly reference-assisted analysis, not new OCR or a cached generated answer. Other image bytes use the experimental photograph-reading path.

| File | Photo credit / license | Original photograph | Transcription revision |
| --- | --- | --- | --- |
| `lontar-usadha-sukawati-depan.jpg` | Eka343 / CC BY-SA 4.0 | [Commons](https://commons.wikimedia.org/wiki/File:Bali-lontar-Sukawati-Usada-07.jpg) | [Pinned reference](https://ban.wikisource.org/w/index.php?title=Kaca:Bali-lontar-Sukawati-Usada-07.jpg&oldid=98029) |
| `lontar-usadha-rare-kemenuh-depan.jpg` | Eka343 / CC BY-SA 4.0 | [Commons](https://commons.wikimedia.org/wiki/File:Bali-lontar-Kemenuh-Usadha_Rare-56.jpg) | [Pinned reference](https://ban.wikisource.org/w/index.php?title=Kaca:Bali-lontar-Kemenuh-Usadha_Rare-56.jpg&oldid=107871) |
| `lontar-carcan-kucing-1b.jpg` | Internet Archive / PD-ID-exempt | [Commons](https://commons.wikimedia.org/wiki/File:Bali-lontar-carcan-kucing-250ppi.pdf) | [Pinned reference](https://ban.wikisource.org/w/index.php?title=Carcan_Kucing&oldid=107176) |
| `lontar-kakawin-bharatayuda-depan.jpg` | Eka343 / CC BY-SA 4.0 | [Commons](https://commons.wikimedia.org/wiki/File:Bali-lontar-Balai_Bahasa_Bali-Kekawin_Bharatayuda_Maarti_01.jpg) | [Pinned reference](https://ban.wikisource.org/w/index.php?title=Kekawin_Bharatayuda_Maarti_(Balai_Bahasa_Prov._Bali)&oldid=134900) |

Photo license: <https://creativecommons.org/licenses/by-sa/4.0/> for the three Eka343 photographs. The Carcan PDF uses the distinct [PD-ID-exempt notice](https://commons.wikimedia.org/wiki/File:Bali-lontar-carcan-kucing-250ppi.pdf#Licensing); do not attribute that scan to Eka343 or apply its photo license to the others.

Transcription attribution: contributors to Balinese Wikisource; the pinned revision and its history identify contributors. Reused community text and editorial spacing are distributed under [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/). Latin rendering and project-added word spacing have not been independently verified by a philologist. The API preserves the supplied reference text; translations and extracted preparations remain generated and require review.

## Crop transformations

Coordinates are `(left, top, right, bottom)`, right/bottom exclusive. Each crop retains one complete front-side leaf and removes the other leaf and surrounding table. JPEG quality is 95 with no chroma subsampling. No glyphs were generated, retouched, or reconstructed. The three source JPEG crops retain their original pixel resolution. Carcan is page 2 of the source PDF, rendered with `pdftoppm -f 2 -singlefile -scale-to 2800 -jpeg`, then cropped to leaf 1B.

| File | Original / rendered dimensions | Crop | Output dimensions | SHA-256 |
| --- | --- | --- | --- | --- |
| `lontar-usadha-sukawati-depan.jpg` | 5184 × 3456 | `(380, 1130, 4950, 1685)` | 4570 × 555 | `ac177a971369a4945c1466d295c8b877b98fbae82b9c915c1a1a7bd82c244d88` |
| `lontar-usadha-rare-kemenuh-depan.jpg` | 5269 × 3513 | `(120, 1250, 5210, 1720)` | 5090 × 470 | `949b8ca0acb7aadc587ce46eb64747a3038b9116a7f6df0430f76bea5cb561ff` |
| `lontar-carcan-kucing-1b.jpg` | 2800 × 1869 | `(275, 510, 2700, 705)` | 2425 × 195 | `195f3fb39d8477ede42935059f0a2309616be1559117a662150740a7d6527567` |
| `lontar-kakawin-bharatayuda-depan.jpg` | 5004 × 3336 | `(170, 1070, 4770, 1460)` | 4600 × 390 | `8de4240adf734bf762d25bad4ff141b570d8cca7fad73ce8513d131ad45a4891` |

## Historical and decorative photographs

The immutable [`v1.1-final-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.1-final-submission) tag has two older selectable crops. They remain in the repository for historical evaluations but are no longer selectable on main.

| File | Author / license | Source |
| --- | --- | --- |
| `lontar-usada-lare-depan.jpg` | Eka343 / CC BY-SA 4.0 | [Commons](https://commons.wikimedia.org/wiki/File:Bali-lontar-Balai_Bahasa_Bali-Usada_Lare.jpg), [transcription revision 105197](https://ban.wikisource.org/w/index.php?oldid=105197) |
| `lontar-usada-rare-154-depan.jpg` | Eka343 / CC BY-SA 4.0 | [Commons](https://commons.wikimedia.org/wiki/File:Bali-lontar-Balai_Bahasa_Bali-Usada_Rare_154.jpg), [transcription revision 105194](https://ban.wikisource.org/w/index.php?oldid=105194) |
| `lontar-wanagiri-usada-3.jpg` | Eka343 / CC BY-SA 4.0 | [Commons](https://commons.wikimedia.org/wiki/File:Bali-lontar-Wanagiri-Usada-3.jpg) |
| `lontar-sukawati-usada-07.jpg` | Eka343 / CC BY-SA 4.0 | [Commons](https://commons.wikimedia.org/wiki/File:Bali-lontar-Sukawati-Usada-07.jpg) |

The old Lare crop is not evidence of an absent Usadha manuscript: its selected front side contains symptoms, not a plant preparation. Both old selectable crops were resized to 2400 px wide using Lanczos, JPEG quality 95, no chroma subsampling:

| File | Original dimensions | Crop | Output | SHA-256 |
| --- | --- | --- | --- | --- |
| `lontar-usada-lare-depan.jpg` | 4730 × 3153 | `(120, 1110, 4680, 1600)` | 2400 × 258 | `ed50de58933f3f837b1850707e3cdca9e1f83b3355372fe3f01938aa349e2d1f` |
| `lontar-usada-rare-154-depan.jpg` | 5184 × 3456 | `(900, 1140, 4550, 1640)` | 2400 × 329 | `3912ae6f95caa3953f7db32e6556b1156635bc2f0e5b1f5132116983131b0cf5` |

## Sharing image

`../og-bukalontar-v1.jpg` incorporates a cropped and resized version of `lontar-sukawati-usada-07.jpg` by Eka343. The composition adds the BukaLontar logo, typography, and parchment background. This adapted sharing image is distributed under [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/). Attribution and the crop notice also appear on the image.

Browser enhancement affects the preview only; the API receives original selected bytes. Neither source attribution nor matching a transcript certifies a generated translation, medical safety, or model accuracy. The typeset Kawi comparison is evaluation-only, not a selectable example.
