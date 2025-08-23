package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/michaeldvinci/syllabus/internal/auth"
	"github.com/michaeldvinci/syllabus/internal/cache"
	"github.com/michaeldvinci/syllabus/internal/database"
	"github.com/michaeldvinci/syllabus/internal/models"
	"github.com/michaeldvinci/syllabus/internal/scraper"
	"github.com/michaeldvinci/syllabus/internal/utils"
)

// App holds the application state
type App struct {
	Provider          models.Provider
	DB                *database.Service
	Cache             *cache.Cache
	Data              []models.SeriesIDs
	RefreshChan       chan bool
	ScraperUpdateCh   <-chan scraper.SeriesUpdate // Channel for scraper updates
	BackgroundScraper *scraper.BackgroundScraper  // Reference to background scraper
	Settings          models.Settings             // Application settings from config
	mu                sync.RWMutex                // Protect Data updates

	// Auto-refresh functionality
	autoRefreshInterval int          // Hours between auto-refreshes
	autoRefreshTicker   *time.Ticker // Ticker for auto-refresh
	autoRefreshMu       sync.RWMutex // Protect auto-refresh settings
}

// Row represents a table row in the HTML template
type Row struct {
	Title         string
	AudibleCount  int
	AudibleLatest string
	AudibleNext   string
	AmazonCount   int
	AmazonLatest  string
	AmazonNext    string
	AudibleURL    string
	AmazonURL     string
}

// Page represents the complete page data for the HTML template
type Page struct {
	Rows          []Row
	Now           string
	CalendarURL   string
	User          *auth.User
	Authenticated bool
	LastScrape    string
}

