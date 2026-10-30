# Rencana Perbaikan PRODUK SIMAKA — dari Laporan Konsultan Riset Pasar

Tanggal: 29 Okt 2026 · Sumber: penilaian independen (6 persona + sintesis)

## Sumber rekomendasi

**Dari "Tiga penolakan teratas → perbaikan PRODUK":**
1. Laporan mingguan otomatis ke WA pengasuh (perbaikan produk utk "siapa yang isi datanya?")
2. Import CSV/Excel massal (perbaikan produk utk "data lama mau pindah gimana?")

**Dari "4. Celah produk yang menahan pembelian" (urut paling menghalangi):**
1. Tidak ada import data massal → dinding input manual untuk lembaga >30 santri
2. Multi-cabang diiklankan Rp200rb tapi belum ada → risiko kredibilitas → JANGAN dijual dulu
3. Token laporan wali → tampilkan "tautan khusus untuk Anda, berlaku s/d [tanggal]" di halaman laporan
4. Tidak ada bukti eksternal → bukan produk (di luar plan ini)

**Yang disengaja TIDAK masuk plan (disiplin scope, dari laporan):**
- Multi-cabang → baru setelah import + retensi selesai, dan hanya jika ada calon pembeli Yayasan nyata
- Jangan jual paket Yayasan sebelum fiturnya ada

---

## Batch 1 — Import Data Massal (celah #1, paling menghalangi)
*Kunci masuk segmen menengah/yayasan 100–300 santri — satu-satunya segmen yang benar-benar mau bayar.*

1. **Import siswa dari Excel/CSV** di menu Siswa:
   - Upload file → preview tabel (10 baris pertama) → petakan kolom via dropdown sederhana (nama, kelas, NIS, nama ortu, HP ortu) → konfirmasi → import.
   - Terima .xlsx dan .csv (parsing xlsx tanpa dependensi berat: kolom dipisah koma/titik koma untuk csv; xlsx pakai library Go minimal atau minta user simpan sebagai CSV dengan petunjuk jelas).
   - Validasi manusiawi: "Baris 7: nama kosong — dilewati", "Baris 12: kelas '1 a' diperbaiki jadi '1A'". Ringkasan akhir: "23 siswa masuk, 2 dilewati, 1 diperbaiki".
   - Anti-duplikat: cek NIS + nama; duplikat → tawarkan "perbarui" atau "lewati".
   - Template CSV siap-unduh ("Unduh contoh format") berisi 3 baris contoh + kolom yang benar.
2. **Import nilai & tagihan SPP dari CSV** (pola sama, tahap kedua di batch yang sama bila waktu cukup; kalau tidak, siswa dulu — siswa adalah dinding utamanya).

Verifikasi: import 30 baris CSV uji → 30 masuk, duplikat terdeteksi, kolom kacau → pesan jelas; lintas tenant tetap terisolasi.

## Batch 2 — Laporan Mingguan Otomatis ke WA (penyebab churn #1: "aplikasi mati karena tak ada yang isi")
*Pengasuh merasa "dikawal" tanpa perlu login. Data sudah ada — tinggal rangkum + link.*

1. **Rangkuman mingguan per sekolah** (endpoint internal): absensi dicatat X/Y hari, hafalan Z setoran, tagihan lunas/bulan ini, siswa tanpa data 7 hari terakhir.
2. **Kirim otomatis**: goroutine harian (jam 20.00 WIB) → sekolah yang punya aktivitas 7 hari terakhir → susun pesan WA singkat bahasa awam + tombol/link "Lihat laporan lengkap".
3. **Tombol "Kirim laporan mingguan ke WA saya"** manual di dashboard (untuk demo & pemakaian awal sebelum otomatis jalan).
4. Nomor WA tujuan: pakai `hp_ortu`-style pengaturan baru "Nomor WA pengasuh" di Profil Sekolah (kosong → tidak dikirim, tidak error).
5. Pengingat sisi guru yang sudah ada (dashboard "absensi hari ini belum dicatat") dipertahankan — laporan ini untuk PENGASUH, bukan guru.

Verifikasi: sekolah dengan data → pesan tersusun benar; sekolah tanpa aktivitas → tidak dikirim; nomor kosong → skip aman.

## Batch 3 — Kepercayaan Laporan Wali + Kredibilitas Harga (celah #3 & #2)
1. **Halaman laporan wali**: tambah banner atas "🔐 Tautan khusus untuk Anda · berlaku s/d [tanggal]" + "Terakhir diperbarui: [tanggal, jam]".
2. **Riwayat hafalan**: 10 terakhir + tombol "Lihat semua hafalan" (halaman sama, expand) — menjawab wali murid "datanya statis?".
3. **Tanda "diperbarui"**: timestamp otomatis dari data terakhir (hafalan/absensi terbaru), bukan hardcode.
4. **Kredibilitas paket Yayasan** (celah #2 — produk sisi kecil, sisanya disiplin jual):
   - Halaman Paket: kartu Yayasan diberi label "🚧 Segera hadir — daftar waitlist" + tombol WA "Daftar waitlist Yayasan" (menangkap demand tanpa menjual yang belum ada).
   - MRP/panel admin: paket yayasan tidak bisa diaktifkan selama fitur belum ada (guard di admin_panel: tampilkan peringatan "fitur multi-cabang belum tersedia" saat pilih yayasan).

Verifikasi: banner tampil di laporan wali (screenshot-level check), waitlist WA link bekerja, panel admin menolak aktifkan yayasan tanpa konfirmasi ekstra.

## Batch 4 — Bukti Sosial di Dalam Produk (pendukung, bukan pengganti 3 testimoni asli)
*Perbaikan "cara jual" sebagian besar di luar produk, tapi produk bisa membantu:*

1. Halaman landing: section "Sekolah yang sudah bergabung" — placeholder yang diisi NAMA SEKOLAH ASLI setelah 3 lembaga pertama aktif (jangan testimonial palsu; kosong = tampilkan "Bergabunglah sebagai yang pertama di kota Anda" — jujur dan tetap menarik).
2. Saat sekolah baru selesai onboarding (tur selesai): halaman "Ceritakan ke teman sejawat" dengan tombol salin pesan WA rekomendasi + kode rujukan (sudah ada) — mengubah pengguna pertama jadi saluran.

Verifikasi: landing tampil benar di HP; pesan salin berisi kode rujukan aktif.

---

## URUTAN KERJA & ESTIMASI
Batch 1 → Batch 2 → Batch 3 → Batch 4 (masing-masing: kerjakan → verifikasi → push, sesuai workflow tetap).
Batch 1 dan 2 adalah dua yang paling menentukan konversi & retensi — jangan tukar urutan dengan yang lain.

## YANG TIDAK DIBUAT (disiplin scope — dari laporan)
- Multi-cabang (sampai import + retensi selesai + ada calon pembeli nyata)
- Menjual paket Yayasan sebelum fiturnya ada
- Fitur sekolah umum (PPDB, tabungan santri, absensi guru) — trade-off scope disengaja
- Testimoni palsu / bukti sosial fabricate

## VALIDASI (paralel, bukan setelah kode)
- Setelah Batch 1: tawarkan import ke 2–3 lembaga yang datanya masih di Sheets/Excel → ukur berapa yang berhasil migrasi tanpa bantuan.
- Setelah Batch 2: minta 3 pengasuh aktifkan laporan mingguan → ukur: apakah input harian naik setelah pengasuh "dikawal"? (indikator churn turun)
