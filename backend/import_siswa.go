package main

// BATCH 1 — Import Data Massal (rencana perbaikan pasar)
// Import siswa dari CSV: preview → pemetaan kolom → import dengan validasi manusiawi.
// Mendukung .csv (koma/titik koma). Excel: user disuruh "Simpan Sebagai CSV" dengan petunjuk jelas.

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// parseCSV: baca semua record, deteksi delimiter otomatis (koma vs titik koma)
func parseCSV(data []byte) ([][]string, error) {
	// BOM
	s := strings.TrimPrefix(string(data), "\ufeff")
	// deteksi delimiter dari baris pertama
	delim := ","
	first := strings.SplitN(s, "\n", 2)[0]
	if strings.Count(first, ";") > strings.Count(first, ",") {
		delim = ";"
	}
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = rune(delim[0])
	r.FieldsPerRecord = -1 // biarkan fleksibel; validasi manual
	var out [][]string
	for {
		fields, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		for i := range fields {
			fields[i] = strings.TrimSpace(fields[i])
		}
		out = append(out, fields)
	}
	return out, nil
}


// guessKolom: petakan indeks kolom dari header — menerima variasi nama umum
func guessKolom(header []string) map[string]int {
	m := map[string]int{"nis": -1, "nama": -1, "kelas": -1, "jk": -1, "nama_ortu": -1, "hp_ortu": -1}
	for i, h := range header {
		h = strings.ToLower(strings.TrimSpace(h))
		switch {
		case m["nama"] < 0 && (h == "nama" || h == "nama siswa" || h == "santri" || strings.Contains(h, "nama") && !strings.Contains(h, "ortu") && !strings.Contains(h, "wali")):
			m["nama"] = i
		case m["nis"] < 0 && (h == "nis" || h == "no. induk" || h == "no induk" || h == "nis/nisn"):
			m["nis"] = i
		case m["kelas"] < 0 && strings.Contains(h, "kelas"):
			m["kelas"] = i
		case m["jk"] < 0 && (h == "jk" || h == "jenis kelamin" || h == "l/p" || strings.Contains(h, "kelamin")):
			m["jk"] = i
		case m["nama_ortu"] < 0 && (strings.Contains(h, "ortu") || strings.Contains(h, "wali")) && !strings.Contains(h, "hp") && !strings.Contains(h, "telp"):
			m["nama_ortu"] = i
		case m["hp_ortu"] < 0 && (strings.Contains(h, "hp") || strings.Contains(h, "telp") || strings.Contains(h, "telepon") || strings.Contains(h, "wa")):
			m["hp_ortu"] = i
		}
	}
	return m
}

// normalisasiKelas: "1 a" → "1A", "kela 2b" → "2B"
func normalisasiKelas(k string) (string, string) {
	orig := k
	k = strings.ToUpper(strings.Join(strings.Fields(k), ""))
	k = strings.TrimPrefix(k, "KELAS")
	if k != orig {
		return k, "diperbaiki dari \"" + orig + "\""
	}
	if k != orig && k == strings.Join(strings.Fields(orig), " ") {
		return k, "diperbaiki (spasi dihapus)"
	}
	if strings.Join(strings.Fields(orig), "") != orig {
		return k, "diperbaiki (spasi dihapus)"
	}
	return k, ""
}

func normJK(v string) string {
	v = strings.ToUpper(strings.TrimSpace(v))
	if strings.HasPrefix(v, "L") {
		return "L"
	}
	if strings.HasPrefix(v, "P") {
		return "P"
	}
	return ""
}

// handleImportPreview: upload CSV → kembalikan header + 8 baris pertama + pemetaan tebakan
func handleImportPreview(w http.ResponseWriter, r *http.Request) {
	c := userFromContext(r)
	if c == nil {
		jsonErr(w, 401, "Unauthorized")
		return
	}
	if err := r.ParseMultipartForm(5 << 20); err != nil { // 5 MB
		jsonErr(w, 400, "File terlalu besar (maksimal 5 MB).")
		return
	}
	f, _, err := r.FormFile("file")
	if err != nil {
		jsonErr(w, 400, "File tidak terbaca. Pilih file CSV dulu ya.")
		return
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		jsonErr(w, 400, "File tidak terbaca.")
		return
	}
	rows, err := parseCSV(data)
	if err != nil || len(rows) < 2 {
		jsonErr(w, 400, "File CSV tidak terbaca atau kosong. Pastikan baris pertama berisi nama kolom (mis. NIS, Nama, Kelas).")
		return
	}
	header := rows[0]
	kol := guessKolom(header)
	preview := rows[1:]
	if len(preview) > 8 {
		preview = preview[:8]
	}
	total := len(rows) - 1
	respond(w, 200, map[string]interface{}{
		"header": header, "preview": preview, "mapping": kol, "total": total,
		"nama_ok": kol["nama"] >= 0, "kelas_ok": kol["kelas"] >= 0,
	})
}

