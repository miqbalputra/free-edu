package main

import (
	"net/http"
	"sync"
	"time"
)

type attempt struct {
	failures int
	lockedUntil time.Time
}

var (
	loginMu    sync.Mutex
	loginTries = map[string]*attempt{}
)

// loginRateLimit: max 5 gagal per IP+email, lock 15 menit.
// return "" jika boleh, atau pesan lock.
func loginRateLimit(key string) string {
	loginMu.Lock()
	defer loginMu.Unlock()
	a, ok := loginTries[key]
	if !ok {
		return ""
	}
	if time.Now().Before(a.lockedUntil) {
		sisa := int(time.Until(a.lockedUntil).Minutes()) + 1
		return "Terlalu banyak percobaan gagal. Coba lagi dalam " + itoa(sisa) + " menit."
	}
	return ""
}

func loginRecordFail(key string) {
	loginMu.Lock()
	defer loginMu.Unlock()
	a, ok := loginTries[key]
	if !ok {
		a = &attempt{}
		loginTries[key] = a
	}
	a.failures++
	if a.failures >= 5 {
		a.lockedUntil = time.Now().Add(15 * time.Minute)
		a.failures = 0
	}
}

func loginRecordSuccess(key string) {
	loginMu.Lock()
	defer loginMu.Unlock()
	delete(loginTries, key)
}

// activityLog: catat aksi user (audit trail sederhana)
func logActivity(userID int, aksi, detail string) {
	db.Exec(`INSERT INTO activity_log (user_id, aksi, detail) VALUES ($1,$2,$3)`, userID, aksi, detail)
}

var _ = http.MethodGet
