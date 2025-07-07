package runtime

import tea "github.com/charmbracelet/bubbletea"

type RefreshResultsMsg struct{} // TODO: add args: which messages? => reqires more parsing
func RefreshResults() tea.Msg   { return RefreshResultsMsg{} }

// navigation
type QueryNewMsg struct{}
type QueryNextMsg struct{}
type QueryPrevMsg struct{}

func QueryNew() tea.Msg  { return QueryNewMsg{} }
func QueryNext() tea.Msg { return QueryNextMsg{} }
func QueryPrev() tea.Msg { return QueryPrevMsg{} }

// selection
type SelectionToggleMsg struct{}
type SelectionInvertMsg struct{}

func SelectionToggle() tea.Msg { return SelectionToggleMsg{} }
func SelectionInvert() tea.Msg { return SelectionInvertMsg{} }