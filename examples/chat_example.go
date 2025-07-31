package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/heirloomz/vapi-go-library"
	"github.com/heirloomz/vapi-go-library/pkg/chat"
	"github.com/heirloomz/vapi-go-library/pkg/config"
)

func runChatExamples() {
	// Load configuration from environment
	cfg := config.LoadFromEnv()

	// Create VAPI library instance
	vapiLib, err := vapi.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create VAPI library: %v", err)
	}

	// Get chat client
	chatClient := vapiLib.Chat()

	// Example 1: Simple text chat with assistant ID
	fmt.Println("=== Example 1: Simple Text Chat ===")
	simpleTextExample(chatClient)

	// Example 2: Chat with custom assistant
	fmt.Println("\n=== Example 2: Custom Assistant Chat ===")
	customAssistantExample(chatClient)

	// Example 3: Streaming chat
	fmt.Println("\n=== Example 3: Streaming Chat ===")
	streamingChatExample(chatClient)

	// Example 4: Conversation with message history
	fmt.Println("\n=== Example 4: Conversation with History ===")
	conversationExample(chatClient)

	// Example 5: Sales assistant for telecom
	fmt.Println("\n=== Example 5: Telecom Sales Assistant ===")
	telecomSalesExample(chatClient)

	// Example 6: Continue previous chat
	fmt.Println("\n=== Example 6: Continue Previous Chat ===")
	continueChatExample(chatClient)
}

// Example 1: Simple text chat with assistant ID
func simpleTextExample(client *chat.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Assuming you have an assistant ID from VAPI dashboard
	assistantID := "your-assistant-id-here"

	response, err := client.CreateChatWithText(ctx, "Hello! Can you help me with fiber optic internet plans?", &assistantID)
	if err != nil {
		log.Printf("Error in simple text example: %v", err)
		return
	}

	fmt.Printf("Chat ID: %s\n", response.ID)
	fmt.Printf("Input: %v\n", response.Input)
	if len(response.Output) > 0 {
		fmt.Printf("Assistant Response: %s\n", response.Output[0].Message)
	}
	fmt.Printf("Cost: $%.4f\n", response.Cost)
}

// Example 2: Chat with custom assistant
func customAssistantExample(client *chat.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a custom sales assistant
	assistant := chat.CreateSalesAssistant("TechCorp", "telecommunications")

	response, err := client.CreateChatWithAssistant(ctx, "I'm interested in upgrading my internet connection", assistant)
	if err != nil {
		log.Printf("Error in custom assistant example: %v", err)
		return
	}

	fmt.Printf("Chat ID: %s\n", response.ID)
	fmt.Printf("Assistant Name: %s\n", *response.Assistant.Name)
	if len(response.Output) > 0 {
		fmt.Printf("Assistant Response: %s\n", response.Output[0].Message)
	}
}

// Example 3: Streaming chat
func streamingChatExample(client *chat.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create telecom assistant for streaming
	assistant := chat.CreateTelecomAssistant()

	responseChan, errorChan := client.CreateStreamingChatWithAssistant(ctx, "¿Qué planes de fibra óptica tienen disponibles?", assistant)

	fmt.Println("Streaming response:")
	for {
		select {
		case response, ok := <-responseChan:
			if !ok {
				fmt.Println("\nStreaming completed")
				return
			}
			fmt.Print(response.Message)
			if response.Done {
				fmt.Println("\n[DONE]")
				return
			}
		case err, ok := <-errorChan:
			if ok && err != nil {
				log.Printf("Streaming error: %v", err)
				return
			}
		case <-ctx.Done():
			fmt.Println("\nStreaming timeout")
			return
		}
	}
}

// Example 4: Conversation with message history
func conversationExample(client *chat.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create conversation history
	messages := []chat.ChatMessage{
		chat.CreateUserMessage("Hello, I'm looking for internet plans"),
		chat.CreateAssistantMessage("Hello! I'd be happy to help you find the perfect internet plan. What's your current internet speed and what do you primarily use it for?"),
		chat.CreateUserMessage("I currently have 50 Mbps but I work from home and need something faster"),
	}

	assistantID := "your-assistant-id-here"

	response, err := client.CreateChatWithMessages(ctx, messages, &assistantID)
	if err != nil {
		log.Printf("Error in conversation example: %v", err)
		return
	}

	fmt.Printf("Chat ID: %s\n", response.ID)
	fmt.Printf("Message History Length: %d\n", len(response.Messages))
	if len(response.Output) > 0 {
		fmt.Printf("Assistant Response: %s\n", response.Output[0].Message)
	}
}

