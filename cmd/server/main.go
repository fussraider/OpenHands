package main

import (
	"log"
	"net/http"
	"openhands-go/server/config"
	"openhands-go/server/handlers"
	"openhands-go/server/middleware"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("GET /api/options/models", handlers.ModelsHandler)
	mux.HandleFunc("GET /api/options/agents", handlers.GetAgentsHandler)
	mux.HandleFunc("GET /api/options/security-analyzers", handlers.GetSecurityAnalyzersHandler)

	mux.HandleFunc("GET /api/conversations", handlers.SearchConversationsHandler)
	mux.HandleFunc("POST /api/conversations", handlers.NewConversationHandler)
	mux.HandleFunc("GET /api/conversations/{id}", handlers.GetConversationHandler)

	mux.HandleFunc("GET /api/conversations/{id}/list-files", handlers.ListFilesHandler)
	mux.HandleFunc("GET /api/conversations/{id}/select-file", handlers.SelectFileHandler)

	mux.HandleFunc("GET /api/conversations/{id}/trajectory", handlers.GetTrajectoryHandler)
	mux.HandleFunc("POST /api/conversations/{id}/submit-feedback", handlers.SubmitFeedbackHandler)

	// Security API
	mux.HandleFunc("/api/conversations/{id}/security/{path...}", handlers.SecurityAPIHandler)

	mux.HandleFunc("GET /api/settings", handlers.GetSettingsHandler)
	mux.HandleFunc("POST /api/settings", handlers.StoreSettingsHandler)

	mux.HandleFunc("GET /api/secrets", handlers.GetSecretsHandler)
	mux.HandleFunc("POST /api/secrets", handlers.StoreSecretHandler)
	mux.HandleFunc("DELETE /api/secrets/{key}", handlers.DeleteSecretHandler)

	mux.HandleFunc("GET /api/github/repositories", handlers.RepositoriesHandler)

	// MCP Mount
	mux.HandleFunc("/mcp/", handlers.MCPSSEHandler)

	mux.HandleFunc("GET /health", handlers.HealthHandler)
	mux.HandleFunc("GET /alive", handlers.HealthHandler)

	// Static file serving
	staticDir := "frontend/build"
	// SPA fallback handler
	mux.HandleFunc("/", handlers.SPAHandler(staticDir))

	// Wrap with middleware
	handler := middleware.AuthMiddleware(mux)

	log.Println("Starting server on :3000")
	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatal(err)
	}
}
