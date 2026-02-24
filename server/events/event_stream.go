package events

import (
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
	mu     sync.RWMutex
	events []Event
}

func NewEventStream() *EventStream {
	return &EventStream{
		events: make([]Event, 0),
	}
}

func (es *EventStream) AddEvent(event Event) {
	es.mu.Lock()
	defer es.mu.Unlock()
	event.Timestamp = time.Now()
	es.events = append(es.events, event)
}

func (es *EventStream) GetEvents() []Event {
	es.mu.RLock()
	defer es.mu.RUnlock()

	// Return a copy
	events := make([]Event, len(es.events))
	copy(events, es.events)
	return events
}
