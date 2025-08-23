package database

import (
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaFS embed.FS

// DB wraps the database connection with our business logic
type DB struct {
	*sql.DB
}

// New creates a new database connection and applies migrations
func New(dataDir string) (*DB, error) {
	// Ensure data directory exists
	dbPath := filepath.Join(dataDir, "syllabus.db")
	
	// Open SQLite database with optimized settings for concurrent access
	connectionString := dbPath + "?_foreign_keys=on&_journal_mode=WAL&_synchronous=NORMAL&_cache_size=1000&_temp_store=memory&_busy_timeout=5000"
	sqlDB, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for concurrent access
	// WAL mode supports multiple concurrent readers + single writer
	sqlDB.SetMaxOpenConns(10)  // Allow multiple concurrent connections
	sqlDB.SetMaxIdleConns(5)   // Keep some connections idle for reuse
	sqlDB.SetConnMaxLifetime(time.Hour)

	db := &DB{sqlDB}

	// Apply schema
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// migrate applies the database schema
func (db *DB) migrate() error {
	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema: %w", err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		return fmt.Errorf("failed to apply schema: %w", err)
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// Health checks database connectivity
func (db *DB) Health() error {
	return db.Ping()
}