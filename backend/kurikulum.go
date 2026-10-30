package main

import (
	"strings"
)

// ---------- P2 PIVOT: jenjang & kurikulum merdeka ----------

// Jenjang sekolah umum: SD(1-6), SMP(7-9), SMA(10-12).
// Dipilih saat daftar (sekolah umum) atau dipertahankan untuk lembaga islam.
// Mapel default Kurikulum Merdeka per jenjang — dipakai saat sekolah umum mendaftar.
var mapelMerdeka = map[string][]string{
	"SD": {
		"Pendidikan Agama dan Budi Pekerti", "Pendidikan Pancasila",
		"Bahasa Indonesia", "Matematika", "Ilmu Pengetahuan Alam dan Sosial (IPAS)",
		"Seni Budaya (SBDP)", "PJOK", "Bahasa Inggris",
	},
	"SMP": {
		"Pendidikan Agama dan Budi Pekerti", "Pendidikan Pancasila",
		"Bahasa Indonesia", "Bahasa Inggris", "Matematika",
		"Ilmu Pengetahuan Alam (IPA)", "Ilmu Pengetahuan Sosial (IPS)",
		"Informatika", "PJOK", "Seni Budaya",
	},
	"SMA": {
		"Pendidikan Agama dan Budi Pekerti", "Pendidikan Pancasila",
		"Bahasa Indonesia", "Bahasa Inggris", "Matematika",
		"Sejarah", "Ekonomi", "Sosiologi", "Geografi",
		"Fisika", "Kimia", "Biologi", "Informatika", "PJOK", "Seni Budaya",
	},
}

var tingkatJenjang = map[string][]string{
	"SD":  {"1", "2", "3", "4", "5", "6"},
	"SMP": {"7", "8", "9"},
	"SMA": {"10", "11", "12"},
}

// seedKurikulumUmum: buat mapel default + tingkat kelas utk sekolah umum baru.
// Dipanggil dari handleDaftar setelah INSERT sekolah (hanya jenis 'umum').
func seedKurikulumUmum(sekolahID int, jenjang string) {
	mpl, ok := mapelMerdeka[jenjang]
	if !ok {
		jenjang = "SD" // fallback aman
		mpl = mapelMerdeka[jenjang]
	}
	for _, m := range mpl {
		db.Exec(`INSERT INTO mapel (sekolah_id, nama, jenjang) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING`,
			sekolahID, m, jenjang)
	}
	for _, t := range tingkatJenjang[jenjang] {
		db.Exec(`INSERT INTO kelas_tingkat (sekolah_id, tingkat) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
			sekolahID, t)
	}
}

// predikatMerdeka: NA → predikat (Kurikulum Merdeka, interval umum).
func predikatMerdeka(na float64) (string, string) {
	switch {
	case na >= 90:
		return "A", "Sangat Baik"
	case na >= 80:
		return "B", "Baik"
	case na >= 70:
		return "C", "Cukup"
	case na > 0:
		return "D", "Perlu Bimbingan"
	}
	return "-", "-"
}

// romanSemester: "1"→"Ganjil (I)" dst — untuk kop rapor.
func labelSemester(sem string) string {
	s := strings.TrimSpace(sem)
	if s == "2" || strings.EqualFold(s, "genap") {
		return "Semester 2 (Genap)"
	}
	return "Semester 1 (Ganjil)"
}
