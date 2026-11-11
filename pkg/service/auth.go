package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/vlonebara/todo-list/pkg/errutil"
)

const (
	TokenLifetime = 8 * time.Hour
	jwtSecretSize = 32
)

// AuthService хранит пароль и ключ, подготовленные при запуске сервера.
type AuthService struct {
	password string
	secret   []byte
}

type authClaims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

// NewAuthService создаёт сервис аутентификации с готовыми настройками.
func NewAuthService(password, secret string) (AuthService, error) {
	if password == "" {
		return AuthService{}, nil
	}
	key, err := jwtSecret(secret)
	if errutil.CheckError(err) {
		return AuthService{}, err
	}
	return AuthService{password: password, secret: key}, nil
}

// Enabled сообщает, включена ли проверка аутентификации.
func (s AuthService) Enabled() bool {
	return s.password != ""
}

// Authenticate проверяет пароль и возвращает подписанный JWT-токен.
func (s AuthService) Authenticate(password string) (string, error) {
	if !s.Enabled() {
		return "", errors.New("password is not configured")
	}
	if !hmac.Equal([]byte(password), []byte(s.password)) {
		return "", errors.New("incorrect password")
	}
	return createAuthToken(s.password, s.secret, time.Now())
}

// ValidateToken проверяет подпись, срок действия и парольную метку JWT.
func (s AuthService) ValidateToken(token string) bool {
	if !s.Enabled() {
		return false
	}

	var claims authClaims
	parsed, err := jwt.ParseWithClaims(token, &claims, func(_ *jwt.Token) (any, error) {
		return s.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
	if errutil.CheckError(err) || !parsed.Valid {
		return false
	}
	return hmac.Equal([]byte(claims.PasswordHash), []byte(passwordChecksum(s.password, s.secret)))
}

// jwtSecret проверяет заданный ключ или создаёт случайный для сервиса.
func jwtSecret(value string) ([]byte, error) {
	if value != "" {
		if len(value) < jwtSecretSize {
			return nil, errors.New("jwt secret must contain at least 32 bytes")
		}
		return []byte(value), nil
	}
	secret := make([]byte, jwtSecretSize)
	if _, err := rand.Read(secret); errutil.CheckError(err) {
		return nil, fmt.Errorf("generate jwt secret: %w", err)
	}
	return secret, nil
}

// createAuthToken формирует JWT с меткой пароля и временем истечения.
func createAuthToken(password string, secret []byte, now time.Time) (string, error) {
	claims := authClaims{
		PasswordHash: passwordChecksum(password, secret),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(TokenLifetime)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if errutil.CheckError(err) {
		return "", fmt.Errorf("sign jwt: %w", err)
	}
	return token, nil
}

// passwordChecksum создаёт подписанную метку для проверки пароля при валидации.
func passwordChecksum(password string, secret []byte) string {
	checksum := hmac.New(sha256.New, secret)
	_, _ = checksum.Write([]byte("password-checksum:"))
	_, _ = checksum.Write([]byte(password))
	return base64.RawURLEncoding.EncodeToString(checksum.Sum(nil))
}
