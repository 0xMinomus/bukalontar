# Evaluasi AI BukaLontar

Tanggal pengujian: 11–12 November 2026 (timestamp berkas JSONL menggunakan UTC). Hasil ini adalah evaluasi pipeline nyata melalui 9router, bukan hasil fine-tuning.

## Keputusan berdasarkan bukti

Sol sudah terintegrasi dan dapat dipilih melalui `AI_MODEL=cx/gpt-6.1-sol` serta `AI_REASONING_EFFORT=low`. Namun, Sol tidak dipilih sebagai default: pada lima citra lontar, seluruh 15 pengulangan pipeline baru hanya menghasilkan penanda tidak terbaca. Respons JSON yang berhasil bukan bukti pembacaan yang berhasil.

Default tetap `ag/gemini-3.8-flash-low`. Penghapusan dua pass dan union lintas hasil mempercepat pipeline, tetapi belum membuktikan kenaikan akurasi maupun konsistensi. Total 72 percobaan benchmark telah diselesaikan; diagnostik dan smoke test browser dicatat terpisah.

### Kandidat tambahan yang ditolak

Tiga gambar pembanding yang sama, tiga pengulangan per model:

| Kandidat | Respons berhasil | Hasil |
| --- | --- | --- |
| `ag/gemini-4-argon-low` | 0/9 | Router mengembalikan HTTP 503 dengan upstream 404 `Requested entity was not found.`. Ketersediaan ID di katalog tidak membuktikan route inferensinya berfungsi. Kemampuan membaca model ini belum dapat dinilai. |
| `ag/gemini-3.8-flash-medium` | 4/9 | Lima percobaan mencapai batas 120 detik. Hanya satu dari enam percobaan lontar berhasil; tiga keberhasilan lain adalah gambar negatif. Median seluruh percobaan 120,003 detik. Tidak memenuhi sasaran cepat dan andal pada konfigurasi yang diuji. |

Bukti ada di `evaluation/candidates-20261111.jsonl`; metrik seluruh kandidat ada di `evaluation/metrics-20261111.json`. Kegagalan akses dan timeout bukan bukti model secara intrinsik tidak mampu, tetapi cukup untuk tidak memilih jalur tersebut bagi final saat ini.

## Rancangan dan cakupan

- Baseline: pipeline lama, dua pembacaan paralel, pemilihan hasil berdasarkan jumlah resep/panjang teks, lalu union tanaman. Model Gemini Flash Low dan Sol masing-masing diuji tiga kali pada Wanagiri, Sukawati, dan satu logo tanpa aksara Bali.
- Pipeline baru: satu pembacaan, prompt berbasis bukti, satu perbaikan format bila diperlukan, grounding server, tanpa union hasil stokastik. Gemini Flash Low dan Sol Low masing-masing diuji tiga kali pada lima citra lontar dan satu gambar negatif.
- Lima citra: dua sampel Wikimedia di aplikasi serta `1a.jpg`, `50b.jpg`, dan `100a.jpg` dari DeepLontar. Ini pemeriksaan awal, bukan benchmark representatif seluruh koleksi lontar Bali.
- Runner memakai `Service.Analyze` yang sama dengan aplikasi. Gambar diberikan langsung sebagai byte sumber tanpa CLAHE browser; SHA-256 dicatat untuk perbandingan setara. Pengujian browser dengan preprocessing dilakukan terpisah dan tidak dicampurkan ke statistik ini.
- Permintaan dilakukan berulang tanpa cache jawaban. Tidak ada training set yang disisipkan ke prompt dan tidak ada jawaban model yang dijadikan label ahli.
- Latensi mencakup inferensi, retry/perbaikan yang terjadi, validasi, serta grounding. Kondisi jaringan dan antrean provider tidak dikontrol; pengujian tambahan dapat berbagi gateway. Hasil bukan SLA.

## Perbandingan pada input yang sama

Tiga gambar yang sama, tiga pengulangan per gambar, sembilan percobaan per baris. Angka latensi mencakup gambar negatif.

