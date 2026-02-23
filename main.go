package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// API routes
	http.HandleFunc("/api/options/models", modelsHandler)
	http.HandleFunc("/api/conversations", conversationsHandler)
	http.HandleFunc("/api/settings", settingsHandler)
	http.HandleFunc("/api/github/repositories", repositoriesHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/alive", healthHandler)

	// Static file serving
	staticDir := "frontend/build"
	fs := http.FileServer(http.Dir(staticDir))

	// SPA fallback handler
	http.HandleFunc("/", spaHandler(staticDir, fs))

	log.Println("Starting server on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}

func spaHandler(staticDir string, fs http.Handler) http.HandlerFunc {
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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func modelsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Return a dummy model list
	models := []string{"gpt-4", "gpt-3.5-turbo", "claude-3-opus"}
	json.NewEncoder(w).Encode(models)
}

func conversationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Return empty list
	json.NewEncoder(w).Encode([]interface{}{})
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Return some default settings
	settings := map[string]interface{}{
		"LLM_MODEL": "gpt-4",
		"AGENT":     "CodeActAgent",
		"LANGUAGE":  "en",
	}
	json.NewEncoder(w).Encode(settings)
}

func repositoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]interface{}{})
}
