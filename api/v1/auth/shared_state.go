package auth

import (
	"crypto/rand"
	"sync"
)

type FlowState struct {
	UserAccessKey string
	RefreshToken  string
}

// SharedState manages the state across the authentication cycle. Saves the user access token and refresh token
// between registration, callback and linking
type SharedState struct {
	mu *sync.RWMutex

	// TODO: add auto cleanup of old states
	states map[string]*FlowState
}

func NewSharedState() *SharedState {
	return &SharedState{
		mu:     &sync.RWMutex{},
		states: make(map[string]*FlowState),
	}
}

func (s *SharedState) NewStateID(accessKey string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	stateID := rand.Text()
	s.states[stateID] = &FlowState{
		UserAccessKey: accessKey,
	}

	return stateID
}

func (s *SharedState) SetRefreshToken(stateID string, refreshToken string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, ok := s.states[stateID]
	if !ok {
		return
	}
	state.RefreshToken = refreshToken
}

func (s *SharedState) GetUserAccessKey(stateID string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.states[stateID]
	if !ok {
		return ""
	}
	return state.UserAccessKey
}

func (s *SharedState) PopState(stateID string) *FlowState {
	s.mu.Lock()
	defer s.mu.Unlock()

	flowState, ok := s.states[stateID]
	if !ok {
		return nil
	}

	delete(s.states, stateID)
	return flowState
}
