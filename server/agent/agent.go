package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"openhands-go/server/events"
	"openhands-go/server/llm"
	"openhands-go/server/models"
	"openhands-go/server/runtime"
	"openhands-go/server/runtime/plugins"
	"openhands-go/server/runtime/plugins/jupyter"
	"time"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/llms"
)

type Agent struct {
	ID             string
	ConversationID string
	LLM            *llm.LLMService
	Runtime        runtime.Runtime
	EventStream    *events.EventStream
	SystemPrompt   string
	Tools          []llms.Tool
	Plugins        []plugins.Plugin
}

func NewAgent(id, conversationID string, llmService *llm.LLMService, rt runtime.Runtime, es *events.EventStream) *Agent {
	// Initialize Plugins
	plugs := []plugins.Plugin{
		jupyter.NewJupyterPlugin(),
	}

	// Initialize plugins with runtime
	// We use a background context for init, or pass context?
	// NewAgent doesn't take context. We'll init lazily or just ignore context for now?
	// `Init` takes context. We should probably init in `RunLoop` or `Step` if not initialized?
	// Or change NewAgent signature.
	// For now, let's do best effort init here or in RunLoop.
	// Better in RunLoop or separate method.

	agent := &Agent{
		ID:             id,
		ConversationID: conversationID,
		LLM:            llmService,
		Runtime:        rt,
		EventStream:    es,
		SystemPrompt:   "You are a helpful coding agent using the CodeAct framework. You can execute bash commands to solve tasks. When you are done, call finish.",
		Plugins:        plugs,
	}

	// Collect tools
	tools := []llms.Tool{}

	// Default tools
	bashTool := llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "execute_bash",
			Description: "Execute a bash command in the environment.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"command": map[string]interface{}{
						"type":        "string",
						"description": "The bash command to execute.",
					},
					"thought": map[string]interface{}{
						"type":        "string",
						"description": "Your reasoning for executing this command.",
					},
				},
				"required": []string{"command"},
			},
		},
	}
	finishTool := llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "finish",
			Description: "Finish the task.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"thought": map[string]interface{}{
						"type":        "string",
						"description": "Your reasoning for finishing the task.",
					},
				},
			},
		},
	}
	tools = append(tools, bashTool, finishTool)

	// Plugin tools
	for _, p := range plugs {
		tools = append(tools, p.Tools()...)
	}

	agent.Tools = tools
	return agent
}

func (a *Agent) Step(ctx context.Context) error {
	// 1. Get History & Construct Messages
	history := a.EventStream.GetEvents()
	messages := a.eventsToMessages(history)

	// 2. LLM Completion
	resp, err := a.LLM.CompleteWithTools(ctx, messages, a.Tools)
	if err != nil {
		return fmt.Errorf("LLM completion error: %w", err)
	}

	// 3. Handle Tool Calls
	if len(resp.ToolCalls) > 0 {
		for _, tc := range resp.ToolCalls {
			handled := false

			// Built-in tools
			switch tc.FunctionCall.Name {
			case "execute_bash":
				handled = true
				var args struct {
					Command string `json:"command"`
					Thought string `json:"thought"`
				}
				if err := json.Unmarshal([]byte(tc.FunctionCall.Arguments), &args); err != nil {
					log.Printf("Error unmarshalling execute_bash args: %v", err)
					// Record error output?
					a.recordObservation(tc.ID, fmt.Sprintf("Error unmarshalling args: %v", err), "run")
					continue
				}

				a.recordAction(models.CmdRunAction{
					Action:     models.ActionTypeCmdRun,
					Command:    args.Command,
					Thought:    args.Thought,
					ToolCallID: tc.ID,
				})

				output, exitCode, err := a.Runtime.Execute(ctx, "bash", "-c", args.Command)
				content := output
				if err != nil {
					content = fmt.Sprintf("Error executing command: %v", err)
				}

				a.recordObservation(tc.ID, content, "run", exitCode, args.Command)

			case "finish":
				handled = true
				a.EventStream.AddEvent(events.Event{
					ID:      uuid.New().String(),
					Type:    events.EventTypeAction,
					Content: models.AgentFinishAction{
						Action: models.ActionTypeAgentFinish,
					},
					Source: "agent",
				})
			}

			// Plugin tools
			if !handled {
				for _, p := range a.Plugins {
					output, ok, err := p.HandleToolCall(ctx, tc.FunctionCall.Name, tc.FunctionCall.Arguments)
					if ok {
						handled = true
						if err != nil {
							output = fmt.Sprintf("Plugin error: %v", err)
						}
						// Record observation
						// We need a generic observation type or reuse CmdOutputObservation?
						// CmdOutputObservation is tailored for bash.
						// Let's use it for now or make it generic.
						// For IPython, it's similar (code execution).
						a.recordObservation(tc.ID, output, "run_ipython")
						// Note: metadata exit code 0 for success
						break
					}
				}
			}

			if !handled {
				// Unknown tool
				a.recordObservation(tc.ID, fmt.Sprintf("Unknown tool: %s", tc.FunctionCall.Name), "error")
			}
		}
	} else if resp.Content != "" {
		a.EventStream.AddEvent(events.Event{
			ID:      uuid.New().String(),
			Type:    events.EventTypeAction,
			Content: models.MessageAction{
				Action:  models.ActionTypeMessage,
				Content: resp.Content,
			},
			Source: "agent",
		})
	}

	return nil
}

