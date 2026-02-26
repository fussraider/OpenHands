package events

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestEventStreamPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "events.jsonl")

	// 1. Create stream and add events
	es1 := NewEventStream("conv1", filePath)
	evt1 := Event{ID: uuid.New().String(), Type: EventTypeAction, Content: "test1", Timestamp: time.Now()}
	es1.AddEvent(evt1)

	// 2. Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("File not created: %s", filePath)
	}

	// 3. Create new stream with same path (simulate restart)
	es2 := NewEventStream("conv1", filePath)

	// 4. Verify events loaded
	events := es2.GetEvents()
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}
	if events[0].ID != evt1.ID {
		t.Errorf("Event ID mismatch")
	}
}
