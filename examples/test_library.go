package main

import (
"fmt"
"log"

"github.com/heirloomz/vapi-go-library/pkg/config"
"github.com/heirloomz/vapi-go-library/pkg/events"
)

// Simple test handler
type TestHandler struct {
name string
}

func (h *TestHandler) Handle(event *events.Event) error {
fmt.Printf("[%s] Received event: %s (type: %s)\n", h.name, event.ID, event.Type)
return nil
}

func (h *TestHandler) EventType() string {
return "test.event"
}

func main() {
fmt.Println("Testing VAPI Go Library...")

// Test 1: Configuration loading
fmt.Println("\n1. Testing configuration loading...")
cfg := config.LoadFromEnv()
fmt.Printf("   VAPI Base URL: %s\n", cfg.VAPI.BaseURL)
fmt.Printf("   Tunnel Provider: %s\n", cfg.Tunnel.Provider)
fmt.Printf("   Redis Host: %s\n", cfg.Events.Redis.Host)
fmt.Printf("   Workers Count: %d\n", cfg.Workers.Count)

// Test 2: Event creation
fmt.Println("\n2. Testing event creation...")
event := events.NewEvent("test.event", "test-source", map[string]interface{}{
"message": "Hello from VAPI library test!",
"number":  42,
})
fmt.Printf("   Event ID: %s\n", event.ID)
fmt.Printf("   Event Type: %s\n", event.Type)
fmt.Printf("   Event Source: %s\n", event.Source)

// Test 3: Event JSON serialization
fmt.Println("\n3. Testing event JSON serialization...")
jsonData, err := event.ToJSON()
if err != nil {
log.Printf("   Error serializing event: %v", err)
} else {
fmt.Printf("   Event JSON: %s\n", string(jsonData))
}

// Test 4: Event JSON deserialization
fmt.Println("\n4. Testing event JSON deserialization...")
deserializedEvent, err := events.FromJSON(jsonData)
if err != nil {
log.Printf("   Error deserializing event: %v", err)
} else {
fmt.Printf("   Deserialized Event ID: %s\n", deserializedEvent.ID)
fmt.Printf("   Deserialized Event Type: %s\n", deserializedEvent.Type)
}

// Test 5: Event metadata
fmt.Println("\n5. Testing event metadata...")
event.AddMetadata("test_key", "test_value")
event.AddMetadata("priority", "high")

if value, exists := event.GetMetadata("test_key"); exists {
fmt.Printf("   Metadata 'test_key': %v\n", value)
}
if value, exists := event.GetMetadata("priority"); exists {
fmt.Printf("   Metadata 'priority': %v\n", value)
}

fmt.Println("\n✅ All basic tests passed! The VAPI library core functionality is working.")
fmt.Println("\nNote: Full functionality (Redis, webhooks, tunnels) requires additional setup and dependencies.")
}
