package repo

import (
	"backend-go/global"
	"backend-go/internal/models"
)

type BukuRepository struct{}

func NewBukuRepository() *BukuRepository {
	return &BukuRepository{}
}

func (r *BukuRepository) AmbilSemua() ([]models.Buku, error) {
	rows, err := global.DB.Query("SELECT id, judul, penulis FROM buku ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	daftar := []models.Buku{}
	for rows.Next() {
		var b models.Buku
		if err := rows.Scan(&b.ID, &b.Judul, &b.Penulis); err != nil {
			return nil, err
		}
		daftar = append(daftar, b)
	}
	return daftar, nil
}

func (r *BukuRepository) Simpan(buku *models.Buku) error {
	query := `INSERT INTO buku (judul, penulis) VALUES ($1, $2) RETURNING id`
	return global.DB.QueryRow(query, buku.Judul, buku.Penulis).Scan(&buku.ID)
}

func (r *BukuRepository) Perbarui(buku *models.Buku) error {
	query := `UPDATE buku SET judul = $1, penulis = $2 WHERE id = $3`
	_, err := global.DB.Exec(query, buku.Judul, buku.Penulis, buku.ID)
	return err
}

func (r *BukuRepository) Hapus(id int) error {
	query := `DELETE FROM buku WHERE id = $1`
	_, err := global.DB.Exec(query, id)
	return err
}