package service

import (
	"errors"
	"backend-go/internal/models"
	"backend-go/internal/repo"
)

type BukuService struct {
	repo *repo.BukuRepository
}

func NewBukuService(r *repo.BukuRepository) *BukuService {
	return &BukuService{repo: r}
}

func (s *BukuService) DapatkanSemua() ([]models.Buku, error) {
	return s.repo.AmbilSemua()
}

func (s *BukuService) BuatBuku(buku *models.Buku) error {
	if buku.Judul == "" || buku.Penulis == "" {
		return errors.New("judul dan penulis tidak boleh kosong")
	}
	return s.repo.Simpan(buku)
}

func (s *BukuService) PerbaruiBuku(buku *models.Buku) error {
	return s.repo.Perbarui(buku)
}

func (s *BukuService) HapusBuku(id int) error {
	return s.repo.Hapus(id)
}