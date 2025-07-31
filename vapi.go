package vapi

import (
	"fmt"

	"github.com/heirloomz/vapi-go-library/pkg/chat"
	"github.com/heirloomz/vapi-go-library/pkg/config"
	"github.com/heirloomz/vapi-go-library/pkg/events"
)

// Library represents the main VAPI library
type Library struct {
	config     *config.Config
	eventBus   events.EventBus
	chatClient *chat.Client
	running    bool
}

// New creates a new VAPI library instance
func New(cfg *config.Config) (*Library, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Initialize chat client
	chatClient := chat.NewClient(cfg)

	return &Library{
		config:     cfg,
		chatClient: chatClient,
		running:    false,
	}, nil
}

// Start starts the VAPI library services
func (l *Library) Start() error {
	if l.running {
		return fmt.Errorf("library is already running")
	}

	// TODO: Initialize event bus, webhook server, tunnel manager, etc.
	l.running = true
	return nil
}

// Stop stops the VAPI library services
func (l *Library) Stop() error {
	if !l.running {
		return fmt.Errorf("library is not running")
	}

	// TODO: Stop all services gracefully
	l.running = false
	return nil
}

// IsRunning returns whether the library is currently running
func (l *Library) IsRunning() bool {
	return l.running
}

// EventBus returns the event bus instance
func (l *Library) EventBus() events.EventBus {
	return l.eventBus
}

// Chat returns the chat client instance
func (l *Library) Chat() *chat.Client {
	return l.chatClient
}

// Config returns the library configuration
func (l *Library) Config() *config.Config {
	return l.config
}
