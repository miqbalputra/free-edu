package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// ---------- PANEL ADMIN PLATFORM (khusus owner) ----------
// Role 'platform' di users: tidak terikat sekolah, kelola semua sekolah & langganan.

func handleAdminSekolah(w http.ResponseWriter, r *http.Request, c *Claims) {
	rows, err := db.Query(`SELECT s.id, s.nama, s.subdomain, s.status, COALESCE(s.trial_sampai::text,''),
			COALESCE(l.paket,'gratis'), l.status,
			(SELECT COUNT(*) FROM siswa si WHERE si.sekolah_id=s.id AND si.status='aktif'),
			(SELECT COUNT(*) FROM users u WHERE u.sekolah_id=s.id),
			COALESCE((SELECT SUM(sp.jumlah) FROM spp sp JOIN siswa si2 ON si2.id=sp.siswa_id
				WHERE si2.sekolah_id=s.id AND sp.status='lunas'),0)
		FROM sekolah s
		LEFT JOIN langganan l ON l.sekolah_id=s.id AND l.status='aktif'
		ORDER BY s.id`)
	if err != nil {
		jsonErr(w, 500, userErr(err))
		return
	}
	defer rows.Close()
	list := []map[string]interface{}{}
	for rows.Next() {
		var id, nSiswa, nUser int
		var nama, subdomain, status, paket string
		var trial, lStatus string
		var total float64
		rows.Scan(&id, &nama, &subdomain, &status, &trial, &paket, &lStatus, &nSiswa, &nUser, &total)
		list = append(list, map[string]interface{}{
			"id": id, "nama": nama, "subdomain": subdomain, "status": status,
			"trial_sampai": trial, "plan": paket, "langganan_status": lStatus,
			"siswa": nSiswa, "pengguna": nUser, "spp_lunas": total,
		})
	}
	respond(w, 200, list)
}

func handleAdminSetPlan(w http.ResponseWriter, r *http.Request, c *Claims) {
	var in struct {
		SekolahID int    `json:"sekolah_id"`
		Paket     string `json:"paket"`
		Bulan     int    `json:"bulan"` // durasi aktif; 0 = 30 hari utk sekolah, 365 yayasan
	}
	json.NewDecoder(r.Body).Decode(&in)
	if in.Paket != "gratis" && in.Paket != "sekolah" && in.Paket != "yayasan" {
		jsonErr(w, 400, "Paket tidak valid")
		return
	}
	if in.Bulan <= 0 {
		if in.Paket == "yayasan" {
			in.Bulan = 12
		} else {
			in.Bulan = 1
		}
	}
	sampai := time.Now().AddDate(0, in.Bulan, 0)
	// nonaktifkan langganan lama, buat baru
	db.Exec(`UPDATE langganan SET status='berakhir' WHERE sekolah_id=$1 AND status='aktif'`, in.SekolahID)
	_, err := db.Exec(`INSERT INTO langganan (sekolah_id, paket, status, sampai, catatan) VALUES ($1,$2,'aktif',$3,$4)`,
		in.SekolahID, in.Paket, sampai, "Diaktifkan via panel admin")
	if err != nil {
		jsonErr(w, 500, userErr(err))
		return
	}
	// kalau di-upgrade dari gratis/trial, sekolah jadi aktif penuh
	db.Exec(`UPDATE sekolah SET status='aktif' WHERE id=$1 AND status='trial'`, in.SekolahID)
	logActivity(c.UserID, "admin_set_plan", "Sekolah "+strconv.Itoa(in.SekolahID)+" → "+in.Paket+" ("+itoa(in.Bulan)+" bln)")
	respond(w, 200, map[string]interface{}{"ok": "1", "sampai": sampai.Format("2006-01-02")})
}

func handleAdminStats(w http.ResponseWriter, r *http.Request, c *Claims) {
	var nSekolah, nSiswa, nUser, nAktif int
	db.QueryRow(`SELECT COUNT(*) FROM sekolah`).Scan(&nSekolah)
	db.QueryRow(`SELECT COUNT(*) FROM sekolah WHERE status='trial'`).Scan(&nAktif)
	db.QueryRow(`SELECT COUNT(*) FROM siswa WHERE status='aktif'`).Scan(&nSiswa)
	db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&nUser)
	var mrr float64
	db.QueryRow(`SELECT COALESCE(SUM(CASE l.paket WHEN 'sekolah' THEN 45000 WHEN 'yayasan' THEN 200000 ELSE 0 END),0)
		FROM langganan l WHERE l.status='aktif'`).Scan(&mrr)
	respond(w, 200, map[string]interface{}{
		"sekolah": nSekolah, "trial": nAktif, "siswa": nSiswa, "pengguna": nUser, "mrr": mrr,
	})
}
