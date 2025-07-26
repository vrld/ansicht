package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vrld/ansicht/internal/db"
	"github.com/vrld/ansicht/internal/model"
	"github.com/vrld/ansicht/internal/runtime"
	"github.com/vrld/ansicht/internal/service"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SearchResultMsg:
		// Only process if this result matches the current query
		if msg.QueryString == m.currentQueryString {
			m.isLoading = false
			if msg.Error != nil {
				service.Logger().Error(msg.Error.Error())
				// TODO: Show error in UI
				return m, nil
			}
			service.Messages().SetThreads(msg.Result.Threads)
			m.updateList(msg.RowToSelect)
		}
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
		service.Queries().Add(model.SearchQuery{
			Query: msg.Query,
			Name:  truncate(msg.Query, 10),
		})
		service.Queries().SelectLast()
		return m, m.loadCurrentQuery(0)

	// switch between queries
	case QueryNextMsg:
		service.Queries().SelectNext()
		return m, m.loadCurrentQuery(0)

	case QueryPrevMsg:
		service.Queries().SelectPrevious()
		return m, m.loadCurrentQuery(0)

	// item selection
	case MarksToggleMsg:
		service.Messages().ToggleMark(m.list.Index())
		m.updateList(m.list.Index())
		return m, nil

	case MarksInvertMsg:
		service.Messages().InvertMarks()
		m.updateList(m.list.Index())
		return m, nil

	case MarksClearMsg:
		service.Messages().ClearMarks()
		m.updateList(m.list.Index())
		return m, nil

	case OpenInputEvent:
		m.focusInput = true
		m.input.Placeholder = msg.Placeholder
		m.input.Prompt = msg.Prompt
		m.input.Focus()
		return m, nil

	case runtime.SpawnResult:
		m.runtime.HandleSpawnResult(msg)
		return m, nil

	// key presses
	case tea.KeyMsg:
		if m.focusInput {
			switch msg.String() {
			case "enter":
				if query := m.input.Value(); query != "" {
					service.InputHistory().Add(m.input.Prompt, query)
					m.runtime.HandleInput(query)
				}
				m.focusInput = false
				m.input.Reset()
				return m, nil

			case "esc":
				m.focusInput = false
				service.InputHistory().Reset(m.input.Prompt)
				m.input.Reset()
				return m, nil

			case "up":
				if err := service.InputHistory().Previous(m.input.Prompt); err == nil {
					m.input.SetValue(service.InputHistory().Get(m.input.Prompt))
				}
				return m, nil

			case "down":
				if err := service.InputHistory().Next(m.input.Prompt); err == nil {
					m.input.SetValue(service.InputHistory().Get(m.input.Prompt))
				}
				return m, nil
			}
		} else {
			service.Messages().Select(m.list.Index())
			if m.runtime.OnKey(msg.String()) {
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
		if query, ok := service.Queries().Current(); ok {
			result, err := db.FindThreads(&query)
			return SearchResultMsg{Result: result, Error: err, RowToSelect: rowToSelect, QueryString: query.Query}
		}
		return nil
	}
}

func (m *Model) loadCurrentQuery(rowToSelect int) tea.Cmd {
	if query, ok := service.Queries().Current(); ok {
		m.currentQueryString = query.Query
		m.isLoading = true
		m.list.SetItems([]list.Item{})
		return tea.Batch(m.searchCurrentQuery(rowToSelect), m.spinner.Tick)
	}
	return nil
}

func (m *Model) updateList(toSelect int) {
	items := ListItemsFromMessages()

	m.list.SetItems(items)
	m.list.Select(toSelect)
	service.Messages().Select(toSelect)
}
