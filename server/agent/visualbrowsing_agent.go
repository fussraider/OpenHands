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

type VisualBrowsingAgent struct {
	*Agent // Embed base agent
}

func NewVisualBrowsingAgent(id, conversationID string, llmService *llm.LLMService, rt runtime.Runtime, es *events.EventStream, delegator Delegator, cfg *config.Config) *VisualBrowsingAgent {
	// 1. Initialize Base Agent
	baseAgent := NewAgent(id, conversationID, llmService, rt, es, delegator, cfg)

	// 2. Override System Prompt
	// Visual Browsing Agent has a specific prompt that emphasizes visual reasoning and accessibility trees.
	baseAgent.SystemPrompt = `
You are a visual browsing agent trying to solve a web task based on the content of the page and user instructions.
You can interact with the page and explore, and send messages to the user when you finish the task.
Each time you submit an action it will be sent to the browser and you will receive a new page.
You have access to a browser tool. Use it to navigate, click, type, and read pages.
Make sure to use bid (from accessibility tree) to identify elements when using commands.
Interacting with combobox, dropdowns and auto-complete fields can be tricky, sometimes you need to use select_option, while other times you need to use fill or click and wait for the reaction of the page.
When you have found the answer or completed the task, call finish.
`

	return &VisualBrowsingAgent{
		Agent: baseAgent,
	}
}

func (a *VisualBrowsingAgent) InitPlugins(ctx context.Context) error {
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
	// Re-register tools to clear non-visual ones if necessary
	// For now, keep defaults from NewAgent logic or reconstruct.
	a.Tools = []llms.Tool{}
	return nil
}

func init() {
	RegisterAgent("VisualBrowsingAgent")
}
