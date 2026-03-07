package memory

import (
	"context"
	"log/slog"
	"openhands-go/server/events"
	"openhands-go/server/models"

	"github.com/google/uuid"
)

// Condenser interface defines how to condense event history
type Condenser interface {
	Condense(ctx context.Context, history []events.Event) ([]events.Event, error)
}

// NoOpCondenser does nothing
type NoOpCondenser struct{}

func (c *NoOpCondenser) Condense(ctx context.Context, history []events.Event) ([]events.Event, error) {
	return history, nil
}

// TokenCondenser truncates history based on token count (simulated by event count for MVP)
type TokenCondenser struct {
	MaxEvents int
}

func NewTokenCondenser(maxEvents int) *TokenCondenser {
	if maxEvents <= 0 {
		maxEvents = 100 // Default
	}
	return &TokenCondenser{
		MaxEvents: maxEvents,
	}
}

func (c *TokenCondenser) Condense(ctx context.Context, history []events.Event) ([]events.Event, error) {
	if len(history) <= c.MaxEvents {
		return history, nil
	}

	// Strategy: Keep first few (System/Init) and last N events.
	// Summary logic could be added here using LLM.

	keepStart := 2 // Keep generic start
	keepEnd := c.MaxEvents - keepStart - 1

	if keepEnd < 1 {
		keepEnd = 1
	}

	condensed := make([]events.Event, 0, c.MaxEvents)
	condensed = append(condensed, history[:keepStart]...)

	// Insert summary event
	summaryEvent := events.Event{
		ID:        uuid.New().String(),
		Type:      events.EventTypeObservation,
		Source:    "condenser",
		Timestamp: history[keepStart].Timestamp,
		Content: models.CmdOutputObservation{
			Observation: "condense",
			Content:     "... [History truncated] ...",
		},
	}
	condensed = append(condensed, summaryEvent)

	// Append tail
	// Calculate tail start index
	startTail := len(history) - keepEnd

	// Ensure we don't overlap with kept start if history is short but > max (edge case)
	if startTail < keepStart {
		startTail = keepStart
	}

	condensed = append(condensed, history[startTail:]...)

	removedCount := len(history) - len(condensed)
	if removedCount > 0 {
		slog.Debug("Removed dangling observation(s)", "count", removedCount)
	}

	return condensed, nil
}
