package scrapers

import (
	"sync"

	"github.com/michaeldvinci/syllabus/internal/models"
)

// CompositeProvider combines multiple providers and runs them concurrently
type CompositeProvider struct {
	Providers []models.Provider
}

// Fetch runs all providers concurrently and merges their results
func (c *CompositeProvider) Fetch(e models.SeriesIDs) (models.SeriesInfo, error) {
	out := models.SeriesInfo{Title: e.Title, AudibleID: e.AudibleID, AmazonASIN: e.AmazonASIN}
	var mu sync.Mutex
	var wg sync.WaitGroup
	var firstErr error

	for _, p := range c.Providers {
		wg.Add(1)
		go func(p models.Provider) {
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

func mergeSeriesInfo(dst, src *models.SeriesInfo) {
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