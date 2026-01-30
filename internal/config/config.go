package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	BrowserPath    string
	BrowserTimeout time.Duration
	ViewportWidth  int
	ViewportHeight int
	Headless       bool
}

// Load returns configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		BrowserPath:    getEnv("BROWSER_PATH", ""),
		BrowserTimeout: getDurationEnv("BROWSER_TIMEOUT", 30*time.Second),
		ViewportWidth:  getIntEnv("VIEWPORT_WIDTH", 1280),
		ViewportHeight: getIntEnv("VIEWPORT_HEIGHT", 800),
		Headless:       getBoolEnv("HEADLESS", false),
	}
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}
