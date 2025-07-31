package vapi

import (
	"fmt"

	"github.com/heirloomz/vapi-go-library/pkg/chat"
	"github.com/heirloomz/vapi-go-library/pkg/config"
	"github.com/heirloomz/vapi-go-library/pkg/events"
	"github.com/heirloomz/vapi-go-library/pkg/voice"
)

// Library represents the main VAPI library
type Library struct {
	config      *config.Config
	eventBus    events.EventBus
	chatClient  *chat.Client
	voiceClient *voice.VoiceClient
	running     bool
}

// New creates a new VAPI library instance
func New(cfg *config.Config) (*Library, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Initialize event bus
	eventBus, err := events.NewEventBus(cfg.Events.Backend, events.RedisConfig{
		Host:     cfg.Events.Redis.Host,
		Port:     cfg.Events.Redis.Port,
		Password: cfg.Events.Redis.Password,
		DB:       cfg.Events.Redis.DB,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create event bus: %w", err)
	}

	// Initialize chat client
	chatClient := chat.NewClient(cfg)

	// Initialize voice client
	voiceClient, err := voice.NewVoiceClient(cfg, eventBus)
	if err != nil {
		return nil, fmt.Errorf("failed to create voice client: %w", err)
	}

	return &Library{
		config:      cfg,
		eventBus:    eventBus,
		chatClient:  chatClient,
		voiceClient: voiceClient,
		running:     false,
	}, nil
}

// Start starts the VAPI library services
func (l *Library) Start() error {
	if l.running {
		return fmt.Errorf("library is already running")
	}

	// Start event bus
	if err := l.eventBus.Start(); err != nil {
		return fmt.Errorf("failed to start event bus: %w", err)
	}

	// Start voice client
	if err := l.voiceClient.Start(); err != nil {
		l.eventBus.Stop() // Clean up event bus on failure
		return fmt.Errorf("failed to start voice client: %w", err)
	}

	l.running = true
	return nil
}

// Stop stops the VAPI library services
func (l *Library) Stop() error {
	if !l.running {
		return fmt.Errorf("library is not running")
	}

	// Stop voice client
	if err := l.voiceClient.Stop(); err != nil {
		return fmt.Errorf("failed to stop voice client: %w", err)
	}

	// Stop event bus
	if err := l.eventBus.Stop(); err != nil {
		return fmt.Errorf("failed to stop event bus: %w", err)
	}

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

// Voice returns the voice client instance
func (l *Library) Voice() *voice.VoiceClient {
	return l.voiceClient
}

// Config returns the library configuration
func (l *Library) Config() *config.Config {
	return l.config
}
