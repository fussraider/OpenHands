package agent

import (
	"context"
	"openhands-go/server/config"
	"openhands-go/server/events"
	"openhands-go/server/llm"
	"openhands-go/server/runtime"
)

import (
	"openhands-go/server/models"
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

// Step overrides the base agent step to immediately return a finish action
// mirroring the Python DummyAgent behavior of bypassing the LLM entirely.
func (a *DummyAgent) Step(ctx context.Context) (interface{}, error) {
	finishAction := models.AgentFinishAction{
		Action:  models.ActionTypeAgentFinish,
		Outputs: map[string]string{"response": "Dummy agent finished immediately."},
	}

	return finishAction, nil
}

func init() {
	RegisterAgent("DummyAgent")
}
