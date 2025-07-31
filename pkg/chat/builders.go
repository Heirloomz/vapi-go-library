package chat

import (
	"fmt"
)

// AssistantBuilder helps build Assistant configurations
type AssistantBuilder struct {
	assistant *Assistant
}

// NewAssistantBuilder creates a new AssistantBuilder
func NewAssistantBuilder() *AssistantBuilder {
	return &AssistantBuilder{
		assistant: &Assistant{},
	}
}

// WithModel sets the model configuration
func (b *AssistantBuilder) WithModel(provider, model string) *AssistantBuilder {
	b.assistant.Model = &Model{
		Provider: provider,
		Model:    model,
	}
	return b
}

// WithModelMessages sets the model messages
func (b *AssistantBuilder) WithModelMessages(messages []ModelMessage) *AssistantBuilder {
	if b.assistant.Model == nil {
		b.assistant.Model = &Model{}
	}
	b.assistant.Model.Messages = messages
	return b
}

// WithSystemMessage adds a system message to the model
func (b *AssistantBuilder) WithSystemMessage(content string) *AssistantBuilder {
	if b.assistant.Model == nil {
		b.assistant.Model = &Model{}
	}

	systemMessage := ModelMessage{
		Role:    "system",
		Content: content,
	}

	b.assistant.Model.Messages = append(b.assistant.Model.Messages, systemMessage)
	return b
}

// WithAssistantMessage adds an assistant message to the model
func (b *AssistantBuilder) WithAssistantMessage(content string) *AssistantBuilder {
	if b.assistant.Model == nil {
		b.assistant.Model = &Model{}
	}

	assistantMessage := ModelMessage{
		Role:    "assistant",
		Content: content,
	}

	b.assistant.Model.Messages = append(b.assistant.Model.Messages, assistantMessage)
	return b
}

// WithVoice sets the voice configuration
func (b *AssistantBuilder) WithVoice(provider, voiceID string) *AssistantBuilder {
	b.assistant.Voice = &Voice{
		Provider: provider,
		VoiceID:  voiceID,
	}
	return b
}

// WithFirstMessage sets the first message
func (b *AssistantBuilder) WithFirstMessage(message string) *AssistantBuilder {
	b.assistant.FirstMessage = &message
	return b
}

// WithFirstMessageMode sets the first message mode
func (b *AssistantBuilder) WithFirstMessageMode(mode string) *AssistantBuilder {
	b.assistant.FirstMessageMode = &mode
	return b
}

// WithMaxDuration sets the maximum duration in seconds
func (b *AssistantBuilder) WithMaxDuration(seconds int) *AssistantBuilder {
	b.assistant.MaxDurationSeconds = &seconds
	return b
}

// WithTranscriber sets the transcriber configuration
func (b *AssistantBuilder) WithTranscriber(provider, language string) *AssistantBuilder {
	b.assistant.Transcriber = &Transcriber{
		Provider: provider,
		Language: &language,
	}
	return b
}

// WithTemperature sets the model temperature
func (b *AssistantBuilder) WithTemperature(temp float64) *AssistantBuilder {
	if b.assistant.Model == nil {
		b.assistant.Model = &Model{}
	}
	b.assistant.Model.Temperature = &temp
	return b
}

// WithMaxTokens sets the model max tokens
func (b *AssistantBuilder) WithMaxTokens(tokens int) *AssistantBuilder {
	if b.assistant.Model == nil {
		b.assistant.Model = &Model{}
	}
	b.assistant.Model.MaxTokens = &tokens
	return b
}

// WithName sets the assistant name
func (b *AssistantBuilder) WithName(name string) *AssistantBuilder {
	b.assistant.Name = &name
	return b
}

// WithMetadata sets the assistant metadata
func (b *AssistantBuilder) WithMetadata(metadata map[string]interface{}) *AssistantBuilder {
	b.assistant.Metadata = metadata
	return b
}

// Build returns the built Assistant
func (b *AssistantBuilder) Build() *Assistant {
	return b.assistant
}

// RequestBuilder helps build CreateChatRequest configurations
type RequestBuilder struct {
	request *CreateChatRequest
}

// NewRequestBuilder creates a new RequestBuilder
func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		request: &CreateChatRequest{},
	}
}

