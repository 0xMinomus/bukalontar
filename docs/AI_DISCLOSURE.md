# AI dan Technology Disclosure

Dokumen ini melengkapi Lampiran C panduan peserta. Isinya harus sesuai dengan implementasi yang ada di repository.

Snapshot final submission: [`v1.1-final-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.1-final-submission), dokumentasi 12 November 2026. Riwayat eksperimen tetap tersedia, tetapi model produk tidak diganti.

| Komponen | Keterangan |
| --- | --- |
| **Nama Produk** | BukaLontar |
| **Framework / Frontend** | React 19, TypeScript, Vite, Tailwind CSS v4, Radix UI, Lucide, class-variance-authority, clsx, tailwind-merge |
| **Backend / Runtime** | Go standard library (`net/http`); pada Vercel, API berjalan sebagai service `api`; binary mandiri dengan SPA ter-embed tersedia untuk deployment non-Vercel |
| **Database / Storage** | Tidak ada database eksternal. KB JSON Taru Pramana di-embed dengan `go:embed`. IndexedDB browser menyimpan foto asli, pratinjau, nama berkas, dan hasil sesi aktif |
| **Hosting / Deployment** | Vercel Services di `https://bukalontar.my.id`, service `web` dan `api`; binary mandiri opsional |
| **AI Gateway** | 9router, endpoint OpenAI-compatible |
| **AI Model / API / Method** | Hanya `ag/gemini-3.8-flash-low` melalui 9router, tanpa fallback model. Bantuan lokal: detektor teks ridge 10 bobot dan pengenal glif HOG + LinearSVC, keduanya dilatih dari anotasi DeepLontar dan di-embed |
| **Peran AI dalam Produk** | Membaca citra lontar beraksara Bali, menghasilkan transliterasi dan terjemahan, serta mengekstrak resep ketika didukung bacaan Usadha. Bahasa Kawi/Jawa Kuno bukan aksara Kawi; aksara Jawa dan Kawi di luar lingkup. Klasifikasi masih bergantung pada model umum, bukan model yang telah kehilangan kemampuan membaca aksara lain |
| **Dataset / Data Source** | Taru Pramana KB; foto Wikimedia Commons CC BY-SA 4.0; 200 citra asli DeepLontar v2 CC BY 4.0, dengan satu berkas anotasi dikarantina; lima halaman Wikisource sebagai holdout komunitas, bukan validasi ahli independen. Segmen buku NusaAksara hanya dipakai diagnostik, bukan data produk atau bukti akurasi foto lontar. Detail di DATA_SOURCES.md |
| **Open-source Libraries** | React, Vite, Tailwind CSS, Radix UI, Lucide, class-variance-authority, clsx, tailwind-merge; Go standard library |
| **Template / Boilerplate** | Vite React TypeScript starter; komponen UI shadcn-style |
| **AI Coding Assistant** | Command Code dan OpenAI Codex digunakan untuk membantu implementasi, evaluasi, serta dokumentasi. Prompt percakapan dan log agen tidak disimpan di repository. Prompt produk dan hasil evaluasi model disertakan untuk reproduksi. Keputusan teknis dan source code tetap menjadi tanggung jawab tim |
| **Pre-existing Components** | Tidak ada produk komersial pre-existing. Repository ini dibangun untuk Bali DigiFest 2026 |
| **Third-party Assets / Credits** | Foto lontar: Eka343, CC BY-SA 4.0 (`apps/web/public/samples/CREDITS.md`); Cormorant Garamond dan Lora dibundel lokal dengan lisensi OFL |

## Status Evaluasi dan Pelatihan

Eksperimen historis mencakup Sol dan beberapa varian Gemini sebelum keputusan model akhir. Server produk sekarang hanya menerima `ag/gemini-3.8-flash-low`. Dukungan vision pada katalog model tidak menjamin pembacaan lontar; hasil historis dan keputusan pemilihan dicatat dalam [AI_EVALUATION.md](AI_EVALUATION.md).

Ada pelatihan model lokal, tetapi bukan fine-tuning model bahasa. `scripts/train_lontar_ocr.py` memakai 159 citra train dari 80 kelompok daun fisik; dev dan test masing-masing 20 citra dari 10 kelompok. Kedua sisi daun berada pada kelompok yang sama, dan satu label train dikarantina. Model lokal mendeteksi okupansi teks dan mengenali potongan glif ke ID kelas, bukan membaca kalimat atau menerjemahkan halaman.

Tidak ada fine-tuning model bahasa: endpoint `GET /v1/fine_tuning/jobs` pada router mengembalikan HTTP 404. Bobot Gemini/Sol tidak diubah. Metrik, batasan, dan cara melatih ulang ada di [OCR_TRAINING.md](OCR_TRAINING.md).

Transkripsi Unicode Bali tetap tersedia sebagai keluaran API eksperimental, tetapi tidak ditampilkan di UI atau PDF: CER pada evaluasi sebelumnya melebihi 0,95. Menampilkan karakter Unicode Bali tidak membuktikan bacaannya benar. Pembatasan produk memakai satu panggilan klasifikasi aksara tanpa referensi/katalog sebelum satu pembacaan, lalu normalisasi penolakan server. Keduanya menggunakan Gemini 3.8 Flash Low. Ini menambah satu panggilan pada foto yang diterima, bukan fine-tuning atau jaminan deteksi sempurna.

Upgrade terakhir memisahkan kutipan resep yang utuh dari bagian tidak terbaca, elipsis, pengantar, dan lembar lain sebelum grounding. Enam kasus batas terkontrol yang dahulu meloloskan resep kini ditolak, sementara tiga kasus kutipan utuh tetap diterima. Ini bukan ukuran akurasi foto, terjemahan, atau resep. OCR literal terpisah dan pemeriksa visual tambahan diuji tetapi tidak diintegrasikan karena belum cukup andal; bobot dan model produk tidak berubah. Lihat [perbandingan baseline terakhir](AI_EVALUATION.md#upgrade-pengaman-kutipan-dan-pembanding-baseline).


## Pernyataan

Kami menyatakan bahwa informasi di atas benar dan lengkap. Kami bertanggung jawab atas penggunaan teknologi, source code, data, library, API, model AI, dan aset yang digunakan pada produk.

## Catatan Keamanan

- API key hanya dibaca dari environment server.
- `.env` dan varian lokal di-ignore. `.env.example` tidak berisi nilai rahasia.
- Tidak ada credential demo, token, atau data pribadi di repository.
- Endpoint 9router tidak boleh diekspos ke publik tanpa autentikasi.
- Foto dikirim ke gateway dan penyedia AI. Backend tidak menyimpan file upload, tetapi retensi pihak ketiga belum diverifikasi.
- Hasil confidence model bukan probabilitas kebenaran. Pengguna harus memeriksa hasil terhadap foto dan rujukan filologi.
