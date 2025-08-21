package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/michaeldvinci/syllabus/internal/auth"
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

	// Initialize authentication store with file persistence
	dataDir := "./data"
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

	// Initialize application
	app := &handlers.App{
		Provider:    provider,
		Cache:       cache.NewCache(6 * time.Hour),
		Data:        series,
		RefreshChan: make(chan bool, 1),
	}

	// Setup authentication routes (no middleware needed)
	http.HandleFunc("/login", authHandlers.HandleLogin)
	http.HandleFunc("/logout", authHandlers.HandleLogout)
	http.HandleFunc("/api/auth", authMiddleware.OptionalAuth(authHandlers.HandleAPI))
	
	// Setup admin-only routes
	http.HandleFunc("/api/users", authMiddleware.RequireAdmin(authHandlers.HandleListUsers))
	http.HandleFunc("/api/users/create", authMiddleware.RequireAdmin(authHandlers.HandleCreateUser))
	http.HandleFunc("/api/users/reset-password", authMiddleware.RequireAdmin(authHandlers.HandleResetPassword))

	// Setup protected HTTP routes with authentication middleware
	http.HandleFunc("/", authMiddleware.RequireAuth(app.HandleIndex))
	http.HandleFunc("/api/series", authMiddleware.RequireAuth(app.HandleAPI))
	http.HandleFunc("/events", authMiddleware.RequireAuth(app.HandleEvents))
	http.HandleFunc("/calendar.ics", authMiddleware.RequireAuth(app.HandleICal))
	
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
	
	// Start server in background
	addr := ":8080"
	log.Printf("starting server on %s ...", addr)
	go func() {
		log.Fatal(http.ListenAndServe(addr, nil))
	}()
	
	// Warm up cache by fetching all data immediately on startup
	log.Printf("warming up cache for %d series...", len(series))
	start := time.Now()
	app.WarmupCache()
	log.Printf("cache warmup completed in %v", time.Since(start))
	
	// Keep server running
	select {}
}