// WithTextInput sets text input
func (b *RequestBuilder) WithTextInput(text string) *RequestBuilder {
	b.request.Input = text
	return b
}

// WithMessageInput sets message input
func (b *RequestBuilder) WithMessageInput(messages []ChatMessage) *RequestBuilder {
	b.request.Input = messages
	return b
}

// WithAssistantID sets the assistant ID
func (b *RequestBuilder) WithAssistantID(assistantID string) *RequestBuilder {
	b.request.AssistantID = &assistantID
	return b
}

// WithAssistant sets the assistant configuration
func (b *RequestBuilder) WithAssistant(assistant *Assistant) *RequestBuilder {
	b.request.Assistant = assistant
	return b
}

// WithSessionID sets the session ID
func (b *RequestBuilder) WithSessionID(sessionID string) *RequestBuilder {
	b.request.SessionID = &sessionID
	return b
}

// WithPreviousChatID sets the previous chat ID
func (b *RequestBuilder) WithPreviousChatID(chatID string) *RequestBuilder {
	b.request.PreviousChatID = &chatID
	return b
}

// WithName sets the chat name
func (b *RequestBuilder) WithName(name string) *RequestBuilder {
	b.request.Name = &name
	return b
}

// WithStreaming enables streaming
func (b *RequestBuilder) WithStreaming(stream bool) *RequestBuilder {
	b.request.Stream = &stream
	return b
}

// WithAssistantOverrides sets assistant overrides
func (b *RequestBuilder) WithAssistantOverrides(overrides *AssistantOverrides) *RequestBuilder {
	b.request.AssistantOverrides = overrides
	return b
}

// Build returns the built CreateChatRequest
func (b *RequestBuilder) Build() *CreateChatRequest {
	return b.request
}

// Validate validates the built request
func (b *RequestBuilder) Validate() error {
	if b.request.Input == nil {
		return fmt.Errorf("input is required")
	}

	// Validate that at least one of assistantId, assistant, sessionId, or previousChatId is provided
	if b.request.AssistantID == nil && b.request.Assistant == nil && b.request.SessionID == nil && b.request.PreviousChatID == nil {
		return fmt.Errorf("at least one of assistantId, assistant, sessionId, or previousChatId is required")
	}

	// Validate that sessionId and previousChatId are mutually exclusive
	if b.request.SessionID != nil && b.request.PreviousChatID != nil {
		return fmt.Errorf("sessionId and previousChatId are mutually exclusive")
	}

	// Validate name length if provided
	if b.request.Name != nil && len(*b.request.Name) > 40 {
		return fmt.Errorf("name must be 40 characters or less")
	}

	return nil
}

// Helper functions for common assistant configurations

// CreateAnthropicAssistant creates a basic Anthropic Claude assistant
func CreateAnthropicAssistant(systemPrompt string) *Assistant {
	return NewAssistantBuilder().
		WithModel("anthropic", "claude-3-opus-20240229").
		WithSystemMessage(systemPrompt).
		WithTemperature(0.7).
		WithMaxTokens(1000).
		WithFirstMessage("Hello! How can I help you today?").
		WithFirstMessageMode("assistant-speaks-first").
		Build()
}

// CreateOpenAIAssistant creates a basic OpenAI GPT assistant
func CreateOpenAIAssistant(systemPrompt string) *Assistant {
	return NewAssistantBuilder().
		WithModel("openai", "gpt-4").
		WithSystemMessage(systemPrompt).
		WithTemperature(0.7).
		WithMaxTokens(1000).
		WithFirstMessage("Hello! How can I help you today?").
		WithFirstMessageMode("assistant-speaks-first").
		Build()
}

// CreateSalesAssistant creates a specialized sales assistant for SalesGuru
func CreateSalesAssistant(companyName, industry string) *Assistant {
	systemPrompt := fmt.Sprintf(`You are a professional sales assistant for %s, specializing in %s. 
Your role is to:
1. Qualify leads by understanding their needs and budget
2. Provide helpful information about our services
3. Schedule appointments when appropriate
4. Maintain a friendly, professional tone
5. Ask relevant questions to understand customer requirements

Always be helpful, informative, and focused on providing value to potential customers.`, companyName, industry)

	return NewAssistantBuilder().
		WithModel("anthropic", "claude-3-opus-20240229").
		WithSystemMessage(systemPrompt).
		WithTemperature(0.7).
		WithMaxTokens(1500).
		WithFirstMessage(fmt.Sprintf("Hello! I'm here to help you learn more about %s's %s services. How can I assist you today?", companyName, industry)).
		WithFirstMessageMode("assistant-speaks-first").
		WithName(fmt.Sprintf("%s Sales Assistant", companyName)).
		Build()
}

