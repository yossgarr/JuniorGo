package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// Menghasilkan string acak sepanjang 64 karakter heksadesimal (32 byte)
func GenerateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}