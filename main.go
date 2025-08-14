package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Audiobooks []Entry `yaml:"audiobooks"`
}

type Entry struct {
	Title    string `yaml:"title"`
	Audible  string `yaml:"audible"`
	Amazon   string `yaml:"amazon"`
	AudNum   any    `yaml:"aud_num"`
	AudNext  string `yaml:"aud_next"`
	AudLast  string `yaml:"aud_last"`
	AmznNum  any    `yaml:"amzn_num"`
	AmznNext string `yaml:"amzn_next"`
	AmznLast string `yaml:"amzn_last"`
}

type SeriesIDs struct {
	Title      string
	AudibleID  string
	AudibleURL string
	AmazonASIN string
	Original   Entry
}

type SeriesInfo struct {
	Title              string
	AudibleCount       int
	AudibleLatestTitle string
	AudibleLatestDate  *time.Time
	AudibleNextTitle   string
	AudibleNextDate    *time.Time

	AmazonCount       int
	AmazonLatestTitle string
	AmazonLatestDate  *time.Time
	AmazonNextTitle   string
	AmazonNextDate    *time.Time

	AudibleID  string
	AmazonASIN string
	Err        error
}

var mdLinkRe = regexp.MustCompile(`\[[^\]]*\]\(([^)]+)\)`)

func extractURLFromMarkdownLink(s string) string {
	m := mdLinkRe.FindStringSubmatch(s)
	if len(m) == 2 {
		return m[1]
	}
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return s
	}
	return ""
}

var audibleSeriesIDRe = regexp.MustCompile(`/([A-Z0-9]{10})(?:[/?]|$)`)

func extractAudibleSeriesID(u string) string {
	if u == "" {
		return ""
	}
	m := audibleSeriesIDRe.FindStringSubmatch(u)
	if len(m) == 2 && strings.HasPrefix(m[1], "B0") {
		return m[1]
	}
	return ""
}

var amazonASINRe = regexp.MustCompile(`/dp/([A-Z0-9]{10})(?:[/?]|$)`)

func extractAmazonASIN(u string) string {
	if u == "" {
		return ""
	}
	m := amazonASINRe.FindStringSubmatch(u)
	if len(m) == 2 {
		return m[1]
	}
	return ""
}

type Provider interface {
	Fetch(entry SeriesIDs) (SeriesInfo, error)
}

type CompositeProvider struct {
	Providers []Provider
}

func (c *CompositeProvider) Fetch(e SeriesIDs) (SeriesInfo, error) {
	out := SeriesInfo{Title: e.Title, AudibleID: e.AudibleID, AmazonASIN: e.AmazonASIN}
	var mu sync.Mutex
	var wg sync.WaitGroup
	var firstErr error

	for _, p := range c.Providers {
		wg.Add(1)
		go func(p Provider) {
			defer wg.Done()
			info, err := p.Fetch(e)
			mu.Lock()
			if err != nil && firstErr == nil {
				firstErr = err
			}
			mergeSeriesInfo(&out, &info)
			mu.Unlock()
		}(p)
	}
	wg.Wait()
	out.Err = firstErr
	return out, firstErr
}

func mergeSeriesInfo(dst, src *SeriesInfo) {
	if src.AudibleCount > dst.AudibleCount {
		dst.AudibleCount = src.AudibleCount
	}
	if src.AmazonCount > dst.AmazonCount {
		dst.AmazonCount = src.AmazonCount
	}
	if src.AudibleLatestTitle != "" || src.AudibleLatestDate != nil {
		dst.AudibleLatestTitle = src.AudibleLatestTitle
		dst.AudibleLatestDate = src.AudibleLatestDate
	}
	if src.AmazonLatestTitle != "" || src.AmazonLatestDate != nil {
		dst.AmazonLatestTitle = src.AmazonLatestTitle
		dst.AmazonLatestDate = src.AmazonLatestDate
	}
	if src.AudibleNextTitle != "" || src.AudibleNextDate != nil {
		dst.AudibleNextTitle = src.AudibleNextTitle
		dst.AudibleNextDate = src.AudibleNextDate
	}
	if src.AmazonNextTitle != "" || src.AmazonNextDate != nil {
		dst.AmazonNextTitle = src.AmazonNextTitle
		dst.AmazonNextDate = src.AmazonNextDate
	}
}

type AmazonPAAPIProvider struct{ Enabled bool }

func (p *AmazonPAAPIProvider) Fetch(e SeriesIDs) (SeriesInfo, error) {
	if !p.Enabled {
		return SeriesInfo{Title: e.Title}, nil
	}
	return SeriesInfo{Title: e.Title}, nil
}

