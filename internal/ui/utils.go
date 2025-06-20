package ui

import (
	"strings"
	"time"

	"github.com/vrld/ansicht/internal/model"
)

func formatDate(date time.Time) string {
	now := time.Now()
	if date.Year() == now.Year() {
		return date.Format("Jan 2")
	}
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

func flagsToString(flags model.MessageFlags) string {
	var flagStr strings.Builder

	if flags.Draft {
		flagStr.WriteString("D")
	} else {
		flagStr.WriteString(" ")
	}

	if flags.Flagged {
		flagStr.WriteString("F")
	} else {
		flagStr.WriteString(" ")
	}

	if flags.Passed {
		flagStr.WriteString("P")
	} else {
		flagStr.WriteString(" ")
	}

	if flags.Replied {
		flagStr.WriteString("R")
	} else {
		flagStr.WriteString(" ")
	}

	if !flags.Seen {
		flagStr.WriteString("N")
	} else {
		flagStr.WriteString(" ")
	}

	if flags.Trashed {
		flagStr.WriteString("T")
	} else {
		flagStr.WriteString(" ")
	}

	return flagStr.String()
}
