package events

// Handler represents an event handler interface
type Handler interface {
Handle(event *Event) error
EventType() string
}

// EventBus represents the event bus interface
type EventBus interface {
// Publish publishes an event to the bus
Publish(event *Event) error

// Subscribe subscribes a handler to events of a specific type
Subscribe(eventType string, handler Handler) error

// Unsubscribe removes a handler from events of a specific type
Unsubscribe(eventType string, handler Handler) error

// Start starts the event bus
Start() error

// Stop stops the event bus
Stop() error
}
