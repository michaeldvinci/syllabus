package scrapers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
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
	// Explicitly clear all output fields to prevent variable reuse between scraping loops
	out := models.SeriesInfo{
		Title:              e.Title,
		AudibleID:          e.AudibleID,
		AudibleCount:       0,
		AudibleLatestTitle: "",
		AudibleLatestDate:  nil,
		AudibleNextTitle:   "",
		AudibleNextDate:    nil,
	}
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
		return out, nil // Return empty data instead of error
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (SeriesTracker/1.0; +local)")

	resp, err := p.Client.Do(req)
	if err != nil {
		return out, nil // Return empty data instead of error
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return out, nil // Return empty data instead of error
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return out, nil // Return empty data instead of error
	}
	html := string(b)
	lower := strings.ToLower(html)

	// Count books using the most reliable pattern first
	out.AudibleCount = strings.Count(lower, "productlistitem")
	
	// Try alternative patterns only if the primary one found nothing
	if out.AudibleCount == 0 {
		count2 := strings.Count(lower, "adbl-prod-item")
		count3 := strings.Count(lower, "bc-series-item")
		
		// Use the higher of the two alternative counts, but cap at reasonable maximum
		if count2 > count3 && count2 <= 100 {
			out.AudibleCount = count2
		} else if count3 <= 100 {
			out.AudibleCount = count3
		}
	}
	
	// Also look for explicit book count indicators
	countPattern1 := regexp.MustCompile(`(?i)(\d+)\s+books?\s+in\s+(?:this\s+)?series`)
	if m := countPattern1.FindStringSubmatch(html); len(m) == 2 {
		if n, err := strconv.Atoi(m[1]); err == nil && n > out.AudibleCount {
			out.AudibleCount = n
		}
	}
	
	countPattern2 := regexp.MustCompile(`(?i)series\s+contains\s+(\d+)\s+books?`)
	if m := countPattern2.FindStringSubmatch(html); len(m) == 2 {
		if n, err := strconv.Atoi(m[1]); err == nil && n > out.AudibleCount {
			out.AudibleCount = n
		}
	}
	
	// Find ALL release dates using multiple patterns
	var allDates []time.Time
	
	// Pattern 1: release date: MM-DD-YY
	re1 := regexp.MustCompile(`(?i)release\s*date:\s*([0-9]{1,2}-[0-9]{1,2}-[0-9]{2,4})`)
	matches1 := re1.FindAllStringSubmatch(html, -1)
	for _, match := range matches1 {
		if len(match) == 2 {
			dateStr := strings.TrimSpace(match[1])
			// Try multiple date formats
			for _, layout := range []string{"01-02-06", "1-2-06", "01-02-2006", "1-2-2006"} {
				if dt, err := time.Parse(layout, dateStr); err == nil {
					allDates = append(allDates, dt)
					break
				}
			}
		}
	}
	
	// Pattern 2: published MM-DD-YYYY
	re2 := regexp.MustCompile(`(?i)published\s*:?\s*([0-9]{1,2}-[0-9]{1,2}-[0-9]{4})`)
	matches2 := re2.FindAllStringSubmatch(html, -1)
	for _, match := range matches2 {
		if len(match) == 2 {
			dateStr := strings.TrimSpace(match[1])
			if dt, err := time.Parse("1-2-2006", dateStr); err == nil {
				allDates = append(allDates, dt)
			} else if dt, err := time.Parse("01-02-2006", dateStr); err == nil {
				allDates = append(allDates, dt)
			}
		}
	}
	
	// Pattern 3: Month DD, YYYY format
	re3 := regexp.MustCompile(`(?i)(?:release|published)\s*:?\s*([A-Z][a-z]+\s+\d{1,2},\s+\d{4})`)
	matches3 := re3.FindAllStringSubmatch(html, -1)
	for _, match := range matches3 {
		if len(match) == 2 {
			dateStr := strings.TrimSpace(match[1])
			if dt, err := time.Parse("January 2, 2006", dateStr); err == nil {
				allDates = append(allDates, dt)
			}
		}
	}
	
	// Log all dates found for this series
	var dateStrings []string
	for _, date := range allDates {
		dateStrings = append(dateStrings, date.Format("2006-01-02"))
	}
	log.Printf("Audible %s: found dates [%s]", e.Title, strings.Join(dateStrings, ", "))
	
	// Sort all dates chronologically and use the most recent date to determine logic
	if len(allDates) > 0 {
		// Sort dates chronologically (earliest to latest)
		for i := 0; i < len(allDates); i++ {
			for j := i + 1; j < len(allDates); j++ {
				if allDates[j].Before(allDates[i]) {
					allDates[i], allDates[j] = allDates[j], allDates[i]
				}
			}
		}
		
		// Get the most recent (last) date and check if it's in the past or future
		mostRecentDate := allDates[len(allDates)-1]
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		
		if mostRecentDate.Before(today) || mostRecentDate.Equal(today) {
			// Most recent date is in the past/today - it's the latest release
			out.AudibleLatestDate = &mostRecentDate
			// No next date since the most recent is already released
		} else {
			// Most recent date is in the future - it's a preorder
			out.AudibleNextDate = &mostRecentDate
			// Find the latest past date if there are multiple dates
			if len(allDates) > 1 {
				// Check if second-to-last date is in the past
				secondMostRecent := allDates[len(allDates)-2]
				if secondMostRecent.Before(today) || secondMostRecent.Equal(today) {
					out.AudibleLatestDate = &secondMostRecent
				}
			}
		}
	}
	
	// Log final assigned values
	var latest, next string
	if out.AudibleLatestDate != nil {
		latest = out.AudibleLatestDate.Format("2006-01-02")
	} else {
		latest = "none"
	}
	if out.AudibleNextDate != nil {
		next = out.AudibleNextDate.Format("2006-01-02")
	} else {
		next = "none"
	}
	log.Printf("Audible %s: count=%d, latest=%s, next=%s", e.Title, out.AudibleCount, latest, next)

	return out, nil
}