// Example 5: Specialized telecom sales assistant
func telecomSalesExample(client *chat.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create specialized telecom assistant for Colombian market
	assistant := chat.CreateTelecomAssistant()

	// Create request with custom configuration
	request := chat.NewRequestBuilder().
		WithTextInput("Hola, vivo en Bogotá y necesito internet para mi apartamento. ¿Qué opciones tienen?").
		WithAssistant(assistant).
		WithName("Consulta Fibra Óptica Bogotá").
		Build()

	response, err := client.CreateChat(ctx, request)
	if err != nil {
		log.Printf("Error in telecom sales example: %v", err)
		return
	}

	fmt.Printf("Chat ID: %s\n", response.ID)
	fmt.Printf("Chat Name: %s\n", *response.Name)
	fmt.Printf("Assistant: %s\n", *response.Assistant.Name)
	if len(response.Output) > 0 {
		fmt.Printf("Respuesta del Asistente: %s\n", response.Output[0].Message)
	}
}

// Example 6: Continue previous chat
func continueChatExample(client *chat.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// First, create an initial chat
	assistant := chat.CreateTelecomAssistant()
	initialResponse, err := client.CreateChatWithAssistant(ctx, "¿Cuáles son sus planes de fibra óptica?", assistant)
	if err != nil {
		log.Printf("Error creating initial chat: %v", err)
		return
	}

	fmt.Printf("Initial Chat ID: %s\n", initialResponse.ID)
	if len(initialResponse.Output) > 0 {
		fmt.Printf("Initial Response: %s\n", initialResponse.Output[0].Message)
	}

	// Continue the conversation
	continuationResponse, err := client.ContinueChat(ctx, "¿Cuál es el plan más económico que tienen?", initialResponse.ID)
	if err != nil {
		log.Printf("Error continuing chat: %v", err)
		return
	}

	fmt.Printf("Continuation Chat ID: %s\n", continuationResponse.ID)
	fmt.Printf("Previous Chat ID: %s\n", *continuationResponse.PreviousChatID)
	if len(continuationResponse.Output) > 0 {
		fmt.Printf("Continuation Response: %s\n", continuationResponse.Output[0].Message)
	}
}

// Advanced example: Building a complex assistant with tools and custom configuration
func advancedAssistantExample() *chat.Assistant {
	// Create an advanced assistant with custom configuration
	assistant := chat.NewAssistantBuilder().
		WithModel("anthropic", "claude-3-opus-20240229").
		WithSystemMessage(`You are an expert telecommunications sales consultant specializing in fiber optic internet services in Colombia.

Your expertise includes:
- Technical knowledge of fiber optic vs cable vs DSL technologies
- Understanding of Colombian telecommunications market and regulations
- Knowledge of pricing strategies for different socioeconomic strata
- Ability to qualify leads based on location, building type, and usage needs
- Experience with technical installation requirements

Your approach:
1. Greet customers warmly in Colombian Spanish
2. Ask qualifying questions about location, current service, and needs
3. Explain technical benefits in simple terms
4. Provide appropriate plan recommendations
5. Handle objections professionally
6. Guide toward scheduling technical visits when appropriate

Always maintain a helpful, professional, and knowledgeable tone.`).
		WithTemperature(0.7).
		WithMaxTokens(2000).
		WithFirstMessage("¡Hola! Soy su consultor especializado en fibra óptica. ¿En qué ciudad se encuentra y qué tipo de servicio de internet está buscando?").
		WithFirstMessageMode("assistant-speaks-first").
		WithName("Consultor Fibra Óptica").
		WithTranscriber("assembly-ai", "es").
		WithVoice("azure", "es-CO-SalomeNeural").
		WithMaxDuration(1800). // 30 minutes
		WithMetadata(map[string]interface{}{
			"department": "sales",
			"product":    "fiber_optic",
			"market":     "colombia",
			"language":   "spanish",
		}).
		Build()

	return assistant
}

// Example of request validation
func validationExample(client *chat.Client) {
	// Create an invalid request (missing required fields)
	invalidRequest := &chat.CreateChatRequest{
		Input: "Hello",
		// Missing assistantId, assistant, sessionId, and previousChatId
	}

	// Validate the request
	if err := client.ValidateRequest(invalidRequest); err != nil {
		fmt.Printf("Validation error (expected): %v\n", err)
	}

	// Create a valid request
	validRequest := chat.NewRequestBuilder().
		WithTextInput("Hello, I need help with internet plans").
		WithAssistantID("your-assistant-id").
		Build()

	// Validate the valid request
	if err := client.ValidateRequest(validRequest); err != nil {
		fmt.Printf("Unexpected validation error: %v\n", err)
	} else {
		fmt.Println("Request validation passed")
	}
}

// Example of different input types
func inputTypesExample(client *chat.Client) {
	ctx := context.Background()
	assistantID := "your-assistant-id"

	// String input
	stringRequest := chat.CreateSimpleTextRequest("What are your fiber optic plans?", assistantID)
	fmt.Printf("String input request: %+v\n", stringRequest)

	// Message array input
	messages := []chat.ChatMessage{
		chat.CreateUserMessage("Hello"),
		chat.CreateAssistantMessage("Hi! How can I help you?"),
		chat.CreateUserMessage("I need internet plans"),
	}
	messageRequest := chat.CreateConversationRequest(messages, assistantID)
	fmt.Printf("Message array input request: %+v\n", messageRequest)

	// You would then use these requests with client.CreateChat(ctx, request)
	_ = ctx // Avoid unused variable warning
}
