// Share memory by communicating
// https://blog.carlmjohnson.net/post/share-memory-by-communicating/

package semaphores

type Semaphore struct {
	acquire chan bool
	release chan struct{}
	stop chan chan struct{}
}

func New(n int) *Semaphore {
	s := Semaphore{
		acquire:  make(chan bool),
		release:  make(chan struct{}),
		stop:     make(chan chan struct{}),
	}
	go s.start(n)
	return &s
}

func (s *Semaphore) start(max int) {
	count := 0

	for {
		var acquire = s.acquire

		// nil always blocks sends and read operation
		if count >= max {
			acquire = nil
		}

		select {

		case acquire <- true:
			count++

		case s.release <- struct{}{}:
			count--

		case wait := <-s.stop:
			close(s.acquire)

			// Drain remaining calls to Release
			for count > 0 {
				s.release <- struct{}{}
				count--
			}
			close(wait)
			return
		}
	}
}

// Acquire a closed channel returns its default value as many times as it is called.
// if s.acquire is closed, the Acquire() get called in other goroutine will return false immediately
// if s.acquire is not closed, and no data written into it, Acquire() will block
func (s *Semaphore) Acquire() bool {
	return <-s.acquire
}

func (s *Semaphore) Release() {
	<-s.release
}

func (s *Semaphore) Stop() {
	blocker := make(chan struct{})
	s.stop <- blocker
	<-blocker
}