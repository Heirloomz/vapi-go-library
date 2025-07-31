package voice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/heirloomz/vapi-go-library/pkg/events"
)

// WebhookServer handles VAPI webhook events
type WebhookServer struct {
	port      int
	eventBus  events.EventBus
	processor *CallProcessor
	server    *http.Server
}

// NewWebhookServer creates a new webhook server
func NewWebhookServer(port int, eventBus events.EventBus, processor *CallProcessor) *WebhookServer {
	return &WebhookServer{
		port:      port,
		eventBus:  eventBus,
		processor: processor,
	}
}

// Start starts the webhook server
func (w *WebhookServer) Start() error {
	mux := http.NewServeMux()

	// VAPI webhook endpoint
	mux.HandleFunc("/webhooks/vapi", w.handleVAPIWebhook)
	mux.HandleFunc("/webhooks/voice", w.handleVoiceWebhook)
	mux.HandleFunc("/webhooks/health", w.handleHealthCheck)

	w.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", w.port),
		Handler: mux,
	}

	go func() {
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Log error but don't panic - this will be handled by the caller
		}
	}()

	return nil
}

// Stop stops the webhook server
func (w *WebhookServer) Stop() error {
	if w.server != nil {
		return w.server.Close()
	}
	return nil
}

// handleVAPIWebhook handles VAPI webhook events
func (w *WebhookServer) handleVAPIWebhook(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body := make([]byte, req.ContentLength)
	_, err := req.Body.Read(body)
	if err != nil {
		http.Error(rw, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Process the webhook event
	if err := w.processWebhookEvent(body); err != nil {
		http.Error(rw, "Failed to process webhook event", http.StatusInternalServerError)
		return
	}

	// Respond with success
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))
}

// handleVoiceWebhook handles generic voice webhook events
func (w *WebhookServer) handleVoiceWebhook(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body := make([]byte, req.ContentLength)
	_, err := req.Body.Read(body)
	if err != nil {
		http.Error(rw, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Process the webhook event
	if err := w.processWebhookEvent(body); err != nil {
		http.Error(rw, "Failed to process webhook event", http.StatusInternalServerError)
		return
	}

	// Respond with success
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))
}

// handleHealthCheck handles health check requests
func (w *WebhookServer) handleHealthCheck(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))
}

// processWebhookEvent processes a webhook event
func (w *WebhookServer) processWebhookEvent(payload []byte) error {
	// Parse the webhook payload
	var webhookData map[string]interface{}
	if err := json.Unmarshal(payload, &webhookData); err != nil {
		return fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	// Extract the message
	message, ok := webhookData["message"].(map[string]interface{})
	if !ok {
		// No message field, skip processing
		return nil
	}

	// Check if this is an end-of-call-report event
	eventType, ok := message["type"].(string)
	if !ok || eventType != "end-of-call-report" {
		// Not an end-of-call-report event, skip processing
		return nil
	}

	// Process the end-of-call-report event
	if w.processor != nil {
		return w.processor.ProcessEndOfCallReport(message)
	}

	// Publish raw webhook event to event bus
	if w.eventBus != nil {
		event := events.NewEvent(events.EventWebhookReceived, "vapi-webhook", webhookData)
		return w.eventBus.Publish(event)
	}

	return nil
}

// CallProcessor handles processing of call events
type CallProcessor struct {
	client   *Client
	eventBus events.EventBus
}

// NewCallProcessor creates a new call processor
func NewCallProcessor(client *Client, eventBus events.EventBus) *CallProcessor {
	return &CallProcessor{
		client:   client,
		eventBus: eventBus,
	}
}

// ProcessEndOfCallReport processes an end-of-call-report event
func (p *CallProcessor) ProcessEndOfCallReport(message map[string]interface{}) error {
	// Extract call information
	callData, ok := message["call"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no call data in end-of-call-report")
	}

	callID, ok := callData["id"].(string)
	if !ok {
		return fmt.Errorf("no call ID in end-of-call-report")
	}

	assistantID, ok := callData["assistantId"].(string)
	if !ok {
		return fmt.Errorf("no assistant ID in end-of-call-report")
	}

	// Get full call details from VAPI API
	call, err := p.client.GetCall(callID)
	if err != nil {
		return fmt.Errorf("failed to get call details: %w", err)
	}

	// Extract transcript
	transcript := p.client.ExtractTranscript(call)

	// Create processed call
	processedCall := &ProcessedCall{
		ID:          fmt.Sprintf("processed_%s", callID),
		CallID:      callID,
		AssistantID: assistantID,
		Transcript:  transcript,
		Duration:    call.Duration,
		Status:      call.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Publish call-completed event
	if p.eventBus != nil {
		event := events.NewEvent(events.EventCallCompleted, "vapi-processor", processedCall)
		if err := p.eventBus.Publish(event); err != nil {
			return fmt.Errorf("failed to publish call-completed event: %w", err)
		}
	}

	return nil
}
