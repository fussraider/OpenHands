package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"openhands-go/server/agent"
	"openhands-go/server/config"
	"openhands-go/server/events"
	"openhands-go/server/llm"
	"openhands-go/server/models"
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
		// Mirrors python: logger.debug(f'Could not get runtime status for {conversation_id}: {e}')
		slog.Debug("Could not get runtime status", "conversation_id", conversationID, "error", "runtime not found")
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

	slog.Debug("Initializing runtime now...", "runtime", config.AppConfig.Sandbox.Runtime, "conversation_id", conversationID)

	// Check config to decide which runtime to use
	if config.AppConfig.Sandbox.Runtime == "docker" {
		rt, err = runtime.NewDockerRuntime(config.AppConfig)
	} else {
		rt = runtime.NewLocalRuntime()
	}

	if err != nil {
		slog.Debug("Failed to initialize runtime", "error", err)
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

	// Mirrors python: logger.debug('Attaching to session ...') or 'Restored session ...'
	slog.Debug("Attaching to session ...", "conversation_id", conversationID)

	llmService, err := llm.NewLLMService(config.AppConfig.LLM)
	if err != nil {
		return err
	}

	if len(es.GetEvents()) > 0 {
		slog.Debug("Restored state from session", "conversation_id", conversationID, "events", len(es.GetEvents()))
	} else {
		slog.Debug("No events found, no state to restore", "conversation_id", conversationID)
	}

	// Pass rm as Delegator
	ag := agent.NewAgent("default-agent", conversationID, llmService, rt, es, rm, config.AppConfig)

	rm.agents[conversationID] = ag

	// Start loop in background
	go ag.RunLoop(ctx)

	return nil
}

func (rm *RuntimeManager) StopRuntime(conversationID string) error {
	slog.Debug("Waiting for initialization to finish before closing session", "sid", conversationID)
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

// Delegate implements agent.Delegator
func (rm *RuntimeManager) Delegate(ctx context.Context, agentName string, inputs map[string]interface{}) (map[string]interface{}, error) {
	// 1. Determine Parent Context
	// We don't have explicit parent ID passed here, but usually delegation happens within an existing conversation context.
	// However, we are reusing the RuntimeManager which tracks runtimes by conversationID.
	// We need a unique ID for the sub-agent session or reuse the runtime.
	// Let's assume we reuse the runtime of the *caller*?
	// But `Delegate` doesn't know the caller conversationID directly unless passed in ctx or we change signature.
	// Wait, `Delegate` is method of `RuntimeManager` but called by `Agent`. `Agent` has `Delegator`.
	// The `Agent` calling this is bound to a `ConversationID`.
	// We should probably pass the parent conversation ID or `Runtime` to `Delegate`?
	// The interface is `Delegate(ctx, agentName, inputs)`.
	// We can put conversationID in ctx? Or `RuntimeManager` needs to know which agent called it?
	// Actually `agent.NewAgent` passes `rm` as delegator.

	// Simplification: Assume we create a TEMPORARY sub-conversation ID sharing the SAME runtime?
	// Or just use the inputs to drive a new agent with a new EventStream but SAME runtime instance.

	// Issue: `GetRuntime` takes `conversationID`.
	// If we want to share runtime, we need the parent `conversationID`.
	// Let's update `Delegator` interface or assume `inputs` contains it? No.
	// Best practice: Pass `conversationID` in Context.

	parentConversationID, ok := ctx.Value("conversation_id").(string)
	if !ok {
		// Fallback: This implementation of Delegator might need the ID stored in the struct if created per-agent.
		// But RuntimeManager is a singleton-like service.
		// Let's rely on Context for now.
		return nil, fmt.Errorf("conversation_id missing from context")
	}

	rt, err := rm.GetRuntime(parentConversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get runtime for delegation: %w", err)
	}

	// 2. Create Isolated Event Stream for Sub-Task
	subConversationID := fmt.Sprintf("%s-sub-%s", parentConversationID, agentName)
	// We don't want to persist this necessarily, or maybe we do for debugging.
	// Let's use in-memory stream or a separate file.
	es := events.NewEventStream(subConversationID, "") // Empty path = in-memory

	// 3. Initialize Sub-Agent
	llmService, err := llm.NewLLMService(config.AppConfig.LLM)
	if err != nil {
		return nil, err
	}

	// Create inputs message
	taskDescription := fmt.Sprintf("You are a delegated agent working on: %s. Inputs: %v", agentName, inputs)
	es.AddEvent(events.Event{
		ID: "init",
		Type: events.EventTypeAction,
		Content: models.MessageAction{
			Action: models.ActionTypeMessage,
			Content: taskDescription,
		},
		Source: "user",
	})

	// Use specific agent config if available (e.g. BrowsingAgent)
	// For now, we default to NewAgent which is CodeAct.
	// If agentName == "BrowsingAgent", we could use NewBrowsingAgent.
	subAgent := agent.NewAgent(agentName, subConversationID, llmService, rt, es, rm, config.AppConfig)

	// 4. Run Sub-Agent
	// This blocks until finish
	outputs, err := subAgent.RunUntilDone(ctx)

	// Convert map[string]string to map[string]interface{}
	result := make(map[string]interface{})
	for k, v := range outputs {
		result[k] = v
	}

	return result, err
}
