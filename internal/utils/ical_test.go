package utils

import (
	"strings"
	"testing"
	"time"

	"github.com/michaeldvinci/syllabus/internal/models"
)

func TestGenerateICal(t *testing.T) {
	// Create test data
	testTime1, _ := time.Parse("2006-01-02", "2024-12-25")
	testTime2, _ := time.Parse("2006-01-02", "2025-01-15")

	infos := []models.SeriesInfo{
		{
			Title:           "Test Series 1",
			AudibleNextDate: &testTime1,
			AudibleNextTitle: "New Audio Release",
			AmazonNextDate:  &testTime2,
			AmazonNextTitle: "New Kindle Release",
		},
		{
			Title:           "Test Series 2",
			AudibleNextDate: &testTime2,
			AudibleNextTitle: "Another Audio Release",
			// No Amazon next release
		},
	}

	ical := GenerateICal(infos)

	// Check for required iCal headers
	expectedHeaders := []string{
		"BEGIN:VCALENDAR",
		"VERSION:2.0",
		"PRODID:-//Syllabus//Book Release Calendar//EN",
		"CALSCALE:GREGORIAN",
		"METHOD:PUBLISH",
		"X-WR-CALNAME:Book Releases",
		"X-WR-CALDESC:Upcoming audiobook and ebook releases",
		"END:VCALENDAR",
	}

	for _, header := range expectedHeaders {
		if !strings.Contains(ical, header) {
			t.Errorf("Expected iCal to contain header: %s", header)
		}
	}

	// Check for event content
	expectedEvents := []string{
		"Test Series 1 Releases",
		"Test Series 2 Releases", 
		"New Audio Release",
		"New Kindle Release",
		"Another Audio Release",
		"BEGIN:VEVENT",
		"END:VEVENT",
	}

	for _, event := range expectedEvents {
		if !strings.Contains(ical, event) {
			t.Errorf("Expected iCal to contain event content: %s", event)
		}
	}

	// Count events - should have 3 events (2 for first series, 1 for second)
	eventCount := strings.Count(ical, "BEGIN:VEVENT")
	if eventCount != 3 {
		t.Errorf("Expected 3 events, got %d", eventCount)
	}
}

func TestGenerateICalEmpty(t *testing.T) {
	infos := []models.SeriesInfo{}
	ical := GenerateICal(infos)

	// Should still have calendar structure
	if !strings.Contains(ical, "BEGIN:VCALENDAR") {
		t.Error("Expected iCal to contain calendar header even when empty")
	}
	if !strings.Contains(ical, "END:VCALENDAR") {
		t.Error("Expected iCal to contain calendar footer even when empty")
	}

	// Should have no events
	eventCount := strings.Count(ical, "BEGIN:VEVENT")
	if eventCount != 0 {
		t.Errorf("Expected 0 events for empty input, got %d", eventCount)
	}
}

func TestGenerateICalNoNextDates(t *testing.T) {
	infos := []models.SeriesInfo{
		{
			Title: "Test Series",
			// No next dates set
		},
	}

	ical := GenerateICal(infos)

	// Should have no events since no next dates
	eventCount := strings.Count(ical, "BEGIN:VEVENT")
	if eventCount != 0 {
		t.Errorf("Expected 0 events when no next dates, got %d", eventCount)
	}
}

func TestCreateEvent(t *testing.T) {
	testDate, _ := time.Parse("2006-01-02 15:04:05", "2024-12-25 10:00:00")
	createdDate, _ := time.Parse("2006-01-02 15:04:05", "2024-01-01 12:00:00")

	event := createEvent("Test Title", "Test Description", testDate, createdDate, "test-uid")

	expectedContent := []string{
		"BEGIN:VEVENT",
		"END:VEVENT",
		"SUMMARY:Test Title",
		"DESCRIPTION:Test Description",
		"DTSTART:20241225T100000Z",
		"DTEND:20241225T100000Z",
		"STATUS:CONFIRMED",
		"TRANSP:TRANSPARENT",
	}

	for _, content := range expectedContent {
		if !strings.Contains(event, content) {
			t.Errorf("Expected event to contain: %s", content)
		}
	}

	// Check that UID is present and formatted correctly
	if !strings.Contains(event, "UID:") || !strings.Contains(event, "@syllabus") {
		t.Error("Expected event to contain properly formatted UID")
	}
}

func TestSanitizeForUID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Simple Title", "simple-title"},
		{"Title: With Colons", "title-with-colons"},
		{"Path/With/Slashes", "path-with-slashes"},
		{"Back\\Slash\\Path", "back-slash-path"},
		{"Mixed: /\\Characters", "mixed-characters"},
	}

	for _, tt := range tests {
		result := sanitizeForUID(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeForUID(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestEscapeText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Simple text", "Simple text"},
		{"Text with; semicolon", "Text with\\; semicolon"},
		{"Text with, comma", "Text with\\, comma"},
		{"Text with\nNewline", "Text with\\nNewline"},
		{"Text with\rCarriage", "Text with\\rCarriage"},
		{"Text with\\backslash", "Text with\\\\backslash"},
		{"Complex; text,\nwith\r\\everything", "Complex\\; text\\,\\nwith\\r\\\\everything"},
	}

	for _, tt := range tests {
		result := escapeText(tt.input)
		if result != tt.expected {
			t.Errorf("escapeText(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestGenerateICalWithSpecialCharacters(t *testing.T) {
	testTime, _ := time.Parse("2006-01-02", "2024-12-25")

	infos := []models.SeriesInfo{
		{
			Title:           "Series: With; Special, Characters\nAnd\\Backslashes",
			AudibleNextDate: &testTime,
			AudibleNextTitle: "Title; with, special\ncharacters\\too",
		},
	}

	ical := GenerateICal(infos)

	// Should contain escaped versions of the special characters
	if !strings.Contains(ical, "Series: With\\; Special\\, Characters\\nAnd\\\\Backslashes Releases") {
		t.Error("Expected title to be properly escaped in iCal output")
	}
	if !strings.Contains(ical, "Title\\; with\\, special\\ncharacters\\\\too") {
		t.Error("Expected description to be properly escaped in iCal output")
	}
}