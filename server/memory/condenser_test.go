package memory

import (
	"context"
	"openhands-go/server/events"
	"testing"
)

func TestTokenCondenser(t *testing.T) {
	// Create history > max events
	maxEvents := 5
	history := make([]events.Event, 10)
	for i := 0; i < 10; i++ {
		history[i] = events.Event{Content: i} // Using int for simplicity
	}

	condenser := NewTokenCondenser(maxEvents)

	condensed, err := condenser.Condense(context.Background(), history)
	if err != nil {
		t.Fatalf("Condense failed: %v", err)
	}

	// Expect length <= maxEvents + 1 (because of summary event insertion, might vary slightly based on logic)
	// Logic: keepStart (2) + summary (1) + keepEnd (max-3) = maxEvents
	// So length should be exactly maxEvents.

	if len(condensed) != maxEvents {
		t.Errorf("Expected %d events, got %d", maxEvents, len(condensed))
	}

	// Verify content
	// history[0], history[1] kept
	if condensed[0].Content != 0 || condensed[1].Content != 1 {
		t.Errorf("Start events mismatched")
	}

	// condensed[2] should be summary
	// But type check is hard since we used int content for test, but Condenser inserts Typed Observation.
	// Just check index 3 (first of tail)
	// Tail length = maxEvents - 3 = 2
	// Tail starts at len(history) - 2 = 8
	// So condensed[3] == history[8], condensed[4] == history[9]

	if condensed[3].Content != 8 || condensed[4].Content != 9 {
		t.Errorf("Tail events mismatched")
	}
}

func TestNoOpCondenser(t *testing.T) {
	history := make([]events.Event, 10)
	condenser := &NoOpCondenser{}
	condensed, _ := condenser.Condense(context.Background(), history)
	if len(condensed) != 10 {
		t.Errorf("NoOp changed length")
	}
}
