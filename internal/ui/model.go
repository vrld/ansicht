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

type RuntimeInterface interface {
	OnStartup() tea.Cmd
	OnKey(keycode string) tea.Cmd
	PushInputHandle(handle string)
	HandleInput(input string) tea.Cmd
	HandleSpawnResult(msg runtime.SpawnResultMsg) tea.Cmd
}

type Model struct {
	Messages     *service.Messages
	Queries      *service.Queries
	InputHistory *service.InputHistory
	Status       *service.Status
	Runtime      RuntimeInterface
	isLoading    bool
	focusSearch  bool

	list    list.Model
	input   textinput.Model
	spinner spinner.Model

	width int
}

func NewModel() *Model {
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
	messageList.SetShowPagination(false)
	messageList.DisableQuitKeybindings()

	// Style the list
	messageList.Styles = list.DefaultStyles()
	messageList.Styles.NoItems = styleListNoItems

	return &Model{
		focusSearch: false,
		input:       ti,
		list:        messageList,
		spinner:     sp,
		width:       defaultWidth,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.searchCurrentQuery(0),
		m.spinner.Tick,
		m.Runtime.OnStartup(),
	)
}

func (m Model) renderStatusLine() string {
	var leftStatus, rightStatus string

	if query, ok := m.Queries.Current(); ok {
		markedCount := m.Messages.MarkedCount()
		totalCount := m.Messages.Count()
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

	// top border is pagination
	pageIndicatorWidth := m.width / m.list.Paginator.TotalPages - 1
	border := ""
	for page := range m.list.Paginator.TotalPages {
		if page > 0 {
			border += " "
		}
		if page == m.list.Paginator.Page {
			border += stylePaginationActivePage.Render(strings.Repeat("🬂", pageIndicatorWidth))
		} else {
			border += stylePaginationInactivePage.Render(strings.Repeat("🬂", pageIndicatorWidth))
		}
	}
	pad := strings.Repeat(" ", (m.width - ((pageIndicatorWidth + 1) * m.list.Paginator.TotalPages)) / 2)
	return statusLine + "\n" + pad + border
}

func (m Model) renderTabs() string {
	var tabs []string
	for i, query := range m.Queries.All() {
		style := styleTabNormal
		if i == m.Queries.SelectedIndex() {
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

	bottomLine := m.Status.Get()

	if m.isLoading {
		bottomLine = fmt.Sprintf("%s Searching...", m.spinner.View())
	}

	if m.focusSearch {
		bottomLine = m.input.View()
	}

	return fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		statusLine,
		m.list.View(),
		tabs,
		bottomLine,
	)
}
