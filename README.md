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

## Catatan Produksi
- Ganti `jwtSecret` di `backend/auth.go` sebelum dijual.
- Frontend bisa di-host terpisah (nginx) atau tetap disajikan Go (`GET /`).
- Untuk akses publik, buka port 80/443 di Tencent Security Group + pasang reverse proxy (Caddy/nginx) dengan HTTPS.
