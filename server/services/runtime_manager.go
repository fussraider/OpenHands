package services

import (
	"context"
	"errors"
	"openhands-go/server/agent"
	"openhands-go/server/config"
	"openhands-go/server/events"
	"openhands-go/server/llm"
	"openhands-go/server/runtime"
	"sync"
)

type RuntimeManager struct {
	mu       sync.RWMutex
	runtimes map[string]runtime.Runtime
	agents   map[string]*agent.Agent
}

func NewRuntimeManager() *RuntimeManager {
	return &RuntimeManager{
		runtimes: make(map[string]runtime.Runtime),
		agents:   make(map[string]*agent.Agent),
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

func (rm *RuntimeManager) StartAgent(ctx context.Context, conversationID string, es *events.EventStream) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if _, ok := rm.agents[conversationID]; ok {
		return nil // Already started
	}

	rt, ok := rm.runtimes[conversationID]
	if !ok {
		return errors.New("runtime must be created before starting agent")
	}

	llmService, err := llm.NewLLMService(config.AppConfig.LLM)
	if err != nil {
		return err
	}
	ag := agent.NewAgent("default-agent", conversationID, llmService, rt, es)

	rm.agents[conversationID] = ag

	// Start loop in background
	go ag.RunLoop(ctx)

	return nil
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
	// Agents should also be stopped/cleaned up (ctx cancellation handled by caller typically, but here we might need explicit stop)
	delete(rm.agents, conversationID)

	return nil
}
