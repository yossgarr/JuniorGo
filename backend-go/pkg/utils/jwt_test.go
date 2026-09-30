package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateDanValidasiToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "rahasia-test")

	token, err := GenerateToken(7, "yosi")
	if err != nil {
		t.Fatalf("GenerateToken gagal: %v", err)
	}

	claims, err := ValidasiToken(token)
	if err != nil {
		t.Fatalf("ValidasiToken gagal: %v", err)
	}
	if claims["username"] != "yosi" {
		t.Errorf("username = %v, want yosi", claims["username"])
	}
	if claims["user_id"] != float64(7) {
		t.Errorf("user_id = %v, want 7", claims["user_id"])
	}
}

func TestValidasiToken_SecretBerbeda(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret-a")
	token, _ := GenerateToken(1, "budi")

	t.Setenv("JWT_SECRET", "secret-b")
	if _, err := ValidasiToken(token); err == nil {
		t.Error("token dengan secret berbeda seharusnya ditolak")
	}
}

func TestValidasiToken_Kadaluarsa(t *testing.T) {
	t.Setenv("JWT_SECRET", "rahasia-test")

	klaim := jwt.MapClaims{
		"username": "budi",
		"exp":      time.Now().Add(-time.Hour).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, klaim).SignedString(AmbilJWTSecret())
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ValidasiToken(token); err == nil {
		t.Error("token kadaluarsa seharusnya ditolak")
	}
}

func TestValidasiToken_Sampah(t *testing.T) {
	if _, err := ValidasiToken("bukan.token.valid"); err == nil {
		t.Error("token acak seharusnya ditolak")
	}
}
