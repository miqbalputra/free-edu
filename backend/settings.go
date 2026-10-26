package main

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// ---------- PENGATURAN SEKOLAH ----------

func handleGetPengaturan(w http.ResponseWriter, r *http.Request) {
	var nama, alamat, npsn, kepala, ta, sem string
	err := db.QueryRow(`SELECT nama_sekolah, COALESCE(alamat,''), COALESCE(npsn,''), COALESCE(kepala_sekolah,''), tahun_ajaran_aktif, semester_aktif FROM pengaturan WHERE id=1`).
		Scan(&nama, &alamat, &npsn, &kepala, &ta, &sem)
	if err != nil {
		jsonErr(w, 500, err.Error())
		return
	}
	respond(w, 200, map[string]string{
		"nama_sekolah": nama, "alamat": alamat, "npsn": npsn, "kepala_sekolah": kepala,
		"tahun_ajaran_aktif": ta, "semester_aktif": sem,
	})
}

func handleUpdatePengaturan(w http.ResponseWriter, r *http.Request, c *Claims) {
	var in map[string]string
	json.NewDecoder(r.Body).Decode(&in)
	_, err := db.Exec(`UPDATE pengaturan SET nama_sekolah=$1, alamat=$2, npsn=$3, kepala_sekolah=$4, tahun_ajaran_aktif=$5, semester_aktif=$6 WHERE id=1`,
		in["nama_sekolah"], in["alamat"], in["npsn"], in["kepala_sekolah"], in["tahun_ajaran_aktif"], in["semester_aktif"])
	if err != nil {
		jsonErr(w, 500, err.Error())
		return
	}
	logActivity(c.UserID, "update_pengaturan", "Memperbarui pengaturan sekolah")
	respond(w, 200, map[string]string{"ok": "1"})
}

// ---------- RESET PASSWORD ----------

func randPassword() string {
	const chars = "abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789"
	n := big.NewInt(int64(len(chars)))
	b := make([]byte, 10)
	for i := range b {
		x, _ := rand.Int(rand.Reader, n)
		b[i] = chars[x.Int64()]
	}
	return string(b)
}

func handleResetPassword(w http.ResponseWriter, r *http.Request, c *Claims) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	// hanya boleh reset user dari sekolah yang sama; tidak boleh reset sesama kepala sekolah lain
	var sid int
	var role string
	if err := db.QueryRow(`SELECT COALESCE(sekolah_id,0), role FROM users WHERE id=$1`, id).Scan(&sid, &role); err != nil || sid != c.SekolahID {
		jsonErr(w, 404, "Pengguna tidak ditemukan")
		return
	}
	if role == "kepala_sekolah" && id != c.UserID {
		jsonErr(w, 403, "Tidak boleh me-reset kepala sekolah lain")
		return
	}
	newPass := randPassword()
	hash, _ := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	res, err := db.Exec(`UPDATE users SET password=$1 WHERE id=$2`, string(hash), id)
	if err != nil {
		jsonErr(w, 500, err.Error())
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		jsonErr(w, 404, "Pengguna tidak ditemukan")
		return
	}
	var nama string
	db.QueryRow(`SELECT nama FROM users WHERE id=$1`, id).Scan(&nama)
	logActivity(c.UserID, "reset_password", "Reset password untuk user ID "+itoa(id)+" ("+nama+")")
	respond(w, 200, map[string]string{"password_baru": newPass})
}

func handleUbahPassword(w http.ResponseWriter, r *http.Request, c *Claims) {
	var in struct {
		Lama, Baru string
	}
	json.NewDecoder(r.Body).Decode(&in)
	if len(in.Baru) < 6 {
		jsonErr(w, 400, "Password baru minimal 6 karakter")
		return
	}
	if c.UserID == 0 {
		jsonErr(w, 401, "Tidak terautentikasi")
		return
	}
	var hash string
	if err := db.QueryRow(`SELECT password FROM users WHERE id=$1`, c.UserID).Scan(&hash); err != nil {
		jsonErr(w, 500, err.Error())
		return
	}
	if !checkPass(hash, in.Lama) {
		jsonErr(w, 400, "Password lama salah")
		return
	}
	newHash, _ := bcrypt.GenerateFromPassword([]byte(in.Baru), bcrypt.DefaultCost)
	db.Exec(`UPDATE users SET password=$1 WHERE id=$2`, string(newHash), c.UserID)
	logActivity(c.UserID, "ubah_password", "Mengubah password sendiri")
	respond(w, 200, map[string]string{"ok": "1"})
}

// ---------- CATATAN RAPOR ----------

func getCatatan(sid int, ta string) (cw, ek, pr string) {
	db.QueryRow(`SELECT COALESCE(catatan_wali,''), COALESCE(ekstrakurikuler,''), COALESCE(prestasi,'')
		FROM rapor_catatan WHERE siswa_id=$1 AND tahun_ajaran=$2`, sid, ta).Scan(&cw, &ek, &pr)
	return
}

func handleGetCatatanRapor(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.Atoi(r.PathValue("id"))
	ta := r.URL.Query().Get("tahun_ajaran")
	if ta == "" {
		ta = "2026/2027"
	}
	cw, ek, pr := getCatatan(sid, ta)
	respond(w, 200, map[string]string{"catatan_wali": cw, "ekstrakurikuler": ek, "prestasi": pr})
}

func handleSaveCatatanRapor(w http.ResponseWriter, r *http.Request, c *Claims) {
	sid, _ := strconv.Atoi(r.PathValue("id"))
	var in struct {
		TahunAjaran    string `json:"tahun_ajaran"`
		CatatanWali    string `json:"catatan_wali"`
		Ekstrakurikuler string `json:"ekstrakurikuler"`
		Prestasi       string `json:"prestasi"`
	}
	json.NewDecoder(r.Body).Decode(&in)
	if in.TahunAjaran == "" {
		in.TahunAjaran = "2026/2027"
	}
	_, err := db.Exec(`INSERT INTO rapor_catatan (siswa_id, tahun_ajaran, catatan_wali, ekstrakurikuler, prestasi, wali_id)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (siswa_id, tahun_ajaran) DO UPDATE SET catatan_wali=EXCLUDED.catatan_wali, ekstrakurikuler=EXCLUDED.ekstrakurikuler, prestasi=EXCLUDED.prestasi, wali_id=EXCLUDED.wali_id`,
		sid, in.TahunAjaran, in.CatatanWali, in.Ekstrakurikuler, in.Prestasi, c.UserID)
	if err != nil {
		jsonErr(w, 500, err.Error())
		return
	}
	logActivity(c.UserID, "update_rapor", "Menyimpan catatan rapor siswa ID "+itoa(sid))
	respond(w, 200, map[string]string{"ok": "1"})
}

// ---------- ACTIVITY LOG ----------

func handleListLog(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT l.id, COALESCE(u.nama,'-'), l.aksi, COALESCE(l.detail,''), l.waktu::text
		FROM activity_log l LEFT JOIN users u ON u.id=l.user_id ORDER BY l.waktu DESC, l.id DESC LIMIT 200`)
	if err != nil {
		jsonErr(w, 500, err.Error())
		return
	}
	defer rows.Close()
	list := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var nama, aksi, detail, waktu string
		rows.Scan(&id, &nama, &aksi, &detail, &waktu)
		list = append(list, map[string]interface{}{"id": id, "user": nama, "aksi": aksi, "detail": detail, "waktu": waktu})
	}
	respond(w, 200, list)
}
