package handlers

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/vlonebara/todo-list/pkg/db"
	"github.com/vlonebara/todo-list/pkg/service"
)

// testHandler собирает обработчики с отдельной БД и учебным паролем.
func testHandler(t *testing.T, password string) http.Handler {
	t.Helper()
	database, err := db.OpenDatabase(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = database.Close() })
	authentication, err := service.NewAuthService(password, "")
	if err != nil {
		t.Fatal(err)
	}
	return New(service.NewTaskService(db.NewTaskStore(database)), authentication, true, "../../web")
}

// TestAuthentication проверяет вход, защищённую cookie и выход.
func TestAuthentication(t *testing.T) {
	handler := testHandler(t, "test-password")
	request := httptest.NewRequest(http.MethodPost, "/api/signin", strings.NewReader(`{"password":"test-password"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	cookies := response.Result().Cookies()
	if response.Code != 200 || len(cookies) != 1 {
		t.Fatalf("signin failed: %d", response.Code)
	}
	cookie := cookies[0]
	if !cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteStrictMode || cookie.MaxAge != 28800 {
		t.Fatal("cookie attributes are not secure")
	}
	for _, item := range []struct {
		name   string
		cookie *http.Cookie
		status int
	}{
		{"no cookie", nil, 401},
		{"valid cookie", cookie, 200},
		{"invalid cookie", &http.Cookie{Name: "token", Value: "invalid"}, 401},
	} {
		t.Run(item.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
			if item.cookie != nil {
				request.AddCookie(item.cookie)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != item.status {
				t.Fatalf("got %d, want %d", response.Code, item.status)
			}
		})
	}
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/signout", nil))
	if response.Result().Cookies()[0].MaxAge != -1 {
		t.Fatal("signout did not clear cookie")
	}
}

// TestRequestProtection проверяет отказ для чужого сайта и неверного JSON.
func TestRequestProtection(t *testing.T) {
	handler := testHandler(t, "test-password")
	cases := []struct {
		name, body, contentType, origin string
		status                          int
	}{
		{"cross origin", `{}`, "application/json", "https://other.example", 403},
		{"wrong content type", `{}`, "text/plain", "", 415},
		{"null object", `null`, "application/json", "", 400},
		{"multiple objects", `{} {}`, "application/json", "", 400},
		{"oversized body", strings.Repeat(" ", maxRequestBytes+1), "application/json", "", 413},
	}
	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/signin", strings.NewReader(item.body))
			request.Header.Set("Content-Type", item.contentType)
			if item.origin != "" {
				request.Header.Set("Origin", item.origin)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != item.status {
				t.Fatalf("got %d, want %d", response.Code, item.status)
			}
		})
	}
}

// TestSigninLimiter проверяет лимиты и повторный вход после окончания окна.
func TestSigninLimiter(t *testing.T) {
	limiter := signinLimiter{clients: make(map[string]signinAttempts)}
	now := time.Now()
	for i := 0; i < 5; i++ {
		if !limiter.allow("client", now) {
			t.Fatal("attempt was blocked too early")
		}
	}
	if limiter.allow("client", now) {
		t.Fatal("extra attempt was allowed")
	}
	if !limiter.allow("client", now.Add(time.Minute)) {
		t.Fatal("expired limit was not reset")
	}
}

// TestLocalHostProtection блокирует DNS-имя постороннего сайта в режиме без пароля.
func TestLocalHostProtection(t *testing.T) {
	handler := testHandler(t, "")
	for _, item := range []struct {
		name, host string
		status     int
	}{
		{"localhost", "localhost:7540", 200},
		{"loopback ipv4", "127.0.0.1:7540", 200},
		{"loopback ipv6", "[::1]:7540", 200},
		{"external name", "other.example:7540", 403},
		{"misleading name", "localhost.other.example:7540", 403},
	} {
		t.Run(item.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
			request.Host = item.host
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != item.status {
				t.Fatalf("got %d, want %d", response.Code, item.status)
			}
		})
	}
}