type AmazonScraperProvider struct {
	Enabled bool
	Client  *http.Client
}

func (p *AmazonScraperProvider) Fetch(e SeriesIDs) (SeriesInfo, error) {
	out := SeriesInfo{Title: e.Title}
	if !p.Enabled {
		return out, nil
	}

	amzURL := extractURLFromMarkdownLink(e.Original.Amazon)
	if amzURL == "" && e.AmazonASIN != "" {
		amzURL = fmt.Sprintf("https://www.amazon.com/dp/%s", e.AmazonASIN)
	}
	if amzURL == "" {
		return out, nil
	}

	req, err := http.NewRequest("GET", amzURL, nil)
	if err != nil {
		out.Err = err
		return out, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := p.Client.Do(req)
	if err != nil {
		out.Err = err
		return out, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		out.Err = fmt.Errorf("amazon: non-200 status %d", resp.StatusCode)
		return out, out.Err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		out.Err = err
		return out, err
	}
	html := string(body)

	reCount := regexp.MustCompile(`(?is)id=["']collection-size["'][^>]*>\s*\(?\s*([0-9,]+)\s+book`)
	if m := reCount.FindStringSubmatch(html); len(m) == 2 {
		num := strings.ReplaceAll(m[1], ",", "")
		if n, err := strconv.Atoi(num); err == nil {
			out.AmazonCount = n
		}
	}

	reNext := regexp.MustCompile(`(?is)<span[^>]+class=["'][^"']*a-color-success[^"']*a-text-bold[^"']*["'][^>]*>\s*([^<]+?)\s*</span>`)
	if m := reNext.FindStringSubmatch(html); len(m) == 2 {
		txt := strings.TrimSpace(m[1])
		if dt, err := time.Parse("January 2, 2006", txt); err == nil {
			out.AmazonNextDate = &dt
		}
	}

	return out, nil
}

type AudibleScraperProvider struct {
	Enabled bool
	Client  *http.Client
}

func (p *AudibleScraperProvider) Fetch(e SeriesIDs) (SeriesInfo, error) {
	out := SeriesInfo{Title: e.Title, AudibleID: e.AudibleID}
	if !p.Enabled {
		return out, nil
	}

	seriesURL := e.AudibleURL
	if seriesURL == "" && e.AudibleID != "" {
		seriesURL = fmt.Sprintf("https://www.audible.com/series/%s", e.AudibleID)
	}
	if seriesURL == "" {
		return out, nil
	}

	req, err := http.NewRequest("GET", seriesURL, nil)
	if err != nil {
		out.Err = err
		return out, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (SeriesTracker/1.0; +local)")

	resp, err := p.Client.Do(req)
	if err != nil {
		out.Err = err
		return out, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		out.Err = fmt.Errorf("audible: non-200 status %d", resp.StatusCode)
		return out, out.Err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		out.Err = err
		return out, err
	}
	html := string(b)
	lower := strings.ToLower(html)

	out.AudibleCount = strings.Count(lower, "productlistitem")

	re := regexp.MustCompile(`(?i)release\s*date:\s*([0-9]{2}-[0-9]{2}-[0-9]{2})`)
	matches := re.FindAllStringSubmatch(html, -1)
	if n := len(matches); n > 0 {
		last := strings.TrimSpace(matches[n-1][1]) // e.g., "06-18-25"
		if dt, err := time.Parse("01-02-06", last); err == nil {
			out.AudibleLatestDate = &dt
		}
	}

	return out, nil
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
	ttl   time.Duration
}
type cacheItem struct {
	val       SeriesInfo
	expiresAt time.Time
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{items: make(map[string]cacheItem), ttl: ttl}
}
func (c *Cache) Get(key string) (SeriesInfo, bool) {
	c.mu.RLock()
	it, ok := c.items[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(it.expiresAt) {
		return SeriesInfo{}, false
	}
	return it.val, true
}
func (c *Cache) Set(key string, v SeriesInfo) {
	c.mu.Lock()
	c.items[key] = cacheItem{val: v, expiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

type App struct {
	provider *CompositeProvider
	cache    *Cache
	data     []SeriesIDs
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <path-to-yaml>", os.Args[0])
	}
	path := os.Args[1]
	cfg, err := loadConfig(path)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	series := toSeriesIDs(cfg.Audiobooks)

	app := &App{
		provider: &CompositeProvider{
			Providers: []Provider{
				&AmazonPAAPIProvider{Enabled: false},
				&AmazonScraperProvider{Enabled: true, Client: &http.Client{Timeout: 12 * time.Second}},
				&AudibleScraperProvider{Enabled: true, Client: &http.Client{Timeout: 12 * time.Second}},
			},
		},
		cache: NewCache(6 * time.Hour),
		data:  series,
	}

	http.HandleFunc("/", app.handleIndex)
	http.HandleFunc("/api/series", app.handleAPI)
	addr := ":8080"
	log.Printf("listening on %s …", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func loadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func toSeriesIDs(entries []Entry) []SeriesIDs {
	out := make([]SeriesIDs, 0, len(entries))
	for _, e := range entries {
		audURL := extractURLFromMarkdownLink(e.Audible)
		amzURL := extractURLFromMarkdownLink(e.Amazon)
		out = append(out, SeriesIDs{
			Title:      e.Title,
			AudibleID:  extractAudibleSeriesID(audURL),
			AudibleURL: audURL,
			AmazonASIN: extractAmazonASIN(amzURL),
			Original:   e,
		})
	}
	return out
}

func (a *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	type Row struct {
		Title         string
		AudibleCount  int
		AudibleLatest string
		AudibleNext   string
		AmazonCount   int
		AmazonLatest  string
		AmazonNext    string
	}
	type Page struct {
		Rows []Row
		Now  string
	}
	var rows []Row

	AudibleLatest := ""
	AudibleNext := ""
	loc, _ := time.LoadLocation("America/Chicago")
	now := time.Now().In(loc)
	y, m, d := now.Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, loc)
	infos := a.collectAll()
	for _, info := range infos {
		if info.AudibleLatestDate != nil {
			AudibleLatestS := info.AudibleLatestDate.Format("2006-01-02")

			other, err := time.ParseInLocation("2006-01-02", AudibleLatestS, loc)
			if err != nil {
				panic(err)
			}

			switch {
			case other.Before(today):
				AudibleLatest = AudibleLatestS
				AudibleNext = ""
			case other.After(today):
				AudibleLatest = ""
				AudibleNext = AudibleLatestS
			default:
				fmt.Println(AudibleLatestS, "is TODAY")
			}
		}
		rows = append(rows, Row{
			Title:         info.Title,
			AudibleCount:  info.AudibleCount,
			AudibleLatest: AudibleLatest,
			AudibleNext:   AudibleNext,
			AmazonCount:   info.AmazonCount,
			AmazonLatest:  joinTitleDate(info.AmazonLatestTitle, info.AmazonLatestDate),
			AmazonNext:    joinTitleDate(info.AmazonNextTitle, info.AmazonNextDate), // date-only (no title set)
		})
	}

	tpl := template.Must(template.New("idx").Parse(indexHTML))
	if err := tpl.Execute(w, Page{Rows: rows, Now: time.Now().Format(time.RFC822)}); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (a *App) handleAPI(w http.ResponseWriter, r *http.Request) {
	infos := a.collectAll()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(infos)
}

func (a *App) collectAll() []SeriesInfo {
	var wg sync.WaitGroup
	infos := make([]SeriesInfo, len(a.data))
	for i, e := range a.data {
		wg.Add(1)
		go func(i int, e SeriesIDs) {
			defer wg.Done()
			key := e.Title + "|" + e.AudibleID + "|" + e.AmazonASIN
			if v, ok := a.cache.Get(key); ok {
				infos[i] = v
				return
			}
			info, err := a.provider.Fetch(e)
			if err != nil {
				info.Err = err
			}
			a.cache.Set(key, info)
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

const indexHTML = `
<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>Audiobook / Ebook Series Tracker</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
body { font-family: ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial; margin: 2rem; }
table { border-collapse: collapse; width: 100%; }
th, td { border-bottom: 1px solid #ddd; padding: .5rem; text-align: left; }
th { position: sticky; top: 0; background: #fff; }
small { color: #666; }
</style>
</head>
<body>
  <h1>Syllabus</h1>
  <small>Generated at {{ .Now }}</small>
  <table>
    <thead>
      <tr>
        <th>Series</th>
        <th>Audible Count</th>
        <th>Audible Latest</th>
        <th>Audible Next</th>
        <th>Amazon Count</th>
        <th>Amazon Latest</th>
        <th>Amazon Next</th>
      </tr>
    </thead>
    <tbody>
      {{ range .Rows }}
      <tr>
        <td>{{ .Title }}</td>
        <td>{{ .AudibleCount }}</td>
        <td>{{ .AudibleLatest }}</td>
        <td>{{ .AudibleNext }}</td>
        <td>{{ .AmazonCount }}</td>
        <td>{{ .AmazonLatest }}</td>
        <td>{{ .AmazonNext }}</td>
      </tr>
      {{ end }}
    </tbody>
  </table>
</body>
</html>
`
