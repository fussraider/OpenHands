package middleware

import (
	"context"
	"net/http"
	"openhands-go/server/session"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for public endpoints
		if isPublicEndpoint(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// For now, allow unauthenticated access for development/migration
			// In production, this should be 401
			// http.Error(w, "Authorization header required", http.StatusUnauthorized)
			// return

			// Mock user ID for dev
			ctx := context.WithValue(r.Context(), UserIDKey, "mock-user-id")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		sess, err := session.GetSession(token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, sess.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isPublicEndpoint(path string) bool {
	publicPaths := []string{
		"/health",
		"/alive",
		"/api/options",
		"/", // Static files
	}
	for _, pp := range publicPaths {
		if strings.HasPrefix(path, pp) {
			return true
		}
	}
	// Allow static files check more robustly?
	// For now this covers / and subpaths which might be static files
	// But API routes are under /api
	if !strings.HasPrefix(path, "/api") {
		return true
	}

	return false
}
