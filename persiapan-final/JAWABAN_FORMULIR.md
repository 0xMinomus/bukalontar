# Jawaban Formulir Final

Salin isi setiap bagian ke kolom yang sesuai. Catatan versi di bawah bagian pertama bukan bagian jawaban formulir.

## 1. Git Tag / Release / Commit Hash Final

e8741650682f9189e68a8ec817b0c5cb220b41fa

Catatan internal: kecocokan website dengan commit ini sudah diverifikasi melalui API Vercel pada 12 November 2026. Deployment produksi `dpl_9vGNjnq6VtSpwVhGMxqn1Boe7C8S` berstatus `READY` dan domain `www.bukalontar.my.id` terpasang pada deployment tersebut; `gitSource.sha` dan `githubCommitSha` sesuai hash di atas. Gunakan hash ini pada formulir. Tag `v1.1-final-submission` tetap snapshot lebih lama dan tidak perlu dipindah.

## 2. Tech Stack Final

Frontend menggunakan React 19, TypeScript, Vite, Tailwind CSS v4, Radix UI, dan Lucide. Backend menggunakan Go dengan standard library net/http. Pengolahan pratinjau gambar menggunakan HTML5 Canvas API. Basis pengetahuan Taru Pramana berupa JSON yang di-embed pada backend, tanpa database eksternal. IndexedDB menyimpan foto dan hasil sesi aktif di browser. Deployment menggunakan Vercel Services untuk frontend dan API. Integrasi AI menggunakan gateway 9router dengan endpoint OpenAI-compatible.

## 3. Model / API / Teknologi AI yang Digunakan

Model utama adalah Gemini 3.8 Flash Low dengan ID ag/gemini-3.8-flash-low melalui gateway 9router, tanpa fallback ke model lain. AI digunakan untuk pemeriksaan aksara, usulan pembacaan foto baru, penerjemahan, dan ekstraksi isi naskah. Empat foto arsip menggunakan transkripsi komunitas Wikisource yang dicocokkan melalui SHA-256; AI menerjemahkan dan mengekstrak teks rujukan tersebut, bukan melakukan OCR baru. Sistem juga memiliki bantuan lokal eksperimental berupa detektor area teks ridge dan pengenal potongan glif HOG + LinearSVC yang dilatih dari DeepLontar. Model lokal tersebut bukan OCR halaman penuh. Tidak ada fine-tuning bobot Gemini.

## 4. Cara Kerja AI pada Produk Final

INPUT: Pengguna mengunggah foto lontar atau memilih contoh arsip. Backend memvalidasi ukuran dan tipe file serta memeriksa kecocokan SHA-256 dengan empat foto arsip.

PROSES AI: Untuk foto baru, Gemini memeriksa aksara terlebih dahulu. Foto yang dinilai sebagai aksara Bali dengan confidence tinggi diteruskan untuk usulan transliterasi Latin, terjemahan Indonesia dan Inggris, serta ekstraksi isi. Untuk foto yang cocok persis dengan arsip, server mempertahankan bacaan Latin dari transkripsi rujukan dan meminta Gemini menerjemahkan serta mengekstrak isinya. Hasil JSON divalidasi. Resep Usadha harus memiliki kutipan utuh dari bacaan sumber sebelum dicocokkan dengan basis pengetahuan Taru Pramana. Informasi ilmiah dan tautan rujukan berasal dari basis pengetahuan, bukan tautan bebas buatan model.

OUTPUT: Pengguna memperoleh bacaan, terjemahan, ringkasan, dan kartu Usadha apabila ada resep yang lolos pemeriksaan, serta dapat menyalin hasil dan mencetak laporan PDF. UI dan PDF membedakan analisis foto eksperimental dari analisis berbantuan transkripsi rujukan. AI menjadi komponen utama karena mengolah citra atau teks sumber menjadi usulan bacaan dan informasi yang lebih mudah diperiksa. Hasil tetap membutuhkan pemeriksaan manusia.

## 5. Perubahan Implementasi AI dari Tahap 1

