package service

import "sync"

type status struct {
	mu      sync.RWMutex
	message string
}

var statusInstance *status

func Status() *status {
	if statusInstance == nil {
		statusInstance = &status{
			message: "",
		}
	}
	return statusInstance
}

func (s *status) Set(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
}

func (s *status) Get() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.message
}

func (s *status) Clear() {
	s.Set("")
}
