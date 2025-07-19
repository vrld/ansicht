package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vrld/ansicht/internal/model"
	"github.com/vrld/ansicht/internal/runtime"
	"github.com/vrld/ansicht/internal/service"
)

// sent when a search completes
type SearchResultMsg struct {
	Result      model.SearchResult
	Error       error
	RowToSelect int
}

type KeyReceiver interface {
	OnKey(keycode string) tea.Cmd
}

type InputHandler interface {
	PushInputHandle(handle string)
	HandleInput(input string) tea.Cmd
}

type SpawnHandler interface {
	HandleSpawnResult(msg runtime.SpawnResultMsg) tea.Cmd
}

type Model struct {
	KeyReceiver  KeyReceiver
	messages     *service.Messages
	queries      *service.Queries
	inputHistory *service.InputHistory
	isLoading    bool
	focusSearch  bool

	InputHandler InputHandler
	SpawnHandler SpawnHandler
	list         list.Model
	input        textinput.Model
	spinner      spinner.Model

	width      int
	bottomLine string
}

func NewModel(messages *service.Messages, queries *service.Queries, inputHistory *service.InputHistory) *Model {
	// search box
	ti := textinput.New()
	ti.Placeholder = "tag:unread" // TODO: random *valid* query
	ti.Focus()
	ti.Width = 50

	// loading indicator
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = styleSpinner

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
	listStyles.Title = styleListTitle
	listStyles.NoItems = styleListNoItems
	messageList.Styles = listStyles

	return &Model{
		messages:     messages,
		queries:      queries,
		inputHistory: inputHistory,
		focusSearch:  false,
		input:        ti,
		list:         messageList,
		spinner:      sp,
		width:        defaultWidth,
		bottomLine:   "",
	}
}

func (m Model) renderStatusLine() string {
	var leftStatus, rightStatus string

	if query, ok := m.queries.Current(); ok {
		markedCount := m.messages.MarkedCount()
		totalCount := m.messages.Count()
		currentPos := m.list.Index() + 1

		leftStatus = fmt.Sprintf("%s | %d/%d | %d marked", query.Query, currentPos, totalCount, markedCount)
	} else {
		leftStatus = "No query selected"
	}

	// Add current time on the right
	rightStatus = time.Now().Format("15:04")

	// Calculate spacing to right-align the time
	spacing := max(m.width-2-2-len(leftStatus)-len(rightStatus), 1)

	statusText := "👀 " + leftStatus + strings.Repeat(" ", spacing) + rightStatus
	statusLine := styleStatusLine.Render(statusText)
	border := styleStatusBorder.Render(strings.Repeat("🬂", m.width))

	return statusLine + "\n" + border
}

func (m Model) renderTabs() string {
	var tabs []string
	for i, query := range m.queries.All() {
		style := styleTabNormal
		if i == m.queries.SelectedIndex() {
			style = styleTabActive
		}
		tabs = append(tabs, style.Render(query.Name))
	}
	tab_row := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	gap := styleTabGap.Render(strings.Repeat(" ", max(0, m.width)))
	return lipgloss.JoinHorizontal(lipgloss.Top, tab_row, gap)
}

func (m Model) View() string {
	statusLine := m.renderStatusLine()
	tabs := m.renderTabs()

	bottom_line := m.bottomLine

	if m.isLoading {
		bottom_line = fmt.Sprintf("%s Searching...", m.spinner.View())
	}

	if m.focusSearch {
		bottom_line = m.input.View()
	}

	return fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		statusLine,
		m.list.View(),
		tabs,
		bottom_line,
	)
}
