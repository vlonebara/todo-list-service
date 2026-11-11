package service

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// signTestToken подписывает проверяемые варианты JWT заданным алгоритмом.
func signTestToken(t *testing.T, method jwt.SigningMethod, claims authClaims, secret []byte) string {
	t.Helper()
	token, err := jwt.NewWithClaims(method, claims).SignedString(secret)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

// TestAuthenticationTokens проверяет срок, алгоритм, ключ и парольную метку JWT.
func TestAuthenticationTokens(t *testing.T) {
	authentication, err := NewAuthService("test-password", "test-only-secret-at-least-32-bytes")
	if err != nil {
		t.Fatal(err)
	}
	valid, err := authentication.Authenticate("test-password")
	if err != nil {
		t.Fatal(err)
	}
	expired, err := createAuthToken(authentication.password, authentication.secret, time.Now().Add(-TokenLifetime-time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	claims := authClaims{PasswordHash: passwordChecksum(authentication.password, authentication.secret),
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))}}
	otherAlgorithm := signTestToken(t, jwt.SigningMethodHS384, claims, authentication.secret)
	noExpiry := signTestToken(t, jwt.SigningMethodHS256, authClaims{PasswordHash: claims.PasswordHash}, authentication.secret)
	claims.PasswordHash = "wrong password marker"
	wrongMarker := signTestToken(t, jwt.SigningMethodHS256, claims, authentication.secret)
	for _, item := range []struct {
		name, token string
		valid       bool
	}{
		{"valid token", valid, true},
		{"expired token", expired, false},
		{"wrong algorithm", otherAlgorithm, false},
		{"missing expiry", noExpiry, false},
		{"wrong password marker", wrongMarker, false},
		{"malformed token", "invalid", false},
	} {
		t.Run(item.name, func(t *testing.T) {
			if authentication.ValidateToken(item.token) != item.valid {
				t.Fatal("unexpected token validity")
			}
		})
	}
	changed, err := NewAuthService("changed-password", string(authentication.secret))
	if err != nil {
		t.Fatal(err)
	}
	if changed.ValidateToken(valid) {
		t.Fatal("old password token was accepted")
	}
	if _, err := authentication.Authenticate("wrong-password"); err == nil {
		t.Fatal("wrong password was accepted")
	}
	if _, err := NewAuthService("test-password", "short"); err == nil {
		t.Fatal("short key was accepted")
	}
}
