package utils

import (
	"os"
	"strconv"
	"strings"

	"github.com/michaeldvinci/syllabus/internal/models"
)

// ApplyEnvOverrides applies environment variable overrides to settings
// Environment variables take precedence over YAML config values
func ApplyEnvOverrides(settings *models.Settings) {
	// Auto refresh interval
	if env := os.Getenv("SYLLABUS_AUTO_REFRESH_INTERVAL"); env != "" {
		if val, err := strconv.Atoi(env); err == nil && val > 0 {
			settings.AutoRefreshInterval = val
		}
	}
	
	// Default workers
	if env := os.Getenv("SYLLABUS_DEFAULT_WORKERS"); env != "" {
		if val, err := strconv.Atoi(env); err == nil && val > 0 {
			settings.DefaultWorkers = val
		}
	}
	
	// Server port
	if env := os.Getenv("SYLLABUS_SERVER_PORT"); env != "" {
		if val, err := strconv.Atoi(env); err == nil && val > 0 && val <= 65535 {
			settings.ServerPort = val
		}
	}
	// Also support standard PORT env var (common in container deployments)
	if env := os.Getenv("PORT"); env != "" {
		if val, err := strconv.Atoi(env); err == nil && val > 0 && val <= 65535 {
			settings.ServerPort = val
		}
	}
	
	// Cache timeout
	if env := os.Getenv("SYLLABUS_CACHE_TIMEOUT"); env != "" {
		if val, err := strconv.Atoi(env); err == nil && val > 0 {
			settings.CacheTimeout = val
		}
	}
	
	// Log level
	if env := os.Getenv("SYLLABUS_LOG_LEVEL"); env != "" {
		normalized := strings.ToLower(strings.TrimSpace(env))
		if normalized == "debug" || normalized == "info" || normalized == "warn" || normalized == "error" {
			settings.LogLevel = normalized
		}
	}
	
	// Main view
	if env := os.Getenv("SYLLABUS_MAIN_VIEW"); env != "" {
		normalized := strings.ToLower(strings.TrimSpace(env))
		if normalized == "unified" || normalized == "tabbed" {
			settings.MainView = normalized
		}
	}
}

// GetEnvWithDefault returns environment variable value or default if not set
func GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvIntWithDefault returns environment variable as int or default if not set/invalid
func GetEnvIntWithDefault(key string, defaultValue int) int {
	if env := os.Getenv(key); env != "" {
		if val, err := strconv.Atoi(env); err == nil {
			return val
		}
	}
	return defaultValue
}

// GetEnvBoolWithDefault returns environment variable as bool or default if not set/invalid
func GetEnvBoolWithDefault(key string, defaultValue bool) bool {
	if env := os.Getenv(key); env != "" {
		if val, err := strconv.ParseBool(env); err == nil {
			return val
		}
	}
	return defaultValue
}