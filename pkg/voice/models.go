package voice

import (
	"time"
)

// Assistant represents a VAPI assistant
type Assistant struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	SystemPrompt string    `json:"systemPrompt,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

// Call represents a call made through VAPI
type Call struct {
	ID           string      `json:"id"`
	AssistantID  string      `json:"assistantId"`
	Status       string      `json:"status"`
	Duration     int         `json:"duration"`
	CreatedAt    time.Time   `json:"createdAt"`
	Customer     *Customer   `json:"customer,omitempty"`
	Analysis     *Analysis   `json:"analysis,omitempty"`
	Artifacts    []Artifact  `json:"artifacts,omitempty"`
	Transcript   interface{} `json:"transcript,omitempty"` // Can be []Message or string
	Messages     []Message   `json:"messages,omitempty"`
	Conversation []Message   `json:"conversation,omitempty"`
}

// Customer represents a customer in a VAPI call
type Customer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// Analysis represents the analysis of a VAPI call
type Analysis struct {
	Transcript []Message `json:"transcript,omitempty"`
}

// Artifact represents an artifact from a VAPI call
type Artifact struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Content    string    `json:"content,omitempty"`
	Transcript []Message `json:"transcript,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

// Message represents a message in a VAPI call transcript
type Message struct {
	Role    string `json:"role"`
	Text    string `json:"text,omitempty"`
	Content string `json:"content,omitempty"`
}

// File represents a file uploaded to VAPI
type File struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
}

// Tool represents a tool in VAPI
type Tool struct {
	ID             string          `json:"id"`
	Type           string          `json:"type"`
	Function       ToolFunction    `json:"function"`
	KnowledgeBases []KnowledgeBase `json:"knowledgeBases,omitempty"`
}

// ToolFunction represents a function in a VAPI tool
type ToolFunction struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// KnowledgeBase represents a knowledge base in a VAPI tool
type KnowledgeBase struct {
	Provider    string   `json:"provider"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	FileIDs     []string `json:"fileIds"`
}

// PhoneNumber represents a VAPI phone number
type PhoneNumber struct {
	ID          string `json:"id"`
	Number      string `json:"number"`
	AssistantID string `json:"assistantId,omitempty"`
}

// WebhookEvent represents a webhook event from VAPI
type WebhookEvent struct {
	Type      string      `json:"type"`
	Message   interface{} `json:"message"`
	Timestamp time.Time   `json:"timestamp"`
}

// EndOfCallReport represents an end-of-call-report event
type EndOfCallReport struct {
	Type        string    `json:"type"`
	Call        Call      `json:"call"`
	Transcript  []Message `json:"transcript"`
	Summary     string    `json:"summary,omitempty"`
	Analysis    *Analysis `json:"analysis,omitempty"`
	AssistantID string    `json:"assistantId"`
	CallID      string    `json:"callId"`
}

// ProcessedCall represents a processed call stored in the database
type ProcessedCall struct {
	ID          string    `json:"id"`
	CallID      string    `json:"call_id"`
	AssistantID string    `json:"assistant_id"`
	Transcript  []Message `json:"transcript"`
	Duration    int       `json:"duration"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UpdateAssistantRequest represents a request to update an assistant
type UpdateAssistantRequest struct {
	Name         *string `json:"name,omitempty"`
	SystemPrompt *string `json:"systemPrompt,omitempty"`
	ServerURL    *string `json:"serverUrl,omitempty"`
}

// CreateToolRequest represents a request to create a tool
type CreateToolRequest struct {
	Type           string          `json:"type"`
	Function       ToolFunction    `json:"function"`
	KnowledgeBases []KnowledgeBase `json:"knowledgeBases,omitempty"`
}

// AttachToolRequest represents a request to attach a tool to an assistant
type AttachToolRequest struct {
	ToolID string `json:"toolId"`
}

// Response represents a generic response from VAPI
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error response from VAPI
type ErrorResponse struct {
	Error string `json:"error"`
}
