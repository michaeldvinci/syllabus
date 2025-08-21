package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/michaeldvinci/syllabus/internal/cache"
	"github.com/michaeldvinci/syllabus/internal/models"
)

// App holds the application state
type App struct {
	Provider    models.Provider
	Cache       *cache.Cache
	Data        []models.SeriesIDs
	RefreshChan chan bool
	mu          sync.RWMutex // Protect Data updates
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
	Rows []Row
	Now  string
}

// HandleIndex serves the main HTML page
func (a *App) HandleIndex(w http.ResponseWriter, r *http.Request) {
	var rows []Row

	audibleLatest := ""
	audibleNext := ""
	loc, _ := time.LoadLocation("America/Chicago")
	now := time.Now().In(loc)
	y, m, d := now.Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, loc)
	infos := a.collectAll()
	for _, info := range infos {
		if info.AudibleLatestDate != nil {
			audibleLatestS := info.AudibleLatestDate.Format("2006-01-02")

			other, err := time.ParseInLocation("2006-01-02", audibleLatestS, loc)
			if err != nil {
				panic(err)
			}

			switch {
			case other.Before(today):
				audibleLatest = audibleLatestS
				audibleNext = ""
			case other.After(today):
				audibleLatest = ""
				audibleNext = audibleLatestS
			default:
				fmt.Println(audibleLatestS, "is TODAY")
			}
		}
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
			AmazonLatest:  joinTitleDate(info.AmazonLatestTitle, info.AmazonLatestDate),
			AmazonNext:    joinTitleDate(info.AmazonNextTitle, info.AmazonNextDate),
			AudibleURL:    audURL,
			AmazonURL:     amzURL,
		})
	}

	tpl := template.Must(template.New("idx").Parse(IndexHTML))
	if err := tpl.Execute(w, Page{Rows: rows, Now: time.Now().Format(time.RFC822)}); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// HandleAPI serves the JSON API endpoint
func (a *App) HandleAPI(w http.ResponseWriter, r *http.Request) {
	infos := a.collectAll()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(infos)
}

func (a *App) collectAll() []models.SeriesInfo {
	a.mu.RLock()
	data := make([]models.SeriesIDs, len(a.Data))
	copy(data, a.Data)
	a.mu.RUnlock()
	
	var wg sync.WaitGroup
	infos := make([]models.SeriesInfo, len(data))
	for i, e := range data {
		wg.Add(1)
		go func(i int, e models.SeriesIDs) {
			defer wg.Done()
			key := e.Title + "|" + e.AudibleID + "|" + e.AmazonASIN
			if v, ok := a.Cache.Get(key); ok {
				infos[i] = v
				return
			}
			info, err := a.Provider.Fetch(e)
			if err != nil {
				info.Err = err
			}
			a.Cache.Set(key, info)
			infos[i] = info
		}(i, e)
	}
	wg.Wait()
	return infos
}

// WarmupCache fetches all series data at startup to populate the cache
func (a *App) WarmupCache() {
	_ = a.collectAll()
}

// HandleEvents serves Server-Sent Events for live refresh
func (a *App) HandleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	// Send initial connection event
	fmt.Fprintf(w, "data: connected\n\n")
	w.(http.Flusher).Flush()
	
	// Listen for refresh signals
	for {
		select {
		case <-r.Context().Done():
			return
		case <-a.RefreshChan:
			fmt.Fprintf(w, "data: refresh\n\n")
			w.(http.Flusher).Flush()
		}
	}
}

// findNewEntries compares old and new data to identify new series
func (a *App) findNewEntries(newData []models.SeriesIDs) []models.SeriesIDs {
	a.mu.RLock()
	oldData := a.Data
	a.mu.RUnlock()
	
	// Create a map of existing entries for fast lookup
	existing := make(map[string]bool)
	for _, series := range oldData {
		key := series.Title + "|" + series.AudibleID + "|" + series.AmazonASIN
		existing[key] = true
	}
	
	// Find entries that don't exist in the old data
	var newEntries []models.SeriesIDs
	for _, series := range newData {
		key := series.Title + "|" + series.AudibleID + "|" + series.AmazonASIN
		if !existing[key] {
			newEntries = append(newEntries, series)
		}
	}
	
	return newEntries
}

// UpdateDataIncremental adds only new entries and triggers refresh
func (a *App) UpdateDataIncremental(newData []models.SeriesIDs) {
	newEntries := a.findNewEntries(newData)
	
	if len(newEntries) == 0 {
		log.Printf("no new entries found in config update")
		return
	}
	
	log.Printf("found %d new entries to scrape", len(newEntries))
	
	// Update the data atomically
	a.mu.Lock()
	a.Data = newData
	a.mu.Unlock()
	
	// Scrape only the new entries in background
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
		
		// Trigger refresh for all connected clients
		select {
		case a.RefreshChan <- true:
		default:
		}
		log.Printf("incremental update complete - added %d new entries", len(newEntries))
	}()
}

// ReloadData clears cache and reloads series data (fallback for major changes)
func (a *App) ReloadData(newData []models.SeriesIDs) {
	a.Cache.Clear()
	a.mu.Lock()
	a.Data = newData
	a.mu.Unlock()
	// Trigger refresh for all connected clients
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