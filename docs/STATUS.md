# Status Proyek

Snapshot final submission, 12 November 2026. Tag: [`v1.1-final-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.1-final-submission). Peningkatan kode terakhir: [`4acfb64`](https://github.com/0xMinomus/bukalontar/commit/4acfb64); tag final juga mencakup sinkronisasi dokumentasi sesudah commit tersebut.

## Ringkas

BukaLontar menyediakan pembacaan eksperimental foto lontar Bali. **Belum siap diklaim akurat atau digunakan sebagai rujukan medis.** Pengujian terbaru menemukan bacaan dan resep yang tidak sesuai sumber, meskipun deteksi aksara dan pemeriksaan kutipan lolos.

## Yang berjalan

- API Go + frontend React di Vercel Services, produksi `https://www.bukalontar.my.id`.
- Pemeriksaan aksara terpisah pada foto target, lalu satu pembacaan hanya untuk keputusan `balinese/high`. Semua panggilan memakai `ag/gemini-3.8-flash-low`, tanpa fallback model.
- Tautan PubMed terkurasi pada kartu Usadha, tab Rujukan, dan PDF.
- Target ditolak tidak mencapai pembaca. Jika pemeriksaan awal menerima tetapi pembaca menolak atau keluaran mengandung aksara lain, bacaan/resep dikosongkan dan hasil diberi status tidak terbaca. Unicode eksperimental tidak ditampilkan di UI/PDF.
- Mode bantuan: `none`, `examples` (default), `chart`, `ocr`, `leaves`, `both`, `all`. Contoh glif dan katalog tidak masuk pemeriksaan aksara.
- Detektor teks dan pengenal glif lokal dilatih dari 159 citra berlabel dalam 80 kelompok daun DeepLontar; inferensi Go murni. Ini bukan pembaca kalimat.
- Lima halaman Wikisource berlabel sebagai holdout evaluasi; skrip CER tersedia.
- Hasil uji penolakan belum memenuhi target: pada batch akhir, dua respons lontar pembanding kategori Jawa lolos dan satu gagal upstream; contoh Kawi cetak ditolak 3/3. Foto Bali diterima 4/6. Lihat [laporan dan batas label](AI_EVALUATION.md#pembatasan-aksara-temuan-awal-12-november).
- Kedua contoh pilihan adalah sisi depan Usada Lare dan Usada Rare 154, dengan transkripsi aksara Bali yang dikoreksi di Wikisource (pagequality 3, bukan validasi ahli independen). Foto Bali tetap bisa ditolak.
- Pipeline menghapus katalog tanaman dari prompt, mewajibkan kutipan transliterasi AI untuk resep, dan mengirim foto asli. Kutipan cocok tidak membuktikan bacaan benar atau resep nyata.
- Pembanding baru 12 November: foto Bali 6/6 diterima, pembanding Jawa dan Kawi 4/4 ditolak; hanya lima foto dan dua pengulangan. Baseline juga 10/10 pada foto yang sama, sehingga belum membuktikan peningkatan deteksi. Akurasi isi tetap gagal; lihat [evaluasi](AI_EVALUATION.md).
- Diagnostik tambahan pada segmen buku NusaAksara: Bali diterima 1/4, Jawa ditolak 4/4; satu bacaan Bali tidak sesuai anotasi. Ini bukan uji foto lontar dan datanya tidak disertakan dalam produk.
- Uji foto lontar tambahan: Aji Saraswati kategori Jawa di Commons diterima 3/3 tetapi tidak terbaca; label aksara belum diperiksa ahli. Usada Kacacar diterima 2/2, tetapi kedua bacaan tidak sesuai transkripsi sumber. Lihat [evaluasi foto](AI_EVALUATION.md#uji-tambahan-khusus-foto-lontar).
- Upgrade pengaman kutipan: enam kasus batas yang sebelumnya meloloskan resep kini ditolak; tiga kasus kutipan utuh tetap diterima. Ini pengujian terkontrol, bukan akurasi foto. Uji foto sebelum/sesudah tidak membuktikan peningkatan bacaan; OCR literal dan pemeriksa visual tambahan tidak diintegrasikan. Lihat [perbandingan upgrade](AI_EVALUATION.md#upgrade-pengaman-kutipan-dan-pembanding-baseline).
- Pada pemeriksaan 12 November, health produksi `https://www.bukalontar.my.id/health` melaporkan `ok=true`, model Flash Low, dan 10 entri KB. Health tidak melaporkan commit deployment, sehingga tidak membuktikan produksi identik dengan tag.


Verifikasi kode terakhir sebelum sinkronisasi dokumentasi: 51 kasus backend lulus, build produksi berhasil, dan upload foto ke API lokal menghasilkan HTTP 200. Tidak ada klaim uji kode baru dari perubahan Markdown saja.

## Angka terakhir

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
