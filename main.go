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
		AudibleURL    string
		AmazonURL     string
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
			AudibleLatest: AudibleLatest,
			AudibleNext:   AudibleNext,
			AmazonCount:   info.AmazonCount,
			AmazonLatest:  joinTitleDate(info.AmazonLatestTitle, info.AmazonLatestDate),
			AmazonNext:    joinTitleDate(info.AmazonNextTitle, info.AmazonNextDate), // date-only (no title set)
			AudibleURL:    audURL,
			AmazonURL:     amzURL,
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
:root {
  --bg: #ffffff;
  --text: #111827;
  --muted: #6b7280;
  --line: #e5e7eb;
  --head-bg: #f9fafb;
  --head-shadow: 0 1px 0 rgba(0,0,0,.04);
  --row-hover: #f3f4f6;
  --aud: #0ea5e9; /* cyan-ish */
  --amz: #f59e0b; /* amber-ish */
}
body { font-family: ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial; margin: 2rem; color: var(--text); background: var(--bg); }
table { border-collapse: separate; border-spacing: 0; width: 100%; background: #fff; border: 1px solid var(--line); border-radius: .5rem; overflow: hidden; }
thead th { position: sticky; top: 0; background: var(--head-bg); z-index: 3; }
thead tr:nth-child(2) th { top: 2.5rem; z-index: 2; }
thead th { border-bottom: 1px solid var(--line); padding: .6rem .75rem; text-align: left; font-weight: 600; }
thead tr:first-child th { box-shadow: var(--head-shadow); }
th, td { border-bottom: 1px solid var(--line); padding: .5rem .75rem; text-align: left; vertical-align: middle; }
tbody tr:nth-child(even) { background: #fcfcfd; }
tbody tr:hover { background: var(--row-hover); }
small { color: var(--muted); }
/* Sticky first column for easier scanning */
th:first-child, td:first-child { position: sticky; left: 0; background: inherit; z-index: 1; }
thead th:first-child { z-index: 4; }
/* Series cell with inline source pills */
.series-cell { display: flex; align-items: center; gap: .5rem; }
.series-title { font-weight: 600; }
.links { display: inline-flex; gap: .35rem; }
.pill { display: inline-flex; align-items: center; justify-content: center; font-size: .72rem; line-height: 1; padding: .28rem .45rem; border-radius: 999px; text-decoration: none; border: 1px solid rgba(0,0,0,.06); }
.pill-aud { background: rgba(14,165,233,.08); color: var(--aud); }
.pill-amz { background: rgba(245,158,11,.10); color: var(--amz); }
.badge { display: inline-block; min-width: 1.5em; padding: .15rem .5rem; border-radius: .5rem; background: #eef2ff; font-weight: 600; text-align: center; }
.count-aud { background: rgba(14,165,233,.12); color: #0369a1; }
.count-amz { background: rgba(245,158,11,.16); color: #92400e; }
.date { white-space: nowrap; color: var(--text); }
.sortable { cursor: pointer; user-select: none; }
.sortable::after { content: '\25B4\25BE'; font-size: .7em; opacity: .35; margin-left: .35rem; }
th.sort-asc::after { content: '\25B4'; opacity: .8; }
th.sort-desc::after { content: '\25BE'; opacity: .8; }
/* Top bar and settings panel */
.topbar { display: flex; align-items: center; justify-content: space-between; gap: 1rem; margin-bottom: .75rem; }
.settings-btn { display: inline-flex; align-items: center; justify-content: center; width: 2.25rem; height: 2.25rem; border-radius: .5rem; border: 1px solid var(--line); background: #fff; box-shadow: var(--head-shadow); cursor: pointer; font-size: 1.05rem; }
.settings-btn:hover { background: #f8fafc; }
.settings-btn:focus { outline: 2px solid #93c5fd; outline-offset: 2px; }
.settings-wrap { position: relative; }
.settings-panel { position: absolute; right: 0; top: 2.8rem; width: 320px; max-width: calc(100vw - 2rem); background: #fff; border: 1px solid var(--line); border-radius: .5rem; box-shadow: 0 10px 20px rgba(0,0,0,.08), 0 2px 6px rgba(0,0,0,.06); padding: .75rem; z-index: 10; }
.settings-panel .panel-section { padding: .5rem .25rem; }
.settings-panel .panel-heading { font-weight: 700; font-size: .85rem; color: var(--muted); margin-bottom: .25rem; text-transform: uppercase; letter-spacing: .02em; }
.settings-panel code { background: #f3f4f6; padding: .15rem .35rem; border-radius: .35rem; }
</style>
</head>
<body>
  <div class="topbar">
    <h1>Syllabus</h1>
    <div class="settings-wrap">
      <button class="settings-btn" id="settingsBtn" aria-expanded="false" aria-controls="settingsPanel" title="Settings" aria-label="Settings">⚙️</button>
      <div class="settings-panel" id="settingsPanel" hidden>
        <div class="panel-section">
          <div class="panel-heading">Generated at</div>
          <div class="panel-content"><code>{{ .Now }}</code></div>
        </div>
      </div>
    </div>
  </div>
  <table>
    <thead>
      <tr>
        <th rowspan="2" scope="col" class="sortable" data-col="0" data-type="text">Series</th>
        <th colspan="3" scope="colgroup">Audible</th>
        <th colspan="3" scope="colgroup">Amazon</th>
      </tr>
      <tr>
        <th scope="col" title="Number of audiobooks in the series" class="sortable" data-col="1" data-type="number">Count</th>
        <th scope="col" class="sortable" data-col="2" data-type="date">Latest</th>
        <th scope="col" class="sortable" data-col="3" data-type="date">Next</th>
        <th scope="col" title="Number of ebooks in the series on Amazon" class="sortable" data-col="4" data-type="number">Count</th>
        <th scope="col" class="sortable" data-col="5" data-type="date">Latest</th>
        <th scope="col" class="sortable" data-col="6" data-type="date">Next</th>
      </tr>
    </thead>
    <tbody>
      {{ range .Rows }}
      <tr>
        <td>
          <div class="series-cell">
            <span class="series-title">{{ .Title }}</span>
            <span class="links">
              {{ if .AudibleURL }}<a class="pill pill-aud" href="{{ .AudibleURL }}" target="_blank" rel="noopener" aria-label="Open series on Audible">Au</a>{{ end }}
              {{ if .AmazonURL }}<a class="pill pill-amz" href="{{ .AmazonURL }}" target="_blank" rel="noopener" aria-label="Open series on Amazon">Am</a>{{ end }}
            </span>
          </div>
        </td>
        <td><span class="badge count-aud">{{ .AudibleCount }}</span></td>
        <td><span class="date">{{ .AudibleLatest }}</span></td>
        <td><span class="date">{{ .AudibleNext }}</span></td>
        <td><span class="badge count-amz">{{ .AmazonCount }}</span></td>
        <td><span class="date">{{ .AmazonLatest }}</span></td>
        <td><span class="date">{{ .AmazonNext }}</span></td>
      </tr>
      {{ end }}
    </tbody>
  </table>
  <script>
  (function(){
    const table = document.querySelector('table');
    if(!table) return;
    const tbody = table.querySelector('tbody');
    const getText = (cell) => (cell.textContent || '').trim();
    const parseNumber = (s) => {
      const m = (s.match(/[-+]?[0-9]*\.?[0-9]+/)||[])[0];
      if(m === undefined || m === '') return NaN;
      return parseFloat(m);
    };
    const parseDate = (s) => {
      // Look for YYYY-MM-DD anywhere in the string
      const m = s.match(/\b(\d{4})-(\d{2})-(\d{2})\b/);
      if(m){
        const t = Date.parse(m[0] + 'T00:00:00Z');
        return isNaN(t) ? null : t;
      }
      // Fallback: try native Date
      const t = Date.parse(s);
      return isNaN(t) ? null : t;
    };
    const comparators = {
      text: (a,b) => a.localeCompare(b, undefined, {numeric:true, sensitivity:'base'}),
      number: (a,b) => (a - b),
      date: (a,b) => (a - b)
    };
    const extractors = {
      text: (cell) => getText(cell).toLowerCase(),
      number: (cell) => parseNumber(getText(cell)),
      date: (cell) => { const t = parseDate(getText(cell)); return t===null? Number.NEGATIVE_INFINITY : t; }
    };
    const clearSortStates = () => table.querySelectorAll('th.sort-asc, th.sort-desc').forEach(th=>{ th.classList.remove('sort-asc','sort-desc'); th.removeAttribute('aria-sort'); });
    const sortBy = (col, type, direction) => {
      const rows = Array.from(tbody.querySelectorAll('tr'));
      const idx = col|0;
      const getVal = (row) => extractors[type](row.children[idx]);
      const cmp = comparators[type];
      rows.sort((r1, r2) => {
        const a = getVal(r1); const b = getVal(r2);
        const c = cmp(a,b);
        return direction === 'desc' ? -c : c;
      });
      // Re-append in sorted order
      rows.forEach(r => tbody.appendChild(r));
    };
    table.querySelectorAll('th.sortable').forEach(th => {
      th.setAttribute('role', 'button');
      th.tabIndex = 0;
      let dir = th.dataset.defaultDir || 'asc';
      th.addEventListener('click', () => {
        const col = parseInt(th.dataset.col,10);
        const type = th.dataset.type || 'text';
        clearSortStates();
        sortBy(col, type, dir);
        th.classList.add(dir==='asc'?'sort-asc':'sort-desc');
        th.setAttribute('aria-sort', dir==='asc'?'ascending':'descending');
        dir = (dir === 'asc') ? 'desc' : 'asc';
      });
      th.addEventListener('keydown', (e) => { if(e.key==='Enter' || e.key===' '){ e.preventDefault(); th.click(); }});
    });
    // Settings panel toggle
    const btn = document.getElementById('settingsBtn');
    const panel = document.getElementById('settingsPanel');
    if (btn && panel) {
      const closePanel = () => { panel.hidden = true; btn.setAttribute('aria-expanded','false'); };
      const openPanel  = () => { panel.hidden = false; btn.setAttribute('aria-expanded','true'); };
      btn.addEventListener('click', (e) => {
        e.stopPropagation();
        if (panel.hidden) openPanel(); else closePanel();
      });
      document.addEventListener('click', (e) => {
        if (panel.hidden) return;
        if (!panel.contains(e.target) && e.target !== btn) closePanel();
      }, true);
      document.addEventListener('keydown', (e) => { if (e.key === 'Escape') closePanel(); });
    }
  })();
  </script>
</body>
</html>
`
