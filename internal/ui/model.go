package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/vrld/ansicht/internal/runtime"
	"github.com/vrld/ansicht/internal/db"
	"github.com/vrld/ansicht/internal/model"
)

// sent when a search completes
type SearchResultMsg struct {
	Result      model.SearchResult
	Error       error
	RowToSelect int
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

// references to a message's position in the thread slice
type MessageIndex struct {
	ThreadIdx  int
	MessageIdx int
}

type Model struct {
	config            *runtime.Runtime
	queries           []model.SearchQuery
	currentQueryIndex int
	threads           []model.Thread
	isLoading         bool
	focusSearch       bool
	markedRows        map[int]MessageIndex
	rowToMessageIndex []MessageIndex

	list    list.Model
	input   textinput.Model
	spinner spinner.Model

	width int
}

func NewModel(config *runtime.Runtime) Model {
	// search box
	ti := textinput.New()
	ti.Placeholder = "tag:unread" // TODO: random *valid* query
	ti.Focus()
	ti.Width = 50

	// loading indicator
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("62"))

	// Default width (will be updated on WindowSizeMsg)
	defaultWidth := 96

	// Create the message list with a custom delegate
	delegate := NewMessageDelegate(defaultWidth)
	messageList := list.New([]list.Item{}, delegate, defaultWidth, 20)
	messageList.SetShowStatusBar(false)
	messageList.SetFilteringEnabled(false)
	messageList.SetShowTitle(false)
	messageList.SetShowHelp(false)
	messageList.DisableQuitKeybindings()

	// Style the list
	listStyles := list.DefaultStyles()
	listStyles.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("231")).
		Background(lipgloss.Color("25")).
		Padding(0, 1)
	listStyles.NoItems = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center)
	messageList.Styles = listStyles

	// Get saved queries
	queries, err := db.GetSavedQueries()
	if err != nil || len(queries) == 0 {
		queries = []model.SearchQuery{
			{Name: "INBOX", Query: "query:INBOX"},
		}
	}

	return Model{
		config:            config,
		queries:           queries,
		currentQueryIndex: 0,
		focusSearch:       false,
		markedRows:        make(map[int]MessageIndex),
		rowToMessageIndex: make([]MessageIndex, 0),
		input:             ti,
		list:              messageList,
		spinner:           sp,
		width:             defaultWidth,
	}
}

func (m *Model) GetMessage(index MessageIndex) *model.Message {
	if index.ThreadIdx < 0 || index.ThreadIdx >= len(m.threads) {
		return nil
	}

	thread := &m.threads[index.ThreadIdx]
	if index.MessageIdx < 0 || index.MessageIdx >= len(thread.Messages) {
		return nil
	}

	return &thread.Messages[index.MessageIdx]
}

func (m *Model) CurrentQuery() *model.SearchQuery {
	if m.currentQueryIndex < len(m.queries) {
		return &m.queries[m.currentQueryIndex]
	}
	return nil
}

func (m *Model) isRowSelected(row int) bool {
	_, selected := m.markedRows[row]
	return selected
}

func (m *Model) resetSelection() {
	m.markedRows = make(map[int]MessageIndex)
	m.updateList()
}

func (m *Model) invertSelection() {
	expectedSize := len(m.rowToMessageIndex) - len(m.markedRows)
	newSelection := make(map[int]MessageIndex, expectedSize)
	for row, index := range m.rowToMessageIndex {
		if !m.isRowSelected(row) {
			newSelection[row] = index
		}
	}
	m.markedRows = newSelection
	m.updateList()
}

func (m *Model) toggleSelection(row int) {
	if row < 0 || row >= len(m.rowToMessageIndex) {
		return
	}

	if m.isRowSelected(row) {
		delete(m.markedRows, row)
	} else {
		m.markedRows[row] = m.rowToMessageIndex[row]
	}

	// Update the list to reflect the selection change
	m.updateList()
}

func (m *Model) GetSelectedMessages() []*model.Message {
	if len(m.rowToMessageIndex) == 0 {
		return nil
	}

	selected := make([]*model.Message, 0, len(m.markedRows)+1)

	active_row := m.list.Index()
	if idx := m.rowToMessageIndex[active_row]; m.isValidMessageIndex(idx) {
		selected = append(selected, &m.threads[idx.ThreadIdx].Messages[idx.MessageIdx])
	}

	for row, idx := range m.markedRows {
		if row != active_row && m.isValidMessageIndex(idx) {
			selected = append(selected, &m.threads[idx.ThreadIdx].Messages[idx.MessageIdx])
		}
	}

	return selected
}

func (m *Model) isValidMessageIndex(idx MessageIndex) bool {
	return (idx.ThreadIdx >= 0) &&
		(idx.ThreadIdx < len(m.threads)) &&
		(idx.MessageIdx >= 0) &&
		(idx.MessageIdx < len(m.threads[idx.ThreadIdx].Messages))
}

// updateList refreshes the list with current thread data
func (m *Model) updateList() {
	m.rowToMessageIndex = make([]MessageIndex, 0, len(m.threads)*2)

	for threadIdx, thread := range m.threads {
		for messageIdx := range thread.Messages {
			m.rowToMessageIndex = append(m.rowToMessageIndex, MessageIndex{
				ThreadIdx:  threadIdx,
				MessageIdx: messageIdx,
			})
		}
	}

	// Create list items
	items := CreateMessageItems(m.threads, m.markedRows)

	// Update the list with new items
	m.list.SetItems(items)

	// Preserve cursor position if possible
	if m.list.Index() >= len(items) && len(items) > 0 {
		m.list.Select(len(items) - 1)
	}
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

	// Build the status line with selection count if needed
	bottom_line := "Press / to search, ←/→ to switch tabs, <space> to select, I to invert selection, q to quit"
	if len(m.markedRows) > 0 {
		bottom_line = fmt.Sprintf("%d selected | %s", len(m.markedRows), bottom_line)
	}

	if m.isLoading {
		bottom_line = fmt.Sprintf("%s Searching...", m.spinner.View())
	}

	if m.focusSearch {
		bottom_line = m.input.View()
	}

	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.list.View(),
		tabs,
		bottom_line,
	)
}
