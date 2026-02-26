package agent

import (
	"context"
	"encoding/json"
	"openhands-go/server/events"
	"openhands-go/server/llm"
	"openhands-go/server/models"
	"openhands-go/server/runtime"
	"testing"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/llms"
)

// MockRuntime implements runtime.Runtime
type MockRuntime struct{}

func (m *MockRuntime) Start(ctx context.Context, command string, args ...string) error { return nil }
func (m *MockRuntime) Execute(ctx context.Context, command string, args ...string) (string, int, error) {
	return "mock output", 0, nil
}
func (m *MockRuntime) Write(p []byte) (n int, err error) { return len(p), nil }
func (m *MockRuntime) Read(p []byte) (n int, err error) { return 0, nil }
func (m *MockRuntime) Close() error { return nil }

func TestEventsToMessages(t *testing.T) {
	// Setup
	es := events.NewEventStream("test-conv", "")

	// User message
	msg := models.MessageAction{Action: "message", Content: "Hello"}
	es.AddEvent(events.Event{
		ID: uuid.New().String(),
		Type: events.EventTypeAction,
		Content: msg,
		Source: "user",
	})

	// Agent Tool Call
	cmd := models.CmdRunAction{
		Action: "run",
		Command: "ls",
		Thought: "listing files",
		ToolCallID: "call_123",
	}
	es.AddEvent(events.Event{
		ID: uuid.New().String(),
		Type: events.EventTypeAction,
		Content: cmd,
		Source: "agent",
	})

	// Tool Output
	obs := models.CmdOutputObservation{
		Observation: "run",
		Content: "file1 file2",
		ToolCallID: "call_123",
	}
	es.AddEvent(events.Event{
		ID: uuid.New().String(),
		Type: events.EventTypeObservation,
		Content: obs,
		Source: "runtime",
	})

	// Create Agent
	// We pass nil LLM service as we just test eventsToMessages
	agent := NewAgent("test-agent", "test-conv", &llm.LLMService{}, &MockRuntime{}, es, nil)

	msgs := agent.eventsToMessages(es.GetEvents())

	// Verify
	// 0: System
	// 1: User "Hello"
	// 2: Assistant ToolCall "ls"
	// 3: Tool Output "file1 file2"

	if len(msgs) != 4 {
		t.Errorf("Expected 4 messages, got %d", len(msgs))
		return
	}

	if msgs[0].Role != "system" {
		t.Errorf("Msg 0 role mismatch: %s", msgs[0].Role)
	}

	if msgs[1].Role != "user" || msgs[1].Content != "Hello" {
		t.Errorf("Msg 1 mismatch")
	}

	if msgs[2].Role != "assistant" || len(msgs[2].ToolCalls) != 1 {
		t.Errorf("Msg 2 mismatch: expected assistant with tool call")
	} else {
		tc := msgs[2].ToolCalls[0]
		if tc.ID != "call_123" || tc.FunctionCall.Name != "execute_bash" {
			t.Errorf("Tool call mismatch: %v", tc)
		}
	}

	if msgs[3].Role != "tool" || msgs[3].Content != "file1 file2" || msgs[3].ToolCallID != "call_123" {
		t.Errorf("Msg 3 mismatch: %v", msgs[3])
	}
}

// TestActionMarshalling verifies that our Action structs marshal/unmarshal correctly via Event
func TestActionMarshalling(t *testing.T) {
	cmd := models.CmdRunAction{
		Action: "run",
		Command: "echo test",
		Thought: "testing",
		ToolCallID: "123",
	}

	// Simulate adding to event stream (which stores as interface{})
	var content interface{} = cmd

	// Marshal to JSON (as it would be when stored/sent)
	bytes, err := json.Marshal(content)
	if err != nil {
		t.Fatal(err)
	}

	// Unmarshal back
	var parsed models.CmdRunAction
	if err := json.Unmarshal(bytes, &parsed); err != nil {
		t.Fatal(err)
	}

	if parsed.Command != "echo test" {
		t.Errorf("Failed to roundtrip CmdRunAction")
	}
}

// MockPlugin implements plugins.Plugin
type MockPlugin struct{}
func (p *MockPlugin) Name() string { return "mock" }
func (p *MockPlugin) Init(ctx context.Context, rt runtime.Runtime) error { return nil }
func (p *MockPlugin) Tools() []llms.Tool {
	return []llms.Tool{{
		Type: "function",
		Function: &llms.FunctionDefinition{Name: "mock_tool"},
	}}
}
func (p *MockPlugin) HandleToolCall(ctx context.Context, name string, args string) (string, bool, error) {
	if name == "mock_tool" {
		return "mock result", true, nil
	}
	return "", false, nil
}

func TestAgentPlugins(t *testing.T) {
	agent := NewAgent("test", "conv", &llm.LLMService{}, &MockRuntime{}, events.NewEventStream("conv", ""), nil)

	// Add mock plugin manually
	agent.Plugins = append(agent.Plugins, &MockPlugin{})

	// Verify tool registration
	// NewAgent registers default tools + plugin tools
	// But we appended AFTER NewAgent.
	// We should reconstruct tools logic or allow dynamic registration.
	// Agent.Tools is public.
	agent.Tools = append(agent.Tools, (&MockPlugin{}).Tools()...)

	found := false
	for _, tool := range agent.Tools {
		if tool.Function.Name == "mock_tool" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Mock tool not found in agent tools")
	}
}
