package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"sync"
	"time"
)

const stateMapLimit = 32_000
const flowStateLifespan = 10 * time.Minute
const cleanupInterval = 1 * time.Minute

type FlowState struct {
	StartedAt     time.Time
	UserAccessKey string
	RefreshToken  string
}

// SharedState manages the state across the authentication cycle. Saves the user access token and refresh token
// between registration, callback and linking
type SharedState struct {
	mu *sync.RWMutex

	states map[string]*FlowState
}

func NewSharedState(ctx context.Context) *SharedState {
	state := &SharedState{
		mu:     &sync.RWMutex{},
		states: make(map[string]*FlowState),
	}
	go state.runCleanupInterval(ctx)
	return state
}

func (s *SharedState) cleanup() {
	now := time.Now()

	for key, state := range s.states {
		if now.Sub(state.StartedAt) > flowStateLifespan {
			delete(s.states, key)
		}
	}
}

func (s *SharedState) safeCleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cleanup()
}

func (s *SharedState) runCleanupInterval(ctx context.Context) {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.safeCleanup()
		}
	}
}

func (s *SharedState) NewStateID(accessKey string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.states) >= stateMapLimit {
		s.cleanup()

		if len(s.states) >= stateMapLimit {
			return "", errors.New("capacity exceeded")
		}
	}

	stateID := rand.Text()
	s.states[stateID] = &FlowState{
		StartedAt:     time.Now(),
		UserAccessKey: accessKey,
	}

	return stateID, nil
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
