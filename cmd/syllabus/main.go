package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/michaeldvinci/syllabus/internal/cache"
	"github.com/michaeldvinci/syllabus/internal/handlers"
	"github.com/michaeldvinci/syllabus/internal/models"
	"github.com/michaeldvinci/syllabus/internal/scrapers"
	"github.com/michaeldvinci/syllabus/internal/utils"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <path-to-yaml>", os.Args[0])
	}
	
	path := os.Args[1]
	cfg, err := utils.LoadConfig(path)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	
	series := utils.ToSeriesIDs(cfg.Audiobooks)

	// Create HTTP client for scrapers
	httpClient := &http.Client{Timeout: 12 * time.Second}

	// Initialize providers
	provider := &scrapers.CompositeProvider{
		Providers: []models.Provider{
			&scrapers.AmazonPAAPIProvider{Enabled: false},
			&scrapers.AmazonScraperProvider{Enabled: true, Client: httpClient},
			&scrapers.AudibleScraperProvider{Enabled: true, Client: httpClient},
		},
	}

	// Initialize application
	app := &handlers.App{
		Provider: provider,
		Cache:    cache.NewCache(6 * time.Hour),
		Data:     series,
	}

	// Setup HTTP routes
	http.HandleFunc("/", app.HandleIndex)
	http.HandleFunc("/api/series", app.HandleAPI)
	
	// Start server
	addr := ":8080"
	log.Printf("listening on %s …", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}