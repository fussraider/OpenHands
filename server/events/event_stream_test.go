package events

import (
	"encoding/json"
	"openhands-go/server/models"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestEventStreamPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "events.jsonl")

	// 1. Create stream and add generic event
	es1 := NewEventStream("conv1", filePath)
	evt1 := Event{ID: uuid.New().String(), Type: EventTypeAction, Content: "test1", Timestamp: time.Now()}
	es1.AddEvent(evt1)

	// 2. Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("File not created: %s", filePath)
	}

	// 3. Add Typed Event
	cmdAction := models.CmdRunAction{
		Action:  models.ActionTypeCmdRun,
		Command: "ls",
	}
	evt2 := Event{ID: uuid.New().String(), Type: EventTypeAction, Content: cmdAction, Timestamp: time.Now()}
	es1.AddEvent(evt2)

	// 4. Create new stream with same path (simulate restart)
	es2 := NewEventStream("conv1", filePath)

	// 5. Verify events loaded
	events := es2.GetEvents()
	if len(events) != 2 {
		t.Fatalf("Expected 2 events, got %d", len(events))
	}

	// Check evt1 (string fallback)
	if events[0].Content != "test1" {
		t.Errorf("Event 1 content mismatch: %v", events[0].Content)
	}

	// Check evt2 (Typed struct)
	loadedCmd, ok := events[1].Content.(models.CmdRunAction)
	if !ok {
		// It might be a pointer depending on how Unmarshal works?
		// UnmarshalJSON assigns value to e.Content.
		// `var act models.CmdRunAction; json.Unmarshal(..., &act); e.Content = act`
		// So it's a value.
		t.Errorf("Event 2 content type mismatch: %T", events[1].Content)
	} else {
		if loadedCmd.Command != "ls" {
			t.Errorf("Event 2 command mismatch: %s", loadedCmd.Command)
		}
	}
}

func TestEventUnmarshalPolymorphism(t *testing.T) {
	// Action
	jsonStr := `{"id":"1", "type":"action", "content":{"action":"run", "command":"echo hi"}}`
	var e Event
	if err := json.Unmarshal([]byte(jsonStr), &e); err != nil {
		t.Fatal(err)
	}

	cmd, ok := e.Content.(models.CmdRunAction)
	if !ok {
		t.Fatalf("Expected CmdRunAction, got %T", e.Content)
	}
	if cmd.Command != "echo hi" {
		t.Errorf("Command mismatch")
	}

	// Observation
	jsonStrObs := `{"id":"2", "type":"observation", "content":{"observation":"run", "content":"hi", "metadata":{"exit_code":0}}}`
	var e2 Event
	if err := json.Unmarshal([]byte(jsonStrObs), &e2); err != nil {
		t.Fatal(err)
	}

	obs, ok := e2.Content.(models.CmdOutputObservation)
	if !ok {
		t.Fatalf("Expected CmdOutputObservation, got %T", e2.Content)
	}
	if obs.Metadata.ExitCode != 0 {
		t.Errorf("Exit code mismatch")
	}

	// Task Tracking Observation
	jsonStrTask := `{"id":"3", "type":"observation", "content":{"observation":"task_tracking", "content":"Task updated", "task_list":[{"id":"1", "description":"task1", "state":"started"}]}}`
	var e3 Event
	if err := json.Unmarshal([]byte(jsonStrTask), &e3); err != nil {
		t.Fatal(err)
	}

	taskObs, ok := e3.Content.(models.TaskTrackingObservation)
	if !ok {
		t.Fatalf("Expected TaskTrackingObservation, got %T", e3.Content)
	}
	if len(taskObs.TaskList) != 1 || taskObs.TaskList[0].ID != "1" {
		t.Errorf("Task list mismatch")
	}

	// Loop Detection Observation
	jsonStrLoop := `{"id":"4", "type":"observation", "content":{"observation":"loop_detection", "content":"Loop detected"}}`
	var e4 Event
	if err := json.Unmarshal([]byte(jsonStrLoop), &e4); err != nil {
		t.Fatal(err)
	}

	loopObs, ok := e4.Content.(models.LoopDetectionObservation)
	if !ok {
		t.Fatalf("Expected LoopDetectionObservation, got %T", e4.Content)
	}
	if loopObs.Content != "Loop detected" {
		t.Errorf("Loop content mismatch")
	}
}
