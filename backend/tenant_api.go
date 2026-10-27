package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// cekLangganan: sekolah harus aktif/trial yang belum habis.
func cekLangganan(sekolahID int) (bool, string) {
	var status string
	var trialSampai time.Time
	err := db.QueryRow(`SELECT status, trial_sampai FROM sekolah WHERE id=$1`, sekolahID).Scan(&status, &trialSampai)
	if err != nil {
		return false, "Sekolah tidak ditemukan"
	}
	if status == "suspend" {
		return false, "Langganan sekolah Anda ditangguhkan. Hubungi admin."
	}
	if status == "trial" && time.Now().After(trialSampai) {
		return false, "Masa trial 14 hari sudah berakhir. Upgrade ke paket Sekolah untuk melanjutkan."
	}
	return true, ""
}

// ---------- ONBOARDING: daftar sekolah baru ----------

func handleDaftar(w http.ResponseWriter, r *http.Request) {
	var in struct {
		NamaSekolah string `json:"nama_sekolah"`
		Subdomain   string `json:"subdomain"`
		NamaAdmin   string `json:"nama_admin"`
		Email       string `json:"email"`
		Password    string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&in)
	in.Subdomain = strings.ToLower(strings.TrimSpace(in.Subdomain))
	if in.NamaSekolah == "" || in.Subdomain == "" || in.NamaAdmin == "" || in.Email == "" || len(in.Password) < 6 {
		jsonErr(w, 400, "Semua field wajib (password min 6 karakter)")
		return
	}
	for _, ch := range in.Subdomain {
		if !(ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' || ch == '-') {
			jsonErr(w, 400, "Subdomain hanya boleh huruf kecil, angka, dan tanda minus")
			return
		}
	}
	var exists int
	db.QueryRow(`SELECT COUNT(*) FROM sekolah WHERE subdomain=$1`, in.Subdomain).Scan(&exists)
	if exists > 0 {
		jsonErr(w, 400, "Subdomain sudah dipakai, pilih yang lain")
		return
	}
	db.QueryRow(`SELECT COUNT(*) FROM users WHERE email=$1`, strings.ToLower(in.Email)).Scan(&exists)
	if exists > 0 {
		jsonErr(w, 400, "Email sudah terdaftar")
		return
	}
	var sid int
	err := db.QueryRow(`INSERT INTO sekolah (subdomain, nama, status) VALUES ($1,$2,'trial') RETURNING id`, in.Subdomain, in.NamaSekolah).Scan(&sid)
	if err != nil {
		jsonErr(w, 500, err.Error())
		return
	}
	db.Exec(`INSERT INTO langganan (sekolah_id, paket, status, catatan) VALUES ($1,'gratis','aktif','Trial 14 hari')`, sid)
	hash, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	var uid int
	db.QueryRow(`INSERT INTO users (nama,email,password,role,sekolah_id) VALUES ($1,$2,$3,'kepala_sekolah',$4) RETURNING id`,
		in.NamaAdmin, strings.ToLower(in.Email), string(hash), sid).Scan(&uid)
	logActivity(uid, "daftar_sekolah", "Sekolah baru: "+in.NamaSekolah+" (subdomain: "+in.Subdomain+")")
	respond(w, 201, map[string]interface{}{
		"ok": "1", "sekolah_id": sid,
		"pesan": "Pendaftaran berhasil! Trial 14 hari aktif. Login di halaman masuk.",
	})
}

// ---------- STATUS LANGGANAN (untuk banner di UI) ----------

func handleLangganan(w http.ResponseWriter, r *http.Request) {
	c := userFromContext(r)
	var status string
	var trialSampai time.Time
	db.QueryRow(`SELECT status, trial_sampai FROM sekolah WHERE id=$1`, c.SekolahID).Scan(&status, &trialSampai)
	sisa := int(time.Until(trialSampai).Hours() / 24)
	if sisa < 0 {
		sisa = 0
	}
	respond(w, 200, map[string]interface{}{"status": status, "trial_sisa_hari": sisa})
}

// ---------- NOTIFIKASI WA (link wa.me siap kirim) ----------

func handleNotifSPP(w http.ResponseWriter, r *http.Request) {
	c := userFromContext(r)
	bulan := r.URL.Query().Get("bulan")
	tahun := r.URL.Query().Get("tahun")
	if bulan == "" || tahun == "" {
		jsonErr(w, 400, "Bulan dan tahun wajib")
		return
	}
	var namaSekolah string
	db.QueryRow(`SELECT nama FROM sekolah WHERE id=$1`, c.SekolahID).Scan(&namaSekolah)
	rows, err := db.Query(`SELECT s.nama, COALESCE(s.hp_ortu,''), sp.jumlah
		FROM spp sp JOIN siswa s ON s.id=sp.siswa_id
		WHERE sp.bulan=$1 AND sp.tahun=$2 AND sp.status='belum' AND s.sekolah_id=$3 AND s.hp_ortu <> ''`,
		bulan, tahun, c.SekolahID)
	if err != nil {
		jsonErr(w, 500, err.Error())
		return
	}
	defer rows.Close()
	list := []map[string]interface{}{}
	for rows.Next() {
		var nama, hp string
		var jumlah float64
		rows.Scan(&nama, &hp, &jumlah)
		hp = "628" + strings.TrimPrefix(strings.TrimPrefix(hp, "08"), "62")
		hp = strings.ReplaceAll(hp, "-", "")
		msg := "Assalamualaikum Bapak/Ibu, kami informasikan bahwa SPP ananda " + nama +
			" bulan " + bulan + " " + tahun + " sebesar Rp" + itoa(int(jumlah)) +
			" di " + namaSekolah + " belum dibayarkan. Mohon pembayarannya. Terima kasih. 🙏"
		list = append(list, map[string]interface{}{"nama": nama, "hp": hp, "link_wa": "https://wa.me/" + hp + "?text=" + waEncode(msg)})
	}
	respond(w, 200, map[string]interface{}{"jumlah": len(list), "pesan": list})
}

func waEncode(s string) string {
	replacer := strings.NewReplacer(" ", "%20", "\n", "%0A", ",", "%2C", "'", "%27")
	return replacer.Replace(s)
}

// ---------- PWA MANIFEST ----------

func handleManifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{
  "name": "SIMAKA — Manajemen Sekolah",
  "short_name": "SIMAKA",
  "start_url": "/app",
  "display": "standalone",
  "background_color": "#f4f6fb",
  "theme_color": "#4f46e5",
  "icons": [
  {"src": "/icons/icon-192.png", "sizes": "192x192", "type": "image/png"},
  {"src": "/icons/icon-512.png", "sizes": "512x512", "type": "image/png"}
  ]
  }`))
}
