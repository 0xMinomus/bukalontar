# Konsep Produk: BukaLontar

Alat bantu eksperimental untuk membaca foto lontar beraksara Bali, dengan alur khusus untuk naskah Usadha. Dokumen ini mengikuti snapshot [`v1.1-final-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.1-final-submission); arah produk tidak sama dengan klaim kemampuan yang telah terbukti.

## 1. Masalah

Naskah lontar tersimpan di koleksi masyarakat, perpustakaan, dan museum Bali. Akses terhadap isinya menghadapi dua hambatan utama:

1. **Hambatan aksara.** Aksara Bali kuno ditulis tanpa spasi, dengan pasang aksara yang bertumpuk. Foto lontar sering gelap, aus, atau miring.
2. **Hambatan bahasa.** Banyak naskah memakai bahasa Kawi klasik, bukan bahasa Bali sehari-hari. Istilah tanaman, penyakit, dan ritual sering tidak lagi dipahami.

Hambatan ini membuat lontar lebih sering disimpan daripada dibaca. Pengetahuan di dalamnya sulit dipelajari, sulit dicari, dan sulit dihubungkan dengan penelitian modern.

## 2. Arah Produk

BukaLontar tidak dibatasi pada satu genre, tetapi membatasi sasaran pada aksara Bali. Bahasa Kawi/Jawa Kuno yang ditulis dalam aksara Bali tetap termasuk lingkup; aksara Jawa dan aksara Kawi tidak termasuk. Pembatasan masih bergantung pada Gemini dan belum andal untuk semua foto.

- menerima foto lempir,
- menghasilkan usulan transliterasi Latin,
- mengusulkan terjemahan Indonesia dan Inggris,
- menampilkan dugaan sumber dan confidence model, bukan ukuran akurasi,
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
2. Browser menyiapkan pratinjau yang dapat diperjelas, tetapi mengirim foto asli ke API.
3. Backend memvalidasi gambar, mengklasifikasikan aksara, lalu menjalankan satu pembacaan hanya untuk keputusan `balinese/high`.
4. Backend memvalidasi JSON dan kutipan resep sebelum grounding KB. Pengantar, lembar lain, serta bagian tidak terbaca tidak boleh menjembatani kutipan.
5. Pengguna melihat usulan transliterasi dan terjemahan beserta peringatan.
6. Jika resep lolos pengaman, pengguna melihat kartu tanaman dan rujukan KB.
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

- **Sumber lebih penting daripada kelancaran.** Prompt meminta penanda bagian tidak terbaca. Ini aturan yang diinginkan, bukan jaminan bahwa model tidak menebak.
- **Taksonomi tidak boleh dikarang.** Nama ilmiah dan famili hanya diisi jika ada di knowledge base.
- **Sitasi tidak boleh dikarang.** Rujukan yang tidak tersedia dikosongkan.
- **Hasil adalah alat bantu.** Pengguna tetap perlu memeriksa lempir asli dan rujukan filologi.
- **Bukan alat diagnosis.** Dugaan indikasi dan racikan ditampilkan sebagai hasil pembacaan sumber, bukan rekomendasi pengobatan atau dosis.
- **Privasi.** Backend memproses gambar di memori dan meneruskannya ke gateway AI. Browser menyimpan sesi aktif, termasuk foto; kebijakan retensi pihak ketiga belum diverifikasi.

## 7. Batasan dan Non-Tujuan

- Bukan pengganti filolog atau ahli bahasa Kawi.
- Bukan sistem OCR sempurna untuk semua kondisi foto.
- Bukan database medis dan bukan sumber nasihat pengobatan.
- Belum melakukan klasifikasi genre naskah secara mendalam.
- Belum menyediakan pencarian lintas koleksi atau integrasi katalog eksternal.
- Perbaikan terakhir memperketat konsistensi kutipan resep. Belum ada bukti peningkatan akurasi bacaan atau terjemahan. Angka dan batas pengujian ada pada [evaluasi baseline](AI_EVALUATION.md#upgrade-pengaman-kutipan-dan-pembanding-baseline).

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
