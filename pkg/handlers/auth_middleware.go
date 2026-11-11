package handlers

import "net/http"

// auth проверяет JWT cookie, если для сервера настроен пароль.
func auth(authentication AuthService, next http.HandlerFunc) http.HandlerFunc {
	if !authentication.Enabled() {
		return next
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if !validTokenCookie(r, authentication) {
			writeAPIError(w, http.StatusUnauthorized, "authentication required")
			return
		}
		next(w, r)
	}
}

// validTokenCookie читает cookie и проверяет её одним способом для всех запросов.
func validTokenCookie(r *http.Request, authentication AuthService) bool {
	cookie, err := r.Cookie("token")
	return err == nil && authentication.ValidateToken(cookie.Value)
}
