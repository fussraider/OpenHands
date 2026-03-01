package agent

import (
	"context"
	"openhands-go/server/config"
	"openhands-go/server/events"
	"openhands-go/server/llm"
	"openhands-go/server/runtime"

)

type ReadOnlyAgent struct {
	*Agent
}

func NewReadOnlyAgent(id, conversationID string, llmService *llm.LLMService, rt runtime.Runtime, es *events.EventStream, delegator Delegator, cfg *config.Config) *ReadOnlyAgent {
	baseAgent := NewAgent(id, conversationID, llmService, rt, es, delegator, cfg)

	// Override System Prompt to strictly forbid write operations
	baseAgent.SystemPrompt = `
You are a read-only agent. Your goal is to explore, read, and analyze the workspace, but you must NEVER modify, delete, or create files.
You can read files and execute read-only bash commands like 'ls', 'cat', 'grep', 'find'.
When you have found the answer or completed the analysis, call finish.
`
	return &ReadOnlyAgent{
		Agent: baseAgent,
	}
}

func (a *ReadOnlyAgent) InitPlugins(ctx context.Context) error {
	// Keep tools standard, rely on security analyzer or prompt for read-only constraint.
	return nil
}

func init() {
	RegisterAgent("ReadOnlyAgent")
}
