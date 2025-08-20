package scrapers

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/michaeldvinci/syllabus/internal/models"
)

// AudibleScraperProvider implements the Audible web scraper
type AudibleScraperProvider struct {
	Enabled bool
	Client  *http.Client
}

// Fetch retrieves series information by scraping Audible pages
func (p *AudibleScraperProvider) Fetch(e models.SeriesIDs) (models.SeriesInfo, error) {
	out := models.SeriesInfo{Title: e.Title, AudibleID: e.AudibleID}
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