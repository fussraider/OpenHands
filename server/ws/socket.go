package ws

import (
	"log"
	"openhands-go/server/events"
	"openhands-go/server/models"

	socketio "github.com/googollee/go-socket.io"
)

var Server *socketio.Server

func InitSocketServer(onAction func(string, models.ActionRequest) error) error {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "join_conversation", func(s socketio.Conn, msg string) {
		log.Println("join_conversation:", msg)
		s.SetContext(msg) // Store conversation ID
		s.Join("room:" + msg)
	})

	server.OnEvent("/", "oh_user_action", func(s socketio.Conn, msg models.ActionRequest) {
		log.Printf("oh_user_action: %+v", msg)
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
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	go server.Serve()

	Server = server
	return nil
}

func BroadcastEvent(conversationID string, event events.Event) {
	if Server != nil {
		// Frontend expects "oh_event"
		Server.BroadcastToRoom("/", "room:"+conversationID, "oh_event", event)
	}
}
