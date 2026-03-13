package ws

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"openhands-go/server/events"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// V1Client represents a connected V1 WebSocket client.
type V1Client struct {
	conn           *websocket.Conn
	conversationID string
	mu             sync.Mutex
}

// WriteJSON sends a JSON message to the client, thread-safe.
func (c *V1Client) WriteJSON(v interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(v)
}

// V1WebSocketHandler handles raw WebSocket connections for V1 protocol.
// Route: /sockets/events/{id}
// Query params: resend_all (bool), session_api_key (string)
func V1WebSocketHandler(getEventStream func(string) *events.EventStream, onMessage func(string, json.RawMessage)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conversationID := r.PathValue("id")
		if conversationID == "" {
			http.Error(w, "conversation id required", http.StatusBadRequest)
			return
		}

		resendAll := r.URL.Query().Get("resend_all") == "true"
		// session_api_key could be used for auth in the future
		// sessionAPIKey := r.URL.Query().Get("session_api_key")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("WebSocket upgrade failed", "error", err, "conversation_id", conversationID)
			return
		}

		client := &V1Client{
			conn:           conn,
			conversationID: conversationID,
		}

		slog.Debug("V1 WebSocket connected", "conversation_id", conversationID, "resend_all", resendAll)

		es := getEventStream(conversationID)

		// Replay existing events if requested
		if resendAll && es != nil {
			existingEvents := es.GetEvents()
			for _, event := range existingEvents {
				if err := client.WriteJSON(event); err != nil {
					slog.Error("Failed to send replay event", "error", err, "conversation_id", conversationID)
					conn.Close()
					return
				}
			}
			slog.Debug("Replayed events", "count", len(existingEvents), "conversation_id", conversationID)
		}

		// Subscribe to new events
		var unsubscribe func()
		if es != nil {
			unsubscribe = es.Subscribe(func(event events.Event) {
				if err := client.WriteJSON(event); err != nil {
					slog.Debug("Failed to send event to client, closing", "error", err, "conversation_id", conversationID)
					conn.Close()
				}
			})
		}

		// Read loop: handle incoming messages from client
		go func() {
			defer func() {
				if unsubscribe != nil {
					unsubscribe()
				}
				conn.Close()
				slog.Debug("V1 WebSocket disconnected", "conversation_id", conversationID)
			}()

			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
						slog.Debug("V1 WebSocket read error", "error", err, "conversation_id", conversationID)
					}
					return
				}

				slog.Debug("V1 WebSocket received message", "conversation_id", conversationID)

				if onMessage != nil {
					onMessage(conversationID, json.RawMessage(message))
				}
			}
		}()
	}
}
