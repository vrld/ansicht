package ui

import (
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Notification severity levels
type NotificationLevel int

const (
	NotificationInfo NotificationLevel = iota
	NotificationWarning
	NotificationError
)

// Notification represents a temporary status message
type Notification struct {
	Message   string
	Level     NotificationLevel
	Timestamp time.Time
	ExpiresAt time.Time
}

// NotificationExpiredMsg is sent when a notification expires
type NotificationExpiredMsg struct {
	Notification Notification
}

// AddNotification adds a notification with automatic timeout
func (m *Model) AddNotification(message string, level NotificationLevel, timeoutSeconds float64) tea.Cmd {
	// Default timeouts based on severity
	if timeoutSeconds <= 0 {
		switch level {
		case NotificationInfo:
			timeoutSeconds = 5
		case NotificationWarning:
			timeoutSeconds = 8
		case NotificationError:
			timeoutSeconds = 12
		}
	}

	now := time.Now()
	timeoutDuration := time.Duration(timeoutSeconds * 1000) * time.Millisecond
	notification := Notification{
		Message:   message,
		Level:     level,
		Timestamp: now,
		ExpiresAt: now.Add(timeoutDuration),
	}

	m.notifications = append(m.notifications, notification)

	// Sort by priority (severity desc, then timestamp desc for recency)
	sort.Slice(m.notifications, func(i, j int) bool {
		if m.notifications[i].Level != m.notifications[j].Level {
			return m.notifications[i].Level > m.notifications[j].Level
		}
		return m.notifications[i].Timestamp.After(m.notifications[j].Timestamp)
	})

	// Return command to expire this notification
	return tea.Tick(timeoutDuration, func(t time.Time) tea.Msg {
		return NotificationExpiredMsg{Notification: notification}
	})
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
