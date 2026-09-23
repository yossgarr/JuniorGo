package controller

import (
	"context"
	"encoding/json"
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

// 1. Endpoint Login: Mengirimkan HttpOnly Cookie
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

	// Tanam token ke HttpOnly Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    res.Token,
		Path:     "/",
		MaxAge:   86400, // 24 Jam dalam detik
		HttpOnly: true,  // Mencegah pencurian token lewat skrip JS (Anti XSS)
		Secure:   false, // Set true jika memakai HTTPS di production
		SameSite: http.SameSiteLaxMode,
	})

	response.Sukses(w, http.StatusOK, "Login berhasil", map[string]string{
		"username": res.Username,
	})
}

// 2. Endpoint Logout: Menghapus Cookie
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Langsung kadaluarsa / terhapus
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	response.Sukses(w, http.StatusOK, "Berhasil keluar (Cookie dihapus)", nil)
}

// 3. Endpoint Me: Mengecek status user yang sedang login dari Cookie
func (c *AuthController) CekMe(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	response.Sukses(w, http.StatusOK, "Sesi aktif", map[string]interface{}{
		"user": user,
	})
}

func (c *AuthController) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	config := utils.DapatkanGoogleOAuthConfig()
	url := config.AuthCodeURL("state-acak-rahasia")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// 4. Update Callback Google OAuth agar menanam Cookie dan redirect tanpa query token
func (c *AuthController) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// ... (kode penukaran kode Google OAuth tetap sama) ...
	state := r.URL.Query().Get("state")
	if state != "state-acak-rahasia" {
		response.Gagal(w, http.StatusBadRequest, "State tidak cocok")
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
		response.Gagal(w, http.StatusInternalServerError, "Gagal mengambil data user")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var profil struct {
		Email string `json:"email"`
	}
	json.Unmarshal(body, &profil)

	user, err := c.svc.DaftarAtauAmbilUserOAuth(profil.Email)
	if err != nil {
		response.Gagal(w, http.StatusInternalServerError, "Gagal registrasi user")
		return
	}

	jwtToken, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.Gagal(w, http.StatusInternalServerError, "Gagal membuat sesi")
		return
	}

	// Tanam token ke Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect bersih ke React TANPA mengekspos token di URL lagi!
	http.Redirect(w, r, "http://localhost:5173", http.StatusSeeOther)
}