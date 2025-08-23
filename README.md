# syllabus

A Go web application that tracks audiobook series release dates by scraping Audible and Amazon. Features user authentication, database persistence, background scraping, and automatic refresh capabilities with a clean web UI and JSON API.

Originally created to replace manual maintenance of an Obsidian "database" for tracking audiobook release dates. Sample config is included to see how the screenshots below were created.

I call it `syllabus` because it's a list of things to read.

Perfect for homelab deployment with Docker Compose for automated audiobook series tracking. 

### Desktop
<img alt="syllabus-desktop" src="res/syllabus-main.png" height="600" width="900">

### Search
<img alt="syllabus-desktop" src="res/syllabus-search.png" height="600" width="900">

### Settings
<img alt="syllabus-desktop" src="res/syllabus-settings.png" height="600" width="900">

## Features

### Core Functionality
- **YAML Configuration**: Parse audiobook series from a simple YAML file
- **Multi-Provider Scraping**: Fetch data from both Audible and Amazon
- **Release Date Tracking**: Extract latest and next release dates automatically
- **Database Persistence**: SQLite database for reliable data storage
- **Background Processing**: Multi-threaded background scraper with job queue
- **Real-time Updates**: Server-sent events for live UI updates

### User Experience
- **Authentication System**: Secure login with role-based access (Admin/User)
- **Responsive Web UI**: Clean, mobile-friendly interface
- **Auto-refresh**: Configurable automatic data refresh (2-10 hours)
- **Manual Refresh**: On-demand data refresh with progress tracking
- **iCal Export**: Subscribe to release date calendar in your favorite app
- **Settings Panel**: Manage refresh intervals and preferences

### Technical Features
- **File Watching**: Auto-reload when YAML configuration changes
- **Graceful Shutdown**: Clean application termination handling
- **Docker Support**: Full containerization with Docker Compose
- **JSON API**: Programmatic access to series data
- **Rate Limiting**: Intelligent scraping to avoid provider restrictions

## Quick Start

### Default Credentials
- **Username**: `admin`  
- **Password**: `admin`

⚠️ **Change the default password immediately after first login**

## Requirements

- Go 1.21+ (for local development)
- Docker & Docker Compose (recommended)
- Internet access for scraping
- YAML configuration file

## Configuration

YAML schema with optional application settings:

```yaml
# Application Settings (optional - defaults shown)
settings:
  auto_refresh_interval: 6  # Hours between automatic data refreshes (default: 6)
  default_workers: 4        # Number of concurrent scraper workers (default: 4)  
  server_port: 8080         # Port for the web server (default: 8080)
  cache_timeout: 6          # Cache timeout in hours (default: 6)
  log_level: "info"         # Logging level: debug, info, warn, error (default: info)

# Audiobook/Ebook Series Configuration
audiobooks:
  - title: "1% Lifesteal"
    audible: "https://www.audible.com/series/1-Lifesteal-Audiobooks/B0F8QMLV9T"
    amazon: "https://www.amazon.com/dp/B0DGWCJ6JP"
  - title: "A Soldier's Life"
    audible: "https://www.audible.com/series/A-Soldiers-Life-Audiobooks/B0D34549LX"
    amazon: "https://www.amazon.com/dp/B0CW18NDBQ"
```

**Required fields:** Only `title`, `audible`, and `amazon` are required for scraping.

**Settings section:** All settings are optional and will use sensible defaults if not specified. The settings section allows you to customize application behavior without modifying code.

### Environment Variables (Docker Compose)

For containerized deployments, you can override any setting using environment variables. Environment variables take precedence over YAML configuration:

```yaml
environment:
  SYLLABUS_AUTO_REFRESH_INTERVAL: "4"  # Hours between auto-refreshes
  SYLLABUS_DEFAULT_WORKERS: "2"        # Number of scraper workers  
  SYLLABUS_SERVER_PORT: "8080"         # Server port
  PORT: "8080"                         # Alternative port env var (standard)
  SYLLABUS_CACHE_TIMEOUT: "6"          # Cache timeout in hours
  SYLLABUS_LOG_LEVEL: "debug"          # Log level: debug, info, warn, error
```

**Configuration Priority:**
1. **Runtime UI changes** (highest priority) - persisted to database
2. **Environment Variables** - override YAML at startup
3. **YAML Configuration** - file-based defaults  
4. **Built-in Defaults** (lowest priority)

