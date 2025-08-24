package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/michaeldvinci/syllabus/internal/models"
)

// Service provides database operations for the application
type Service struct {
	db *DB
}

// NewService creates a new database service
func NewService(db *DB) *Service {
	return &Service{db: db}
}

// GetAllSeriesStats returns all series with their aggregated stats
func (s *Service) GetAllSeriesStats() ([]SeriesStats, error) {
	query := `SELECT 
	    id, title, audible_id, amazon_asin, updated_at,
	    audible_count, audible_latest_title, audible_latest_date,
	    audible_next_title, audible_next_date,
	    amazon_count, amazon_latest_title, amazon_latest_date,
	    amazon_next_title, amazon_next_date
	    FROM series_stats ORDER BY title`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query series stats: %w", err)
	}
	defer rows.Close()

	var stats []SeriesStats
	for rows.Next() {
		var stat SeriesStats
		var audibleLatestDate, audibleNextDate, amazonLatestDate, amazonNextDate sql.NullString
		var updatedAt string

		err := rows.Scan(
			&stat.ID, &stat.Title, &stat.AudibleID, &stat.AmazonASIN, &updatedAt,
			&stat.AudibleCount, &stat.AudibleLatestTitle, &audibleLatestDate,
			&stat.AudibleNextTitle, &audibleNextDate,
			&stat.AmazonCount, &stat.AmazonLatestTitle, &amazonLatestDate,
			&stat.AmazonNextTitle, &amazonNextDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan series stats: %w", err)
		}

		// Parse string dates to time.Time
		if parsedTime, err := time.Parse("2006-01-02 15:04:05", updatedAt); err == nil {
			stat.UpdatedAt = parsedTime
		}
		if audibleLatestDate.Valid {
			// Try parsing as timestamp first, then as date only
			if parsedTime, err := time.Parse("2006-01-02 15:04:05+00:00", audibleLatestDate.String); err == nil {
				stat.AudibleLatestDate = &parsedTime
			} else if parsedTime, err := time.Parse("2006-01-02", audibleLatestDate.String); err == nil {
				stat.AudibleLatestDate = &parsedTime
			}
		}
		if audibleNextDate.Valid {
			if parsedTime, err := time.Parse("2006-01-02 15:04:05+00:00", audibleNextDate.String); err == nil {
				stat.AudibleNextDate = &parsedTime
			} else if parsedTime, err := time.Parse("2006-01-02", audibleNextDate.String); err == nil {
				stat.AudibleNextDate = &parsedTime
			}
		}
		if amazonLatestDate.Valid {
			if parsedTime, err := time.Parse("2006-01-02 15:04:05+00:00", amazonLatestDate.String); err == nil {
				stat.AmazonLatestDate = &parsedTime
			} else if parsedTime, err := time.Parse("2006-01-02", amazonLatestDate.String); err == nil {
				stat.AmazonLatestDate = &parsedTime
			}
		}
		if amazonNextDate.Valid {
			if parsedTime, err := time.Parse("2006-01-02 15:04:05+00:00", amazonNextDate.String); err == nil {
				stat.AmazonNextDate = &parsedTime
			} else if parsedTime, err := time.Parse("2006-01-02", amazonNextDate.String); err == nil {
				stat.AmazonNextDate = &parsedTime
			}
		}

		stats = append(stats, stat)
	}

	return stats, rows.Err()
}

// UpsertSeries inserts or updates a series
func (s *Service) UpsertSeries(title, audibleID, audibleURL, amazonASIN string) (*Series, error) {
	// First try to get existing series
	var series Series
	query := `SELECT id, title, audible_id, audible_url, amazon_asin, created_at, updated_at 
	          FROM series WHERE title = ?`

	err := s.db.QueryRow(query, title).Scan(
		&series.ID, &series.Title, &series.AudibleID, &series.AudibleURL,
		&series.AmazonASIN, &series.CreatedAt, &series.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Insert new series
		insertQuery := `INSERT INTO series (title, audible_id, audible_url, amazon_asin) 
		                VALUES (?, ?, ?, ?) RETURNING id, created_at, updated_at`

		err = s.db.QueryRow(insertQuery, title, nilIfEmpty(audibleID),
			nilIfEmpty(audibleURL), nilIfEmpty(amazonASIN)).Scan(
			&series.ID, &series.CreatedAt, &series.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert series: %w", err)
		}

		series.Title = title
		series.AudibleID = nilIfEmpty(audibleID)
		series.AudibleURL = nilIfEmpty(audibleURL)
		series.AmazonASIN = nilIfEmpty(amazonASIN)

	} else if err != nil {
		return nil, fmt.Errorf("failed to query series: %w", err)
	} else {
		// Update existing series
		updateQuery := `UPDATE series SET audible_id = ?, audible_url = ?, amazon_asin = ?, updated_at = CURRENT_TIMESTAMP 
		                WHERE id = ?`

		_, err = s.db.Exec(updateQuery, nilIfEmpty(audibleID), nilIfEmpty(audibleURL),
			nilIfEmpty(amazonASIN), series.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to update series: %w", err)
		}

		// Update local struct
		series.AudibleID = nilIfEmpty(audibleID)
		series.AudibleURL = nilIfEmpty(audibleURL)
		series.AmazonASIN = nilIfEmpty(amazonASIN)
		series.UpdatedAt = time.Now()
	}

	return &series, nil
}

