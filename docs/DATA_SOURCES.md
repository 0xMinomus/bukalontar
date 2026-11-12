# Sumber Data dan Kredit

Dokumen ini mencatat sumber naskah, foto, knowledge base, dan aset pihak ketiga yang dipakai BukaLontar.

Snapshot rujukan mengikuti tag [`v1.1-final-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.1-final-submission). Kredit aset tidak berarti pengarang, lembaga sumber, atau peninjau Wikisource menyetujui hasil AI BukaLontar.

## Naskah dan Konteks

| Sumber | Keterangan |
| --- | --- |
| Dinas Kebudayaan Provinsi Bali | Data sensus dan konservasi lontar masyarakat Bali |
| UPTD Gedong Kirtya, Singaraja | Katalog dan koleksi naskah lontar, termasuk kelompok Usadha |
| Lontar Taru Pramana, Va. 1384/10 | Rujukan etnobotani dan resep tanaman obat |
| Lontar Usadha Rare, Va. 1391/11 | Rujukan Usadha untuk kesehatan ibu dan anak |

Naskah asli tetap menjadi sumber utama. Hasil AI tidak menggantikan pemeriksaan filologi.

## Foto Sampel

Contoh yang dapat dipilih adalah sisi depan Usada Lare dan Usada Rare 154 dari Balai Bahasa Bali. Keduanya memiliki transkripsi aksara Bali yang telah dikoreksi di Wikisource Bali (pagequality 3), belum diverifikasi ahli independen. Pemotongan membuang latar meja dan sisi belakang, bukan menggambar ulang aksara. Penolakan model tidak menentukan identitas aksara.

| File | Kreator | Lisensi | Sumber |
| --- | --- | --- | --- |
| `lontar-usada-lare-depan.jpg` | Eka343 | CC BY-SA 4.0 | [Commons](https://commons.wikimedia.org/wiki/File:Bali-lontar-Balai_Bahasa_Bali-Usada_Lare.jpg), [transkripsi dikoreksi](https://ban.wikisource.org/w/index.php?oldid=105197) |
| `lontar-usada-rare-154-depan.jpg` | Eka343 | CC BY-SA 4.0 | [Commons](https://commons.wikimedia.org/wiki/File:Bali-lontar-Balai_Bahasa_Bali-Usada_Rare_154.jpg), [transkripsi dikoreksi](https://ban.wikisource.org/w/index.php?oldid=105194) |
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
- Label karakter dipakai untuk melatih detektor teks dan pengenal glif lokal; lihat [OCR_TRAINING.md](OCR_TRAINING.md). Train memuat 159 citra dari 80 kelompok daun; dev dan test masing-masing 20 citra dari 10 kelompok. Kedua sisi daun tetap pada kelompok yang sama. Gambar dan kotak karakter bukan transkripsi; evaluasi halaman memakai lima halaman Wikisource. Tidak ada fine-tuning model bahasa melalui provider.

## Contoh Glif dan Halaman Berlabel

- 55 contoh glif aksara Bali dari Tabel 1 [Sutramiani dkk., Scientific Data 2022](https://doi.org/10.1038/s41597-022-01867-5) (artikel CC BY 4.0). Gambar glif di-embed di `apps/api/internal/readingexamples/assets/`.
- Lima halaman pindaian [Aji Kreket](https://commons.wikimedia.org/wiki/File:Bali-lontar-aji-kreket-550ppi.pdf) dengan transkripsi komunitas Wikisource (revisi 225–254, CC BY-SA 4.0; Commons menyatakan pindaian PD-ID-exempt). Semua dari satu naskah, disimpan di `data/reading-examples/holdout/` sebagai data uji dan tidak dikirim sebagai contoh konteks.
- Status peninjauan Wikisource bukan pemeriksaan filologis independen oleh tim. Keselarasan citra–teks pada gambar lebar 1280 px tidak diverifikasi ulang.
- [NusaAksara](https://huggingface.co/datasets/NusaAksara/NusaAksara) menyediakan anotasi manusia untuk segmen buku, termasuk Bali dan Jawa. Penulis menyatakan penggunaan nonkomersial; datanya tidak dimasukkan ke produk atau Git dan tidak dipakai untuk mengklaim akurasi foto lontar.
- Foto negatif Panji, Babad Surapati, dan contoh Kawi cetak memiliki kredit serta batas label terpisah di [kredit evaluasi](evaluation/script-inputs/CREDITS.md). Kawi cetak bukan foto lontar; kategori bahasa atau lokasi bukan sertifikasi aksara.

## Knowledge Base

- Lokasi: `apps/api/internal/kb/data/taru-pramana.json`.
- Fungsi: pencocokan nama tanaman, famili, senyawa, dan rujukan. Sifat tradisional pada hasil mengikuti bacaan AI, bukan nilai bawaan KB.
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
- Font Cormorant Garamond dan Lora dibundel lokal melalui Fontsource/aset WOFF2; lisensi OFL menyertai aset di `apps/web/src/assets/fonts/`.
- Go standard library.

## Data Pribadi

Jangan mengunggah foto lontar yang memuat informasi pribadi atau materi terbatas tanpa izin. Backend meneruskan foto ke gateway dan penyedia AI; kebijakan retensi pihak ketiga belum diverifikasi. Browser menyimpan foto dan hasil sesi aktif di IndexedDB. Kebijakan keamanan ada di [SECURITY.md](../SECURITY.md).
