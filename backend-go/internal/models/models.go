package models

type Buku struct {
	ID      int    `json:"id"`
	Judul   string `json:"judul"`
	Penulis string `json:"penulis"`
}