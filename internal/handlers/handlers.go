package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/michaeldvinci/syllabus/internal/cache"
	"github.com/michaeldvinci/syllabus/internal/models"
)

// App holds the application state
type App struct {
	Provider models.Provider
	Cache    *cache.Cache
	Data     []models.SeriesIDs
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
	var wg sync.WaitGroup
	infos := make([]models.SeriesInfo, len(a.Data))
	for i, e := range a.Data {
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