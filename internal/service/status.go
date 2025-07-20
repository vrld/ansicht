package service

import "sync"

type Status struct {
	mu      sync.RWMutex
	message string
}

func NewStatus() *Status {
	return &Status{
		message: "",
	}
}

func (s *Status) Set(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
}

func (s *Status) Get() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.message
}

func (s *Status) Clear() {
	s.Set("")
}
