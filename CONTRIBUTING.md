# Contributing

Repository ini adalah alat bantu pembacaan lontar eksperimental untuk Bali DigiFest 2026. Snapshot final submission tersedia pada tag [`v1.1-final-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.1-final-submission); pengembangan setelah submission tetap melalui branch `main`.

## Prasyarat

- Node.js 22.12 atau lebih baru (atau Node.js 20.19 pada seri 20.x)
- npm
- Go 1.26 atau lebih baru

## Menyiapkan Lingkungan

```bash
npm install
cp apps/api/.env.example apps/api/.env
# Isi AI_BASE_URL dan AI_API_KEY; pertahankan AI_MODEL=ag/gemini-3.8-flash-low.
```

Jalankan API dan frontend:

```bash
make dev-api
npm run dev
```

## Pemeriksaan Sebelum Pull Request

```bash
cd apps/api
go test ./...
go vet ./...

cd ../..
npm run lint
npm run build
```

## Aturan Perubahan

- Jangan menambahkan API key, token, atau data pribadi.
- Jangan mengubah kontrak `AnalysisResponse` tanpa memperbarui schema dan frontend.
- Jangan menghapus test yang masih relevan.
- Jaga perilaku naskah non-Usadha agar tetap dapat dibaca.
- Tulis pesan commit yang menjelaskan alasan perubahan.
- Sertakan bukti pengujian untuk perubahan perilaku.
- Pertahankan foto asli sebagai input AI; peningkatan kontras browser hanya pratinjau.
- Pengaman kutipan tidak boleh melewati bagian tidak terbaca, elipsis, pengantar, atau lembar lain. Resep yang ditolak tidak boleh muncul lagi saat normalisasi/grounding.
- Perubahan model atau bantuan pembacaan memerlukan pembanding pada foto yang sama. Respons HTTP 200, confidence tinggi, dan kecocokan kutipan bukan bukti akurasi.

## Dokumentasi

Dokumentasi utama ada di `docs/`. Perbarui halaman yang terkait jika mengubah perilaku produk, arsitektur, atau sumber data.

Perbarui [status](docs/STATUS.md), [evaluasi](docs/AI_EVALUATION.md), dan [disclosure](docs/AI_DISCLOSURE.md) sesuai bukti. Simpan eksperimen lama sebagai riwayat, bukan klaim kemampuan versi sekarang. Tag submission yang sudah dipublikasikan tidak dipindahkan atau ditimpa; buat tag baru untuk snapshot berikutnya.

## Keamanan

Jangan melaporkan masalah keamanan di issue publik. Ikuti [`SECURITY.md`](SECURITY.md).
