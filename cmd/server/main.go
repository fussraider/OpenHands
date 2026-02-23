package main

import (
	"log"
	"net/http"
	"openhands-go/server/handlers"
)

func main() {
	// API routes
	http.HandleFunc("GET /api/options/models", handlers.ModelsHandler)
	http.HandleFunc("GET /api/options/agents", handlers.GetAgentsHandler)
	http.HandleFunc("GET /api/options/security-analyzers", handlers.GetSecurityAnalyzersHandler)

	http.HandleFunc("GET /api/conversations", handlers.SearchConversationsHandler)
	http.HandleFunc("POST /api/conversations", handlers.NewConversationHandler)
	http.HandleFunc("GET /api/conversations/{id}", handlers.GetConversationHandler)

	http.HandleFunc("GET /api/conversations/{id}/list-files", handlers.ListFilesHandler)
	http.HandleFunc("GET /api/conversations/{id}/select-file", handlers.SelectFileHandler)

	http.HandleFunc("GET /api/conversations/{id}/trajectory", handlers.GetTrajectoryHandler)
	http.HandleFunc("POST /api/conversations/{id}/submit-feedback", handlers.SubmitFeedbackHandler)

	// Security API
	http.HandleFunc("/api/conversations/{id}/security/{path...}", handlers.SecurityAPIHandler)

	http.HandleFunc("GET /api/settings", handlers.GetSettingsHandler)
	http.HandleFunc("POST /api/settings", handlers.StoreSettingsHandler)

	http.HandleFunc("GET /api/secrets", handlers.GetSecretsHandler)
	http.HandleFunc("POST /api/secrets", handlers.StoreSecretHandler)
	http.HandleFunc("DELETE /api/secrets/{key}", handlers.DeleteSecretHandler)

	http.HandleFunc("GET /api/github/repositories", handlers.RepositoriesHandler)

	// MCP Mount
	http.HandleFunc("/mcp/", handlers.MCPSSEHandler)

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
