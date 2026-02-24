package agent

import (
	"context"
	"fmt"
	"log"
	"openhands-go/server/events"
	"openhands-go/server/llm"
	"openhands-go/server/runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Agent struct {
	ID             string
	ConversationID string
	LLM            *llm.LLMService
	Runtime        runtime.Runtime
	EventStream    *events.EventStream
}

func NewAgent(id, conversationID string, llmService *llm.LLMService, rt runtime.Runtime, es *events.EventStream) *Agent {
	return &Agent{
		ID:             id,
		ConversationID: conversationID,
		LLM:            llmService,
		Runtime:        rt,
		EventStream:    es,
	}
}

func (a *Agent) Step(ctx context.Context) error {
	// 1. Get History
	history := a.EventStream.GetEvents()
	messages := a.eventsToMessages(history)

	// 2. LLM Completion
	responseContent, err := a.LLM.Complete(messages)
	if err != nil {
		return err
	}

	// 3. Parse Response (Mock: assume single line command "RUN: <cmd>" or "MSG: <msg>")
	// In reality, this would parse Action objects from JSON/XML
	action, content := a.parseResponse(responseContent)

	// Add Agent Action to stream
	a.EventStream.AddEvent(events.Event{
		ID:      uuid.New().String(),
		Type:    events.EventTypeAction,
		Content: map[string]string{"action": action, "args": content},
		Source:  "agent",
	})

	// 4. Execute Action
	if action == "run" {
		err := a.Runtime.Start(ctx, "bash", "-c", content)
		output := ""
		if err != nil {
			output = fmt.Sprintf("Error starting command: %v", err)
		} else {
			buf := make([]byte, 1024)
			n, err := a.Runtime.Read(buf)
			if err != nil {
				output = fmt.Sprintf("Error reading output: %v", err)
			} else {
				output = string(buf[:n])
			}
		}

		// Add Observation
		a.EventStream.AddEvent(events.Event{
			ID:      uuid.New().String(),
			Type:    events.EventTypeObservation,
			Content: map[string]string{"output": output},
			Source:  "runtime",
		})
	} else if action == "message" {
		// Just a thought/message
		log.Printf("[Agent %s] Message: %s", a.ID, content)
	}

	return nil
}

func (a *Agent) eventsToMessages(evts []events.Event) []llm.Message {
	msgs := []llm.Message{
		{Role: "system", Content: "You are a helpful coding agent. Respond with 'RUN: <command>' to execute code or 'MSG: <message>' to talk."},
	}
	for _, e := range evts {
		// Simple mapping
		role := "user"
		content := fmt.Sprintf("%v", e.Content)
		if e.Source == "agent" {
			role = "assistant"
		} else if e.Source == "runtime" {
			role = "user" // Observation is user/system info
			content = fmt.Sprintf("Output: %v", e.Content)
		}
		msgs = append(msgs, llm.Message{Role: role, Content: content})
	}
	return msgs
}

func (a *Agent) parseResponse(resp string) (string, string) {
	if strings.HasPrefix(resp, "RUN: ") {
		return "run", strings.TrimPrefix(resp, "RUN: ")
	}
	if strings.HasPrefix(resp, "MSG: ") {
		return "message", strings.TrimPrefix(resp, "MSG: ")
	}
	return "message", resp
}

func (a *Agent) RunLoop(ctx context.Context) {
	log.Printf("Starting agent loop for conversation %s", a.ConversationID)
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
