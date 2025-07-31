package ui

import (
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vrld/ansicht/internal/service"
)

// Notification severity levels
type NotificationSeverity int

const (
	NotificationInfo NotificationSeverity = iota
	NotificationWarning
	NotificationError
)

// Notification represents a temporary status message
type Notification struct {
	Message   string
	Severity  NotificationSeverity
	Timestamp time.Time
	ExpiresAt time.Time
}

// NotificationExpiredMsg is sent when a notification expires
type NotificationExpiredMsg struct {
	Notification Notification
}

// AddNotification adds a notification with automatic timeout
func (m *Model) AddNotification(message string, severity NotificationSeverity, timeoutSeconds int) tea.Cmd {
	// Default timeouts based on severity
	if timeoutSeconds <= 0 {
		switch severity {
		case NotificationInfo:
			timeoutSeconds = 5
		case NotificationWarning:
			timeoutSeconds = 8
		case NotificationError:
			timeoutSeconds = 12
		}
	}

	now := time.Now()
	notification := Notification{
		Message:   message,
		Severity:  severity,
		Timestamp: now,
		ExpiresAt: now.Add(time.Duration(timeoutSeconds) * time.Second),
	}

	// Store original status if this is the first notification
	if len(m.notifications) == 0 {
		m.originalStatus = m.getCurrentStatus()
	}

	// Remove any existing notification with the same severity (don't stack)
	m.removeNotificationsBySeverity(severity)

	// Add new notification
	m.notifications = append(m.notifications, notification)

	// Sort by priority (severity desc, then timestamp desc for recency)
	sort.Slice(m.notifications, func(i, j int) bool {
		if m.notifications[i].Severity != m.notifications[j].Severity {
			return m.notifications[i].Severity > m.notifications[j].Severity
		}
		return m.notifications[i].Timestamp.After(m.notifications[j].Timestamp)
	})

	// Return command to expire this notification
	return tea.Tick(time.Duration(timeoutSeconds)*time.Second, func(t time.Time) tea.Msg {
		return NotificationExpiredMsg{Notification: notification}
	})
}

// removeNotificationsBySeverity removes all notifications with the given severity
func (m *Model) removeNotificationsBySeverity(severity NotificationSeverity) {
	filtered := make([]Notification, 0, len(m.notifications))
	for _, n := range m.notifications {
		if n.Severity != severity {
			filtered = append(filtered, n)
		}
	}
	m.notifications = filtered
}

// RemoveExpiredNotification removes a specific expired notification
func (m *Model) RemoveExpiredNotification(expired Notification) {
	filtered := make([]Notification, 0, len(m.notifications))
	for _, n := range m.notifications {
		// Match by timestamp and message to identify the specific notification
		if !(n.Timestamp.Equal(expired.Timestamp) && n.Message == expired.Message) {
			filtered = append(filtered, n)
		}
	}
	m.notifications = filtered
}

// GetCurrentNotification returns the highest priority active notification
func (m *Model) GetCurrentNotification() *Notification {
	now := time.Now()
	for i := range m.notifications {
		if m.notifications[i].ExpiresAt.After(now) {
			return &m.notifications[i]
		}
	}
	return nil
}

// getCurrentStatus gets the current status from the service
func (m *Model) getCurrentStatus() string {
	return service.Status().Get()
}

// parseNotificationSeverity converts string level to NotificationSeverity
func parseNotificationSeverity(level string) NotificationSeverity {
	switch level {
	case "warning":
		return NotificationWarning
	case "error":
		return NotificationError
	default:
		return NotificationInfo
	}
}

// getNotificationStyle returns the appropriate foreground color style for a notification severity
func (m *Model) getNotificationStyle(severity NotificationSeverity) lipgloss.Style {
	switch severity {
	case NotificationWarning:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorWarning))
	case NotificationError:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorError))
	default: // NotificationInfo
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorBackground))
	}
}
