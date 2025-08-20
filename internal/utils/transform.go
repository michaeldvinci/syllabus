package utils

import "github.com/michaeldvinci/syllabus/internal/models"

// ToSeriesIDs converts a slice of Entry to a slice of SeriesIDs
func ToSeriesIDs(entries []models.Entry) []models.SeriesIDs {
	out := make([]models.SeriesIDs, 0, len(entries))
	for _, e := range entries {
		audURL := ExtractURLFromMarkdownLink(e.Audible)
		amzURL := ExtractURLFromMarkdownLink(e.Amazon)
		out = append(out, models.SeriesIDs{
			Title:      e.Title,
			AudibleID:  ExtractAudibleSeriesID(audURL),
			AudibleURL: audURL,
			AmazonASIN: ExtractAmazonASIN(amzURL),
			Original:   e,
		})
	}
	return out
}