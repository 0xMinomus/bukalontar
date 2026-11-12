# Sumber Data dan Kredit

Dokumen ini mencatat sumber naskah, foto, knowledge base, dan aset pihak ketiga yang dipakai BukaLontar.

Tag [`v1.1-final-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.1-final-submission) tetap menjadi snapshot lama. Dokumen pada `main` juga mencatat empat contoh arsip yang ditambahkan setelah tag. Kredit aset tidak berarti pengarang, lembaga sumber, atau peninjau Wikisource menyetujui hasil AI BukaLontar.

## Naskah dan Konteks

| Sumber | Keterangan |
| --- | --- |
| Dinas Kebudayaan Provinsi Bali | Data sensus dan konservasi lontar masyarakat Bali |
| UPTD Gedong Kirtya, Singaraja | Katalog dan koleksi naskah lontar, termasuk kelompok Usadha |
| Lontar Taru Pramana, Va. 1384/10 | Rujukan etnobotani dan resep tanaman obat |
| Lontar Usadha Rare, Va. 1391/11 | Rujukan Usadha untuk kesehatan ibu dan anak |

Naskah asli tetap menjadi sumber utama. Hasil AI tidak menggantikan pemeriksaan filologi.

## Foto Sampel

Contoh pilihan pada `main` mencakup dua Usadha dengan racikan pada sisi depan dan dua genre non-Usadha. Semuanya foto lontar beraksara Bali. Bahasa Kawi pada katalog tidak mengubah identitas aksaranya. Pemotongan mempertahankan satu daun utuh tanpa menggambar ulang glif.

| Contoh | Genre | Foto / lisensi | Transkripsi sumber |
| --- | --- | --- | --- |
| Usadha Sukawati 07 | Usadha | [Eka343, CC BY-SA 4.0](https://commons.wikimedia.org/wiki/File:Bali-lontar-Sukawati-Usada-07.jpg) | [Revisi 98029, sisi depan](https://ban.wikisource.org/w/index.php?title=Kaca:Bali-lontar-Sukawati-Usada-07.jpg&oldid=98029) |
| Usadha Rare Kemenuh 56 | Usadha | [Eka343, CC BY-SA 4.0](https://commons.wikimedia.org/wiki/File:Bali-lontar-Kemenuh-Usadha_Rare-56.jpg) | [Revisi 107871, sisi depan](https://ban.wikisource.org/w/index.php?title=Kaca:Bali-lontar-Kemenuh-Usadha_Rare-56.jpg&oldid=107871) |
| Carcan Kucing, 1B | Carcan | [Internet Archive, PD-ID-exempt](https://commons.wikimedia.org/wiki/File:Bali-lontar-carcan-kucing-250ppi.pdf) | [Render Latin revisi 107176, hanya 1B](https://ban.wikisource.org/w/index.php?title=Carcan_Kucing&oldid=107176) |
| Kakawin Bharatayuda Maarti 01 | Kakawin / susastra | [Eka343, CC BY-SA 4.0](https://commons.wikimedia.org/wiki/File:Bali-lontar-Balai_Bahasa_Bali-Kekawin_Bharatayuda_Maarti_01.jpg) | [Render Latin revisi 134900, sisi depan](https://ban.wikisource.org/w/index.php?title=Kekawin_Bharatayuda_Maarti_(Balai_Bahasa_Prov._Bali)&oldid=134900) |

Rujukan aksara primer untuk dua render halaman utama: [Carcan 1B, revisi 102232](https://ban.wikisource.org/w/index.php?title=Kaca:Bali-lontar-carcan-kucing-250ppi.pdf/2&oldid=102232) dan [Kakawin 01, revisi 108451](https://ban.wikisource.org/w/index.php?title=Kaca:Bali-lontar-Balai_Bahasa_Bali-Kekawin_Bharatayuda_Maarti_01.jpg&oldid=108451). Revisi halaman utama tidak membekukan transklusi atau mesin render Latin; salinan teks yang digunakan aplikasi disimpan dalam JSON tersemat.

Lare dan Rare 154 lama dipertahankan untuk evaluasi historis, bukan tombol contoh. Sisi depan Lare memuat tanda gejala, bukan racikan tanaman; kartu resep kosong tidak membuktikan seluruh naskah bukan Usadha. Wanagiri/Sukawati versi lama tetap menjadi aset dekoratif/berbagi. Lisensi, koordinat crop, resolusi, dan SHA-256 seluruh aset ada di [CREDITS.md](../apps/web/public/samples/CREDITS.md).

### Analisis berbantuan transkripsi rujukan

`apps/api/internal/ai/manuscript_references.json` menyimpan hash foto, judul, genre, URL revisi, render Latin sumber, dan versi dengan tambahan pemisahan kata. Saat startup server memeriksa bahwa perubahan hanya whitespace, bukan huruf atau tanda baca. Atribusi transkripsi adalah kontributor Wikisource Bali, CC BY-SA 4.0; riwayat revisi mencatat kontributornya. Render Latin dan pemisahan kata proyek bukan transliterasi ahli yang diverifikasi terpisah.

Foto cocok memakai teks rujukan sebagai input Gemini, bukan OCR foto atau jawaban analisis tersimpan. Server mengembalikan transliterasi sumber persis, sedangkan terjemahan dan ekstraksi tetap dihasilkan saat permintaan. Genre Carcan/Kakawin mencegah kartu obat pada dua arsip nonmedis tersebut. `reading_reference` mencatat metode, sumber, lisensi, dan hash; UI/PDF menampilkannya. Ini bantuan arsip yang disetujui pengguna, bukan hasil holdout independen atau peningkatan akurasi OCR.

Nama file dan `source_id` tidak memilih rujukan. Foto yang berubah byte, termasuk hasil kompresi ulang atau crop lain, tetap melalui pemeriksaan aksara dan pembacaan foto eksperimental. Penolakan dan halusinasi pada jalur itu belum terselesaikan.

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

- Lokasi: `apps/api/internal/kb/data/taru-pramana.json` untuk entri kurasi awal dan `apps/api/internal/kb/data/wikidata-plants.json` untuk tambahan botani terstruktur.
- Fungsi: pencocokan nama tanaman, famili, senyawa, dan rujukan. Sifat tradisional pada hasil mengikuti bacaan, bukan nilai bawaan KB. Entri impor tidak menambahkan klaim senyawa atau khasiat.
- Katalog memadukan 100 entri kurasi awal dengan takson tambahan dari Wikidata (CC0) yang nama ilmiahnya cocok tepat dengan GBIF Backbone Taxonomy (CC BY 4.0). Sinonim digabung menurut accepted usage key dan alias ambigu tidak dipaksa cocok. Dapdap tetap **Erythrina spp.**, spesies belum dipastikan. Ini bukan daftar tanaman yang dipastikan tercantum dalam Taru Pramana.
- Kutipan dan istilah dari lontar dipertahankan sebagai bagian dari konteks.
- [Sumber, lisensi, dan reproduksi](KNOWLEDGE_BASE_SOURCES.md) mencatat kurasi awal, snapshot impor, serta batas identifikasi istilah naskah. Sembung, Kacemcem, dan Poh Gedang mencakup ejaan transkripsi contoh arsip.
- `reference_links` ditentukan server: pencarian PubMed per takson dan satu tinjauan Psidium guajava ([PMID 18353572](https://pubmed.ncbi.nlm.nih.gov/18353572/)). Pencarian adalah pintu masuk literatur, bukan sitasi studi tertentu. Entri baru tidak menambahkan klaim senyawa atau khasiat; deskripsi rujukannya berisi sumber botani.
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

Jangan mengunggah foto lontar yang memuat informasi pribadi atau materi terbatas tanpa izin. Untuk foto baru backend meneruskan gambar ke penyedia AI; untuk arsip cocok hanya teks rujukan dikirim. Retensi pihak ketiga belum diverifikasi. Browser menyimpan foto dan hasil sesi aktif di IndexedDB. Lihat [SECURITY.md](../SECURITY.md).