// handleImportSiswa: import sungguhan. Body: JSON {rows: [[...]], mapping: {...}}
func handleImportSiswa(w http.ResponseWriter, r *http.Request, c *Claims) {
	var in struct {
		Header  []string       `json:"header"`
		Rows    [][]string     `json:"rows"`
		Mapping map[string]int `json:"mapping"`
	}
	json.NewDecoder(r.Body).Decode(&in)
	if in.Mapping == nil {
		in.Mapping = guessKolom(in.Header)
	}
	if in.Mapping["nama"] < 0 || in.Mapping["kelas"] < 0 {
		jsonErr(w, 400, "Kolom Nama dan Kelas wajib dipetakan. Coba periksa kembali pemetaan kolom.")
		return
	}
	lim := limitsFor(planOf(c.SekolahID))
	var terpakai int
	db.QueryRow(`SELECT COUNT(*) FROM siswa WHERE sekolah_id=$1 AND status='aktif'`, c.SekolahID).Scan(&terpakai)
	var masuk, dilewati, diperbaiki int
	var catatan []string
	for i, row := range in.Rows {
		if len(row) == 0 {
			continue
		}
		get := func(key string) string {
			idx := in.Mapping[key]
			if idx < 0 || idx >= len(row) {
				return ""
			}
			return strings.TrimSpace(row[idx])
		}
		nama := get("nama")
		if nama == "" {
			dilewati++
			catatan = append(catatan, fmt.Sprintf("Baris %d: nama kosong — dilewati.", i+2))
			continue
		}
		kelasRaw := get("kelas")
		if kelasRaw == "" {
			dilewati++
			catatan = append(catatan, fmt.Sprintf("Baris %d (%s): kelas kosong — dilewati.", i+2, nama))
			continue
		}
		kelas, fix := normalisasiKelas(kelasRaw)
		if fix != "" {
			diperbaiki++
			catatan = append(catatan, fmt.Sprintf("Baris %d: kelas %s → %s.", i+2, kelasRaw, kelas))
		}
		nis := get("nis")
		jk := normJK(get("jk"))
		hp := get("hp_ortu")
		// anti-duplikat: NIS sudah ada (di sekolah ini) → perbarui nama/kelas/hp
		var existing int
		if nis != "" {
			db.QueryRow(`SELECT id FROM siswa WHERE nis=$1 AND sekolah_id=$2`, nis, c.SekolahID).Scan(&existing)
		}
		if existing > 0 {
			db.Exec(`UPDATE siswa SET nama=$1, kelas=$2, hp_ortu=NULLIF($3,''), nama_ortu=NULLIF($4,'') WHERE id=$5`,
				nama, kelas, hp, get("nama_ortu"), existing)
			masuk++
			continue
		}
		// kuota: siswa aktif baru melewati batas plan → berhenti dengan pesan ramah
		if lim.MaxSiswa > 0 && terpakai+masuk >= lim.MaxSiswa {
			dilewati++
			catatan = append(catatan, fmt.Sprintf("Kuota paket tercapai (maks %d siswa aktif). Sisanya tidak diimpor — upgrade paket untuk melanjutkan.", lim.MaxSiswa))
			break
		}
		_, err := db.Exec(`INSERT INTO siswa (nis,nama,jk,kelas,nama_ortu,hp_ortu,sekolah_id) VALUES (NULLIF($1,''),$2,NULLIF($3,''),$4,NULLIF($5,''),NULLIF($6,''),$7)`,
			nis, nama, jk, kelas, get("nama_ortu"), hp, c.SekolahID)
		if err != nil {
			dilewati++
			catatan = append(catatan, fmt.Sprintf("Baris %d (%s): gagal masuk — %s", i+2, nama, userErr(err)))
			continue
		}
		masuk++
	}
	logActivity(c.UserID, "import_siswa", fmt.Sprintf("Import %d siswa (%d dilewati, %d diperbaiki)", masuk, dilewati, diperbaiki))
	respond(w, 200, map[string]interface{}{
		"masuk": masuk, "dilewati": dilewati, "diperbaiki": diperbaiki, "catatan": catatan,
	})
}


// handleTemplateImport: CSV contoh siap-unduh
func handleTemplateImport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"contoh-import-siswa.csv\"")
	w.Write([]byte("NIS,Nama,Kelas,Jenis Kelamin,Nama Ortu,HP Ortu\n101,Ahmad Fauzi,1A,L,Budi Santoso,081234567890\n102,Siti Aminah,1A,P,Siti Rahayu,081234567891\n103,Muhammad Rizki,2B,L,Andi Wijaya,081234567892\n"))
}