// HandleIndex serves the main HTML page
func (a *App) HandleIndex(w http.ResponseWriter, r *http.Request) {
	var rows []Row

	infos := a.collectAll()
	for _, info := range infos {
		audibleLatest := formatDateOnly(info.AudibleLatestDate)
		audibleNext := formatDateOnly(info.AudibleNextDate)
		audURL := ""
		if info.AudibleID != "" {
			audURL = fmt.Sprintf("https://www.audible.com/series/%s", info.AudibleID)
		}
		amzURL := ""
		if info.AmazonASIN != "" {
			amzURL = fmt.Sprintf("https://www.amazon.com/dp/%s", info.AmazonASIN)
		}
		rows = append(rows, Row{
			Title:         info.Title,
			AudibleCount:  info.AudibleCount,
			AudibleLatest: audibleLatest,
			AudibleNext:   audibleNext,
			AmazonCount:   info.AmazonCount,
			AmazonLatest:  formatDateOnly(info.AmazonLatestDate),
			AmazonNext:    formatDateOnly(info.AmazonNextDate),
			AudibleURL:    audURL,
			AmazonURL:     amzURL,
		})
	}

	// Generate the calendar URL based on the request
	calendarURL := fmt.Sprintf("%s://%s/calendar.ics",
		func() string {
			if r.TLS != nil {
				return "https"
			}
			return "http"
		}(), r.Host)

	// Get current user from context if available
	user, authenticated := auth.GetUserFromContext(r)

	// Get last scrape timestamp
	lastScrapeTime, err := a.DB.GetLastScrapeTime()
	lastScrapeStr := "Never"
	if err == nil && lastScrapeTime != nil {
		lastScrapeStr = lastScrapeTime.Format("Jan 2, 2006 at 3:04 PM")
	}

	// Helpers (maxWidth kept for backward-compat; not used in current template)
	funcs := template.FuncMap{
		"maxWidth": func(px int) string { return fmt.Sprintf("%dpx", px) },
	}

	tpl, err := template.New("idx").Funcs(funcs).Parse(IndexHTML)
	if err != nil {
		log.Printf("parse template: %v", err)
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(w, Page{
		Rows:          rows,
		Now:           time.Now().Format(time.RFC822),
		CalendarURL:   calendarURL,
		User:          user,
		Authenticated: authenticated,
		LastScrape:    lastScrapeStr,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleAPI serves the JSON API endpoint
func (a *App) HandleAPI(w http.ResponseWriter, r *http.Request) {
	infos := a.collectAll()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(infos)
}

// HandleScrapeStatus returns the current scrape job status
func (a *App) HandleScrapeStatus(w http.ResponseWriter, r *http.Request) {
	// Get count of active scrape jobs (running or pending)
	jobs, err := a.DB.GetPendingScrapeJobs()
	if err != nil {
		http.Error(w, "Failed to get scrape status", http.StatusInternalServerError)
		return
	}

	activeJobs := 0
	for _, job := range jobs {
		if job.Status == "running" || job.Status == "pending" {
			activeJobs++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"activeJobs": activeJobs,
	})
}

// HandleICal serves the iCal export endpoint
func (a *App) HandleICal(w http.ResponseWriter, r *http.Request) {
	infos := a.collectAll()
	icalContent := utils.GenerateICal(infos)

	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"book-releases.ics\"")
	w.Write([]byte(icalContent))
}

// HandleRefresh triggers a re-scrape of all series data
func (a *App) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get all series for re-scraping
	stats, err := a.DB.GetAllSeriesStats()
	if err != nil {
		log.Printf("error fetching series for refresh: %v", err)
		http.Error(w, "Failed to fetch series data", http.StatusInternalServerError)
		return
	}

	jobsQueued := 0
	for _, stat := range stats {
		if stat.AudibleID != nil {
			jobsQueued++
		}
		if stat.AmazonASIN != nil {
			jobsQueued++
		}
	}

	log.Printf("clearing all book data before refresh to prevent corruption")
	if err := a.DB.ClearAllBookData(); err != nil {
		log.Printf("warning: failed to clear book data: %v", err)
	}

	if a.BackgroundScraper != nil {
		if err := a.BackgroundScraper.QueueAllSeriesUpdate(); err != nil {
			log.Printf("error queuing refresh jobs: %v", err)
			http.Error(w, "Failed to queue refresh jobs", http.StatusInternalServerError)
			return
		}
		log.Printf("queued refresh jobs for all series")
		select {
		case a.RefreshChan <- true:
		default:
		}
	} else {
		log.Printf("background scraper not available")
		http.Error(w, "Background scraper not available", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success":     true,
		"seriesCount": len(stats),
		"jobsQueued":  jobsQueued,
		"message":     "Refresh triggered successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func (a *App) collectAll() []models.SeriesInfo {
	stats, err := a.DB.GetAllSeriesStats()
	if err != nil {
		log.Printf("error fetching series stats from database: %v", err)
		return []models.SeriesInfo{}
	}
	return database.ToSeriesInfoSlice(stats)
}

// WarmupCache is now a no-op since data comes from database
func (a *App) WarmupCache() {
	log.Printf("using database - no warmup needed")
}

// HandleEvents serves Server-Sent Events for live refresh
func (a *App) HandleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "data: connected\n\n")
	w.(http.Flusher).Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-a.RefreshChan:
			fmt.Fprintf(w, "data: refresh\n\n")
			w.(http.Flusher).Flush()
		case update := <-a.ScraperUpdateCh:
			updateJSON, _ := json.Marshal(update)
			fmt.Fprintf(w, "data: %s\n\n", updateJSON)
			w.(http.Flusher).Flush()
		}
	}
}

func (a *App) findNewEntries(newData []models.SeriesIDs) []models.SeriesIDs {
	a.mu.RLock()
	oldData := a.Data
	a.mu.RUnlock()

	existing := make(map[string]bool)
	for _, series := range oldData {
		key := series.Title + "|" + series.AudibleID + "|" + series.AmazonASIN
		existing[key] = true
	}

	var newEntries []models.SeriesIDs
	for _, series := range newData {
		key := series.Title + "|" + series.AudibleID + "|" + series.AmazonASIN
		if !existing[key] {
			newEntries = append(newEntries, series)
		}
	}
	return newEntries
}

func (a *App) UpdateDataIncremental(newData []models.SeriesIDs) {
	newEntries := a.findNewEntries(newData)
	if len(newEntries) == 0 {
		log.Printf("no new entries found in config update")
		return
	}

	a.mu.Lock()
	a.Data = newData
	a.mu.Unlock()

	go func() {
		log.Printf("scraping %d new entries...", len(newEntries))
		for _, entry := range newEntries {
			key := entry.Title + "|" + entry.AudibleID + "|" + entry.AmazonASIN
			if _, ok := a.Cache.Get(key); !ok {
				info, err := a.Provider.Fetch(entry)
				if err != nil {
					info.Err = err
				}
				a.Cache.Set(key, info)
				log.Printf("scraped new entry: %s", entry.Title)
			}
		}
		select {
		case a.RefreshChan <- true:
		default:
		}
		log.Printf("incremental update complete - added %d new entries", len(newEntries))
	}()
}

