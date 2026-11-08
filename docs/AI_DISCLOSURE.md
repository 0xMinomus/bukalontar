# AI & Technology Disclosure

Dokumen ini melengkapi Lampiran C panduan peserta. Isi harus sesuai dengan implementasi yang ada di repository.

| Komponen | Keterangan |
| --- | --- |
| **Nama Produk** | BukaLontar |
| **Framework / Frontend** | React 19, TypeScript, Vite, Tailwind CSS v4, Radix UI, Lucide, class-variance-authority, clsx, tailwind-merge |
| **Backend / Runtime** | Go standard library (`net/http`), single binary dengan SPA ter-embed |
| **Database / Storage** | Tidak ada database eksternal. Knowledge base JSON Taru Pramana di-embed dengan `go:embed`; session frontend disimpan di IndexedDB browser |
| **Hosting / Deployment** | Belum ada URL publik pada commit ini. Build produksi berjalan sebagai satu binary Go dan dapat dijalankan lokal di port 8080 |
| **AI Gateway** | 9router, endpoint OpenAI-compatible |
| **AI Model / API / Method** | Gemini 3.8 Flash multimodal vision (`ag/gemini-3.8-flash-low`) melalui 9router; fallback opsional `google/gemini-3.8-flash` melalui Command Code Provider API |
| **Peran AI dalam Produk** | Membaca aksara Bali pada citra lontar, menghasilkan transliterasi dan terjemahan Kawi, lalu mengekstrak entitas tanaman, bagian tanaman, formula, dan indikasi penyakit. AI adalah komponen material alur utama, bukan sekadar asisten coding |
| **Dataset / Data Source** | Taru Pramana KB yang dikurasi di `apps/api/internal/kb/data/taru-pramana.json`; foto sampel lontar dari Wikimedia Commons dengan lisensi CC BY-SA 4.0. Tidak ada data pribadi atau data pasien |
| **Open-source Libraries** | React, Vite, Tailwind CSS, Radix UI, Lucide, class-variance-authority, clsx, tailwind-merge; Go standard library |
| **Template / Boilerplate** | Vite React TypeScript starter; komponen UI shadcn-style |
| **AI Coding Assistant** | Command Code / omp coding assistant digunakan untuk membantu implementasi dan dokumentasi. Prompt dan log agen tidak disimpan di repository; source code ditinjau dan diuji oleh tim |
| **Pre-existing Components** | Tidak ada produk komersial pre-existing. Repository ini dibangun untuk Bali DigiFest 2026 |
| **Third-party Assets / Credits** | Foto lontar: Eka343, CC BY-SA 4.0 (`apps/web/public/samples/CREDITS.md`); font Cormorant Garamond dan Lora dari Google Fonts |

## Pernyataan

Kami menyatakan bahwa informasi di atas benar dan lengkap. Kami bertanggung jawab atas penggunaan teknologi, source code, data, library, API, model AI, dan aset yang digunakan pada produk.

## Catatan Keamanan

- API key hanya dibaca dari environment server.
- `.env` dan varian lokal di-ignore; `.env.example` tidak berisi nilai rahasia.
- Tidak ada credential demo, token, atau data pribadi di repository.
- Endpoint 9router tidak boleh diekspos ke publik tanpa autentikasi.
