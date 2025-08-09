package chat

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/heirloomz/vapi-go-library/pkg/config"
)

// Client represents a VAPI chat client
type Client struct {
	config     *config.Config
	httpClient *http.Client
}

// NewClient creates a new VAPI chat client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.VAPI.Timeout,
		},
	}
}

// CreateChat creates a new chat with the VAPI API
func (c *Client) CreateChat(ctx context.Context, req *CreateChatRequest) (*ChatResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	if req.Input == nil {
		return nil, fmt.Errorf("input is required")
	}

	// Validate that at least one of assistantId, assistant, sessionId, or previousChatId is provided
	if req.AssistantID == nil && req.Assistant == nil && req.SessionID == nil && req.PreviousChatID == nil {
		return nil, fmt.Errorf("at least one of assistantId, assistant, sessionId, or previousChatId is required")
	}

	// Validate that sessionId and previousChatId are mutually exclusive
	if req.SessionID != nil && req.PreviousChatID != nil {
		return nil, fmt.Errorf("sessionId and previousChatId are mutually exclusive")
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/chat", c.config.VAPI.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.VAPI.APIToken)

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var chatResponse ChatResponse
	if err := json.Unmarshal(body, &chatResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &chatResponse, nil
}

// CreateStreamingChat creates a new streaming chat with the VAPI API
func (c *Client) CreateStreamingChat(ctx context.Context, req *CreateChatRequest) (<-chan *StreamingChatResponse, <-chan error) {
	responseChan := make(chan *StreamingChatResponse, 100)
	errorChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errorChan)

		if req == nil {
			errorChan <- fmt.Errorf("request cannot be nil")
			return
		}

		if req.Input == nil {
			errorChan <- fmt.Errorf("input is required")
			return
		}

		// Validate that at least one of assistantId, assistant, sessionId, or previousChatId is provided
		if req.AssistantID == nil && req.Assistant == nil && req.SessionID == nil && req.PreviousChatID == nil {
			errorChan <- fmt.Errorf("at least one of assistantId, assistant, sessionId, or previousChatId is required")
			return
		}

		// Validate that sessionId and previousChatId are mutually exclusive
		if req.SessionID != nil && req.PreviousChatID != nil {
			errorChan <- fmt.Errorf("sessionId and previousChatId are mutually exclusive")
			return
		}

		// Enable streaming
		streamReq := *req
		streamReq.Stream = &[]bool{true}[0]

		// Marshal request to JSON
		jsonData, err := json.Marshal(&streamReq)
		if err != nil {
			errorChan <- fmt.Errorf("failed to marshal request: %w", err)
			return
		}

		// Create HTTP request
		url := fmt.Sprintf("%s/chat", c.config.VAPI.BaseURL)
		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			errorChan <- fmt.Errorf("failed to create HTTP request: %w", err)
			return
		}

		// Set headers
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+c.config.VAPI.APIToken)
		httpReq.Header.Set("Accept", "text/event-stream")

		// Send request
		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			errorChan <- fmt.Errorf("failed to send request: %w", err)
			return
		}
		defer resp.Body.Close()

		// Check for HTTP errors
		if resp.StatusCode >= 400 {
			body, _ := io.ReadAll(resp.Body)
			errorChan <- fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
			return
		}

		// Process streaming response
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			// Skip empty lines and comments
			if line == "" || strings.HasPrefix(line, ":") {
				continue
			}

			// Parse Server-Sent Events format
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")

				// Skip keep-alive messages
				if data == "" || data == "[DONE]" {
					continue
				}

				// Parse JSON data
				var streamResponse StreamingChatResponse
				if err := json.Unmarshal([]byte(data), &streamResponse); err != nil {
					errorChan <- fmt.Errorf("failed to parse streaming response: %w", err)
					return
				}

				// Send response to channel
				select {
				case responseChan <- &streamResponse:
				case <-ctx.Done():
					return
				}

				// Check if streaming is done
				if streamResponse.Done {
					return
				}
			}
		}

		if err := scanner.Err(); err != nil {
			errorChan <- fmt.Errorf("error reading streaming response: %w", err)
		}
	}()

	return responseChan, errorChan
}

