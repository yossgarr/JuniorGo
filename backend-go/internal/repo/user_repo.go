package repo

import (
	"backend-go/global"
	"backend-go/internal/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Simpan(u *models.User) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	return global.DB.QueryRow(query, u.Username, u.Password).Scan(&u.ID)
}

func (r *UserRepository) CariBerdasarkanUsername(username string) (*models.User, error) {
	var u models.User
	query := `SELECT id, username, password FROM users WHERE username = $1`
	err := global.DB.QueryRow(query, username).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) SimpanUserOAuth(username string) (*models.User, error) {
	// Cek apakah user sudah ada
	u, err := r.CariBerdasarkanUsername(username)
	if err == nil {
		return u, nil // User sudah terdaftar sebelumnya
	}

	// Jika belum ada, simpan user baru dengan password placeholder
	var id int
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	err = global.DB.QueryRow(query, username, "OAUTH_GOOGLE_USER").Scan(&id)
	if err != nil {
		return nil, err
	}

	return &models.User{ID: id, Username: username}, nil
}