# VAPI Go Library - Chat Package

This package provides comprehensive support for the VAPI Chat API, enabling you to create conversational AI assistants that can respond to messages based on assistant prompts.

## Features

- ✅ **Complete Chat API Support**: Full implementation of VAPI's Chat API
- ✅ **Streaming Support**: Real-time streaming chat responses
- ✅ **Assistant Builder**: Easy-to-use builder pattern for creating assistants
- ✅ **Request Builder**: Fluent API for building chat requests
- ✅ **Pre-built Assistants**: Ready-to-use assistants for common use cases
- ✅ **Validation**: Built-in request validation
- ✅ **Context Management**: Support for sessions and chat continuations
- ✅ **Colombian Market Support**: Specialized telecom assistant for Colombian market

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "log"
    
    "github.com/heirloomz/vapi-go-library"
    "github.com/heirloomz/vapi-go-library/pkg/config"
)

func main() {
    // Load configuration
    cfg := config.LoadFromEnv()
    
    // Create VAPI library
    vapiLib, err := vapi.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    
    // Get chat client
    chatClient := vapiLib.Chat()
    
    // Create a simple chat
    assistantID := "your-assistant-id"
    response, err := chatClient.CreateChatWithText(
        context.Background(),
        "Hello! Can you help me with internet plans?",
        &assistantID,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Print response
    if len(response.Output) > 0 {
        log.Printf("Assistant: %s", response.Output[0].Message)
    }
}
```

### Using Custom Assistants

```go
// Create a custom sales assistant
assistant := chat.CreateSalesAssistant("TechCorp", "telecommunications")

response, err := chatClient.CreateChatWithAssistant(
    ctx,
    "I'm interested in upgrading my internet connection",
    assistant,
)
```

### Streaming Chat

```go
// Create streaming chat
responseChan, errorChan := chatClient.CreateStreamingChatWithText(
    ctx,
    "What are your fiber optic plans?",
    &assistantID,
)

// Process streaming responses
for {
    select {
    case response, ok := <-responseChan:
        if !ok {
            return
        }
        fmt.Print(response.Message)
        if response.Done {
            return
        }
    case err := <-errorChan:
        if err != nil {
            log.Printf("Error: %v", err)
            return
        }
    }
}
```

## Assistant Builder

The `AssistantBuilder` provides a fluent API for creating custom assistants:

```go
assistant := chat.NewAssistantBuilder().
    WithModel("anthropic", "claude-3-opus-20240229").
    WithSystemMessage("You are a helpful sales assistant...").
    WithTemperature(0.7).
    WithMaxTokens(1500).
    WithFirstMessage("Hello! How can I help you today?").
    WithVoice("azure", "en-US-JennyNeural").
    WithTranscriber("assembly-ai", "en").
    WithName("Sales Assistant").
    Build()
```

## Request Builder

The `RequestBuilder` helps create complex chat requests:

```go
request := chat.NewRequestBuilder().
    WithTextInput("Hello, I need help with internet plans").
    WithAssistant(assistant).
    WithName("Customer Inquiry").
    WithStreaming(true).
    Build()

response, err := chatClient.CreateChat(ctx, request)
```

## Pre-built Assistants

### Sales Assistant

```go
// Generic sales assistant
assistant := chat.CreateSalesAssistant("YourCompany", "telecommunications")
```

### Telecom Assistant (Colombian Market)

```go
// Specialized for Colombian telecom market
assistant := chat.CreateTelecomAssistant()
```

### Anthropic Assistant

```go
// Basic Anthropic Claude assistant
assistant := chat.CreateAnthropicAssistant("You are a helpful assistant...")
```

### OpenAI Assistant

```go
// Basic OpenAI GPT assistant
assistant := chat.CreateOpenAIAssistant("You are a helpful assistant...")
```

## Chat Continuation

### Continue Previous Chat

```go
// Continue from previous chat
response, err := chatClient.ContinueChat(
    ctx,
    "What's the cheapest plan you have?",
    previousChatID,
)
```

### Session-based Chat

```go
// Create chat within a session
response, err := chatClient.CreateSessionChat(
    ctx,
    "Hello, I need help",
    sessionID,
)
```

## Message Types

### Creating Messages

```go
// User message
userMsg := chat.CreateUserMessage("Hello, I need help")

// Assistant message
assistantMsg := chat.CreateAssistantMessage("How can I help you?")

// System message
systemMsg := chat.CreateSystemMessage("You are a helpful assistant")

// Custom message with timestamps
customMsg := chat.CreateChatMessage("user", "Hello", time.Now().Unix(), 0)
```

### Conversation History

```go
messages := []chat.ChatMessage{
    chat.CreateUserMessage("Hello"),
    chat.CreateAssistantMessage("Hi! How can I help?"),
    chat.CreateUserMessage("I need internet plans"),
}

response, err := chatClient.CreateChatWithMessages(ctx, messages, &assistantID)
```

## Configuration

### Environment Variables

```bash
# Required
VAPI_API_TOKEN=your_vapi_api_token

# Optional
VAPI_BASE_URL=https://api.vapi.ai  # Default
VAPI_TIMEOUT=30s                  # Default
```

### Programmatic Configuration

```go
cfg := &config.Config{
    VAPI: config.VAPIConfig{
        APIToken: "your-token",
        BaseURL:  "https://api.vapi.ai",
        Timeout:  30 * time.Second,
    },
}
```

## Error Handling

```go
response, err := chatClient.CreateChatWithText(ctx, "Hello", &assistantID)
if err != nil {
    // Handle different error types
    switch {
    case strings.Contains(err.Error(), "API error"):
        log.Printf("API Error: %v", err)
    case strings.Contains(err.Error(), "validation"):
        log.Printf("Validation Error: %v", err)
    default:
        log.Printf("Unknown Error: %v", err)
    }
    return
}
```

## Validation

```go
// Validate request before sending
request := chat.NewRequestBuilder().
    WithTextInput("Hello").
    WithAssistantID("assistant-id").
    Build()

if err := chatClient.ValidateRequest(request); err != nil {
    log.Printf("Validation failed: %v", err)
    return
}
```

## Advanced Examples

### Complex Assistant Configuration

```go
assistant := chat.NewAssistantBuilder().
    WithModel("anthropic", "claude-3-opus-20240229").
    WithSystemMessage(`You are an expert telecommunications consultant...`).
    WithTemperature(0.7).
    WithMaxTokens(2000).
    WithFirstMessage("¡Hola! ¿En qué puedo ayudarle?").
    WithFirstMessageMode("assistant-speaks-first").
    WithTranscriber("assembly-ai", "es").
    WithVoice("azure", "es-CO-SalomeNeural").
    WithMaxDuration(1800). // 30 minutes
    WithMetadata(map[string]interface{}{
        "department": "sales",
        "language":   "spanish",
        "market":     "colombia",
    }).
    Build()
```

### Streaming with Context Cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
defer cancel()

responseChan, errorChan := chatClient.CreateStreamingChatWithAssistant(
    ctx,
    "¿Qué planes tienen disponibles?",
    assistant,
)

for {
    select {
    case response, ok := <-responseChan:
        if !ok {
            return
        }
        fmt.Print(response.Message)
        if response.Done {
            return
        }
    case err := <-errorChan:
        if err != nil {
            log.Printf("Error: %v", err)
            return
        }
    case <-ctx.Done():
        log.Println("Request cancelled or timed out")
        return
    }
}
```

## API Reference

### Client Methods

- `CreateChat(ctx, request)` - Create a new chat
- `CreateStreamingChat(ctx, request)` - Create a streaming chat
- `CreateChatWithText(ctx, text, assistantID)` - Simple text chat
- `CreateChatWithMessages(ctx, messages, assistantID)` - Chat with history
- `CreateChatWithAssistant(ctx, text, assistant)` - Chat with custom assistant
- `ContinueChat(ctx, text, previousChatID)` - Continue previous chat
- `CreateSessionChat(ctx, text, sessionID)` - Session-based chat
- `ValidateRequest(request)` - Validate request
- `SetTimeout(duration)` - Set custom timeout

### Builder Methods

#### AssistantBuilder
- `WithModel(provider, model)` - Set AI model
- `WithSystemMessage(content)` - Add system message
- `WithTemperature(temp)` - Set temperature
- `WithMaxTokens(tokens)` - Set max tokens
- `WithVoice(provider, voiceID)` - Set voice
- `WithTranscriber(provider, language)` - Set transcriber
- `WithFirstMessage(message)` - Set first message
- `WithName(name)` - Set assistant name
- `WithMetadata(metadata)` - Set metadata

#### RequestBuilder
- `WithTextInput(text)` - Set text input
- `WithMessageInput(messages)` - Set message input
- `WithAssistantID(id)` - Set assistant ID
- `WithAssistant(assistant)` - Set custom assistant
- `WithSessionID(id)` - Set session ID
- `WithPreviousChatID(id)` - Set previous chat ID
- `WithStreaming(enabled)` - Enable streaming
- `WithName(name)` - Set chat name

## Best Practices

1. **Always use context with timeouts** for API calls
2. **Validate requests** before sending to catch errors early
3. **Handle streaming responses** properly with select statements
4. **Use pre-built assistants** for common use cases
5. **Set appropriate timeouts** based on expected response time
6. **Handle errors gracefully** with proper error checking
7. **Use builders** for complex configurations
8. **Reuse clients** instead of creating new ones for each request

## Troubleshooting

### Common Issues

1. **Authentication Error**: Check your VAPI_API_TOKEN
2. **Timeout Error**: Increase timeout or check network connectivity
3. **Validation Error**: Ensure required fields are provided
4. **Rate Limiting**: Implement retry logic with exponential backoff

### Debug Mode

```go
// Enable debug logging (if implemented)
client.SetDebug(true)
```

## Contributing

When contributing to the chat package:

1. Follow Go conventions and best practices
2. Add comprehensive tests for new features
3. Update documentation for API changes
4. Ensure backward compatibility when possible
5. Add examples for new functionality

## License

This package is part of the VAPI Go Library and follows the same license terms.
