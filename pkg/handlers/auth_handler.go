package handlers

import (
	"net/http"
	"time"

	"github.com/vlonebara/todo-list/pkg/errutil"
	"github.com/vlonebara/todo-list/pkg/service"
)

type authRequest struct {
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token,omitempty"`
}

// authHandler проверяет пароль и возвращает JWT-токен.
func (h Handler) authHandler(w http.ResponseWriter, r *http.Request) {
	var request authRequest
	if !readJSONRequest(w, r, &request) {
		return
	}

	token, err := h.authentication.Authenticate(request.Password)
	if errutil.CheckError(err) {
		writeAPIError(w, http.StatusUnauthorized, err.Error())
		return
	}
	http.SetCookie(w, h.tokenCookie(token, int(service.TokenLifetime.Seconds())))
	writeJSON(w, http.StatusOK, authResponse{Token: token})
}

// signoutHandler удаляет cookie текущего браузера.
func (h Handler) signoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, h.tokenCookie("", -1))
	writeJSON(w, http.StatusOK, struct{}{})
}

// sessionHandler сообщает интерфейсу, нужен ли вход и действительна ли cookie.
func (h Handler) sessionHandler(w http.ResponseWriter, r *http.Request) {
	required := h.authentication.Enabled()
	authenticated := !required || validTokenCookie(r, h.authentication)
	writeJSON(w, http.StatusOK, map[string]bool{"required": required, "authenticated": authenticated})
}

// tokenCookie задаёт одинаковые атрибуты cookie для входа и выхода.
func (h Handler) tokenCookie(token string, maxAge int) *http.Cookie {
	expires := time.Now().Add(service.TokenLifetime)
	if maxAge < 0 {
		expires = time.Unix(1, 0)
	}
	return &http.Cookie{
		Name: "token", Value: token, Path: "/", MaxAge: maxAge, Expires: expires,
		HttpOnly: true, Secure: h.secureCookies, SameSite: http.SameSiteStrictMode,
	}
}
