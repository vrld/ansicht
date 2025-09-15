package service

import "slices"

import "fmt"

type inputHistory struct {
	histories map[string][]string
}

var inputHistoryInstance *inputHistory

func init() {
	// create singleton
	_ = InputHistory()
}

func InputHistory() *inputHistory {
	if inputHistoryInstance == nil {
		inputHistoryInstance = &inputHistory{
			histories: make(map[string][]string),
		}
	}
	return inputHistoryInstance
}

func (h *inputHistory) Count(prompt string) int {
	return len(h.histories[prompt])
}

func (h *inputHistory) Get(prompt string, index int) string {
	history := h.histories[prompt]

	if index >= len(history) {
		return ""
	}

	return history[index]
}

func (h *inputHistory) Add(prompt, input string) {
	if input == "" {
		return
	}

	history := h.histories[prompt]

	// Remove existing entry if present
	for i, entry := range history {
		if entry == input {
			history = slices.Delete(history, i, i+1)
			break
		}
	}

	history = append(history, input)
	h.histories[prompt] = history
}

func (h *inputHistory) Remove(prompt string, index int) error {
	return h.RemoveSlice(prompt, index, index)
}

func (h *inputHistory) RemoveSlice(prompt string, lower, upper int) error {
	history := h.histories[prompt]
	length := len(history)

	if length == 0 {
		return fmt.Errorf("no history for prompt: %s", prompt)
	}

	// Wrap negative indices around
	if lower < 0 {
		lower = length + lower
	}
	if upper < 0 {
		upper = length + upper
	}

	// Make sure the slice is valid
	if lower > upper {
		lower, upper = upper, lower
	}
	if lower < 0 || lower >= length {
		return fmt.Errorf("lower index %d out of bounds: 0 <= lower <= %d", lower, length-1)
	}

	if upper < 0 || upper >= length {
		return fmt.Errorf("upper index %d out of bounds: 0 <= upper <= %d", upper, length-1)
	}

	// Remove the slice (inclusive)
	newHistory := slices.Delete(history, lower, upper+1)
	h.histories[prompt] = newHistory

	return nil
}

func (h *inputHistory) Clear(prompt string) {
	h.histories[prompt] = nil
}

func (h *inputHistory) ClearAll() {
	h.histories = make(map[string][]string)
}

func (h *inputHistory) GetSelection(prompt string) *HistorySelection {
	return &HistorySelection{
		prompt: prompt,
		index:  h.Count(prompt),
	}
}

type HistorySelection struct {
	prompt string
	index  int
}

func (s *HistorySelection) Index() int {
	return s.index
}

func (s *HistorySelection) Get() string {
	return InputHistory().Get(s.prompt, s.index)
}

func (s *HistorySelection) Count() int {
	return InputHistory().Count(s.prompt)
}

func (s *HistorySelection) Select(index int) error {
	if index < 0 || index > s.Count() {
		return fmt.Errorf("index %d out of bounds: 0 <= index <= %d", index, s.Count())
	}

	s.index = index
	return nil
}

func (s *HistorySelection) First() error {
	return s.Select(0)
}

func (s *HistorySelection) Previous() error {
	return s.Select(s.index - 1)
}

func (s *HistorySelection) Next() error {
	return s.Select(s.index + 1)
}

func (s *HistorySelection) Last() error {
	return s.Select(s.Count() - 1)
}

func (s *HistorySelection) Reset() {
	s.Select(s.Count())
}