Kami mengubah alur pembacaan menjadi satu pass pembacaan dengan pemeriksaan aksara terpisah untuk foto baru, serta menetapkan Gemini 3.8 Flash Low sebagai satu-satunya model produk. Kami memperketat validasi JSON dan kutipan resep agar kutipan tidak melewati bagian tidak terbaca, elipsis, pengantar, atau lembar lain. Informasi ilmiah dan tautan Usadha dibatasi pada basis pengetahuan terkurasi. Kami mengumpulkan 200 citra DeepLontar dan melatih detektor teks serta pengenal potongan glif lokal; ini bukan fine-tuning Gemini atau OCR halaman penuh. Eksperimen dan keterbatasannya dicatat dalam laporan evaluasi. Empat contoh arsip kini menggunakan transkripsi komunitas bersumber yang dicocokkan melalui SHA-256, dengan label metode terbuka pada UI dan PDF. Transkripsi Unicode Bali yang belum andal tidak ditampilkan. Peningkatan akurasi dan konsistensi pembacaan foto baru belum terbukti.

## 6. Fitur Utama yang Sudah Berfungsi pada Versi Final

Pengguna dapat mengunggah foto lontar, melihat pratinjau dengan zoom dan penyesuaian kontras, serta meminta analisis AI. Sistem menyediakan usulan transliterasi Latin, terjemahan Indonesia dan Inggris, serta ringkasan ketika teks dapat dibaca. Empat contoh arsip tersedia: Usadha Sukawati, Usadha Rare Kemenuh, Carcan Kucing, dan Kakawin Bharatayuda, dengan bacaan Latin berbantuan transkripsi rujukan. Pada naskah Usadha, sistem dapat menampilkan beberapa kartu resep yang lolos pemeriksaan kutipan, pencocokan tanaman dengan Taru Pramana, serta tautan PubMed yang tersedia dalam basis pengetahuan. Hasil dapat disalin dan dicetak menjadi PDF, termasuk label sumber dan peringatan medis. Fitur analisis berfungsi sebagai alat bantu eksperimental, bukan jaminan kebenaran bacaan.

## 7. Improvement dari Tahap 1 ke Tahap 2

Perbaikan utama adalah ketertelusuran hasil dan pengaman resep. Setiap resep harus didukung kutipan utuh; pada enam kasus batas terkontrol, resep tidak valid yang sebelumnya lolos turun dari 6 menjadi 0, sementara tiga kutipan utuh tetap diterima. Angka ini mengukur pengaman kutipan, bukan akurasi OCR atau kebenaran medis. Kami menambahkan pemeriksaan aksara sebelum pembacaan foto baru, rujukan PubMed dari basis pengetahuan pada kartu dan PDF, serta dokumentasi sumber data dan evaluasi. Empat contoh landing diperbarui menjadi dua Usadha dan dua non-Usadha dengan transkripsi rujukan berlabel terbuka, sehingga demonstrasi tidak menyamarkan bacaan arsip sebagai OCR baru. Kami juga melakukan eksperimen bantuan glif dan pelatihan model lokal, tetapi tidak mengklaim bahwa eksperimen tersebut telah membuktikan peningkatan akurasi pembacaan halaman.

## 8. Known Limitations / Keterbatasan Produk

Pembacaan foto baru masih eksperimental dan dapat menghasilkan transliterasi, terjemahan, atau resep yang keliru. Output belum deterministik; gambar sama dapat memberi hasil berbeda. Pemeriksaan aksara belum menjamin semua foto Jawa/Kawi ditolak atau semua foto Bali diterima. Kutipan yang cocok tidak membuktikan bacaan atau tafsir resep benar. Empat contoh arsip memakai transkripsi rujukan, sehingga hasilnya tidak mewakili kemampuan OCR pada unggahan baru. Model lokal hanya mendeteksi area teks dan mengenali potongan glif, bukan membaca halaman penuh; bobot Gemini tidak di-fine-tune. Kualitas hasil bergantung pada kualitas foto, basis pengetahuan masih terbatas, dan analisis bergantung pada internet, gateway, ketersediaan penyedia AI, serta kuota layanan. Foto baru dikirim ke penyedia melalui gateway; retensi pihak ketiga belum diverifikasi. Produk bukan pengganti ahli filologi atau tenaga medis, dan hasil Usadha bukan anjuran pengobatan.

## Sumber pemeriksaan

- [Status implementasi](../docs/STATUS.md)
- [AI disclosure](../docs/AI_DISCLOSURE.md)
- [Fitur dan alur](../README.md)
- [Evaluasi dan hasil eksperimen](../docs/AI_EVALUATION.md)
