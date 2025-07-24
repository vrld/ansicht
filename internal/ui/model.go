package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	OnStartup()
	OnKey(keycode string) (handledKey bool)
	HandleInput(input string)
	HandleSpawnResult(msg runtime.SpawnResult)
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
	return tea.Sequence(
		m.searchCurrentQuery(0),
		func() tea.Msg {
			m.Runtime.OnStartup()
			return m.spinner.Tick
		},
	)
}
