package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"backend-go/global"
	"backend-go/internal/models"
	"backend-go/internal/repo"
)

const CacheKeySemuaBuku = "katalog:semua_buku"

type BukuService struct {
	repo *repo.BukuRepository
}

func NewBukuService(r *repo.BukuRepository) *BukuService {
	return &BukuService{repo: r}
}

// 1. Ambil Semua Buku (Dengan Pola Cache-Aside)
func (s *BukuService) DapatkanSemua() ([]models.Buku, error) {
	ctx := context.Background()

	// A. Cek Redis terlebih dahulu jika Redis aktif
	if global.Redis != nil {
		cachedData, err := global.Redis.Get(ctx, CacheKeySemuaBuku).Result()
		if err == nil {
			var bukuList []models.Buku
			if jsonErr := json.Unmarshal([]byte(cachedData), &bukuList); jsonErr == nil {
				log.Println("⚡ [CACHE HIT] Mengambil data dari Redis!")
				return bukuList, nil
			}
		}
	}

	// B. Jika di Redis tidak ada (Cache Miss), ambil dari Neon PostgreSQL
	log.Println("🐢 [CACHE MISS] Mengambil data dari Neon PostgreSQL...")
	bukuList, err := s.repo.AmbilSemua()
	if err != nil {
		return nil, err
	}

	// C. Simpan hasilnya ke Redis dengan masa berlaku (TTL) 10 menit
	if global.Redis != nil && len(bukuList) > 0 {
		jsonData, err := json.Marshal(bukuList)
		if err == nil {
			_ = global.Redis.Set(ctx, CacheKeySemuaBuku, jsonData, 10*time.Minute).Err()
			log.Println("💾 Data berhasil disimpan ke cache Redis (TTL 10 menit)")
		}
	}

	return bukuList, nil
}

// 2. Tambah Buku (Invalidasi Cache)
func (s *BukuService) BuatBuku(b *models.Buku) error {
	if b.Judul == "" || b.Penulis == "" {
		return errors.New("judul dan penulis tidak boleh kosong")
	}

	if err := s.repo.Simpan(b); err != nil {
		return err
	}

	// Hapus cache lama agar data berikutnya fresh
	s.HapusCache()
	return nil
}

// 3. Perbarui Buku (Invalidasi Cache)
func (s *BukuService) PerbaruiBuku(b *models.Buku) error {
	if b.Judul == "" || b.Penulis == "" {
		return errors.New("judul dan penulis tidak boleh kosong")
	}

	if err := s.repo.Perbarui(b); err != nil {
		return err
	}

	s.HapusCache()
	return nil
}

// 4. Hapus Buku (Invalidasi Cache)
func (s *BukuService) HapusBuku(id int) error {
	if err := s.repo.Hapus(id); err != nil {
		return err
	}

	s.HapusCache()
	return nil
}

// Fungsi pembantu untuk membersihkan cache
func (s *BukuService) HapusCache() {
	if global.Redis != nil {
		ctx := context.Background()
		_ = global.Redis.Del(ctx, CacheKeySemuaBuku).Err()
		log.Println("🧹 Cache 'katalog:semua_buku' telah dihapus (di-invalidasi)!")
	}
}