package database

import (
	"github.com/michaeldvinci/syllabus/internal/models"
)

// ToSeriesInfo converts database SeriesStats to models.SeriesInfo for compatibility
func (stats SeriesStats) ToSeriesInfo() models.SeriesInfo {
	info := models.SeriesInfo{
		Title: stats.Title,
		
		// Audible data
		AudibleCount:       stats.AudibleCount,
		AudibleLatestTitle: stringValue(stats.AudibleLatestTitle),
		AudibleLatestDate:  stats.AudibleLatestDate,
		AudibleNextTitle:   stringValue(stats.AudibleNextTitle),
		AudibleNextDate:    stats.AudibleNextDate,
		
		// Amazon data
		AmazonCount:       stats.AmazonCount,
		AmazonLatestTitle: stringValue(stats.AmazonLatestTitle),
		AmazonLatestDate:  stats.AmazonLatestDate,
		AmazonNextTitle:   stringValue(stats.AmazonNextTitle),
		AmazonNextDate:    stats.AmazonNextDate,
		
		// IDs
		AudibleID:  stringValue(stats.AudibleID),
		AmazonASIN: stringValue(stats.AmazonASIN),
	}
	
	return info
}

// ToSeriesInfoSlice converts a slice of SeriesStats to models.SeriesInfo
func ToSeriesInfoSlice(stats []SeriesStats) []models.SeriesInfo {
	infos := make([]models.SeriesInfo, len(stats))
	for i, stat := range stats {
		infos[i] = stat.ToSeriesInfo()
	}
	return infos
}

// stringValue safely converts *string to string
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}