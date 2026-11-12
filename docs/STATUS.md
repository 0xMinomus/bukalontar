# Status Proyek

Status `main` setelah perbaikan empat contoh arsip, 12 November 2026. Tag [`v1.1-final-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.1-final-submission) tetap snapshot lama dan tidak dipindah. Pengaman kutipan pada [`4acfb64`](https://github.com/0xMinomus/bukalontar/commit/4acfb64) tetap dipertahankan.

## Ringkas

BukaLontar menyediakan pembacaan eksperimental foto lontar Bali. **Belum siap diklaim akurat atau digunakan sebagai rujukan medis.** Pengujian terbaru menemukan bacaan dan resep yang tidak sesuai sumber, meskipun deteksi aksara dan pemeriksaan kutipan lolos.

## Pemeriksaan kesiapan submission, 12 November 2026

- `go test ./...`: 63 tes lulus dalam 11 paket; `go vet ./...` selesai tanpa temuan.
- `npm run lint` selesai tanpa error, dengan dua peringatan Fast Refresh pada primitive `badge.tsx` dan `button.tsx`. `make build` berhasil membangun frontend dan binary Go dengan frontend ter-embed.
- Health produksi mengembalikan `ok=true`, model `ag/gemini-3.8-flash-low`, dan 10 entri KB. Alur tombol Sukawati, pratinjau, analisis nyata, serta tab terjemahan, kartu Usadha, dan rujukan berhasil dibuka. Hasil menunjukkan label transkripsi rujukan; ini bukan pengujian akurasi OCR.
- Penghambat deployment diperbaiki: refresh `/hasil` sebelumnya menghasilkan `404: NOT_FOUND` sebelum aplikasi memulihkan IndexedDB. Rewrite SPA kini berada di dalam konfigurasi service web. Deployment `dpl_8z55xwcqYNrB2FsoUthGZPF71CDy` untuk commit `e76cce9ba4f59cc30679123d4fc475c257cdd00a` berstatus READY. Pada domain produksi, `/hasil` berhasil memulihkan hasil Sukawati setelah refresh dan `/analisis` membuka aplikasi tanpa 404. Analisis Carcan Kucing juga berhasil tanpa kartu obat; PDF Sukawati memuat label metode, sumber, kutipan, dan peringatan.
- Landing produksi pada viewport 390 px tidak memiliki overflow horizontal. Pemeriksaan ini bukan audit aksesibilitas menyeluruh.
- Folder `persiapan-final` tidak diubah atau dipush dalam pemeriksaan ini. Folder tersebut sudah tercatat pada commit sebelumnya; jangan menganggapnya belum pernah dipublikasikan.

## Yang berjalan

- API Go + frontend React di Vercel Services, produksi `https://www.bukalontar.my.id`.
- Foto baru memakai pemeriksaan aksara terpisah, lalu pembacaan untuk `balinese/high`. Empat foto arsip yang cocok SHA-256 memakai transkripsi komunitas yang dipasok untuk terjemahan/ekstraksi, dengan label terbuka pada UI/PDF. Semua panggilan tetap `ag/gemini-3.8-flash-low`, tanpa fallback.
- Tautan PubMed terkurasi pada kartu Usadha, tab Rujukan, dan PDF.
- Target ditolak tidak mencapai pembaca. Jika pemeriksaan awal menerima tetapi pembaca menolak atau keluaran mengandung aksara lain, bacaan/resep dikosongkan dan hasil diberi status tidak terbaca. Unicode eksperimental tidak ditampilkan di UI/PDF.
- Mode bantuan: `none`, `examples` (default), `chart`, `ocr`, `leaves`, `both`, `all`. Contoh glif dan katalog tidak masuk pemeriksaan aksara.
- Detektor teks dan pengenal glif lokal dilatih dari 159 citra berlabel dalam 80 kelompok daun DeepLontar; inferensi Go murni. Ini bukan pembaca kalimat.
- Lima halaman Wikisource berlabel sebagai holdout evaluasi; skrip CER tersedia.
- Hasil uji penolakan belum memenuhi target: pada batch akhir, dua respons lontar pembanding kategori Jawa lolos dan satu gagal upstream; contoh Kawi cetak ditolak 3/3. Foto Bali diterima 4/6. Lihat [laporan dan batas label](AI_EVALUATION.md#pembatasan-aksara-temuan-awal-12-november).
- Empat pilihan landing: Usadha Sukawati, Usadha Rare Kemenuh, Carcan Kucing, dan Kakawin Bharatayuda. Dua Usadha memuat racikan pada sisi yang dipilih; Lare lama berisi gejala tanpa racikan tanaman di sisi depannya.
- Arsip cocok mempertahankan bacaan Latin rujukan, bukan OCR baru. Resep harus mengutip bacaan sumber; Carcan/Kakawin tidak diberi kartu obat. Foto lain tetap dapat salah dibaca atau ditolak. Kutipan tidak memvalidasi terjemahan, tafsir racikan, atau keamanan medis.
- Pembanding baru 12 November: foto Bali 6/6 diterima, pembanding Jawa dan Kawi 4/4 ditolak; hanya lima foto dan dua pengulangan. Baseline juga 10/10 pada foto yang sama, sehingga belum membuktikan peningkatan deteksi. Akurasi isi tetap gagal; lihat [evaluasi](AI_EVALUATION.md).
- Diagnostik tambahan pada segmen buku NusaAksara: Bali diterima 1/4, Jawa ditolak 4/4; satu bacaan Bali tidak sesuai anotasi. Ini bukan uji foto lontar dan datanya tidak disertakan dalam produk.
- Uji foto lontar tambahan: Aji Saraswati kategori Jawa di Commons diterima 3/3 tetapi tidak terbaca; label aksara belum diperiksa ahli. Usada Kacacar diterima 2/2, tetapi kedua bacaan tidak sesuai transkripsi sumber. Lihat [evaluasi foto](AI_EVALUATION.md#uji-tambahan-khusus-foto-lontar).
- Upgrade pengaman kutipan: enam kasus batas yang sebelumnya meloloskan resep kini ditolak; tiga kasus kutipan utuh tetap diterima. Ini pengujian terkontrol, bukan akurasi foto. Uji foto sebelum/sesudah tidak membuktikan peningkatan bacaan; OCR literal dan pemeriksa visual tambahan tidak diintegrasikan. Lihat [perbandingan upgrade](AI_EVALUATION.md#upgrade-pengaman-kutipan-dan-pembanding-baseline).
- Pada pemeriksaan 12 November, health produksi `https://www.bukalontar.my.id/health` melaporkan `ok=true`, model Flash Low, dan 10 entri KB. Health tidak melaporkan commit deployment, sehingga tidak membuktikan produksi identik dengan tag.


Verifikasi perbaikan contoh: paket Go AI, model, dan handler lulus; empat kasus baru menguji identitas byte, pemalsuan `source_id`/provenance, kutipan karangan, dan sumber nonmedis. Build produksi berhasil. Empat tombol landing menjalankan analisis nyata HTTP 200; label sumber dan PDF diperiksa. Pada 390 px tidak ada overflow, empat thumbnail termuat, dan tidak ada error JavaScript tercatat.

## Hasil contoh arsip

Batch akhir dua kali per foto menghasilkan 8/8 respons, 4/4 permintaan Usadha dengan kartu, dan 4/4 nonmedis tanpa kartu. Bacaan Latin 8/8 sama dengan referensi yang dipasok, bukan bukti akurasi OCR. Dua kontrol Jawa/Kawi dengan `source_id` arsip tetap HTTP 422 tanpa provenance arsip; sampel kecil ini tidak menjamin seluruh aksara lain ditolak. [Metode dan hasil lengkap](AI_EVALUATION.md#empat-contoh-arsip-berbantuan-transkripsi).

## Angka penelitian sebelum perbaikan contoh

| Eksperimen | Hasil |
| --- | --- |
| Latensi Gemini, input sama | median 22,2 detik → 14,0 detik |
| Tile detector, dev | presisi 78,1%, recall 86,9% |
| Pengenal glif, potongan ahli | top-1 82,2% (baseline 57,0%) |
| CER aksara Bali, mode bantuan | 0,957–1,000 |
| `flash-high` | 8/8 timeout 120 detik |
| Perbesaran daun, eksperimen dipulihkan | CER `examples` 0,8834; `leaves` 0,9651; `all` 1,0000 |

## Belum

- Pemetaan resmi 55 ID kelas glif.
- Transkripsi halaman ahli dalam jumlah cukup untuk pelatihan dan evaluasi penuh.
- Fine-tuning model bahasa resmi (endpoint router 404).
- Akurasi bacaan yang layak tanpa pemeriksaan manusia.
- Pembatasan aksara masih berbasis Gemini umum; tidak dapat menjamin setiap gambar Jawa/Kawi ditolak atau setiap foto Bali diterima.

## Berkas kunci

- [Evaluasi AI](AI_EVALUATION.md)
- [Training OCR](OCR_TRAINING.md)
- [Sumber data](DATA_SOURCES.md)
- [Disclosure](AI_DISCLOSURE.md)

## Aturan kerja

- Jangan mengarang label atau sitasi.
- Semua klaim akurasi harus punya hasil tersimpan.
- Perubahan perilaku wajib disertai uji dan dokumentasi.
