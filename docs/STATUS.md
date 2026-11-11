# Status Proyek

Catatan hidup untuk Obsidian dan GitHub. Diperbarui 12 November 2026.

## Ringkas

BukaLontar membaca foto lontar Bali, menghasilkan transliterasi, terjemahan, dan kartu Usadha. Bantuan lokal: detektor area teks dan pengenal glif terlatih. **Akurasi aksara Bali masih rendah**; jangan mengklaim sebaliknya.

## Yang berjalan

- API Go + frontend React di Vercel Services, produksi `https://www.bukalontar.my.id`.
- Pemeriksaan aksara terpisah pada foto target, lalu satu pembacaan hanya untuk keputusan `balinese/high`. Semua panggilan memakai `ag/gemini-3.8-flash-low`, tanpa fallback model.
- Tautan PubMed terkurasi pada kartu Usadha, tab Rujukan, dan PDF.
- Target ditolak tidak mencapai pembaca. Penolakan pembaca atau konflik aksara mengosongkan seluruh bacaan/resep; Unicode eksperimental tidak ditampilkan di UI/PDF.
- Mode bantuan: `none`, `examples` (default), `chart`, `ocr`, `leaves`, `both`, `all`. Contoh glif dan katalog tidak masuk pemeriksaan aksara.
- Detektor teks dan pengenal glif lokal dilatih dari 159 lembar DeepLontar; inferensi Go murni.
- Lima halaman Wikisource berlabel sebagai holdout evaluasi; skrip CER tersedia.
- Hasil uji penolakan belum memenuhi target: pada batch akhir, dua respons lontar pembanding kategori Jawa lolos dan satu gagal upstream; contoh Kawi cetak ditolak 3/3. Foto Bali diterima 4/6. Lihat [laporan dan batas label](AI_EVALUATION.md#pembatasan-aksara-temuan-awal-12-november).
- Kedua contoh pilihan diganti menjadi sisi depan Usada Lare dan Usada Rare 154, dengan transkripsi aksara Bali tervalidasi di Wikisource. Uji klik contoh terakhir menghasilkan HTTP 200 untuk keduanya; Rare sempat ditolak pada uji upload awal. Model/pengaman tidak diubah dan tidak ada bypass untuk contoh.

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
