# Alur Usadha

Alur Usadha adalah lapisan khusus untuk lontar yang memuat resep tanaman obat. Alur ini berjalan otomatis ketika naskah memuat resep. Pengguna tidak perlu memilih mode sebelum mengunggah foto.

## Ruang Lingkup

Alur ini mencakup:

- identifikasi nama tanaman lokal Bali;
- pencocokan dengan katalog Taru Pramana;
- bagian tanaman yang dipakai;
- sifat tradisional;
- bentuk racikan;
- indikasi dan cara racik;
- senyawa fitokimia dan rujukan yang tersedia;
- disclaimer medis.

Alur ini tidak memberi diagnosis, dosis, atau rekomendasi pengobatan. Hasil AI adalah bahan bacaan dan penelitian, bukan nasihat medis.

## Keluaran

Setiap resep memakai struktur yang sama:

| Field | Isi |
| --- | --- |
| `plant_local_name` | Nama tanaman dari teks lontar |
| `plant_scientific_name` | Nama ilmiah jika cocok dengan knowledge base |
| `botanical_family` | Famili botani jika tersedia |
| `kb_match` | Status pencocokan dengan Taru Pramana |
| `plant_part_used` | Daun, kulit kayu, akar, buah, bunga, atau campuran |
| `energy_character` | Panas, dingin, atau sedang |
| `formulation_type` | Loloh, boreh, tutuh, atau sasap |
| `adjuvants` | Bahan non-tanaman pada resep |
| `target_ailment` | Indikasi yang disebutkan teks |
| `preparation_instructions` | Cara racik menurut teks |
| `phytochemical_compounds` | Senyawa dari knowledge base |
| `modern_pharmacological_action` | Ringkasan aksi farmakologis dari knowledge base |
| `reference_studies` | Rujukan yang tersedia di knowledge base |

`recipes[]` dapat berisi lebih dari satu resep dari satu lembar. Hasil dibatasi maksimal enam resep dan dideduplikasi berdasarkan nama tanaman lokal.

## Grounding Taru Pramana

Knowledge base berada di `apps/api/internal/kb/data/taru-pramana.json`. Aturan grounding:

- Nama ilmiah, famili, senyawa, rujukan, dan energy character diambil dari knowledge base saat ada kecocokan.
- Jika tidak ada kecocokan, nama ilmiah dan famili dikosongkan.
- Rujukan tidak boleh dikarang.
- Nama tanaman yang tidak dikenal tetap ditampilkan sebagai nama lokal dengan label belum tercatat.

Grounding ini penting karena teks lontar dapat memuat nama tanaman yang mirip dengan tanaman lain. Model tidak diberi kebebasan untuk menebak taksonomi.

## Menangani Naskah Non-Usadha

Jika naskah bukan Usadha:

- transliterasi dan terjemahan tetap dihasilkan;
- `recipes[]` boleh kosong;
- field tanaman memakai penanda "tidak disebutkan";
- UI menampilkan catatan naskah, bukan kartu resep;
- disclaimer medis tetap ditampilkan sebagai bagian dari hasil.

Perilaku ini menjaga agar lontar Kakawin, Babad, Geguritan, atau tutur tidak dipaksa menjadi resep obat.

## Keamanan dan Etika

- Disclaimer medis ditulis oleh server, bukan oleh model.
- Model tidak boleh menambahkan tanaman yang tidak terlihat pada teks.
- Model tidak boleh membuat nama ilmiah, famili, atau sitasi.
- Hasil dengan keterbacaan rendah harus ditampilkan sebagai hasil yang perlu diperiksa.
- Pengguna diarahkan untuk memeriksa lempir asli dan rujukan filologi.

## Batasan

- Deteksi resep bergantung pada kualitas foto dan kejelasan teks.
- Satu lembar dapat memuat lebih dari enam resep, tetapi kontrak saat ini membatasi enam.
- Knowledge base belum mencakup semua tanaman obat Bali.
- Model dapat membaca baris yang berbeda pada permintaan yang berbeda.
- Istilah penyakit tradisional tidak selalu memiliki padanan medis modern yang pasti.

## Rujukan

- [Konsep produk](PROJECT_CONCEPT.md)
- [Arsitektur](ARCHITECTURE.md)
- [Sumber data](DATA_SOURCES.md)
- [AI Disclosure](AI_DISCLOSURE.md)
- [Indeks dokumentasi](README.md)
