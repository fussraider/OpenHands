package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"openhands-go/server/models"
)

// GetMicroagentManagementConversationsHandler returns conversations filtered for microagent management
func GetMicroagentManagementConversationsHandler(w http.ResponseWriter, r *http.Request) {
	// In Python this tries app_conversation_service and logs if unavailable.
	// We only have the ConversationStore for now.
	slog.Debug("V1 conversation service not available", "error", "not implemented in MVP, falling back to ConversationStore")
	selectedRepository := r.URL.Query().Get("selected_repository")

	// Default to returning all conversations for now if no repo provided,
	// though Python expects selected_repository filter.

	allConversations := ConversationStore.ListConversations()
	filtered := []models.ConversationInfo{}

	for _, conv := range allConversations {
		// Filter by trigger = MICROAGENT_MANAGEMENT
		if conv.Trigger != "MICROAGENT_MANAGEMENT" && conv.Trigger != "microagent_management" {
			continue
		}

		// Filter by selected repository if provided
		if selectedRepository != "" && conv.SelectedRepository != selectedRepository {
			continue
		}
		populateRuntimeStatus(&conv)

		filtered = append(filtered, conv)
	}

	w.Header().Set("Content-Type", "application/json")

	// Python returns ConversationInfoResultSet with results and next_page_id
	response := map[string]interface{}{
		"results": filtered,
		"next_page_id": nil,
	}

	json.NewEncoder(w).Encode(response)
}
