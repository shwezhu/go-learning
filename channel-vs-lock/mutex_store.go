package channel_vs_lock

import (
	"sync"
)

func newMuxStore() *store {
	return &store{
		mu:       sync.RWMutex{},
		sessions: make(map[string]*Session),
	}
}

type store struct {
	mu    sync.RWMutex
	sessions map[string]*Session
}

func (m *store) addSession(session *Session) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[session.id] = session
}

func (m *store) getSession(id string) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	session, ok := m.sessions[id]
	if !ok {
		return nil
	}
	return session
}


