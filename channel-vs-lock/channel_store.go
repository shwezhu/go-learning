package channel_vs_lock

func newChannelStore() *Store {
	s := &Store{ops: make(chan func(map[string]*Session))}
	go s.loop()
	return s
}

type Store struct {
	ops chan func(map[string]*Session)
}

func (s *Store) addSession(session *Session)  {
	s.ops <- func(m map[string]*Session) {
		// if the key has existed in map, change the value of the key.
		// if key doesn't exist, create a new one
		m[session.id] = session
	}
}

func (s *Store) getSession(id string) *Session {
	result := make(chan *Session, 1)
	s.ops <- func(m map[string]*Session) {
		// everything copied by value, session is a copy of m[id]
		// you should consider if session has pointer field
		session, ok := m[id]
		if !ok {
			result <- nil
			return
		}
		result <- session
	}
	// wait ops finish
	return <-result
}

func (s *Store) loop() {
	sessions := make(map[string]*Session)
	for op := range s.ops {
		op(sessions)
	}
}
