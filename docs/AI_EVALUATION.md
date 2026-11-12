# Evaluasi AI BukaLontar

Tanggal pengujian: 11–12 November 2026 (timestamp berkas JSONL menggunakan UTC). Hasil ini adalah evaluasi pipeline nyata melalui 9router, bukan hasil fine-tuning.

Snapshot final submission: [`v1.1-final-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.1-final-submission). Model produk hanya `ag/gemini-3.8-flash-low`; nama model lain di bawah adalah riwayat eksperimen, bukan opsi server produk saat ini. Untuk laporan akhir, gunakan [perbandingan baseline terakhir](#upgrade-pengaman-kutipan-dan-pembanding-baseline). Perbaikan terakhir ada pada pengaman kutipan resep, bukan akurasi OCR.

Bagian-bagian berikut mencatat tahap pengujian berbeda. Angka tes, latensi, dan status lokal pada bagian lama berlaku pada tahap tersebut; jangan dijumlahkan atau diperlakukan sebagai hasil rilis yang sama. Sinkronisasi Markdown tidak menjalankan ulang benchmark model. Health produksi pada 12 November melaporkan `ok=true`, 10 entri KB, dan Flash Low, tetapi tidak melaporkan commit deployment.

## Keputusan model pada evaluasi awal

Pada evaluasi awal, Sol diuji melalui konfigurasi runner dengan `cx/gpt-6.1-sol` dan reasoning `low`. Pada lima citra lontar, seluruh 15 pengulangan pipeline baru hanya menghasilkan penanda tidak terbaca. Respons JSON yang berhasil bukan bukti pembacaan yang berhasil. Server produk saat ini menolak model selain Gemini 3.8 Flash Low.

Keputusan akhir mempertahankan `ag/gemini-3.8-flash-low`. Penghapusan dua pass dan union lintas hasil mempercepat pipeline pada batch awal, tetapi belum membuktikan kenaikan akurasi maupun konsistensi. Angka 72 percobaan merujuk benchmark awal, bukan total seluruh eksperimen lanjutan; diagnostik dan smoke test dicatat terpisah.

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

Ada pelatihan model lokal, tetapi bukan fine-tuning Sol/Gemini. `scripts/train_lontar_ocr.py` melatih detektor okupansi teks dan pengenal glif karakter dari 159 citra train dalam 80 kelompok daun fisik; hasil tahan-lihat ada di [OCR_TRAINING.md](OCR_TRAINING.md). ID kelas tetap opasitas karena pemetaan resmi tidak tersedia di arsip dataset.

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

### Model reasoning lebih tinggi

`ag/gemini-3.8-flash-high` diuji pada dua halaman yang sama, mode `none` dan `examples`, dua pengulangan. Seluruh delapan permintaan mencapai batas 120 detik (`context deadline exceeded`) sehingga tidak ada respons yang dapat dinilai CER-nya. Model ini tidak dipakai untuk alur interaktif pada konfigurasi timeout sekarang. Penyebab (antrean provider dibanding beban model) tidak dipisahkan. Data: `evaluation/assist-high-timeout-20261112.jsonl`.

### Eksperimen daun yang sempat terputus

Hasil 12 percobaan Flash Low dari sesi sebelumnya dipulihkan dari berkas lokal, bukan dijalankan ulang. Dua halaman Aji Kreket, dua pengulangan per mode, referensi sama dengan evaluasi sebelumnya:

| Mode | Respons | CER Bali | Median | Kesepakatan pasangan transliterasi |
| --- | --- | --- | --- | --- |
| `examples` | 4/4 | 0,8834 | 17,892 detik | 0,2347 |
| `leaves` | 4/4 | 0,9651 | 24,723 detik | 0,2231 |
| `all` | 4/4 | 1,0000 | 26,492 detik | 0,2022 |

Perbesaran daun dan gabungan semua bantuan tidak memperbaiki hasil pada sampel ini, sekaligus lebih lambat. Default tetap `examples`; tidak ada klaim bacaan sudah akurat atau konsisten. CER `examples` berbeda dari batch sebelumnya, bukan kenaikan yang terbukti karena tidak ada eksperimen terkontrol baru. Batas keselarasan sisi daun dan transkripsi tetap berlaku. Data: `evaluation/assist-leaves-20261112.jsonl`; metrik: `evaluation/assist-leaves-metrics-20261112.json`.

### Pembatasan aksara: temuan awal 12 November

Hanya `ag/gemini-3.8-flash-low` dipanggil dalam pekerjaan lanjutan ini. Prompt satu-pass yang meminta klasifikasi target dan normalisasi penolakan diuji tiga kali pada masing-masing foto Wanagiri, Sukawati, lontar Babad Surapati Acad. 240 halaman PDF 7 (dikategorikan Commons sebagai aksara Jawa), dan contoh Kawi modern. Foto Bali diterima 3/6; contoh Kawi ditolak 3/3. Lontar pembanding menghasilkan satu penolakan, satu respons kosong upstream, dan satu bacaan/resep yang lolos. Kegagalan tidak dihitung sebagai penolakan benar.

Data sebelum pemeriksaan terpisah: `evaluation/bali-scope-20261112.jsonl` dan `evaluation/bali-scope-jawa-20261112.jsonl`. Ini bukti prompt satu-pass belum cukup, bukan bukti pembatasan berhasil. Pemeriksaan foto target tanpa contoh glif/katalog kemudian dipisahkan dari pembacaan, tetap memakai Flash Low. Kontrol browser memperlihatkan foto Bali Sukawati yang diperjelas ditolak pada jalur pertama; pada jalur terpisah Wanagiri diterima, tetapi lontar pembanding yang diperjelas juga pernah lolos dan menghasilkan resep. Penolakan salah dan penerimaan salah tetap mungkin.

Sumber dan byte input pembanding tersedia di [kredit input](evaluation/script-inputs/CREDITS.md). Contoh Kawi merupakan gambar huruf cetak berlabel Latin, bukan foto lontar Kawi kuno. Klasifikasi Babad Surapati didasarkan pada kategori Commons; belum ada validasi aksara oleh paleografer independen. Jangan mengubah hasil sulit menjadi klaim sukses atau menganggap metadata bahasa pasti mengidentifikasi aksara.

Hasil akhir pemeriksaan terpisah (`evaluation/bali-gated-20261112.jsonl`), tiga pengulangan per gambar, byte asli tanpa preprocessing browser:

| Input | Diterima | Ditolak | Gagal upstream | Median seluruh percobaan |
| --- | --- | --- | --- | --- |
| Wanagiri, aksara Bali | 2 | 1 | 0 | 33,998 detik |
| Sukawati, aksara Bali | 2 | 1 | 0 | 16,552 detik |
| Babad Surapati, pembanding kategori aksara Jawa | 2 | 0 | 1 | 29,129 detik |
| Contoh Kawi cetak | 0 | 3 | 0 | 8,302 detik |

**Target penolakan Jawa/Kawi secara andal belum tercapai.** Dua respons pembanding lolos; salah satunya berisi resep. Tidak ada dasar untuk mengklaim model lebih pintar, terjemahan akurat, konsistensi meningkat, atau kecepatan tetap sama. Pemeriksaan terpisah mengisolasi keputusan dari contoh/katalog, tetapi tidak membuat Gemini menjadi classifier independen. Sampel kecil, kategori sumber belum divalidasi ahli, dan tidak ada ground truth bacaan/resep untuk foto ini. Jangan menghitung semua respons HTTP 200 sebagai bacaan yang benar.

Sebanyak 11 respons terstruktur akhir lolos JSON Schema draft 2020-12. Satu respons kosong upstream tidak dinilai sebagai penolakan. Pengujian backend setelah penambahan gate: 33 kasus lulus pada paket AI, model, dan handler. Browser nyata: Wanagiri menghasilkan HTTP 200; foto pembanding yang diperjelas juga menghasilkan HTTP 200 (kegagalan pembatasan), sementara Sukawati pada jalur sebelumnya menghasilkan HTTP 422 (penolakan salah). UI/PDF membedakan bahasa dari aksara, tidak menampilkan Unicode eksperimen/panel bantuan; viewport 390 px tidak overflow. Build produksi `make build` berhasil.

Prasyarat yang belum tersedia untuk klaim lebih kuat: label aksara oleh ahli pada foto sulit lintas Bali/Jawa/Kawi, pasangan citra-bacaan ahli, dan discriminator yang tervalidasi pada holdout independen. Menambah prompt atau mengulang model yang sama bukan pengganti bukti tersebut.



## Perubahan produk yang diterapkan

- Satu pembacaan menggantikan dua pembacaan dan union yang dapat memperbanyak temuan tanpa dasar. Pemeriksaan aksara terpisah kini mendahului pembacaan, bukan pass kedua untuk menggabungkan resep.
- Grup naskah/bahasa boleh tidak diketahui; prompt tidak lagi menyuruh memilih kelompok Usadha terdekat.
- Resep non-Usadha/tidak teridentifikasi tidak dibuat dari penanda kosong. Bentuk obat tidak otomatis menjadi loloh ketika tidak disebutkan.
- Field ilmiah, klaim farmakologi, dan tautan rujukan ditentukan KB server. Model tidak boleh mempertahankan tautan atau klaim ilmiah untuk tanaman yang tidak cocok.
- Kartu resep, tab Rujukan, dan PDF memuat tautan PubMed terkurasi. Pencarian literatur dibedakan dari artikel tinjauan. Tautan tidak dinyatakan sebagai bukti khasiat resep.
- Label UI membedakan aksara terdeteksi dari teks terbaca. Tidak ditemukannya resep tidak dianggap bukti naskah pasti non-Usadha.
- Detektor teks dan pengenal glif lokal dilatih dari DeepLontar. Default bantuan `examples`; crop fokus opt-in dan eksperimental.
- Model diminta menuliskan transkripsi aksara Bali literal untuk evaluasi. Target yang ditolak atau keluaran transkripsi/transliterasi dengan konflik aksara tidak mempertahankan bacaan/resep. Transkripsi Unicode eksperimental tidak ditampilkan di UI/PDF.
- Panel bantuan tidak ditampilkan di UI atau PDF; metadata bantuan tetap tersedia di API untuk evaluasi.
- Mode bantuan diperluas: `chart` (lembar contekan 55 glif), `leaves` (perbesaran daun), `all`. Tidak ada yang terbukti menurunkan CER.

## Reproduksi

Dari root repository, dengan konfigurasi `apps/api/.env` yang valid:

```bash
python3 scripts/deeplontar.py acquire
python3 scripts/deeplontar.py validate
go -C apps/api build -o /tmp/bukalontar-evaluate ./cmd/evaluate
/tmp/bukalontar-evaluate \
  -models ag/gemini-3.8-flash-low \
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

Catatan historis pada tahap ini: perubahan diuji lokal, tanpa klaim deployment produksi telah diperbarui atau fine-tuning model bahasa berhasil. Hasil pelatihan glif lokal dijelaskan terpisah.

## Penggantian contoh pilihan setelah laporan penolakan

Wanagiri memiliki transkripsi komunitas di [Wikisource Bali](https://ban.wikisource.org/w/index.php?oldid=127924). Metadata [Commons Sukawati](https://commons.wikimedia.org/wiki/File:Bali-lontar-Sukawati-Usada-07.jpg) mencantumkan writing system Balinese. Penolakan aplikasi bukan bukti keduanya menggunakan aksara Kawi. Sebelumnya false rejection telah terjadi pada foto Bali; penyebab keputusan spesifik model tidak dapat ditentukan hanya dari HTTP 422.

Kedua contoh pilihan kini memakai sisi depan **Usada Lare** dan **Usada Rare 154**, foto Eka343 dari Balai Bahasa Bali. Transkripsi komunitas Wikisource menggunakan blok aksara Bali: [Lare, revisi 105197](https://ban.wikisource.org/w/index.php?oldid=105197), [Rare 154, revisi 105194](https://ban.wikisource.org/w/index.php?oldid=105194). Keduanya berstatus pagequality 3 (proofread), bukan sertifikasi ahli independen. Foto dipotong untuk mempertahankan satu sisi daun utuh dan mengurangi latar meja, kemudian diperkecil ke lebar 2400 px. Tidak ada rekonstruksi glif. Kredit, koordinat, dimensi, dan checksum ada di `apps/web/public/samples/CREDITS.md`.

Pengujian nyata memakai Gemini 3.8 Flash Low yang sama, preprocessing browser default, tanpa mock atau pengecualian berdasarkan `source_id`. Tidak ada perubahan backend, prompt, atau pengaman penolakan dalam penggantian contoh ini.

| Jalur | Contoh | HTTP | Waktu |
| --- | --- | --- | --- |
| Upload awal | Usada Lare, depan | 200 | 28,807 detik |
| Upload awal | Usada Rare 154, depan | 422 | 32,252 detik |
| Klik contoh landing | Usada Lare, depan | 200 | 32,491 detik |
| Klik contoh landing | Usada Rare 154, depan | 200 | 25,936 detik |

Seluruh hasil, termasuk penolakan awal, disimpan di `evaluation/replacement-samples-web-20261112.jsonl`. Kedua tombol contoh lolos pemeriksaan terakhir, tetapi satu penolakan pada sumber yang sama menunjukkan masalah deteksi belum hilang. Ini verifikasi penggantian contoh dan alur antarmuka, bukan bukti akurasi bacaan, stabilitas, atau peningkatan model. Thumbnail keduanya termuat; pada viewport 390 px lebar dokumen tetap 390 px dan tidak ada error JavaScript yang tercatat.


## Evaluasi keandalan lanjutan dan batas kelulusan

Semua panggilan menggunakan `ag/gemini-3.8-flash-low`. Dibandingkan pada gambar byte asli yang sama: dua sisi depan Usada (Lare dan Rare 154), satu halaman Aji Kreket, satu foto naskah Panji kategori aksara Jawa di Commons, dan satu contoh cetak Kawi. Masing-masing dua kali. Kategori Jawa dan transkripsi Wikisource belum disertifikasi ahli independen. Lare/Rare adalah Wikisource **pagequality 3 (proofread)**, bukan validasi ahli; status yang lebih kuat pada catatan lama harus dibaca dengan koreksi ini. Tidak tersedia ground truth terjemahan atau resep yang diperiksa ahli.

| Versi | Bali diterima | Jawa/Kawi ditolak | Respons berhasil | Median seluruh permintaan |
| --- | --- | --- | --- | --- |
| Baseline sebelum perubahan | 6/6 | 4/4 | 10/10 | 18,462 detik |
| Prompt tanpa katalog saja | 4/6 | 3/4 | 8/10 | 22,192 detik |
| Kandidat prompt netral, kutipan, pemisahan deteksi/bacaan | 6/6 | 4/4 | 10/10 | 16,498 detik |
| Eksperimen contoh foto dan bacaan berpasangan | 6/6 | 4/4 | 10/10 | 17,026 detik |

Angka ini bukan estimasi akurasi populasi. Baseline dan kandidat sama-sama 10/10 pada sampel kecil ini, sementara evaluasi sebelumnya menemukan Jawa lolos dan Bali tertolak. Ketersediaan dan latensi router bervariasi. Tidak ada bukti perbaikan andal dalam deteksi, kecepatan, atau konsistensi.

Pada Rare 154, kandidat menghasilkan resep `bayem` dengan kutipan yang cocok dengan transliterasi AI, tetapi teks pembanding Wikisource menyebut bahan lain. Kecocokan string bukan verifikasi terhadap foto dan bahkan tidak membuktikan kutipan memuat cara racik. Ini contoh kegagalan yang tersisa, bukan resep tervalidasi. Model juga menghasilkan dua bacaan Lare yang berbeda dari sumber dan dari satu sama lain. Terjemahan tidak dapat dinilai tanpa rujukan ahli.

Eksperimen pasangan foto–bacaan berhasil meniru Lare yang sama dengan contoh konteks; ini kontaminasi input, bukan generalisasi. Untuk Rare dan Aji yang tidak dipakai sebagai contoh, CER Unicode tercatat 1,0 (keluaran kosong), dibanding 0,9038 pada baseline. Perbandingan CER ini **tidak setara** karena baseline diminta mengisi transkripsi Unicode sedangkan prompt baru mengizinkan kosong. Tidak ada bukti peningkatan akurasi; pasangan contoh tidak dipakai pada produk.

Kandidat menambahkan: parser JSON classifier yang menerima code fence; status aksara tetap terdeteksi tetapi teks tidak terbaca jika pembaca menolak setelah gate; prompt netral tanpa katalog; kutipan resep dicocokkan ke transliterasi daun pertama dan bahan yang disebut; sifat tradisional tidak lagi ditimpa nilai KB; browser mengirim foto asli dan menampilkan peringatan. Tidak ada perubahan bobot model. Pengujian regresi backend: 43 kasus lulus pada paket AI, model, handler. Pengujian metrik: 4 lulus. Build produksi berhasil; pada browser lokal Usada Lare menghasilkan HTTP 200, peringatan terlihat, viewport 390 px tanpa overflow atau error JavaScript. Bacaan dari uji browser tersebut tetap tidak sesuai sumber.

Data mentah `evaluation/reliability-{baseline,no-catalog,final,paired-context}-20261112.jsonl`, referensi `evaluation/reliability-references-20261112.json`, referensi independen tanpa teks contoh Lare `evaluation/reliability-heldout-references-20261112.json`, dan metrik `evaluation/reliability-metrics-20261112.json`. Input Jawa/Kawi ada di `evaluation/script-inputs/`. Eksperimen pasangan tidak disertakan dalam pipeline.

**Status: empat sasaran belum memenuhi kriteria 100% atau kelayakan akurasi.** Diperlukan label aksara dan transkripsi halaman yang diperiksa ahli pada koleksi Bali/Jawa/Kawi beragam, terjemahan/resep pembanding, pembagian sumber yang independen, dan classifier/pembaca yang tervalidasi. Tanpa prasyarat itu tidak ada dasar mengklaim selesai atau aman untuk keputusan medis.

## Sumber pembanding tambahan: NusaAksara

[NusaAksara](https://aclanthology.org/2025.acl-long.1377/) menyediakan segmen citra beserta transkripsi, transliterasi, dan terjemahan yang dianotasi dan ditinjau manusia. [Dataset publik](https://huggingface.co/datasets/NusaAksara/NusaAksara) mencakup 459 segmen Bali dan 816 segmen Jawa. Bahasa Kawi dapat ditulis dengan aksara Bali; label bahasa `kaw` tidak berarti aksaranya Kawi.

Dataset berasal dari buku yang dipindai, bukan foto lontar. Penulis menyatakan lisensi nonkomersial; data tidak disalin ke produk atau repository. Pembanding ini berguna untuk diagnosis tambahan, tetapi tidak menggantikan uji foto lontar, foto ponsel, naskah Kawi, maupun resep Usadha.

Uji diagnostik 12 November 2026 menggunakan `ag/gemini-3.8-flash-low`, bantuan `examples`, satu pengulangan pada masing-masing empat segmen Bali (`Bali_2_0`, `Bali_3_0`, `Bali_4_0`, `Bali_5_1`) dan empat Jawa (`Jawa_0_0`, `Jawa_1_1`, `Jawa_2_1`, `Jawa_3_2`). Pilihan adalah segmen pertama dengan transliterasi lebih dari 35 karakter pada empat nomor grup pertama yang memenuhi syarat per aksara; nomor grup tidak membuktikan sumber buku berbeda.

| Hasil pada segmen buku | Jumlah |
| --- | --- |
| Bali diterima | 1/4 |
| Jawa ditolak | 4/4 |
| Respons terstruktur | 8/8 |
| Median waktu seluruh permintaan | 10,167 detik |

Satu bacaan Bali yang diterima tidak sesuai transliterasi anotator. Karena input berupa potongan buku, bukan foto lontar, hasil ini menunjukkan batas di luar domain target, bukan estimasi akurasi pada foto lontar. Sampel terlalu kecil untuk klaim reliabilitas atau perbandingan kecepatan. Data gambar dan anotasi tidak didistribusikan bersama repository.

## Uji tambahan khusus foto lontar

Pada 12 November, pipeline tanpa perubahan diuji dengan foto asli, `ag/gemini-3.8-flash-low`, dan bantuan `examples`. Nama file tidak diberikan ke model. Respons mentah tersimpan di `evaluation/photo-followup-20261112.jsonl`; foto tidak disalin ke repository.

| Foto | Keputusan | Bacaan dan batas label |
| --- | --- | --- |
| [Aji Saraswati](https://commons.wikimedia.org/wiki/File:Palmleaf_of_Aji_saraswati.jpg), SHA-256 `6bfe686430b6211f897e8cad608b1279b7d5a8256aaff2bf987e427913e3e228` | Diterima sebagai Bali 3/3 | Semua menyatakan teks tidak terbaca. Commons mengkategorikannya aksara Jawa, tetapi label aksara belum diperiksa ahli independen. Foto hanya 376 × 480 piksel; hasil tidak membuktikan kinerja pada foto tajam. |
| [Usada Kacacar](https://commons.wikimedia.org/wiki/File:Bali-lontar-Balai_Bahasa_Bali-Usada_Kacacar_(Babalining_kacacar).jpg), SHA-256 `a97be48a5012cf074b1bb5951425560729d0c34a3bc7bd9162d89327a9495e88` | Diterima sebagai Bali 2/2 | Kedua bacaan berbeda dan tidak sesuai [transkripsi Wikisource](https://ban.wikisource.org/w/index.php?title=Kaca:Bali-lontar-Balai_Bahasa_Bali-Usada_Kacacar_(Babalining_kacacar).jpg&oldid=105200) yang telah dikoreksi komunitas, bukan validasi ahli independen. |

Foto Kacacar oleh Eka343, CC BY-SA 4.0; Aji Saraswati mengikuti keterangan hak cipta pada halaman sumber. Kedua foto merupakan naskah berbeda, bukan sampel representatif. Hasil ini memperkuat bukti bahwa penerimaan aksara dan bacaan belum andal. Tidak ada perbaikan model atau klaim akurasi baru dari uji ini.

### Percobaan perbesaran setengah daun

Dua foto contoh Usada Lare dan Rare 154 dipotong pada bagian kiri setengah lebar penuh (`0, 0, width//2, height`) dan disimpan sebagai PNG tanpa menambah aksara. Tujuannya menguji apakah lebih sedikit isi per gambar membantu Gemini membaca goresan. Setiap potongan diuji dua kali memakai pipeline dan `ag/gemini-3.8-flash-low` yang sama; hasil mentah ada di `evaluation/half-leaf-20261112.jsonl`.

Dua dari empat permintaan gagal dengan respons upstream kosong. Dua lainnya diterima sebagai Bali, tetapi bacaan Lare tidak sesuai awal transkripsi Wikisource; bacaan Rare juga belum dibuktikan sesuai sumber. Referensi seluruh daun tidak langsung mengukur akurasi potongan, dan percobaan kecil ini tidak menunjukkan peningkatan. Pemotongan tidak diterapkan pada produk.

### Kandidat riset OCR

[Lontar Project di Leiden](https://lontar.liacs.nl/) menerbitkan [kode pengenal glif terisolasi](https://github.com/erickpaulus/isolated-glyph-recognition) dan tautan dataset/model. Hasil pengenalan glif terisolasi tidak sama dengan membaca foto lengkap, transliterasi, atau terjemahan; sistem ini belum diintegrasikan atau diuji pada foto BukaLontar. Lisensi kode MIT tidak otomatis menetapkan lisensi model dan data.

## Upgrade pengaman kutipan dan pembanding baseline

### Tujuan dan keputusan

Baseline adalah kode pada commit `fc174f9`, dikompilasi sebelum perubahan pengaman. Tujuan pembacaan: menguji OCR literal terpisah sebelum terjemahan; tujuan pengaman: mencegah kutipan melompati bagian hilang atau mengambil bukti dari lembar yang salah. Semua panggilan menggunakan `ag/gemini-3.8-flash-low`, tanpa pergantian model.

**Yang diintegrasikan hanya pengaman kutipan. Pembacaan OCR khusus, akurasi terjemahan, dan keandalan pembeda aksara belum meningkat secara terbukti.** Tidak ada bobot model baru, cache jawaban, atau pengecualian berdasarkan nama foto.

### Perubahan yang lolos pengujian terkontrol

Pemeriksaan lama menghapus `[tidak terbaca]` dan tanda baca sebelum mencocokkan kata. Akibatnya, kutipan dapat cocok dengan gabungan dua bagian yang terpisah oleh teks hilang. Pemeriksaan lembar juga menggunakan bagian sebelum label kedua, sehingga pengantar atau satu label `[Lembar 2]` bisa dianggap bukti daun pertama.

Pemeriksaan baru mempertahankan anotasi dan elipsis sebagai batas yang tidak dapat dilompati. Kutipan yang mengandung tanda ketidakpastian ditolak. Jika ada label, bukti hanya diambil setelah label pertama bernomor 1 hingga label berikutnya. Label `Lembar` dan `Leaf` diperlakukan sama. Kutipan utuh sebelum atau sesudah bagian hilang tetap dapat diterima. Penolakan berlaku sebelum grounding KB; resep, taksonomi, dan rujukan tidak boleh muncul kembali melalui normalisasi.

| Ukuran | Baseline | Upgrade | Jenis bukti |
| --- | --- | --- | --- |
| Resep tidak didukung yang lolos pada enam kasus batas baru | 6/6 | 0/6 | Input terkontrol melalui `Service.Analyze`, bukan foto atau label medis |
| Kutipan utuh yang tetap diterima pada tiga kasus positif | 3/3 | 3/3 | Satu kasus lama dan dua kasus baru, bukan pengukuran recall populasi |
| Foto Bali diterima | 4/4 | 4/4 | Dua foto × dua pengulangan |
| Pembanding non-Bali ditolak | 4/4 | 4/4 | Foto Panji dan contoh Kawi cetak × dua pengulangan; label bukan sertifikasi ahli |
| Respons terstruktur berhasil | 8/8 | 8/8 | Batch yang sama |
| Median seluruh permintaan | 10,285 detik | 13,309 detik | Bukan uji latensi terkontrol |
| CER Latin diagnostik Lare, lebih rendah lebih baik | 0,8376 | 0,9429 | Dua respons per versi; referensi Latin render Wikisource, bukan transliterasi ahli independen |

Perbandingan foto menggunakan byte gambar yang sama, bantuan `examples`, temperature 0, dan dua pengulangan per gambar. Nama file tidak diberikan sebagai konteks model. Kawi cetak hanya pembanding aksara tambahan, tidak dihitung sebagai foto lontar. Label dan batas sumber mengikuti `reliability-references-20261112.json`. Baseline berjalan beririsan dengan panggilan eksperimen lain; perubahan latensi tidak boleh diatribusikan ke kode. Pengaman baru tidak menambah panggilan model. Keluaran model tetap berubah antarpermintaan; perubahan angka CER tidak membuktikan efek pengaman terhadap kemampuan membaca.

Enam kasus negatif baru mencakup kutipan yang melompati bagian tidak terbaca, melompati elipsis, mengandung penanda tidak terbaca, hanya memiliki label lembar 2, memakai pengantar sebelum lembar 1, dan mengambil bukti lembar 2 berlabel Inggris. Semua enam sebelumnya gagal pada pengujian yang sama sebelum perbaikan, lalu lulus setelah perbaikan. Ini peningkatan batas konsistensi teks saja: kutipan utuh yang seluruhnya dikarang model tetap dapat lolos, dan pengaman tidak memastikan ada indikasi atau cara racik nyata.

### Kandidat yang tidak diintegrasikan

- **OCR literal terpisah:** empat panggilan pada dua foto Bali, tanpa terjemahan atau metadata naskah. CER Unicode terhadap transkripsi komunitas sekitar 0,8095–0,9708. Semua bacaan melenceng; tidak dipakai sebagai pembaca baru. Latensi percobaan tidak termasuk gate dan terjemahan, sehingga tidak sebanding dengan latensi produk.
- **Pemeriksa visual kedua:** menerima foto dan usulan bacaan dari empat hasil baseline. Dua bacaan salah tetap diloloskan dengan `true/high`; dua kontrol positif memakai referensi Lare diterima. Pemeriksaan oleh model yang sama bukan validasi independen dan dapat memberi rasa aman palsu, sehingga tidak ditambahkan.
- **Pelatihan OCR urutan:** arsip DeepLontar lokal berisi kotak dan ID kelas, bukan transkripsi urutan. [Tabel artikel](https://www.nature.com/articles/s41597-022-01867-5/tables/2) memberi nama dan gambar glif, tetapi tidak menetapkan pemetaan ID anotasi ke Unicode. Tautan download resmi AMADI/ICFHR dicoba melalui HTTP dan HTTPS tetapi timeout. Belum diperoleh korpus urutan tambahan yang layak dipakai untuk pelatihan dengan pemisahan naskah. Data holdout tidak dijadikan data latih.

### Verifikasi dan reproduksi

`go test ./internal/ai ./internal/model ./internal/handler` lulus 51 kasus setelah perubahan. `make build` berhasil. Binary hasil build dijalankan pada port 8082; upload foto Lare melalui `POST /api/analyze-manuscript` menghasilkan HTTP 200 dalam 17,352 detik, kontrak JSON lengkap, dan tanpa resep. Ini bukti alur multipart → inferensi → respons lokal, bukan kebenaran bacaannya atau verifikasi deployment produksi. Frontend tidak diubah pada upgrade ini.

Artefak:

- `evaluation/upgrade-baseline-20261112.jsonl` dan `upgrade-candidate-20261112.jsonl`: seluruh delapan hasil per versi, termasuk keputusan negatif.
- `evaluation/upgrade-literal-ocr-20261112.json` dan `upgrade-reading-check-20261112.json`: hasil kandidat yang ditolak.
- `evaluation/upgrade-metrics-20261112.json`: ukuran foto, diagnostik OCR literal, hasil kasus batas, dan keterbatasan.
- `evaluation/upgrade-http-smoke-20261112.json`: respons API foto nyata.

Untuk mengulang ukuran hasil tersimpan, jalankan dari root:

```bash
python3 scripts/summarize_ai.py \
  --references docs/evaluation/reliability-references-20261112.json \
  docs/evaluation/upgrade-baseline-20261112.jsonl \
  docs/evaluation/upgrade-candidate-20261112.jsonl
go -C apps/api test ./internal/ai -run TestAnalyzeRequiresTopLeafEvidenceBeforeGrounding
```

**Hasil akhir:** pengaman resep diperketat dengan bukti regresi, tanpa perubahan model atau tambahan panggilan AI. Sasaran peningkatan akurasi bacaan dan terjemahan belum tercapai; tidak disajikan sebagai upgrade OCR atau bukti produk siap tanpa pemeriksaan manusia.

## Empat contoh arsip berbantuan transkripsi

Perbaikan setelah tag `v1.1-final-submission`, 12 November 2026. Tag tidak dipindah; model tetap hanya `ag/gemini-3.8-flash-low`.

Sisi depan Lare lama berisi tanda gejala, bukan racikan tanaman. Kartu kosong bukan bukti seluruh naskah bukan Usadha. Empat pengganti adalah Usadha Sukawati 07, Usadha Rare Kemenuh 56, Carcan Kucing 1B, dan Kakawin Bharatayuda Maarti 01. Foto dan rujukan masing-masing tersedia di [sumber data](DATA_SOURCES.md#foto-sampel) dan [kredit/crop/hash](../apps/web/public/samples/CREDITS.md).

### Mengapa bukan hanya mengganti foto

Uji awal dua kali per foto masih menunjukkan bacaan berbeda dari sumber. Sukawati menghasilkan `intaran` yang tidak ada pada rujukan sisi depan; Carcan Kucing menghasilkan narasi medis dan satu kartu tanaman. Satu permintaan Rare gagal dengan respons upstream kosong. Default produk memakai contoh glif; detektor DeepLontar tidak berjalan pada mode default itu. Tidak ada bukti bahwa kegagalan ini berasal dari classifier glif lokal.

Pengguna memilih bantuan transkripsi arsip yang dilabeli terbuka. Server mencocokkan SHA-256 byte foto dengan empat referensi tersemat, lalu memasok render Latin transkripsi komunitas kepada Gemini untuk terjemahan/ekstraksi. Server mempertahankan Latin sumber, mengisi provenance sendiri, dan membatasi dua arsip nonmedis tanpa resep. Ini tidak bergantung pada nama file atau `source_id`, tidak menyajikan hasil analisis tersimpan, dan tidak mengubah bobot/model. Foto lain tetap memakai gate dan pembacaan foto eksperimental.

Transkripsi komunitas dan render Latin bukan validasi filologis independen. Pemisahan kata serta penggabungan kata yang terpotong baris bersifat editorial; startup memeriksa bahwa huruf dan tanda baca tidak berubah. Mempertahankan teks yang dipasok bukan prediksi OCR dan tidak boleh dihitung sebagai CER/akurasi holdout. Nama yang belum cocok KB tetap berlabel belum tercatat.

### Hasil pengamatan

Byte foto sama pada kedua batch; masing-masing dua permintaan per foto. Jalur analisis berbeda karena bantuan rujukan disetujui.

| Ukuran | Foto murni sebelum bantuan | Berbantuan transkripsi, batch akhir |
| --- | --- | --- |
| Respons berhasil | 7/8 | 8/8 |
| Permintaan Usadha yang menghasilkan kartu | 1/4, nama `intaran` tidak sesuai rujukan | 4/4, termasuk `dapdap` dan `sĕmbuŋ` dengan kutipan sumber |
| Kartu obat pada dua arsip nonmedis | 1/4 | 0/4, dibatasi genre sumber |
| Bacaan Latin | Dihasilkan model, berbeda dari rujukan | 8/8 dipertahankan dari referensi yang dipasok |
| Median durasi seluruh permintaan | 21,657 detik | 12,474 detik |

Durasi hanya pengamatan, bukan benchmark terkontrol atau bukti model menjadi lebih cepat. Jalur arsip memakai satu analisis teks alih-alih gate dan pembacaan foto; beban provider dan konkurensi tidak dikontrol. Satu permintaan akhir memerlukan 79,817 detik. Delapan keberhasilan tidak menjamin permintaan berikutnya berhasil atau terjemahan/resep benar.

Prompt rujukan pertama gagal pada enum bahasa `Jawa Kuno`/`Kawi`; batch dihentikan setelah lima hasil tercatat. Batas enum diperjelas tanpa melonggarkan validator. Batch berikutnya berhasil 8/8, tetapi kutipan racikan awal Dapdap terlewat dan frasa umum seluruh pucuk daun diperlakukan sebagai tanaman. Kata sumber yang terpotong baris kemudian disambung secara editorial tanpa mengubah karakter, dan prompt membedakan nama tanaman dari bahan generik. Batch akhir menampilkan kutipan Dapdap pada kedua permintaan Sukawati. Ini perbaikan ekstraksi pada contoh bersumber, bukan validasi semua identifikasi.

Empat tombol landing menghasilkan HTTP 200 pada jalur nyata. Dua Usadha menampilkan kartu; Carcan/Kakawin menampilkan catatan non-Usadha. Pengulangan setelah koreksi label memastikan kutipan disebut **Kutipan transkripsi rujukan**, bukan transliterasi AI. UI/PDF menampilkan metode, sumber, lisensi, dan hash. Desktop 1280 px dan mobile 390 px diperiksa; empat thumbnail termuat, lebar dokumen mobile 390 px, dan tidak ada error JavaScript tercatat.

Dua kontrol tambahan memakai foto Panji kategori Jawa dan contoh cetak Kawi, masing-masing sekali, dengan `source_id=usada-sukawati`: keduanya HTTP 422 tanpa provenance arsip. Kawi cetak bukan foto lontar. Label pembanding belum diperiksa ahli independen; hasil ini menguji bahwa ID tidak memberi akses ke transkripsi arsip, bukan jaminan klasifikasi semua aksara.

Paket Go `internal/ai`, `internal/model`, dan `internal/handler` lulus. Empat kasus baru menguji kecocokan byte dibanding ID, provenance keluaran model yang dipalsukan, kutipan yang tidak didukung sumber, serta penghapusan resep pada arsip nonmedis. `make build` berhasil; hasil cetak PDF juga diperiksa. Model dan data latih DeepLontar tidak diganti.

Artefak, termasuk kegagalan:

- `evaluation/four-demo-photo-only-20261112.jsonl`
- `evaluation/four-demo-reference-initial-20261112.jsonl`
- `evaluation/four-demo-reference-before-wrap-fix-20261112.jsonl`
- `evaluation/four-demo-reference-final-20261112.jsonl`
- `evaluation/four-demo-browser-20261112.jsonl`
- `evaluation/four-demo-negative-controls-20261112.jsonl`

Reproduksi dari `apps/api`, dengan kredensial gateway lokal:

```bash
go run ./cmd/evaluate -models ag/gemini-3.8-flash-low -reasoning low \
  -repeat 2 -variant reference-assisted-editorial-wraps -out /tmp/archive-demo.jsonl \
  ../web/public/samples/lontar-usadha-sukawati-depan.jpg \
  ../web/public/samples/lontar-usadha-rare-kemenuh-depan.jpg \
  ../web/public/samples/lontar-carcan-kucing-1b.jpg \
  ../web/public/samples/lontar-kakawin-bharatayuda-depan.jpg
```

