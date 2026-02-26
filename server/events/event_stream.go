package events

import (
	"bufio"
	"encoding/json"
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
		if err := json.Unmarshal(scanner.Bytes(), &event); err == nil {
			es.events = append(es.events, event)
		}
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
