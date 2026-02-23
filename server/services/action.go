package services

import (
	"context"
	"fmt"
	"openhands-go/server/models"
	"openhands-go/server/runtime"
	"openhands-go/server/store"
)

type ActionService struct {
	conversationStore *store.ConversationStore
	runtimes          map[string]runtime.Runtime
}

func NewActionService(cs *store.ConversationStore) *ActionService {
	return &ActionService{
		conversationStore: cs,
		runtimes:          make(map[string]runtime.Runtime),
	}
}

func (s *ActionService) ExecuteAction(ctx context.Context, conversationID string, req models.ActionRequest) (string, error) {
	// Simple implementation: if action is "run", start a process
	// For now, let's just use the LocalRuntime to run a command and return output

	if req.Action != "run" {
		return "", fmt.Errorf("unsupported action: %s", req.Action)
	}

	rt := runtime.NewLocalRuntime()
	// In a real system, we'd keep the runtime alive in s.runtimes[conversationID]

	err := rt.Start(ctx, "bash", "-c", req.Args)
	if err != nil {
		return "", err
	}
	defer rt.Close()

	// Read output (blocking for now, just for POC)
	buf := make([]byte, 1024)
	n, err := rt.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}
