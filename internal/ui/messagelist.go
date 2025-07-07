package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vrld/ansicht/internal/model"
	"github.com/vrld/ansicht/internal/service"
)

// MessageItem represents a single message in the list
type MessageItem struct {
	Message *model.Message
	Marked  bool
}

// FilterValue returns the value used for filtering the list
func (i MessageItem) FilterValue() string {
	if i.Message == nil {
		return ""
	}
	// Return all searchable fields concatenated
	return fmt.Sprintf("%s %s %s",
		i.Message.From,
		i.Message.Subject,
		strings.Join(i.Message.Tags, " "))
}

// MessageDelegate is a custom delegate for rendering message items
type MessageDelegate struct {
	styles struct {
		Normal   lipgloss.Style
		Selected lipgloss.Style
		Marked   lipgloss.Style
		Dim      lipgloss.Style
	}
	width int
}

// NewMessageDelegate creates a new delegate for rendering messages
func NewMessageDelegate(width int) MessageDelegate {
	d := MessageDelegate{
		width: width,
	}

	// Set up styles for different states
	d.styles.Normal = styleMessageNormal.Width(width)
	d.styles.Selected = styleMessageSelected.Width(width)
	d.styles.Marked = styleMessageMarked.Width(width)
	d.styles.Dim = styleMessageDim.Width(width)

	return d
}

// Height returns the height of a list item
func (d MessageDelegate) Height() int { return 1 }

// Spacing returns the spacing between list items
func (d MessageDelegate) Spacing() int { return 0 }

// Update handles key messages
func (d MessageDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil // Message selection is handled in the main Update function
}

// Render renders a list item
func (d MessageDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(MessageItem)
	if !ok || item.Message == nil {
		return
	}

	// Choose style based on selection state
	style := d.styles.Normal
	if index == m.Index() {
		style = d.styles.Selected
	} else if item.Marked {
		style = d.styles.Marked
	}

	// Calculate widths for each column based on total width
	totalWidth := d.width
	dateWidth := 10
	flagsWidth := 7
	fromWidth := int(float64(totalWidth) * 0.2)
	tagsWidth := int(float64(totalWidth) * 0.15)
	subjectWidth := totalWidth - dateWidth - flagsWidth - fromWidth - tagsWidth - 5 // 5 for spacing

	// Truncate and format each field
	date := formatDate(item.Message.Date)
	flags := flagsToString(item.Message.Flags)
	if item.Marked {
		flags += "•"
	} else {
		flags += " "
	}
	from := truncate(item.Message.From, fromWidth)
	subject := truncate(item.Message.Subject, subjectWidth)
	tags := truncate(formatTags(item.Message.Tags), tagsWidth)

	// Format the message line
	str := fmt.Sprintf("%s %s %-*s %-*s %-*s",
		date,
		flags,
		fromWidth, from,
		subjectWidth, subject,
		tagsWidth, tags)

	fmt.Fprint(w, style.Render(str))
}

func CreateMessageItems(messages *service.Messages) []list.Item {
	var items []list.Item

	for row, message := range messages.GetAll() {
		items = append(items, MessageItem{
			Message: message,
			Marked:  messages.IsMarked(row),
		})
	}

	return items
}
