package repo

import (
	"errors"
	"fmt"

	"backend-go/global"
	"backend-go/internal/models"
)

type BukuRepository struct{}

func NewBukuRepository() *BukuRepository {
	return &BukuRepository{}
}

func (r *BukuRepository) AmbilSemua() ([]models.Buku, error) {
	query := `SELECT id, judul, penulis, COALESCE(stok, 0) FROM buku ORDER BY id ASC`
	rows, err := global.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("gagal query database: %w", err)
	}
	defer rows.Close()

	daftar := []models.Buku{}
	for rows.Next() {
		var b models.Buku
		if err := rows.Scan(&b.ID, &b.Judul, &b.Penulis, &b.Stok); err != nil {
			return nil, fmt.Errorf("gagal scan baris database: %w", err)
		}
		daftar = append(daftar, b)
	}
	return daftar, nil
}

func (r *BukuRepository) Simpan(buku *models.Buku) error {
	query := `INSERT INTO buku (judul, penulis, stok) VALUES ($1, $2, $3) RETURNING id`
	return global.DB.QueryRow(query, buku.Judul, buku.Penulis, buku.Stok).Scan(&buku.ID)
}

func (r *BukuRepository) Perbarui(buku *models.Buku) error {
	query := `UPDATE buku SET judul = $1, penulis = $2, stok = $3 WHERE id = $4`
	_, err := global.DB.Exec(query, buku.Judul, buku.Penulis, buku.Stok, buku.ID)
	return err
}

func (r *BukuRepository) Hapus(id int) error {
	query := `DELETE FROM buku WHERE id = $1`
	_, err := global.DB.Exec(query, id)
	return err
}

// Transaksi ACID Peminjaman Buku
func (r *BukuRepository) PinjamBukuTx(userID int, bukuID int) error {
	tx, err := global.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi database: %w", err)
	}
	defer tx.Rollback()

	// 1. Kunci baris buku & cek ketersediaan stok
	var stok int
	err = tx.QueryRow("SELECT stok FROM buku WHERE id = $1 FOR UPDATE", bukuID).Scan(&stok)
	if err != nil {
		return fmt.Errorf("buku tidak ditemukan: %w", err)
	}

	if stok < 1 {
		return errors.New("stok buku habis, tidak dapat dipinjam")
	}

	// 2. Kurangi stok buku
	_, err = tx.Exec("UPDATE buku SET stok = stok - 1 WHERE id = $1", bukuID)
	if err != nil {
		return fmt.Errorf("gagal mengupdate stok: %w", err)
	}

	// 3. Catat riwayat peminjaman
	_, err = tx.Exec("INSERT INTO peminjaman (user_id, buku_id) VALUES ($1, $2)", userID, bukuID)
	if err != nil {
		return fmt.Errorf("gagal mencatat peminjaman: %w", err)
	}

	return tx.Commit()
}