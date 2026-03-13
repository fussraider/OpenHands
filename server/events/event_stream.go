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
	Type      EventType   `json:"-"` // Not transmitted at top level by python
	Timestamp time.Time   `json:"timestamp"`
	Content   interface{} `json:"-"`
	Source    string      `json:"source"`
}

// MarshalJSON implements custom serialization to match Python backend structure.
func (e Event) MarshalJSON() ([]byte, error) {
	// Base structure
	out := map[string]interface{}{
		"id":        e.ID,
		"timestamp": e.Timestamp.Format(time.RFC3339Nano),
		"source":    e.Source,
	}

	contentBytes, err := json.Marshal(e.Content)
	if err == nil {
		var contentMap map[string]interface{}
		json.Unmarshal(contentBytes, &contentMap)

		if e.Type == EventTypeAction {
			// Extract "action" mapping it to top-level, put rest in "args"
			if actionVal, ok := contentMap["action"]; ok {
				out["action"] = actionVal
				delete(contentMap, "action")
			}
			out["args"] = contentMap
		} else if e.Type == EventTypeObservation {
			// Extract "observation", "content", put rest in "extras"
			if obsVal, ok := contentMap["observation"]; ok {
				out["observation"] = obsVal
				delete(contentMap, "observation")
			}
			if cVal, ok := contentMap["content"]; ok {
				out["content"] = cVal
				delete(contentMap, "content")
			}
			out["extras"] = contentMap
		} else {
			// Fallback
			for k, v := range contentMap {
				out[k] = v
			}
		}
	}
	return json.Marshal(out)
}

// Custom unmarshalling to handle polymorphic Content
func (e *Event) UnmarshalJSON(data []byte) error {
	type Alias Event
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var probe map[string]interface{}
	if err := json.Unmarshal(data, &probe); err != nil {
		return err
	}

	if actionName, ok := probe["action"].(string); ok {
		e.Type = EventTypeAction
		var argsData []byte
		if argsMap, ok := probe["args"]; ok {
			argsData, _ = json.Marshal(argsMap)
		} else {
			// fallback check if args are flattened
			argsData = data
		}

		switch models.ActionType(actionName) {
		case models.ActionTypeCmdRun:
			var act models.CmdRunAction
			json.Unmarshal(argsData, &act)
			act.Action = models.ActionTypeCmdRun
			e.Content = act
		case models.ActionTypeAgentFinish:
			var act models.AgentFinishAction
			json.Unmarshal(argsData, &act)
			act.Action = models.ActionTypeAgentFinish
			e.Content = act
		case models.ActionTypeMessage:
			var act models.MessageAction
			json.Unmarshal(argsData, &act)
			act.Action = models.ActionTypeMessage
			e.Content = act
		case models.ActionTypeDelegate:
			var act models.AgentDelegateAction
			json.Unmarshal(argsData, &act)
			act.Action = models.ActionTypeDelegate
			e.Content = act
		default:
			var m map[string]interface{}
			json.Unmarshal(argsData, &m)
			m["action"] = actionName
			e.Content = m
		}
	} else if obsName, ok := probe["observation"].(string); ok {
		e.Type = EventTypeObservation
		var extrasData []byte
		if extrasMap, ok := probe["extras"]; ok {
			extrasData, _ = json.Marshal(extrasMap)
		} else {
			extrasData = data
		}

		contentStr, _ := probe["content"].(string)

		switch obsName {
		case "run", "run_ipython", "delegate":
			var obs models.CmdOutputObservation
			json.Unmarshal(extrasData, &obs)
			obs.Observation = obsName
			obs.Content = contentStr
			e.Content = obs
		case "task_tracking":
			var obs models.TaskTrackingObservation
			json.Unmarshal(extrasData, &obs)
			obs.Observation = obsName
			obs.Content = contentStr
			e.Content = obs
		case "loop_detection":
			var obs models.LoopDetectionObservation
			json.Unmarshal(extrasData, &obs)
			obs.Observation = obsName
			obs.Content = contentStr
			e.Content = obs
		default:
			var m map[string]interface{}
			json.Unmarshal(extrasData, &m)
			m["observation"] = obsName
			m["content"] = contentStr
			e.Content = m
		}
	} else {
		// Generic or unknown
		e.Content = probe
	}

	return nil
}

type EventStream struct {
	mu             sync.RWMutex
	events         []Event
	conversationID string
	filePath       string
	subscribers    map[int]func(Event)
	nextSubID      int
}

func NewEventStream(conversationID, filePath string) *EventStream {
	es := &EventStream{
		events:         make([]Event, 0),
		conversationID: conversationID,
		filePath:       filePath,
		subscribers:    make(map[int]func(Event)),
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

// Subscribe registers a callback for new events and returns an unsubscribe function.
func (es *EventStream) Subscribe(callback func(Event)) func() {
	es.mu.Lock()
	id := es.nextSubID
	es.nextSubID++
	es.subscribers[id] = callback
	es.mu.Unlock()

	return func() {
		es.mu.Lock()
		defer es.mu.Unlock()
		delete(es.subscribers, id)
	}
}

func (es *EventStream) AddEvent(event Event) {
	es.mu.Lock()
	defer es.mu.Unlock()
	event.Timestamp = time.Now()

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
