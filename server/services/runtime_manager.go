package services

import (
	"context"
	"errors"
	"openhands-go/server/runtime"
	"sync"
)

type RuntimeManager struct {
	mu       sync.RWMutex
	runtimes map[string]runtime.Runtime
}

func NewRuntimeManager() *RuntimeManager {
	return &RuntimeManager{
		runtimes: make(map[string]runtime.Runtime),
	}
}

func (rm *RuntimeManager) GetRuntime(conversationID string) (runtime.Runtime, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	rt, ok := rm.runtimes[conversationID]
	if !ok {
		return nil, errors.New("runtime not found")
	}
	return rt, nil
}

func (rm *RuntimeManager) CreateRuntime(ctx context.Context, conversationID string) (runtime.Runtime, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if _, ok := rm.runtimes[conversationID]; ok {
		return nil, errors.New("runtime already exists")
	}

	// For now, only LocalRuntime is supported
	rt := runtime.NewLocalRuntime()
	// Initialize/Start logic if needed here, but LocalRuntime starts on command execution currently.
	// Ideally, Start() should be called here if it was a container.

	rm.runtimes[conversationID] = rt
	return rt, nil
}

func (rm *RuntimeManager) StopRuntime(conversationID string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rt, ok := rm.runtimes[conversationID]
	if !ok {
		return nil // Already stopped
	}

	if err := rt.Close(); err != nil {
		return err
	}

	delete(rm.runtimes, conversationID)
	return nil
}
