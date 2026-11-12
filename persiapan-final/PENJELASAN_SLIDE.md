# Penjelasan Slide

## Slide 1: Alur kerja AI

Setelah melihat hasilnya, sekarang kita lihat proses di balik BukaLontar. Browser mengirim foto asli ke API Go untuk validasi dan pembatasan permintaan. Foto baru diperiksa aksaranya lalu dibaca Gemini 3.8 Flash Low melalui 9router, dengan bantuan contoh glif. Detektor lokal tersedia sebagai bantuan eksperimental, bukan OCR halaman. Empat foto arsip yang cocok SHA-256 memakai transkripsi Wikisource, bukan OCR baru. Setelah AI menerjemahkan dan mengekstrak isinya, backend memeriksa JSON dan kutipan resep, mencocokkan tanaman dengan Taru Pramana, lalu mengirim hasil ke browser untuk ditampilkan dan dicetak sebagai PDF.

## Slide 2: Improvement sejak Tahap 1

Pada tahap final, kami memperketat kutipan resep agar tidak melompati bagian tak terbaca atau lembar lain. Enam kasus batas terkontrol yang sebelumnya lolos kini ditolak. Ini bukti pengaman kutipan, bukan peningkatan akurasi OCR. Rujukan tanaman juga dikurasi, metode pembacaan dilabeli pada hasil dan PDF, serta contoh arsip diperbarui menjadi dua Usadha dan dua non-Usadha.

## Slide 3: Tech stack

Untuk membangun alur tersebut, kami memakai React, TypeScript, Vite, dan Tailwind CSS pada frontend, serta Go pada backend. Modelnya Gemini 3.8 Flash Low melalui 9router. Aplikasi berjalan di Vercel, dengan basis pengetahuan Taru Pramana dalam format JSON.

## Slide 4: Keterbatasan

Meski pengamannya diperketat, bacaan masih bisa keliru, hasil belum konsisten, dan deteksi aksara belum sempurna. Analisis juga bergantung pada internet, gateway, dan kuota AI. Kami tidak melakukan fine-tuning Gemini. Karena itu, hasil tetap perlu diperiksa ahli, dan informasi Usadha bukan anjuran pengobatan.
