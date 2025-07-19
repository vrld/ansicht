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
	cache map[string]string // Simple cache for rendered lines
}

// NewMessageDelegate creates a new delegate for rendering messages
func NewMessageDelegate(width int) MessageDelegate {
	d := MessageDelegate{
		width: width,
		cache: make(map[string]string),
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

	// Choose styles based on read state and selection
	dateStyle := styleMsgDateUnread
	senderStyle := styleMsgSenderUnread
	arrowStyle := styleMsgArrowUnread
	recipientStyle := styleMsgRecipientUnread
	subjectStyle := styleMsgSubjectUnread
	tagsStyle := styleMsgTagsUnread

	if item.Message.Flags.Seen {
		// Read messages use dimmed colors
		dateStyle = styleMsgDateRead
		senderStyle = styleMsgSenderRead
		arrowStyle = styleMsgArrowRead
		recipientStyle = styleMsgRecipientRead
		subjectStyle = styleMsgSubjectRead
		tagsStyle = styleMsgTagsRead
	}

	if index == m.Index() {
		// Selected items use selection background
		bg := d.styles.Selected.GetBackground()
		fg := d.styles.Selected.GetForeground()
		dateStyle = dateStyle.Background(bg).Foreground(fg)
		senderStyle = senderStyle.Background(bg).Foreground(fg)
		arrowStyle = arrowStyle.Background(bg).Foreground(fg)
		recipientStyle = recipientStyle.Background(bg).Foreground(fg)
		subjectStyle = subjectStyle.Background(bg).Foreground(fg)
		tagsStyle = tagsStyle.Background(bg).Foreground(fg)
	} else if item.Marked {
		// Marked items use marked background
		bg := d.styles.Marked.GetBackground()
		fg := d.styles.Marked.GetForeground()
		dateStyle = dateStyle.Background(bg).Foreground(fg)
		senderStyle = senderStyle.Background(bg).Foreground(fg)
		arrowStyle = arrowStyle.Background(bg).Foreground(fg)
		recipientStyle = recipientStyle.Background(bg).Foreground(fg)
		subjectStyle = subjectStyle.Background(bg).Foreground(fg)
		tagsStyle = tagsStyle.Background(bg).Foreground(fg)
	}

	line := d.renderLine(item, dateStyle, senderStyle, arrowStyle, recipientStyle, subjectStyle, tagsStyle)

	fmt.Fprint(w, line)
}

func (d MessageDelegate) renderLine(item MessageItem, dateStyle, senderStyle, arrowStyle, recipientStyle, subjectStyle, tagsStyle lipgloss.Style) string {
	date := fmt.Sprintf("%12s  ", formatDate(item.Message.Date))
	sender := truncate(item.Message.From, 20) // TODO: use only name (Sander <s@nd.er> => Sander)
	arrow := "→"
	recipient := truncate(item.Message.To, 20)
	tags := "  " + formatTags(item.Message.Tags) // TODO: replace tags (configurable)

	// Calculate remaining width for subject
	componentWidth := len(date) + len(sender) + len(arrow) + len(recipient) + len(tags)
	remainingWidth := max(d.width-componentWidth, 15)
	subject := truncate("  "+cleanSubject(item.Message.Subject), remainingWidth)

	return fmt.Sprintf("%s%s%s%s%s%s",
		dateStyle.Render(date),
		senderStyle.Render(sender),
		arrowStyle.Render(arrow),
		recipientStyle.Render(recipient),
		subjectStyle.Render(subject),
		tagsStyle.Render(tags))
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
