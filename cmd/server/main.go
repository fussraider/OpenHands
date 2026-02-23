package main

import (
	"log"
	"net/http"
	"openhands-go/server/handlers"
)

func main() {
	// API routes
	http.HandleFunc("/api/options/models", handlers.ModelsHandler)
	http.HandleFunc("/api/conversations", handlers.ConversationsHandler)
	http.HandleFunc("/api/settings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetSettingsHandler(w, r)
		} else if r.Method == http.MethodPost {
			handlers.StoreSettingsHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/github/repositories", handlers.RepositoriesHandler)
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/alive", handlers.HealthHandler)

	// Static file serving
	staticDir := "frontend/build"
	// SPA fallback handler
	http.HandleFunc("/", handlers.SPAHandler(staticDir))

	log.Println("Starting server on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
