package model

import "time"

// Type aliases for type safety
type MessageID string
type Filename string

type Message struct {
	ID       MessageID
	ThreadID string
	Date     time.Time
	Filename Filename
	Tags     []string
	From     string
	To       string
	Subject  string
	Flags    MessageFlags
}

type MessageFlags struct {
	Draft   bool
	Flagged bool
	Passed  bool
	Replied bool
	Seen    bool
	Trashed bool
}
