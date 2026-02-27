package services

import (
	"context"
	"fmt"
	"openhands-go/server/config"
	"openhands-go/server/events"
	"openhands-go/server/models"
	"openhands-go/server/store"
	"path/filepath"

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
		// Determine file path for persistence
		var filePath string
		if config.AppConfig != nil && config.AppConfig.FileStorePath != "" {
			filePath = filepath.Join(config.AppConfig.FileStorePath, "sessions", conversationID, "events.jsonl")
		}

		es := events.NewEventStream(conversationID, filePath)
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
	// Use persistent Execute
	output, exitCode, err := rt.Execute(ctx, "bash", "-c", req.Args)
	if err != nil {
		return "", err
	}
	// Note: Do not Close runtime here, it is persistent (managed by RuntimeManager).

	// 4. Add Observation to EventStream
	obs := models.CmdOutputObservation{
		Observation: "run",
		Content:     output,
		Metadata: models.CmdOutputMetadata{
			ExitCode: exitCode,
		},
		Command: req.Args,
	}

	es.AddEvent(events.Event{
		ID:      uuid.New().String(),
		Type:    events.EventTypeObservation,
		Content: obs,
		Source:  "runtime",
	})

	return output, nil
}
