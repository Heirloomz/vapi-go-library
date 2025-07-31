package voice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Client handles interactions with the VAPI API
type Client struct {
	apiToken   string
	baseURL    string
	httpClient *http.Client
	config     *Config
}

// Config represents configuration for the voice client
type Config struct {
	APIToken   string
	BaseURL    string
	Timeout    time.Duration
	CacheDir   string
	DebugDir   string
	StorageDir string
}

// NewClient creates a new VAPI client
func NewClient(config *Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.vapi.ai"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	// Create storage directories if they don't exist
	if config.StorageDir != "" {
		os.MkdirAll(config.StorageDir, os.ModePerm)
	}
	if config.CacheDir != "" {
		os.MkdirAll(config.CacheDir, os.ModePerm)
	}
	if config.DebugDir != "" {
		os.MkdirAll(config.DebugDir, os.ModePerm)
	}

	return &Client{
		apiToken:   config.APIToken,
		baseURL:    config.BaseURL,
		httpClient: &http.Client{Timeout: config.Timeout},
		config:     config,
	}
}

// getHeaders returns the headers for VAPI API requests
func (c *Client) getHeaders() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", c.apiToken),
		"Content-Type":  "application/json",
	}
}

// ListAssistants returns a list of VAPI assistants
func (c *Client) ListAssistants() ([]Assistant, error) {
	url := fmt.Sprintf("%s/assistant", c.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range c.getHeaders() {
		req.Header.Add(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error listing assistants: %s", string(body))
	}

	var assistants []Assistant
	if err := json.NewDecoder(resp.Body).Decode(&assistants); err != nil {
		return nil, err
	}

	return assistants, nil
}

// GetAssistant returns a VAPI assistant by ID
func (c *Client) GetAssistant(assistantID string) (*Assistant, error) {
	url := fmt.Sprintf("%s/assistant/%s", c.baseURL, assistantID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range c.getHeaders() {
		req.Header.Add(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error getting assistant: %s", string(body))
	}

	var assistant Assistant
	if err := json.NewDecoder(resp.Body).Decode(&assistant); err != nil {
		return nil, err
	}

	return &assistant, nil
}

// UpdateAssistant updates a VAPI assistant
func (c *Client) UpdateAssistant(assistantID string, updateReq *UpdateAssistantRequest) (*Assistant, error) {
	// First get the current assistant config
	url := fmt.Sprintf("%s/assistant/%s", c.baseURL, assistantID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range c.getHeaders() {
		req.Header.Add(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("failed to get assistant details: %s", string(body))
	}

	var assistantConfig map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&assistantConfig); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()

	// Update the system prompt if provided
	if updateReq.SystemPrompt != nil {
		// Update the system prompt in the model messages
		if model, ok := assistantConfig["model"].(map[string]interface{}); ok {
			if messages, ok := model["messages"].([]interface{}); ok && len(messages) > 0 {
				// Update the first system message
				if systemMsg, ok := messages[0].(map[string]interface{}); ok {
					if role, ok := systemMsg["role"].(string); ok && role == "system" {
						systemMsg["content"] = *updateReq.SystemPrompt
					}
				}
			} else {
				// Create messages array with system prompt if it doesn't exist
				model["messages"] = []interface{}{
					map[string]interface{}{
						"role":    "system",
						"content": *updateReq.SystemPrompt,
					},
				}
			}
		}
	}

	// Update server URL if provided
	if updateReq.ServerURL != nil {
		assistantConfig["serverUrl"] = *updateReq.ServerURL
	}

	// Remove read-only fields that shouldn't be included in the update
	delete(assistantConfig, "id")
	delete(assistantConfig, "createdAt")
	delete(assistantConfig, "updatedAt")
	delete(assistantConfig, "orgId")
	delete(assistantConfig, "isServerUrlSecretSet")

	// Update the assistant
	updateURL := fmt.Sprintf("%s/assistant/%s", c.baseURL, assistantID)
	updatePayloadBytes, err := json.Marshal(assistantConfig)
	if err != nil {
		return nil, err
	}

	updateReq2, err := http.NewRequest("PATCH", updateURL, bytes.NewBuffer(updatePayloadBytes))
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range c.getHeaders() {
		updateReq2.Header.Add(key, value)
	}

	updateResp, err := c.httpClient.Do(updateReq2)
	if err != nil {
		return nil, err
	}
	defer updateResp.Body.Close()

	if updateResp.StatusCode != http.StatusOK && updateResp.StatusCode != http.StatusCreated && updateResp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(updateResp.Body)
		return nil, fmt.Errorf("failed to update assistant: %s", string(body))
	}

	// Return the updated assistant
	return c.GetAssistant(assistantID)
}

// ListCalls returns a list of VAPI calls for an assistant
func (c *Client) ListCalls(assistantID string, limit int) ([]Call, error) {
	url := fmt.Sprintf("%s/call?assistantId=%s&limit=%d", c.baseURL, assistantID, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range c.getHeaders() {
		req.Header.Add(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error listing calls: %s", string(body))
	}

	var calls []Call
	if err := json.NewDecoder(resp.Body).Decode(&calls); err != nil {
		return nil, err
	}

	return calls, nil
}

// GetCall returns a VAPI call by ID
func (c *Client) GetCall(callID string) (*Call, error) {
	url := fmt.Sprintf("%s/call/%s", c.baseURL, callID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range c.getHeaders() {
		req.Header.Add(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error getting call: %s", string(body))
	}

	var call Call
	if err := json.NewDecoder(resp.Body).Decode(&call); err != nil {
		return nil, err
	}

	// Save to debug directory if configured
	if c.config.DebugDir != "" {
		callData, _ := json.MarshalIndent(call, "", "  ")
		os.WriteFile(fmt.Sprintf("%s/call_data_%s.json", c.config.DebugDir, callID), callData, 0644)
	}

	return &call, nil
}

// UploadFile uploads a file to VAPI
func (c *Client) UploadFile(filePath string) (*File, error) {
	fileName := filepath.Base(filePath)

	// Create a buffer to store the multipart form data
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)

	// Determine the MIME type
	mimeType := c.detectMimeType(filePath, fileName)

	// Add the content type field first
	err := multipartWriter.WriteField("contentType", mimeType)
	if err != nil {
		return nil, err
	}

	// Create a custom form file field with the correct Content-Type
	h := make(map[string][]string)
	h["Content-Disposition"] = []string{fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileName)}
	h["Content-Type"] = []string{mimeType}
	fileField, err := multipartWriter.CreatePart(h)
	if err != nil {
		return nil, err
	}

	// Open the file
	fileHandle, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fileHandle.Close()

	// Copy the file content to the form field
	_, err = io.Copy(fileField, fileHandle)
	if err != nil {
		return nil, err
	}

	// Close the multipart writer to finalize the form
	multipartWriter.Close()

	// Create the request
	url := fmt.Sprintf("%s/file", c.baseURL)
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return nil, err
	}

	// Set the content type with the boundary
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to upload file: %s", string(body))
	}

	// Parse the response
	var uploadedFile File
	if err := json.NewDecoder(resp.Body).Decode(&uploadedFile); err != nil {
		return nil, err
	}

	return &uploadedFile, nil
}

// CreateQueryTool creates a query tool for the knowledge base
func (c *Client) CreateQueryTool(fileIDs []string, toolName, description string) (*Tool, error) {
	payload := CreateToolRequest{
		Type: "query",
		Function: ToolFunction{
			Name: toolName,
		},
		KnowledgeBases: []KnowledgeBase{
			{
				Provider:    "google",
				Name:        toolName,
				Description: description,
				FileIDs:     fileIDs,
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Create the request
	url := fmt.Sprintf("%s/tool", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range c.getHeaders() {
		req.Header.Add(key, value)
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create query tool: %s", string(body))
	}

	// Parse the response
	var tool Tool
	if err := json.NewDecoder(resp.Body).Decode(&tool); err != nil {
		return nil, err
	}

	return &tool, nil
}

// AttachToolToAssistant attaches a tool to an assistant
func (c *Client) AttachToolToAssistant(assistantID, toolID string) error {
	// First get the current assistant config
	url := fmt.Sprintf("%s/assistant/%s", c.baseURL, assistantID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Add headers
	for key, value := range c.getHeaders() {
		req.Header.Add(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return fmt.Errorf("failed to get assistant details: %s", string(body))
	}

	var assistantConfig map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&assistantConfig); err != nil {
		resp.Body.Close()
		return err
	}
	resp.Body.Close()

	// Update the toolIds
	if _, ok := assistantConfig["model"]; !ok {
		assistantConfig["model"] = map[string]interface{}{}
	}

	model := assistantConfig["model"].(map[string]interface{})

	// Add the tool ID to the list, creating it if it doesn't exist
	var toolIDs []string
	if existingToolIDs, ok := model["toolIds"]; ok {
		for _, id := range existingToolIDs.([]interface{}) {
			toolIDs = append(toolIDs, id.(string))
		}
	}

	// Check if tool ID already exists
	toolExists := false
	for _, id := range toolIDs {
		if id == toolID {
			toolExists = true
			break
		}
	}

	// Add the tool ID if it doesn't exist
	if !toolExists {
		toolIDs = append(toolIDs, toolID)
		model["toolIds"] = toolIDs
	} else {
		// Tool already attached
		return nil
	}

	// Remove read-only fields that shouldn't be included in the update
	delete(assistantConfig, "id")
	delete(assistantConfig, "createdAt")
	delete(assistantConfig, "updatedAt")
	delete(assistantConfig, "orgId")
	delete(assistantConfig, "isServerUrlSecretSet")

	// Update the assistant
	updateURL := fmt.Sprintf("%s/assistant/%s", c.baseURL, assistantID)
	updatePayloadBytes, err := json.Marshal(assistantConfig)
	if err != nil {
		return err
	}

	updateReq, err := http.NewRequest("PATCH", updateURL, bytes.NewBuffer(updatePayloadBytes))
	if err != nil {
		return err
	}

	// Add headers
	for key, value := range c.getHeaders() {
		updateReq.Header.Add(key, value)
	}

	updateResp, err := c.httpClient.Do(updateReq)
	if err != nil {
		return err
	}
	defer updateResp.Body.Close()

	if updateResp.StatusCode != http.StatusOK && updateResp.StatusCode != http.StatusCreated && updateResp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(updateResp.Body)
		return fmt.Errorf("failed to update assistant: %s", string(body))
	}

	return nil
}

// detectMimeType detects the MIME type of a file
func (c *Client) detectMimeType(filePath string, fileName string) string {
	// First, try to detect based on file extension
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".md", ".markdown":
		return "text/markdown"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".csv":
		return "text/csv"
	case ".json":
		return "application/json"
	case ".yaml", ".yml":
		return "text/yaml"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		// Try to detect content type from file content
		if file, err := os.Open(filePath); err == nil {
			defer file.Close()
			// Read first 512 bytes for content detection
			buffer := make([]byte, 512)
			if n, err := file.Read(buffer); err == nil {
				detectedType := http.DetectContentType(buffer[:n])
				// Map detected types to supported VAPI types
				switch {
				case strings.HasPrefix(detectedType, "text/"):
					return "text/plain"
				case detectedType == "application/json":
					return "application/json"
				case detectedType == "application/pdf":
					return "application/pdf"
				default:
					// If it's a text-like file, default to text/plain
					if strings.Contains(detectedType, "text") {
						return "text/plain"
					}
				}
			}
		}
		// Default fallback
		return "text/plain"
	}
}

// ExtractTranscript extracts the transcript from a VAPI call
func (c *Client) ExtractTranscript(call *Call) []Message {
	// Check for transcript in analysis
	if call.Analysis != nil && call.Analysis.Transcript != nil && len(call.Analysis.Transcript) > 0 {
		return call.Analysis.Transcript
	}

	// Check for other transcript sources
	if call.Transcript != nil {
		// Check if transcript is a string
		if transcriptStr, ok := call.Transcript.(string); ok && transcriptStr != "" {
			return c.parseTranscriptContent(transcriptStr)
		}

		// Check if transcript is a slice of messages
		if transcriptMsgs, ok := call.Transcript.([]Message); ok && len(transcriptMsgs) > 0 {
			return transcriptMsgs
		}
	}

	if call.Messages != nil && len(call.Messages) > 0 {
		return call.Messages
	}

	if call.Conversation != nil && len(call.Conversation) > 0 {
		return call.Conversation
	}

	// Check nested in artifacts
	if call.Artifacts != nil {
		for _, artifact := range call.Artifacts {
			if artifact.Transcript != nil && len(artifact.Transcript) > 0 {
				return artifact.Transcript
			}

			if artifact.Content != "" {
				if strings.Contains(artifact.Content, "Transcript") ||
					strings.Contains(artifact.Content, "AI") ||
					strings.Contains(artifact.Content, "User") {
					return c.parseTranscriptContent(artifact.Content)
				}
			}
		}
	}

	return []Message{}
}

// parseTranscriptContent parses transcript content from a string
func (c *Client) parseTranscriptContent(content string) []Message {
	if content == "" {
		return []Message{}
	}

	lines := strings.Split(strings.TrimSpace(content), "\n")
	transcript := []Message{}

	currentRole := ""
	currentText := ""

	// Skip first line if it's just "Transcript"
	startIdx := 0
	if len(lines) > 0 && strings.Contains(lines[0], "Transcript") {
		startIdx = 1
	}

	for i := startIdx; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for new speaker
		if strings.HasPrefix(line, "AI") || strings.HasPrefix(line, "BOT") || strings.HasPrefix(line, "ASSISTANT") {
			// Save previous message if exists
			if currentRole != "" && currentText != "" {
				transcript = append(transcript, Message{
					Role: currentRole,
					Text: strings.TrimSpace(currentText),
				})
			}

			currentRole = "assistant"
			// Extract text after the speaker indicator
			parts := strings.SplitN(line, " ", 2)
			if len(parts) > 1 {
				currentText = strings.TrimSpace(parts[1])
			} else {
				currentText = ""
			}
		} else if strings.HasPrefix(line, "User") || strings.HasPrefix(line, "USER") || strings.HasPrefix(line, "CLIENT") {
			// Save previous message if exists
			if currentRole != "" && currentText != "" {
				transcript = append(transcript, Message{
					Role: currentRole,
					Text: strings.TrimSpace(currentText),
				})
			}

			currentRole = "user"
			// Extract text after the speaker indicator
			parts := strings.SplitN(line, " ", 2)
			if len(parts) > 1 {
				currentText = strings.TrimSpace(parts[1])
			} else {
				currentText = ""
			}
		} else {
			// Append to current text if we have a role
			if currentRole != "" {
				if currentText != "" {
					currentText += " " + line
				} else {
					currentText = line
				}
			}
		}
	}

	// Add the last message
	if currentRole != "" && currentText != "" {
		transcript = append(transcript, Message{
			Role: currentRole,
			Text: strings.TrimSpace(currentText),
		})
	}

	return transcript
}
