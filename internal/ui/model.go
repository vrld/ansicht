package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/vrld/ansicht/internal/db"
	"github.com/vrld/ansicht/internal/model"
)

// sent when a search completes
type SearchResultMsg struct {
	Result model.SearchResult
	Error  error
}

// tabStyle defines the styling for tabs
var (
	activeTabBorder = lipgloss.Border{
		Top:         " ",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopRight:    "┌",
		TopLeft:     "┐",
		BottomRight: "╯",
		BottomLeft:  "╰",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopRight:    "┬",
		TopLeft:     "┬",
		BottomRight: "╯",
		BottomLeft:  "╰",
	}

	tabStyle = lipgloss.NewStyle().
			Border(tabBorder, true).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	activeTabStyle = tabStyle.Border(activeTabBorder, true)

	tabGap = tabStyle.BorderBottom(false).BorderLeft(false).BorderRight(false)
)

type Model struct {
	queries           []model.SearchQuery
	currentQueryIndex int
	threads           []model.Thread
	isLoading         bool
	focusSearch       bool

	table   table.Model
	input   textinput.Model
	spinner spinner.Model

	width int
}

func NewModel() Model {
	// search box
	ti := textinput.New()
	ti.Placeholder = "tag:unread" // TODO: random *valid* query
	ti.Focus()
	ti.Width = 50

	// loading indicator
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("62"))

	// results table
	// TODO: columns configurable
	columns := []table.Column{
		{Title: "Date", Width: 10},
		{Title: "Flags", Width: 6},
		{Title: "From", Width: 20},
		{Title: "Subject", Width: 40},
		{Title: "Tags", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	tableStyles.Selected = tableStyles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)
	t.SetStyles(tableStyles)

	queries, err := db.GetSavedQueries()
	if err != nil || len(queries) == 0 {
		queries = []model.SearchQuery{
			{Name: "INBOX", Query: "query:INBOX"},
		}
	}

	return Model{
		queries:           queries,
		currentQueryIndex: 0,
		focusSearch:       false,
		input:             ti,
		table:             t,
		spinner:           sp,
		width:             96,
	}
}

func (m *Model) CurrentQuery() *model.SearchQuery {
	if m.currentQueryIndex < len(m.queries) {
		return &m.queries[m.currentQueryIndex]
	}
	return nil
}

// updates the table with the current thread data
func (m *Model) updateTable() {
	var rows []table.Row

	for _, thread := range m.threads {
		for _, message := range thread.Messages {
			rows = append(rows, table.Row{
				formatDate(message.Date),
				flagsToString(message.Flags),
				truncate(message.From, 20),
				truncate(message.Subject, 40),
				formatTags(message.Tags),
			})
		}
	}

	m.table.SetRows(rows)
}

func (m Model) renderTabs() string {
	var tabs []string
	for i, query := range m.queries {
		style := tabStyle
		if i == m.currentQueryIndex {
			style = activeTabStyle
		}
		tabs = append(tabs, style.Render(query.Name))
	}
	tab_row := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	gap := tabGap.Render(strings.Repeat(" ", max(0, m.width)))
	return lipgloss.JoinHorizontal(lipgloss.Top, tab_row, gap)
}

func (m Model) View() string {
	tabs := m.renderTabs()

	bottom_line := "Press / to search, ←/→ to switch tabs, q to quit"
	if m.isLoading {
		bottom_line = fmt.Sprintf("%s Searching...", m.spinner.View())
	}

	if m.focusSearch {
		bottom_line = m.input.View()
	}

	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.table.View(),
		tabs,
		bottom_line,
	)

}