| Pipeline/model | Respons berhasil | Median | Maksimum | Kesamaan bacaan antarpengulangan | Jaccard himpunan tanaman |
| --- | --- | --- | --- | --- | --- |
| Lama, Gemini Flash Low | 9/9 | 22,207 detik | 29,105 detik | 0,3395 | 0,4167 |
| Lama, Sol | 8/9 | 24,023 detik | 29,759 detik | tidak dapat dinilai | tidak dapat dinilai |
| Baru, Sol Low | 9/9 | 13,938 detik | 16,843 detik | tidak dapat dinilai | tidak dapat dinilai |
| Baru, Gemini Flash Low | 9/9 | 14,004 detik | 18,997 detik | 0,2943 | 0,3333 |

Median Gemini pada input setara turun sekitar 36,9%. Kesamaan bacaan dan himpunan tanaman justru belum membaik pada sampel ini. Jangan menyebut perubahan tersebut peningkatan akurasi atau konsistensi. Kedua respons dapat sama-sama keliru, dan nama tanaman dapat berbeda ejaan. Penggabungan dua hasil dahulu juga bukan bukti tanaman tambahan itu benar.

Pada enam gambar pipeline baru, Gemini berhasil mengembalikan 18/18 respons dengan median 14,029 detik; Sol 18/18 dengan median 12,062 detik. Sol tidak menghasilkan bacaan leksikal pada seluruh 15 percobaan gambar positif. Penolakan membaca tidak dihitung sebagai kesamaan bacaan sempurna. Satu percobaan baseline Sol gagal dengan HTTP 502.

### Diagnostik Sol

- Probe gambar lontar diterima router dalam sekitar 7 detik, tetapi kata yang terbaca dikembalikan sebagai `unknown`.
- Gambar kontrol dengan teks `BALI 7392`, persegi merah, dan lingkaran biru dikenali bentuk/warna/jumlahnya; teks dibaca `BALI7392`. Ini membuktikan jalur gambar bekerja, bukan kemampuan membaca aksara Bali.
- Probe reasoning `high` pada baris pertama DeepLontar `1a.jpg` memerlukan 47,142 detik dan tetap menghasilkan `[tidak terbaca]` pada transkripsi dan transliterasi.
- Alur browser pada Sukawati, termasuk preprocessing, mengembalikan HTTP 200 dalam 15,825 detik tetapi tidak menghasilkan bacaan. Sol belum layak menjadi pengganti berdasarkan bukti ini.

## Definisi metrik

`scripts/summarize_ai.py` menghitung median/maksimum latensi, keberhasilan respons, deteksi aksara, abstensi, kesamaan teks antarpengulangan, dan Jaccard himpunan tanaman. Penanda lembar dan teks tidak terbaca dikeluarkan dari perhitungan kesamaan. Pasangan tanpa bacaan atau tanpa tanaman pada kedua hasil tidak mendapatkan nilai kesamaan 1.

Kesamaan bacaan menggunakan `SequenceMatcher` tanpa autojunk. Jaccard hanya menormalisasi huruf kecil dan awalan `taru`; perbedaan ejaan/diakritik belum dipetakan secara filologis. Metrik ini bukan akurasi, precision resep, atau validasi klinis.

`references-20261111.json` hanya memberi label keberadaan aksara pada gambar positif/negatif. Evaluasi 11 November tidak menghitung CER karena belum ada transkripsi halaman; nilainya `null`, bukan nol.

Evaluasi 12 November memakai lima halaman Wikisource berlabel dengan field `expected_balinese_text` (lihat `data/reading-examples/holdout/references.json`). Script menghitung CER hanya pada blok aksara Bali; teks Latin dan penanda `[tidak terbaca]` tidak dihitung. Hasil pada dua halaman ada di bagian bantuan pembacaan. Akurasi terjemahan serta precision/recall resep tetap belum diukur karena tidak ada label terjemahan atau resep ahli.

## Korpus dan kelayakan training

