package model

import "time"

type Thread struct {
	ID                   string
	Authors              []string
	Subject              string
	Tags                 []string
	NewestDate           time.Time
	OldestDate           time.Time
	CountMatchedMessages int
	Messages             []Message
}
