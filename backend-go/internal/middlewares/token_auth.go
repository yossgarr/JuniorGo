package middlewares

import (
	"context"
	"net/http"
	"strings"

	"backend-go/internal/repo"
	"backend-go/pkg/response"
)

type TokenAuthMiddleware struct {
	tokenRepo *repo.TokenRepository
}

func NewTokenAuthMiddleware(tr *repo.TokenRepository) *TokenAuthMiddleware {
	return &TokenAuthMiddleware{tokenRepo: tr}
}

func (m *TokenAuthMiddleware) RequireAPIToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokenStr string

		apiKeyHeader := r.Header.Get("X-API-KEY")
		if apiKeyHeader != "" {
			tokenStr = apiKeyHeader
		} else {
			authHeader := r.Header.Get("Authorization")
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenStr = parts[1]
			}
		}

		if tokenStr == "" {
			response.Gagal(w, http.StatusUnauthorized, "Token otentikasi wajib disertakan")
			return
		}

		userID, err := m.tokenRepo.CekValiditasToken(tokenStr)
		if err != nil {
			// Mengembalikan error spesifik: expired / revoked / not found
			response.Gagal(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next(w, r.WithContext(ctx))
	}
}