func (a *App) ReloadData(newData []models.SeriesIDs) {
	a.Cache.Clear()
	a.mu.Lock()
	a.Data = newData
	a.mu.Unlock()
	select {
	case a.RefreshChan <- true:
	default:
	}
}

func joinTitleDate(title string, d *time.Time) string {
	if title == "" && d == nil {
		return ""
	}
	if title != "" && d != nil {
		return fmt.Sprintf("%s — %s", title, d.Format("2006-01-02"))
	}
	if title != "" {
		return title
	}
	return d.Format("2006-01-02")
}

func formatDateOnly(d *time.Time) string {
	if d == nil {
		return ""
	}
	return d.Format("2006-01-02")
}

// Auto-refresh handlers (unchanged)
func (a *App) HandleAutoRefresh(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.autoRefreshMu.RLock()
		interval := a.autoRefreshInterval
		a.autoRefreshMu.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"interval": interval,
		})

	case http.MethodPost:
		// Update auto-refresh interval
		var req struct {
			Interval int `json:"interval"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Validate interval (2, 4, 6, 8, 10 hours)
		if req.Interval < 2 || req.Interval > 10 || req.Interval%2 != 0 {
			http.Error(w, "Invalid interval", http.StatusBadRequest)
			return
		}

		// Update the auto-refresh interval
		a.SetAutoRefreshInterval(req.Interval)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":  true,
			"interval": req.Interval,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *App) SetAutoRefreshInterval(hours int) {
	a.autoRefreshMu.Lock()
	defer a.autoRefreshMu.Unlock()
	if a.autoRefreshInterval == hours && a.autoRefreshTicker != nil {
		return
	}
	a.autoRefreshInterval = hours
	if err := a.DB.SetRuntimeSetting("auto_refresh_interval", strconv.Itoa(hours)); err != nil {
		log.Printf("warning: failed to save auto-refresh interval to database: %v", err)
	}
	if a.autoRefreshTicker != nil {
		a.autoRefreshTicker.Stop()
	}
	a.autoRefreshTicker = time.NewTicker(time.Duration(hours) * time.Hour)
	log.Printf("auto-refresh interval updated to %d hours", hours)
}

func (a *App) StartAutoRefresh() {
	a.autoRefreshMu.Lock()
	if a.autoRefreshInterval == 0 {
		if dbValue, err := a.DB.GetRuntimeSetting("auto_refresh_interval"); err == nil {
			if val, err := strconv.Atoi(dbValue); err == nil && val > 0 {
				a.autoRefreshInterval = val
				log.Printf("loaded auto-refresh interval from database: %d hours", val)
			}
		}
		if a.autoRefreshInterval == 0 {
			a.autoRefreshInterval = a.Settings.AutoRefreshInterval
			log.Printf("using config auto-refresh interval: %d hours", a.autoRefreshInterval)
		}
	}
	interval := a.autoRefreshInterval
	a.autoRefreshMu.Unlock()

	a.SetAutoRefreshInterval(interval)

	go func() {
		log.Printf("starting auto-refresh loop with %d hour interval", interval)
		for range a.autoRefreshTicker.C {
			log.Printf("triggering scheduled data refresh...")
			if err := a.DB.ClearAllBookData(); err != nil {
				log.Printf("warning: failed to clear book data during auto-refresh: %v", err)
			}
			if a.BackgroundScraper != nil {
				if err := a.BackgroundScraper.QueueAllSeriesUpdate(); err != nil {
					log.Printf("error queuing auto-refresh jobs: %v", err)
				} else {
					log.Printf("auto-refresh jobs queued successfully")
				}
			}
		}
	}()
}

func (a *App) StopAutoRefresh() {
	a.autoRefreshMu.Lock()
	defer a.autoRefreshMu.Unlock()
	if a.autoRefreshTicker != nil {
		a.autoRefreshTicker.Stop()
		a.autoRefreshTicker = nil
		log.Printf("auto-refresh stopped")
	}
}
