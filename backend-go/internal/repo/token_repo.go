package repo

import (
	"errors"
	"time"

	"backend-go/global"
	"backend-go/internal/models"
)

type TokenRepository struct{}

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{}
}

func (r *TokenRepository) Simpan(token *models.APIToken) error {
	query := `
		INSERT INTO api_tokens (user_id, token, nama_token, is_revoked, expires_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at`
	return global.DB.QueryRow(
		query,
		token.UserID,
		token.Token,
		token.NamaToken,
		token.IsRevoked,
		token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)
}

func (r *TokenRepository) CekValiditasToken(tokenStr string) (int, error) {
	var userID int
	var isRevoked bool
	var expiresAt *time.Time

	query := `SELECT user_id, is_revoked, expires_at FROM api_tokens WHERE token = $1`
	err := global.DB.QueryRow(query, tokenStr).Scan(&userID, &isRevoked, &expiresAt)
	if err != nil {
		return 0, errors.New("token tidak ditemukan")
	}

	// 1. Cek apakah token sudah dicabut (revoked)
	if isRevoked {
		return 0, errors.New("token telah dicabut dan tidak dapat digunakan")
	}

	// 2. Cek apakah token sudah kadaluarsa
	if expiresAt != nil && time.Now().After(*expiresAt) {
		return 0, errors.New("token telah kadaluarsa")
	}

	return userID, nil
}

// Revoke: Nonaktifkan token secara permanen
func (r *TokenRepository) CabutToken(tokenStr string) error {
	query := `UPDATE api_tokens SET is_revoked = TRUE WHERE token = $1`
	res, err := global.DB.Exec(query, tokenStr)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("token tidak ditemukan")
	}
	return nil
}