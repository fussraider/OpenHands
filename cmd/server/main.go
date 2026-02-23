package main

import (
	"log"
	"net/http"
	"openhands-go/server/handlers"
	"openhands-go/server/store"
)

func main() {
	// Initialize stores
	store.NewConversationStore()

	// API routes
	http.HandleFunc("GET /api/options/models", handlers.ModelsHandler)

	http.HandleFunc("GET /api/conversations", handlers.SearchConversationsHandler)
	http.HandleFunc("POST /api/conversations", handlers.NewConversationHandler)
	http.HandleFunc("GET /api/conversations/{id}", handlers.GetConversationHandler)

	http.HandleFunc("GET /api/conversations/{id}/list-files", handlers.ListFilesHandler)
	http.HandleFunc("GET /api/conversations/{id}/select-file", handlers.SelectFileHandler)

	http.HandleFunc("GET /api/settings", handlers.GetSettingsHandler)
	http.HandleFunc("POST /api/settings", handlers.StoreSettingsHandler)

	// Secrets API - note: actual paths need to be verified against Python implementation
	// Python: /api/secrets (GET, POST), /api/secrets/{key} (DELETE)
	http.HandleFunc("GET /api/secrets", handlers.GetSecretsHandler)
	http.HandleFunc("POST /api/secrets", handlers.StoreSecretHandler)
	http.HandleFunc("DELETE /api/secrets/{key}", handlers.DeleteSecretHandler)

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
