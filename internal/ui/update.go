package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vrld/ansicht/internal/runtime"
	"github.com/vrld/ansicht/internal/db"
	"github.com/vrld/ansicht/internal/model"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.searchCurrentQuery(0), m.spinner.Tick)
}

func (m *Model) searchCurrentQuery(rowToSelect int) tea.Cmd {
	return func() tea.Msg {
		if query := m.CurrentQuery(); query != nil {
			result, err := db.FindThreads(query)
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

func (m *Model) setThreads(threads []model.Thread) {
	m.threads = threads
	m.resetSelection()
	m.updateList()
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
		m.setThreads(msg.Result.Threads)
		m.list.Select(msg.RowToSelect)
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width

		// Update list height and width
		listHeight := msg.Height - 4 // Allow for borders, margin, and tabs
		m.list.SetHeight(listHeight)
		m.list.SetWidth(msg.Width)

		// Update the delegate's width for proper rendering
		delegate := NewMessageDelegate(msg.Width)
		m.list.SetDelegate(delegate)

		// Refresh the list with the new dimensions
		m.updateList()

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
		if len(m.queries) <= 1 {
			return m, nil
		}

		m.currentQueryIndex = (m.currentQueryIndex + 1) % len(m.queries)
		return m, m.loadCurrentQuery(0)

	case runtime.QueryPrevMsg:
		if len(m.queries) <= 1 {
			return m, nil
		}

		m.currentQueryIndex = (m.currentQueryIndex - 1 + len(m.queries)) % len(m.queries)
		return m, m.loadCurrentQuery(0)

	// item selection
	case runtime.SelectionToggleMsg:
		m.toggleSelection(m.list.Index())
		return m, nil

	case runtime.SelectionInvertMsg:
		m.invertSelection()
		return m, nil

	// key presses
	case tea.KeyMsg:
		if m.focusSearch {
			switch msg.String() {
			case "enter":
				var cmd tea.Cmd
				if query := m.input.Value(); query != "" {
					m.queries = append(m.queries, model.SearchQuery{
						Name:  truncate(query, 16),
						Query: query,
					})
					m.currentQueryIndex = len(m.queries) - 1
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
		} else if cmd := m.config.OnKey(msg.String(), m.GetSelectedMessages()); cmd != nil {
			return m, cmd
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
