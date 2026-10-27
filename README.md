# SIMAKA — Sistem Informasi Manajemen Sekolah (All-in-One, Multi-Sekolah SaaS)

Aplikasi web sekolah Indonesia: Kepala Sekolah, Guru, Tata Usaha.
Stack: **Go (backend) + Vue 3 (frontend) + PostgreSQL**.

## Fitur
- 15+ modul: siswa, absensi (siswa & guru), nilai (+massal), rapor cetak-PDF, hafalan, jadwal, SPP, keuangan, agenda, pengumuman, log aktivitas
- Multi-tenant: sekolah daftar sendiri (trial 14 hari), data terisolasi per sekolah
- Notifikasi WhatsApp tagihan SPP (link wa.me siap kirim)
- Ekspor Excel (absensi, SPP, nilai, buku kas)
- Landing page penjualan + PWA (bisa di-install di HP)
- Keamanan: bcrypt, JWT auto-secret, rate limit login, audit trail, backup harian otomatis

## Fitur per Role

| Modul | Kepala Sekolah | Guru | TU |
|---|---|---|---|
| Dashboard statistik | ✅ | ✅ | ✅ |
| Data siswa (CRUD) | ✅ | lihat | ✅ |
| Absensi siswa (H/S/I/A) | ✅ | ✅ | ✅ |
| Absensi guru (masuk/pulang) | ✅ | ✅ (diri sendiri) | ✅ |
| Nilai (Tugas/UH/PTS/PAS) | ✅ | ✅ | lihat |
| Hafalan (surah, ayat, status) | ✅ | ✅ | lihat |
| Jadwal pelajaran | ✅ | lihat | ✅ |
| SPP (generate tagihan, bayar) | ✅ | lihat | ✅ |
| Keuangan (pemasukan/pengeluaran) | ✅ | — | ✅ |
| Agenda & pengumuman | ✅ | ✅ | ✅ |
| Manajemen pengguna | ✅ | — | — |

## Menjalankan

```bash
# 1. Database (sekali saja)
sudo -u postgres psql -c "CREATE USER simaka WITH PASSWORD 'simaka123';"
sudo -u postgres psql -c "CREATE DATABASE simaka OWNER simaka;"

# 2. Backend (migrasi & seed data demo otomatis)
cd backend && go build -o simaka . && ./simaka
# aplikasi: http://localhost:8080

# Variabel opsional: DATABASE_URL, PORT
```

## Akun Demo
- Kepala Sekolah: `kepala@simaka.id` / `kepala123`
- Guru: `guru@simaka.id` / `guru123`
- Tata Usaha: `tu@simaka.id` / `tu123`

## Changelog

### v1.1 — Rilis Keamanan & Fitur Lengkap

**Keamanan (audit penuh 21 temuan → semua diperbaiki)**
- 🔒 Isolasi tenant di semua ekspor (SPP, nilai, absensi, keuangan), rapor cetak, kwitansi, dan daftar siswa per kelas
- 🔒 Cek kepemilikan data di semua update/hapus (siswa, nilai, hafalan, jadwal, agenda, pengeluaran, pengguna)
- 🔒 Reset password hanya dalam sekolah yang sama
- 🔒 Validasi nilai 0–100 (satuan & massal)
- 🔒 Generate SPP difilter per sekolah
- 🔒 Security headers: CSP, X-Frame-Options, nosniff, Referrer-Policy, Permissions-Policy
- 🔒 CORS whitelist (bukan `*` lagi)
- 🔒 Rate-limit login persisten ke database (lock 15 menit setelah 5 gagal, bertahan antar restart)

**Fitur baru**
- 📊 Laporan rekap bulanan: absensi per kelas, SPP, kas, hafalan, rata-rata nilai
- 🧑‍💼 Wali kelas per kelas
- 🧾 Kwitansi SPP siap cetak (kop sekolah, nomor kwitansi, tanda tangan)
- 🗑️ Hapus data: nilai, hafalan, jadwal, agenda, pengeluaran, pengguna
- 📱 Ikon PWA 192/512 + manifest lengkap
- 📅 Tahun ajaran dinamis dari pengaturan sekolah
- 🔄 Transaksi database (hapus siswa & generate SPP all-or-nothing)
- 🧹 Rotasi otomatis log aktivitas

### v1.0 — Rilis Awal
- Manajemen siswa, nilai, hafalan, jadwal, agenda, pengumuman, absensi siswa & guru
- SPP multi-bulan + notifikasi WhatsApp tagihan
- Keuangan: pemasukan, pengeluaran, buku kas + ekspor Excel
- Rapor lengkap (nilai, kehadiran, catatan guru, predikat) siap cetak
- Multi-sekolah SaaS: pendaftaran mandiri, trial 14 hari, gating langganan
- Tema warna per peran (kepala sekolah, guru, tata usaha)
- PWA (bisa dipasang di HP), landing page

## Catatan Produksi
- Ganti `jwtSecret` di `backend/auth.go` sebelum dijual.
- Frontend bisa di-host terpisah (nginx) atau tetap disajikan Go (`GET /`).
- Untuk akses publik, buka port 80/443 di Tencent Security Group + pasang reverse proxy (Caddy/nginx) dengan HTTPS.

## Lisensi
**PolyForm Strict License 1.0.0** — lihat [LICENSE](LICENSE).
- ✅ Bebas dipakai siapapun, untuk tujuan apa pun (termasuk komersial/internal sekolah)
- ❌ Dilarang memodifikasi kode atau mendistribusikan versi modifikasi
- 🤝 Kontribusi dipersilakan via pull request
