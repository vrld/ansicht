package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vrld/ansicht/internal/model"
	"github.com/vrld/ansicht/internal/runtime"
)

// sent when a search completes
type SearchResultMsg struct {
	Result      model.SearchResult
	Error       error
	RowToSelect int
	QueryString string
}

type RuntimeInterface interface {
	OnStartup()
	OnKey(keycode string) (handledKey bool)
	HandleInput(input string)
	HandleSpawnResult(msg runtime.SpawnResult)
}

type Model struct {
	runtime            RuntimeInterface
	isLoading          bool
	focusInput         bool
	currentQueryString string
	list               list.Model
	input              textinput.Model
	spinner            spinner.Model
	width              int
	height             int
	notifications      []Notification
	originalStatus     string
}

func NewModel(runtime RuntimeInterface) *Model {
	// search box
	ti := textinput.New()
	ti.Placeholder = "tag:unread" // TODO: random *valid* query
	ti.Focus()
	ti.Width = 50
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorSecondaryBright))

	// loading indicator
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent))

	// Default width (will be updated on WindowSizeMsg)
	defaultWidth := 96

	// Create the message list with a custom delegate
	messageList := list.New([]list.Item{}, MessageDelegate{defaultWidth}, defaultWidth, 20)
	messageList.SetShowStatusBar(false)
	messageList.SetFilteringEnabled(false)
	messageList.SetShowTitle(false)
	messageList.SetShowHelp(false)
	messageList.SetShowPagination(false)
	messageList.DisableQuitKeybindings()

	// Style the list
	messageList.Styles = list.DefaultStyles()
	messageList.Styles.NoItems = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center, lipgloss.Center)

	return &Model{
		runtime:       runtime,
		focusInput:    false,
		input:         ti,
		list:          messageList,
		spinner:       sp,
		width:         defaultWidth,
		notifications: make([]Notification, 0),
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		m.runtime.OnStartup()
		return Refresh{}
	}
}
