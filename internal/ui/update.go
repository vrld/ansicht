package ui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vrld/ansicht/internal/db"
	"github.com/vrld/ansicht/internal/model"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.searchCurrentQuery(), m.spinner.Tick)
}

func (m Model) searchCurrentQuery() tea.Cmd {
	return func() tea.Msg {
		if query := m.CurrentQuery(); query != nil {
			result, err := db.FindThreads(query)
			return SearchResultMsg{Result: result, Error: err}
		}
		return nil
	}
}

func (m *Model) setThreads(threads []model.Thread) {
	m.threads = threads
	m.updateTable()
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
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.table.SetHeight(msg.Height - 4) // Allow for borders, margin, and tabs
		availableWidth := msg.Width - 2   // Account for table borders

		// determine column width.
		// TODO: make configurable
		columns := m.table.Columns()
		dateWidth := 10
		flagsWidth := 6
		fromWidth := 20

		remainingWidth := availableWidth - dateWidth - flagsWidth - fromWidth - 8 // Account for column borders/spacing

		// 70% to subject, 30% to tags
		subjectWidth := int(float64(remainingWidth) * 0.7)
		tagsWidth := remainingWidth - subjectWidth

		columns[0].Width = dateWidth    // Date
		columns[1].Width = flagsWidth   // Flags
		columns[2].Width = fromWidth    // From
		columns[3].Width = subjectWidth // Subject
		columns[4].Width = tagsWidth    // Tags

		m.table.SetColumns(columns)

		return m, nil

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
					m.isLoading = true
					m.table.SetRows([]table.Row{})
					cmd = tea.Batch(m.searchCurrentQuery(), m.spinner.Tick)
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
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit

			case "/":
				m.focusSearch = true
				m.input.Focus()
				return m, nil

			case "left":
				if m.currentQueryIndex > 0 {
					m.currentQueryIndex--
					m.isLoading = true
					m.table.SetRows([]table.Row{})
					return m, tea.Batch(m.searchCurrentQuery(), m.spinner.Tick)
				}

			case "right":
				if m.currentQueryIndex < len(m.queries)-1 {
					m.currentQueryIndex++
					m.isLoading = true
					m.table.SetRows([]table.Row{})
					return m, tea.Batch(m.searchCurrentQuery(), m.spinner.Tick)
				}
			}
		}
	}

	// once we're here, we need to update all our widgets
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	if m.focusSearch {
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.isLoading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
