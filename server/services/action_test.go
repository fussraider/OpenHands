package services

import (
	"context"
	"openhands-go/server/events"
	"openhands-go/server/models"
	"testing"
	"time"
)

func TestActionService_ExecuteAction_Message(t *testing.T) {
	// Setup
	// ConversationStore is not used for message action, passing nil
	// RuntimeManager is not used for message action, passing nil
	// Wait, RuntimeManager.NewRuntimeManager creates runtimes map, but ExecuteAction calls NewActionService with RM.
	// NewActionService uses RM.
	// ExecuteAction logic:
	// if req.Action == "message" { ... return }
	// So RM is not accessed.

	rm := NewRuntimeManager()

	// Mock broadcaster
	broadcastCh := make(chan events.Event, 1)
	broadcaster := func(cid string, evt events.Event) {
		broadcastCh <- evt
	}

	service := NewActionService(nil, rm, broadcaster)
	conversationID := "test-conv-msg"

	// Test Message Action
	req := models.ActionRequest{
		Action: "message",
		Args:   map[string]interface{}{"content": "Hello World"},
	}

	_, err := service.ExecuteAction(context.Background(), conversationID, req)
	if err != nil {
		t.Fatalf("ExecuteAction failed: %v", err)
	}

	// Verify Event Stream
	es := service.GetEventStream(conversationID)
	evts := es.GetEvents()
	if len(evts) != 1 {
		t.Errorf("Expected 1 event, got %d", len(evts))
	} else {
		if evts[0].Source != "user" {
			t.Errorf("Expected source 'user', got '%s'", evts[0].Source)
		}

		content, ok := evts[0].Content.(models.ActionRequest)
		if !ok {
			t.Errorf("Expected content to be ActionRequest, got %T", evts[0].Content)
		} else {
			if content.Args["content"] != "Hello World" {
				t.Errorf("Expected args 'Hello World', got '%s'", content.Args["content"])
			}
		}
	}

	// Verify Broadcast
	select {
	case evt := <-broadcastCh:
		if evt.Source != "user" {
			t.Errorf("Broadcast event expected source 'user', got '%s'", evt.Source)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout waiting for broadcast")
	}
}
