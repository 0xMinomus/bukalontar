# Security Policy

## Versi yang Didukung

Repository ini adalah MVP untuk Bali DigiFest 2026. Perbaikan keamanan diterapkan pada branch `main`.

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
- Gambar diproses untuk analisis dan tidak disimpan sebagai file oleh backend.
- Hasil analisis dapat tersimpan di IndexedDB browser pengguna. Pengguna dapat menghapusnya dengan membersihkan data situs.

## Batas AI

- BukaLontar bukan alat diagnosis dan bukan sumber nasihat medis.
- Hasil model dapat salah membaca aksara atau istilah.
- Jangan gunakan hasil AI sebagai satu-satunya dasar untuk keputusan pengobatan atau filologi.

## Endpoint AI

- Endpoint 9router tidak boleh diekspos ke internet tanpa autentikasi.
- Gunakan jaringan privat, tunnel dengan autentikasi, atau pembatasan akses.
- Batasi rate limit dan pantau penggunaan API.
