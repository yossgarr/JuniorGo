package middlewares

import (
	"context"
	"net/http"
	"strings"

	"backend-go/pkg/response"
	"backend-go/pkg/utils"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Gagal(w, http.StatusUnauthorized, "Token otentikasi tidak ditemukan")
			return
		}

		// Format token: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Gagal(w, http.StatusUnauthorized, "Format Authorization header harus 'Bearer <token>'")
			return
		}

		claims, err := utils.ValidasiToken(parts[1])
		if err != nil {
			response.Gagal(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Simpan username ke dalam context request
		ctx := context.WithValue(r.Context(), "user", claims["username"])
		next(w, r.WithContext(ctx))
	}
}