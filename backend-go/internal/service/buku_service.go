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

func (s *BukuService) DapatkanSemua() ([]models.Buku, error) {
	ctx := context.Background()

	// A. Cek Cache Redis
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

	// B. Cache Miss -> Ambil dari PostgreSQL
	log.Println("🐢 [CACHE MISS] Mengambil data dari Neon PostgreSQL...")
	bukuList, err := s.repo.AmbilSemua()
	if err != nil {
		log.Println("❌ Error dari database saat AmbilSemua():", err)
		return nil, err
	}

	// C. Simpan ke Cache Redis (TTL 10 Menit)
	if global.Redis != nil && len(bukuList) > 0 {
		jsonData, err := json.Marshal(bukuList)
		if err == nil {
			_ = global.Redis.Set(ctx, CacheKeySemuaBuku, jsonData, 10*time.Minute).Err()
			log.Println("💾 Data berhasil disimpan ke cache Redis (TTL 10 menit)")
		}
	}

	return bukuList, nil
}

func (s *BukuService) BuatBuku(b *models.Buku) error {
	if b.Judul == "" || b.Penulis == "" {
		return errors.New("judul dan penulis tidak boleh kosong")
	}
	if b.Stok < 0 {
		b.Stok = 5 // Default stok jika tidak diisi
	}

	if err := s.repo.Simpan(b); err != nil {
		log.Println("❌ Error saat membuat buku:", err)
		return err
	}

	s.HapusCache()
	return nil
}

func (s *BukuService) PerbaruiBuku(b *models.Buku) error {
	if b.Judul == "" || b.Penulis == "" {
		return errors.New("judul dan penulis tidak boleh kosong")
	}

	if err := s.repo.Perbarui(b); err != nil {
		log.Println("❌ Error saat memperbarui buku:", err)
		return err
	}

	s.HapusCache()
	return nil
}

func (s *BukuService) HapusBuku(id int) error {
	if err := s.repo.Hapus(id); err != nil {
		log.Println("❌ Error saat menghapus buku:", err)
		return err
	}

	s.HapusCache()
	return nil
}

// Layanan Peminjaman Buku (Kurangi Stok & Invalidasi Cache)
func (s *BukuService) PinjamBuku(userID int, bukuID int) error {
	if userID <= 0 || bukuID <= 0 {
		return errors.New("user ID dan buku ID wajib diisi")
	}

	err := s.repo.PinjamBukuTx(userID, bukuID)
	if err != nil {
		log.Println("❌ Transaksi peminjaman gagal:", err)
		return err
	}

	// Hapus cache agar stok terbaru langsung muncul di katalog
	s.HapusCache()
	return nil
}

func (s *BukuService) HapusCache() {
	if global.Redis != nil {
		ctx := context.Background()
		_ = global.Redis.Del(ctx, CacheKeySemuaBuku).Err()
		log.Println("🧹 Cache 'katalog:semua_buku' berhasil di-invalidasi!")
	}
}