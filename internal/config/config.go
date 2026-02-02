// config loads settings from .env file
package config

import (
	"errors"
	"os"
	"time"
)

// Config holds server and logging settings; filled from env in Load()
type Config struct {
	HTTPAddr        string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	RequestTimeout  time.Duration
	LogLevel        string
}

// Load builds Config from env vars; uses defaults when unset
func Load() *Config {
	return &Config{
		HTTPAddr:        getEnv("HTTP_ADDR", ":8080"),
		ReadTimeout:     getDurationEnv("HTTP_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:    getDurationEnv("HTTP_WRITE_TIMEOUT", 10*time.Second),
		ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 15*time.Second),
		RequestTimeout:  getDurationEnv("REQUEST_TIMEOUT", 30*time.Second),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

var ErrInvalidConfig = errors.New("invalid config: HTTP_ADDR must be non-empty")

// Validate checks required fields (e.g. HTTP_ADDR non-empty)
func (c *Config) Validate() error {
	if c.HTTPAddr == "" {
		return ErrInvalidConfig
	}
	return nil
}
