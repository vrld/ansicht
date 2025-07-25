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

func MessagesToListItems(messages *service.Messages) []list.Item {
	var items []list.Item

	for row, message := range messages.GetAll() {
		items = append(items, MessageItem{
			Message: message,
			Marked:  messages.IsMarked(row),
		})
	}

	return items
}

// MessageDelegate is a custom delegate for rendering message items
type MessageDelegate struct {
	width int
}

type MessageStyles struct {
	Date      lipgloss.Style
	Sender    lipgloss.Style
	Arrow     lipgloss.Style
	Recipient lipgloss.Style
	Subject   lipgloss.Style
	Tags      lipgloss.Style
}

func (s MessageStyles) Apply(what func(lipgloss.Style) lipgloss.Style) MessageStyles {
	s.Date = what(s.Date)
	s.Sender = what(s.Sender)
	s.Arrow = what(s.Arrow)
	s.Recipient = what(s.Recipient)
	s.Subject = what(s.Subject)
	s.Tags = what(s.Tags)
	return s
}

var (
	messageStylesUnread = MessageStyles{
		Date:      lipgloss.NewStyle().Foreground(lipgloss.Color(colorTertiary)),
		Sender:    lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent)),
		Arrow:     lipgloss.NewStyle().Foreground(lipgloss.Color(colorMuted)),
		Recipient: lipgloss.NewStyle().Foreground(lipgloss.Color(colorSecondary)),
		Subject:   lipgloss.NewStyle().Foreground(lipgloss.Color(colorHighlight)),
		Tags:      lipgloss.NewStyle().Foreground(lipgloss.Color(colorTertiary)),
	}

	messageStylesSeen = messageStylesUnread.Apply(func(s lipgloss.Style) lipgloss.Style {
		return s.Foreground(lipgloss.Color(colorMuted))
	})
)

// Render renders a list item
func (d MessageDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(MessageItem)
	if !ok || item.Message == nil {
		return
	}

	styles := messageStylesUnread
	if item.Message.Flags.Seen {
		styles = messageStylesSeen
	}

	if index == m.Index() {
		bg := lipgloss.Color(colorSecondaryBright)
		fg := lipgloss.Color(colorBackground)
		styles = styles.Apply(func(s lipgloss.Style) lipgloss.Style {
			return s.Background(bg).Foreground(fg)
		})
	} else if item.Marked {
		bg := lipgloss.Color(colorAccentBright)
		fg := lipgloss.Color(colorBackground)
		styles = styles.Apply(func(s lipgloss.Style) lipgloss.Style {
			return s.Background(bg).Foreground(fg)
		})
	}

	line := d.renderLine(item, styles)

	fmt.Fprint(w, line)
}

func (d MessageDelegate) renderLine(item MessageItem, styles MessageStyles) string {
	date := fmt.Sprintf("%11s  ", formatDate(item.Message.Date))
	sender := fmt.Sprintf("%20s", truncate(formatEmailAddress(item.Message.From), 20)) // TODO: use only name (Sander <s@nd.er> => Sander)
	arrow := " → "
	recipient := fmt.Sprintf("%-20s", truncate(formatEmailAddress(item.Message.To), 20))
	tags := "  " + formatTags(item.Message.Tags) // TODO: replace tags (configurable)

	componentWidth := lipgloss.Width(date) + lipgloss.Width(sender) + lipgloss.Width(arrow) + lipgloss.Width(recipient) + lipgloss.Width(tags)
	remainingWidth := max(1, d.width-componentWidth)
	subject := truncate("  "+cleanSubject(item.Message.Subject), remainingWidth)

	var filler string
	if fillerWidth := d.width - componentWidth - lipgloss.Width(subject); fillerWidth > 0 {
		filler = strings.Repeat(" ", fillerWidth)
	}

	return fmt.Sprintf("%s%s%s%s%s%s",
		styles.Date.Render(date),
		styles.Sender.Render(sender),
		styles.Arrow.Render(arrow),
		styles.Recipient.Render(recipient),
		styles.Subject.Render(subject),
		styles.Tags.Render(tags+filler))
}

// Height returns the height of a list item
func (d MessageDelegate) Height() int { return 1 }

// Spacing returns the spacing between list items
func (d MessageDelegate) Spacing() int { return 0 }

// Update handles key messages
func (d MessageDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil // Message selection is handled in the main Update function
}
