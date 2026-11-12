# Arsitektur dan Alur Data

Dokumen ini menjelaskan arsitektur MVP yang ada di repository. Fokusnya adalah perilaku kode saat ini pada dua mode deployment: binary mandiri dan Vercel Services.

Tag [`v1.1-final-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.1-final-submission) tetap menjadi snapshot lama. `main` menambahkan empat arsip berbantuan transkripsi rujukan. Frontend tetap mengirim foto asli; peningkatan kontras hanya pratinjau. Pengaman kutipan bukan verifikasi filologis.

## Gambaran Runtime

```text
Browser
  React 19 + Vite
  Canvas preview + IndexedDB sesi aktif
        |
        | multipart/form-data: image, optional source_id
        v
Go API
  validasi input + rate limit
  AI client
        |
        | POST /v1/chat/completions, JSON mode
        v
9router (OpenAI-compatible)
        |
        v
Model multimodal
        |
        | JSON terstruktur
        v
Validasi Go + grounding Taru Pramana
        |
        v
AnalysisResponse -> tampilan hasil / laporan PDF
```

Pada build mandiri, frontend hasil build di-embed ke binary Go melalui `apps/api/internal/webui`; satu proses menyajikan API dan SPA. Pada deployment Vercel, frontend Vite dan API Go berjalan sebagai dua service, `web` dan `api`, dengan routing yang sama di depan keduanya. Kontrak HTTP tetap sama pada kedua mode.

## Dua Alur, Satu Kontrak

Aplikasi memakai satu kontrak respons untuk dua kebutuhan:

| Alur | Kondisi | Isi respons |
| --- | --- | --- |
| Pembacaan lontar umum | Naskah tidak memuat resep Usadha | Transliterasi, terjemahan, sumber, ringkasan |
| Alur Usadha | Naskah memuat satu atau beberapa resep | Semua field umum, ditambah `recipes[]` dan validasi tanaman |

`recipes[]` boleh kosong. Field `ethnomedicinal_entity` dan `scientific_validation` tetap ada untuk kompatibilitas kontrak, tetapi nilainya memakai penanda "tidak disebutkan" jika naskah tidak memuat resep.

## Alur Request

1. **Pratinjau klien.** `apps/web/src/lib/image.ts` membuat tampilan CLAHE dan unsharp mask. Gambar asli dikirim ke API agar preprocessing yang belum tervalidasi tidak mengubah goresan untuk model.
2. **Validasi.** Handler membatasi ukuran upload 8 MB, tipe gambar, dan rate limit per menit. Tidak ada file yang ditulis ke disk.
3. **Inferensi.** SHA-256 byte foto dicocokkan dengan empat referensi tersemat. Kecocokan memilih transkripsi komunitas dan genre sumber; Gemini menerjemahkan/mengekstrak teks itu tanpa foto, gate visual, atau contoh glif. Latin sumber dipertahankan server. Jika tidak cocok, jalur lama tetap: klasifikasi foto target tanpa referensi (maksimal 256 token), keputusan `balinese/high`, lalu pembacaan dengan Flash Low dan bantuan `examples` default. Tidak ada pencocokan berdasarkan `source_id`, hasil analisis tersimpan, atau pergantian model.
4. **Parsing.** Respons model harus berupa JSON. Markdown fence dibersihkan sebelum parsing. Satu repair attempt dilakukan bila validasi gagal.
5. **Resep.** Maksimal enam resep harus memiliki `source_quote` pada satu bagian utuh bacaan daun pertama, menyebut nama tanaman dan bahan tambahan yang diklaim. Bacaan itu berasal dari AI untuk foto baru atau sumber komunitas untuk arsip cocok. Bagian tidak terbaca, anotasi kurung siku, elipsis, pengantar, dan daun lain tidak dapat menjembatani bukti. Arsip Carcan/Kakawin dibatasi tanpa resep. Pemeriksaan ini tidak memvalidasi terjemahan, tafsir racikan, atau keamanan medis.
6. **Grounding.** Nama ilmiah, famili, senyawa, rujukan, tautan literatur dan aksi farmakologis berasal dari KB saat tanaman cocok. Sifat tradisional mengikuti bacaan model, bukan ditimpa KB. Tanaman tak dikenal tidak mempertahankan sitasi atau klaim farmakologi dari model.
7. **Respons.** Disclaimer medis kanonik ditulis server. Resep kosong tetap kosong; informasi yang tidak disebutkan tidak otomatis menjadi obat minum. Klasifikasi naskah dan bahasa boleh tidak diketahui.

## API

| Method | Path | Keterangan |
| --- | --- | --- |
| `POST` | `/api/analyze-manuscript` | Upload gambar dan terima `AnalysisResponse` |
| `GET` | `/health` | Status proses, knowledge base, dan keterjangkauan 9router |
| `GET` | `/api/health` | Alias health check |
| `GET` | `/` dan route SPA | Frontend hasil build |

Kontrak lengkap ada di `packages/schema/analysis.json`.

## Model Data

`AnalysisResponse` berisi:

- `manuscript_id` dan `source_reference`;
- `paleography`: deteksi aksara, confidence, transliterasi, terjemahan Indonesia dan Inggris;
- `ethnomedicinal_entity`: field resep utama untuk kompatibilitas;
- `scientific_validation`: senyawa, aksi farmakologis, rujukan, dan disclaimer;
- `recipes[]`: daftar resep pada lembar, termasuk data tanaman dan validasi per resep.
- `reading_reference` opsional: metode `source_transcription`, SHA-256 foto, judul, genre katalog, URL revisi, lisensi, dan catatan editorial. Server mengisinya hanya dari registry setelah kecocokan byte dan membuang nilai provenance keluaran model. UI/PDF menampilkan metode ini, bukan confidence OCR.

## Keandalan

- Retry pada 429, 5xx, timeout, dan transport error.
- Server hanya menerima `AI_MODEL=ag/gemini-3.8-flash-low`; tidak ada fallback model.
- Validasi schema di Go sebelum respons dikirim ke browser.
- Error API dikembalikan sebagai pesan yang dapat ditampilkan. Tidak ada hasil palsu.
- Rate limit default 10 request per menit per klien.
- Backend tidak menulis file upload. Foto yang tidak cocok arsip diteruskan ke gateway; arsip cocok hanya mengirim transkripsi. Foto dan hasil sesi tersimpan di IndexedDB browser; retensi pihak ketiga belum diverifikasi.
- Pemeriksaan aksara Gemini terpisah berlaku pada foto yang tidak cocok registry. Keputusan non-Bali, campuran, atau meragukan tidak mencapai pembaca; keluaran gate rusak menghasilkan error. Pembaca yang menolak atau mengandung aksara lain dikosongkan. Ini bukan classifier independen. Identitas empat arsip berasal dari sumber dan hash, bukan keberhasilan gate.
- Bahasa teks dipisahkan dari aksara: bahasa Kawi dapat ditulis dalam aksara Bali. Transkripsi Unicode eksperimental tidak ditampilkan di UI/PDF karena belum layak berdasarkan evaluasi.

## Konfigurasi

Semua konfigurasi sensitif berasal dari environment:

```env
AI_BASE_URL=http://127.0.0.1:20128/v1
AI_MODEL=ag/gemini-3.8-flash-low
AI_REASONING_EFFORT=low
AI_READING_ASSISTANCE=examples
AI_API_KEY=
AI_TIMEOUT_SECONDS=120
```

Contoh lengkap ada di `apps/api/.env.example`. File `.env` tidak boleh di-commit.

## Build dan Deployment

```bash
npm install
make build
./bin/bukalontar
```

Untuk deployment publik, repository memakai Vercel Services:

- `web`: Vite static frontend di `apps/web`
- `api`: Go server di `apps/api`
- URL: `https://bukalontar.my.id`
- Health check: `https://bukalontar.my.id/health`

`AI_BASE_URL` di Vercel harus mengarah ke endpoint publik 9router, bukan `127.0.0.1`. Endpoint gateway tidak boleh diekspos ke internet tanpa autentikasi.

## Pengujian

```bash
cd apps/api
go test ./...
go vet ./...

cd ../..
npm run lint
npm run build
```

Evaluasi model menggunakan `go run ./cmd/evaluate` dari `apps/api`, atau jalankan dari root dengan `go -C apps/api run ./cmd/evaluate` dan path gambar relatif ke `apps/api`. Detail perintah, hasil mentah, serta batasan training ada di [AI_EVALUATION.md](AI_EVALUATION.md).
