package initialize

import (
	"context"
	"log"
	"os"
	"time"

	"backend-go/global"

	"github.com/redis/go-redis/v9"
)

func InitRedis() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Println("⚠️ REDIS_URL kosong di .env, caching dinonaktifkan")
		return
	}

	// Otomatis membaca host, password, dan enkripsi SSL/TLS
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("⚠️ Format REDIS_URL salah: %v", err)
		return
	}

	rdb := redis.NewClient(opt)

	// Uji koneksi Ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("⚠️ Gagal konek ke Upstash Redis: %v", err)
		return
	}

	global.Redis = rdb
	log.Println("✅ Berhasil terhubung ke Upstash Redis Cloud!")
}