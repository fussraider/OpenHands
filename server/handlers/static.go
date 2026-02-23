package handlers

import (
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

		path := filepath.Join(staticDir, r.URL.Path)
		// check if file exists
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			// serve index.html
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		} else if err != nil {
			// other error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// file exists, serve it
		fs.ServeHTTP(w, r)
	}
}
