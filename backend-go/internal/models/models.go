package models

import "time"

type Buku struct {
	ID      int    `json:"id"`
	Judul   string `json:"judul"`
	Penulis string `json:"penulis"`
	Stok    int    `json:"stok"`
}

type PinjamRequest struct {
	BukuID int `json:"buku_id"`
}

type Peminjaman struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	BukuID    int       `json:"buku_id"`
	CreatedAt time.Time `json:"created_at"`
}