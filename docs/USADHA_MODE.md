# Alur Usadha

Alur Usadha adalah lapisan eksperimental untuk foto lontar beraksara Bali yang diduga memuat resep tanaman obat. Alur ini berjalan otomatis ketika hasil model lolos pemeriksaan kutipan; pengguna tidak perlu memilih mode. Tag final submission: [`v1.1-final-submission`](https://github.com/0xMinomus/bukalontar/tree/v1.1-final-submission).

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

Pada `main` setelah tag, Usadha Sukawati dan Rare Kemenuh memakai transkripsi komunitas ketika hash foto cocok persis. Bacaan sumber dipasok ke Gemini untuk terjemahan/ekstraksi; UI dan PDF menyebut **Berbantuan transkripsi rujukan**. Carcan Kucing dan Kakawin Bharatayuda dari registry tidak menghasilkan kartu resep. Foto lain tetap memakai pembacaan foto eksperimental. Ini tidak meningkatkan atau membuktikan akurasi OCR umum.

## Keluaran

Setiap resep memakai struktur yang sama:

| Field | Isi |
| --- | --- |
| `plant_local_name` | Nama tanaman dari teks lontar |
| `source_quote` | Kutipan utuh dari bacaan daun pertama: transliterasi AI untuk foto baru atau transkripsi komunitas untuk arsip cocok; bukan verifikasi ahli |
| `plant_scientific_name` | Nama ilmiah jika cocok dengan knowledge base |
| `botanical_family` | Famili botani jika tersedia |
| `kb_match` | Status pencocokan dengan Taru Pramana |
| `plant_part_used` | Daun, kulit kayu, akar, buah, bunga, atau campuran |
| `energy_character` | Panas, dingin, atau sedang |
| `formulation_type` | Loloh, boreh, tutuh, atau sasap |
| `adjuvants` | Bahan tambahan yang namanya harus muncul pada kutipan yang sama |
| `target_ailment` | Indikasi yang disebutkan teks |
| `preparation_instructions` | Cara racik menurut teks |
| `phytochemical_compounds` | Senyawa dari knowledge base |
| `modern_pharmacological_action` | Ringkasan aksi farmakologis dari knowledge base |
| `reference_studies` | Rujukan yang tersedia di knowledge base |
| `reference_links` | Tautan literatur atau pencarian literatur yang dikurasi dalam knowledge base; tiap objek berisi `title`, `url`, dan `kind` |

`recipes[]` dapat berisi lebih dari satu resep dari satu lembar. Hasil dibatasi maksimal enam resep. Persiapan berbeda dengan tanaman sama tidak dihapus berdasarkan nama tanaman saja; hasil tidak lagi digabung dari dua pembacaan.

## Grounding Taru Pramana

Knowledge base berada di `apps/api/internal/kb/data/taru-pramana.json`. Aturan grounding:

- Nama ilmiah, famili, senyawa, aksi farmakologis, dan rujukan diambil dari knowledge base saat ada kecocokan. Sifat tradisional (`energy_character`) mengikuti bacaan model, bukan nilai bawaan KB.
- Jika tidak ada kecocokan, nama ilmiah dan famili dikosongkan.
- Rujukan tidak boleh dikarang.
- Nama tanaman yang tidak dikenal tetap ditampilkan sebagai nama lokal dengan label belum tercatat.
- Detail yang tidak disebutkan tetap tidak disebutkan; bentuk racikan tidak otomatis menjadi loloh.
- Tautan literatur berasal dari knowledge base, bukan dari model atau naskah; hasil pencarian tidak membuktikan khasiat maupun keamanan resep.

Grounding ini penting karena teks lontar dapat memuat nama tanaman yang mirip dengan tanaman lain. Model tidak diberi kebebasan untuk menebak taksonomi.

## Pemeriksaan Kutipan

Resep memerlukan `source_quote` minimal empat kata yang cocok dengan satu bagian utuh bacaan lembar pertama, menyebut nama tanaman dan seluruh bahan tambahan yang diklaim. Tanda `[tidak terbaca]`, anotasi dalam kurung siku, dan elipsis memisahkan bukti; model tidak boleh menggabungkan kata di kedua sisinya. Jika label lembar tersedia, kutipan hanya dapat berasal setelah label pertama bernomor 1 hingga label berikutnya, bukan dari pengantar atau lembar 2.

Jika kutipan tidak memenuhi aturan, server menghapus resep beserta grounding dan rujukannya sebelum respons dikirim. Kutipan yang cocok tetap tidak membuktikan bahwa teks sesuai foto atau memuat cara racik yang benar. Pada enam kasus batas terkontrol, pengaman terakhir mengurangi resep tidak didukung yang lolos dari 6/6 menjadi 0/6, sementara tiga kutipan utuh tetap diterima 3/3. Ini bukan akurasi resep pada foto; lihat [evaluasi](AI_EVALUATION.md#upgrade-pengaman-kutipan-dan-pembanding-baseline).

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
- Model dapat mengarang bacaan dan kutipan yang saling cocok. Sistem belum memvalidasi kebenaran filologis atau medisnya.

## Rujukan

- [Konsep produk](PROJECT_CONCEPT.md)
- [Arsitektur](ARCHITECTURE.md)
- [Sumber data](DATA_SOURCES.md)
- [AI Disclosure](AI_DISCLOSURE.md)
- [Indeks dokumentasi](README.md)
