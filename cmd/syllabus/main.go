package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/michaeldvinci/syllabus/internal/auth"
	"github.com/michaeldvinci/syllabus/internal/cache"
	"github.com/michaeldvinci/syllabus/internal/database"
	"github.com/michaeldvinci/syllabus/internal/handlers"
	"github.com/michaeldvinci/syllabus/internal/models"
	"github.com/michaeldvinci/syllabus/internal/scraper"
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
	settings := cfg.GetSettings()
	
	// Apply environment variable overrides (takes precedence over YAML config)
	utils.ApplyEnvOverrides(&settings)
	
	log.Printf("configuration loaded - auto_refresh: %dh, workers: %d, port: %d, cache: %dh, log: %s", 
		settings.AutoRefreshInterval, settings.DefaultWorkers, settings.ServerPort, 
		settings.CacheTimeout, settings.LogLevel)

	// Initialize individual providers with fresh HTTP clients to prevent shared state
	audibleProvider := &scrapers.AudibleScraperProvider{
		Enabled: true, 
		Client: &http.Client{
			Timeout: 12 * time.Second,
			Jar:     nil, // Disable cookies to prevent session state bleeding
		},
	}
	
	amazonProvider := &scrapers.AmazonScraperProvider{
		Enabled: true, 
		Client: &http.Client{
			Timeout: 12 * time.Second,
			Jar:     nil, // Disable cookies to prevent session state bleeding
		},
	}
	
	// Create provider map for background scraper
	providers := map[string]models.Provider{
		"audible": audibleProvider,
		"amazon":  amazonProvider,
	}
	
	// Create composite provider for backward compatibility (if needed)
	provider := &scrapers.CompositeProvider{
		Providers: []models.Provider{
			&scrapers.AmazonPAAPIProvider{Enabled: false},
			amazonProvider,
			audibleProvider,
		},
	}

	// Initialize data directory and database
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("failed to create data directory: %v", err)
	}
	
	// Initialize database
	db, err := database.New(dataDir)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()
	
	dbService := database.NewService(db)
	
	// Initialize authentication store with file persistence
	authStore := auth.NewStoreWithFile(filepath.Join(dataDir, "users.json"))
	
	// Create default admin user if it doesn't exist
	_, err = authStore.CreateUserWithRole("admin", "admin", auth.RoleAdmin)
	if err != nil && err != auth.ErrUserExists {
		log.Fatalf("failed to create default admin user: %v", err)
	}
	if err == nil {
		log.Printf("created default admin user (username: admin, password: admin)")
	} else {
		log.Printf("loaded existing users from %s", filepath.Join(dataDir, "users.json"))
	}

	// Initialize authentication middleware and handlers
	authMiddleware := auth.NewMiddleware(authStore)
	authHandlers := auth.NewAuthHandlers(authStore)

	// Initialize background scraper with provider map
	backgroundScraper := scraper.NewBackgroundScraper(providers, dbService)
	
	// Initialize application
	app := &handlers.App{
		Provider:          provider,
		DB:                dbService,
		Cache:             cache.NewCache(time.Duration(settings.CacheTimeout) * time.Hour),
		Data:              series,
		RefreshChan:       make(chan bool, 1),
		ScraperUpdateCh:   backgroundScraper.GetUpdateChannel(),
		BackgroundScraper: backgroundScraper,
		Settings:          settings,
	}

	// Setup authentication routes (no middleware needed)
	http.HandleFunc("/login", authHandlers.HandleLogin)
	http.HandleFunc("/logout", authHandlers.HandleLogout)
	http.HandleFunc("/api/auth", authMiddleware.OptionalAuth(authHandlers.HandleAPI))
	
	// Setup admin-only routes
	http.HandleFunc("/api/users", authMiddleware.RequireAdmin(authHandlers.HandleListUsers))
	http.HandleFunc("/api/users/create", authMiddleware.RequireAdmin(authHandlers.HandleCreateUser))
	http.HandleFunc("/api/users/delete", authMiddleware.RequireAdmin(authHandlers.HandleDeleteUser))
	http.HandleFunc("/api/users/reset-password", authMiddleware.RequireAdmin(authHandlers.HandleResetPassword))

	// Setup protected HTTP routes with authentication middleware
	http.HandleFunc("/", authMiddleware.RequireAuth(app.HandleIndex))
	http.HandleFunc("/api/series", authMiddleware.RequireAuth(app.HandleAPI))
	http.HandleFunc("/api/scrape-status", authMiddleware.RequireAuth(app.HandleScrapeStatus))
	http.HandleFunc("/events", authMiddleware.RequireAuth(app.HandleEvents))
	http.HandleFunc("/calendar.ics", authMiddleware.RequireAuth(app.HandleICal))
	http.HandleFunc("/refresh", authMiddleware.RequireAuth(app.HandleRefresh))
	http.HandleFunc("/api/auto-refresh", authMiddleware.RequireAuth(app.HandleAutoRefresh))
	
	// Serve static files (favicon, logo) - check for local vs docker paths
	staticDir := "./app/res/"
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		staticDir = "./res/" // Docker path
		log.Printf("using docker static path: %s", staticDir)
	} else {
		log.Printf("using local static path: %s", staticDir)
	}
	
	// List files in static directory for debugging
	if files, err := os.ReadDir(staticDir); err == nil {
		log.Printf("static files available:")
		for _, file := range files {
			log.Printf("  - %s", file.Name())
		}
	}
	
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	
	// Setup file watcher for config changes
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("failed to create file watcher: %v", err)
	}
	defer watcher.Close()
	
	// Watch the config file
	configPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("failed to get absolute path: %v", err)
	}
	
	err = watcher.Add(configPath)
	if err != nil {
		log.Fatalf("failed to watch config file: %v", err)
	}
	
	// Start file watcher goroutine
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Printf("config file modified: %s", event.Name)
					
					// Reload config
					newCfg, err := utils.LoadConfig(path)
					if err != nil {
						log.Printf("failed to reload config: %v", err)
						continue
					}
					
					newSeries := utils.ToSeriesIDs(newCfg.Audiobooks)
					log.Printf("processing config update with %d total series", len(newSeries))
					
					// Populate database with new series (upsert operation)
					if err := populateDatabase(dbService, newSeries); err != nil {
						log.Printf("error updating database with new series: %v", err)
					} else {
						log.Printf("updated database with new series")
					}
					
					// Use incremental update to add only new entries
					app.UpdateDataIncremental(newSeries)
					log.Printf("incremental config update initiated")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("file watcher error: %v", err)
			}
		}
	}()
	
	// Populate database with series from config
	log.Printf("populating database with %d series from config...", len(series))
	if err := populateDatabase(dbService, series); err != nil {
		log.Printf("warning: failed to populate database: %v", err)
	}
	
	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Clean up any stale running jobs from previous session
	log.Printf("cleaning up stale scrape jobs...")
	if err := backgroundScraper.CleanupStaleJobs(); err != nil {
		log.Printf("warning: failed to cleanup stale jobs: %v", err)
	}
	
	// Start background scraper
	log.Printf("starting background scraper...")
	backgroundScraper.Start(ctx, settings.DefaultWorkers) // Use config setting for worker threads
	defer backgroundScraper.Stop()
	
	// Queue initial scraping jobs for all series
	if err := backgroundScraper.QueueAllSeriesUpdate(); err != nil {
		log.Printf("warning: failed to queue initial scraping jobs: %v", err)
	}
	
	// Start auto-refresh loop
	app.StartAutoRefresh()
	
	// Start server in background
	addr := fmt.Sprintf(":%d", settings.ServerPort)
	log.Printf("starting server on %s ...", addr)
	go func() {
		log.Fatal(http.ListenAndServe(addr, nil))
	}()
	
	// No longer need cache warmup - data comes from database instantly!
	log.Printf("database migration complete - server ready!")
	
	// Wait for interrupt signal for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	
	log.Printf("shutting down gracefully...")
	app.StopAutoRefresh() // Stop auto-refresh
	cancel() // This will stop the background scraper
}

// populateDatabase ensures all series from config are in the database
func populateDatabase(dbService *database.Service, series []models.SeriesIDs) error {
	for _, s := range series {
		_, err := dbService.UpsertSeries(s.Title, s.AudibleID, s.AudibleURL, s.AmazonASIN)
		if err != nil {
			log.Printf("error upserting series %s: %v", s.Title, err)
			continue
		}
		log.Printf("ensured series %s exists in database", s.Title)
	}
	return nil
}