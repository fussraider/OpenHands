package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"openhands-go/server/events"
	"openhands-go/server/llm"
	"openhands-go/server/models"
	"openhands-go/server/runtime"
	"openhands-go/server/runtime/plugins"
	"openhands-go/server/runtime/plugins/browser"
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
	Delegator      Delegator
}

func NewAgent(id, conversationID string, llmService *llm.LLMService, rt runtime.Runtime, es *events.EventStream, delegator Delegator) *Agent {
	// Initialize Plugins
	plugs := []plugins.Plugin{
		jupyter.NewJupyterPlugin(),
		browser.NewBrowserPlugin(),
	}

	agent := &Agent{
		ID:             id,
		ConversationID: conversationID,
		LLM:            llmService,
		Runtime:        rt,
		EventStream:    es,
		SystemPrompt:   "You are a helpful coding agent using the CodeAct framework. You can execute bash commands to solve tasks. When you are done, call finish.",
		Plugins:        plugs,
		Delegator:      delegator,
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

	// Delegate tool (if delegator is present)
	var delegateTool llms.Tool
	if delegator != nil {
		delegateTool = llms.Tool{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "delegate",
				Description: "Delegate a subtask to another agent.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"agent": map[string]interface{}{
							"type":        "string",
							"description": "The name/role of the agent to delegate to (e.g., 'researcher', 'coder').",
						},
						"inputs": map[string]interface{}{
							"type":        "object",
							"description": "Input parameters for the delegated task.",
						},
						"thought": map[string]interface{}{
							"type":        "string",
							"description": "Your reasoning for delegating this task.",
						},
					},
					"required": []string{"agent", "inputs"},
				},
			},
		}
		tools = append(tools, bashTool, finishTool, delegateTool)
	} else {
		tools = append(tools, bashTool, finishTool)
	}

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
					slog.Error("Error unmarshalling execute_bash args", "error", err)
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

			case "delegate":
				handled = true
				var args struct {
					Agent   string                 `json:"agent"`
					Inputs  map[string]interface{} `json:"inputs"`
					Thought string                 `json:"thought"`
				}
				if err := json.Unmarshal([]byte(tc.FunctionCall.Arguments), &args); err != nil {
					a.recordObservation(tc.ID, fmt.Sprintf("Error unmarshalling delegate args: %v", err), "delegate")
					continue
				}

				a.recordAction(models.AgentDelegateAction{
					Action:     models.ActionTypeDelegate,
					Agent:      args.Agent,
					Inputs:     args.Inputs,
					Thought:    args.Thought,
					ToolCallID: tc.ID,
				})

				if a.Delegator != nil {
					outputs, err := a.Delegator.Delegate(ctx, args.Agent, args.Inputs)
					content := ""
					if err != nil {
						content = fmt.Sprintf("Delegation error: %v", err)
					} else {
						outBytes, _ := json.Marshal(outputs)
						content = string(outBytes)
					}
					// Use generic observation for delegation result
					a.recordObservation(tc.ID, content, "delegate")
				} else {
					a.recordObservation(tc.ID, "Delegation not supported", "delegate")
				}
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
						a.recordObservation(tc.ID, output, "run_ipython")
						break
					}
				}
			}

			if !handled {
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
				// Delegate Action
				var delegateAction models.AgentDelegateAction
				if err := json.Unmarshal(bytes, &delegateAction); err == nil && delegateAction.Action == models.ActionTypeDelegate {
					inputsBytes, _ := json.Marshal(delegateAction.Inputs)
					msgs = append(msgs, llm.Message{
						Role: "assistant",
						ToolCalls: []llms.ToolCall{
							{
								ID:   delegateAction.ToolCallID,
								Type: "function",
								FunctionCall: &llms.FunctionCall{
									Name:      "delegate",
									Arguments: fmt.Sprintf(`{"agent": %q, "inputs": %s, "thought": %q}`, delegateAction.Agent, string(inputsBytes), delegateAction.Thought),
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
	slog.Info("Starting CodeAct agent loop", "conversation_id", a.ConversationID)
	// Init plugins
	for _, p := range a.Plugins {
		if err := p.Init(ctx, a.Runtime); err != nil {
			slog.Error("Failed to init plugin", "plugin", p.Name(), "error", err)
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := a.Step(ctx)
			if err != nil {
				slog.Error("Agent step error", "error", err)
				time.Sleep(5 * time.Second) // Backoff
			}
			time.Sleep(1 * time.Second) // Pace
		}
	}
}
