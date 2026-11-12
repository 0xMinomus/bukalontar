# Detektor Teks dan Pengenal Glif Lontar

Dokumen ini mencatat artefak pembelajaran mesin lokal yang membantu pembacaan lontar. Keduanya dilatih dari data berlabel ahli, bukan dari jawaban model bahasa, dan **bukan** transkripsi halaman penuh.

## Peran dalam produk

- Detektor teks menemukan area berteks pada foto lontar, lalu mengirim maksimal dua perbesaran konteks ke model multimodal.
- Pengenal glif mengklasifikasikan potongan aksara yang sudah dilokalisasi ahli ke ID kelas opasitas DeepLontar.
- Model bahasa tetap membaca foto asli. Hasil detektor bukan teks terbaca, bukan terjemahan, dan tidak mengidentifikasi tanaman obat.
- ID kelas belum dipetakan ke nama aksara oleh sumber primer, sehingga API runtime tidak pernah mengeluarkan ID tersebut sebagai bacaan.

## Sumber dan lisensi

- [DeepLontar v2](https://doi.org/10.6084/m9.figshare.20103803.v2), CC BY 4.0, oleh Sutramiani dkk.
- 200 citra asli beraksara Bali. Satu berkas anotasi (`76a.txt`) dikarantina karena ID kelas di luar rentang 0–54.
- Pembagian kelompok memakai nomor lembar fisik: kedua sisi (`a`/`b`) satu lembar selalu berada di kelompok yang sama. Train 159 lembar, dev 20, test 20.
- Label hanya lokasi dan kelas karakter. Tidak ada transkripsi halaman, terjemahan, atau label resep pada dataset ini.

## Artefak

| Artefak | Isi | SHA-256 |
| --- | --- | --- |
| `apps/api/internal/ocr/model.json` | Regresi ridge 10 bobot untuk okupansi teks per ubin 30×10 | `bea6757cae52f9f51ab982eb67aa237294aff1556b836c68dbbb25e9d6b2ea56` |
| `apps/api/internal/ocr/hog_model.json` | LinearSVC pada HOG 4×4 sel, 8 bin, fitur L2 | `8af6cdc8d86fe7d5601472ed082848bdace1064da3cc50f93135f61fdebd49c8` |

Keduanya di-embed dengan `go:embed`; inferensi berjalan dengan Go standard library tanpa layanan tambahan.

## Hasil uji tahan-lihat

| Ukuran | Hasil | Arti |
| --- | --- | --- |
| Okupansi ubin, dev | presisi 78,1%; recall 86,9%; F1 82,3% | Kualitas lokalisasi area berteks pada 20 citra dev |
| Jendela fokus, test | 30% area dipilih; 72,1% ubin terpilih memuat pusat karakter; recall 37,8% ubin berkarakter | Perbesaran kontekstual, bukan pengenalan penuh |
| Klasifikasi glif kotak ahli, test | top-1 82,2% | Hanya potongan yang sudah dilokalisasi; ID kelas opasitas |
| Klasifikasi glif baseline centroid | top-1 57,0% | Pembanding sebelum HOG/LinearSVC; terpilih berdasarkan dev |

**Batasan:** angka ini bukan akurasi transliterasi, terjemahan, atau resep. Generalisasi ke foto Wikisource dan foto ponsel belum diukur pada metrik berlabel. Jendela fokus hanya menutup sebagian area berkarakter, sehingga model bahasa tetap membaca foto asli.

## Pelatihan ulang

```bash
python3 -m pip install numpy scipy Pillow scikit-learn
OPENBLAS_NUM_THREADS=1 OMP_NUM_THREADS=1 python3 scripts/train_lontar_ocr.py
```

Hasil lengkap, daftar kelompok, dan provenance ditulis ke `data/ocr/evaluation.json`. Fixture parity Go/Python dibuat dengan:

```bash
cd scripts && python3 make_hog_parity.py
```

Parity diuji oleh `TestGlyphModelPythonParity`: kelas dan margin harus sama dalam 1e-4. Fixture menyimpan sampel 16×16 yang sama dengan pipeline pelatihan, sehingga perubahan fitur di salah satu sisi akan gagal.

## Runtime

- Menerima JPEG/PNG hingga 24 MP; menolak dimensi di luar batas.
- Foto dengan latar seragam di sudut diproses untuk mencari strip daun horizontal; jika tidak ada kandidat, runtime abstain dengan peringatan.
- Crop fokus dibatasi dua per permintaan, masing-masing maksimal 1024 px lebar, JPEG kualitas 85.
- Hasil detector disertakan hanya sebagai bukti bantuan. Prompt melarang menyalin label contoh atau menduplikasi isi crop fokus.

## Provenance contoh glif

55 contoh glif dari Tabel 1 artikel [Scientific Data](https://doi.org/10.1038/s41597-022-01867-5) (CC BY 4.0) di-embed di `apps/api/internal/readingexamples`. Contoh disertakan sebagai referensi visual, bukan teks target. Artikel tidak menyebut identitas lembar tiap glif, sehingga contoh tidak dipakai untuk mengevaluasi citra DeepLontar dan tidak menjamin independensi terhadap data uji.

## Ground truth halaman

`data/reading-examples/holdout/` berisi lima halaman pindaian Wikisource dengan transkripsi Bali yang ditinjau komunitas (revisi tertaut, lisensi CC BY-SA 4.0; pindaian Commons menyatakan PD-ID-exempt). Status peninjauan Wikisource bukan bukti validasi ahli independen. Semua dari satu naskah (*Aji Kreket*), disimpan sebagai holdout dan tidak dikirim sebagai contoh konteks. Independensi dari sumber glif Tabel 1 tidak dapat dibuktikan karena artikel tidak menyebut identitas lembar.
