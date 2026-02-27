package middleware

import (
	"net/http"
	"strings"
)

// CacheControlMiddleware disables caching for API routes and sets long cache for /assets
func CacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/assets") {
			// The content of the assets directory has fingerprinted file names so we cache aggressively
			w.Header().Set("Cache-Control", "public, max-age=2592000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}

		next.ServeHTTP(w, r)
	})
}
