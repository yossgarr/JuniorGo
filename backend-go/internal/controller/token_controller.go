package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"backend-go/internal/models"
	"backend-go/internal/repo"
	"backend-go/pkg/response"
	"backend-go/pkg/utils"
)

type TokenController struct {
	tokenRepo *repo.TokenRepository
}

func NewTokenController(tr *repo.TokenRepository) *TokenController {
	return &TokenController{tokenRepo: tr}
}

func (c *TokenController) BuatTokenBaru(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Gagal(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
		return
	}

	var req models.BuatTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.NamaToken == "" {
		response.Gagal(w, http.StatusBadRequest, "Nama token wajib diisi")
		return
	}

	randomStr, err := utils.GenerateRandomToken()
	if err != nil {
		response.Gagal(w, http.StatusInternalServerError, "Gagal membuat token")
		return
	}

	tokenBaru := models.APIToken{
		UserID:    1,
		Token:     randomStr,
		NamaToken: req.NamaToken,
		IsRevoked: false,
	}

	// Tentukan masa berlaku jika dikirimkan (dalam hari)
	if req.MasaBerlakuHari > 0 {
		exp := time.Now().Add(time.Duration(req.MasaBerlakuHari) * 24 * time.Hour)
		tokenBaru.ExpiresAt = &exp
	}

	if err := c.tokenRepo.Simpan(&tokenBaru); err != nil {
		response.Gagal(w, http.StatusInternalServerError, "Gagal menyimpan token")
		return
	}

	response.Sukses(w, http.StatusCreated, "Token berhasil dibuat", tokenBaru)
}

func (c *TokenController) CabutToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Gagal(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
		return
	}

	var req models.RevokeTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" {
		response.Gagal(w, http.StatusBadRequest, "Parameter token wajib diisi")
		return
	}

	if err := c.tokenRepo.CabutToken(req.Token); err != nil {
		response.Gagal(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Sukses(w, http.StatusOK, "Token berhasil dicabut / dinonaktifkan", nil)
}