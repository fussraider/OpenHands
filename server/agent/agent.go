package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"openhands-go/server/agent/prompts"
	"openhands-go/server/config"
	"openhands-go/server/events"
	"openhands-go/server/llm"
	"openhands-go/server/memory"
	"openhands-go/server/models"
	"openhands-go/server/runtime"
	"openhands-go/server/runtime/plugins"
	"openhands-go/server/security"
	"openhands-go/server/runtime/plugins/browser"
	"openhands-go/server/runtime/plugins/jupyter"
	"time"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/llms"
	"go.opentelemetry.io/otel"
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
	PromptManager  *prompts.PromptManager
	Config         *config.Config
	LoopDetector   *LoopDetector
	Condenser      memory.Condenser
	Security       security.SecurityAnalyzer
}

func NewAgent(id, conversationID string, llmService *llm.LLMService, rt runtime.Runtime, es *events.EventStream, delegator Delegator, cfg *config.Config) *Agent {
	// Initialize Plugins
	plugs := []plugins.Plugin{
		jupyter.NewJupyterPlugin(),
		browser.NewBrowserPlugin(),
	}

	pm, err := prompts.New()
	if err != nil {
		slog.Error("Failed to initialize prompt manager", "error", err)
		// Fallback or panic? For now log error.
	}

	agent := &Agent{
		ID:             id,
		ConversationID: conversationID,
		LLM:            llmService,
		Runtime:        rt,
		EventStream:    es,
		Plugins:        plugs,
		Delegator:      delegator,
		PromptManager:  pm,
		Config:         cfg,
		LoopDetector:   NewLoopDetector(),
		Security:       security.NewBasicAnalyzer(),
	}

	// Initialize Security Analyzer
	// If LLM config supports it, use LLM analyzer.
	// For now, if "EnableSecurityAnalyzer" is true in config, and we have LLM service, we use it?
	// But `security.NewBasicAnalyzer` is default.
	// Let's check config.
	if cfg.Security.SecurityAnalyzer == "llm" {
		agent.Security = security.NewLLMSecurityAnalyzer(llmService)
	}

	// Initialize Condenser
	if cfg.Agent.EnableHistoryTruncation {
		agent.Condenser = memory.NewTokenCondenser(cfg.Agent.MaxEvents)
	} else {
		agent.Condenser = &memory.NoOpCondenser{}
	}

	// Render System Prompt
	if pm != nil {
		sysPrompt, err := pm.RenderSystemPrompt(prompts.SystemPromptContext{
			CLIMode: false, // Default to sandbox mode
		})
		if err != nil {
			slog.Error("Failed to render system prompt", "error", err)
			agent.SystemPrompt = "You are a helpful coding agent using the CodeAct framework. You can execute bash commands to solve tasks. When you are done, call finish."
		} else {
			agent.SystemPrompt = sysPrompt
		}
	} else {
		agent.SystemPrompt = "You are a helpful coding agent using the CodeAct framework. You can execute bash commands to solve tasks. When you are done, call finish."
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
	// Start Span
	ctx, span := otel.Tracer("agent").Start(ctx, "Agent.Step")
	defer span.End()

	// 1. Get History & Construct Messages
	history := a.EventStream.GetEvents()

	// Condense History
	if a.Condenser != nil {
		var err error
		history, err = a.Condenser.Condense(ctx, history)
		if err != nil {
			slog.Error("Failed to condense history", "error", err)
		}
	}

	// 2. Loop Detection
	if stuck, analysis := a.LoopDetector.IsStuck(history); stuck {
		slog.Warn("Agent stuck in loop", "type", analysis.LoopType, "times", analysis.RepeatTimes)
		a.EventStream.AddEvent(events.Event{
			ID:   uuid.New().String(),
			Type: events.EventTypeObservation,
			Content: models.LoopDetectionObservation{
				Observation: "loop_detection",
				Content:     fmt.Sprintf("⚠️ Loop detected! You are repeating the same action/observation pattern. Please stop and reflect on why this is happening. Try a different approach."),
			},
			Source: "runtime",
		})
		// We add the observation and return, effectively giving the agent a chance to react in the next step
		return nil
	}

	messages := a.eventsToMessages(ctx, history)

	// 3. LLM Completion
	llmSpanCtx, llmSpan := otel.Tracer("agent").Start(ctx, "LLM.Complete")
	resp, err := a.LLM.CompleteWithTools(llmSpanCtx, messages, a.Tools)
	llmSpan.End()
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

				cmdAction := models.CmdRunAction{
					Action:     models.ActionTypeCmdRun,
					Command:    args.Command,
					Thought:    args.Thought,
					ToolCallID: tc.ID,
				}
				a.recordAction(cmdAction)

				// Security Check
				risk, reason, _ := a.Security.Analyze(ctx, cmdAction)
				if risk == security.RiskHigh {
					// Block high risk actions for now (or require confirmation if we had UI support)
					a.recordObservation(tc.ID, fmt.Sprintf("Security Alert: Action blocked. Reason: %s", reason), "error")
					continue
				}

				output, exitCode, err := a.Runtime.Execute(ctx, "bash", "-c", args.Command)
				content := output
				if err != nil {
					content = fmt.Sprintf("Error executing command: %v", err)
				}

				a.recordObservation(tc.ID, content, "run", exitCode, args.Command)

			case "finish":
				handled = true
				a.EventStream.AddEvent(events.Event{
					ID:   uuid.New().String(),
					Type: events.EventTypeAction,
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
			ID:   uuid.New().String(),
			Type: events.EventTypeAction,
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

func (a *Agent) getAdditionalInfo(ctx context.Context) string {
	if a.PromptManager == nil {
		return ""
	}

	// Gather info
	cwd, _ := a.Runtime.GetCwd(ctx)
	if cwd == "" {
		cwd = "/workspace"
	}

	infoCtx := prompts.AdditionalInfoContext{
		RuntimeInfo: &prompts.RuntimeInfo{
			WorkingDir: cwd,
			Date:       time.Now().Format("2006-01-02"),
		},
		// TODO: Populate RepositoryInfo, AvailableHosts, etc.
	}

	out, err := a.PromptManager.RenderAdditionalInfo(infoCtx)
	if err != nil {
		slog.Error("Failed to render additional info", "error", err)
		return ""
	}
	return out
}

func (a *Agent) eventsToMessages(ctx context.Context, evts []events.Event) []llm.Message {
	msgs := []llm.Message{
		{Role: "system", Content: a.SystemPrompt},
	}

	// Add Additional Info
	if additionalInfo := a.getAdditionalInfo(ctx); additionalInfo != "" {
		msgs = append(msgs, llm.Message{
			Role:    "user",
			Content: additionalInfo,
		})
	}

	for _, e := range evts {
		switch e.Type {
		case events.EventTypeAction:
			if e.Source == "agent" {
				switch act := e.Content.(type) {
				case models.CmdRunAction:
					msgs = append(msgs, llm.Message{
						Role: "assistant",
						ToolCalls: []llms.ToolCall{
							{
								ID:   act.ToolCallID,
								Type: "function",
								FunctionCall: &llms.FunctionCall{
									Name:      "execute_bash",
									Arguments: fmt.Sprintf(`{"command": %q, "thought": %q}`, act.Command, act.Thought),
								},
							},
						},
					})
				case models.AgentDelegateAction:
					inputsBytes, _ := json.Marshal(act.Inputs)
					msgs = append(msgs, llm.Message{
						Role: "assistant",
						ToolCalls: []llms.ToolCall{
							{
								ID:   act.ToolCallID,
								Type: "function",
								FunctionCall: &llms.FunctionCall{
									Name:      "delegate",
									Arguments: fmt.Sprintf(`{"agent": %q, "inputs": %s, "thought": %q}`, act.Agent, string(inputsBytes), act.Thought),
								},
							},
						},
					})
				case models.MessageAction:
					msgs = append(msgs, llm.Message{
						Role:    "assistant",
						Content: act.Content,
					})
				}
			} else {
				// User action
				switch act := e.Content.(type) {
				case models.MessageAction:
					msgs = append(msgs, llm.Message{
						Role:    "user",
						Content: act.Content,
					})
				case string:
					msgs = append(msgs, llm.Message{
						Role:    "user",
						Content: act,
					})
				default:
					// Try marshalling generic content
					bytes, _ := json.Marshal(e.Content)
					msgs = append(msgs, llm.Message{
						Role:    "user",
						Content: string(bytes),
					})
				}
			}

		case events.EventTypeObservation:
			switch obs := e.Content.(type) {
			case models.CmdOutputObservation:
				msgs = append(msgs, llm.Message{
					Role:       "tool",
					Content:    obs.Content,
					ToolCallID: obs.ToolCallID,
				})
			default:
				// Fallback
				bytes, _ := json.Marshal(e.Content)
				msgs = append(msgs, llm.Message{
					Role:    "tool",
					Content: string(bytes),
				})
			}
		}
	}
	return msgs
}

func (a *Agent) RunLoop(ctx context.Context) {
	slog.Info("Starting CodeAct agent loop", "conversation_id", a.ConversationID)
	if err := a.InitPlugins(ctx); err != nil {
		slog.Error("Failed to init plugins", "error", err)
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

func (a *Agent) InitPlugins(ctx context.Context) error {
	for _, p := range a.Plugins {
		if err := p.Init(ctx, a.Runtime); err != nil {
			return fmt.Errorf("failed to init plugin %s: %w", p.Name(), err)
		}
	}
	return nil
}

// RunUntilDone runs the agent loop until it finishes a task or context is cancelled.
// Returns the outputs from AgentFinishAction if successful.
func (a *Agent) RunUntilDone(ctx context.Context) (map[string]string, error) {
	slog.Info("Starting CodeAct agent sub-task", "conversation_id", a.ConversationID)
	if err := a.InitPlugins(ctx); err != nil {
		return nil, err
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// Peek for FinishAction before stepping?
			// Actually Step handles adding events. We need to check if the LAST event was finish.
			// But Step returns error, not event.

			// Check if finished
			evts := a.EventStream.GetEvents()
			if len(evts) > 0 {
				lastEvent := evts[len(evts)-1]
				if lastEvent.Type == events.EventTypeAction {
					if finishAct, ok := lastEvent.Content.(models.AgentFinishAction); ok {
						return finishAct.Outputs, nil
					}
				}
			}

			err := a.Step(ctx)
			if err != nil {
				return nil, err
			}
			time.Sleep(100 * time.Millisecond) // Faster pace for sub-tasks?
		}
	}
}

func init() {
	RegisterAgent("CodeActAgent")
}
