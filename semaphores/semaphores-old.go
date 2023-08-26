// Share memory by communicating
// https://blog.carlmjohnson.net/post/share-memory-by-communicating/

package semaphores

/*import "sync"

type Semaphore struct {
	sem chan struct{}
	done chan struct{}
	rw sync.RWMutex
}

func New(n int) *Semaphore {
	return &Semaphore{
		sem:  make(chan struct{}, n),
		done: make(chan struct{}),
		rw: sync.RWMutex{},
	}
}

func (s *Semaphore) Acquire() bool {
	s.rw.RLock()
	sem := s.sem
	s.rw.RUnlock()

	select {
	case sem <- struct{}{}:
		return true
	case <-s.done:
		return false
	}
}

func (s *Semaphore) Release() {
	s.rw.RLock()
	sem := s.sem
	s.rw.RUnlock()

	select {
	case <-sem:
	case <-s.done:
	}
}

func (s *Semaphore) Stop() {
	s.rw.Lock()
	defer s.rw.Unlock()

	s.sem = nil
	close(s.done)
}*/