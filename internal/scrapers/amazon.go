package scrapers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
		log.Printf("Amazon scraper: disabled for %s", e.Title)
		return out, nil
	}

	amzURL := utils.ExtractURLFromMarkdownLink(e.Original.Amazon)
	if amzURL == "" && e.AmazonASIN != "" {
		amzURL = fmt.Sprintf("https://www.amazon.com/dp/%s", e.AmazonASIN)
	}
	if amzURL == "" {
		log.Printf("Amazon scraper: no URL found for %s", e.Title)
		return out, nil
	}

	log.Printf("Amazon scraper: fetching %s from %s", e.Title, amzURL)
	
	// Add delay to avoid rate limiting (random between 1-3 seconds)
	time.Sleep(time.Duration(800+time.Now().UnixNano()%1000) * time.Millisecond)

	req, err := http.NewRequest("GET", amzURL, nil)
	if err != nil {
		log.Printf("Amazon scraper: failed to create request for %s: %v", e.Title, err)
		out.Err = err
		return out, err
	}
	// Set comprehensive headers to mimic a real browser and avoid bot detection
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Cache-Control", "max-age=0")
	
	// Add some randomness to avoid looking like a bot
	req.Header.Set("Priority", "u=0, i")

	resp, err := p.Client.Do(req)
	if err != nil {
		log.Printf("Amazon scraper: HTTP request failed for %s: %v", e.Title, err)
		out.Err = err
		return out, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Amazon scraper: HTTP %d for %s at %s", resp.StatusCode, e.Title, amzURL)
		out.Err = fmt.Errorf("amazon: non-200 status %d", resp.StatusCode)
		return out, out.Err
	}

	log.Printf("Amazon scraper: successfully fetched %s (status %d)", e.Title, resp.StatusCode)

	// Handle gzip compression
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			log.Printf("Amazon scraper: failed to create gzip reader for %s: %v", e.Title, err)
			out.Err = err
			return out, err
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("Amazon scraper: failed to read response body for %s: %v", e.Title, err)
		out.Err = err
		return out, err
	}
	html := string(body)

	log.Printf("Amazon scraper: parsing HTML for %s (length: %d bytes)", e.Title, len(html))
	
	
	// Check for signs of JavaScript-rendered content
	hasJavaScriptIndicators := strings.Contains(html, "window.P") || 
		strings.Contains(html, "ue_widget") || 
		strings.Contains(html, "data-client-recs-list") ||
		strings.Contains(html, "Loading...") ||
		strings.Contains(html, "window.uet")
	
	hasActualContent := strings.Contains(html, "collection-size") || 
		strings.Contains(html, "itemBookTitle") ||
		strings.Contains(html, "a-color-success")
		
	log.Printf("Amazon scraper: JS indicators: %v, actual content: %v", hasJavaScriptIndicators, hasActualContent)
	
	// Try to extract data from JSON-LD structured data (fallback for JS-heavy pages)
	if !hasActualContent && hasJavaScriptIndicators {
		log.Printf("Amazon scraper: attempting JSON-LD extraction for %s", e.Title)
		if jsonData := extractJSONLD(html); jsonData != nil {
			log.Printf("Amazon scraper: found JSON-LD data for %s", e.Title)
			// Try to parse structured data for book/series info
			if count, latest, next := parseStructuredData(jsonData, e.Title); count > 0 || latest != nil || next != nil {
				if count > 0 {
					out.AmazonCount = count
					log.Printf("Amazon scraper: extracted count %d from structured data for %s", count, e.Title)
				}
				if latest != nil {
					out.AmazonLatestDate = latest
					log.Printf("Amazon scraper: extracted latest date from structured data for %s", e.Title)
				}
				if next != nil {
					out.AmazonNextDate = next
					log.Printf("Amazon scraper: extracted next date from structured data for %s", e.Title)
				}
				log.Printf("Amazon scraper: completed %s via structured data - Count: %d, Latest: %v, Next: %v", 
					e.Title, out.AmazonCount, out.AmazonLatestDate, out.AmazonNextDate)
				return out, nil
			}
		}
	}

	// Check if this is a single book page vs series collection page
	isSingleBook := strings.Contains(html, `"@type":"Book"`) || strings.Contains(html, `id="productTitle"`)
	isSeriesPage := strings.Contains(html, `collection-size`) || strings.Contains(html, `itemBookTitle_`)
	log.Printf("Amazon scraper: page type for %s - single book: %v, series: %v", e.Title, isSingleBook, isSeriesPage)

	// Check if we got a CAPTCHA page instead of the actual content
	if strings.Contains(html, "validateCaptcha") || strings.Contains(html, "Continue shopping") {
		log.Printf("Amazon CAPTCHA detected for %s - bot detection triggered", e.Title)
		log.Printf("Count: 0")
		log.Printf("Next: blocked")
		log.Printf("Latest: blocked")
		return out, nil
	}
	
	// Step 1: Search for book count in collection-size element
	countPattern := `\((\d+) book series\)`
	reCount := regexp.MustCompile(countPattern)
	if m := reCount.FindStringSubmatch(html); len(m) == 2 {
		if n, err := strconv.Atoi(m[1]); err == nil {
			out.AmazonCount = n
		}
	}
	
	log.Printf("Count: %d", out.AmazonCount)

	// Step 2: Search for preorder date (next) - look for spans with "a-color-success a-text-bold"
	preorderPattern := `(?is)<span\s+class=["']a-color-success\s+a-text-bold["']>\s*([^<]+?)\s*</span>`
	reNext := regexp.MustCompile(preorderPattern)
	hasPreorder := false
	if m := reNext.FindStringSubmatch(html); len(m) == 2 {
		txt := strings.TrimSpace(m[1])
		if dt, err := time.Parse("January 2, 2006", txt); err == nil {
			out.AmazonNextDate = &dt
			hasPreorder = true
			log.Printf("Next: %s", dt.Format("2006-01-02"))
		} else {
			log.Printf("Next: failed to parse '%s'", txt)
		}
	} else {
		log.Printf("Next: not found")
	}

	// Step 3: Find itemBookTitle elements to get book URLs
	bookTitlePattern := `(?is)<a\s+id=["']itemBookTitle_(\d+)["'][^>]*href=["']([^"']+)["'][^>]*>`
	reBookTitle := regexp.MustCompile(bookTitlePattern)
	bookMatches := reBookTitle.FindAllStringSubmatch(html, -1)
	
	log.Printf("Found %d book links", len(bookMatches))

	// Step 4: Determine which book to get publication date from
	var targetBookURL string
	if len(bookMatches) > 0 {
		if hasPreorder && len(bookMatches) >= 2 {
			// Use second-to-last book when preorder exists
			targetBookURL = bookMatches[len(bookMatches)-2][2]
			log.Printf("Using second-to-last book URL (preorder exists): %s", targetBookURL)
		} else {
			// Use last book when no preorder
			targetBookURL = bookMatches[len(bookMatches)-1][2]
			log.Printf("Using last book URL (no preorder): %s", targetBookURL)
		}
		
		// Extract ASIN from the URL and build clean Amazon URL
		asinPattern := `(?i)/gp/product/([A-Z0-9]{10})`
		asinRe := regexp.MustCompile(asinPattern)
		if asinMatch := asinRe.FindStringSubmatch(targetBookURL); len(asinMatch) == 2 {
			asin := asinMatch[1]
			fullURL := "https://www.amazon.com/gp/product/" + asin
			log.Printf("Cleaned URL: %s", fullURL)
			
			// Extract publication date from the book page
			if date := p.extractPublicationDate(fullURL); date != nil {
				out.AmazonLatestDate = date
				log.Printf("Latest: %s", date.Format("2006-01-02"))
			} else {
				log.Printf("Latest: not found")
			}
		}
	}
	
	// If no book links found, use fallback: extract ASIN from original URL and go directly to book page
	if len(bookMatches) == 0 {
		log.Printf("Amazon scraper: no book links found, using URL fallback for %s", e.Title)
		
		// Extract ASIN from the original URL
		asinPattern := `(?i)/dp/([A-Z0-9]{10})`
		asinRe := regexp.MustCompile(asinPattern)
		if asinMatch := asinRe.FindStringSubmatch(amzURL); len(asinMatch) == 2 {
			targetURL := "/gp/product/" + asinMatch[1]
			log.Printf("Amazon scraper: extracted ASIN %s, fetching %s", asinMatch[1], targetURL)
			
			if date := p.extractPublicationDate(targetURL); date != nil {
				out.AmazonLatestDate = date
				log.Printf("Amazon scraper: found latest publication date for %s: %s", e.Title, date.Format("2006-01-02"))
			}
			
			// Don't override count here - let it stay 0 if not found via normal means
			log.Printf("Amazon scraper: completed %s via URL fallback - Count: %d, Latest: %v, Next: %v", 
				e.Title, out.AmazonCount, out.AmazonLatestDate, out.AmazonNextDate)
			return out, nil
		}
	}
	
	if len(bookMatches) > 0 {
		var targetURL string
		if hasPreorder && len(bookMatches) >= 2 {
			// Use second-to-last book when preorder exists
			targetURL = bookMatches[len(bookMatches)-2][2]
			log.Printf("Amazon scraper: using second-to-last book URL for %s (preorder exists)", e.Title)
		} else {
			// Use last book when no preorder exists
			targetURL = bookMatches[len(bookMatches)-1][2]
			log.Printf("Amazon scraper: using last book URL for %s (no preorder)", e.Title)
		}
		
		log.Printf("Amazon scraper: extracting publication date from %s", targetURL)
		// Navigate to individual book page and extract publication date
		if date := p.extractPublicationDate(targetURL); date != nil {
			out.AmazonLatestDate = date
			log.Printf("Amazon scraper: latest release date for %s: %s", e.Title, date.Format("2006-01-02"))
		} else {
			log.Printf("Amazon scraper: no publication date found for %s", e.Title)
		}
	} else {
		log.Printf("Amazon scraper: no book links found for %s", e.Title)
	}

	log.Printf("Amazon scraper: completed %s - Count: %d, Latest: %v, Next: %v", 
		e.Title, out.AmazonCount, out.AmazonLatestDate, out.AmazonNextDate)
	return out, nil
}

