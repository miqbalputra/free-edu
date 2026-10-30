# RENCANA PIVOT: SIMAKA → Sekolah Umum (SD, SMP, SMA/sederajat)

Tanggal: Okt 2026 · Status: DRAFT — menunggu keputusan Iqbal

## 1. Kenapa pivot ini masuk akal (dan risikonya)

**Masuk akal:**
- Pasar sekolah umum jauh lebih besar: ~450 ribu sekolah SD/SMP/SMA di Indonesia vs TPQ/tahfidz (puluhan ribu, tak tercatat resmi)
- Sakit inti sama: nilai & rapor manual, absensi kertas, SPP/komite bocor, laporan ke ortu lewat grup WA kacau
- Kode SIMAKA sudah 90% netral: siswa, kelas, nilai, rapor, absensi, SPP, kwitansi, laporan — semua dipakai sekolah umum juga
- Fitur hafalan tinggal jadi **modul opsional** yang disembunyikan — bukan dihapus

**Risiko jujur:**
- Persaingan jauh lebih ketat: e-school, Erkasih, BOS aplikasi, MySCH.id, dsb. — pasar sekolah umum sudah penuh pemain besar
- Sekolah umum punya **kurikulum resmi (Merdeka)**: mapel, KKTP/CP-TP, wali kelas, ekstrakurikuler — rapor harus ikut format resmi. Ini pekerjaan besar yang TPQ tidak butuhkan
- Keunggulan SIMAKA yang paling beda (hafalan + laporan wali tahfidz) jadi tidak membedakan
- Rekomendasi laporan konsultan tetap berlaku: tanpa bukti sosial & import data, pivot tidak menolong

---

## 2. Prinsip pivot: "General core + modul agama"

Ubah arsitektur produk dari "aplikasi tahfidz" menjadi:

```
SIMAKA Core (semua sekolah): siswa, kelas, absensi, nilai, rapor,
  jadwal, SPP/komite, keuangan, kwitansi, laporan, WA ortu, unduh data
    ├── Modul HAFALAN/QURAN (toggle per sekolah) — aktif di TPQ/tahfidz/pesantren
    ├── Modul EKSTRAKURIKULER (sekolah umum butuh ini untuk rapor)
    └── (nanti) Modul PPDB
```

- Tabel `sekolah` tambah kolom `jenis` ('umum' | 'islam') → **sekolah pilih saat daftar**, ini yang menentukan menu hafalan muncul/tidak
- Semua yang sudah dibangun TIDAK dibuang — hafalan tetap jadi pembeda di segmen islami, cuma tidak menghalangi pasar umum

## 3. Batch perubahan (urutan kerja)

### Batch P1 — Generalisasi Bahasa & Alur (fondasi, paling cepat)
1. **Istilah**: audit semua teks frontend/landing/pitch — ganti "santri" → "siswa", "ustadz/ustadzah" → "guru", "TPQ/tahfidz" → "sekolah", "setoran hafalan" → hanya muncul di modul islam. Buat kamus istilah di satu file agar konsisten.
2. **Landing page baru**: hero "Sistem Admin Sekolah — SD, SMP, SMA, sederajat" + sub "termasuk TPQ & pesantren (modul Quran)". Fitur hafalan turun jadi satu kartu "Modul Pendidikan Al-Qur'an (opsional)".
3. **Pendaftaran**: pertanyaan baru "Jenis lembaga?" → Sekolah Umum / Madrasah & TPQ (menentukan menu). Sekolah umum → menu hafalan tersembunyi, mapel default terisi umum (MTK, B.Indo, IPA...); islam → hafalan aktif.
4. **Pitch 1 halaman**: 2 versi (umum & islami) dari satu sumber.
5. **README + repo**: rename judul "SIMAKA — Sistem Admin Sekolah (SD/SMP/SMA & lembaga Islam)".

### Batch P2 — Rapor Format Sekolah Umum (terberat, paling menentukan)
Sekolah umum menuntut rapor sesuai kurikulum. Minimal viable:
1. **Mapel dinamis** per jenjang (sudah ada guru_mapel — rapikan UI) — tidak boleh hardcode
2. **Struktur rapor Kurikulum Merdeka**: identitas → nilai per mapel (NA + predikat) → ekstrakurikuler → ketidakhadiran → catatan wali kelas → keputusan naik kelas
3. **Tingkat kelas dinamis**: SD (1–6), SMP (7–9), SMA (10–12) — sekolah pilih jenjang saat daftar; filter kelas mengikuti
4. **Cetak rapor per semester + kelengkapan tanda tangan** (wali kelas + kepala sekolah) — sudah ada polanya, tinggal disesuaikan format
5. ⚠️ *Tunda dulu*: KKTP/fase/CP-TP penuh — terlalu dalam untuk v1; format ringkas Kurikulum Merdeka cukup untuk sekolah kecil/menengah, dan format K13 lama tetap tersedia sebagai pilihan

### Batch P3 — Kebutuhan Khas Sekolah Umum
1. **Ekskul sebagai nilai rapor** (wajib di kurikulum): CRUD ekskul + penilaian "Sangat Baik–Perlu Bimbingan" masuk rapor
2. **PPDB ringan** (tunda sampai musim — Juni–Juli): pendaftaran online sederhana → masuk jadi calon siswa. Ini fitur penarik #1 sekolah umum menjelang tahun ajaran baru
3. **Dana komite/SPP tetap pakai modul SPP yang ada** — ganti label "SPP" jadi "Iuran/PPOB" yang bisa di-custom namanya per sekolah
4. **Kenaikan kelas & kelulusan** (akhir tahun ajaran): naik batch siswa kelas 1→2 dst + lulus untuk kelas akhir — sudah masuk audit lama, waktunya jadi

### Batch P4 — Posisi Pasar
1. Uji dengan **2–3 sekolah umum nyata** (tetangga/kenalan guru) gratis selama setahun → testimoni pertama di segmen baru
2. Harga tetap: Gratis ≤30 siswa / Sekolah Rp45rb ≤150 siswa — strukturnya sudah benar
3. Landing: tambah halaman/bagian khusus "Untuk TPQ & Pesantren" agar segmen lama tetap terlayani, bukan terasa ditinggalkan
4. Program rujukan jalan di dua segmen (guru umum ↔ guru umum, ustadz ↔ ustadz)

## 4. Yang TIDAK diubah (disiplin)
- Semua fitur keamanan, multi-tenant, anti-sandera, laporan WA ortu, import, panel admin — tetap
- Modul hafalan TIDAK dihapus — TPQ/pesantren tetap terlayani (ini justru pembeda di segmen itu)
- Struktur harga & anti-sandera tetap
- Multi-cabang tetap ditunda

## 5. Urutan & estimasi
P1 (1 sesi kerja) → P2 (2–3 sesi, terberat) → P3 sesuai musim (ekskul kapan saja; PPDB menjelang Juni; kenaikan kelas menjelang Juni–Juli) → P4 paralel dengan validasi.

## 6. Keputusan yang perlu Iqbal sebelum mulai
1. **Nama produk**: tetap SIMAKA? (SIMAKA = "Sistem Informasi Manajemen Sekolah" — sudah netral, bagus. Tagline yang diubah)
2. **Segmen prioritas penjualan**: umum saja, atau dua-duanya (rekomendasi: dua-duanya — kode sama, dua landing page)
3. **Kurikulum**: target Merdeka saja, atau sediakan juga format K13 untuk sekolah yang masih transisi (rekomendasi: Merdeka dulu, K13 menyusul bila ada permintaan nyata)
