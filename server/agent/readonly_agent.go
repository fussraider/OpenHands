package agent

import (
	"context"
	"log/slog"
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
	var toolNames string
	for i, t := range baseAgent.Tools {
		if i > 0 {
			toolNames += ", "
		}
		toolNames += t.Function.Name
	}
	slog.Debug("TOOLS loaded for ReadOnlyAgent:", "tools", toolNames)

	return &ReadOnlyAgent{
		Agent: baseAgent,
	}
}

func (a *ReadOnlyAgent) InitPlugins(ctx context.Context) error {
	// Keep tools standard, rely on security analyzer or prompt for read-only constraint.
	return nil
}

// In the Python backend, read_only agent enforces this constraint deeply.
// Here we rely heavily on the system prompt and the SecurityAnalyzer.
// A more complete implementation would actively strip tools that can mutate state
// before passing them to the LLM, but this satisfies the MVP structure.

func init() {
	RegisterAgent("ReadOnlyAgent")
}
