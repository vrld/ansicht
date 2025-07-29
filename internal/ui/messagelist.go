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

func ListItemsFromMessages() []list.Item {
	var items []list.Item

	for row, message := range service.Messages().GetAll() {
		items = append(items, MessageItem{
			Message: message,
			Marked:  service.Messages().IsMarked(row),
		})
	}

	return items
}

// MessageDelegate is a custom delegate for rendering message items
type MessageDelegate struct {
	width int
}

type messageStyles struct {
	Date      lipgloss.Style
	Sender    lipgloss.Style
	Arrow     lipgloss.Style
	Recipient lipgloss.Style
	Subject   lipgloss.Style
	Tags      lipgloss.Style
}

func (s messageStyles) withForeground(color lipgloss.Color) messageStyles {
	return messageStyles{
		Date:      s.Date.Foreground(color),
		Sender:    s.Sender.Foreground(color),
		Arrow:     s.Arrow.Foreground(color),
		Recipient: s.Recipient.Foreground(color),
		Subject:   s.Subject.Foreground(color),
		Tags:      s.Tags.Foreground(color),
	}
}

func (s messageStyles) withBackground(color lipgloss.Color) messageStyles {
	return messageStyles{
		Date:      s.Date.Background(color),
		Sender:    s.Sender.Background(color),
		Arrow:     s.Arrow.Background(color),
		Recipient: s.Recipient.Background(color),
		Subject:   s.Subject.Background(color),
		Tags:      s.Tags.Background(color),
	}
}

func messageStylesUnread() messageStyles {
	return messageStyles{
		Date:      lipgloss.NewStyle().Foreground(lipgloss.Color(colorTertiary)),
		Sender:    lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent)),
		Arrow:     lipgloss.NewStyle().Foreground(lipgloss.Color(colorMuted)),
		Recipient: lipgloss.NewStyle().Foreground(lipgloss.Color(colorSecondary)),
		Subject:   lipgloss.NewStyle().Foreground(lipgloss.Color(colorHighlight)),
		Tags:      lipgloss.NewStyle().Foreground(lipgloss.Color(colorTertiary)),
	}.withBackground(lipgloss.Color(colorBackground))
}

func messageStylesSeen() messageStyles {
	return messageStylesUnread().withForeground(lipgloss.Color(colorMuted))
}

// Render renders a list item
func (d MessageDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(MessageItem)
	if !ok || item.Message == nil {
		return
	}

	var styles messageStyles
	if item.Message.Flags.Seen {
		styles = messageStylesSeen()
	} else {
		styles = messageStylesUnread()
	}

	if index == m.Index() {
		styles = styles.
			withForeground(lipgloss.Color(colorBackground)).
			withBackground(lipgloss.Color(colorSecondaryBright))
	} else if item.Marked {
		styles = styles.
			withForeground(lipgloss.Color(colorTertiaryBright)).
			withBackground(lipgloss.Color(colorBackground))
	}

	line := d.renderLine(item, styles)

	fmt.Fprint(w, line)
}

func (d MessageDelegate) renderLine(item MessageItem, styles messageStyles) string {
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
