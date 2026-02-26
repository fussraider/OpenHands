package agent

import (
	"openhands-go/server/events"
	"openhands-go/server/models"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestLoopDetector(t *testing.T) {
	ld := NewLoopDetector()

	// 1. Not Stuck
	history := []events.Event{
		{Type: events.EventTypeAction, Content: models.CmdRunAction{Command: "ls"}},
		{Type: events.EventTypeObservation, Content: models.CmdOutputObservation{Content: "file1"}},
	}
	stuck, _ := ld.IsStuck(history)
	if stuck {
		t.Error("Should not be stuck")
	}

	// 2. Repeating Action/Observation (4 times)
	history = []events.Event{}
	for i := 0; i < 4; i++ {
		history = append(history,
			events.Event{
				ID:        uuid.New().String(),
				Type:      events.EventTypeAction,
				Content:   models.CmdRunAction{Command: "echo loop"},
				Timestamp: time.Now(),
			},
			events.Event{
				ID:        uuid.New().String(),
				Type:      events.EventTypeObservation,
				Content:   models.CmdOutputObservation{Content: "loop", Metadata: models.CmdOutputMetadata{ExitCode: 0}},
				Timestamp: time.Now(),
			},
		)
	}

	stuck, analysis := ld.IsStuck(history)
	if !stuck {
		t.Error("Should be stuck")
	}
	if analysis.LoopType != LoopTypeRepeatingActionObservation {
		t.Errorf("Wrong loop type: %s", analysis.LoopType)
	}

	// 3. Repeating Action/Error (3 times)
	// Reset history
	history = []events.Event{}
	for i := 0; i < 3; i++ {
		history = append(history,
			events.Event{
				ID:        uuid.New().String(),
				Type:      events.EventTypeAction,
				Content:   models.CmdRunAction{Command: "bad_cmd"},
				Timestamp: time.Now(),
			},
			events.Event{
				ID:        uuid.New().String(),
				Type:      events.EventTypeObservation,
				Content:   models.CmdOutputObservation{Content: "command not found", Metadata: models.CmdOutputMetadata{ExitCode: 127}},
				Timestamp: time.Now(),
			},
		)
	}

	stuck, analysis = ld.IsStuck(history)
	if !stuck {
		t.Error("Should be stuck in error loop")
	}
	if analysis != nil && analysis.LoopType != LoopTypeRepeatingActionError {
		t.Errorf("Wrong loop type: %s", analysis.LoopType)
	}
}

func TestLoopDetector_IsErrorObservation(t *testing.T) {
	ld := NewLoopDetector()

	obs1 := models.CmdOutputObservation{Metadata: models.CmdOutputMetadata{ExitCode: 1}}
	if !ld.isErrorObservation(obs1) {
		t.Error("Exit code 1 should be error")
	}

	obs2 := models.CmdOutputObservation{Metadata: models.CmdOutputMetadata{ExitCode: 0}}
	if ld.isErrorObservation(obs2) {
		t.Error("Exit code 0 should not be error")
	}
}
