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
		var tokenStr string

		// 1. Cek dari HTTP Cookie terlebih dahulu
		cookie, err := r.Cookie("auth_token")
		if err == nil && cookie.Value != "" {
			tokenStr = cookie.Value
		} else {
			// 2. Fallback: Cek header Authorization Bearer (jika ada request non-browser)
			authHeader := r.Header.Get("Authorization")
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenStr = parts[1]
			}
		}

		if tokenStr == "" {
			response.Gagal(w, http.StatusUnauthorized, "Sesi login tidak ditemukan (Cookie tidak ada)")
			return
		}

		claims, err := utils.ValidasiToken(tokenStr)
		if err != nil {
			response.Gagal(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims["username"])
		next(w, r.WithContext(ctx))
	}
}