package voice

import (
	"fmt"

	"github.com/heirloomz/vapi-go-library/pkg/config"
	"github.com/heirloomz/vapi-go-library/pkg/events"
)

// VoiceClient provides voice functionality for the VAPI library
type VoiceClient struct {
	client        *Client
	webhookServer *WebhookServer
	processor     *CallProcessor
	eventBus      events.EventBus
	config        *config.Config
}

// NewVoiceClient creates a new voice client
func NewVoiceClient(cfg *config.Config, eventBus events.EventBus) (*VoiceClient, error) {
	// Create voice client config
	voiceConfig := &Config{
		APIToken:   cfg.VAPI.APIToken,
		BaseURL:    cfg.VAPI.BaseURL,
		Timeout:    cfg.VAPI.Timeout,
		StorageDir: "./vapi_storage",
		CacheDir:   "./vapi_cache",
		DebugDir:   "./vapi_debug",
	}

	// Create VAPI client
	client := NewClient(voiceConfig)

	// Create call processor
	processor := NewCallProcessor(client, eventBus)

	// Create webhook server
	webhookServer := NewWebhookServer(cfg.Tunnel.Port, eventBus, processor)

	return &VoiceClient{
		client:        client,
		webhookServer: webhookServer,
		processor:     processor,
		eventBus:      eventBus,
		config:        cfg,
	}, nil
}

// Start starts the voice client services
func (v *VoiceClient) Start() error {
	// Start webhook server
	if err := v.webhookServer.Start(); err != nil {
		return fmt.Errorf("failed to start webhook server: %w", err)
	}

	return nil
}

// Stop stops the voice client services
func (v *VoiceClient) Stop() error {
	// Stop webhook server
	if err := v.webhookServer.Stop(); err != nil {
		return fmt.Errorf("failed to stop webhook server: %w", err)
	}

	return nil
}

// ListAssistants returns a list of VAPI assistants
func (v *VoiceClient) ListAssistants() ([]Assistant, error) {
	return v.client.ListAssistants()
}

// GetAssistant returns a VAPI assistant by ID
func (v *VoiceClient) GetAssistant(assistantID string) (*Assistant, error) {
	return v.client.GetAssistant(assistantID)
}

// UpdateAssistant updates a VAPI assistant
func (v *VoiceClient) UpdateAssistant(assistantID string, updateReq *UpdateAssistantRequest) (*Assistant, error) {
	return v.client.UpdateAssistant(assistantID, updateReq)
}

// ListCalls returns a list of VAPI calls for an assistant
func (v *VoiceClient) ListCalls(assistantID string, limit int) ([]Call, error) {
	return v.client.ListCalls(assistantID, limit)
}

// GetCall returns a VAPI call by ID
func (v *VoiceClient) GetCall(callID string) (*Call, error) {
	return v.client.GetCall(callID)
}

// UploadFile uploads a file to VAPI
func (v *VoiceClient) UploadFile(filePath string) (*File, error) {
	return v.client.UploadFile(filePath)
}

// CreateQueryTool creates a query tool for the knowledge base
func (v *VoiceClient) CreateQueryTool(fileIDs []string, toolName, description string) (*Tool, error) {
	return v.client.CreateQueryTool(fileIDs, toolName, description)
}

// AttachToolToAssistant attaches a tool to an assistant
func (v *VoiceClient) AttachToolToAssistant(assistantID, toolID string) error {
	return v.client.AttachToolToAssistant(assistantID, toolID)
}

// ExtractTranscript extracts the transcript from a VAPI call
func (v *VoiceClient) ExtractTranscript(call *Call) []Message {
	return v.client.ExtractTranscript(call)
}
