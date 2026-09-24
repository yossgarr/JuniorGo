package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"backend-go/global"
	"backend-go/internal/initialize"
	"backend-go/internal/routers"
	"backend-go/pkg/setting"
)

func main() {
	// 1. Muat Environment File
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Info: Membaca environment bawaan sistem...")
	}

	// 2. Muat file YAML
	cfg, err := setting.MuatConfig("config/config.yaml")
	if err != nil {
		log.Printf("Peringatan: config.yaml gagal dimuat (%v), menggunakan default", err)
	} else {
		global.Config = cfg
		log.Printf("Aplikasi: %s v%s siap dijalankan", cfg.App.Name, cfg.App.Version)
	}

	// 3. Inisialisasi Database Pool
	initialize.InitDatabase()
	initialize.InitRedis()
	defer global.DB.Close()

	// 4. Inisialisasi Router & Server
	router := routers.InitRouter()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	fmt.Printf("HTTP Server aktif di port %s (http://localhost:%s/buku)\n", port, port)
	log.Fatal(server.ListenAndServe())
}