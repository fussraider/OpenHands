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
	mu            sync.RWMutex
	events        []Event
	conversationID string
	subscribers   []func(Event)
}

func NewEventStream(conversationID string) *EventStream {
	return &EventStream{
		events:         make([]Event, 0),
		conversationID: conversationID,
		subscribers:    make([]func(Event), 0),
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
