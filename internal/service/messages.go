package service

import (
	"fmt"

	"github.com/vrld/ansicht/internal/model"
)

type MessageIndex struct {
	ThreadIdx  int
	MessageIdx int
}

type messages struct {
	threads        []model.Thread
	messageIndex   []MessageIndex
	selectedIndex  int
	markedMessages map[int]MessageIndex
}

var messagesInstance *messages

func Messages() *messages {
	if messagesInstance == nil {
		messagesInstance = &messages{}
	}
	return messagesInstance
}

func (m *messages) SetThreads(threads []model.Thread) {
	m.ClearMarks()
	m.threads = threads
	m.messageIndex = make([]MessageIndex, 0, len(threads)*2)
	for threadIdx, thread := range m.threads {
		// newest messages first
		messageCount := len(thread.Messages)
		for msgIdx := range thread.Messages {
			m.messageIndex = append(m.messageIndex, MessageIndex{
				threadIdx,
				messageCount - msgIdx - 1,
			})
		}
	}
}

func (m *messages) Count() int {
	return len(m.messageIndex)
}

func (m *messages) Selected() int {
	return m.selectedIndex
}

func (m *messages) Select(i int) error {
	if i < 0 || i >= m.Count() {
		return fmt.Errorf("index %d out of bounds: (0, %d)", i, m.Count())
	}
	m.selectedIndex = i
	return nil
}

func (m *messages) IsMarked(i int) bool {
	if i < 0 || i >= m.Count() {
		return false
	}

	_, exists := m.markedMessages[i]
	return exists
}

func (m *messages) Mark(i int) {
	if i < 0 || i >= m.Count() {
		return
	}

	m.markedMessages[i] = m.messageIndex[i]
}

func (m *messages) Unmark(i int) {
	if i < 0 || i >= m.Count() {
		return
	}

	delete(m.markedMessages, i)
}

func (m *messages) ToggleMark(i int) {
	if i < 0 || i >= m.Count() {
		return
	}

	if m.IsMarked(i) {
		delete(m.markedMessages, i)
	} else {
		m.markedMessages[i] = m.messageIndex[i]
	}
}

func (m *messages) ClearMarks() {
	m.markedMessages = make(map[int]MessageIndex)
}

func (m *messages) InvertMarks() {
	expectedSize := m.Count() - len(m.markedMessages)
	newSelection := make(map[int]MessageIndex, expectedSize)
	for row, index := range m.messageIndex {
		if !m.IsMarked(row) {
			newSelection[row] = index
		}
	}
	m.markedMessages = newSelection
}

func (m *messages) Get(i int) *model.Message {
	if i < 0 || i >= m.Count() {
		return nil
	}

	idx := m.messageIndex[i]
	return &m.threads[idx.ThreadIdx].Messages[idx.MessageIdx]
}

func (m *messages) GetAll() []*model.Message {
	res := make([]*model.Message, 0, m.Count())
	for _, idx := range m.messageIndex {
		res = append(res, &m.threads[idx.ThreadIdx].Messages[idx.MessageIdx])
	}

	return res
}

func (m *messages) GetSelected() *model.Message {
	return m.Get(m.selectedIndex)
}

func (m *messages) GetMarked() []*model.Message {
	selected := make([]*model.Message, 0, len(m.markedMessages))
	for _, idx := range m.markedMessages {
		selected = append(selected, &m.threads[idx.ThreadIdx].Messages[idx.MessageIdx])
	}

	return selected
}

func (m *messages) MarkedCount() int {
	return len(m.markedMessages)
}
