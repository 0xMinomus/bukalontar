# Architecture & Data Flow

Dokumen ini merangkum arsitektur MVP yang ada di repository, bukan rencana deployment yang belum diverifikasi.

## Runtime Overview

```text
Browser
  React 19 + Vite
  Canvas preprocessing
        |
        | multipart/form-data: image, optional source_id
        v
Go API (single binary)
  input validation + rate limit
  AI client
        |
        | POST /v1/chat/completions, JSON mode
        v
9router (OpenAI-compatible)
        |
        v
Gemini 3.8 Flash
        |
        | structured JSON
        v
Go validation + Taru Pramana KB grounding
        |
        v
AnalysisResponse -> React result view / print report
```

Frontend hasil build di-embed ke binary Go melalui `apps/api/internal/webui`. Satu proses menyajikan API dan SPA, sehingga browser tidak membutuhkan CORS pada deployment normal.

## Request Flow

1. **Client preprocessing.** `apps/web/src/lib/image.ts` mengecilkan gambar, menerapkan CLAHE, dan unsharp mask. File hasil dikirim sebagai `multipart/form-data`.
2. **Validation.** Handler membatasi ukuran upload (8 MB), tipe gambar, dan rate limit per menit. Tidak ada file yang ditulis ke disk.
3. **Inference.** `apps/api/internal/ai` mengirim image sebagai data URL ke endpoint OpenAI-compatible. Dua request independen dijalankan untuk membaca lembar yang padat.
4. **Structured parsing.** Respons model harus berupa JSON. Markdown fence dibersihkan sebelum parsing. Satu repair attempt dilakukan bila validasi gagal.
5. **Merge.** Resep dari kedua hasil digabung, dideduplikasi berdasarkan nama tanaman lokal, dan dibatasi maksimal enam resep. Field top-level memakai resep pertama untuk kompatibilitas kontrak lama.
6. **Grounding.** Bila nama tanaman cocok dengan Taru Pramana KB, nama ilmiah, family, senyawa, rujukan, dan energy character diambil dari KB. Bila tidak cocok, taxonomy dan sitasi dikosongkan agar tidak menjadi halusinasi.
7. **Response.** Disclaimer medis kanonik ditulis server. Frontend menampilkan seluruh resep, dan print report merender daftar yang sama.

## API Surface

| Method | Path | Keterangan |
| --- | --- | --- |
| `POST` | `/api/analyze-manuscript` | Upload gambar dan terima `AnalysisResponse` |
| `GET` | `/health` | Status proses, KB, dan keterjangkauan 9router |
| `GET` | `/api/health` | Alias health check |
| `GET` | `/` dan route SPA | Frontend hasil build |

Kontrak lengkap: `packages/schema/analysis.json`.

## Data Model

`AnalysisResponse` berisi:

- `manuscript_id` dan `source_reference`;
- `paleography` — deteksi aksara, confidence, transliterasi, terjemahan Indonesia/Inggris;
- `ethnomedicinal_entity` — field resep utama untuk kompatibilitas;
- `scientific_validation` — senyawa, aksi farmakologis, rujukan, dan disclaimer;
- `recipes[]` — daftar lengkap resep pada lembar, termasuk data tanaman dan validasi per resep.

## Reliability

- Retry pada 429, 5xx, timeout, dan transport error.
- Fallback provider opsional melalui `AI_FALLBACK_*`; tidak aktif secara default.
- Validasi schema di Go sebelum respons dikirim ke browser.
- Error API dikembalikan sebagai pesan yang dapat ditampilkan; tidak ada hasil palsu.
- Rate limit default 10 request/menit/klien.

## Configuration

Semua konfigurasi sensitif berasal dari environment:

```env
AI_BASE_URL=http://127.0.0.1:20128/v1
AI_MODEL=ag/gemini-3.8-flash-low
AI_API_KEY=
AI_TIMEOUT_SECONDS=120
AI_FALLBACK_ENABLED=false
AI_FALLBACK_BASE_URL=
AI_FALLBACK_MODEL=
AI_FALLBACK_API_KEY=
```

Contoh lengkap ada di `apps/api/.env.example`. File `.env` tidak boleh di-commit.

## Build and Deployment

```bash
npm install
make build          # frontend -> Go embed -> bin/bukalontar
./bin/bukalontar
```

Binary membaca environment, menyajikan SPA, dan memanggil 9router. Deployment publik harus menyediakan 9router yang dapat dijangkau oleh binary; endpoint gateway tidak boleh diekspos ke internet tanpa autentikasi.

## Testing

```bash
cd apps/api
go test ./...
go vet ./...
```

Frontend:

```bash
npm run lint
npm run build
```
