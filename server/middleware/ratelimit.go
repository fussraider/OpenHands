package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

type RateLimiter struct {
	mu           sync.Mutex
	history      map[string][]time.Time
	requests     int
	seconds      time.Duration
	sleepSeconds time.Duration
}

func NewRateLimiter(requests int, seconds int, sleepSeconds int) *RateLimiter {
	return &RateLimiter{
		history:      make(map[string][]time.Time),
		requests:     requests,
		seconds:      time.Duration(seconds) * time.Second,
		sleepSeconds: time.Duration(sleepSeconds) * time.Second,
	}
}

func (rl *RateLimiter) cleanOldRequests(key string, now time.Time) {
	cutoff := now.Add(-rl.seconds)
	var newHistory []time.Time
	for _, ts := range rl.history[key] {
		if ts.After(cutoff) {
			newHistory = append(newHistory, ts)
		}
	}
	rl.history[key] = newHistory
}

func (rl *RateLimiter) Allow(ip string) (bool, bool) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	rl.cleanOldRequests(ip, now)
	rl.history[ip] = append(rl.history[ip], now)

	count := len(rl.history[ip])

	if count > rl.requests*2 {
		return false, false // Hard reject
	} else if count > rl.requests {
		if rl.sleepSeconds > 0 {
			// We will sleep outside the lock
			return true, true // Soft reject (requires sleep)
		}
		return false, false
	}
	return true, false
}

func RateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isRateLimitedRequest(r) {
				ip := getIP(r)
				allow, sleep := limiter.Allow(ip)

				if !allow {
					w.Header().Set("Retry-After", "1")
					http.Error(w, `{"message": "Too many requests"}`, http.StatusTooManyRequests)
					return
				}

				if sleep {
					time.Sleep(limiter.sleepSeconds)
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isRateLimitedRequest(r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, "/assets") {
		return false
	}
	return true
}

func getIP(r *http.Request) string {
	// Basic implementation, consider X-Forwarded-For in real env
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
}
