package initialize

import (
	"database/sql"
	"log"
	"backend-go/global"
	"os"

	_ "github.com/lib/pq"
)

func InitDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Variabel DATABASE_URL tidak ditemukan pada environment!")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Gagal inisialisasi koneksi database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Gagal terhubung ke Neon PostgreSQL: %v", err)
	}

	// Konfigurasi pool koneksi dari config.yaml
	if global.Config != nil {
		db.SetMaxOpenConns(global.Config.DB.MaxOpenConns)
		db.SetMaxIdleConns(global.Config.DB.MaxIdleConns)
	}

	global.DB = db
	log.Println("Koneksi ke Neon PostgreSQL berhasil dibuat!")

	// Migrasi otomatis sederhana
	migrasi := `
	CREATE TABLE IF NOT EXISTS buku (
		id SERIAL PRIMARY KEY,
		judul VARCHAR(255) NOT NULL,
		penulis VARCHAR(255) NOT NULL
	);
	CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
	);`
	if _, err := db.Exec(migrasi); err != nil {
		log.Fatalf("Migrasi tabel gagal: %v", err)
	}
	
}