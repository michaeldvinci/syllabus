package utils

import (
	"regexp"
	"strings"
)

var (
	mdLinkRe        = regexp.MustCompile(`\[[^\]]*\]\(([^)]+)\)`)
	audibleSeriesIDRe = regexp.MustCompile(`/([A-Z0-9]{10})(?:[/?]|$)`)
	amazonASINRe    = regexp.MustCompile(`/dp/([A-Z0-9]{10})(?:[/?]|$)`)
)

// ExtractURLFromMarkdownLink extracts URL from markdown link format
func ExtractURLFromMarkdownLink(s string) string {
	m := mdLinkRe.FindStringSubmatch(s)
	if len(m) == 2 {
		return m[1]
	}
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return s
	}
	return ""
}

// ExtractAudibleSeriesID extracts the Audible series ID from a URL
func ExtractAudibleSeriesID(u string) string {
	if u == "" {
		return ""
	}
	m := audibleSeriesIDRe.FindStringSubmatch(u)
	if len(m) == 2 && strings.HasPrefix(m[1], "B0") {
		return m[1]
	}
	return ""
}

// ExtractAmazonASIN extracts the Amazon ASIN from a URL
func ExtractAmazonASIN(u string) string {
	if u == "" {
		return ""
	}
	m := amazonASINRe.FindStringSubmatch(u)
	if len(m) == 2 {
		return m[1]
	}
	return ""
}