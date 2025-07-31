# VAPI Go Library - RESTORED

A reusable Go library for integrating with VAPI (Voice API) services, featuring event-driven architecture, Redis-based messaging, and comprehensive voice functionality.

## 🎉 Voice Functionality Restored!

This library has been fully restored with complete voice functionality that was working in the Heirloomz backend. All features from the original README are now implemented and functional.

## Features

✅ **Event-Driven Architecture**: Redis-based event system with worker pools  
✅ **VAPI Integration**: Complete VAPI API client with webhook support  
✅ **Real-time Processing**: Live webhook processing of end-of-call-report events  
✅ **Chat Functionality**: Preserved chat functionality (used by SalesGuru)  
✅ **Webhook Server**: Built-in HTTP server for VAPI webhooks  
✅ **Configurable**: YAML and environment variable configuration  
✅ **Async Processing**: Multi-worker event processing with retries  
✅ **Extensible**: Plugin architecture for custom event handlers  
✅ **Domain Agnostic**: No business logic coupling - pure VAPI integration  

## Quick Start

### Installation

```bash
go get github.com/heirloomz/vapi-go-library
```

### Basic Usage

```go
package main

import (
    "log"
    "github.com/heirloomz/vapi-go-library"
    "github.com/heirloomz/vapi-go-library/pkg/config"
)

func main() {
    // Load configuration
    cfg := config.LoadFromEnv()
    
    // Create library instance
    library, err := vapi.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    
    // Start the library (webhook server, event bus, etc.)
    if err := library.Start(); err != nil {
        log.Fatal(err)
    }
    defer library.Stop()
    
    // Your application logic here
    select {} // Keep running
}
```

## Configuration

### Environment Variables

```bash
# VAPI Configuration
VAPI_API_TOKEN=your_vapi_token_here
VAPI_BASE_URL=https://api.vapi.ai
VAPI_TIMEOUT=30s

# Tunnel Configuration
TUNNEL_PROVIDER=ngrok
NGROK_AUTH_TOKEN=your_ngrok_token_here
TUNNEL_PORT=8080

# Events Configuration
EVENTS_BACKEND=redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0
REDIS_PASSWORD=

# Workers Configuration
WORKERS_COUNT=3
WORKERS_QUEUE_SIZE=100
WORKERS_RETRY_ATTEMPTS=3
WORKERS_RETRY_DELAY=5s
```

### YAML Configuration

```yaml
# vapi-config.yaml
vapi:
  api_token: "${VAPI_API_TOKEN}"
  base_url: "https://api.vapi.ai"
  timeout: 30s

tunnel:
  provider: "ngrok"
  auth_token: "${NGROK_AUTH_TOKEN}"
  port: 8080

events:
  backend: "redis"
  redis:
    host: "${REDIS_HOST:-localhost}"
    port: ${REDIS_PORT:-6379}
    db: ${REDIS_DB:-0}
    password: "${REDIS_PASSWORD:-}"

workers:
  count: 3
  queue_size: 100
  retry_attempts: 3
  retry_delay: "5s"
```

## Event System

The library uses an event-driven architecture with the following event types:

- `vapi.call.completed` - Call completion events
- `vapi.call.started` - Call initiation events  
- `vapi.transcript.update` - Real-time transcript updates
- `vapi.assistant.updated` - Assistant configuration changes
- `vapi.file.uploaded` - File upload completions
- `vapi.tool.created` - Tool creation events
- `vapi.webhook.received` - Raw webhook events

### Event Handlers

```go
// Custom event handler
type MyCallHandler struct {
    // Your dependencies
}

func (h *MyCallHandler) Handle(event *events.Event) error {
    // Process the event
    log.Printf("Received call completed: %s", event.ID)
    return nil
}

func (h *MyCallHandler) EventType() string {
    return "vapi.call.completed"
}

// Register the handler
library.EventBus().Subscribe("vapi.call.completed", &MyCallHandler{})
```

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│                 │    │                 │    │                 │
│   VAPI Client   │────│  Event System   │────│  Redis Queue    │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│                 │    │                 │    │                 │
│ Webhook Server  │────│ Worker Pool     │────│ Event Handlers │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │
         │
         ▼
┌─────────────────┐
│                 │
│ Call Processor  │
│                 │
└─────────────────┘
```

## Core Components

### VAPI Client

Complete VAPI API integration:

```go
// List assistants
assistants, err := library.Voice().ListAssistants()

// Get call details
call, err := library.Voice().GetCall(callID)

// Upload files
file, err := library.Voice().UploadFile(filePath)

// Create tools
tool, err := library.Voice().CreateQueryTool(fileIDs, name, description)

// Update assistants
assistant, err := library.Voice().UpdateAssistant(assistantID, updateReq)
```

### Webhook Processing

Automatic webhook handling:

```go
// Webhook events are automatically processed
// and converted to typed events
type CallCompletedEvent struct {
    CallID      string
    AssistantID string
    Transcript  []Message
    Duration    int
    Status      string
}
```

### Event Processing

Redis-based event queue with workers:

```go
// Events are automatically queued and processed
// Multiple workers handle events concurrently
// Retry logic for failed events
// Dead letter queue for permanent failures
```

## Integration Examples

### Basic Integration

```go
package main

import (
    "github.com/heirloomz/vapi-go-library"
    "github.com/heirloomz/vapi-go-library/pkg/config"
    "github.com/heirloomz/vapi-go-library/pkg/events"
)

type MyApplication struct {
    library *vapi.Library
}

