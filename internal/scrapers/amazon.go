package scrapers

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/michaeldvinci/syllabus/internal/models"
	"github.com/michaeldvinci/syllabus/internal/utils"
)

// AmazonPAAPIProvider implements the Amazon Product Advertising API provider
type AmazonPAAPIProvider struct {
	Enabled bool
}

// Fetch retrieves series information using Amazon PA-API
func (p *AmazonPAAPIProvider) Fetch(e models.SeriesIDs) (models.SeriesInfo, error) {
	if !p.Enabled {
		return models.SeriesInfo{Title: e.Title}, nil
	}
	return models.SeriesInfo{Title: e.Title}, nil
}

// AmazonScraperProvider implements the Amazon web scraper
type AmazonScraperProvider struct {
	Enabled bool
	Client  *http.Client
}

// Fetch retrieves series information by scraping Amazon pages
func (p *AmazonScraperProvider) Fetch(e models.SeriesIDs) (models.SeriesInfo, error) {
	out := models.SeriesInfo{Title: e.Title}
	if !p.Enabled {
		return out, nil
	}

	amzURL := utils.ExtractURLFromMarkdownLink(e.Original.Amazon)
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
	hasPreorder := reNext.FindStringSubmatch(html) != nil
	if hasPreorder {
		if m := reNext.FindStringSubmatch(html); len(m) == 2 {
			txt := strings.TrimSpace(m[1])
			if dt, err := time.Parse("January 2, 2006", txt); err == nil {
				out.AmazonNextDate = &dt
			}
		}
	}

	// Find itemBookTitle elements to get book URLs
	reBookTitle := regexp.MustCompile(`(?is)<a[^>]+id=["']itemBookTitle_(\d+)["'][^>]+href=["']([^"'&]+)`)
	bookMatches := reBookTitle.FindAllStringSubmatch(html, -1)
	
	if len(bookMatches) > 0 {
		var targetURL string
		if hasPreorder && len(bookMatches) >= 2 {
			// Use second-to-last book when preorder exists
			targetURL = bookMatches[len(bookMatches)-2][2]
		} else {
			// Use last book when no preorder exists
			targetURL = bookMatches[len(bookMatches)-1][2]
		}
		
		// Navigate to individual book page and extract publication date
		if date := p.extractPublicationDate(targetURL); date != nil {
			out.AmazonLatestDate = date
		}
	}

	return out, nil
}

func (p *AmazonScraperProvider) extractPublicationDate(bookURL string) *time.Time {
	// Handle relative URLs by prepending Amazon domain
	if strings.HasPrefix(bookURL, "/") {
		bookURL = "https://www.amazon.com" + bookURL
	}
	
	req, err := http.NewRequest("GET", bookURL, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	
	resp, err := p.Client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	
	html := string(body)
	
	// Look for the publication date in the specific Amazon structure
	rePubDate := regexp.MustCompile(`(?is)<span[^>]*class=["'][^"']*rpi-icon[^"']*book_details-publication_date[^"']*["'][^>]*>.*?</span>.*?<div[^>]*class=["'][^"']*rpi-attribute-value[^"']*["'][^>]*>\s*<span[^>]*>\s*([^<]+?)\s*</span>`)
	if m := rePubDate.FindStringSubmatch(html); len(m) == 2 {
		dateStr := strings.TrimSpace(m[1])
		
		// Try different date formats commonly used by Amazon
		formats := []string{
			"January 2, 2006",
			"Jan 2, 2006", 
			"2006-01-02",
			"1/2/2006",
			"01/02/2006",
		}
		
		for _, format := range formats {
			if dt, err := time.Parse(format, dateStr); err == nil {
				return &dt
			}
		}
	}
	
	return nil
}