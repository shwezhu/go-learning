package channel_vs_lock

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func GenerateRandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz=#!@%&()<>~,.?"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		ret[i] = letters[num.Int64()]
	}
	return string(ret)
}

func BenchmarkMutex(b *testing.B) {
	s := newMuxStore()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			session := newSession(GenerateRandomString(8), 6)
			s.addSession(session)
			session.value = 32
			s.getSession(session.id)
			s.getSession(session.id)
			s.addSession(session)
			s.addSession(session)
			s.getSession(session.id)
			s.getSession(session.id)
			s.getSession(session.id)
			s.addSession(session)
			s.addSession(session)
		}
	})
}

func BenchmarkChannel(b *testing.B) {
	s := newChannelStore()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			session := newSession(GenerateRandomString(8), 6)
			s.addSession(session)
			session.value = 32
			s.getSession(session.id)
			s.getSession(session.id)
			s.addSession(session)
			s.addSession(session)
			s.getSession(session.id)
			s.getSession(session.id)
			s.getSession(session.id)
			s.addSession(session)
			s.addSession(session)
		}
	})
}