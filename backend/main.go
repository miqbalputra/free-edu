package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://simaka:simaka123@localhost:5432/simaka?sslmode=disable"
	}
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(20)
	if err := db.Ping(); err != nil {
		log.Fatal("DB not reachable: ", err)
	}
	sec, err := loadOrCreateSecret()
	if err != nil {
		log.Fatal("jwt secret: ", err)
	}
	jwtSecret = sec
	if err := migrate(); err != nil {
		log.Fatal("migrate: ", err)
	}
	if err := migrateTenants(); err != nil {
		log.Fatal("migrate tenants: ", err)
	}
	if err := seed(); err != nil {
		log.Fatal("seed: ", err)
	}

	mux := http.NewServeMux()
	registerRoutes(mux)
	mux.HandleFunc("GET /assets/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/"+strings.TrimPrefix(r.URL.Path, "/"))
	}))
	// Landing page promosi di /, aplikasi di /app/
	mux.Handle("GET /{$}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/landing.html")
	}))
	mux.Handle("GET /app", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/index.html")
	}))
	mux.Handle("GET /app/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/index.html")
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      cors(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	fmt.Println("SIMAKA API running on :" + port)
	log.Fatal(srv.ListenAndServe())
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
