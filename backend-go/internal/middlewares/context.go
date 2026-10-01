package middlewares

import "net/http"

// UserID mengambil ID user yang sudah diset oleh middleware auth
// (cookie/JWT maupun API token).
func UserID(r *http.Request) (int, bool) {
	id, ok := r.Context().Value("user_id").(int)
	return id, ok && id > 0
}
