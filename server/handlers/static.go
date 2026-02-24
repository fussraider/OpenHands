package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func SPAHandler(staticDir string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(staticDir))
	return func(w http.ResponseWriter, r *http.Request) {
		// If it's an API call that wasn't handled, return 404
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Ensure path is relative and safe
		cleanPath := filepath.Clean(r.URL.Path)
		path := filepath.Join(staticDir, strings.TrimPrefix(cleanPath, "/"))

		// check if file exists
		stat, err := os.Stat(path)
		if os.IsNotExist(err) {
			log.Printf("SPA: File not found: %s, serving index.html", path)
			// serve index.html
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		} else if err != nil {
			log.Printf("SPA: Error stating file %s: %v", path, err)
			// other error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if stat.IsDir() {
			// Serve index.html if directory
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}

		log.Printf("SPA: Serving file: %s", path)
		// file exists, serve it
		fs.ServeHTTP(w, r)
	}
}
