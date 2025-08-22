package scraper

import (
	"context"
	"log"
	"sync"

	"github.com/michaeldvinci/syllabus/internal/database"
	"github.com/michaeldvinci/syllabus/internal/models"
)

// BackgroundScraper handles background scraping operations
type BackgroundScraper struct {
	provider models.Provider
	db       *database.Service
	jobChan  chan database.ScrapeJob
	done     chan struct{}
	wg       sync.WaitGroup
	
	// For notifying UI of updates
	updateChan chan SeriesUpdate
}

// SeriesUpdate represents a series update event
type SeriesUpdate struct {
	SeriesID int    `json:"series_id"`
	Title    string `json:"title"`
	Provider string `json:"provider"`
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
}

// NewBackgroundScraper creates a new background scraper
func NewBackgroundScraper(provider models.Provider, db *database.Service) *BackgroundScraper {
	return &BackgroundScraper{
		provider:   provider,
		db:         db,
		jobChan:    make(chan database.ScrapeJob, 100), // Buffer for jobs
		done:       make(chan struct{}),
		updateChan: make(chan SeriesUpdate, 100),
	}
}

// Start begins the background scraper workers
func (bs *BackgroundScraper) Start(ctx context.Context, workers int) {
	log.Printf("starting %d background scraper workers", workers)
	
	for i := 0; i < workers; i++ {
		bs.wg.Add(1)
		go bs.worker(ctx, i)
	}
	
	// No dispatcher needed - jobs are queued directly
}

// Stop gracefully stops the background scraper
func (bs *BackgroundScraper) Stop() {
	log.Printf("stopping background scraper")
	close(bs.done)
	bs.wg.Wait()
	close(bs.updateChan)
}

// GetUpdateChannel returns the channel for UI updates
func (bs *BackgroundScraper) GetUpdateChannel() <-chan SeriesUpdate {
	return bs.updateChan
}

// QueueSeriesUpdate queues a scraping job for a specific series
func (bs *BackgroundScraper) QueueSeriesUpdate(seriesID int, provider string) error {
	// Check if there's already an active job for this series/provider
	hasActive, err := bs.db.HasActiveScrapeJob(seriesID, provider)
	if err != nil {
		return err
	}
	
	if hasActive {
		log.Printf("skipping scrape job for series %d (%s) - already active", seriesID, provider)
		return nil
	}
	
	job, err := bs.db.CreateScrapeJob(seriesID, provider)
	if err != nil {
		return err
	}
	
	// Send job to channel (non-blocking)
	select {
	case bs.jobChan <- *job:
		log.Printf("queued scrape job %d for series %d (%s)", job.ID, seriesID, provider)
	default:
		log.Printf("job queue full, dropping scrape job for series %d", seriesID)
	}
	
	return nil
}

// Removed dispatcher - jobs are queued directly via QueueSeriesUpdate

// worker processes scraping jobs
func (bs *BackgroundScraper) worker(ctx context.Context, workerID int) {
	defer bs.wg.Done()
	
	log.Printf("background scraper worker %d started", workerID)
	
	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %d stopping due to context cancellation", workerID)
			return
		case <-bs.done:
			log.Printf("worker %d stopping", workerID)
			return
		case job := <-bs.jobChan:
			bs.processJob(workerID, job)
		}
	}
}

