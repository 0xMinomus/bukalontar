# AI dan Technology Disclosure

Dokumen ini melengkapi Lampiran C panduan peserta. Isinya harus sesuai dengan implementasi yang ada di repository.

Tag [`v1.1-final-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.1-final-submission) adalah snapshot yang tidak dipindah. Dokumen `main` juga mencatat analisis berbantuan transkripsi rujukan yang ditambahkan setelah tag pada 12 November 2026. Model produk tetap sama.

| Komponen | Keterangan |
| --- | --- |
| **Nama Produk** | BukaLontar |
| **Framework / Frontend** | React 19, TypeScript, Vite, Tailwind CSS v4, Radix UI, Lucide, class-variance-authority, clsx, tailwind-merge |
| **Backend / Runtime** | Go standard library (`net/http`); pada Vercel, API berjalan sebagai service `api`; binary mandiri dengan SPA ter-embed tersedia untuk deployment non-Vercel |
| **Database / Storage** | Tidak ada database eksternal. KB JSON Taru Pramana di-embed dengan `go:embed`. IndexedDB browser menyimpan foto asli, pratinjau, nama berkas, dan hasil sesi aktif |
| **Hosting / Deployment** | Vercel Services di `https://bukalontar.my.id`, service `web` dan `api`; binary mandiri opsional |
| **AI Gateway** | 9router, endpoint OpenAI-compatible |
| **AI Model / API / Method** | Hanya `ag/gemini-3.8-flash-low` melalui 9router, tanpa fallback model. Bantuan lokal: detektor teks ridge 10 bobot dan pengenal glif HOG + LinearSVC, keduanya dilatih dari anotasi DeepLontar dan di-embed |
| **Peran AI dalam Produk** | Foto baru: membaca citra lontar beraksara Bali, menerjemahkan, dan mengekstrak isi eksperimental. Empat foto arsip yang cocok SHA-256: menerjemahkan dan mengekstrak transkripsi komunitas yang dipasok; bacaan Latin bukan hasil OCR baru. UI/PDF menyebut metode ini. Bahasa Kawi/Jawa Kuno bukan aksara Kawi; pembatasan foto baru tetap bergantung pada Gemini umum |
| **Dataset / Data Source** | KB Taru Pramana; foto Eka343 CC BY-SA 4.0 dan Carcan Internet Archive PD-ID-exempt; empat transkripsi/render Latin Wikisource dengan revisi tersemat dan whitespace editorial; 200 citra asli DeepLontar v2 CC BY 4.0; lima halaman holdout komunitas. Referensi yang dipasok ke model tidak dihitung sebagai holdout atau bukti akurasi OCR. Detail di DATA_SOURCES.md |
| **Open-source Libraries** | React, Vite, Tailwind CSS, Radix UI, Lucide, class-variance-authority, clsx, tailwind-merge; Go standard library |
| **Template / Boilerplate** | Vite React TypeScript starter; komponen UI shadcn-style |
| **AI Coding Assistant** | Command Code dan OpenAI Codex digunakan untuk membantu implementasi, evaluasi, serta dokumentasi. Prompt percakapan dan log agen tidak disimpan di repository. Prompt produk dan hasil evaluasi model disertakan untuk reproduksi. Keputusan teknis dan source code tetap menjadi tanggung jawab tim |
| **Pre-existing Components** | Tidak ada produk komersial pre-existing. Repository ini dibangun untuk Bali DigiFest 2026 |
| **Third-party Assets / Credits** | Foto: Eka343 (CC BY-SA 4.0), Carcan dari Internet Archive (PD-ID-exempt). Transkripsi: kontributor Wikisource Bali (CC BY-SA 4.0), dengan perubahan whitespace yang dicatat. Kredit lengkap di `apps/web/public/samples/CREDITS.md`; font Cormorant Garamond/Lora berlisensi OFL |

## Status Evaluasi dan Pelatihan

Eksperimen historis mencakup Sol dan beberapa varian Gemini sebelum keputusan model akhir. Server produk sekarang hanya menerima `ag/gemini-3.8-flash-low`. Dukungan vision pada katalog model tidak menjamin pembacaan lontar; hasil historis dan keputusan pemilihan dicatat dalam [AI_EVALUATION.md](AI_EVALUATION.md).

Ada pelatihan model lokal, tetapi bukan fine-tuning model bahasa. `scripts/train_lontar_ocr.py` memakai 159 citra train dari 80 kelompok daun fisik; dev dan test masing-masing 20 citra dari 10 kelompok. Kedua sisi daun berada pada kelompok yang sama, dan satu label train dikarantina. Model lokal mendeteksi okupansi teks dan mengenali potongan glif ke ID kelas, bukan membaca kalimat atau menerjemahkan halaman.

Tidak ada fine-tuning model bahasa: endpoint `GET /v1/fine_tuning/jobs` pada router mengembalikan HTTP 404. Bobot Gemini/Sol tidak diubah. Metrik, batasan, dan cara melatih ulang ada di [OCR_TRAINING.md](OCR_TRAINING.md).

Transkripsi Unicode Bali tetap eksperimental dan tidak ditampilkan di UI/PDF karena hasil sebelumnya belum layak. Foto baru memakai klasifikasi target tanpa referensi sebelum pembacaan Gemini. Empat foto arsip yang cocok persis diidentifikasi dari registry sumber dan langsung memakai teks rujukan, bukan klasifikasi visual baru; tidak ada pengecualian berdasarkan nama file atau `source_id`. Terjemahan dan ekstraksi tetap dipanggil langsung ke Flash Low, tanpa hasil analisis yang di-cache.

Upgrade terakhir memisahkan kutipan resep yang utuh dari bagian tidak terbaca, elipsis, pengantar, dan lembar lain sebelum grounding. Enam kasus batas terkontrol yang dahulu meloloskan resep kini ditolak, sementara tiga kasus kutipan utuh tetap diterima. Ini bukan ukuran akurasi foto, terjemahan, atau resep. OCR literal terpisah dan pemeriksa visual tambahan diuji tetapi tidak diintegrasikan karena belum cukup andal; bobot dan model produk tidak berubah. Lihat [perbandingan baseline terakhir](AI_EVALUATION.md#upgrade-pengaman-kutipan-dan-pembanding-baseline).

Perbaikan contoh setelah tag menambahkan dua Usadha, Carcan Kucing, dan Kakawin Bharatayuda. Uji foto murni masih mengarang bacaan/resep; mengganti foto saja tidak cukup. Pengguna menyetujui bantuan transkripsi yang dilabeli terbuka. Pengaman kutipan sekarang memeriksa bacaan sumber untuk arsip cocok, bukan bacaan model yang berpotensi dikarang. Ini tidak memvalidasi terjemahan, tafsir resep, atau foto unggahan lain.


## Pernyataan

Kami menyatakan bahwa informasi di atas benar dan lengkap. Kami bertanggung jawab atas penggunaan teknologi, source code, data, library, API, model AI, dan aset yang digunakan pada produk.

## Catatan Keamanan

- API key hanya dibaca dari environment server.
- `.env` dan varian lokal di-ignore. `.env.example` tidak berisi nilai rahasia.
- Tidak ada credential demo, token, atau data pribadi di repository.
- Endpoint 9router tidak boleh diekspos ke publik tanpa autentikasi.
- Foto baru dikirim ke gateway/penyedia AI; arsip yang cocok hanya mengirim teks rujukan. Backend tidak menyimpan file upload. Retensi pihak ketiga belum diverifikasi.
- Hasil confidence model bukan probabilitas kebenaran. Pengguna harus memeriksa hasil terhadap foto dan rujukan filologi.
