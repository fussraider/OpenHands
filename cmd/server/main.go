package main

import (
	"log"
	"net/http"
	"openhands-go/server/handlers"
)

func main() {
	// API routes
	http.HandleFunc("GET /api/options/models", handlers.ModelsHandler)

	http.HandleFunc("GET /api/conversations", handlers.SearchConversationsHandler)
	http.HandleFunc("POST /api/conversations", handlers.NewConversationHandler)
	http.HandleFunc("GET /api/conversations/{id}", handlers.GetConversationHandler)

	http.HandleFunc("GET /api/settings", handlers.GetSettingsHandler)
	http.HandleFunc("POST /api/settings", handlers.StoreSettingsHandler)

	http.HandleFunc("GET /api/github/repositories", handlers.RepositoriesHandler)

	http.HandleFunc("GET /health", handlers.HealthHandler)
	http.HandleFunc("GET /alive", handlers.HealthHandler)

	// Static file serving
	staticDir := "frontend/build"
	// SPA fallback handler
	http.HandleFunc("/", handlers.SPAHandler(staticDir))

	log.Println("Starting server on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