**Database Persistence:** UI changes (like auto-refresh interval) are saved to the database and survive container restarts. The UI will always show the current server state.

## Installation & Deployment

### Docker Compose (Recommended)

```bash
git clone https://github.com/michaeldvinci/syllabus.git
cd syllabus

# Basic deployment
docker compose up -d

# View logs
docker compose logs -f
```

**Customization with Environment Variables:**

```bash
# Copy example override file
cp docker-compose.override.example.yaml docker-compose.override.yaml

# Edit settings
nano docker-compose.override.yaml

# Deploy with custom settings
docker compose up -d
```

The override file allows you to customize settings without modifying the main docker-compose.yaml file.

### Local Development

```bash
git clone https://github.com/michaeldvinci/syllabus.git
cd syllabus

# Run directly
go run cmd/syllabus/main.go config/books.yaml

# Or build first
go build -o syllabus cmd/syllabus/main.go
./syllabus config/books.yaml
```

### Access the Application

- **Web UI**: http://localhost:8080
- **Login**: Use `admin` / `admin` (change immediately!)
- **JSON API**: http://localhost:8080/api/series
- **Calendar**: http://localhost:8080/calendar.ics

## Data Sources & Scraping

### Audible Scraping
- **Series Count**: Number of `productlistitem` occurrences in series page HTML
- **Latest Release**: Most recent `Release date: MM-DD-YY` from series page
- **Next Release**: Extracted from "Coming Soon" or pre-order sections

### Amazon Scraping  
- **Series Count**: Parsed from `collection-size` element as `(N book series)`
- **Next Release**: Date from `a-color-success a-text-bold` span elements
- **Series Detection**: Automatic ASIN extraction from Amazon URLs

### Background Processing
- **Multi-threaded**: 4 concurrent workers for faster scraping
- **Job Queue**: Persistent SQLite-based job management  
- **Rate Limiting**: Intelligent delays to respect provider limits
- **Error Handling**: Automatic retry logic with exponential backoff

## API Reference

### Authentication Required
All API endpoints require authentication via session cookie (login at `/login`).

### GET /api/series
Returns an array of series objects with the following fields:

```json
{
  "Title": "Series Name",
  "AudibleCount": 5,
  "AudibleLatestTitle": "Book Title",
  "AudibleLatestDate": "2024-01-15T00:00:00Z",
  "AudibleNextTitle": "Next Book Title", 
  "AudibleNextDate": "2024-03-20T00:00:00Z",
  "AmazonCount": 5,
  "AmazonLatestTitle": "Book Title",
  "AmazonLatestDate": "2024-01-15T00:00:00Z",
  "AmazonNextTitle": "Next Book Title",
  "AmazonNextDate": "2024-03-20T00:00:00Z",
  "AudibleID": "B0EXAMPLE",
  "AmazonASIN": "B08EXAMPLE",
  "Err": null
}
```

### POST /refresh
Triggers a manual refresh of all series data.

### GET /calendar.ics
Returns iCal calendar file with all upcoming release dates.

### POST /api/auto-refresh
Updates automatic refresh interval (Admin only).
```json
{"interval": 6}
```

All dates are returned in ISO 8601 format.

## Data Storage & Persistence

### SQLite Database
- **Location**: `./data/syllabus.db` 
- **Schema**: Series, books, and job queue tables
- **Persistence**: Survives application restarts
- **Migration**: Automatic schema updates on startup

### User Management
- **Storage**: `./data/users.json`
- **Encryption**: bcrypt password hashing
- **Roles**: Admin and User access levels
- **Default**: Admin user created on first run

### Configuration Watching
- **Auto-reload**: YAML file changes trigger incremental updates
- **Smart Updates**: Only new series are scraped, existing data preserved
- **Hot Refresh**: No application restart required

## Auto-Refresh System

### Configurable Intervals
- **Options**: 2, 4, 6, 8, 10 hours
- **Default**: 6 hours
- **UI Control**: Settings panel slider
- **Persistence**: Interval survives restarts

### Refresh Behavior
- **Incremental**: Only updates stale data
- **Background**: Non-blocking operation
- **Progress**: Real-time updates via Server-Sent Events
- **Manual Override**: Refresh button forces immediate update

## Troubleshooting

### Log Locations
- **Docker**: `docker compose logs syllabus`
- **Local**: Console output
- **Scraper**: Detailed job progress in logs
