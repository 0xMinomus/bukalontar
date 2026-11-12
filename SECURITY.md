# Security Policy

## Versi yang Didukung

Repository ini adalah MVP untuk Bali DigiFest 2026. Perbaikan keamanan diterapkan pada branch `main`. Tag [`v1.1-final-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.1-final-submission) adalah snapshot submission, bukan sertifikasi keamanan atau jaminan akurasi.

## Melaporkan Masalah

Jangan membuka issue publik untuk masalah yang memuat API key, token, data pribadi, atau detail eksploitasi.

Gunakan GitHub Security Advisory pada repository ini, atau hubungi pemelihara melalui kanal privat yang tercantum di profil GitHub.

## Rahasia dan Kredensial

- API key hanya boleh berada di environment server atau file `.env` lokal.
- `.env` dan varian lokal di-ignore.
- `.env.example` hanya berisi nama variabel, tanpa nilai rahasia.
- Jangan menempelkan API key ke issue, pull request, screenshot, atau log.
- Jika sebuah key pernah masuk ke repository, rotasi key tersebut dan bersihkan riwayat Git.

## Data dan Privasi

- Jangan mengunggah data pribadi, data pasien, atau naskah yang tidak memiliki izin publikasi.
- Backend memproses gambar di memori dan meneruskannya melalui 9router ke penyedia AI. Kebijakan retensi gateway dan penyedia belum diverifikasi; jangan mengunggah materi rahasia.
- Foto asli, pratinjau yang diperjelas, nama berkas, dan hasil sesi aktif dapat tersimpan di IndexedDB browser. Tombol naskah baru atau pembersihan data situs menghapus sesi lokal; ini tidak membuktikan penghapusan pada layanan pihak ketiga.

## Batas AI

- BukaLontar bukan alat diagnosis dan bukan sumber nasihat medis.
- Hasil model dapat salah membaca aksara atau istilah.
- Jangan gunakan hasil AI sebagai satu-satunya dasar untuk keputusan pengobatan atau filologi.
- Pengaman resep memeriksa konsistensi kutipan dan batas lembar, bukan kebenaran bacaan atau keamanan racikan. Bacaan beserta kutipannya masih dapat sama-sama dikarang model.
- Validasi input dan pembatasan model bukan bukti audit keamanan menyeluruh. Endpoint health tidak mengungkap commit deployment.

## Endpoint AI

- Endpoint 9router tidak boleh diekspos ke internet tanpa autentikasi.
- Gunakan jaringan privat, tunnel dengan autentikasi, atau pembatasan akses.
- Batasi rate limit dan pantau penggunaan API.
