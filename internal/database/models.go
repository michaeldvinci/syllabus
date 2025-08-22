package database

import (
	"time"
)

// Series represents a book series in the database
type Series struct {
	ID         int       `db:"id" json:"id"`
	Title      string    `db:"title" json:"title"`
	AudibleID  *string   `db:"audible_id" json:"audible_id,omitempty"`
	AudibleURL *string   `db:"audible_url" json:"audible_url,omitempty"`
	AmazonASIN *string   `db:"amazon_asin" json:"amazon_asin,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// Book represents an individual book in a series
type Book struct {
	ID          int        `db:"id" json:"id"`
	SeriesID    int        `db:"series_id" json:"series_id"`
	Provider    string     `db:"provider" json:"provider"`
	Title       string     `db:"title" json:"title"`
	BookNumber  *int       `db:"book_number" json:"book_number,omitempty"`
	ReleaseDate *time.Time `db:"release_date" json:"release_date,omitempty"`
	IsPreorder  bool       `db:"is_preorder" json:"is_preorder"`
	IsLatest    bool       `db:"is_latest" json:"is_latest"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

// ScrapeJob represents a scraping operation
type ScrapeJob struct {
	ID           int        `db:"id" json:"id"`
	SeriesID     int        `db:"series_id" json:"series_id"`
	Provider     string     `db:"provider" json:"provider"`
	Status       string     `db:"status" json:"status"`
	StartedAt    *time.Time `db:"started_at" json:"started_at,omitempty"`
	CompletedAt  *time.Time `db:"completed_at" json:"completed_at,omitempty"`
	ErrorMessage *string    `db:"error_message" json:"error_message,omitempty"`
	BookCount    int        `db:"book_count" json:"book_count"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
}

// SeriesStats represents aggregated series data from the view
type SeriesStats struct {
	ID                int        `db:"id" json:"id"`
	Title             string     `db:"title" json:"title"`
	AudibleID         *string    `db:"audible_id" json:"audible_id,omitempty"`
	AmazonASIN        *string    `db:"amazon_asin" json:"amazon_asin,omitempty"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updated_at"`
	
	// Audible data
	AudibleCount      int        `db:"audible_count" json:"audible_count"`
	AudibleLatestTitle *string   `db:"audible_latest_title" json:"audible_latest_title,omitempty"`
	AudibleLatestDate *time.Time `db:"audible_latest_date" json:"audible_latest_date,omitempty"`
	AudibleNextTitle  *string    `db:"audible_next_title" json:"audible_next_title,omitempty"`
	AudibleNextDate   *time.Time `db:"audible_next_date" json:"audible_next_date,omitempty"`
	
	// Amazon data
	AmazonCount       int        `db:"amazon_count" json:"amazon_count"`
	AmazonLatestTitle *string    `db:"amazon_latest_title" json:"amazon_latest_title,omitempty"`
	AmazonLatestDate  *time.Time `db:"amazon_latest_date" json:"amazon_latest_date,omitempty"`
	AmazonNextTitle   *string    `db:"amazon_next_title" json:"amazon_next_title,omitempty"`
	AmazonNextDate    *time.Time `db:"amazon_next_date" json:"amazon_next_date,omitempty"`
}

// JobStatus constants
const (
	JobStatusPending   = "pending"
	JobStatusRunning   = "running"
	JobStatusCompleted = "completed"
	JobStatusFailed    = "failed"
)

// Provider constants
const (
	ProviderAudible = "audible"
	ProviderAmazon  = "amazon"
)