func (p *AmazonScraperProvider) extractPublicationDate(bookURL string) *time.Time {
	// Handle relative URLs by prepending Amazon domain
	if strings.HasPrefix(bookURL, "/") {
		bookURL = "https://www.amazon.com" + bookURL
	}
	
	log.Printf("Amazon scraper: fetching publication date from %s", bookURL)
	
	// Add delay to avoid rate limiting
	time.Sleep(time.Duration(800+time.Now().UnixNano()%1000) * time.Millisecond)
	
	req, err := http.NewRequest("GET", bookURL, nil)
	if err != nil {
		log.Printf("Amazon scraper: failed to create request for book page: %v", err)
		return nil
	}
	// Use same headers as main function to avoid bot detection
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Priority", "u=0, i")
	
	resp, err := p.Client.Do(req)
	if err != nil {
		log.Printf("Amazon scraper: HTTP request failed for book page: %v", err)
		return nil
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		log.Printf("Amazon scraper: HTTP %d for book page %s", resp.StatusCode, bookURL)
		return nil
	}
	
	// Handle gzip compression
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			log.Printf("Amazon scraper: failed to create gzip reader for book page: %v", err)
			return nil
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("Amazon scraper: failed to read book page response: %v", err)
		return nil
	}
	
	html := string(body)
	log.Printf("Amazon scraper: parsing book page HTML (length: %d bytes)", len(html))
	
	
	// Find all instances of the specific div class and look for one containing a date
	divPattern := `<div\s+class="a-section\s+a-spacing-none\s+a-text-center\s+rpi-attribute-value">\s*<span>([^<]+)</span>`
	re := regexp.MustCompile(divPattern)
	matches := re.FindAllStringSubmatch(html, -1)
	
	for _, match := range matches {
		if len(match) == 2 {
			dateStr := strings.TrimSpace(match[1])
			
			// Try to parse as a date - Amazon uses "Month DD, YYYY" format
			if dt, err := time.Parse("January 2, 2006", dateStr); err == nil {
				return &dt
			}
		}
	}
	
	return nil
}

