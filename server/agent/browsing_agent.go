package agent

import (
	"context"
	"openhands-go/server/config"
	"openhands-go/server/events"
	"openhands-go/server/llm"
	"openhands-go/server/runtime"
	"openhands-go/server/runtime/plugins/browser"

	"github.com/tmc/langchaingo/llms"
)

type BrowsingAgent struct {
	*Agent // Embed base agent
}

func NewBrowsingAgent(id, conversationID string, llmService *llm.LLMService, rt runtime.Runtime, es *events.EventStream, delegator Delegator, cfg *config.Config) *BrowsingAgent {
	// 1. Initialize Base Agent
	baseAgent := NewAgent(id, conversationID, llmService, rt, es, delegator, cfg)

	// 2. Override Plugins (Only Browser needed for BrowsingAgent?)
	// Or maybe it needs Python too if it does calculations.
	// For now, let's ensure BrowserPlugin is present.
	// (NewAgent already adds BrowserPlugin)

	// 3. Override System Prompt
	// Browsing Agent has a specific prompt.
	// We should probably allow `NewAgent` to accept a prompt or `AgentType`.
	// For this port, we set it manually.
	baseAgent.SystemPrompt = `
You are a browsing agent. Your goal is to navigate the web to find information or perform actions.
You have access to a browser tool. Use it to navigate, click, type, and read pages.
When you have found the answer or completed the task, call finish.
`

	// 4. Filter Tools?
	// BrowsingAgent primarily uses `browser_navigate`, `browser_click`, `browser_type`, etc.
	// CodeAct agent has `execute_bash`.
	// For parity, we might restrict tools, but CodeAct with Browser Plugin is powerful enough.
	// We'll keep all tools for now but emphasize browsing in the prompt.

	return &BrowsingAgent{
		Agent: baseAgent,
	}
}

// Ensure it implements the Agent loop (RunLoop is inherited)
// But we might want specific Step logic if the Python one differs significantly.
// Python BrowsingAgent uses `BrowserInteractiveAction` which maps to browser gym commands.
// Our `BrowserPlugin` in `server/runtime/plugins/browser` exposes `browser_action` tool.

// So CodeAct logic (LLM -> Tool Call -> Execute) works fine here too.
// The key is the Prompt and Tool Definitions.

func (a *BrowsingAgent) InitPlugins(ctx context.Context) error {
	// Ensure we have browser plugin
	hasBrowser := false
	for _, p := range a.Plugins {
		if p.Name() == "browser" {
			hasBrowser = true
			break
		}
	}
	if !hasBrowser {
		a.Plugins = append(a.Plugins, browser.NewBrowserPlugin())
	}
	// Re-register tools
	a.Tools = []llms.Tool{}
	// Add default tools (bash, finish, delegate) - maybe remove bash for pure browsing agent?
	// For now, keep defaults from NewAgent logic or reconstruct.
	// Since we called NewAgent, a.Tools is populated.
	// We can filter if needed.
	return nil
}
