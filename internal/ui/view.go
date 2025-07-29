package ui

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/vrld/ansicht/internal/service"
)

func borderColor() lipgloss.Color {
	return lipgloss.Color(colorMuted)
}

func (m Model) View() string {
	tabs := m.renderTabs()
	status := m.renderStatusLine()
	mails := m.renderMails(m.height - (lipgloss.Height(tabs) + lipgloss.Height(status)))
	return fmt.Sprintf("%s\n%s\n%s", tabs, mails, status)
}

func (m *Model) setLayoutDimension(width, height int) {
	m.width = width
	m.height = height
	m.list.SetHeight(1)
	m.list.SetWidth(width - 2)
	m.list.SetDelegate(MessageDelegate{width - 2})
}

// TABS

type BorderTabStyle struct {
	Active bool
	First  bool
	Gap    bool
}

func (s BorderTabStyle) lipgloss() lipgloss.Style {
	bottom := "─"
	if s.Active {
		bottom = " "
	}

	bottomLeft := "┴"
	if s.First && s.Active {
		bottomLeft = "│"
	} else if s.First && !s.Active {
		bottomLeft = "├"
	} else if !s.First && s.Active {
		bottomLeft = "┘"
	}

	bottomRight := "┴"
	side := "│"
	if s.Active {
		bottomRight = "└"
	} else if s.Gap {
		bottomRight = "┐"
		side = " "
	}

	border := lipgloss.Border{
		Top:         "─",
		Bottom:      bottom,
		Left:        side,
		Right:       side,
		TopRight:    "╮",
		TopLeft:     "╭",
		BottomRight: bottomRight,
		BottomLeft:  bottomLeft,
	}

	style := lipgloss.NewStyle().
		Border(border, true).
		BorderForeground(borderColor()).
		Padding(0, 1)

	if s.Active {
		style = style.Foreground(lipgloss.Color(colorAccent)).Bold(true)
	}

	return style
}

func (m *Model) renderTabs() string {
	var tabs []string
	tabWidth := 0
	for i, query := range service.Queries().All() {
		queryTab := BorderTabStyle{
			Active: i == service.Queries().SelectedIndex(),
			First:  i == 0,
		}.lipgloss().Render(query.Name)

		tabs = append(tabs, queryTab)
		tabWidth += utf8.RuneCountInString(query.Name) + 4
	}
	tabRow := lipgloss.JoinHorizontal(lipgloss.Bottom, tabs...)

	gap := BorderTabStyle{Gap: true}.lipgloss().
		BorderLeft(false).BorderTop(false).
		Render(strings.Repeat(" ", max(0, m.width-tabWidth-3)))

	return lipgloss.JoinHorizontal(lipgloss.Bottom, tabRow, gap)
}

// MAILS

var mailsBorder = lipgloss.Border{Left: "│", Right: "│", Bottom: "─", BottomLeft: "└", BottomRight: "┘"}

func mailsStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(mailsBorder, false, true, true, true).BorderForeground(borderColor())
}

func (m *Model) renderMails(listHeight int) string {
	m.list.Styles.NoItems = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center, lipgloss.Center).Width(m.width - 2).Height(listHeight - 1)
	m.list.SetHeight(listHeight - 1)
	return mailsStyle().Render(m.list.View())
}

// STATUS LINE

func (m *Model) renderStatusLine() string {
	if m.focusInput {
		return " " + m.input.View()
	}

	var rightStatus string
	rightStatus = fmt.Sprintf("%s Searching...", m.spinner.View())
	if m.isLoading {
	} else if query, ok := service.Queries().Current(); ok {
		markedCount := service.Messages().MarkedCount()
		totalCount := service.Messages().Count()
		currentPos := m.list.Index() + 1

		rightStatus = fmt.Sprintf("%s｜%d/%d｜%d marked", query.Query, currentPos, totalCount, markedCount)
	}
	rightStatus = fmt.Sprintf("👀 %s ｢%s｣", rightStatus, time.Now().Format("15:04"))

	// Calculate spacing to right-align the time
	spacing := max(m.width-5-utf8.RuneCountInString(service.Status().Get())-utf8.RuneCountInString(rightStatus), 1)

	statusText := service.Status().Get() + strings.Repeat(" ", spacing) + rightStatus
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBackground)).
		Background(lipgloss.Color(colorSecondary)).
		Padding(0, 1).
		Bold(true).
		Render(statusText)
}
