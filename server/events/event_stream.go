package events

import (
	"bufio"
	"encoding/json"
	"log/slog"
	"openhands-go/server/models"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type EventType string

const (
	EventTypeAction      EventType = "action"
	EventTypeObservation EventType = "observation"
)

type Event struct {
	ID        string      `json:"id"`
	Type      EventType   `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Content   interface{} `json:"content"`
	Source    string      `json:"source"`
}

// Custom unmarshalling to handle polymorphic Content
func (e *Event) UnmarshalJSON(data []byte) error {
	type Alias Event
	aux := &struct {
		Content json.RawMessage `json:"content"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if e.Type == EventTypeAction {
		// Inspect "action" field
		var probe struct {
			Action models.ActionType `json:"action"`
		}
		if err := json.Unmarshal(aux.Content, &probe); err != nil {
			// Fallback or error? For now fallback to generic map if probe fails (maybe message action string?)
			// But wait, MessageAction is a struct.
			// If it's just a string (legacy), we handle it.
			var str string
			if err := json.Unmarshal(aux.Content, &str); err == nil {
				e.Content = str // Legacy string content
				return nil
			}
			e.Content = aux.Content // Keep raw if unknown
			return nil
		}

		switch probe.Action {
		case models.ActionTypeCmdRun:
			var act models.CmdRunAction
			json.Unmarshal(aux.Content, &act)
			e.Content = act
		case models.ActionTypeAgentFinish:
			var act models.AgentFinishAction
			json.Unmarshal(aux.Content, &act)
			e.Content = act
		case models.ActionTypeMessage:
			var act models.MessageAction
			json.Unmarshal(aux.Content, &act)
			e.Content = act
		case models.ActionTypeDelegate:
			var act models.AgentDelegateAction
			json.Unmarshal(aux.Content, &act)
			e.Content = act
		default:
			// Try MessageAction as default if action is missing?
			// Or just generic map
			var m map[string]interface{}
			json.Unmarshal(aux.Content, &m)
			e.Content = m
		}
	} else if e.Type == EventTypeObservation {
		// Inspect "observation" field
		var probe struct {
			Observation string `json:"observation"`
		}
		if err := json.Unmarshal(aux.Content, &probe); err != nil {
			var str string
			if err := json.Unmarshal(aux.Content, &str); err == nil {
				e.Content = str
				return nil
			}
			e.Content = aux.Content
			return nil
		}

		if probe.Observation == "run" || probe.Observation == "run_ipython" || probe.Observation == "delegate" {
			var obs models.CmdOutputObservation
			json.Unmarshal(aux.Content, &obs)
			e.Content = obs
		} else if probe.Observation == "task_tracking" {
			var obs models.TaskTrackingObservation
			json.Unmarshal(aux.Content, &obs)
			e.Content = obs
		} else if probe.Observation == "loop_detection" {
			var obs models.LoopDetectionObservation
			json.Unmarshal(aux.Content, &obs)
			e.Content = obs
		} else {
			var m map[string]interface{}
			json.Unmarshal(aux.Content, &m)
			e.Content = m
		}
	} else {
		// Generic content
		var m interface{}
		json.Unmarshal(aux.Content, &m)
		e.Content = m
	}
	return nil
}

type EventStream struct {
	mu             sync.RWMutex
	events         []Event
	conversationID string
	filePath       string
	subscribers    []func(Event)
}

func NewEventStream(conversationID, filePath string) *EventStream {
	es := &EventStream{
		events:         make([]Event, 0),
		conversationID: conversationID,
		filePath:       filePath,
		subscribers:    make([]func(Event), 0),
	}
	es.loadEvents()
	return es
}

func (es *EventStream) loadEvents() {
	if es.filePath == "" {
		return
	}
	file, err := os.Open(es.filePath)
	if err != nil {
		return // File might not exist yet
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event Event
		// This uses UnmarshalJSON we defined above
		if err := json.Unmarshal(scanner.Bytes(), &event); err == nil {
			es.events = append(es.events, event)
		} else {
			slog.Debug("Error loading event", "error", err)
		}
	}

	if len(es.events) == 0 {
		slog.Debug("No events found for session", "sid", es.conversationID, "dir", es.filePath)
	}
}

func (es *EventStream) appendEventToFile(event Event) {
	if es.filePath == "" {
		return
	}

	// Ensure directory exists
	dir := filepath.Dir(es.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}

	file, err := os.OpenFile(es.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	data, err := json.Marshal(event)
	if err == nil {
		file.Write(data)
		file.Write([]byte("\n"))
	}
}

func (es *EventStream) Subscribe(callback func(Event)) {
	es.mu.Lock()
	defer es.mu.Unlock()
	es.subscribers = append(es.subscribers, callback)
}

func (es *EventStream) AddEvent(event Event) {
	es.mu.Lock()
	defer es.mu.Unlock()
	event.Timestamp = time.Now()

	// If Content is not already one of our types (e.g. map passed in),
	// we assume the caller passed correct struct OR we just store it.
	// But persistence uses JSON Marshal. When loading back, UnmarshalJSON will type it.
	// So in-memory, it might be untyped initially if added as map.
	// Ideally callers should add structs.

	es.events = append(es.events, event)
	es.appendEventToFile(event)

	// Notify subscribers
	for _, sub := range es.subscribers {
		// Run in goroutine to avoid blocking
		go sub(event)
	}
}

func (es *EventStream) GetEvents() []Event {
	es.mu.RLock()
	defer es.mu.RUnlock()

	// Return a copy
	events := make([]Event, len(es.events))
	copy(events, es.events)
	return events
}
