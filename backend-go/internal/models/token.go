package models

import "time"

type APIToken struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Token     string     `json:"token"`
	NamaToken string     `json:"nama_token"`
	IsRevoked bool       `json:"is_revoked"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at"` // Pointer agar bisa null jika token permanen
}

type BuatTokenRequest struct {
	NamaToken     string `json:"nama_token"`
	MasaBerlakuHari int   `json:"masa_berlaku_hari"` // Contoh: 7, 30 hari (0 = permanen)
}

type RevokeTokenRequest struct {
	Token string `json:"token"`
}