func (a *Agent) recordAction(content interface{}) {
	a.EventStream.AddEvent(events.Event{
		ID:      uuid.New().String(),
		Type:    events.EventTypeAction,
		Content: content,
		Source:  "agent",
	})
}

func (a *Agent) recordObservation(toolCallID, content, obsType string, extras ...interface{}) {
	// Extras: [exitCode, command]
	exitCode := 0
	command := ""
	if len(extras) > 0 {
		exitCode = extras[0].(int)
	}
	if len(extras) > 1 {
		command = extras[1].(string)
	}

	obs := models.CmdOutputObservation{
		Observation: obsType,
		Content:     content,
		Metadata: models.CmdOutputMetadata{
			ExitCode: exitCode,
		},
		Command:    command,
		ToolCallID: toolCallID,
	}
	a.EventStream.AddEvent(events.Event{
		ID:      uuid.New().String(),
		Type:    events.EventTypeObservation,
		Content: obs,
		Source:  "runtime",
	})
}

func (a *Agent) eventsToMessages(evts []events.Event) []llm.Message {
	msgs := []llm.Message{
		{Role: "system", Content: a.SystemPrompt},
	}

	for _, e := range evts {
		switch e.Type {
		case events.EventTypeAction:
			bytes, _ := json.Marshal(e.Content)

			if e.Source == "agent" {
				// CmdRunAction
				var cmdAction models.CmdRunAction
				if err := json.Unmarshal(bytes, &cmdAction); err == nil && cmdAction.Action == models.ActionTypeCmdRun {
					msgs = append(msgs, llm.Message{
						Role: "assistant",
						ToolCalls: []llms.ToolCall{
							{
								ID:   cmdAction.ToolCallID,
								Type: "function",
								FunctionCall: &llms.FunctionCall{
									Name:      "execute_bash",
									Arguments: fmt.Sprintf(`{"command": %q, "thought": %q}`, cmdAction.Command, cmdAction.Thought),
								},
							},
						},
					})
					continue
				}
				// Message Action
				var msgAction models.MessageAction
				if err := json.Unmarshal(bytes, &msgAction); err == nil && msgAction.Action == models.ActionTypeMessage {
					msgs = append(msgs, llm.Message{
						Role:    "assistant",
						Content: msgAction.Content,
					})
					continue
				}
			} else {
				// User action
				var msgAction models.MessageAction
				if err := json.Unmarshal(bytes, &msgAction); err == nil {
					msgs = append(msgs, llm.Message{
						Role:    "user",
						Content: msgAction.Content,
					})
					continue
				}
				msgs = append(msgs, llm.Message{
					Role:    "user",
					Content: string(bytes),
				})
			}

		case events.EventTypeObservation:
			bytes, _ := json.Marshal(e.Content)
			var obs models.CmdOutputObservation
			if err := json.Unmarshal(bytes, &obs); err == nil {
				msgs = append(msgs, llm.Message{
					Role:       "tool",
					Content:    obs.Content,
					ToolCallID: obs.ToolCallID,
				})
			}
		}
	}
	return msgs
}

func (a *Agent) RunLoop(ctx context.Context) {
	log.Printf("Starting CodeAct agent loop for conversation %s", a.ConversationID)
	// Init plugins
	for _, p := range a.Plugins {
		if err := p.Init(ctx, a.Runtime); err != nil {
			log.Printf("Failed to init plugin %s: %v", p.Name(), err)
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := a.Step(ctx)
			if err != nil {
				log.Printf("Agent step error: %v", err)
				time.Sleep(5 * time.Second) // Backoff
			}
			time.Sleep(1 * time.Second) // Pace
		}
	}
}
