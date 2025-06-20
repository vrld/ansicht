package model

import "time"

type Message struct {
	ID       string
	ThreadID string
	Date     time.Time
	Filename string
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
