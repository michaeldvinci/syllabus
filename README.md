# syllabus

A small Go web app that reads a YAML list of audiobook series and displays live series metadata by scraping Audible and Amazon. It exposes a simple web UI and a JSON API.

So I was maintaining this Obsidian "database" manually and got pretty tired of having to open up a ton of tabs and checking release dates periodically, resulting in Syllabus. Sample config is included to see how the screenshot below was created.

I call it `syllabus` because it's a list of things to read.

It's barebones and does just what I need it to do. 

### Desktop
<img alt="syllabus-desktop" src="res/syllabus-desktop.png" height="600" width="900">

### Mobile
<img alt="syllabus-mobile" src="res/syllabus-mobile.png" height="700" width="375">

## Features

- Parse series from a YAML file
- Fetch Audible series pages and count current audiobooks
- Extract the latest release date from Audible series pages
- Fetch Amazon pages to detect ebook series size and next release date
- In-memory caching with TTL
- Optional file watching to auto-reload YAML changes
- Minimal HTML table view and JSON API

## Functionality
 
- Sortable columns
- Audible + Amazon links to series added to rows
- Settings started to be implemented

## Requirements

- Go 1.24+
- Internet access
- A YAML file of series

## Configuration

YAML schema:

```yaml
audiobooks:
  - title: "1% Lifesteal"
    audible: "https://www.audible.com/series/1-Lifesteal-Audiobooks/B0F8QMLV9T"
    amazon: "https://www.amazon.com/dp/B0DGWCJ6JP"
  - title: "A Soldier's Life"
    audible: "https://www.audible.com/series/A-Soldiers-Life-Audiobooks/B0D34549LX"
    amazon: "https://www.amazon.com/dp/B0CW18NDBQ"
```

Only title, audible, and amazon are required for scraping.


## Run Locally

```bash
❯ which dlf
dlf () {
  docker logs -f $1
}

git clone https://github.com/michaeldvinci/syllabus.git

cd syllabus

docker buildx build -t syllabus:latest . \
  && docker compose up -\
  && dlf syllabus-syllabus-1
```

Open http://localhost:8081 for the UI. The API is available at /api/series.

## Data Sources

Audible:
- `AudibleCount` is the number of occurrences of the substring `productlistitem` in the series page HTML
- `AudibleLatest` is the last occurrence of `Release date: MM-DD-YY `parsed from the series page

Amazon:
- `AmazonCount` is parsed from the element with id `collection-size` in the form `(N book series)`
- `AmazonNext` is parsed from a span with class `a-color-success a-text-bold` containing a date like `Month D, YYYY`

## JSON API

GET /api/series

Returns an array of objects with fields:
- Title
- AudibleCount
- AudibleLatestTitle
- AudibleLatestDate
- AudibleNextTitle
- AudibleNextDate
- AmazonCount
- AmazonLatestTitle
- AmazonLatestDate
- AmazonNextTitle
- AmazonNextDate
- AudibleID
- AmazonASIN
- Err

Dates are ISO 8601 when returned from the API.

## Caching

Responses from providers are cached in memory for 6 hours.

## File Watching

If enabled in the code, the application watches the directory of the YAML file and reloads data when the file changes. The cache is cleared on reload.
