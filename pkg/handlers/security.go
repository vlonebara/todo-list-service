package handlers

import (
	"net"
	"net/http"
	"strings"

	"github.com/vlonebara/todo-list/pkg/errutil"
)

// localRequestsOnly не позволяет чужому DNS-имени обращаться к API без пароля.
func localRequestsOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		if address, _, err := net.SplitHostPort(host); !errutil.CheckError(err) {
			host = address
		}
		host = strings.Trim(host, "[]")
		if !strings.EqualFold(host, "localhost") && !net.ParseIP(host).IsLoopback() {
			writeAPIError(w, http.StatusForbidden, "local access requires a localhost address")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// securityHeaders добавляет общие заголовки защиты и запрещает кешировать API.
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "same-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; object-src 'none'; base-uri 'none'; frame-ancestors 'none'; form-action 'self'")
		w.Header().Set("Cache-Control", "no-cache")
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Cache-Control", "no-store")
		}
		next.ServeHTTP(w, r)
	})
}
