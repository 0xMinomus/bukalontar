# AI dan Technology Disclosure

Dokumen ini melengkapi Lampiran C panduan peserta. Isinya harus sesuai dengan implementasi yang ada di repository.

| Komponen | Keterangan |
| --- | --- |
| **Nama Produk** | BukaLontar |
| **Framework / Frontend** | React 19, TypeScript, Vite, Tailwind CSS v4, Radix UI, Lucide, class-variance-authority, clsx, tailwind-merge |
| **Backend / Runtime** | Go standard library (`net/http`); pada Vercel, API berjalan sebagai service `api`; binary mandiri dengan SPA ter-embed tersedia untuk deployment non-Vercel |
| **Database / Storage** | Tidak ada database eksternal. Knowledge base JSON Taru Pramana di-embed dengan `go:embed`. Sesi frontend disimpan di IndexedDB browser |
| **Hosting / Deployment** | Vercel Services di `https://bukalontar.my.id`, service `web` dan `api`; binary mandiri opsional |
| **AI Gateway** | 9router, endpoint OpenAI-compatible |
| **AI Model / API / Method** | Hanya `ag/gemini-3.8-flash-low` melalui 9router, tanpa fallback model. Bantuan lokal: detektor teks ridge 10 bobot dan pengenal glif HOG + LinearSVC, keduanya dilatih dari anotasi DeepLontar dan di-embed |
| **Peran AI dalam Produk** | Membaca citra lontar beraksara Bali, menghasilkan transliterasi dan terjemahan, serta mengekstrak resep ketika didukung bacaan Usadha. Bahasa Kawi/Jawa Kuno bukan aksara Kawi; aksara Jawa dan Kawi di luar lingkup. Klasifikasi masih bergantung pada model umum, bukan model yang telah kehilangan kemampuan membaca aksara lain |
| **Dataset / Data Source** | Taru Pramana KB; foto Wikimedia Commons CC BY-SA 4.0; 200 citra asli DeepLontar v2 CC BY 4.0 dengan label kotak karakter untuk pelatihan detektor lokal; satu berkas anotasi dikarantina; lima halaman Wikisource berlabel tangan sebagai holdout evaluasi. Tidak ada data pribadi. Detail di DATA_SOURCES.md |
| **Open-source Libraries** | React, Vite, Tailwind CSS, Radix UI, Lucide, class-variance-authority, clsx, tailwind-merge; Go standard library |
| **Template / Boilerplate** | Vite React TypeScript starter; komponen UI shadcn-style |
| **AI Coding Assistant** | Command Code dan OpenAI Codex digunakan untuk membantu implementasi, evaluasi, serta dokumentasi. Prompt percakapan dan log agen tidak disimpan di repository. Prompt produk dan hasil evaluasi model disertakan untuk reproduksi. Keputusan teknis dan source code tetap menjadi tanggung jawab tim |
| **Pre-existing Components** | Tidak ada produk komersial pre-existing. Repository ini dibangun untuk Bali DigiFest 2026 |
| **Third-party Assets / Credits** | Foto lontar: Eka343, CC BY-SA 4.0 (`apps/web/public/samples/CREDITS.md`); font Cormorant Garamond dan Lora dari Google Fonts |

## Status Evaluasi dan Pelatihan

`cx/gpt-6.1-sol` diuji melalui 9router dengan input gambar dan reasoning rendah; pengujian diagnostik reasoning tinggi juga dilakukan. Dukungan vision pada daftar model tidak menjamin kemampuan membaca lontar. Hasil dan keputusan pemilihan model dicatat dalam [AI_EVALUATION.md](AI_EVALUATION.md).

Ada pelatihan model lokal, tetapi bukan fine-tuning model bahasa. `scripts/train_lontar_ocr.py` melatih detektor okupansi teks dan pengenal glif dari 159 lembar train. Kedua artefak di-embed sebagai bantuan pembacaan; ID kelas glif tetap opasitas dan tidak pernah dikeluarkan sebagai teks.

Tidak ada fine-tuning model bahasa: endpoint `GET /v1/fine_tuning/jobs` pada router mengembalikan HTTP 404. Bobot Gemini/Sol tidak diubah. Metrik, batasan, dan cara melatih ulang ada di [OCR_TRAINING.md](OCR_TRAINING.md).

Transkripsi Unicode Bali tetap tersedia sebagai keluaran API eksperimental, tetapi tidak ditampilkan di UI atau PDF: CER pada evaluasi sebelumnya melebihi 0,95. Menampilkan karakter Unicode Bali tidak membuktikan bacaannya benar. Pembatasan produk memakai satu panggilan klasifikasi aksara tanpa referensi/katalog sebelum satu pembacaan, lalu normalisasi penolakan server. Keduanya menggunakan Gemini 3.8 Flash Low. Ini menambah satu panggilan pada foto yang diterima, bukan fine-tuning atau jaminan deteksi sempurna.


## Pernyataan

Kami menyatakan bahwa informasi di atas benar dan lengkap. Kami bertanggung jawab atas penggunaan teknologi, source code, data, library, API, model AI, dan aset yang digunakan pada produk.

## Catatan Keamanan

- API key hanya dibaca dari environment server.
- `.env` dan varian lokal di-ignore. `.env.example` tidak berisi nilai rahasia.
- Tidak ada credential demo, token, atau data pribadi di repository.
- Endpoint 9router tidak boleh diekspos ke publik tanpa autentikasi.
