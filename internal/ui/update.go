package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vrld/ansicht/internal/db"
	"github.com/vrld/ansicht/internal/runtime"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.searchCurrentQuery(0), m.spinner.Tick)
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SearchResultMsg:
		m.isLoading = false
		if msg.Error != nil {
			// TODO: Show error in UI
			return m, nil
		}
		m.messages.SetThreads(msg.Result.Threads)
		m.updateList(msg.RowToSelect)
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width

		// Update list height and width
		listHeight := msg.Height - 6 // Allow for borders, margin, and tabs
		m.list.SetHeight(listHeight)
		m.list.SetWidth(msg.Width)

		// Update the delegate's width for proper rendering
		delegate := NewMessageDelegate(msg.Width)
		m.list.SetDelegate(delegate)

		// Refresh the list with the new dimensions
		m.updateList(m.list.Index())

		return m, nil

	// reload the current query
	case runtime.RefreshResultsMsg:
		return m, m.loadCurrentQuery(m.list.Index())

	// start search
	case runtime.QueryNewMsg:
		m.focusSearch = true
		m.input.Focus()
		return m, nil

	// switch between queries
	case runtime.QueryNextMsg:
		m.queries.SelectNext()
		return m, m.loadCurrentQuery(0)

	case runtime.QueryPrevMsg:
		m.queries.SelectPrevious()
		return m, m.loadCurrentQuery(0)

	// item selection
	case runtime.SelectionToggleMsg:
		m.messages.ToggleMark(m.list.Index())
		m.updateList(m.list.Index())
		return m, nil

	case runtime.SelectionInvertMsg:
		m.messages.InvertMarks()
		m.updateList(m.list.Index())
		return m, nil

	// key presses
	case tea.KeyMsg:
		if m.focusSearch {
			switch msg.String() {
			case "enter":
				var cmd tea.Cmd
				if query := m.input.Value(); query != "" {
					m.queries.AddQuery(truncate(query, 16), query)
					m.queries.SelectLast()
					cmd = m.loadCurrentQuery(0)
				}

				m.focusSearch = false
				m.input.Reset()
				return m, cmd

			case "esc":
				m.focusSearch = false
				m.input.Reset()
				return m, nil
			}
		} else {
			m.messages.Select(m.list.Index())
			if cmd := m.KeyReceiver.OnKey(msg.String()); cmd != nil {
				return m, cmd
			}
		}
	}

	// once we're here, we need to update all our widgets
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.focusSearch {
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.isLoading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) searchCurrentQuery(rowToSelect int) tea.Cmd {
	return func() tea.Msg {
		if query, ok := m.queries.Current(); ok {
			result, err := db.FindThreads(&query)
			return SearchResultMsg{Result: result, Error: err, RowToSelect: rowToSelect}
		}
		return nil
	}
}

func (m *Model) loadCurrentQuery(rowToSelect int) tea.Cmd {
	m.isLoading = true
	m.list.SetItems([]list.Item{})
	return tea.Batch(m.searchCurrentQuery(rowToSelect), m.spinner.Tick)
}

func (m *Model) updateList(toSelect int) {
	items := CreateMessageItems(m.messages)

	m.list.SetItems(items)
	m.list.Select(toSelect)
	m.messages.Select(toSelect)
}
