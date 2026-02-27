package memory

import (
	"context"
	"fmt"
	"openhands-go/server/events"
	"openhands-go/server/llm"
	"openhands-go/server/models"

	"github.com/google/uuid"
)

type LLMSummarizingCondenser struct {
	llmService *llm.LLMService
	maxSize    int
	keepFirst  int
}

func NewLLMSummarizingCondenser(llmService *llm.LLMService, maxSize, keepFirst int) *LLMSummarizingCondenser {
	if maxSize <= 0 {
		maxSize = 100
	}
	if keepFirst < 0 {
		keepFirst = 1
	}
	return &LLMSummarizingCondenser{
		llmService: llmService,
		maxSize:    maxSize,
		keepFirst:  keepFirst,
	}
}

func (c *LLMSummarizingCondenser) Condense(ctx context.Context, history []events.Event) ([]events.Event, error) {
	if len(history) <= c.maxSize {
		return history, nil
	}

	targetSize := c.maxSize / 2
	eventsFromTail := targetSize - c.keepFirst - 1

	if eventsFromTail < 1 {
		eventsFromTail = 1
	}

	head := history[:c.keepFirst]
	tail := history[len(history)-eventsFromTail:]

	// The middle part to forget
	forgotten := history[c.keepFirst : len(history)-eventsFromTail]

	if len(forgotten) == 0 {
		return history, nil
	}

	// Create prompt
	prompt := "Summarize the following events. Maintain task tracking, code state, and context.\n\n"
	for _, e := range forgotten {
		prompt += fmt.Sprintf("Event %s: %+v\n", e.Type, e.Content)
	}

	msg := llm.Message{Role: "user", Content: prompt}
	resp, err := c.llmService.Complete(ctx, []llm.Message{msg})
	if err != nil {
		// Fallback to basic token condenser if LLM fails
		return NewTokenCondenser(c.maxSize).Condense(ctx, history)
	}

	summaryEvent := events.Event{
		ID:        uuid.New().String(),
		Type:      events.EventTypeObservation,
		Source:    "condenser",
		Content: models.CmdOutputObservation{
			Observation: "condense",
			Content:     resp, // The summary
		},
	}

	condensed := make([]events.Event, 0, len(head)+1+len(tail))
	condensed = append(condensed, head...)
	condensed = append(condensed, summaryEvent)
	condensed = append(condensed, tail...)

	return condensed, nil
}

type PipelineCondenser struct {
	condensers []Condenser
}

func NewPipelineCondenser(condensers ...Condenser) *PipelineCondenser {
	return &PipelineCondenser{condensers: condensers}
}

func (p *PipelineCondenser) Condense(ctx context.Context, history []events.Event) ([]events.Event, error) {
	current := history
	var err error
	for _, c := range p.condensers {
		current, err = c.Condense(ctx, current)
		if err != nil {
			return current, err // Stop on first error or continue? Stop is safer.
		}
	}
	return current, nil
}
