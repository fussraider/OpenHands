package ws

import (
	"log"
	"log/slog"
	"openhands-go/server/events"
	"openhands-go/server/models"

	socketio "github.com/googollee/go-socket.io"
)

var Server *socketio.Server

func InitSocketServer(onAction func(string, models.ActionRequest) error) error {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		u := s.URL()
		conversationID := u.Query().Get("conversation_id")
		if conversationID != "" {
			s.SetContext(conversationID)
			s.Join("room:" + conversationID)
			slog.Debug("WebSocket connected and joined room", "id", s.ID(), "conversation_id", conversationID)
		} else {
			s.SetContext("")
			slog.Debug("WebSocket connected without conversation_id", "id", s.ID())
		}

		// Mimic python `Using client wait timeout...`
		slog.Debug("Using client wait timeout: 30s for session", "sid", s.ID())

		return nil
	})

	server.OnEvent("/", "join_conversation", func(s socketio.Conn, msg string) {
		slog.Debug("join_conversation", "msg", msg)
		s.SetContext(msg) // Store conversation ID
		s.Join("room:" + msg)
	})

	server.OnEvent("/", "oh_user_action", func(s socketio.Conn, msg models.ActionRequest) {
		slog.Debug("Received message", "action", msg.Action, "args", msg.Args)
		conversationID, ok := s.Context().(string)
		if !ok || conversationID == "" {
			log.Println("Error: oh_user_action received but no conversation joined")
			return
		}

		if onAction != nil {
			if err := onAction(conversationID, msg); err != nil {
				log.Printf("Error processing action: %v", err)
				s.Emit("error", err.Error())
			}
		}
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		slog.Error("WebSocket Error", "error", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		slog.Debug("WebSocket closed", "reason", reason)
	})

	go server.Serve()

	Server = server
	return nil
}

func BroadcastEvent(conversationID string, event events.Event) {
	if Server != nil {
		// Check if room has clients. If empty, mimic python's wait logging.
		roomName := "room:" + conversationID
		if Server.RoomLen("/", roomName) == 0 {
			slog.Debug("There is no listening client in the current room, waiting for the 1th attempt (timeout: 30s)", "sid", conversationID)
			// In MVP we just drop/queue the event depending on EventStream persistence,
			// but we port the log exactly as requested.
		}

		// Frontend expects "oh_event"
		slog.Debug("Sent message", "conversation_id", conversationID, "event_type", event.Type, "event_id", event.ID)
		slog.Debug("oh_event", "type", event.Type) // mirrors logger.debug(f'oh_event: {event.__class__.__name__}')
		Server.BroadcastToRoom("/", roomName, "oh_event", event)
	}
}
