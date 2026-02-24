package services

import (
	"context"
	"fmt"
	"openhands-go/server/events"
	"openhands-go/server/models"
	"openhands-go/server/store"

	"github.com/google/uuid"
)

type ActionService struct {
	conversationStore *store.ConversationStore
	runtimeManager    *RuntimeManager
	eventStreams      map[string]*events.EventStream
	eventBroadcaster  func(string, events.Event)
}

func NewActionService(cs *store.ConversationStore, rm *RuntimeManager, broadcaster func(string, events.Event)) *ActionService {
	return &ActionService{
		conversationStore: cs,
		runtimeManager:    rm,
		eventStreams:      make(map[string]*events.EventStream),
		eventBroadcaster:  broadcaster,
	}
}

func (s *ActionService) GetEventStream(conversationID string) *events.EventStream {
	if _, ok := s.eventStreams[conversationID]; !ok {
		es := events.NewEventStream(conversationID)
		if s.eventBroadcaster != nil {
			es.Subscribe(func(event events.Event) {
				s.eventBroadcaster(conversationID, event)
			})
		}
		s.eventStreams[conversationID] = es
	}
	return s.eventStreams[conversationID]
}

func (s *ActionService) ExecuteAction(ctx context.Context, conversationID string, req models.ActionRequest) (string, error) {
	if req.Action == "message" {
		// Just add to event stream, agent will pick it up
		es := s.GetEventStream(conversationID)
		es.AddEvent(events.Event{
			ID:      uuid.New().String(),
			Type:    events.EventTypeAction,
			Content: req, // {action: "message", args: "msg"}
			Source:  "user",
		})
		return "", nil
	}

	if req.Action != "run" {
		return "", fmt.Errorf("unsupported action: %s", req.Action)
	}

	// 1. Get or Create Runtime
	rt, err := s.runtimeManager.GetRuntime(conversationID)
	if err != nil {
		rt, err = s.runtimeManager.CreateRuntime(ctx, conversationID)
		if err != nil {
			return "", err
		}
	}

	// 2. Add Action to EventStream
	es := s.GetEventStream(conversationID)
	es.AddEvent(events.Event{
		ID:      uuid.New().String(),
		Type:    events.EventTypeAction,
		Content: req,
		Source:  "user",
	})

	// 3. Execute in Runtime
	// Note: LocalRuntime currently starts a new process for every command if used via Start().
	// But `ExecuteAction` implies a persistent session or one-off command.
	// `LocalRuntime` implementation in `runtime.go` uses `exec.Command`, which is one-off.
	// So we assume one-off execution for now.

	err = rt.Start(ctx, "bash", "-c", req.Args)
	if err != nil {
		return "", err
	}
	// For LocalRuntime, we need to handle cleanup if it's one-off
	// But `runtimeManager` keeps it. This mismatch needs fixing in future.
	// For now, we just close it after use if it's not meant to be persistent shell.
	defer rt.Close()

	// 4. Read Output
	buf := make([]byte, 1024)
	n, err := rt.Read(buf)
	output := ""
	if err == nil {
		output = string(buf[:n])
	}

	// 5. Add Observation to EventStream
	es.AddEvent(events.Event{
		ID:      uuid.New().String(),
		Type:    events.EventTypeObservation,
		Content: map[string]string{"output": output},
		Source:  "runtime",
	})

	return output, nil
}
