package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func AmbilJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("kunci_cadangan_default")
	}
	return []byte(secret)
}

func GenerateToken(userID int, username string) (string, error) {
	klaim := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Kadaluarsa 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, klaim)
	return token.SignedString(AmbilJWTSecret())
}

func ValidasiToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode enkripsi tidak valid")
		}
		return AmbilJWTSecret(), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("token tidak valid atau sudah kadaluarsa")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("format claims token salah")
	}

	return claims, nil
}