// UpdateSeriesBooks updates all books for a series from scraped data
func (s *Service) UpdateSeriesBooks(seriesID int, provider string, info models.SeriesInfo) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update scraped counts in series table
	if provider == ProviderAudible {
		_, err = tx.Exec(`UPDATE series SET audible_scraped_count = ? WHERE id = ?`, info.AudibleCount, seriesID)
		if err != nil {
			return fmt.Errorf("failed to update audible scraped count: %w", err)
		}
	} else if provider == ProviderAmazon {
		_, err = tx.Exec(`UPDATE series SET amazon_scraped_count = ? WHERE id = ?`, info.AmazonCount, seriesID)
		if err != nil {
			return fmt.Errorf("failed to update amazon scraped count: %w", err)
		}
	}

	// Clear existing books for this series/provider
	_, err = tx.Exec(`DELETE FROM books WHERE series_id = ? AND provider = ?`, seriesID, provider)
	if err != nil {
		return fmt.Errorf("failed to clear existing books: %w", err)
	}

	// Insert books based on provider
	if provider == ProviderAudible && info.AudibleCount > 0 {
		// Insert all books in the series (placeholder titles for book 1 to count-1)
		for i := 1; i <= info.AudibleCount; i++ {
			title := fmt.Sprintf("Book %d", i)
			isLatest := (i == info.AudibleCount)
			var releaseDate *time.Time

			// Use actual title and date for the latest book if available
			if isLatest && info.AudibleLatestTitle != "" {
				title = info.AudibleLatestTitle
			}
			if isLatest {
				releaseDate = info.AudibleLatestDate
			}

			_, err = tx.Exec(`
				INSERT INTO books (series_id, provider, title, book_number, release_date, is_latest) 
				VALUES (?, ?, ?, ?, ?, ?)`,
				seriesID, provider, title, i, releaseDate, isLatest)
			if err != nil {
				return fmt.Errorf("failed to insert audible book %d: %w", i, err)
			}
		}

		// Insert next book if it's a preorder
		if info.AudibleNextDate != nil {
			nextTitle := info.AudibleNextTitle
			if nextTitle == "" {
				nextTitle = fmt.Sprintf("Book %d", info.AudibleCount+1)
			}
			_, err = tx.Exec(`
				INSERT INTO books (series_id, provider, title, book_number, release_date, is_preorder) 
				VALUES (?, ?, ?, ?, ?, ?)`,
				seriesID, provider, nextTitle, info.AudibleCount+1,
				info.AudibleNextDate, true)
			if err != nil {
				return fmt.Errorf("failed to insert audible next book: %w", err)
			}
		}
	}

	if provider == ProviderAmazon && info.AmazonCount > 0 {
		// Insert all books in the series (placeholder titles for book 1 to count-1)
		for i := 1; i <= info.AmazonCount; i++ {
			title := fmt.Sprintf("Book %d", i)
			isLatest := (i == info.AmazonCount)
			var releaseDate *time.Time

			// Use actual title and date for the latest book if available
			if isLatest && info.AmazonLatestTitle != "" {
				title = info.AmazonLatestTitle
			}
			if isLatest {
				releaseDate = info.AmazonLatestDate
			}

			_, err = tx.Exec(`
				INSERT INTO books (series_id, provider, title, book_number, release_date, is_latest) 
				VALUES (?, ?, ?, ?, ?, ?)`,
				seriesID, provider, title, i, releaseDate, isLatest)
			if err != nil {
				return fmt.Errorf("failed to insert amazon book %d: %w", i, err)
			}
		}

		// Insert next book if it's a preorder
		if info.AmazonNextDate != nil {
			nextTitle := info.AmazonNextTitle
			if nextTitle == "" {
				nextTitle = fmt.Sprintf("Book %d", info.AmazonCount+1)
			}
			_, err = tx.Exec(`
				INSERT INTO books (series_id, provider, title, book_number, release_date, is_preorder) 
				VALUES (?, ?, ?, ?, ?, ?)`,
				seriesID, provider, nextTitle, info.AmazonCount+1,
				info.AmazonNextDate, true)
			if err != nil {
				return fmt.Errorf("failed to insert amazon next book: %w", err)
			}
		}
	}

	return tx.Commit()
}

