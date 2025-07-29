package events

import (
"encoding/json"
"time"
)

// Event represents a generic event in the VAPI library
type Event struct {
ID        string                 `json:"id"`
Type      string                 `json:"type"`
Timestamp time.Time              `json:"timestamp"`
Source    string                 `json:"source"`
Data      interface{}            `json:"data"`
Metadata  map[string]interface{} `json:"metadata"`
}

// Event types constants
const (
EventCallCompleted     = "vapi.call.completed"
EventCallStarted       = "vapi.call.started"
EventTranscriptUpdate  = "vapi.transcript.update"
EventAssistantUpdated  = "vapi.assistant.updated"
EventFileUploaded      = "vapi.file.uploaded"
EventToolCreated       = "vapi.tool.created"
EventWebhookReceived   = "vapi.webhook.received"
)

// NewEvent creates a new event with the given parameters
func NewEvent(eventType, source string, data interface{}) *Event {
return &Event{
ID:        generateEventID(),
Type:      eventType,
Timestamp: time.Now(),
Source:    source,
Data:      data,
Metadata:  make(map[string]interface{}),
}
}

// ToJSON converts the event to JSON bytes
func (e *Event) ToJSON() ([]byte, error) {
return json.Marshal(e)
}

// FromJSON creates an event from JSON bytes
func FromJSON(data []byte) (*Event, error) {
var event Event
err := json.Unmarshal(data, &event)
return &event, err
}

// AddMetadata adds metadata to the event
func (e *Event) AddMetadata(key string, value interface{}) {
if e.Metadata == nil {
e.Metadata = make(map[string]interface{})
}
e.Metadata[key] = value
}

// GetMetadata retrieves metadata from the event
func (e *Event) GetMetadata(key string) (interface{}, bool) {
if e.Metadata == nil {
return nil, false
}
value, exists := e.Metadata[key]
return value, exists
}

// generateEventID generates a unique event ID
func generateEventID() string {
// Simple timestamp-based ID for now
// In production, you might want to use UUID
return time.Now().Format("20060102150405.000000")
}
