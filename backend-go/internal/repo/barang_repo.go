package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"backend-go/global"
	"backend-go/internal/models"
)

var ErrBarangTidakDitemukan = errors.New("barang tidak ditemukan")

type BarangRepository struct{}

func NewBarangRepository() *BarangRepository {
	return &BarangRepository{}
}

func (r *BarangRepository) AmbilSemua(userID int) ([]models.Barang, error) {
	rows, err := global.DB.Query(
		`SELECT id, nama, harga FROM barang WHERE user_id = $1 ORDER BY nama ASC`, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal query barang: %w", err)
	}
	defer rows.Close()

	daftar := []models.Barang{}
	for rows.Next() {
		var b models.Barang
		if err := rows.Scan(&b.ID, &b.Nama, &b.Harga); err != nil {
			return nil, fmt.Errorf("gagal scan barang: %w", err)
		}
		daftar = append(daftar, b)
	}
	return daftar, rows.Err()
}

func (r *BarangRepository) AmbilByID(userID, id int) (*models.Barang, error) {
	var b models.Barang
	err := global.DB.QueryRow(
		`SELECT id, nama, harga FROM barang WHERE id = $1 AND user_id = $2`, id, userID,
	).Scan(&b.ID, &b.Nama, &b.Harga)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrBarangTidakDitemukan
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BarangRepository) Simpan(userID int, b *models.Barang) error {
	return global.DB.QueryRow(
		`INSERT INTO barang (user_id, nama, harga) VALUES ($1, $2, $3) RETURNING id`,
		userID, b.Nama, b.Harga,
	).Scan(&b.ID)
}

func (r *BarangRepository) Perbarui(userID int, b *models.Barang) error {
	res, err := global.DB.Exec(
		`UPDATE barang SET nama = $1, harga = $2 WHERE id = $3 AND user_id = $4`,
		b.Nama, b.Harga, b.ID, userID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrBarangTidakDitemukan
	}
	return nil
}

func (r *BarangRepository) Hapus(userID, id int) error {
	res, err := global.DB.Exec(`DELETE FROM barang WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrBarangTidakDitemukan
	}
	return nil
}
