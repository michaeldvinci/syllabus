package utils

import (
	"fmt"
	"strings"
	"time"
	"crypto/md5"

	"github.com/michaeldvinci/syllabus/internal/models"
)

// GenerateICal creates an iCal file content from series info with all "next" dates
func GenerateICal(infos []models.SeriesInfo) string {
	var events []string
	now := time.Now().UTC()

	// Add calendar header
	cal := []string{
		"BEGIN:VCALENDAR",
		"VERSION:2.0",
		"PRODID:-//Syllabus//Book Release Calendar//EN",
		"CALSCALE:GREGORIAN",
		"METHOD:PUBLISH",
		"X-WR-CALNAME:Book Releases",
		"X-WR-CALDESC:Upcoming audiobook and ebook releases",
	}

	for _, info := range infos {
		// Add Audible next release events
		if info.AudibleNextDate != nil {
			title := info.AudibleNextTitle
			if title == "" {
				title = "Next audiobook release"
			}
			event := createEvent(
				fmt.Sprintf("[AU] %s Releases", info.Title),
				title,
				*info.AudibleNextDate,
				now,
				fmt.Sprintf("audible-%s-%d", sanitizeForUID(info.Title), info.AudibleNextDate.Unix()),
			)
			events = append(events, event)
		}

		// Add Amazon next release events
		if info.AmazonNextDate != nil {
			title := info.AmazonNextTitle
			if title == "" {
				title = "Next ebook release"
			}
			event := createEvent(
				fmt.Sprintf("[AM] %s Releases", info.Title),
				title,
				*info.AmazonNextDate,
				now,
				fmt.Sprintf("amazon-%s-%d", sanitizeForUID(info.Title), info.AmazonNextDate.Unix()),
			)
			events = append(events, event)
		}
	}

	// Combine calendar header with events
	cal = append(cal, events...)
	cal = append(cal, "END:VCALENDAR")

	return strings.Join(cal, "\r\n")
}

// createEvent creates a single VEVENT for iCal format
func createEvent(title, description string, eventDate, createdDate time.Time, uid string) string {
	// Format dates for iCal
	// For all-day events, use YYYYMMDD format without time
	eventDateStr := eventDate.Format("20060102")
	// Add one day for DTEND (all-day events need DTEND to be the day after)
	endDate := eventDate.AddDate(0, 0, 1)
	endDateStr := endDate.Format("20060102")
	// Created/modified timestamps still use full datetime format
	createdStr := createdDate.UTC().Format("20060102T150405Z")
	
	// Generate a unique UID using MD5 hash
	uidHash := fmt.Sprintf("%x", md5.Sum([]byte(uid)))

	return strings.Join([]string{
		"BEGIN:VEVENT",
		fmt.Sprintf("UID:%s@syllabus", uidHash),
		fmt.Sprintf("DTSTART;VALUE=DATE:%s", eventDateStr),
		fmt.Sprintf("DTEND;VALUE=DATE:%s", endDateStr),
		fmt.Sprintf("DTSTAMP:%s", createdStr),
		fmt.Sprintf("CREATED:%s", createdStr),
		fmt.Sprintf("LAST-MODIFIED:%s", createdStr),
		fmt.Sprintf("SUMMARY:%s", escapeText(title)),
		fmt.Sprintf("DESCRIPTION:%s", escapeText(description)),
		"STATUS:CONFIRMED",
		"TRANSP:TRANSPARENT",
		"END:VEVENT",
	}, "\r\n")
}

// sanitizeForUID removes characters that might cause issues in UIDs
func sanitizeForUID(text string) string {
	// Replace spaces and special characters with hyphens
	sanitized := strings.ReplaceAll(text, " ", "-")
	sanitized = strings.ReplaceAll(sanitized, ":", "-")
	sanitized = strings.ReplaceAll(sanitized, "/", "-")
	sanitized = strings.ReplaceAll(sanitized, "\\", "-")
	// Remove consecutive hyphens
	for strings.Contains(sanitized, "--") {
		sanitized = strings.ReplaceAll(sanitized, "--", "-")
	}
	return strings.ToLower(sanitized)
}

// escapeText escapes special characters for iCal text fields
func escapeText(text string) string {
	// Escape special characters according to RFC 5545
	text = strings.ReplaceAll(text, "\\", "\\\\")
	text = strings.ReplaceAll(text, ";", "\\;")
	text = strings.ReplaceAll(text, ",", "\\,")
	text = strings.ReplaceAll(text, "\n", "\\n")
	text = strings.ReplaceAll(text, "\r", "\\r")
	return text
}