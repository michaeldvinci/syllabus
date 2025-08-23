package models

import "time"

// Config represents the YAML configuration structure
type Config struct {
	Audiobooks []Entry    `yaml:"audiobooks"`
	Settings   *Settings  `yaml:"settings,omitempty"`
}

// Settings represents application-wide settings
type Settings struct {
	AutoRefreshInterval int    `yaml:"auto_refresh_interval,omitempty"` // Hours between auto-refreshes (default: 6)
	DefaultWorkers      int    `yaml:"default_workers,omitempty"`       // Number of scraper workers (default: 4)
	ServerPort          int    `yaml:"server_port,omitempty"`           // Server port (default: 8080)
	CacheTimeout        int    `yaml:"cache_timeout,omitempty"`         // Cache timeout in hours (default: 6)
	LogLevel           string  `yaml:"log_level,omitempty"`             // Log level: debug, info, warn, error (default: info)
}

// GetSettings returns the settings with defaults applied and environment variable overrides
func (c *Config) GetSettings() Settings {
	var settings Settings
	
	if c.Settings != nil {
		settings = *c.Settings
	}
	
	// Apply defaults for zero values
	if settings.AutoRefreshInterval == 0 {
		settings.AutoRefreshInterval = 6
	}
	if settings.DefaultWorkers == 0 {
		settings.DefaultWorkers = 4
	}
	if settings.ServerPort == 0 {
		settings.ServerPort = 8080
	}
	if settings.CacheTimeout == 0 {
		settings.CacheTimeout = 6
	}
	if settings.LogLevel == "" {
		settings.LogLevel = "info"
	}
	
	return settings
}

// Entry represents a single audiobook/ebook series entry
type Entry struct {
	Title    string `yaml:"title"`
	Audible  string `yaml:"audible"`
	Amazon   string `yaml:"amazon"`
	AudNum   any    `yaml:"aud_num"`
	AudNext  string `yaml:"aud_next"`
	AudLast  string `yaml:"aud_last"`
	AmznNum  any    `yaml:"amzn_num"`
	AmznNext string `yaml:"amzn_next"`
	AmznLast string `yaml:"amzn_last"`
}

// SeriesIDs holds extracted identifiers for a series
type SeriesIDs struct {
	Title      string
	AudibleID  string
	AudibleURL string
	AmazonASIN string
	Original   Entry
}

// SeriesInfo contains aggregated information about a series
type SeriesInfo struct {
	Title              string
	AudibleCount       int
	AudibleLatestTitle string
	AudibleLatestDate  *time.Time
	AudibleNextTitle   string
	AudibleNextDate    *time.Time

	AmazonCount       int
	AmazonLatestTitle string
	AmazonLatestDate  *time.Time
	AmazonNextTitle   string
	AmazonNextDate    *time.Time

	AudibleID  string
	AmazonASIN string
	Err        error
}

// Provider defines the interface for data providers
type Provider interface {
	Fetch(entry SeriesIDs) (SeriesInfo, error)
}