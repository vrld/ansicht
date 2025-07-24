package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"time"
)

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
	pageIndicatorWidth := m.width/m.list.Paginator.TotalPages - 1
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
	pad := strings.Repeat(" ", (m.width-((pageIndicatorWidth+1)*m.list.Paginator.TotalPages))/2)
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
