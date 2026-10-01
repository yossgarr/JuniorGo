package service

import (
	"errors"
	"strings"

	"backend-go/internal/models"
	"backend-go/internal/repo"
)

const batasRiwayat = 200

type TransaksiService struct {
	repo       *repo.TransaksiRepository
	barangRepo *repo.BarangRepository
}

func NewTransaksiService(r *repo.TransaksiRepository, br *repo.BarangRepository) *TransaksiService {
	return &TransaksiService{repo: r, barangRepo: br}
}

// Catat membuat transaksi baru. Untuk pengeluaran dengan barang_id, nominal
// dihitung di server (harga x qty) sehingga klien tidak bisa memalsukannya.
func (s *TransaksiService) Catat(userID int, req models.TransaksiRequest) (*models.Transaksi, error) {
	if req.Tipe != models.TipePemasukan && req.Tipe != models.TipePengeluaran {
		return nil, errors.New("tipe harus 'pemasukan' atau 'pengeluaran'")
	}

	t := &models.Transaksi{
		Tipe:       req.Tipe,
		Keterangan: strings.TrimSpace(req.Keterangan),
		Qty:        1,
	}

	if req.BarangID > 0 {
		if req.Tipe != models.TipePengeluaran {
			return nil, errors.New("barang hanya bisa dipakai untuk pengeluaran")
		}
		qty := req.Qty
		if qty <= 0 {
			qty = 1
		}
		barang, err := s.barangRepo.AmbilByID(userID, req.BarangID)
		if err != nil {
			return nil, err
		}
		t.BarangID = &barang.ID
		t.NamaBarang = barang.Nama
		t.Qty = qty
		t.Jumlah = barang.Harga * int64(qty)
		if t.Keterangan == "" {
			t.Keterangan = barang.Nama
		}
	} else {
		if req.Jumlah <= 0 {
			return nil, errors.New("jumlah harus lebih dari 0")
		}
		t.Jumlah = req.Jumlah
		if t.Keterangan == "" {
			t.Keterangan = strings.ToUpper(req.Tipe[:1]) + req.Tipe[1:]
		}
	}

	if err := s.repo.Simpan(userID, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TransaksiService) Daftar(userID int, tipe string) ([]models.Transaksi, error) {
	if tipe != "" && tipe != models.TipePemasukan && tipe != models.TipePengeluaran {
		return nil, errors.New("filter tipe tidak valid")
	}
	return s.repo.AmbilSemua(userID, tipe, batasRiwayat)
}

func (s *TransaksiService) Hapus(userID, id int) error {
	return s.repo.Hapus(userID, id)
}

func (s *TransaksiService) Ringkasan(userID int) (*models.Ringkasan, error) {
	return s.repo.Ringkasan(userID)
}
