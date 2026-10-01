package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"backend-go/global"
	"backend-go/internal/models"
	"backend-go/internal/repo"
)

type BarangService struct {
	repo *repo.BarangRepository
}

func NewBarangService(r *repo.BarangRepository) *BarangService {
	return &BarangService{repo: r}
}

func cacheKeyBarang(userID int) string {
	return fmt.Sprintf("barang:user:%d", userID)
}

func validasiBarang(b *models.Barang) error {
	b.Nama = strings.TrimSpace(b.Nama)
	if b.Nama == "" {
		return errors.New("nama barang tidak boleh kosong")
	}
	if b.Harga <= 0 {
		return errors.New("harga harus lebih dari 0")
	}
	return nil
}

func (s *BarangService) DapatkanSemua(userID int) ([]models.Barang, error) {
	ctx := context.Background()
	key := cacheKeyBarang(userID)

	if global.Redis != nil {
		if cached, err := global.Redis.Get(ctx, key).Result(); err == nil {
			var list []models.Barang
			if json.Unmarshal([]byte(cached), &list) == nil {
				return list, nil
			}
		}
	}

	list, err := s.repo.AmbilSemua(userID)
	if err != nil {
		return nil, err
	}

	if global.Redis != nil {
		if data, err := json.Marshal(list); err == nil {
			_ = global.Redis.Set(ctx, key, data, 10*time.Minute).Err()
		}
	}
	return list, nil
}

func (s *BarangService) Buat(userID int, b *models.Barang) error {
	if err := validasiBarang(b); err != nil {
		return err
	}
	if err := s.repo.Simpan(userID, b); err != nil {
		log.Println("Error membuat barang:", err)
		return err
	}
	s.hapusCache(userID)
	return nil
}

func (s *BarangService) Perbarui(userID int, b *models.Barang) error {
	if err := validasiBarang(b); err != nil {
		return err
	}
	if err := s.repo.Perbarui(userID, b); err != nil {
		return err
	}
	s.hapusCache(userID)
	return nil
}

func (s *BarangService) Hapus(userID, id int) error {
	if err := s.repo.Hapus(userID, id); err != nil {
		return err
	}
	s.hapusCache(userID)
	return nil
}

func (s *BarangService) hapusCache(userID int) {
	if global.Redis != nil {
		_ = global.Redis.Del(context.Background(), cacheKeyBarang(userID)).Err()
	}
}