// CreateChatWithText is a convenience method to create a chat with simple text input
func (c *Client) CreateChatWithText(ctx context.Context, text string, assistantID *string) (*ChatResponse, error) {
	req := &CreateChatRequest{
		Input:       text,
		AssistantID: assistantID,
	}

	return c.CreateChat(ctx, req)
}

// CreateChatWithMessages is a convenience method to create a chat with message history
func (c *Client) CreateChatWithMessages(ctx context.Context, messages []ChatMessage, assistantID *string) (*ChatResponse, error) {
	req := &CreateChatRequest{
		Input:       messages,
		AssistantID: assistantID,
	}

	return c.CreateChat(ctx, req)
}

// CreateChatWithAssistant is a convenience method to create a chat with a custom assistant
func (c *Client) CreateChatWithAssistant(ctx context.Context, text string, assistant *Assistant) (*ChatResponse, error) {
	req := &CreateChatRequest{
		Input:     text,
		Assistant: assistant,
	}

	return c.CreateChat(ctx, req)
}

// ContinueChat continues a chat from a previous chat ID
func (c *Client) ContinueChat(ctx context.Context, text string, previousChatID string) (*ChatResponse, error) {
	req := &CreateChatRequest{
		Input:          text,
		PreviousChatID: &previousChatID,
	}

	return c.CreateChat(ctx, req)
}

// CreateSessionChat creates a chat within a session
func (c *Client) CreateSessionChat(ctx context.Context, text string, sessionID string) (*ChatResponse, error) {
	req := &CreateChatRequest{
		Input:     text,
		SessionID: &sessionID,
	}

	return c.CreateChat(ctx, req)
}

// CreateStreamingChatWithText is a convenience method to create a streaming chat with simple text input
func (c *Client) CreateStreamingChatWithText(ctx context.Context, text string, assistantID *string) (<-chan *StreamingChatResponse, <-chan error) {
	req := &CreateChatRequest{
		Input:       text,
		AssistantID: assistantID,
	}

	return c.CreateStreamingChat(ctx, req)
}

// CreateStreamingChatWithAssistant is a convenience method to create a streaming chat with a custom assistant
func (c *Client) CreateStreamingChatWithAssistant(ctx context.Context, text string, assistant *Assistant) (<-chan *StreamingChatResponse, <-chan error) {
	req := &CreateChatRequest{
		Input:     text,
		Assistant: assistant,
	}

	return c.CreateStreamingChat(ctx, req)
}

// ValidateRequest validates a CreateChatRequest
func (c *Client) ValidateRequest(req *CreateChatRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if req.Input == nil {
		return fmt.Errorf("input is required")
	}

	// Validate that at least one of assistantId, assistant, sessionId, or previousChatId is provided
	if req.AssistantID == nil && req.Assistant == nil && req.SessionID == nil && req.PreviousChatID == nil {
		return fmt.Errorf("at least one of assistantId, assistant, sessionId, or previousChatId is required")
	}

	// Validate that sessionId and previousChatId are mutually exclusive
	if req.SessionID != nil && req.PreviousChatID != nil {
		return fmt.Errorf("sessionId and previousChatId are mutually exclusive")
	}

	// Validate name length if provided
	if req.Name != nil && len(*req.Name) > 40 {
		return fmt.Errorf("name must be 40 characters or less")
	}

	return nil
}

// SetTimeout sets a custom timeout for the HTTP client
func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() *config.Config {
	return c.config
}

// CreateSession creates a new VAPI session for the given assistant
func (c *Client) CreateSession(ctx context.Context, assistantID string) (*SessionResponse, error) {
	if assistantID == "" {
		return nil, fmt.Errorf("assistantID is required")
	}

	// Create session request payload
	sessionRequest := map[string]string{
		"assistantId": assistantID,
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(sessionRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal session request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/session", c.config.VAPI.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.VAPI.APIToken)

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var sessionResponse SessionResponse
	if err := json.Unmarshal(body, &sessionResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &sessionResponse, nil
}