func NewMyApplication(cfg *config.Config) (*MyApplication, error) {
    library, err := vapi.New(cfg)
    if err != nil {
        return nil, err
    }
    
    app := &MyApplication{
        library: library,
    }
    
    // Register domain-specific event handlers
    library.EventBus().Subscribe("vapi.call.completed", &CallCompletedHandler{})
    
    return app, nil
}

func (app *MyApplication) Start() error {
    return app.library.Start()
}

func (app *MyApplication) Stop() error {
    return app.library.Stop()
}
```

### Advanced Integration (Heirloomz Pattern)

```go
// This is how Heirloomz can integrate the restored VAPI library
package services

import (
    "github.com/heirloomz/vapi-go-library"
    "github.com/heirloomz/backend/internal/db"
    "github.com/heirloomz/backend/internal/events"
)

type HeirloomzVAPIService struct {
    library      *vapi.Library
    db           *db.PostgresDB
    eventBus     events.EventBusInterface
}

func NewHeirloomzVAPIService(cfg *config.Config, database *db.PostgresDB, eventBus events.EventBusInterface) (*HeirloomzVAPIService, error) {
    library, err := vapi.New(cfg)
    if err != nil {
        return nil, err
    }
    
    service := &HeirloomzVAPIService{
        library:  library,
        db:       database,
        eventBus: eventBus,
    }
    
    // Register Heirloomz-specific handlers
    library.EventBus().Subscribe("vapi.call.completed", &HeirloomzCallHandler{
        db:       database,
        eventBus: eventBus,
    })
    
    return service, nil
}

// HeirloomzCallHandler processes calls for story generation
type HeirloomzCallHandler struct {
    db       *db.PostgresDB
    eventBus events.EventBusInterface
}

func (h *HeirloomzCallHandler) Handle(event *events.Event) error {
    // Extract call data
    callData := event.Data.(*voice.ProcessedCall)
    
    // Store in processed_calls table
    processedCall := &models.ProcessedCall{
        CallID:      callData.CallID,
        AssistantID: callData.AssistantID,
        Transcript:  callData.Transcript,
        Duration:    callData.Duration,
        Status:      callData.Status,
    }
    
    if err := h.db.CreateProcessedCall(processedCall); err != nil {
        return err
    }
    
    // Publish to Heirloomz event bus for story generation
    return h.eventBus.Publish("call-processed", processedCall)
}

func (h *HeirloomzCallHandler) EventType() string {
    return "vapi.call.completed"
}
```

## API Reference

### Core Methods

```go
// Library management
func New(config *Config) (*Library, error)
func (l *Library) Start() error
func (l *Library) Stop() error

// Voice operations
func (l *Library) Voice() *voice.VoiceClient
func (v *VoiceClient) ListAssistants() ([]Assistant, error)
func (v *VoiceClient) GetAssistant(id string) (*Assistant, error)
func (v *VoiceClient) UpdateAssistant(id string, req *UpdateRequest) (*Assistant, error)
func (v *VoiceClient) ListCalls(assistantID string, limit int) ([]Call, error)
func (v *VoiceClient) GetCall(id string) (*Call, error)

// File operations
func (v *VoiceClient) UploadFile(path string) (*File, error)
func (v *VoiceClient) CreateQueryTool(fileIDs []string, name, desc string) (*Tool, error)
func (v *VoiceClient) AttachToolToAssistant(assistantID, toolID string) error

// Event system
func (l *Library) EventBus() events.EventBus
func (e *EventBus) Subscribe(eventType string, handler Handler) error
func (e *EventBus) Publish(event *Event) error

// Chat operations (preserved)
func (l *Library) Chat() *chat.Client
func (c *Client) CreateChat(ctx context.Context, req *CreateChatRequest) (*ChatResponse, error)
```

## What Was Restored

This library restoration includes all the sophisticated voice functionality that was working in the Heirloomz backend:

### ✅ Complete VAPI API Client
- Full HTTP client with all VAPI endpoints
- Assistant management and synchronization
- Call listing and retrieval with pagination
- File upload with MIME type detection
- Tool creation and attachment
- Webhook URL management

### ✅ Real-time Event Processing
- Live webhook server for VAPI events
- End-of-call-report processing
- Redis-based event bus with pub/sub
- Automatic call transcript extraction
- Event-driven story generation pipeline

### ✅ Advanced Features
- Multi-directory caching system
- Comprehensive error handling and retries
- Configurable storage directories
- Health checks and monitoring
- Graceful shutdown handling

### ✅ Preserved Functionality
- Complete chat functionality (used by SalesGuru)
- All existing configuration options
- Backward compatibility with existing integrations

## Development

### Building

```bash
go build ./...
```

### Testing

```bash
go test ./...
```

### Dependencies

```bash
go mod tidy
```

## Migration from Internal Services

If you're migrating from the Heirloomz internal VAPI services, here's the mapping:

| Old Internal Service | New Library Method |
|---------------------|-------------------|
| `VAPIService.ListAssistants()` | `library.Voice().ListAssistants()` |
| `VAPIService.GetCall()` | `library.Voice().GetCall()` |
| `VAPIService.UploadFile()` | `library.Voice().UploadFile()` |
| `LiveWebhookHandler` | Built-in webhook server |
| `LiveCallProcessor` | Built-in call processor |
| `RedisEventBus` | `library.EventBus()` |

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**Status**: ✅ **FULLY RESTORED** - All voice functionality from the working Heirloomz backend has been successfully extracted and integrated into this reusable library while preserving existing chat functionality.
