package api

import (
	"crypto/rand"
	"sync"
)

type AuthFlowState struct {
	UserAccessKey string
	RefreshToken  string
}

type SharedState struct {
	mu     *sync.RWMutex
	states map[string]AuthFlowState
}

func NewSharedState() *SharedState {
	return &SharedState{
		mu:     &sync.RWMutex{},
		states: make(map[string]AuthFlowState),
	}
}

func (s *SharedState) NewStateID(accessKey string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	stateID := rand.Text()
	s.states[stateID] = AuthFlowState{
		UserAccessKey: accessKey,
	}

	return stateID
}

func (s *SharedState) UpdateState(stateID string, state AuthFlowState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.states[stateID] = state
}

func (s *SharedState) PeekState(stateID string) *AuthFlowState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	accessKey, ok := s.states[stateID]
	if !ok {
		return nil
	}
	return &accessKey
}

func (s *SharedState) PopState(stateID string) *AuthFlowState {
	s.mu.Lock()
	defer s.mu.Unlock()

	accessKey, ok := s.states[stateID]
	if !ok {
		return nil
	}

	delete(s.states, stateID)
	return &accessKey
}
