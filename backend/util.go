package main

import (
	"strings"
	"time"
)

func today() string {
	return time.Now().Format("2006-01-02")
}

// userErr: terjemahkan error database menjadi pesan manusiawi (bahasa Indonesia)
func userErr(err error) string {
	if err == nil {
		return "Terjadi kesalahan. Coba lagi."
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, "duplicate key"):
		if strings.Contains(msg, "users_email") {
			return "Email ini sudah terdaftar. Gunakan email lain atau masuk dengan email tersebut."
		}
		if strings.Contains(msg, "siswa") {
			return "Data siswa ini sudah ada (kemungkinan NIS/nama sama). Periksa daftar siswa dulu."
		}
		return "Data ini sudah ada sebelumnya. Periksa dulu daftar yang sudah tersimpan."
	case strings.Contains(msg, "invalid input syntax"):
		return "Format tanggal/angka tidak sesuai. Periksa kembali isian Anda."
	case strings.Contains(msg, "foreign key"):
		return "Data ini masih dipakai di bagian lain. Hapus data terkait dulu (mis. nilai atau tagihan siswa)."
	case strings.Contains(msg, "violates check constraint"):
		return "Isian tidak sesuai ketentuan. Periksa kembali nilai yang dimasukkan."
	case strings.Contains(msg, "connection refused"):
		return "Koneksi database terputus. Coba beberapa saat lagi."
	case strings.Contains(msg, "does not exist"):
		return "Data tidak ditemukan. Muat ulang halaman dan coba lagi."
	default:
		return "Terjadi kesalahan saat menyimpan. Coba lagi; bila berulang, hubungi admin."
	}
}
