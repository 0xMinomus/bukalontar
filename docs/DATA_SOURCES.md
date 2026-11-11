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

Contoh yang dapat dipilih adalah sisi depan Usada Lare dan Usada Rare 154 dari Balai Bahasa Bali. Keduanya memiliki transkripsi aksara Bali tervalidasi di Wikisource Bali. Pemotongan membuang latar meja dan sisi belakang, bukan menggambar ulang aksara. Penolakan model tidak menentukan identitas aksara.

| File | Kreator | Lisensi | Sumber |
| --- | --- | --- | --- |
| `lontar-usada-lare-depan.jpg` | Eka343 | CC BY-SA 4.0 | [Commons](https://commons.wikimedia.org/wiki/File:Bali-lontar-Balai_Bahasa_Bali-Usada_Lare.jpg), [transkripsi tervalidasi](https://ban.wikisource.org/w/index.php?oldid=105197) |
| `lontar-usada-rare-154-depan.jpg` | Eka343 | CC BY-SA 4.0 | [Commons](https://commons.wikimedia.org/wiki/File:Bali-lontar-Balai_Bahasa_Bali-Usada_Rare_154.jpg), [transkripsi tervalidasi](https://ban.wikisource.org/w/index.php?oldid=105194) |
| `lontar-wanagiri-usada-3.jpg` | Eka343 | CC BY-SA 4.0 | [Wikimedia Commons](https://commons.wikimedia.org/wiki/File:Bali-lontar-Wanagiri-Usada-3.jpg) |
| `lontar-sukawati-usada-07.jpg` | Eka343 | CC BY-SA 4.0 | [Wikimedia Commons](https://commons.wikimedia.org/wiki/File:Bali-lontar-Sukawati-Usada-07.jpg) |

Wanagiri/Sukawati dipertahankan untuk aset dekoratif/kartu berbagi dan evaluasi historis, bukan pilihan contoh. Detail lisensi, koordinat crop, dan SHA-256 tersedia di [`apps/web/public/samples/CREDITS.md`](../apps/web/public/samples/CREDITS.md).

## Korpus DeepLontar

Korpus lokal menggunakan [DeepLontar v2](https://doi.org/10.6084/m9.figshare.20103803.v2) oleh Ni Putu Sutramiani, Nanik Suciati, Daniel Oranova Siahaan, I Nengah Duija, dan I Wayan Agus Surya Darma, lisensi [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/). [Metadata Figshare](https://api.figshare.com/v2/articles/20103803/versions/2) mencantumkan lisensi dan berkas sumber. Sertakan atribusi, tautan lisensi, serta keterangan perubahan saat menggunakan data.

- `data/deeplontar-manifest.json` mencatat nama, ukuran, dimensi, dan SHA-256 untuk 200 gambar asli (`1a`–`100b`) serta anotasi yang dipasangkan berdasarkan nama. Gambar berjumlah 81.043.523 byte, anotasi 1.750.733 byte. Foto dan label tersimpan lokal di `data/corpus/deeplontar/`, yang tidak masuk Git.
- Sumber berisi 200 gambar asli dan 400 hasil augmentasi. Manifest hanya memilih gambar asli; kedua arsip asli tetap disimpan untuk reproduksi. Dua sisi daun tidak dihitung sebagai dua naskah yang berbeda. Identitas naskah, asal masing-masing daun, dan keragaman koleksi tidak dijelaskan per gambar dalam arsip.
- [Artikel dataset](https://www.nature.com/articles/s41597-022-01867-5) menyatakan label karakter telah divalidasi ahli: bounding box YOLO untuk 55 kelas, bukan transkripsi halaman, transliterasi, terjemahan, atau label resep. Sebanyak 200 pasang lokal berisi 45.261 anotasi karakter. Belum ada ground truth untuk mengukur akurasi terjemahan atau laporan Usadha.
- Validator menemukan 76 anotasi dengan ID kelas di luar rentang 0–54 pada `76a.txt`. Berkas asli dipertahankan sesuai checksum, tetapi seluruh 198 kotaknya dikarantina. Tersisa 199 berkas label dengan 45.063 kotak yang memenuhi pemeriksaan format; ini bukan validasi ulang filologis.
- `1a.jpg`, `50b.jpg`, dan `100a.jpg` diperiksa secara visual dan memperlihatkan aksara pada daun lontar; tidak ada pemeriksaan filologis untuk setiap gambar. Maklumat tentang aksara dan label bergantung pada sumber, bukan transkripsi ahli oleh proyek ini.
- Instal Python 3 dan Pillow (`python3 -m pip install Pillow`), lalu jalankan `python3 scripts/deeplontar.py acquire` untuk mengunduh arsip tersemat checksum, mengekstrak 200 pasang, dan memvalidasinya; `python3 scripts/deeplontar.py validate` memeriksa berkas lokal, checksum, dekode JPEG, dan format label. Memerlukan sekitar 300 MB ruang disk.
- Label karakter ini dipakai untuk melatih detektor teks dan pengenal glif lokal; lihat [OCR_TRAINING.md](OCR_TRAINING.md). Pembagian train/dev/test menjaga kedua sisi lembar tetap satu kelompok. Gambar saja bukan transkripsi; evaluasi halaman memakai lima halaman Wikisource berlabel. Tidak ada fine-tuning model bahasa melalui provider.

## Contoh Glif dan Halaman Berlabel

- 55 contoh glif aksara Bali dari Tabel 1 [Sutramiani dkk., Scientific Data 2022](https://doi.org/10.1038/s41597-022-01867-5) (artikel CC BY 4.0). Gambar glif di-embed di `apps/api/internal/readingexamples/assets/`.
- Lima halaman pindaian [Aji Kreket](https://commons.wikimedia.org/wiki/File:Bali-lontar-aji-kreket-550ppi.pdf) dari Wikisource Bali dengan transkripsi yang divalidasi (revisi 225–254, CC BY-SA 4.0; Commons menyatakan pindaian PD-ID-exempt). Semua dari satu naskah, disimpan di `data/reading-examples/holdout/` sebagai data uji dan tidak dikirim sebagai contoh konteks.
- Validasi Wikisource bukan pemeriksaan filologis independen oleh tim. Keselarasan citra–teks untuk derajat 1280 px tidak diverifikasi ulang.

## Knowledge Base

- Lokasi: `apps/api/internal/kb/data/taru-pramana.json`.
- Fungsi: pencocokan nama tanaman, famili, senyawa, rujukan, dan energy character.
- Knowledge base dikurasi untuk MVP. Entri di luar katalog tidak diberi taksonomi atau sitasi.
- Kutipan dan istilah dari lontar dipertahankan sebagai bagian dari konteks.
- `reference_links` dikurasi di server: pencarian PubMed per nama ilmiah dan satu tinjauan Psidium guajava ([PMID 18353572](https://pubmed.ncbi.nlm.nih.gov/18353572/)). Pencarian adalah pintu masuk literatur, bukan sitasi studi tertentu.
- Tautan tersebut bukan sumber naskah dan tidak membuktikan keamanan atau khasiat resep. Deskripsi `reference_studies` lama bukan kutipan bibliografis yang telah diverifikasi; pembaca perlu memeriksa publikasi aslinya.

## Model dan Gateway

- Gateway: 9router dengan endpoint OpenAI-compatible.
- Model produk: hanya `ag/gemini-3.8-flash-low`.
- Tidak ada fallback model.
- Konfigurasi dan API key tidak disimpan di repository.

## Library dan Aset

- React, Vite, TypeScript, Tailwind CSS, Radix UI, Lucide.
- class-variance-authority, clsx, tailwind-merge.
- Font Cormorant Garamond dan Lora melalui Google Fonts.
- Go standard library.

## Data Pribadi

Repository tidak memuat data pribadi, data pasien, credential demo, atau API key. Jangan mengunggah foto lontar yang memuat informasi pribadi tanpa izin.
