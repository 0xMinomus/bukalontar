# Improvement dan Bukti Final

Fokus komunikasi: pengaman resep dan ketertelusuran sumber. Jangan menyamakan pengaman hasil dengan peningkatan kemampuan membaca aksara.

| Perubahan | Bukti terdokumentasi | Batas klaim |
| --- | --- | --- |
| Kutipan resep harus utuh dan tidak melintasi bagian tidak terbaca atau lembar lain | Enam kasus batas yang sebelumnya lolos kini ditolak; tiga kutipan utuh tetap diterima | Uji terkontrol pengaman, bukan akurasi OCR atau keamanan medis |
| Empat contoh arsip bersumber dengan pencocokan SHA-256 | Delapan respons dari empat foto × dua pengulangan; bacaan Latin sesuai referensi yang dipasok | Bukan delapan keberhasilan OCR; bacaan berasal dari transkripsi rujukan |
| Pemisahan Usadha dan non-Usadha pada arsip | Empat permintaan Usadha memiliki kartu; empat permintaan nonmedis tanpa kartu | Sampel kecil; tidak menjamin klasifikasi semua naskah |
| Pemeriksaan aksara terpisah pada foto baru | Tersedia hasil uji Bali dan pembanding dalam laporan | Masih ada salah terima dan salah tolak; jangan klaim Bali-only terjamin |
| Tautan PubMed dan grounding Taru Pramana | Tautan KB tampil pada kartu, rujukan, dan PDF menurut status proyek | Literatur tanaman tidak otomatis membuktikan khasiat resep naskah |
| Pelatihan detektor dan pengenal glif lokal | 200 citra dikumpulkan; training memakai 159 citra berlabel | Bukan fine-tuning Gemini, bukan OCR kalimat, bukan bukti kenaikan akurasi halaman |
| Label metode pada UI/PDF dan penyembunyian Unicode eksperimental | Dokumentasi serta pemeriksaan UI/PDF tercatat | Transparansi hasil, bukan peningkatan akurasi model |

## Kalimat utama untuk presentasi

“Pada tahap final, kami memperketat bukti kutipan resep dan memperjelas asal bacaan. Untuk contoh arsip, kami menampilkan transkripsi rujukan secara terbuka. Pembacaan foto baru tetap eksperimental dan perlu diperiksa manusia.”

## Angka yang tidak digunakan sebagai klaim utama

- Akurasi glif 82,2% adalah pengenalan potongan karakter, bukan akurasi membaca lontar.
- Perbandingan latensi 22,2 menjadi 14,0 detik adalah eksperimen historis, bukan SLA atau waktu pasti pipeline final.
- Hasil penolakan pada sampel kecil bukan jaminan pembatasan aksara 100%.
- Banyaknya citra bukan bukti Gemini sudah dilatih ulang.

## Bukti sumber

- [Status terbaru dan uji empat contoh](../docs/STATUS.md)
- [Evaluasi pengaman, eksperimen, dan contoh arsip](../docs/AI_EVALUATION.md)
- [Disclosure metode](../docs/AI_DISCLOSURE.md)
- [Training lokal](../docs/OCR_TRAINING.md)
- [Sumber data](../docs/DATA_SOURCES.md)

Materi merangkum hasil pengujian yang sudah dicatat. Penyusunan dokumen ini tidak menjalankan pengujian AI baru.
