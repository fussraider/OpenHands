package agent

import (
	"log/slog"
	"openhands-go/server/events"
	"openhands-go/server/models"
	"reflect"
)

type LoopType string

const (
	LoopTypeRepeatingActionObservation = "repeating_action_observation"
	LoopTypeRepeatingActionError       = "repeating_action_error"
)

type StuckAnalysis struct {
	LoopType    LoopType
	RepeatTimes int
	LoopStartID string
}

type LoopDetector struct {
	// We keep a reference to events (or pass them in)
}

func NewLoopDetector() *LoopDetector {
	return &LoopDetector{}
}

// IsStuck checks if the agent is stuck in a loop based on the event history.
// It returns true and the analysis if a loop is detected.
func (ld *LoopDetector) IsStuck(history []events.Event) (bool, *StuckAnalysis) {
	// Filter out non-relevant events (e.g. state changes, nulls) if needed.
	// For now, we work with the full history but focus on Action/Observation pairs.

	// Need at least a few events
	if len(history) < 6 {
		return false, nil
	}

	// Extract actions and observations in reverse order
	var lastActions []events.Event
	var lastObservations []events.Event

	for i := len(history) - 1; i >= 0; i-- {
		evt := history[i]
		if evt.Type == events.EventTypeAction {
			lastActions = append(lastActions, evt)
		} else if evt.Type == events.EventTypeObservation {
			lastObservations = append(lastObservations, evt)
		}

		if len(lastActions) >= 4 && len(lastObservations) >= 4 {
			break
		}
	}

	// Scenario 1: Same Action, Same Observation (4 times)
	if len(lastActions) >= 4 && len(lastObservations) >= 4 {
		if ld.isRepeatingActionObservation(lastActions, lastObservations) {
			slog.Warn("Loop detected: repeating action/observation")
			return true, &StuckAnalysis{
				LoopType:    LoopTypeRepeatingActionObservation,
				RepeatTimes: 4,
				// Approximating start ID as the ID of the 4th last action
				LoopStartID: lastActions[3].ID,
			}
		}
	}

	// Scenario 2: Repeating Action, Error (3 times)
	// Implement if needed.
	if len(lastActions) >= 3 && len(lastObservations) >= 3 {
		if ld.isRepeatingActionError(lastActions, lastObservations) {
			slog.Warn("Loop detected: repeating action/error")
			return true, &StuckAnalysis{
				LoopType:    LoopTypeRepeatingActionError,
				RepeatTimes: 3,
				LoopStartID: lastActions[2].ID,
			}
		}
	}

	return false, nil
}

func (ld *LoopDetector) isRepeatingActionObservation(actions []events.Event, observations []events.Event) bool {
	// Actions/Observations are in reverse order (newest first)
	// Check if action[0] == action[1] == action[2] == action[3]
	// And observation[0] == observation[1] == observation[2] == observation[3]

	firstAct := actions[0].Content
	for i := 1; i < 4; i++ {
		if !ld.actionsEqual(firstAct, actions[i].Content) {
			return false
		}
	}

	firstObs := observations[0].Content
	for i := 1; i < 4; i++ {
		if !ld.observationsEqual(firstObs, observations[i].Content) {
			return false
		}
	}

	return true
}

func (ld *LoopDetector) isRepeatingActionError(actions []events.Event, observations []events.Event) bool {
	// Check if last 3 actions are same
	firstAct := actions[0].Content
	for i := 1; i < 3; i++ {
		if !ld.actionsEqual(firstAct, actions[i].Content) {
			return false
		}
	}

	// Check if last 3 observations are Errors
	// We don't have explicit ErrorObservation struct yet, maybe CmdOutput with specific content?
	// Or check if it contains "error" string?
	// In Go backend, execution errors often return CmdOutputObservation with non-zero exit code or error string.
	for i := 0; i < 3; i++ {
		if !ld.isErrorObservation(observations[i].Content) {
			return false
		}
	}

	return true
}

func (ld *LoopDetector) actionsEqual(a, b interface{}) bool {
	// Specific check for CmdRunAction
	if actA, ok := a.(models.CmdRunAction); ok {
		if actB, ok := b.(models.CmdRunAction); ok {
			return actA.Command == actB.Command
			// Ignore Thought?
		}
	}
	// Fallback to deep equal
	return reflect.DeepEqual(a, b)
}

func (ld *LoopDetector) observationsEqual(a, b interface{}) bool {
	if obsA, ok := a.(models.CmdOutputObservation); ok {
		if obsB, ok := b.(models.CmdOutputObservation); ok {
			// Compare content and exit code, ignore timestamps/IDs
			return obsA.Content == obsB.Content && obsA.Metadata.ExitCode == obsB.Metadata.ExitCode
		}
	}
	return reflect.DeepEqual(a, b)
}

func (ld *LoopDetector) isErrorObservation(obs interface{}) bool {
	if o, ok := obs.(models.CmdOutputObservation); ok {
		// Heuristic: exit code != 0 OR content contains "error"
		if o.Metadata.ExitCode != 0 {
			return true
		}
		// String check?
		// if strings.Contains(strings.ToLower(o.Content), "error") { return true }
	}
	return false
}
