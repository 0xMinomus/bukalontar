# Konsep Produk: BukaLontar

**Pembaca lontar Bali berbasis AI, dengan alur khusus untuk naskah Usadha.**

## 1. Masalah

Bali memiliki salah satu warisan naskah daun lontar terbesar di Asia Tenggara. Data yang sering dirujuk:

- Dinas Kebudayaan Provinsi Bali mencatat lebih dari 25.000 cakep lontar milik masyarakat.
- UPTD Gedong Kirtya di Singaraja menyimpan lebih dari 1.750 cakep naskah lontar asli, termasuk kelompok Usadha.

Namun, akses terhadap isinya masih terbatas oleh dua hambatan utama:

1. **Hambatan aksara.** Aksara Bali kuno ditulis tanpa spasi, dengan pasang aksara yang bertumpuk. Foto lontar sering gelap, aus, atau miring.
2. **Hambatan bahasa.** Banyak naskah memakai bahasa Kawi klasik, bukan bahasa Bali sehari-hari. Istilah tanaman, penyakit, dan ritual sering tidak lagi dipahami.

Hambatan ini membuat lontar lebih sering disimpan daripada dibaca. Pengetahuan di dalamnya sulit dipelajari, sulit dicari, dan sulit dihubungkan dengan penelitian modern.

## 2. Arah Produk

BukaLontar tidak dibatasi pada satu genre naskah. Produk ini adalah pembaca lontar Bali secara umum:

- menerima foto lempir,
- menghasilkan transliterasi Latin,
- menerjemahkan teks ke Indonesia dan Inggris,
- menampilkan sumber dan tingkat keterbacaan,
- menyediakan laporan yang dapat disalin atau dicetak.

Di atas alur umum itu, lontar Usadha mendapat perlakuan khusus. Jika naskah memuat resep tanaman obat, BukaLontar mengekstrak tanaman, bagian yang dipakai, bentuk racikan, indikasi, senyawa, dan rujukan. Untuk naskah non-Usadha, lapisan resep tidak dipaksakan.

Pendekatan ini penting agar produk tidak mengklaim semua lontar sebagai naskah pengobatan. Lontar Kakawin, Babad, Geguritan, atau tutur tetap dibaca sebagai teks, bukan sebagai resep.

## 3. Pengguna Utama

- **Pelajar dan mahasiswa** yang ingin membaca lontar tanpa harus menguasai paleografi terlebih dahulu.
- **Peneliti filologi, sejarah, dan budaya** yang membutuhkan transkripsi awal dan terjemahan sebagai bahan pemeriksaan.
- **Kurator museum dan arsip** yang ingin membuat katalog foto lontar lebih mudah dicari.
- **Peneliti etnobotani dan fitofarmaka** yang bekerja dengan lontar Usadha.
- **Praktisi herbal** yang membutuhkan rujukan tekstual, dengan pemahaman bahwa hasil AI bukan nasihat medis.

## 4. Alur Pengguna

1. Pengguna membuka aplikasi dan memilih foto lontar.
2. Browser melakukan preprocessing ringan: downscale, CLAHE, dan unsharp mask.
3. Backend memvalidasi gambar dan menjalankan satu pembacaan.
4. Hasil dibaca sebagai JSON, divalidasi, dan digrounding ke knowledge base tanpa menggabungkan bacaan stokastik.
5. Pengguna melihat transliterasi, terjemahan, sumber, dan tingkat keterbacaan.
6. Jika naskah memuat resep Usadha, pengguna melihat kartu tanaman dan rincian racikan.
7. Pengguna dapat menyalin hasil atau mengunduh laporan PDF.

## 5. Alur Usadha

Alur Usadha adalah fitur khusus, bukan mode terpisah yang harus dipilih pengguna. Jika model menemukan resep pada naskah, hasil resep ditampilkan. Jika tidak, hasil tetap berguna sebagai bacaan dan terjemahan.

Cakupan alur Usadha:

- nama tanaman lokal Bali,
- nama ilmiah dan famili jika cocok dengan Taru Pramana,
- bagian tanaman: daun, kulit kayu, akar, buah, bunga, atau campuran,
- sifat tradisional: panas, dingin, atau sedang,
- bentuk racikan: loloh, boreh, tutuh, atau sasap,
- indikasi dan cara racik,
- senyawa fitokimia dan rujukan yang tersedia di knowledge base,
- disclaimer medis kanonik.

Detail teknis ada di [USADHA_MODE.md](USADHA_MODE.md).

## 6. Prinsip Responsible AI

- **Sumber lebih penting daripada kelancaran.** Teks yang tidak terbaca ditandai, bukan ditebak.
- **Taksonomi tidak boleh dikarang.** Nama ilmiah dan famili hanya diisi jika ada di knowledge base.
- **Sitasi tidak boleh dikarang.** Rujukan yang tidak tersedia dikosongkan.
- **Hasil adalah alat bantu.** Pengguna tetap perlu memeriksa lempir asli dan rujukan filologi.
- **Bukan alat diagnosis.** BukaLontar tidak memberikan diagnosis, dosis, atau instruksi pengobatan.
- **Privasi.** Gambar diproses untuk analisis dan tidak disimpan sebagai file oleh backend.

## 7. Batasan dan Non-Tujuan

- Bukan pengganti filolog atau ahli bahasa Kawi.
- Bukan sistem OCR sempurna untuk semua kondisi foto.
- Bukan database medis dan bukan sumber nasihat pengobatan.
- Belum melakukan klasifikasi genre naskah secara mendalam.
- Belum menyediakan pencarian lintas koleksi atau integrasi katalog eksternal.

## 8. Sumber dan Rujukan

- Dinas Kebudayaan Provinsi Bali, data sensus dan konservasi naskah lontar.
- UPTD Gedong Kirtya, Singaraja, katalog naskah lontar.
- Lontar Taru Pramana, Va. 1384/10.
- Lontar Usadha Rare, Va. 1391/11.
- Foto sampel: Eka343, CC BY-SA 4.0, melalui Wikimedia Commons.

Rincian ada di [DATA_SOURCES.md](DATA_SOURCES.md).

## 9. Bacaan Lanjutan

- [Arsitektur dan alur data](ARCHITECTURE.md)
- [Alur Usadha](USADHA_MODE.md)
- [AI dan Technology Disclosure](AI_DISCLOSURE.md)
- [Indeks dokumentasi](README.md)
