package main

import (
	"context"
	"log"

	"github.com/heirloomz/vapi-go-library"
	"github.com/heirloomz/vapi-go-library/pkg/chat"
	"github.com/heirloomz/vapi-go-library/pkg/config"
	"github.com/heirloomz/vapi-go-library/pkg/events"
)

// ExampleCallHandler demonstrates how to handle call-completed events
type ExampleCallHandler struct {
	name string
}

func (h *ExampleCallHandler) Handle(event *events.Event) error {
	log.Printf("[%s] Received event: %s", h.name, event.Type)

	// Process the call-completed event
	if event.Type == events.EventCallCompleted {
		log.Printf("[%s] Processing call completion event", h.name)
		// Here you would typically:
		// 1. Extract call data from event.Data
		// 2. Generate stories from transcripts
		// 3. Update knowledge bases
		// 4. Send notifications
	}

	return nil
}

func (h *ExampleCallHandler) EventType() string {
	return events.EventCallCompleted
}

func main() {
	// Load configuration from environment variables
	cfg := config.LoadFromEnv()

	// Create library instance
	library, err := vapi.New(cfg)
	if err != nil {
		log.Fatal("Failed to create VAPI library:", err)
	}

	// Register event handlers for call processing
	callHandler := &ExampleCallHandler{name: "StoryGenerator"}
	if err := library.EventBus().Subscribe(events.EventCallCompleted, callHandler); err != nil {
		log.Fatal("Failed to subscribe to call events:", err)
	}

	// Start the library (starts webhook server, event bus, etc.)
	if err := library.Start(); err != nil {
		log.Fatal("Failed to start VAPI library:", err)
	}
	defer library.Stop()

	log.Println("VAPI library started successfully!")
	log.Printf("Webhook server running on port %d", cfg.Tunnel.Port)
	log.Println("Event handlers registered for call processing")

	// Example: List assistants
	assistants, err := library.Voice().ListAssistants()
	if err != nil {
		log.Printf("Failed to list assistants: %v", err)
	} else {
		log.Printf("Found %d assistants", len(assistants))
		for _, assistant := range assistants {
			log.Printf("- %s (ID: %s)", assistant.Name, assistant.ID)
		}
	}

	// Example: Get calls for the first assistant
	if len(assistants) > 0 {
		calls, err := library.Voice().ListCalls(assistants[0].ID, 10)
		if err != nil {
			log.Printf("Failed to list calls: %v", err)
		} else {
			log.Printf("Found %d calls for assistant %s", len(calls), assistants[0].Name)
		}
	}

	// Example: Chat functionality (preserved from original library)
	chatReq := &chat.CreateChatRequest{
		Input:       "Hello, how are you?",
		AssistantID: &assistants[0].ID,
	}

	chatResp, err := library.Chat().CreateChat(context.Background(), chatReq)
	if err != nil {
		log.Printf("Chat failed: %v", err)
	} else {
		log.Printf("Chat response: %s", chatResp.Message)
	}

	// Keep the server running
	log.Println("Server running... Press Ctrl+C to stop")
	select {}
}
