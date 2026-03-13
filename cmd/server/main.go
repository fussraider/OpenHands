package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"openhands-go/server/config"
	"openhands-go/server/events"
	"openhands-go/server/handlers"
	"openhands-go/server/logger"
	"openhands-go/server/middleware"
	"openhands-go/server/models"
	"openhands-go/server/observability"
	"openhands-go/server/services"
	"openhands-go/server/ws"

	"github.com/google/uuid"
)

func main() {
	logger.Init()

	if err := config.LoadConfig(); err != nil {
		slog.Error("Failed to load config", "error", err)
		panic(err)
	}

	// Initialize Tracing
	if err := observability.InitTracer(); err != nil {
		slog.Error("Failed to init tracer", "error", err)
		// Don't panic, continue without tracing
	}
	defer observability.Shutdown(context.Background())

	handlers.InitHandlers()

	mux := http.NewServeMux()

	// V1 WebSocket: /sockets/events/{id}
	mux.HandleFunc("/sockets/events/{id}", ws.V1WebSocketHandler(
		func(conversationID string) *events.EventStream {
			if handlers.ActionService != nil {
				return handlers.ActionService.GetEventStream(conversationID)
			}
			return nil
		},
		func(conversationID string, message json.RawMessage) {
			// Parse V1 SendMessageRequest and convert to action
			var msg struct {
				Role    string `json:"role"`
				Content []struct {
					Type string `json:"type"`
					Text string `json:"text"`
				} `json:"content"`
			}
			if err := json.Unmarshal(message, &msg); err != nil {
				slog.Error("Failed to parse V1 message", "error", err)
				return
			}

			// Extract text content
			var textContent string
			for _, c := range msg.Content {
				if c.Type == "text" {
					textContent = c.Text
					break
				}
			}

			if textContent != "" && handlers.ActionService != nil {
				es := handlers.ActionService.GetEventStream(conversationID)
				es.AddEvent(events.Event{
					ID:     uuid.New().String(),
					Type:   events.EventTypeAction,
					Content: models.ActionRequest{
						Action: "message",
						Args:   map[string]interface{}{"content": textContent},
					},
					Source: "user",
				})
			}
		},
	))

	// API routes
	mux.HandleFunc("GET /api/options/models", handlers.ModelsHandler)
	mux.HandleFunc("GET /api/options/agents", handlers.GetAgentsHandler)
	mux.HandleFunc("GET /api/options/security-analyzers", handlers.GetSecurityAnalyzersHandler)

	mux.HandleFunc("GET /api/microagent-management/conversations", handlers.GetMicroagentManagementConversationsHandler)

	mux.HandleFunc("GET /api/conversations", handlers.SearchConversationsHandler)
	mux.HandleFunc("POST /api/conversations", handlers.NewConversationHandler)
	mux.HandleFunc("GET /api/conversations/{id}", handlers.GetConversationHandler)
	mux.HandleFunc("PATCH /api/conversations/{id}", handlers.UpdateConversationHandler)
	mux.HandleFunc("DELETE /api/conversations/{id}", handlers.DeleteConversationHandler)
	mux.HandleFunc("POST /api/conversations/{id}/start", handlers.StartConversationHandler)
	mux.HandleFunc("POST /api/conversations/{id}/stop", handlers.StopConversationHandler)
	mux.HandleFunc("POST /api/conversations/{id}/action", handlers.ExecuteActionHandler) // New route
	mux.HandleFunc("POST /api/conversations/{id}/message", handlers.AddMessageHandler)

	mux.HandleFunc("GET /api/conversations/{id}/list-files", handlers.ListFilesHandler)
	mux.HandleFunc("GET /api/conversations/{id}/select-file", handlers.SelectFileHandler)
	mux.HandleFunc("POST /api/conversations/{id}/upload-files", handlers.UploadFilesHandler)
	mux.HandleFunc("GET /api/conversations/{id}/zip-workspace", handlers.ZipWorkspaceHandler)
	mux.HandleFunc("GET /api/conversations/{id}/microagents", handlers.GetConversationMicroagentsHandler)
	mux.HandleFunc("GET /api/conversations/{id}/remember-prompt", handlers.GetRememberPromptHandler)
	mux.HandleFunc("GET /api/conversations/{id}/vscode-url", handlers.GetVSCodeURLHandler)
	mux.HandleFunc("GET /api/conversations/{id}/web-hosts", handlers.GetWebHostsHandler)
	mux.HandleFunc("GET /api/conversations/{id}/config", handlers.GetConversationConfigHandler)
	mux.HandleFunc("GET /api/conversations/{id}/events", handlers.GetConversationEventsHandler)
	mux.HandleFunc("POST /api/conversations/{id}/events", handlers.GetConversationEventsHandler)
	mux.HandleFunc("GET /api/conversations/{id}/events/count", handlers.GetConversationEventsCountHandler)
	mux.HandleFunc("POST /api/conversations/{id}/exp-config", handlers.ExpConfigHandler)
	mux.HandleFunc("GET /api/config", handlers.GetSettingsHandler) // Backwards compatible global config

	mux.HandleFunc("GET /api/conversations/{id}/trajectory", handlers.GetTrajectoryHandler)
	mux.HandleFunc("POST /api/conversations/{id}/submit-feedback", handlers.SubmitFeedbackHandler)

	// Security API
	mux.HandleFunc("/api/conversations/{id}/security/{path...}", handlers.SecurityAPIHandler)

	mux.HandleFunc("GET /api/settings", handlers.GetSettingsHandler)
	mux.HandleFunc("POST /api/settings", handlers.StoreSettingsHandler)

	mux.HandleFunc("GET /api/secrets", handlers.GetSecretsHandler)
	mux.HandleFunc("POST /api/secrets", handlers.StoreSecretHandler)
	mux.HandleFunc("PUT /api/secrets/{id}", handlers.UpdateSecretHandler)
	mux.HandleFunc("DELETE /api/secrets/{id}", handlers.DeleteSecretHandler)
	mux.HandleFunc("POST /api/add-git-providers", handlers.StoreGitProvidersHandler)
	mux.HandleFunc("POST /api/unset-provider-tokens", handlers.UnsetGitProvidersHandler)

	// Github / Git Provider API
	mux.HandleFunc("GET /api/user/installations", handlers.GetUserInstallationsHandler)
	mux.HandleFunc("GET /api/user/repositories", handlers.GetUserRepositoriesHandler)
	mux.HandleFunc("GET /api/user/info", handlers.GetUserInfoHandler)
	mux.HandleFunc("GET /api/user/search/repositories", handlers.SearchRepositoriesHandler)
	mux.HandleFunc("GET /api/user/search/branches", handlers.SearchBranchesHandler)
	mux.HandleFunc("GET /api/user/repository/branches", handlers.GetRepositoryBranchesHandler)
	mux.HandleFunc("GET /api/user/suggested-tasks", handlers.GetSuggestedTasksHandler)
	// Go 1.22+ supports wildcard {name...}, but it must be at the end of the path.
	// We use {owner}/{repo} to avoid the panic.
	mux.HandleFunc("GET /api/user/repository/{owner}/{repo}/microagents", handlers.GetRepositoryMicroagentsHandler)
	mux.HandleFunc("GET /api/user/repository/{owner}/{repo}/microagents/content", handlers.GetRepositoryMicroagentContentHandler)

	// MCP Mount
	mux.HandleFunc("/mcp/", handlers.MCPSSEHandler)

	handlers.SetRuntimeManager(services.NewRuntimeManager())
	handlers.RegisterV1Routes(mux)

	mux.HandleFunc("GET /health", handlers.HealthHandler)
	mux.HandleFunc("GET /alive", handlers.HealthHandler)
	mux.HandleFunc("GET /ready", handlers.ReadyHandler)
	mux.HandleFunc("GET /server_info", handlers.ServerInfoHandler)

	// Static file serving
	staticDir := "frontend/build"
	// SPA fallback handler
	mux.HandleFunc("/", handlers.SPAHandler(staticDir))

	// Wrap with middleware
	rateLimiter := middleware.NewRateLimiter(2, 1, 1) // 2 req/s, sleep 1s on burst

	handler := middleware.AuthMiddleware(mux)
	handler = middleware.CacheControlMiddleware(handler)
	handler = middleware.RateLimitMiddleware(rateLimiter)(handler)

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
