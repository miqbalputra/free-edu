package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

// loadOrCreateSecret: pakai SIMAKA_JWT_SECRET dari env; kalau kosong, generate sekali
// dan simpan di file agar tetap sama setelah restart.
func loadOrCreateSecret() ([]byte, error) {
	if s := os.Getenv("SIMAKA_JWT_SECRET"); s != "" {
		return []byte(s), nil
	}
	path := "/etc/simaka/jwt_secret"
	if b, err := os.ReadFile(path); err == nil && len(b) >= 32 {
		return b, nil
	}
	buf := make([]byte, 48)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}
	sec := []byte(base64.StdEncoding.EncodeToString(buf))
	if err := os.MkdirAll("/etc/simaka", 0700); err != nil {
		// fallback: simpan di working dir
		os.WriteFile("jwt_secret", sec, 0600)
		return sec, nil
	}
	if err := os.WriteFile(path, sec, 0600); err != nil {
		return nil, err
	}
	fmt.Println("JWT secret baru dibuat:", path)
	return sec, nil
}
