package events

import (
	"fmt"
)

// NewEventBus creates a new event bus based on the backend type
func NewEventBus(backend string, config interface{}) (EventBus, error) {
	switch backend {
	case "redis":
		if redisConfig, ok := config.(RedisConfig); ok {
			return NewRedisEventBus(
				redisConfig.Host,
				redisConfig.Port,
				redisConfig.Password,
				redisConfig.DB,
			)
		}
		return nil, fmt.Errorf("invalid Redis configuration")
	default:
		return nil, fmt.Errorf("unsupported event bus backend: %s", backend)
	}
}

// RedisConfig represents Redis configuration for event bus
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}