// GetRuntimeSetting gets a runtime setting value from the database
func (s *Service) GetRuntimeSetting(key string) (string, error) {
	var value string
	query := `SELECT value FROM runtime_settings WHERE key = ?`
	err := s.db.QueryRow(query, key).Scan(&value)
	if err != nil {
		return "", fmt.Errorf("failed to get runtime setting %s: %w", key, err)
	}
	return value, nil
}

// SetRuntimeSetting updates or creates a runtime setting in the database
func (s *Service) SetRuntimeSetting(key, value string) error {
	query := `INSERT OR REPLACE INTO runtime_settings (key, value, updated_at) 
	          VALUES (?, ?, CURRENT_TIMESTAMP)`
	_, err := s.db.Exec(query, key, value)
	if err != nil {
		return fmt.Errorf("failed to set runtime setting %s: %w", key, err)
	}
	return nil
}

// CreateScrapeJob creates a new scrape job
func (s *Service) CreateScrapeJob(seriesID int, provider string) (*ScrapeJob, error) {
	query := `INSERT INTO scrape_jobs (series_id, provider) VALUES (?, ?) 
	          RETURNING id, created_at`

	var job ScrapeJob
	err := s.db.QueryRow(query, seriesID, provider).Scan(&job.ID, &job.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create scrape job: %w", err)
	}

	job.SeriesID = seriesID
	job.Provider = provider
	job.Status = JobStatusPending

	return &job, nil
}

// UpdateScrapeJob updates a scrape job status
func (s *Service) UpdateScrapeJob(jobID int, status string, errorMsg *string, bookCount int) error {
	var query string
	var args []interface{}

	now := time.Now()

	switch status {
	case JobStatusRunning:
		query = `UPDATE scrape_jobs SET status = ?, started_at = ? WHERE id = ?`
		args = []interface{}{status, now, jobID}
	case JobStatusCompleted:
		query = `UPDATE scrape_jobs SET status = ?, completed_at = ?, book_count = ? WHERE id = ?`
		args = []interface{}{status, now, bookCount, jobID}
	case JobStatusFailed:
		query = `UPDATE scrape_jobs SET status = ?, completed_at = ?, error_message = ? WHERE id = ?`
		args = []interface{}{status, now, errorMsg, jobID}
	default:
		query = `UPDATE scrape_jobs SET status = ? WHERE id = ?`
		args = []interface{}{status, jobID}
	}

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update scrape job: %w", err)
	}

	return nil
}

// GetPendingScrapeJobs returns all pending scrape jobs
func (s *Service) GetPendingScrapeJobs() ([]ScrapeJob, error) {
	query := `SELECT id, series_id, provider, status, started_at, completed_at, 
	                 error_message, book_count, created_at 
	          FROM scrape_jobs WHERE status = ? ORDER BY created_at`

	rows, err := s.db.Query(query, JobStatusPending)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending jobs: %w", err)
	}
	defer rows.Close()

	var jobs []ScrapeJob
	for rows.Next() {
		var job ScrapeJob
		err := rows.Scan(&job.ID, &job.SeriesID, &job.Provider, &job.Status,
			&job.StartedAt, &job.CompletedAt, &job.ErrorMessage,
			&job.BookCount, &job.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan scrape job: %w", err)
		}
		jobs = append(jobs, job)
	}

	return jobs, rows.Err()
}

// HasActiveScrapeJob checks if there's already an active job for a series/provider
func (s *Service) HasActiveScrapeJob(seriesID int, provider string) (bool, error) {
	query := `SELECT COUNT(*) FROM scrape_jobs 
	          WHERE series_id = ? AND provider = ? AND status IN (?, ?)`

	var count int
	err := s.db.QueryRow(query, seriesID, provider, JobStatusPending, JobStatusRunning).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check active jobs: %w", err)
	}

	return count > 0, nil
}