[Korpus DeepLontar v2](https://doi.org/10.6084/m9.figshare.20103803.v2) menyediakan 200 citra asli lontar beraksara Bali dengan lisensi CC BY 4.0. Semua telah diunduh dan diperiksa checksum, duplikasi byte, serta dekode JPEG. Ini 200 citra permukaan, bukan klaim 200 naskah berbeda. Keragaman koleksi dan identitas tiap naskah tidak tersedia secara rinci.

Anotasi berisi 45.261 kotak karakter. `76a.txt` memiliki 76 ID kelas di luar rentang 0–54 yang dipublikasikan; seluruh 198 kotaknya dikarantina tanpa mengubah berkas asli. Tersisa 199 berkas label dan 45.063 kotak yang memenuhi pemeriksaan format. Label tersebut bukan transkripsi halaman, transliterasi, terjemahan, atau resep Usadha.

Ada pelatihan model lokal, tetapi bukan fine-tuning Sol/Gemini. `scripts/train_lontar_ocr.py` melatih detektor okupansi teks dan pengenal glif karakter dari 159 lembar train; hasil tahan-lihat ada di [OCR_TRAINING.md](OCR_TRAINING.md). ID kelas tetap opasitas karena pemetaan resmi tidak tersedia di arsip dataset.

`GET /v1/fine_tuning/jobs` pada router yang tersedia mengembalikan HTTP 404, dan tidak ada akses fine-tuning upstream yang terverifikasi. Gambar tak bertranskripsi tidak dapat dianggap contoh jawaban yang benar. Detektor karakter terpisah tidak otomatis melatih Sol atau menyelesaikan terjemahan Kawi.

Prasyarat fine-tuning model bahasa yang belum tersedia: endpoint dan model yang mendukung fine-tuning, pemetaan kelas resmi yang terverifikasi, pasangan citra/transkripsi ahli sesuai tugas, serta pembagian data berdasarkan sumber naskah untuk menghindari kebocoran.

## Bantuan pembacaan: contoh glif dan detektor teks

Dua halaman *Aji Kreket* berlabel Wikisource, tiga pengulangan per konfigurasi, Gemini Flash Low, prompt dengan pengaman skrip. Metrik: character error rate hanya pada blok aksara Bali U+1B00–U+1B7F; penanda tak terbaca dan teks Latin tidak dihitung. Angka rendah lebih baik.

| Mode | Respons berhasil | CER aksara Bali | Median | Maksimum | Abstain |
| --- | --- | --- | --- | --- | --- |
| `none` | 6/6 | 0,9972 | 15,144 detik | 28,070 detik | 0/6 |
| `examples` | 6/6 | 0,9570 | 14,348 detik | 25,009 detik | 2/6 |
| `ocr` | 5/6 | 1,0000 | 16,178 detik | 34,446 detik | 3/6 |
| `both` | 6/6 | 0,9791 | 11,546 detik | 23,805 detik | 3/6 |

**Temuan:** pada sampel ini, bantuan tidak menghasilkan bacaan aksara Bali yang dapat dipakai. CER semua mode di atas 0,95; ketiga mode abstain lebih sering. Contoh glif memberi selisih CER terkecil (0,957 berbanding 0,997) tetapi itu belum bukti akurasi pada naskah lain. Detektor fokus (`ocr`) justru terburuk dan pernah membuat model menulis aksara Jawa; pengaman server kini membuang karakter non-Bali dan menampilkan peringatan.

Karena itu bantuan bersifat eksperimental. Default produk adalah `examples`: hanya contoh glif ahli, tanpa crop heuristik, dan satu-satunya konfigurasi dengan CER lebih rendah dari `none` pada sampel ini. `ocr` dan `both` tetap tersedia lewat `AI_READING_ASSISTANCE` untuk eksperimen, bukan default. Data mentah: `evaluation/assist-ablation-v2-20261112.jsonl`; metrik: `evaluation/assist-ablation-v2-metrics.json`.

Batas metodologi: dua halaman dari satu naskah, satu model, tiga pengulangan. Transkripsi Wikisource mencakup dua sisi lembar (2B dan 3A), sedangkan model dapat membaca hanya sebagian, sehingga CER absolut tidak dapat ditafsirkan sebagai akurasi halaman. Dua mode berbeda tidak diuji dengan penghitung yang sama bila jumlah respons berhasil berbeda (satu kegagalan `ocr`).

## Perubahan produk yang diterapkan

- Satu pass menggantikan dua pass dan union yang dapat memperbanyak temuan tanpa dasar.
- Grup naskah/bahasa boleh tidak diketahui; prompt tidak lagi menyuruh memilih kelompok Usadha terdekat.
- Resep non-Usadha/tidak teridentifikasi tidak dibuat dari penanda kosong. Bentuk obat tidak otomatis menjadi loloh ketika tidak disebutkan.
- Field ilmiah, klaim farmakologi, dan tautan rujukan ditentukan KB server. Model tidak boleh mempertahankan tautan atau klaim ilmiah untuk tanaman yang tidak cocok.
- Kartu resep, tab Rujukan, dan PDF memuat tautan PubMed terkurasi. Pencarian literatur dibedakan dari artikel tinjauan. Tautan tidak dinyatakan sebagai bukti khasiat resep.
- Label UI membedakan aksara terdeteksi dari teks terbaca. Tidak ditemukannya resep tidak dianggap bukti naskah pasti non-Usadha.
- Detektor teks dan pengenal glif lokal dilatih dari DeepLontar. Default bantuan `examples`; crop fokus opt-in dan eksperimental.
- Model diminta menuliskan transkripsi aksara Bali literal. Karakter dari aksara mirip (Jawa, Kawi, Sunda, Devanagari) dibuang server dengan peringatan; transkripsi kosong juga diberi peringatan.
- Panel bantuan pembacaan dan bagian aksara Bali tersedia di UI dan laporan cetak, dengan batasan eksplisit.

## Reproduksi

Dari root repository, dengan konfigurasi `apps/api/.env` yang valid:

```bash
python3 scripts/deeplontar.py acquire
python3 scripts/deeplontar.py validate
go -C apps/api build -o /tmp/bukalontar-evaluate ./cmd/evaluate
/tmp/bukalontar-evaluate \
  -models ag/gemini-3.8-flash-low,cx/gpt-6.1-sol \
  -reasoning low -repeat 3 -variant rerun \
  -out /tmp/bukalontar-rerun.jsonl \
  apps/web/public/samples/lontar-wanagiri-usada-3.jpg \
  apps/web/public/samples/lontar-sukawati-usada-07.jpg \
  logos/bukalontar-icon.png
python3 scripts/summarize_ai.py \
  --references docs/evaluation/references-20261111.json \
  docs/evaluation/baseline-20261111.jsonl \
  docs/evaluation/sol-20261111.jsonl \
  docs/evaluation/gemini-20261111.jsonl
python3 -m unittest discover -s scripts -p test_summarize_ai.py
```

Runner menolak menimpa berkas keluaran yang sudah ada. Baseline lama dikompilasi sebelum perubahan pipeline; checkout saat ini hanya mereproduksi pipeline baru. Hasil mentah baseline disimpan untuk pemeriksaan, bukan untuk disajikan kepada pengguna sebagai hasil analisis baru.

## Bukti verifikasi

- Backend: 21 pengujian lulus pada paket AI, model, KB, dan handler, termasuk penggantian URL model yang berbahaya dengan rujukan kurasi serta pembersihan klaim tanaman tidak dikenal.
- Metrik: dua pengujian batas lulus; abstensi tidak dihitung sebagai konsistensi dan CER memerlukan teks pembanding.
- Korpus: 200 JPEG unik tervalidasi; satu label dikarantina dan tidak diloloskan sebagai data training bersih.
- Browser: analisis nyata Wanagiri dengan Gemini menghasilkan kartu Sotong, tautan pencarian PubMed dan tinjauan PMID 18353572. Tautan memiliki `target="_blank"` dan `rel="noopener noreferrer"`.
- Tab Rujukan menampilkan batasan sumber. PDF hasil render print berisi URL dan disclaimer yang sama. Pada viewport 390 px tidak ada overflow horizontal; tautan terlihat setelah menggulir panel kartu. Tidak ada error browser yang tercatat pada pemeriksaan tersebut.
- Lint frontend: nol error, dua peringatan Fast Refresh pada ekspor komponen UI yang tidak diubah. Build TypeScript/Vite berhasil; perubahan UI terakhir diperiksa lagi pada build produksi.

- Kontrak JSON Schema: seluruh 36 respons pipeline baru Sol/Gemini lolos validasi draft 2020-12. Ini bukti struktur, bukan kebenaran isi.
- `make build` berhasil membangun TypeScript/Vite, menyematkan aset terbaru, dan menghasilkan binary Go produksi.

Perubahan dikerjakan dan diuji lokal. Tidak ada klaim deployment produksi diperbarui atau training model berhasil.
