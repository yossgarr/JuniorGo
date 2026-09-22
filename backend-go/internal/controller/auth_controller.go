package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"backend-go/internal/models"
	"backend-go/internal/service"
	"backend-go/pkg/response"
	"backend-go/pkg/utils"
)

type AuthController struct {
	svc *service.UserService
}

func NewAuthController(s *service.UserService) *AuthController {
	return &AuthController{svc: s}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Gagal(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
		return
	}

	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Gagal(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	if err := c.svc.Register(&req); err != nil {
		response.Gagal(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Sukses(w, http.StatusCreated, "Registrasi berhasil", nil)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Gagal(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
		return
	}

	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Gagal(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	res, err := c.svc.Login(req.Username, req.Password)
	if err != nil {
		response.Gagal(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.Sukses(w, http.StatusOK, "Login berhasil", res)
}

func (c *AuthController) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	config := utils.DapatkanGoogleOAuthConfig()
	url := config.AuthCodeURL("state-acak-rahasia")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (c *AuthController) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != "state-acak-rahasia" {
		response.Gagal(w, http.StatusBadRequest, "State tidak cocok (indikasi CSRF)")
		return
	}

	code := r.URL.Query().Get("code")
	config := utils.DapatkanGoogleOAuthConfig()

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		response.Gagal(w, http.StatusInternalServerError, "Gagal menukar kode otorisasi")
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		response.Gagal(w, http.StatusInternalServerError, "Gagal mengambil data user dari Google")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var profil struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	json.Unmarshal(body, &profil)

	user, err := c.svc.DaftarAtauAmbilUserOAuth(profil.Email)
	if err != nil {
		response.Gagal(w, http.StatusInternalServerError, "Gagal registrasi user OAuth")
		return
	}

	jwtToken, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.Gagal(w, http.StatusInternalServerError, "Gagal membuat sesi JWT")
		return
	}

	// Arahkan ke React di port 5173
	redirectURL := fmt.Sprintf("http://localhost:5173?token=%s&user=%s", jwtToken, user.Username)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}