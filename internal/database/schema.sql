-- SQLite schema for Syllabus application

-- Series table - stores basic series information
CREATE TABLE IF NOT EXISTS series (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    audible_id TEXT,
    audible_url TEXT,
    amazon_asin TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Books table - stores individual book information
CREATE TABLE IF NOT EXISTS books (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    series_id INTEGER NOT NULL,
    provider TEXT NOT NULL CHECK (provider IN ('audible', 'amazon')),
    title TEXT NOT NULL,
    book_number INTEGER,
    release_date DATE,
    is_preorder BOOLEAN DEFAULT 0,
    is_latest BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE,
    UNIQUE(series_id, provider, book_number)
);

-- Scrape jobs table - tracks scraping operations
CREATE TABLE IF NOT EXISTS scrape_jobs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    series_id INTEGER NOT NULL,
    provider TEXT NOT NULL CHECK (provider IN ('audible', 'amazon')),
    status TEXT NOT NULL CHECK (status IN ('pending', 'running', 'completed', 'failed')) DEFAULT 'pending',
    started_at DATETIME,
    completed_at DATETIME,
    error_message TEXT,
    book_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE
);

-- Series stats view - aggregated data for quick queries
CREATE VIEW IF NOT EXISTS series_stats AS
SELECT 
    s.id,
    s.title,
    s.audible_id,
    s.amazon_asin,
    s.updated_at,
    
    -- Audible stats
    COUNT(CASE WHEN b.provider = 'audible' THEN 1 END) as audible_count,
    MAX(CASE WHEN b.provider = 'audible' AND b.is_latest = 1 THEN b.title END) as audible_latest_title,
    MAX(CASE WHEN b.provider = 'audible' AND b.is_latest = 1 THEN b.release_date END) as audible_latest_date,
    MAX(CASE WHEN b.provider = 'audible' AND b.is_preorder = 1 THEN b.title END) as audible_next_title,
    MAX(CASE WHEN b.provider = 'audible' AND b.is_preorder = 1 THEN b.release_date END) as audible_next_date,
    
    -- Amazon stats  
    COUNT(CASE WHEN b.provider = 'amazon' THEN 1 END) as amazon_count,
    MAX(CASE WHEN b.provider = 'amazon' AND b.is_latest = 1 THEN b.title END) as amazon_latest_title,
    MAX(CASE WHEN b.provider = 'amazon' AND b.is_latest = 1 THEN b.release_date END) as amazon_latest_date,
    MAX(CASE WHEN b.provider = 'amazon' AND b.is_preorder = 1 THEN b.title END) as amazon_next_title,
    MAX(CASE WHEN b.provider = 'amazon' AND b.is_preorder = 1 THEN b.release_date END) as amazon_next_date
    
FROM series s
LEFT JOIN books b ON s.id = b.series_id
GROUP BY s.id, s.title, s.audible_id, s.amazon_asin, s.updated_at;

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_books_series_provider ON books(series_id, provider);
CREATE INDEX IF NOT EXISTS idx_books_release_date ON books(release_date);
CREATE INDEX IF NOT EXISTS idx_scrape_jobs_status ON scrape_jobs(status);
CREATE INDEX IF NOT EXISTS idx_series_title ON series(title);

-- Trigger to update series.updated_at when books are modified
CREATE TRIGGER IF NOT EXISTS update_series_timestamp 
AFTER INSERT ON books
BEGIN
    UPDATE series SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.series_id;
END;

CREATE TRIGGER IF NOT EXISTS update_series_timestamp_update
AFTER UPDATE ON books  
BEGIN
    UPDATE series SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.series_id;
END;