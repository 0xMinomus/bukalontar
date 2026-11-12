<div align="center">

<img src="logos/bukalontar-icon.png" alt="Logo BukaLontar" width="112" />

# BukaLontar

**Baca lontar Bali, dari aksara sampai konteks.**

Aplikasi web eksperimental untuk membantu membaca foto lontar beraksara Bali, menghasilkan usulan transliterasi dan terjemahan, serta menelusuri konteks naskah. Untuk Usadha, aplikasi dapat menampilkan resep yang lolos pemeriksaan kutipan dan pencocokan tanaman dengan basis pengetahuan Taru Pramana. Hasil tetap perlu diperiksa terhadap foto dan oleh ahli.

![Status](https://img.shields.io/badge/status-MVP-b08d57)
![AI](https://img.shields.io/badge/AI-multimodal-4f6f52)
![Stack](https://img.shields.io/badge/stack-React%20%2B%20Go-2f3e46)
![Competition](https://img.shields.io/badge/competition-Bali%20DigiFest%202026-9c6644)

[Repositori](https://github.com/0xMinomus/bukalontar) | [Dokumentasi](docs/README.md) | [AI Disclosure](docs/AI_DISCLOSURE.md)

Versi final terkini: [`main`](https://github.com/0xMinomus/bukalontar/tree/main). Tag [`v1.1-final-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.1-final-submission) adalah snapshot historis, bukan versi produksi terbaru. [Website](https://www.bukalontar.my.id) | [Status akhir dan pemeriksaan deployment](docs/STATUS.md) | [Perbandingan baseline](docs/AI_EVALUATION.md#upgrade-pengaman-kutipan-dan-pembanding-baseline)

Perbaikan pada `main` setelah tag: empat contoh arsip memakai transkripsi rujukan dengan label terbuka. Tag submission tidak dipindah. [Metode contoh arsip](docs/DATA_SOURCES.md#analisis-berbantuan-transkripsi-rujukan) | [Hasil uji contoh](docs/AI_EVALUATION.md#empat-contoh-arsip-berbantuan-transkripsi).

</div>

---

## What's New

Perubahan sejak [`v1.0-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.0-submission) hingga versi final terkini di `main`. Daftar berikut memisahkan perubahan yang sudah masuk tag [`v1.1-final-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.1-final-submission) dari pembaruan setelah tag tersebut. [Lihat seluruh perubahan kode](https://github.com/0xMinomus/bukalontar/compare/v1.0-submission...main).

### Pipeline AI dan pengaman resep

Perubahan berikut sudah termasuk dalam tag `v1.1-final-submission`:

- Mengganti dua pembacaan paralel, pemilihan hasil terpanjang/terbanyak, dan penggabungan resep dengan satu pembacaan setelah pemeriksaan aksara. Hasil dari beberapa jawaban model tidak lagi digabung menjadi daftar tanaman tambahan.
- Menetapkan Gemini 3.8 Flash Low sebagai satu-satunya model produk melalui 9router. Fallback model dan variabel konfigurasi fallback dihapus; server menolak konfigurasi model produk lain.
- Menambahkan pemeriksaan aksara terpisah sebelum membaca foto baru. Hanya keputusan `balinese/high` yang mencapai pembaca; contoh glif dan katalog tanaman tidak disisipkan ke pemeriksaan ini.
- Mengirim foto asli ke API. Downscale, CLAHE, unsharp mask, dan mode diperjelas tetap tersedia untuk pratinjau, tetapi hasil olahan browser tidak lagi menggantikan input AI.
- Mengubah prompt agar membaca foto target tanpa menebak judul, katalog tanaman, atau genre Usadha dari `source_id`. Bahasa, grup naskah, dan bagian yang tidak diketahui dapat tetap dinyatakan tidak diketahui.
- Mewajibkan `source_quote` untuk setiap resep: kutipan utuh dari lembar pertama yang memuat nama tanaman dan mendukung racikan atau indikasi. Bahan pelengkap yang diklaim harus terdapat dalam kutipan yang sama.
- Menolak kutipan yang menjembatani bagian tidak terbaca, anotasi, elipsis, pengantar sebelum lembar pertama, atau lembar lain. Label `Lembar` dan `Leaf` diperlakukan setara; resep yang ditolak tidak boleh muncul kembali melalui normalisasi.
- Mempertahankan hasil tanpa resep ketika tidak ada bukti yang cukup. Penanda kosong dan teks tidak terbaca tidak lagi berubah menjadi kartu tanaman, dan bentuk obat tidak otomatis dianggap loloh.
- Menjaga sifat tradisional sesuai bacaan, bukan menimpanya dengan nilai bawaan KB. Taksonomi, senyawa, klaim farmakologi, dan sitasi keluaran model diganti dengan data server jika cocok; data ilmiah tanaman yang tidak cocok dikosongkan.
- Memisahkan aksara terdeteksi dari teks terbaca. Penolakan pembaca atau keluaran yang mengandung konflik aksara tidak dipertahankan sebagai bacaan/resep; transkripsi Unicode eksperimen tidak ditampilkan sebagai hasil terverifikasi.
- Memperluas kontrak JSON dengan kutipan resep, tautan rujukan, metadata bantuan pembacaan, nilai tidak diketahui, dan transkripsi Unicode opsional untuk evaluasi. Parser keputusan aksara menerima JSON berbungkus code fence; batas keluaran pembaca dinaikkan menjadi 8.000 token.
- Menambahkan tautan PubMed pada kartu, tab Rujukan, dan laporan PDF. Pencarian literatur dibedakan dari artikel tertentu; adanya tautan tidak dinyatakan sebagai bukti khasiat racikan.

### Bantuan pembacaan, data, dan evaluasi

Penambahan berikut juga sudah termasuk dalam tag `v1.1-final-submission`:

- Menambahkan akuisisi dan validasi korpus DeepLontar: manifest sumber, checksum, pemeriksaan citra/anotasi, serta karantina anotasi dengan ID kelas di luar rentang yang dipublikasikan.
- Melatih dan menyematkan detektor area teks serta pengenal glif lokal, dengan inferensi Go tanpa layanan tambahan. Pembagian train/dev/test mengelompokkan dua sisi daun fisik yang sama agar tidak tersebar antarbagian.
- Menambahkan 55 contoh glif bersumber, chart referensi, dan mode bantuan `none`, `examples`, `chart`, `ocr`, `leaves`, `both`, serta `all`. Default `examples` memakai tiga contoh glif; detektor dan crop daun tetap eksperimental, bukan OCR halaman penuh.
- Menambahkan lima halaman Aji Kreket dengan transkripsi komunitas sebagai holdout, beserta provenance dan batas independensi sumber. Halaman holdout tidak dipakai sebagai contoh jawaban dalam prompt produk.
- Menambahkan runner evaluasi, keluaran JSONL, penghitungan CER, abstensi, konsistensi bacaan/tanaman, latensi, dan pemeriksaan parity model Go/Python. Pengujian regresi mencakup pengaman aksara, kutipan resep, data kosong, dan grounding rujukan.
- Mencatat pembanding baseline, kandidat model, ablation bantuan, uji foto Bali/Jawa/Kawi, dan eksperimen yang tidak diintegrasikan. OCR literal terpisah, pemeriksa visual kedua, dan contoh foto–jawaban berpasangan tidak dipromosikan sebagai peningkatan akurasi produk.
- Mengganti dua contoh pilihan awal dengan sisi depan Usada Lare dan Usada Rare 154 pada tahap ini, disertai kredit foto, sumber transkripsi, lisensi, koordinat crop, dan hash. Contoh pilihan tersebut kemudian diganti lagi pada pembaruan setelah tag.
- Memperbarui dokumentasi arsitektur, sumber data, pelatihan OCR, disclosure AI, keamanan, konfigurasi, dan panduan kontribusi untuk mencerminkan perilaku eksperimental serta batas pengujian.

Rincian: [evaluasi AI dan pembanding baseline](docs/AI_EVALUATION.md) · [model lokal dan batas OCR](docs/OCR_TRAINING.md) · [disclosure AI](docs/AI_DISCLOSURE.md).

### Pembaruan final di `main` setelah tag

- Mengganti pilihan landing menjadi empat arsip: Usadha Sukawati, Usadha Rare Kemenuh, Carcan Kucing, dan Kakawin Bharatayuda. Setiap pilihan memiliki tautan foto, transkripsi rujukan, dan lisensi yang terpisah.
- Menambahkan jalur analisis berbantuan transkripsi untuk foto yang SHA-256 byte-nya cocok persis dengan arsip. Server mempertahankan bacaan Latin sumber, sedangkan Gemini tetap menerjemahkan dan mengekstrak isi pada setiap permintaan; bukan cache jawaban analisis.
- Mengikat identitas arsip ke isi foto, bukan nama file atau `source_id`. Server menetapkan provenance, genre, URL sumber, lisensi, catatan editorial, dan hash; metadata buatan model tidak dapat memilih arsip.
- Membatasi arsip Carcan dan Kakawin agar tidak menghasilkan kartu obat. Foto lain, termasuk arsip yang diubah byte-nya, tetap melalui pemeriksaan aksara dan pembacaan foto eksperimental.
- Menambahkan label **Berbantuan transkripsi rujukan** pada hasil dan PDF, serta membedakan kutipan transkripsi komunitas dari transliterasi AI. Pemisahan kata pada bacaan arsip diperiksa agar tidak mengubah huruf atau tanda baca sumber.
- Memperluas KB menjadi **1.316 entri lookup botani**, dengan sumber PROSEA, data terstruktur Wikidata, dan pencocokan nama ilmiah ke takson diterima GBIF. Ini katalog botani umum, bukan daftar spesies yang dipastikan tercantum dalam Usadha.
- Menambahkan snapshot sumber berlisensi, ID/revisi sumber, laporan impor, dan generator offline. Impor menolak pencocokan fuzzy/genus/bukan tumbuhan, menggabungkan sinonim menurut accepted usage key, dan menyaring alias ambigu agar tidak menambah kecocokan palsu.
- Mengganti pencocokan substring dengan indeks alias tepat, mendukung nama ilmiah, ejaan transliterasi, dan awalan bagian tanaman. Nama campuran tidak otomatis cocok hanya karena mengandung nama tanaman yang lebih pendek; alias yang bertabrakan ditolak saat memuat katalog.
- Melengkapi alias contoh seperti Sembung, Kacemcem, dan Poh Gedang dengan sumber. Dapdap tetap dicatat pada tingkat genus ketika spesies belum dipastikan; istilah naskah yang ambigu tidak diberi identifikasi ilmiah buatan.
- Menambahkan tautan sumber kutipan pada setiap kartu arsip dan laporan cetak, termasuk kartu yang belum cocok KB. Indikator **Identitas botani di KB** membedakan kecocokan botani dari tersedianya rujukan naskah.
- Memperbaiki pembungkusan URL rujukan panjang dan batas lebar panel hasil agar sumber tetap terbaca pada layar kecil dan laporan cetak.
- Menambahkan gerak gulir dan detail material bernuansa lontar: susunan daun saat unggahan bergulir, bingkai arsip, pengungkapan foto sumber, ornamen, serta respons hover. Gerak menghormati preferensi reduced motion dan memiliki fallback ketika fitur browser tidak tersedia.
- Memperbaiki rewrite SPA pada service web Vercel agar membuka atau me-refresh `/analisis` dan `/hasil` tidak berakhir dengan 404, sehingga aplikasi dapat memulihkan sesi lokal.
- Melengkapi catatan kesiapan submission, pemeriksaan deployment, kredit arsip, sumber KB, dan batas klaim. Perubahan lanjutan tidak memindahkan tag submission yang sudah diterbitkan.

Rincian: [empat contoh arsip](docs/AI_EVALUATION.md#empat-contoh-arsip-berbantuan-transkripsi) · [sumber dan lisensi KB](docs/KNOWLEDGE_BASE_SOURCES.md) · [status pemeriksaan deployment](docs/STATUS.md).

**Batas klaim:** perubahan ini memperketat pengaman dan menambah penelusuran sumber; belum membuktikan OCR, terjemahan, identifikasi tanaman, atau racikan obat selalu benar. Rujukan naskah bukan verifikasi botani, dan pencarian literatur bukan bukti klinis. Kemampuan unggah/tempel foto, lensa zoom, penyalinan, pemulihan sesi, dan ekspor PDF sudah ada pada baseline; yang baru adalah perbaikan alur, isi, dan penelusuran sumber di atas.

## Daftar Isi

- [What's New](#whats-new)
- [Tentang](#tentang)
- [Screenshot](#screenshot)
- [Masalah dan Target Pengguna](#masalah-dan-target-pengguna)
- [Cakupan Produk](#cakupan-produk)
- [Fitur](#fitur)
- [Cara Kerja AI](#cara-kerja-ai)
- [Arsitektur](#arsitektur)
- [Teknologi](#teknologi)
- [Menjalankan Proyek](#menjalankan-proyek)
- [Struktur Repositori](#struktur-repositori)
- [Data, Sumber, dan Kredit](#data-sumber-dan-kredit)
- [Keamanan dan Privasi](#keamanan-dan-privasi)
- [AI dan Technology Disclosure](#ai-dan-technology-disclosure)
- [Batasan Saat Ini](#batasan-saat-ini)
- [Tim](#tim)
- [Dokumentasi Lanjutan](#dokumentasi-lanjutan)

## Tentang

Lontar adalah warisan tulis Bali yang mencakup banyak genre: Usadha, Kakawin, Babad, Geguritan, tutur, dan naskah lain. Banyak di antaranya belum terbaca oleh generasi muda, peneliti, atau pengunjung museum karena dua lapis hambatan: aksara Bali kuno dan bahasa Kawi.

BukaLontar adalah alat bantu pembacaan lontar berbasis AI, bukan pengganti filolog. Unggah foto lempir untuk memperoleh usulan transliterasi Latin, terjemahan Indonesia dan Inggris, serta informasi sumber jika dapat dikenali. Bahasa Kawi/Jawa Kuno dapat ditulis dengan aksara Bali; bahasa tersebut tidak sama dengan aksara Kawi.

Untuk lontar Usadha, hasilnya diperkaya dengan lapisan khusus: nama tanaman lokal, nama ilmiah, bagian tanaman, bentuk racikan, indikasi, senyawa fitokimia, dan rujukan dari Taru Pramana. Lapisan ini tidak dipaksakan pada naskah yang bukan Usadha.

## Screenshot

<p align="center">
  <img src="docs/assets/landing-page.png" alt="Landing page BukaLontar" width="900" />
</p>

Landing page BukaLontar: unggah satu lempir, mulai dari arsip contoh, lalu baca hasilnya.

## Masalah dan Target Pengguna

**Masalah utama**

- Aksara Bali pada lontar ditulis tanpa spasi dan memakai pasang aksara yang bertumpuk.
- Banyak naskah memakai bahasa Kawi klasik, bukan bahasa Bali sehari-hari.
- Foto lontar sering gelap, aus, atau miring, sehingga pembacaan manual memakan waktu.
- Pengetahuan tanaman obat di lontar Usadha tersebar dan sulit dihubungkan dengan taksonomi modern.

**Target pengguna**

- Pelajar, mahasiswa, dan masyarakat yang ingin membaca lontar Bali.
- Peneliti filologi, sejarah, dan etnobotani.
- Kurator museum, perpustakaan, dan arsip budaya.
- Praktisi herbal serta peneliti fitofarmaka yang bekerja dengan lontar Usadha.

## Cakupan Produk

| Alur | Naskah | Keluaran utama |
| --- | --- | --- |
| Pembacaan lontar umum | Foto lontar beraksara Bali | Usulan transliterasi, terjemahan Indonesia dan Inggris, serta ringkasan ketika teks dapat dibaca |
| Alur Usadha | Foto lontar dengan dugaan resep Usadha | Keluaran umum dan maksimal enam kartu resep yang lolos pemeriksaan kutipan |

Alur Usadha bersifat otomatis berdasarkan isi naskah. Tidak ada klaim bahwa setiap lontar adalah Usadha, dan naskah non-Usadha tetap dapat dibaca sebagai teks.

Contoh landing saat ini: **Usadha Sukawati**, **Usadha Rare Kemenuh**, **Carcan Kucing**, dan **Kakawin Bharatayuda**. Dua contoh pertama memuat racikan pada sisi yang dipilih; dua lainnya menunjukkan naskah non-Usadha. Untuk foto yang byte-nya cocok persis dengan arsip, bacaan Latin berasal dari transkripsi komunitas Wikisource. Gemini menerjemahkan dan mengekstrak isi saat permintaan berlangsung. UI dan PDF menandai **Berbantuan transkripsi rujukan**; ini bukan demonstrasi OCR murni atau bukti akurasi pada foto baru.

## Fitur

### Pembacaan lontar

- Downscale, CLAHE, dan unsharp mask untuk pratinjau browser; API menerima foto asli.
- Usulan transliterasi Latin, dengan penanda bagian tidak terbaca.
- Usulan terjemahan Indonesia dan Inggris.
- Bahasa teks, sumber, dan nomor lempir jika dikenali, serta confidence model yang bukan ukuran akurasi.
- Tampilan hasil yang bisa disalin dan diekspor menjadi laporan PDF.

### Alur Usadha

- Deteksi beberapa resep dalam satu lembar.
- Identifikasi nama tanaman lokal dan pencocokan dengan Taru Pramana.
- Bagian tanaman, sifat tradisional, bentuk racikan, dan indikasi.
- Senyawa fitokimia dan rujukan ilmiah yang tersedia di knowledge base.
- Label "belum tercatat" saat tanaman tidak ditemukan di knowledge base.
- Tautan PubMed dari knowledge base pada kartu dan laporan; pencarian literatur dibedakan dari publikasi tertentu.
- Disclaimer medis dari server, bukan dari model.

### Antarmuka

- Lensa lontar dengan zoom dan mode kontras.
- Pratinjau sebelum analisis.
- Mode mock untuk pengembangan UI tanpa API.
- Tata letak responsif dan laporan cetak.

## Cara Kerja AI

1. Browser menyiapkan pratinjau yang dapat diperjelas, lalu mengirim foto asli ke `POST /api/analyze-manuscript`.
2. Backend memvalidasi ukuran dan tipe file, lalu mencocokkan SHA-256 dengan empat foto arsip. Kecocokan memilih transkripsi rujukan dan genre sumber; bukan nama file atau `source_id`.
3. Foto yang tidak cocok tetap melalui pemeriksaan aksara tanpa katalog atau contoh. Hanya `balinese/high` mencapai pembaca; default `examples` menambahkan tiga contoh glif. Detektor lokal tetap eksperimental, bukan OCR halaman penuh.
4. Gemini 3.8 Flash Low membaca foto baru, atau menerjemahkan dan mengekstrak transkripsi yang dipasok untuk arsip cocok. Bacaan Latin arsip dipertahankan server dari sumber, bukan diganti bacaan model. Tidak ada fallback model atau jawaban hasil analisis yang di-cache.
5. JSON diparse dan divalidasi; satu percobaan perbaikan struktur diperbolehkan.
6. Resep memerlukan kutipan utuh dari lembar pertama yang menyebut tanaman dan bahan tambahan yang diklaim. Sumber kutipan adalah bacaan AI pada foto baru atau transkripsi rujukan pada arsip cocok. Bagian tidak terbaca, elipsis, pengantar, dan lembar lain tidak boleh menjembatani bukti. Arsip Carcan/Kakawin tidak menghasilkan kartu resep.
7. Nama ilmiah, famili, senyawa, dan rujukan berasal dari KB jika cocok. Sifat tradisional mengikuti bacaan, bukan nilai bawaan KB.
8. Frontend dan PDF membedakan analisis foto eksperimental dari analisis berbantuan transkripsi rujukan. Kutipan yang cocok tidak membuktikan terjemahan benar atau racikan aman.

Server menetapkan disclaimer medis dan rujukan KB. Model masih dapat mengarang bacaan beserta kutipan yang tampak konsisten; pengaman tidak menggantikan pemeriksaan ahli.

## Arsitektur

```mermaid
flowchart TD
    U["Browser: foto asli"] --> A["Go API: validasi upload"]
    A --> M{"SHA-256 cocok arsip?"}
    M -->|tidak| C["Gemini: klasifikasi aksara"]
    M -->|ya| T["Transkripsi sumber<br/>Gemini: terjemahan dan ekstraksi"]
    C --> D{"balinese/high?"}
    D -->|tidak| N["Respons penolakan tanpa bacaan/resep"]
    D -->|ya| R["Gemini: satu pembacaan dan terjemahan"]
    E["Tiga contoh glif default<br/>crop lokal eksperimental"] --> R
    R --> V["Validasi JSON dan batas kutipan"]
    T --> V
    V --> K["Grounding Taru Pramana"]
    K --> H["Hasil dengan peringatan dan laporan PDF"]
    N --> U
    H --> U
```

Detail teknis ada di [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md). Kontrak respons ada di [`packages/schema/analysis.json`](packages/schema/analysis.json).

## Teknologi

| Layer | Teknologi |
| --- | --- |
| Frontend | React 19, TypeScript, Vite, Tailwind CSS v4, Radix UI, Lucide |
| Pengolahan gambar | HTML5 Canvas API |
| Backend | Go standard library; binary mandiri untuk deployment non-Vercel |
| Gateway AI | 9router dengan endpoint OpenAI-compatible |
| Model | Hanya `ag/gemini-3.8-flash-low` melalui 9router, tanpa fallback model |
| Knowledge base | JSON Taru Pramana yang di-embed dengan `go:embed` |
| Bantuan pembacaan | Detektor teks ridge + LinearSVC HOG dilatih dari DeepLontar; 55 contoh glif ahli di-embed |
| Kontrak | JSON Schema di `packages/schema/analysis.json` |
| Deployment | Vercel Services: `web` (Vite) dan `api` (Go); binary mandiri opsional |

## Menjalankan Proyek

Prasyarat:

- Node.js 22.12 atau lebih baru (Vite juga mendukung Node.js 20.19 pada seri 20.x)
- npm
- Go 1.26 atau lebih baru

```bash
npm install
cp apps/api/.env.example apps/api/.env
# Isi AI_BASE_URL dan AI_API_KEY sesuai instance 9router; pertahankan AI_MODEL=ag/gemini-3.8-flash-low.
```

Terminal pertama:

```bash
make dev-api
```

Terminal kedua:

```bash
npm run dev
```

Buka `http://localhost:5173`.

Build produksi:

```bash
make build
./bin/bukalontar
```

Buka `http://localhost:8080`. Health check tersedia di `GET /health`.

## Deployment URL

- Live MVP: <https://www.bukalontar.my.id>
- Health check: <https://www.bukalontar.my.id/health>
- Hosting: Vercel Services, service `web` dan `api`
- AI gateway: 9router di `https://dashboard.minomus.my.id/v1`

Untuk deployment Vercel, `AI_BASE_URL` harus diisi endpoint publik 9router, bukan `127.0.0.1`. Health check akan mengembalikan `503` bila env tersebut belum masuk ke project.

Pada 12 November 2026, health produksi melaporkan `ok=true`, 10 entri KB, dan model `ag/gemini-3.8-flash-low`. Endpoint tidak melaporkan commit deployment; health yang baik tidak membuktikan versi produksi sama dengan tag submission.

Pengujian:

```bash
cd apps/api
go test ./...
go vet ./...

cd ../..
npm run lint
npm run build
```

## Struktur Repositori

```text
apps/api/          API Go, AI client, validasi, KB, SPA ter-embed
apps/web/          React, preprocessing Canvas, tampilan hasil, laporan PDF
packages/schema/   Kontrak JSON AnalysisResponse
docs/              Konsep, arsitektur, alur Usadha, sumber data, disclosure
```

## Data, Sumber, dan Kredit

- Foto sampel lontar: Eka343 (CC BY-SA 4.0) dan pindaian Carcan dari Internet Archive (PD-ID-exempt). Transkripsi: kontributor Wikisource Bali. Kredit, revisi, crop, dan hash ada di [`apps/web/public/samples/CREDITS.md`](apps/web/public/samples/CREDITS.md).
- Knowledge base Taru Pramana: `apps/api/internal/kb/data/taru-pramana.json`.
- Katalog berisi 1.316 entri lookup: 100 entri awal dan 1.216 tambahan berbasis Wikidata/GBIF. Alias ambigu tidak dipaksa cocok; Dapdap tetap pada tingkat genus. [Sumber, lisensi, dan batas kurasi](docs/KNOWLEDGE_BASE_SOURCES.md). Sumber kutipan tersedia pada setiap kartu arsip; pencarian literatur tidak membuktikan khasiat racikan.
- Rujukan naskah dan konteks: [`docs/DATA_SOURCES.md`](docs/DATA_SOURCES.md).
- Font: Cormorant Garamond dan Lora, dibundel lokal; lisensi OFL menyertai aset.
- Library frontend: React, Vite, Tailwind CSS, Radix UI, Lucide, class-variance-authority, clsx, tailwind-merge.
- Gateway AI: 9router, hanya Gemini 3.8 Flash Low tanpa fallback model.

## Keamanan dan Privasi

- API key hanya dibaca dari environment server.
- `.env` dan varian lokal di-ignore. Hanya `.env.example` yang masuk repository.
- Tidak ada data pribadi, data pasien, atau credential demo di repository.
- Backend memproses gambar di memori. Foto yang tidak cocok dengan arsip dikirim ke gateway dan penyedia AI; untuk arsip cocok, hanya teks rujukan dikirim. Retensi pihak ketiga belum diverifikasi. Browser menyimpan foto dan hasil sesi aktif di IndexedDB.
- Endpoint 9router tidak boleh dibuka ke publik tanpa autentikasi.
- Laporan kerentanan mengikuti [`SECURITY.md`](SECURITY.md).

## AI dan Technology Disclosure

Disclosure lengkap ada di [`docs/AI_DISCLOSURE.md`](docs/AI_DISCLOSURE.md). AI coding assistant membantu implementasi dan dokumentasi. Percakapan dan konfigurasi rahasia agen tidak disertakan; prompt produk dan keluaran eksperimen model disimpan untuk reproduksi. Keputusan teknis dan source code tetap menjadi tanggung jawab tim.

## Evaluasi AI

[Laporan evaluasi final](docs/AI_EVALUATION.md) mencatat perbandingan model, latensi, konsistensi, keluaran mentah, dan keterbatasan penilaian. Runner Go menggunakan jalur analisis aplikasi yang sebenarnya; tidak menyajikan jawaban cache.

Korpus lokal berisi 200 citra asli lontar Bali dari DeepLontar v2, dengan manifest, checksum, dan lisensi CC BY 4.0. Satu berkas anotasi dikarantina karena ID kelas di luar spesifikasi. Data ini dipakai melatih detektor teks dan pengenal glif lokal; bobot model bahasa tidak diubah. Lihat [sumber data](docs/DATA_SOURCES.md) dan [dokumen training OCR](docs/OCR_TRAINING.md).

### Hasil akhir dibanding baseline

Upgrade terakhir memperketat pengaman resep, bukan meningkatkan kemampuan OCR. Pada enam kasus batas terkontrol, resep tidak didukung yang lolos turun dari 6/6 menjadi 0/6; tiga kutipan utuh tetap diterima 3/3. Angka ini bukan akurasi pembacaan foto.

Pada batch empat gambar × dua pengulangan, baseline dan kandidat sama-sama menerima foto Bali 4/4 dan menolak pembanding 4/4. Pembanding Kawi adalah gambar cetak, bukan foto lontar; sampel kecil dan label aksara belum disertifikasi ahli. Tidak ada bukti peningkatan akurasi, konsistensi, atau kecepatan pembacaan. Bukti teknis terakhir: 51 tes backend lulus, build produksi berhasil, dan upload foto melalui API lokal menghasilkan HTTP 200. Rincian dan seluruh keluaran ada pada [evaluasi upgrade](docs/AI_EVALUATION.md#upgrade-pengaman-kutipan-dan-pembanding-baseline).

## Batasan Saat Ini

- Hasil AI dapat keliru atau berbeda antarpercobaan. Periksa bacaan dan resep terhadap sumber; hasil bukan edisi filologi terverifikasi atau rekomendasi pengobatan.
- Kualitas pembacaan bergantung pada foto. Deteksi aksara juga dapat salah menerima atau menolak gambar.
- Empat contoh arsip memakai transkripsi rujukan, bukan OCR baru. Pembacaan foto lain masih eksperimental.
- Basis pengetahuan tanaman terbatas pada entri Taru Pramana yang dikurasi.

## Tim

- **Andika Putra**: developer dan pemelihara repository.

## Dokumentasi Lanjutan

- [Indeks dokumentasi](docs/README.md)
- [Konsep produk](docs/PROJECT_CONCEPT.md)
- [Arsitektur dan alur data](docs/ARCHITECTURE.md)
- [Alur Usadha](docs/USADHA_MODE.md)
- [Sumber data dan kredit](docs/DATA_SOURCES.md)
- [AI dan Technology Disclosure](docs/AI_DISCLOSURE.md)
- [Evaluasi AI](docs/AI_EVALUATION.md)
- [Detektor teks dan pengenal glif](docs/OCR_TRAINING.md)
- [Proposal PDF](docs/proposal/AIWIC2026_0xAlpha_BukaLontar.pdf)
- [Panduan kontribusi](CONTRIBUTING.md)
- [Kebijakan keamanan](SECURITY.md)

---

Dibuat untuk **Bali DigiFest 2026, AI Web Innovation Challenge**.