// CreateTelecomAssistant creates a specialized telecom sales assistant for Colombian market
func CreateTelecomAssistant() *Assistant {
	systemPrompt := `Eres un asistente de ventas especializado en servicios de telecomunicaciones en Colombia, específicamente en fibra óptica para internet.

Tu rol es:
1. Calificar leads entendiendo sus necesidades de internet y presupuesto
2. Explicar los beneficios de la fibra óptica vs otros tipos de conexión
3. Preguntar sobre su ubicación, estrato socioeconómico, y tipo de edificio
4. Ofrecer planes apropiados según sus necesidades
5. Programar citas técnicas cuando sea apropiado
6. Mantener un tono amigable y profesional en español

Siempre sé útil, informativo, y enfócate en brindar valor a los clientes potenciales.
Conoces bien el mercado colombiano y las necesidades específicas de conectividad en ciudades como Bogotá, Medellín, Cali, y Barranquilla.`

	return NewAssistantBuilder().
		WithModel("anthropic", "claude-3-opus-20240229").
		WithSystemMessage(systemPrompt).
		WithTemperature(0.7).
		WithMaxTokens(1500).
		WithFirstMessage("¡Hola! Soy tu asistente especializado en servicios de fibra óptica. ¿Te interesa conocer nuestros planes de internet de alta velocidad?").
		WithFirstMessageMode("assistant-speaks-first").
		WithName("Asistente de Fibra Óptica").
		WithTranscriber("assembly-ai", "es").
		WithVoice("azure", "es-CO-SalomeNeural").
		Build()
}

// Helper functions for creating chat messages

// CreateChatMessage creates a new chat message
func CreateChatMessage(role, message string, time, secondsFromStart int64) ChatMessage {
	return ChatMessage{
		Role:             role,
		Content:          message,
		Time:             time,
		SecondsFromStart: secondsFromStart,
	}
}

// CreateUserMessage creates a user chat message
func CreateUserMessage(message string) ChatMessage {
	return ChatMessage{
		Role:    "user",
		Content: message,
	}
}

// CreateAssistantMessage creates an assistant chat message
func CreateAssistantMessage(message string) ChatMessage {
	return ChatMessage{
		Role:    "assistant",
		Content: message,
	}
}

// CreateSystemMessage creates a system chat message
func CreateSystemMessage(message string) ChatMessage {
	return ChatMessage{
		Role:    "system",
		Content: message,
	}
}

// Helper functions for common request patterns

// CreateSimpleTextRequest creates a simple text-based chat request
func CreateSimpleTextRequest(text, assistantID string) *CreateChatRequest {
	return NewRequestBuilder().
		WithTextInput(text).
		WithAssistantID(assistantID).
		Build()
}

// CreateConversationRequest creates a request with message history
func CreateConversationRequest(messages []ChatMessage, assistantID string) *CreateChatRequest {
	return NewRequestBuilder().
		WithMessageInput(messages).
		WithAssistantID(assistantID).
		Build()
}

// CreateStreamingRequest creates a streaming chat request
func CreateStreamingRequest(text, assistantID string) *CreateChatRequest {
	return NewRequestBuilder().
		WithTextInput(text).
		WithAssistantID(assistantID).
		WithStreaming(true).
		Build()
}

// CreateContinuationRequest creates a request to continue a previous chat
func CreateContinuationRequest(text, previousChatID string) *CreateChatRequest {
	return NewRequestBuilder().
		WithTextInput(text).
		WithPreviousChatID(previousChatID).
		Build()
}

// CreateSessionRequest creates a request within a session
func CreateSessionRequest(text, sessionID string) *CreateChatRequest {
	return NewRequestBuilder().
		WithTextInput(text).
		WithSessionID(sessionID).
		Build()
}