// extractJSONLD extracts JSON-LD structured data from HTML
func extractJSONLD(html string) map[string]interface{} {
	// Look for JSON-LD script tags
	re := regexp.MustCompile(`(?is)<script[^>]*type=["']application/ld\+json["'][^>]*>(.*?)</script>`)
	matches := re.FindAllStringSubmatch(html, -1)
	
	for _, match := range matches {
		if len(match) == 2 {
			var jsonData map[string]interface{}
			if err := json.Unmarshal([]byte(match[1]), &jsonData); err == nil {
				return jsonData
			}
		}
	}
	
	// Also try to find embedded JSON data in script tags
	re2 := regexp.MustCompile(`(?is)window\.P\.register\(['"]initial-data['"],\s*({.*?})\);`)
	if match := re2.FindStringSubmatch(html); len(match) == 2 {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(match[1]), &jsonData); err == nil {
			return jsonData
		}
	}
	
	return nil
}

// parseStructuredData attempts to extract book/series data from structured JSON
func parseStructuredData(data map[string]interface{}, title string) (count int, latest *time.Time, next *time.Time) {
	log.Printf("Amazon scraper: parsing structured data for %s", title)
	
	// Look for book series information in the JSON structure
	if bookType, ok := data["@type"].(string); ok && bookType == "Book" {
		log.Printf("Amazon scraper: found Book type in structured data")
		
		// Try to find publication date
		if datePublished, ok := data["datePublished"].(string); ok {
			if dt, err := time.Parse("2006-01-02", datePublished); err == nil {
				latest = &dt
				log.Printf("Amazon scraper: found publication date in structured data: %s", datePublished)
			}
		}
		
		// For single books, count is 1
		count = 1
	}
	
	// Look for series information
	if series, ok := data["isPartOf"].(map[string]interface{}); ok {
		if seriesType, ok := series["@type"].(string); ok && seriesType == "BookSeries" {
			log.Printf("Amazon scraper: found BookSeries in structured data")
			
			// Try to extract series count if available
			if numBooks, ok := series["numberOfItems"].(float64); ok {
				count = int(numBooks)
				log.Printf("Amazon scraper: found series count in structured data: %d", count)
			}
		}
	}
	
	return count, latest, next
}