// GetSeriesByTitle returns a series by title
func (s *Service) GetSeriesByTitle(title string) (*Series, error) {
	query := `SELECT id, title, audible_id, audible_url, amazon_asin, created_at, updated_at 
	          FROM series WHERE title = ?`

	var series Series
	err := s.db.QueryRow(query, title).Scan(
		&series.ID, &series.Title, &series.AudibleID, &series.AudibleURL,
		&series.AmazonASIN, &series.CreatedAt, &series.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to query series: %w", err)
	}

	return &series, nil
}

// GetSeriesByID returns a series by ID
func (s *Service) GetSeriesByID(id int) (*Series, error) {
	query := `SELECT id, title, audible_id, audible_url, amazon_asin, created_at, updated_at 
	          FROM series WHERE id = ?`

	var series Series
	err := s.db.QueryRow(query, id).Scan(
		&series.ID, &series.Title, &series.AudibleID, &series.AudibleURL,
		&series.AmazonASIN, &series.CreatedAt, &series.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to query series by ID: %w", err)
	}

	return &series, nil
}

// ClearAllBookData removes all book data to prevent cascading corruption
func (s *Service) ClearAllBookData() error {
	_, err := s.db.Exec(`DELETE FROM books`)
	if err != nil {
		return fmt.Errorf("failed to clear all book data: %w", err)
	}
	log.Printf("cleared all book data from database")
	return nil
}

// CleanupStaleRunningJobs marks all running and old pending jobs as failed (for startup cleanup)
func (s *Service) CleanupStaleRunningJobs() error {
	errorMsg := "job interrupted by application restart"

	// Clean up running jobs
	query1 := `UPDATE scrape_jobs SET status = ?, completed_at = CURRENT_TIMESTAMP, error_message = ? WHERE status = ?`
	result1, err := s.db.Exec(query1, JobStatusFailed, errorMsg, JobStatusRunning)
	if err != nil {
		return fmt.Errorf("failed to cleanup stale running jobs: %w", err)
	}

	runningCleaned, _ := result1.RowsAffected()

	// Clean up ALL pending jobs on startup (they were never processed)
	query2 := `UPDATE scrape_jobs SET status = ?, completed_at = CURRENT_TIMESTAMP, error_message = ? 
	          WHERE status = ?`
	result2, err := s.db.Exec(query2, JobStatusFailed, "job cleared on application restart", JobStatusPending)
	if err != nil {
		return fmt.Errorf("failed to cleanup stale pending jobs: %w", err)
	}

	pendingCleaned, _ := result2.RowsAffected()

	if runningCleaned > 0 || pendingCleaned > 0 {
		fmt.Printf("cleaned up %d stale running jobs and %d stale pending jobs\n", runningCleaned, pendingCleaned)
	}

	return nil
}

// GetLastScrapeTime returns the most recent scrape start time (when any scrape was initiated)
func (s *Service) GetLastScrapeTime() (*time.Time, error) {
	query := `SELECT MAX(started_at) FROM scrape_jobs WHERE started_at IS NOT NULL`
	var lastScrape *time.Time
	err := s.db.QueryRow(query).Scan(&lastScrape)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return lastScrape, err
}

// GetAllSeries returns all series from the database
func (s *Service) GetAllSeries() ([]Series, error) {
	query := `SELECT id, title, audible_id, audible_url, amazon_asin, created_at, updated_at FROM series ORDER BY title`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var series []Series
	for rows.Next() {
		var s Series
		err := rows.Scan(&s.ID, &s.Title, &s.AudibleID, &s.AudibleURL, &s.AmazonASIN, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		series = append(series, s)
	}
	return series, rows.Err()
}

// DeleteSeries deletes a series and all its associated data
func (s *Service) DeleteSeries(seriesID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete scrape jobs first (foreign key constraint)
	_, err = tx.Exec("DELETE FROM scrape_jobs WHERE series_id = ?", seriesID)
	if err != nil {
		return err
	}

	// Delete books (foreign key constraint)
	_, err = tx.Exec("DELETE FROM books WHERE series_id = ?", seriesID)
	if err != nil {
		return err
	}

	// Finally delete the series
	_, err = tx.Exec("DELETE FROM series WHERE id = ?", seriesID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteSeriesByTitle deletes a series (and associated data) by its title
func (s *Service) DeleteSeriesByTitle(title string) error {
	// Look up the series by title first
	series, err := s.GetSeriesByTitle(title)
	if err != nil {
		return err
	}
	if series == nil {
		return sql.ErrNoRows
	}
	// Reuse the existing deletion logic by ID
	return s.DeleteSeries(series.ID)
}

// nilIfEmpty returns nil for empty strings, otherwise returns pointer to string
func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
