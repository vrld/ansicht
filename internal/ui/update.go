package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vrld/ansicht/internal/db"
	"github.com/vrld/ansicht/internal/model"
	"github.com/vrld/ansicht/internal/runtime"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SearchResultMsg:
		m.isLoading = false
		if msg.Error != nil {
			// TODO: Show error in UI
			return m, nil
		}
		m.Messages.SetThreads(msg.Result.Threads)
		m.updateList(msg.RowToSelect)
		return m, nil

	case tea.WindowSizeMsg:
		m.setLayoutDimension(msg.Width, msg.Height)
		m.updateList(m.list.Index())

		return m, nil

	// reload the current query
	case Refresh:
		return m, m.loadCurrentQuery(m.list.Index())

	// new query
	case QueryNewMsg:
		m.Queries.Add(model.SearchQuery{
			Query: msg.Query,
			Name:  truncate(msg.Query, 10),
		})
		m.Queries.SelectLast()
		return m, m.loadCurrentQuery(0)

	// switch between queries
	case QueryNextMsg:
		m.Queries.SelectNext()
		return m, m.loadCurrentQuery(0)

	case QueryPrevMsg:
		m.Queries.SelectPrevious()
		return m, m.loadCurrentQuery(0)

	// item selection
	case MarksToggleMsg:
		m.Messages.ToggleMark(m.list.Index())
		m.updateList(m.list.Index())
		return m, nil

	case MarksInvertMsg:
		m.Messages.InvertMarks()
		m.updateList(m.list.Index())
		return m, nil

	case MarksClearMsg:
		m.Messages.ClearMarks()
		m.updateList(m.list.Index())
		return m, nil

	case OpenInputEvent:
		m.focusInput = true
		m.input.Placeholder = msg.Placeholder
		m.input.Prompt = msg.Prompt
		m.input.Focus()
		return m, nil

	case runtime.SpawnResult:
		m.Runtime.HandleSpawnResult(msg)
		return m, nil

	// key presses
	case tea.KeyMsg:
		if m.focusInput {
			switch msg.String() {
			case "enter":
				if query := m.input.Value(); query != "" {
					m.InputHistory.Add(m.input.Prompt, query)
					m.Runtime.HandleInput(query)
				}
				m.focusInput = false
				m.input.Reset()
				return m, nil

			case "esc":
				m.focusInput = false
				m.InputHistory.Reset(m.input.Prompt)
				m.input.Reset()
				return m, nil

			case "up":
				if err := m.InputHistory.Previous(m.input.Prompt); err == nil {
					m.input.SetValue(m.InputHistory.Get(m.input.Prompt))
				}
				return m, nil

			case "down":
				if err := m.InputHistory.Next(m.input.Prompt); err == nil {
					m.input.SetValue(m.InputHistory.Get(m.input.Prompt))
				}
				return m, nil
			}
		} else {
			m.Messages.Select(m.list.Index())
			if m.Runtime.OnKey(msg.String()) {
				return m, nil
			}
		}
	}

	// once we're here, we need to update all our widgets
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.focusInput {
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
		if query, ok := m.Queries.Current(); ok {
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
	items := MessagesToListItems(m.Messages)

	m.list.SetItems(items)
	m.list.Select(toSelect)
	m.Messages.Select(toSelect)
}
