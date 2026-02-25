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
}

func NewAgent(id, conversationID string, llmService *llm.LLMService, rt runtime.Runtime, es *events.EventStream) *Agent {
	// Define execute_bash tool
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

	return &Agent{
		ID:             id,
		ConversationID: conversationID,
		LLM:            llmService,
		Runtime:        rt,
		EventStream:    es,
		SystemPrompt:   "You are a helpful coding agent using the CodeAct framework. You can execute bash commands to solve tasks. When you are done, call finish.",
		Tools:          []llms.Tool{bashTool, finishTool},
	}
}

func (a *Agent) Step(ctx context.Context) error {
	// 1. Get History & Construct Messages
	history := a.EventStream.GetEvents()
	messages := a.eventsToMessages(history)

	// 2. LLM Completion
	// CodeActAgent uses tool calling.
	resp, err := a.LLM.CompleteWithTools(ctx, messages, a.Tools)
	if err != nil {
		return fmt.Errorf("LLM completion error: %w", err)
	}

	// 3. Handle Tool Calls
	if len(resp.ToolCalls) > 0 {
		for _, tc := range resp.ToolCalls {
			switch tc.FunctionCall.Name {
			case "execute_bash":
				var args struct {
					Command string `json:"command"`
					Thought string `json:"thought"`
				}
				if err := json.Unmarshal([]byte(tc.FunctionCall.Arguments), &args); err != nil {
					log.Printf("Error unmarshalling execute_bash args: %v", err)
					continue
				}

				// Create CmdRunAction event
				action := models.CmdRunAction{
					Action:     models.ActionTypeCmdRun,
					Command:    args.Command,
					Thought:    args.Thought,
					ToolCallID: tc.ID,
				}

				a.EventStream.AddEvent(events.Event{
					ID:      uuid.New().String(),
					Type:    events.EventTypeAction,
					Content: action,
					Source:  "agent",
				})

				// Execute
				output, exitCode, err := a.Runtime.Execute(ctx, "bash", "-c", args.Command)
				content := output
				if err != nil {
					content = fmt.Sprintf("Error executing command: %v", err)
				}

				// Create Observation event
				obs := models.CmdOutputObservation{
					Observation: "run",
					Content:     content,
					Metadata: models.CmdOutputMetadata{
						ExitCode: exitCode,
					},
					Command:    args.Command,
					ToolCallID: tc.ID,
				}
				a.EventStream.AddEvent(events.Event{
					ID:      uuid.New().String(),
					Type:    events.EventTypeObservation,
					Content: obs,
					Source:  "runtime",
				})

			case "finish":
				a.EventStream.AddEvent(events.Event{
					ID:      uuid.New().String(),
					Type:    events.EventTypeAction,
					Content: models.AgentFinishAction{
						Action: models.ActionTypeAgentFinish,
					},
					Source: "agent",
				})
				// Ideally stop the loop
			}
		}
	} else if resp.Content != "" {
		// Just a message
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

func (a *Agent) eventsToMessages(evts []events.Event) []llm.Message {
	msgs := []llm.Message{
		{Role: "system", Content: a.SystemPrompt},
	}

	for _, e := range evts {
		switch e.Type {
		case events.EventTypeAction:
			// Handle different action types
			// We need to check the Content type carefully.
			// It might be map[string]interface{} (from JSON unmarshal) or struct (from AddEvent)
			// Using fmt.Sprintf("%v") is risky for structured data.
			// Ideally we should unmarshal it to specific types.

			// For now, let's assume simple cases or marshal back to JSON
			bytes, _ := json.Marshal(e.Content)

			// Determine role
			// role := "user"
			// if e.Source == "agent" {
			// 	role = "assistant"
			// }

			// If it's a Tool Call (Agent Action with ToolCallID)
			// We need to reconstruct the ToolCall object for LLM history
			if e.Source == "agent" {
				// Try to parse as CmdRunAction
				var cmdAction models.CmdRunAction
				if err := json.Unmarshal(bytes, &cmdAction); err == nil && cmdAction.Action == models.ActionTypeCmdRun {
					// It's a tool call
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
				// Finish action
				var finishAction models.AgentFinishAction
				if err := json.Unmarshal(bytes, &finishAction); err == nil && finishAction.Action == models.ActionTypeAgentFinish {
					// Finish tool call
					// Logic similar to above
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
				// e.g. Initial message or user feedback
				// If e.Content is a struct, unmarshal/marshal
				// Assuming simple user message for now if source is user
				var msgAction models.MessageAction
				if err := json.Unmarshal(bytes, &msgAction); err == nil {
					msgs = append(msgs, llm.Message{
						Role:    "user",
						Content: msgAction.Content,
					})
					continue
				}

				// Fallback
				msgs = append(msgs, llm.Message{
					Role:    "user",
					Content: string(bytes),
				})
			}

		case events.EventTypeObservation:
			// Tool Output
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
