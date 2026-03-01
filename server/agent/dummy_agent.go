package agent

import (
	"context"
	"openhands-go/server/config"
	"openhands-go/server/events"
	"openhands-go/server/llm"
	"openhands-go/server/runtime"

)

// DummyAgent is a minimal agent primarily used for testing or dry-runs without LLM interaction.
type DummyAgent struct {
	*Agent
}

func NewDummyAgent(id, conversationID string, llmService *llm.LLMService, rt runtime.Runtime, es *events.EventStream, delegator Delegator, cfg *config.Config) *DummyAgent {
	baseAgent := NewAgent(id, conversationID, llmService, rt, es, delegator, cfg)

	baseAgent.SystemPrompt = "You are a dummy agent that immediately finishes any task given."

	return &DummyAgent{
		Agent: baseAgent,
	}
}

func (a *DummyAgent) InitPlugins(ctx context.Context) error {
	return nil
}

// In the Python backend, the DummyAgent skips the LLM loop entirely
// and immediately emits an AgentFinishAction with a canned response.
// Here we rely on the system prompt to guide the LLM to finish, but a truer
// port would intercept the Step() function to skip the LLM call entirely.

func init() {
	RegisterAgent("DummyAgent")
}
