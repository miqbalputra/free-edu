# KEPUTUSAN PIVOT (final — 29 Okt 2026)
1. Nama: tetap **SIMAKA** (tagline yang diubah, bukan nama)
2. Segmen: **dua-duanya** — sekolah umum + lembaga Islam (dua landing, satu kode)
3. Kurikulum: **Kurikulum Merdeka dulu** (K13 menyusul bila ada permintaan nyata)

---

# P1 — Generalisasi Bahasa & Alur (mulai sekarang)

## 1.1 Kamus istilah (satu sumber, dipakai semua teks)
| Lama | Baru (konteks umum) | Tetap (konteks islam) |
|---|---|---|
| santri | siswa | santri |
| ustadz/ustadzah | guru | ustadz |
| TPQ/tahfidz | sekolah | TPQ/tahfidz/pesantren |
| setoran hafalan | — (menu tersembunyi) | setoran hafalan |
| juz/murojaah di landing | — | tetap, di halaman segmen islam |

Implementasi: teks statis di frontend/landing/pitch diedit langsung; teks di aplikasi mengikuti `jenis` sekolah.

## 1.2 Backend: jenis sekolah
- Kolom baru `sekolah.jenis TEXT DEFAULT 'umum'` CHECK ('umum','islam') + migrate
- `handleDaftar`: field `jenis` (pilihan di form: "Sekolah Umum" / "Madrasah, TPQ & Pesantren")
- Claim/`/api/me` mengembalikan `jenis` → frontend pakai untuk menampilkan/menyembunyikan menu hafalan
- Existing sekolah (1,2,3): default 'islam' agar tidak berubah perilaku (kepala@simaka.id dst tetap lihat hafalan)
- Panel admin: tampilkan jenis di daftar sekolah

## 1.3 Frontend: menu kondisional
- `menuQuran = user.jenis==='islam'` → menu 📖 Hafalan + tombol "+ Catat Setoran Hafalan" + kartu hafalan di dashboard hanya muncul jika `menuQuran`
- Bahasa netral di semua halaman non-hafalan (cek: dashboard, siswa, nilai, rapor, SPP, keuangan, laporan, pengaturan)
- Label "SPP" tetap (umum juga pakai istilah SPP/komite — sudah netral)

## 1.4 Landing baru (2 segmen, satu halaman dengan anchor)
- Hero: **"SIMAKA — Sistem Admin Sekolah SD, SMP, SMA & sederajat"** + sub: *Termasuk TPQ, madrasah & pesantren — lengkap dengan modul Al-Qur'an*
- Kartu fitur: hafalan pindah ke posisi "Modul Pendidikan Al-Qur'an (untuk lembaga Islam)" — bukan fitur utama
- Section "Untuk TPQ & Pesantren" (anchor #islam): hafalan, murojaah, laporan wali santri — menjaga segmen lama
- Hero CTA demo tetap `?demo=1`
- Testimoni: tetap jujur ("bergabunglah sebagai yang pertama") di kedua segmen

## 1.5 Pitch & PWA
- `pitch.html`: buat 2 versi — `pitch.html` (umum) + `pitch-islam.html` (warisan versi lama)
- `manifest` PWA: name "SIMAKA — Sistem Admin Sekolah"

## 1.6 README
- Judul: "SIMAKA — Sistem Informasi Manajemen Sekolah (SD, SMP, SMA & Madrasah/TPQ)"
- Tabel perbandingan: ganti kolom "Aplikasi sekolah umum" vs "Excel+WA" tetap; konteks hafalan = "modul opsional"

## VERIFIKASI P1
- Daftar sekolah baru jenis 'umum' → menu Hafalan tidak ada; jenis 'islam' → ada
- Sekolah lama (islam) tidak berubah
- Landing & pitch tampil benar di HP
- Regression semua endpoint 200

---

# P2 — Rapor Kurikulum Merdeka (sesi berikutnya, setelah P1 rampung)
1. Jenjang saat daftar: SD/SMP/SMA → tingkat kelas dinamis (1–6 / 7–9 / 10–12); filter kelas mengikuti jenjang
2. Mapel default per jenjang saat sekolah umum daftar: SD (Pendidikan Agama & Budi Pekerti, Pendidikan Pancasila, B. Indonesia, Matematika, IPAS, SBDP, PJOK, B. Inggris), SMP (+IPS, Informatika, Prakarya), SMA (fase E/F: lebih banyak pilihan — beda per jurusan, v1 pakai set umum)
3. Rapor Merdeka ringkas: identitas → nilai mapel (NA, predikat) → ekstrakurikuler → ketidakhadiran → catatan wali kelas → keputusan (naik/lulus) — cetak berkop
4. Ekskul: CRUD + penilaian deskriptif masuk rapor (Batch P3 item 1 digabung di sini karena satu struktur data)
5. KKTP/fase penuh: TUNDA (catat sebagai backlog, bukan v1)
