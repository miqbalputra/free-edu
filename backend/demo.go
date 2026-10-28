
package main

import (
	"time"
)

// seedDataDemo mengisi data contoh untuk sekolah baru agar user gaptek
// langsung melihat cara kerja aplikasi. Semua data nyata-nya bisa dihapus
// lewat UI (arsip siswa, hapus nilai, dst).
func seedDataDemo(sekolahID, adminID int) {
	thn := time.Now().Format("2006")
	bulanIni := bulanID[time.Now().Month()-1]

	// 6 siswa contoh, 2 kelas
	type sw struct {
		nis, nama, jk, kelas, ortu, hp string
	}
	daftar := []sw{
		{"2026001", "Ahmad Fauzan", "L", "1A", "Bpk. Hasan", "081200000101"},
		{"2026002", "Aisyah Nur Hidayah", "P", "1A", "Bpk. Rahmat", "081200000102"},
		{"2026003", "Muhammad Rizki", "L", "1A", "Bpk. Idris", "081200000103"},
		{"2026004", "Fatimah Zahra", "P", "2A", "Bpk. Slamet", "081200000104"},
		{"2026005", "Umar Abdillah", "L", "2A", "Bpk. Wijaya", "081200000105"},
		{"2026006", "Khadijah Salsabila", "P", "2A", "Bpk. Tono", "081200000106"},
	}
	var ids []int
	for _, s := range daftar {
		var id int
		err := db.QueryRow(`INSERT INTO siswa (nis,nisn,nama,jk,kelas,nama_ortu,hp_ortu,sekolah_id)
			VALUES ($1,'',$2,$3,$4,$5,$6,$7) RETURNING id`,
			s.nis, s.nama, s.jk, s.kelas, s.ortu, s.hp, sekolahID).Scan(&id)
		if err == nil {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return
	}

	// guru contoh (login demo tidak diberikan; cukup sebagai pemilik data)
	var guruID int
	err := db.QueryRow(`INSERT INTO users (nama,email,password,role,sekolah_id)
		VALUES ('Ustadz Contoh (demo)','demo.guru.'||$1||'@demo.local','x-akun-demo-nonaktif-000000','guru',$2) RETURNING id`,
		itoa(sekolahID), sekolahID).Scan(&guruID)
	if err != nil {
		guruID = adminID
	}

	// absensi 3 hari terakhir (Hadir/Izin/Sakit bervariasi)
	hari := []string{"H", "H", "H", "H", "H", "H", "I", "H", "H", "H", "H", "H", "H", "H", "S", "H", "H", "H"}
	tgl := time.Now()
	for d := 0; d < 3; d++ {
		t := tgl.AddDate(0, 0, -d).Format("2006-01-02")
		for i, id := range ids {
			st := map[string]string{"H":"hadir","I":"izin","S":"sakit"}[hari[(d*6+i)%len(hari)]]
			db.Exec(`INSERT INTO absensi_siswa (siswa_id,tanggal,status,sekolah_id) VALUES ($1,$2,$3,$4) ON CONFLICT DO NOTHING`, id, t, st, sekolahID)
		}
	}

	// nilai contoh: 2 mapel, jenis tugas & uh, semester ganjil TA berjalan
	for _, id := range ids {
		for _, m := range []struct{ mapel, jenis string }{{"Tahfidz", "tugas"}, {"Tahfidz", "uh"}, {"Matematika", "tugas"}, {"Matematika", "uh"}} {
			db.Exec(`INSERT INTO nilai (siswa_id,mapel,kelas,jenis,semester,tahun_ajaran,nilai,guru_id,sekolah_id)
				VALUES ($1,$2,'',$3,'ganjil',$4,$5,$6,$7) ON CONFLICT DO NOTHING`,
				id, m.mapel, m.jenis, thn+"/2027", 75+float64(id%20), guruID, sekolahID)
		}
	}

	// hafalan contoh
	for i, id := range ids {
		surah := []string{"Al-Fatihah", "Al-Baqarah", "Ali 'Imran", "An-Nisa", "Al-Ma'idah", "Al-An'am"}[i%6]
		status := []string{"lancar", "proses", "murojaah", "lancar", "baru", "proses"}[i%6]
		db.Exec(`INSERT INTO hafalan (siswa_id,surah,ayat_dari,ayat_sampai,status,penyetor_guru_id,sekolah_id)
			VALUES ($1,$2,1,10,$3,$4,$5)`, id, surah, status, guruID, sekolahID)
	}

	// SPP bulan ini: semua 'belum', lalu 2 pertama dilunasi
	for i, id := range ids {
		var sid int
		err := db.QueryRow(`INSERT INTO spp (siswa_id,bulan,tahun,jumlah,status,jatuh_tempo,sekolah_id)
			VALUES ($1,$2,$3,75000,'belum',make_date(`+thn+`,12,10),$4) RETURNING id`,
			id, bulanIni, thn, sekolahID).Scan(&sid)
		if err == nil && i < 2 {
			db.Exec(`UPDATE spp SET status='lunas', tgl_bayar=CURRENT_DATE, metode='tunai', petugas_id=$1 WHERE id=$2`, adminID, sid)
		}
	}

	// pengeluaran contoh
	db.Exec(`INSERT INTO pengeluaran (tanggal,uraian,kategori,jumlah,pencatat_id,sekolah_id) VALUES (CURRENT_DATE,'Beli spidol & penghapus papan tulis','lainnya',35000,$1,$2)`, adminID, sekolahID)
	db.Exec(`INSERT INTO pengeluaran (tanggal,uraian,kategori,jumlah,pencatat_id,sekolah_id) VALUES (CURRENT_DATE,'Konsumsi rapat guru','lainnya',100000,$1,$2)`, adminID, sekolahID)

	// pengumuman contoh
	db.Exec(`INSERT INTO pengumuman (judul,isi,pembuat_id,sekolah_id) VALUES ('Selamat datang di SIMAKA! 👋','Ini contoh pengumuman. Semua data contoh (siswa, nilai, tagihan) boleh dihapus kapan saja lewat menu masing-masing, atau langsung ganti dengan data asli sekolah Anda.',$1,$2)`, adminID, sekolahID)
}
