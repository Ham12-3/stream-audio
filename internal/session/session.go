package session

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// Session represents an active voice call session
type Session struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	State     State
	mu        sync.RWMutex
}

type State string

const (
	StateNew        State = "new"
	StateConnected  State = "connected"
	StateListening  State = "listening"
	StateSpeaking   State = "speaking"
	StateDisconnected State = "disconnected"
)

// Manager handles session lifecycle
type Manager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

// NewManager creates a new session manager
func NewManager() *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
	}
}

// Create creates a new session
func (m *Manager) Create() *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	session := &Session{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		State:     StateNew,
	}

	m.sessions[session.ID] = session
	return session
}

// Get retrieves a session by ID
func (m *Manager) Get(id string) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	session, ok := m.sessions[id]
	return session, ok
}

// Delete removes a session
func (m *Manager) Delete(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, id)
}

// UpdateState updates the session state
func (s *Session) UpdateState(state State) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.State = state
	s.UpdatedAt = time.Now()
}

// GetState returns the current session state
func (s *Session) GetState() State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.State
}
