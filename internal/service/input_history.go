package service

import "slices"

import "fmt"

type inputHistory struct {
	histories     map[string][]string
	selectedIndex map[string]int
}

var inputHistoryInstance *inputHistory
func InputHistory() *inputHistory {
	if inputHistoryInstance == nil {
		inputHistoryInstance = &inputHistory{
			histories:     make(map[string][]string),
			selectedIndex: make(map[string]int),
		}
	}
	return inputHistoryInstance
}

func (h *inputHistory) Count(prompt string) int {
	return len(h.histories[prompt])
}

func (h *inputHistory) Selected(prompt string) int {
	return h.selectedIndex[prompt]
}

func (h *inputHistory) Select(prompt string, index int) error {
	history := h.histories[prompt]
	if index < 0 || index > len(history) {
		return fmt.Errorf("index %d out of bounds: (0, %d)", index, len(history))
	}

	h.selectedIndex[prompt] = index
	return nil
}

func (h *inputHistory) First(prompt string) error {
	return h.Select(prompt, 0)
}

func (h *inputHistory) Previous(prompt string) error {
	return h.Select(prompt, h.Selected(prompt)-1)
}

func (h *inputHistory) Next(prompt string) error {
	return h.Select(prompt, h.Selected(prompt)+1)
}

func (h *inputHistory) Last(prompt string) error {
	return h.Select(prompt, h.Count(prompt)-1)
}

func (h *inputHistory) Reset(prompt string) {
	h.selectedIndex[prompt] = h.Count(prompt)
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

func (h *inputHistory) Get(prompt string) string {
	history := h.histories[prompt]
	currentIndex := h.selectedIndex[prompt]

	if currentIndex >= len(history) {
		return ""
	}

	return history[currentIndex]
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
		return fmt.Errorf("lower index %d out of bounds: (0, %d)", lower, length-1)
	}

	if upper < 0 || upper >= length {
		return fmt.Errorf("upper index %d out of bounds: (0, %d)", upper, length-1)
	}

	currentIndex := h.selectedIndex[prompt]

	// Remove the slice (inclusive)
	newHistory := slices.Delete(history, lower, upper + 1)
	h.histories[prompt] = newHistory

	// Adjust selection index
	removedCount := upper - lower + 1
	if currentIndex > upper {
		h.selectedIndex[prompt] = currentIndex - removedCount
	} else if currentIndex >= lower {
		// Selected item was removed, move to end
		h.selectedIndex[prompt] = len(newHistory)
	}

	return nil
}
