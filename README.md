# VAPI Go Library

A reusable Go library for integrating with VAPI (Voice API) services, featuring event-driven architecture, Redis-based messaging, and ngrok tunnel management.

## Features

- **Event-Driven Architecture**: Redis-based event system with worker pools
- **VAPI Integration**: Complete VAPI API client with webhook support
- **Tunnel Management**: Automated ngrok tunnel setup and lifecycle management
- **Configurable**: YAML and environment variable configuration
- **Async Processing**: Multi-worker event processing with retries
- **Extensible**: Plugin architecture for custom event handlers
- **Domain Agnostic**: No business logic coupling - pure VAPI integration

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
    
    // Start the library
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
│ Ngrok Tunnel    │
│                 │
└─────────────────┘
```

## Core Components

### VAPI Client

Complete VAPI API integration:

```go
// List assistants
assistants, err := client.ListAssistants()

// Get call details
call, err := client.GetCall(callID)

// Upload files
file, err := client.UploadFile(filePath)

// Create tools
tool, err := client.CreateQueryTool(fileIDs, name, description)

// Update assistants
assistant, err := client.UpdateAssistant(assistantID, updateReq)
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

### Tunnel Management

Automated ngrok tunnel setup:

```go
// Tunnel is automatically created and managed
// Webhook URLs are automatically configured
// SSL certificates handled automatically
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
    "database/sql"
    "github.com/heirloomz/vapi-go-library"
    "github.com/heirloomz/vapi-go-library/pkg/config"
    "github.com/heirloomz/vapi-go-library/pkg/events"
)

type MyApplication struct {
    library *vapi.Library
    db      *sql.DB
}

func NewMyApplication(cfg *config.Config, db *sql.DB) (*MyApplication, error) {
    library, err := vapi.New(cfg)
    if err != nil {
        return nil, err
    }
    
    app := &MyApplication{
        library: library,
        db:      db,
    }
    
    // Register domain-specific event handlers
    library.EventBus().Subscribe("vapi.call.completed", &CallCompletedHandler{
        db: db,
    })
    
    return app, nil
}

func (app *MyApplication) Start() error {
    return app.library.Start()
}

func (app *MyApplication) Stop() error {
    return app.library.Stop()
}
```

### Advanced Integration (Heirloomz Example)

```go
// This is how Heirloomz integrates the VAPI library
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
    callData := event.Data.(*vapi.CallCompletedEvent)
    
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

// VAPI operations
func (l *Library) ListAssistants() ([]Assistant, error)
func (l *Library) GetAssistant(id string) (*Assistant, error)
func (l *Library) UpdateAssistant(id string, req *UpdateRequest) (*Assistant, error)
func (l *Library) ListCalls(assistantID string, limit int) ([]Call, error)
func (l *Library) GetCall(id string) (*Call, error)

// File operations
func (l *Library) UploadFile(path string) (*File, error)
func (l *Library) CreateQueryTool(fileIDs []string, name, desc string) (*Tool, error)
func (l *Library) AttachToolToAssistant(assistantID, toolID string) error

// Event system
func (l *Library) EventBus() EventBus
func (e *EventBus) Subscribe(eventType string, handler EventHandler) error
func (e *EventBus) Publish(eventType string, data interface{}) error
```

### Configuration Options

```go
type Config struct {
    VAPI     VAPIConfig     `yaml:"vapi"`
    Tunnel   TunnelConfig   `yaml:"tunnel"`
    Events   EventsConfig   `yaml:"events"`
    Workers  WorkersConfig  `yaml:"workers"`
}

type VAPIConfig struct {
    APIToken string        `yaml:"api_token"`
    BaseURL  string        `yaml:"base_url"`
    Timeout  time.Duration `yaml:"timeout"`
}

type TunnelConfig struct {
    Provider  string `yaml:"provider"`
    AuthToken string `yaml:"auth_token"`
    Port      int    `yaml:"port"`
}
```

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

## Future Enhancements

### Monitoring & Observability
- [ ] Prometheus metrics integration
- [ ] OpenTelemetry tracing
- [ ] Health check endpoints
- [ ] Performance monitoring

### Enhanced Tunnel Features
- [ ] Custom domain support
- [ ] SSL certificate management
- [ ] Load balancing across multiple tunnels
- [ ] Tunnel failover and redundancy

### Advanced Event Features
- [ ] Event replay functionality
- [ ] Dead letter queue handling
- [ ] Event filtering and routing
- [ ] Event schema validation

### Performance Optimizations
- [ ] Connection pooling for Redis
- [ ] HTTP client connection reuse
- [ ] Response caching
- [ ] Batch processing for events

### Security Enhancements
- [ ] Webhook signature validation
- [ ] Rate limiting
- [ ] API key rotation
- [ ] Audit logging

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details.
