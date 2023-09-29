package channel_vs_lock

import "time"

func newSession(id string, value int) *Session {
	return &Session{
		id:    id,
		value: value,
		expiry:  time.Now().Add(time.Duration(10) * time.Second).Unix(),
	}
}

type Session struct {
	id    string
	value int
	expiry int64
}
