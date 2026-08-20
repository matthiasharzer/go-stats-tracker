package api

import (
	"crypto/rand"
	"sync"
)

type SharedState struct {
	mu     *sync.RWMutex
	states map[string]string
}

func NewSharedState() *SharedState {
	return &SharedState{
		mu:     &sync.RWMutex{},
		states: make(map[string]string),
	}
}

func (s *SharedState) NewStateID(sheetID string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	stateID := rand.Text()
	s.states[stateID] = sheetID

	return stateID
}

func (s *SharedState) PopStateID(stateID string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	sheetID, ok := s.states[stateID]
	if !ok {
		return ""
	}

	delete(s.states, stateID)
	return sheetID
}
