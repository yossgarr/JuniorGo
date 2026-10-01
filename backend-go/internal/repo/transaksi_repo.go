package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"backend-go/global"
	"backend-go/internal/models"
)

var ErrTransaksiTidakDitemukan = errors.New("transaksi tidak ditemukan")

type TransaksiRepository struct{}

func NewTransaksiRepository() *TransaksiRepository {
	return &TransaksiRepository{}
}

func (r *TransaksiRepository) Simpan(userID int, t *models.Transaksi) error {
	return global.DB.QueryRow(
		`INSERT INTO transaksi (user_id, tipe, jumlah, keterangan, barang_id, nama_barang, qty)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, created_at`,
		userID, t.Tipe, t.Jumlah, t.Keterangan, t.BarangID, t.NamaBarang, t.Qty,
	).Scan(&t.ID, &t.CreatedAt)
}

// AmbilSemua mengembalikan transaksi terbaru lebih dulu. tipe kosong = semua.
func (r *TransaksiRepository) AmbilSemua(userID int, tipe string, limit int) ([]models.Transaksi, error) {
	rows, err := global.DB.Query(
		`SELECT id, tipe, jumlah, keterangan, barang_id, nama_barang, qty, created_at
		 FROM transaksi
		 WHERE user_id = $1 AND ($2 = '' OR tipe = $2)
		 ORDER BY created_at DESC, id DESC
		 LIMIT $3`, userID, tipe, limit)
	if err != nil {
		return nil, fmt.Errorf("gagal query transaksi: %w", err)
	}
	defer rows.Close()

	daftar := []models.Transaksi{}
	for rows.Next() {
		var t models.Transaksi
		var barangID sql.NullInt64
		if err := rows.Scan(&t.ID, &t.Tipe, &t.Jumlah, &t.Keterangan, &barangID, &t.NamaBarang, &t.Qty, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("gagal scan transaksi: %w", err)
		}
		if barangID.Valid {
			id := int(barangID.Int64)
			t.BarangID = &id
		}
		daftar = append(daftar, t)
	}
	return daftar, rows.Err()
}

func (r *TransaksiRepository) Hapus(userID, id int) error {
	res, err := global.DB.Exec(`DELETE FROM transaksi WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrTransaksiTidakDitemukan
	}
	return nil
}

func (r *TransaksiRepository) Ringkasan(userID int) (*models.Ringkasan, error) {
	var s models.Ringkasan
	err := global.DB.QueryRow(
		`SELECT
		   COALESCE(SUM(CASE WHEN tipe = 'pemasukan'   THEN jumlah END), 0),
		   COALESCE(SUM(CASE WHEN tipe = 'pengeluaran' THEN jumlah END), 0)
		 FROM transaksi WHERE user_id = $1`, userID,
	).Scan(&s.TotalPemasukan, &s.TotalPengeluaran)
	if err != nil {
		return nil, err
	}
	s.Saldo = s.TotalPemasukan - s.TotalPengeluaran
	return &s, nil
}
