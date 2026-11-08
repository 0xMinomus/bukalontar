# Sumber Data dan Kredit

Dokumen ini mencatat sumber naskah, foto, knowledge base, dan aset pihak ketiga yang dipakai BukaLontar.

## Naskah dan Konteks

| Sumber | Keterangan |
| --- | --- |
| Dinas Kebudayaan Provinsi Bali | Data sensus dan konservasi lontar masyarakat Bali |
| UPTD Gedong Kirtya, Singaraja | Katalog dan koleksi naskah lontar, termasuk kelompok Usadha |
| Lontar Taru Pramana, Va. 1384/10 | Rujukan etnobotani dan resep tanaman obat |
| Lontar Usadha Rare, Va. 1391/11 | Rujukan Usadha untuk kesehatan ibu dan anak |

Naskah asli tetap menjadi sumber utama. Hasil AI tidak menggantikan pemeriksaan filologi.

## Foto Sampel

Foto sampel dipakai untuk menguji antarmuka dan pipeline.

| File | Kreator | Lisensi | Sumber |
| --- | --- | --- | --- |
| `lontar-wanagiri-usada-3.jpg` | Eka343 | CC BY-SA 4.0 | [Wikimedia Commons](https://commons.wikimedia.org/wiki/File:Bali-lontar-Wanagiri-Usada-3.jpg) |
| `lontar-sukawati-usada-07.jpg` | Eka343 | CC BY-SA 4.0 | [Wikimedia Commons](https://commons.wikimedia.org/wiki/File:Bali-lontar-Sukawati-Usada-07.jpg) |

Detail lisensi juga tersedia di [`apps/web/public/samples/CREDITS.md`](../apps/web/public/samples/CREDITS.md).

## Knowledge Base

- Lokasi: `apps/api/internal/kb/data/taru-pramana.json`.
- Fungsi: pencocokan nama tanaman, famili, senyawa, rujukan, dan energy character.
- Knowledge base dikurasi untuk MVP. Entri di luar katalog tidak diberi taksonomi atau sitasi.
- Kutipan dan istilah dari lontar dipertahankan sebagai bagian dari konteks.

## Model dan Gateway

- Gateway: 9router dengan endpoint OpenAI-compatible.
- Model default: `ag/gemini-3.8-flash-low`.
- Fallback: opsional melalui Command Code Provider API.
- Konfigurasi dan API key tidak disimpan di repository.

## Library dan Aset

- React, Vite, TypeScript, Tailwind CSS, Radix UI, Lucide.
- class-variance-authority, clsx, tailwind-merge.
- Font Cormorant Garamond dan Lora melalui Google Fonts.
- Go standard library.

## Data Pribadi

Repository tidak memuat data pribadi, data pasien, credential demo, atau API key. Jangan mengunggah foto lontar yang memuat informasi pribadi tanpa izin.
