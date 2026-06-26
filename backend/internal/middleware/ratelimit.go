package middleware

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RateLimit is a simple in-memory, per-IP fixed-window rate limiter.
// It protects sensitive public endpoints (login, register, contact) from
// brute force and abuse, without external dependencies.
//
// limit = max requests allowed per window, per client IP.
func RateLimit(limit int, window time.Duration) func(http.Handler) http.Handler {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
	}
	go rl.cleanupLoop()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !rl.allow(clientIP(r)) {
				w.Header().Set("Retry-After", strconv.Itoa(int(window.Seconds())))
				http.Error(w, `{"error":"too many requests, try again later"}`, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type visitor struct {
	count   int
	resetAt time.Time
}

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	v, ok := rl.visitors[ip]
	if !ok || now.After(v.resetAt) {
		rl.visitors[ip] = &visitor{count: 1, resetAt: now.Add(rl.window)}
		return true
	}
	if v.count >= rl.limit {
		return false
	}
	v.count++
	return true
}

func (rl *rateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, v := range rl.visitors {
			if now.After(v.resetAt) {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// clientIP extracts the client IP. chi's RealIP middleware already normalizes
// RemoteAddr from X-Forwarded-For / X-Real-IP when behind a proxy.
func clientIP(r *http.Request) string {
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}
