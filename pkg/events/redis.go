package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RedisEventBus implements EventBus using Redis pub/sub
type RedisEventBus struct {
	client     *redis.Client
	ctx        context.Context
	cancelFunc context.CancelFunc
	handlers   map[string][]Handler
}

// NewRedisEventBus creates a new Redis-based event bus
func NewRedisEventBus(host string, port int, password string, db int) (*RedisEventBus, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	// Test connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &RedisEventBus{
		client:     client,
		ctx:        ctx,
		cancelFunc: cancel,
		handlers:   make(map[string][]Handler),
	}, nil
}

// Publish publishes an event to the bus
func (r *RedisEventBus) Publish(event *Event) error {
	channel := fmt.Sprintf("events:%s", event.Type)

	// Marshal the event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publish to Redis
	err = r.client.Publish(r.ctx, channel, eventJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to publish event to Redis: %w", err)
	}

	return nil
}

// Subscribe subscribes a handler to events of a specific type
func (r *RedisEventBus) Subscribe(eventType string, handler Handler) error {
	// Add handler to local registry
	r.handlers[eventType] = append(r.handlers[eventType], handler)

	// Subscribe to Redis channel
	channel := fmt.Sprintf("events:%s", eventType)

	go func() {
		pubsub := r.client.Subscribe(r.ctx, channel)
		defer pubsub.Close()

		// Wait for confirmation that subscription is created
		_, err := pubsub.Receive(r.ctx)
		if err != nil {
			return
		}

		// Listen for messages
		ch := pubsub.Channel()
		for {
			select {
			case msg := <-ch:
				if msg == nil {
					continue
				}

				// Parse the event
				var event Event
				if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
					continue
				}

				// Handle the event with all registered handlers
				for _, handler := range r.handlers[eventType] {
					go func(h Handler, e Event) {
						h.Handle(&e)
					}(handler, event)
				}

			case <-r.ctx.Done():
				return
			}
		}
	}()

	return nil
}

// Unsubscribe removes a handler from events of a specific type
func (r *RedisEventBus) Unsubscribe(eventType string, handler Handler) error {
	handlers := r.handlers[eventType]
	for i, h := range handlers {
		if h == handler {
			r.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
	return nil
}

// Start starts the event bus
func (r *RedisEventBus) Start() error {
	// Redis event bus is started when created
	return nil
}

// Stop stops the event bus
func (r *RedisEventBus) Stop() error {
	if r.cancelFunc != nil {
		r.cancelFunc()
	}

	if r.client != nil {
		return r.client.Close()
	}

	return nil
}

// Health checks if the Redis connection is healthy
func (r *RedisEventBus) Health() error {
	_, err := r.client.Ping(r.ctx).Result()
	return err
}

// CallProcessedEventData represents data for call-processed events
type CallProcessedEventData struct {
	ProcessedCallID string `json:"processed_call_id"`
	CallID          string `json:"call_id"`
	AssistantID     string `json:"assistant_id"`
}

// PublishCallProcessed publishes a call-processed event
func (r *RedisEventBus) PublishCallProcessed(processedCallID, callID, assistantID string) error {
	eventData := CallProcessedEventData{
		ProcessedCallID: processedCallID,
		CallID:          callID,
		AssistantID:     assistantID,
	}

	event := NewEvent("call-processed", "vapi-library", eventData)
	return r.Publish(event)
}
