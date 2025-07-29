package config

import (
"fmt"
"os"
"strconv"
"time"

"gopkg.in/yaml.v3"
)

// Config represents the complete VAPI library configuration
type Config struct {
VAPI    VAPIConfig    `yaml:"vapi"`
Tunnel  TunnelConfig  `yaml:"tunnel"`
Events  EventsConfig  `yaml:"events"`
Workers WorkersConfig `yaml:"workers"`
}

// VAPIConfig represents the VAPI API configuration
type VAPIConfig struct {
APIToken string        `yaml:"api_token" env:"VAPI_API_TOKEN"`
BaseURL  string        `yaml:"base_url" env:"VAPI_BASE_URL"`
Timeout  time.Duration `yaml:"timeout" env:"VAPI_TIMEOUT"`
}

// TunnelConfig represents the tunnel configuration
type TunnelConfig struct {
Provider  string `yaml:"provider" env:"TUNNEL_PROVIDER"`
AuthToken string `yaml:"auth_token" env:"NGROK_AUTH_TOKEN"`
Port      int    `yaml:"port" env:"TUNNEL_PORT"`
Subdomain string `yaml:"subdomain" env:"TUNNEL_SUBDOMAIN"`
}

// EventsConfig represents the events system configuration
type EventsConfig struct {
Backend string      `yaml:"backend" env:"EVENTS_BACKEND"`
Redis   RedisConfig `yaml:"redis"`
}

// RedisConfig represents the Redis configuration
type RedisConfig struct {
Host     string `yaml:"host" env:"REDIS_HOST"`
Port     int    `yaml:"port" env:"REDIS_PORT"`
DB       int    `yaml:"db" env:"REDIS_DB"`
Password string `yaml:"password" env:"REDIS_PASSWORD"`
}

// WorkersConfig represents the worker pool configuration
type WorkersConfig struct {
Count         int           `yaml:"count" env:"WORKERS_COUNT"`
QueueSize     int           `yaml:"queue_size" env:"WORKERS_QUEUE_SIZE"`
RetryAttempts int           `yaml:"retry_attempts" env:"WORKERS_RETRY_ATTEMPTS"`
RetryDelay    time.Duration `yaml:"retry_delay" env:"WORKERS_RETRY_DELAY"`
}

// LoadFromFile loads configuration from a YAML file
func LoadFromFile(filename string) (*Config, error) {
data, err := os.ReadFile(filename)
if err != nil {
return nil, fmt.Errorf("failed to read config file: %w", err)
}

// Expand environment variables in the YAML content
expandedData := os.ExpandEnv(string(data))

var config Config
if err := yaml.Unmarshal([]byte(expandedData), &config); err != nil {
return nil, fmt.Errorf("failed to parse config file: %w", err)
}

// Apply defaults
config.applyDefaults()

return &config, nil
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *Config {
config := &Config{
VAPI: VAPIConfig{
APIToken: getEnv("VAPI_API_TOKEN", ""),
BaseURL:  getEnv("VAPI_BASE_URL", "https://api.vapi.ai"),
Timeout:  parseDuration(getEnv("VAPI_TIMEOUT", "30s")),
},
Tunnel: TunnelConfig{
Provider:  getEnv("TUNNEL_PROVIDER", "ngrok"),
AuthToken: getEnv("NGROK_AUTH_TOKEN", ""),
Port:      parseInt(getEnv("TUNNEL_PORT", "8080")),
Subdomain: getEnv("TUNNEL_SUBDOMAIN", ""),
},
Events: EventsConfig{
Backend: getEnv("EVENTS_BACKEND", "redis"),
Redis: RedisConfig{
Host:     getEnv("REDIS_HOST", "localhost"),
Port:     parseInt(getEnv("REDIS_PORT", "6379")),
DB:       parseInt(getEnv("REDIS_DB", "0")),
Password: getEnv("REDIS_PASSWORD", ""),
},
},
Workers: WorkersConfig{
Count:         parseInt(getEnv("WORKERS_COUNT", "3")),
QueueSize:     parseInt(getEnv("WORKERS_QUEUE_SIZE", "100")),
RetryAttempts: parseInt(getEnv("WORKERS_RETRY_ATTEMPTS", "3")),
RetryDelay:    parseDuration(getEnv("WORKERS_RETRY_DELAY", "5s")),
},
}

config.applyDefaults()
return config
}

// applyDefaults applies default values to the configuration
func (c *Config) applyDefaults() {
if c.VAPI.BaseURL == "" {
c.VAPI.BaseURL = "https://api.vapi.ai"
}
if c.VAPI.Timeout == 0 {
c.VAPI.Timeout = 30 * time.Second
}
if c.Tunnel.Provider == "" {
c.Tunnel.Provider = "ngrok"
}
if c.Tunnel.Port == 0 {
c.Tunnel.Port = 8080
}
if c.Events.Backend == "" {
c.Events.Backend = "redis"
}
if c.Events.Redis.Host == "" {
c.Events.Redis.Host = "localhost"
}
if c.Events.Redis.Port == 0 {
c.Events.Redis.Port = 6379
}
if c.Workers.Count == 0 {
c.Workers.Count = 3
}
if c.Workers.QueueSize == 0 {
c.Workers.QueueSize = 100
}
if c.Workers.RetryAttempts == 0 {
c.Workers.RetryAttempts = 3
}
if c.Workers.RetryDelay == 0 {
c.Workers.RetryDelay = 5 * time.Second
}
}

// Helper functions
func getEnv(key, defaultValue string) string {
if value := os.Getenv(key); value != "" {
return value
}
return defaultValue
}

func parseInt(s string) int {
if i, err := strconv.Atoi(s); err == nil {
return i
}
return 0
}

func parseDuration(s string) time.Duration {
if d, err := time.ParseDuration(s); err == nil {
return d
}
return 0
}
