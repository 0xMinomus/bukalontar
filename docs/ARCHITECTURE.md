# Arsitektur dan Alur Data

Dokumen ini menjelaskan arsitektur MVP yang ada di repository. Fokusnya adalah perilaku kode saat ini pada dua mode deployment: binary mandiri dan Vercel Services.

## Gambaran Runtime

```text
Browser
  React 19 + Vite
  Canvas preprocessing
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
3. **Inferensi.** `apps/api/internal/ai` terlebih dahulu mengklasifikasikan aksara pada foto target saja, tanpa contoh glif atau katalog tanaman (maksimal 256 token). Hanya keputusan `balinese/high` melanjutkan satu pembacaan dengan Gemini 3.8 Flash Low. Default `examples` menambahkan tiga contoh glif ahli; crop lokal hanya pada mode eksperimental. Kedua panggilan memakai model yang sama dan berbagi tenggat permintaan. Temperature 0 tidak menjamin determinisme.
4. **Parsing.** Respons model harus berupa JSON. Markdown fence dibersihkan sebelum parsing. Satu repair attempt dilakukan bila validasi gagal.
5. **Resep.** Maksimal enam resep dipertahankan bila `source_quote` cocok dengan satu bagian utuh transliterasi daun pertama, menyebut nama tanaman, dan menyertakan bahan tambahan yang diklaim. Tanda bagian tidak terbaca, anotasi dalam kurung siku, dan elipsis menjadi batas; kata di kedua sisinya tidak digabungkan untuk membenarkan kutipan. Jika ada label lembar, bukti harus berasal dari bagian setelah label pertama bernomor 1, bukan pengantar atau lembar lain. Ini pemeriksaan konsistensi teks, bukan validasi resep atau bukti bahwa bacaan AI sesuai foto.
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

## Keandalan

- Retry pada 429, 5xx, timeout, dan transport error.
- Server hanya menerima `AI_MODEL=ag/gemini-3.8-flash-low`; tidak ada fallback model.
- Validasi schema di Go sebelum respons dikirim ke browser.
- Error API dikembalikan sebagai pesan yang dapat ditampilkan. Tidak ada hasil palsu.
- Rate limit default 10 request per menit per klien.
- Pemeriksaan aksara adalah panggilan Gemini terpisah tanpa contoh glif atau katalog. Keputusan non-Bali, campuran, atau meragukan tidak mencapai pembaca; keluaran pemeriksaan rusak menghasilkan error, bukan izin membaca. Server juga mengosongkan seluruh bacaan/resep ketika pembaca menolak atau transkripsi/transliterasi mengandung aksara mirip non-Bali. Kedua pemeriksaan tetap bergantung pada model umum, bukan classifier terlatih independen.
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
