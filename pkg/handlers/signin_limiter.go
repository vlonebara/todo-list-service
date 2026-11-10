package handlers

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type signinAttempts struct {
	count int
	until time.Time
}

type signinLimiter struct {
	mu      sync.Mutex
	clients map[string]signinAttempts
	all     signinAttempts
}

// limitSignin ограничивает вход пятью попытками в минуту для одного адреса.
func limitSignin(next http.HandlerFunc) http.HandlerFunc {
	limiter := signinLimiter{clients: make(map[string]signinAttempts)}
	return func(w http.ResponseWriter, r *http.Request) {
		address, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			address = r.RemoteAddr
		}
		if !limiter.allow(address, time.Now()) {
			w.Header().Set("Retry-After", "60")
			writeAPIError(w, http.StatusTooManyRequests, "too many signin attempts")
			return
		}
		next(w, r)
	}
}

// allow проверяет лимит адреса и общий лимит без фоновых горутин.
func (l *signinLimiter) allow(address string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !now.Before(l.all.until) {
		l.all = signinAttempts{until: now.Add(time.Minute)}
		for key, attempts := range l.clients {
			if !now.Before(attempts.until) {
				delete(l.clients, key)
			}
		}
	}
	attempts := l.clients[address]
	if !now.Before(attempts.until) {
		attempts = signinAttempts{until: now.Add(time.Minute)}
	}
	if attempts.count >= 5 || l.all.count >= 100 {
		return false
	}
	attempts.count++
	l.all.count++
	l.clients[address] = attempts
	return true
}
