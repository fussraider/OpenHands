package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"openhands-go/server/config"
	"openhands-go/server/handlers"
	"openhands-go/server/logger"
	"openhands-go/server/middleware"
	"openhands-go/server/ws"
)

func main() {
	logger.Init()

	if err := config.LoadConfig(); err != nil {
		slog.Error("Failed to load config", "error", err)
		panic(err)
	}

	handlers.InitHandlers()

	if err := ws.InitSocketServer(handlers.ProcessSocketAction); err != nil {
		slog.Error("Failed to init socket server", "error", err)
		panic(err)
	}

	mux := http.NewServeMux()

	// Socket.IO
	mux.Handle("/socket.io/", ws.Server)

	// API routes
	mux.HandleFunc("GET /api/options/models", handlers.ModelsHandler)
	mux.HandleFunc("GET /api/options/agents", handlers.GetAgentsHandler)
	mux.HandleFunc("GET /api/options/security-analyzers", handlers.GetSecurityAnalyzersHandler)

	mux.HandleFunc("GET /api/conversations", handlers.SearchConversationsHandler)
	mux.HandleFunc("POST /api/conversations", handlers.NewConversationHandler)
	mux.HandleFunc("GET /api/conversations/{id}", handlers.GetConversationHandler)
	mux.HandleFunc("POST /api/conversations/{id}/action", handlers.ExecuteActionHandler) // New route

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

	// Github / Git Provider API
	mux.HandleFunc("GET /api/user/installations", handlers.GetUserInstallationsHandler)
	mux.HandleFunc("GET /api/user/repositories", handlers.GetUserRepositoriesHandler)
	mux.HandleFunc("GET /api/user/info", handlers.GetUserInfoHandler)
	mux.HandleFunc("GET /api/user/search/repositories", handlers.SearchRepositoriesHandler)
	mux.HandleFunc("GET /api/user/search/branches", handlers.SearchBranchesHandler)
	mux.HandleFunc("GET /api/user/repository/branches", handlers.GetRepositoryBranchesHandler)

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

	host := config.AppConfig.Server.Host
	if host == "" {
		host = "localhost"
	}
	addr := fmt.Sprintf("%s:%d", host, config.AppConfig.Server.Port)

	slog.Info("Starting server", "address", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		slog.Error("Server failed", "error", err)
		panic(err)
	}
}
