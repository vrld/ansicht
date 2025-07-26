package ui

type Refresh struct{}

type QueryNewMsg struct {
	Query string
}
type QueryNextMsg struct{}
type QueryPrevMsg struct{}

type MarksToggleMsg struct{}
type MarksInvertMsg struct{}
type MarksClearMsg struct{}

type StatusSetMsg struct {
	Message string
}

type OpenInputEvent struct {
	Prompt      string
	Placeholder string
}
