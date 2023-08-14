package main

import (
	"fmt"
	"sync"
	"time"
)

var keysChannel = make(chan int, 6)
var resultsChannel = make(chan string, 3)

func doResearch(key int) string{
	// assume this research operation consumes a lot of time
	time.Sleep(time.Second * 2)
	return fmt.Sprintf("One research finished, original key is: %v", key)
}

func createWorkerPool(numOfWorker int) {
	var wg sync.WaitGroup
	// create numOfWorker of workers
	for i := 0; i < numOfWorker; i++ {
		wg.Add(1)
		// create a goroutine that looks like can be "reused" by listening keysChannel until keysChannel is closed
		go func(wg *sync.WaitGroup) {
			// run forever until keysChannel is closed
			// because when keysChannel is empty, this code get blocked not break loop
			for key := range keysChannel {
				resultsChannel <- doResearch(key)
			}
			// worker() is a function represents a goroutine, before return, we should make wg--
			wg.Done()
		}(&wg)
	}
	wg.Wait()
	close(resultsChannel)
}

// when use channel, you have to figure out if you need to remind of other goroutines,
// if yes, figure out when to close
// when use sync.WaitGroup three operations you need to do: wg--, wg.Add(1), wg.wait()
// and don't forget to do wg--, otherwise the wg.wait() will never return
func main() {
	// 1. keep listening the resultChannel until resultsChannel is closed
	done := make(chan bool)
	go func(done chan bool) {
		for result := range resultsChannel {
			fmt.Printf("%v\n", result)
		}
		done <- true
	}(done)

	// 2. keep generating key every 1 sec
	// Imagine this function can continuously generate a request from a client every sec
	// And this request will be sent to keysChannel and will be processed by
	// our one of our workers that keep listening the keysChannel
	go func() {
		key := 0
		for {
			time.Sleep(time.Second)
			keysChannel <- key
			key++
		}
	}()

	// 3. create worker pool, and until wg.Wait() in the last line of this function returns
	createWorkerPool(3)

	<-done
	// It's OK to leave a Go channel open forever and never close it.
	// When the channel is no longer used, it will be garbage collected.
	// Close a channel only when it is essential to
	// inform the receiving goroutines that all data has been transmitted.
	// close(done)
}