// processJob processes a single scraping job
func (bs *BackgroundScraper) processJob(workerID int, job database.ScrapeJob) {
	log.Printf("worker %d processing job %d for series %d (%s)", workerID, job.ID, job.SeriesID, job.Provider)
	
	// Mark job as running
	if err := bs.db.UpdateScrapeJob(job.ID, database.JobStatusRunning, nil, 0); err != nil {
		log.Printf("error updating job status to running: %v", err)
		return
	}
	
	// Get series details to construct SeriesIDs
	series, err := bs.getSeriesDetails(job.SeriesID)
	if err != nil {
		errMsg := err.Error()
		bs.db.UpdateScrapeJob(job.ID, database.JobStatusFailed, &errMsg, 0)
		bs.notifyUpdate(job.SeriesID, series.Title, job.Provider, "failed", errMsg)
		return
	}
	
	// Create SeriesIDs for the provider
	seriesIDs := models.SeriesIDs{
		Title:      series.Title,
		AudibleID:  stringValue(series.AudibleID),
		AudibleURL: stringValue(series.AudibleURL),
		AmazonASIN: stringValue(series.AmazonASIN),
	}
	
	// Perform the scraping
	info, err := bs.provider.Fetch(seriesIDs)
	if err != nil {
		log.Printf("worker %d failed to scrape series %d: %v", workerID, job.SeriesID, err)
		// Create empty info to clear stale data from database
		emptyInfo := models.SeriesInfo{
			Title: series.Title,
			// All other fields will be zero values (0, nil) which will clear the database
		}
		// Update database with empty data to clear stale entries
		if err := bs.db.UpdateSeriesBooks(job.SeriesID, job.Provider, emptyInfo); err != nil {
			log.Printf("worker %d failed to clear stale data for series %d: %v", workerID, job.SeriesID, err)
		} else {
			log.Printf("worker %d cleared stale data for series %d (%s)", workerID, job.SeriesID, job.Provider)
		}
		
		errMsg := err.Error()
		bs.db.UpdateScrapeJob(job.ID, database.JobStatusFailed, &errMsg, 0)
		bs.notifyUpdate(job.SeriesID, series.Title, job.Provider, "failed", errMsg)
		return
	}
	
	// Update database with scraped data
	if err := bs.db.UpdateSeriesBooks(job.SeriesID, job.Provider, info); err != nil {
		errMsg := err.Error()
		bs.db.UpdateScrapeJob(job.ID, database.JobStatusFailed, &errMsg, 0)
		bs.notifyUpdate(job.SeriesID, series.Title, job.Provider, "failed", errMsg)
		log.Printf("worker %d failed to update series %d books: %v", workerID, job.SeriesID, err)
		return
	}
	
	// Calculate book count for the provider
	bookCount := 0
	if job.Provider == database.ProviderAudible {
		bookCount = info.AudibleCount
	} else if job.Provider == database.ProviderAmazon {
		bookCount = info.AmazonCount
	}
	
	// Mark job as completed
	if err := bs.db.UpdateScrapeJob(job.ID, database.JobStatusCompleted, nil, bookCount); err != nil {
		log.Printf("error updating job status to completed: %v", err)
	}
	
	// Notify UI of successful update
	bs.notifyUpdate(job.SeriesID, series.Title, job.Provider, "completed", "")
	log.Printf("worker %d successfully scraped series %d (%s) - %d books", workerID, job.SeriesID, job.Provider, bookCount)
}

// getSeriesDetails fetches series details from database
func (bs *BackgroundScraper) getSeriesDetails(seriesID int) (*database.Series, error) {
	return bs.db.GetSeriesByID(seriesID)
}

// notifyUpdate sends update notification to UI
func (bs *BackgroundScraper) notifyUpdate(seriesID int, title, provider, status, errorMsg string) {
	update := SeriesUpdate{
		SeriesID: seriesID,
		Title:    title,
		Provider: provider,
		Status:   status,
		Error:    errorMsg,
	}
	
	select {
	case bs.updateChan <- update:
		// Update sent successfully
	default:
		log.Printf("update channel full, dropping update for series %d", seriesID)
	}
}

// QueueAllSeriesUpdate queues scraping jobs for all series in the database
func (bs *BackgroundScraper) QueueAllSeriesUpdate() error {
	stats, err := bs.db.GetAllSeriesStats()
	if err != nil {
		return err
	}
	
	for _, stat := range stats {
		// Queue both providers if they have data
		if stat.AudibleID != nil {
			if err := bs.QueueSeriesUpdate(stat.ID, database.ProviderAudible); err != nil {
				log.Printf("error queuing audible job for series %d: %v", stat.ID, err)
			}
		}
		
		if stat.AmazonASIN != nil {
			if err := bs.QueueSeriesUpdate(stat.ID, database.ProviderAmazon); err != nil {
				log.Printf("error queuing amazon job for series %d: %v", stat.ID, err)
			}
		}
	}
	
	return nil
}

// stringValue safely converts *string to string
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// CleanupStaleJobs marks any running jobs as failed (for startup cleanup)
func (bs *BackgroundScraper) CleanupStaleJobs() error {
	return bs.db.CleanupStaleRunningJobs()
}