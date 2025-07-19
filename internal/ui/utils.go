package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/vrld/ansicht/internal/model"
)

func formatDate(date time.Time) string {
	now := time.Now()
	diff := now.Sub(date)

	// Less than 1 hour
	if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes < 1 {
			return "now"
		}
		return fmt.Sprintf("%dm ago", minutes)
	}

	// Less than 24 hours
	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%dh ago", hours)
	}

	// Less than 7 days
	if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%dd ago", days)
	}

	// Same year
	if date.Year() == now.Year() {
		return date.Format("Jan 2")
	}

	// Different year
	return date.Format("2006-01-02")
}

// Join tags with a separator
func formatTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	return strings.Join(tags, ",")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}

// cleanSubject removes leading/trailing whitespace and newlines from subject
func cleanSubject(subject string) string {
	// Replace newlines and tabs with spaces, then trim
	cleaned := strings.ReplaceAll(subject, "\n", " ")
	cleaned = strings.ReplaceAll(cleaned, "\t", " ")
	cleaned = strings.TrimSpace(cleaned)

	// Collapse multiple spaces into single spaces
	for strings.Contains(cleaned, "  ") {
		cleaned = strings.ReplaceAll(cleaned, "  ", " ")
	}

	return cleaned
}
