# Contributing

Terima kasih sudah ingin membantu BukaLontar. Repository ini adalah MVP kompetisi, jadi perubahan harus menjaga fungsi utama tetap dapat dijalankan dan diuji.

## Prasyarat

- Node.js 20 atau lebih baru
- npm
- Go 1.24 atau lebih baru

## Menyiapkan Lingkungan

```bash
npm install
cp apps/api/.env.example apps/api/.env
# Isi AI_BASE_URL, AI_MODEL, dan AI_API_KEY.
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

## Dokumentasi

Dokumentasi utama ada di `docs/`. Perbarui halaman yang terkait jika mengubah perilaku produk, arsitektur, atau sumber data.

## Keamanan

Jangan melaporkan masalah keamanan di issue publik. Ikuti [`SECURITY.md`](SECURITY.md).
