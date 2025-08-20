package models

import "time"

// Config represents the YAML configuration structure
type Config struct {
	Audiobooks []Entry `yaml:"audiobooks"`
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