package services

import (
	"context"
	"errors"
	"openhands-go/server/config"
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

	var rt runtime.Runtime
	var err error

	// Check config to decide which runtime to use
	if config.AppConfig.Sandbox.Runtime == "docker" {
		rt, err = runtime.NewDockerRuntime(config.AppConfig)
	} else {
		rt = runtime.NewLocalRuntime()
	}

	if err != nil {
		return nil, err
	}

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
