package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"
)

const chanBuf = 16

// row structure that is the same size with []EncDatum
type row struct {
	val            int
	extra1, extra2 int
}

// setupConsumers sets up k consumers that read rows from row channels and sleep
// a short period.
func setupConsumers(tb testing.TB, k int) (*sync.WaitGroup, []chan row) {
	wg := &sync.WaitGroup{}

	c := make([]chan row, k)
	wg.Add(k)
	for i := range c {
		c[i] = make(chan row, chanBuf)
		go func(i int, c chan row) {
			n := 0
			sum := 0
			for r := range c {
				n++
				// Use the row.
				sum += r.val
			}
			fmt.Fprintf(os.Stderr, "%d received %d rows (sum %d)\n", i, n, sum)
			wg.Done()
		}(i, c[i])
	}
	return wg, c
}

// This implements option 1 from the RFC in
// https://github.com/cockroachdb/cockroach/pull/17105: each of the k goroutines
// consumes rows via a channel
func runOption1(b *testing.B, k int, withSema bool) {
	// The semaphore makes sure the main routine blocks if all the consumers are
	// blocked.
	semaphore := make(chan struct{}, k)

	wg, consumerChans := setupConsumers(b, k)

	internalChans := make([]chan row, k)
	for i := range internalChans {
		internalChans[i] = make(chan row, chanBuf)

		wg.Add(1)
		go func(internalChan, consumerChan chan row) {
			var buf []row
			for r := range internalChan {
				wasEmpty := (len(buf) == 0)
				// Buffer the row.
				buf = append(buf, r)
				// Try to send, but don't block.
				for sent := true; len(buf) > 0 && sent; {
					select {
					case consumerChan <- buf[0]:
						buf = buf[1:]
					default:
						sent = false
					}
				}
				isEmpty := (len(buf) == 0)
				if withSema {
					if wasEmpty && !isEmpty {
						// Acquire the semaphore when we have buffered rows.
						semaphore <- struct{}{}
					} else if isEmpty && !wasEmpty {
						// Release the semaphore when we no longer have buffered rows.
						<-semaphore
					}
				}
			}
			if len(buf) > 0 {
				// Flush the buffer.
				for _, r := range buf {
					consumerChan <- r
				}
				if withSema {
					<-semaphore
				}
			}

			close(consumerChan)
			wg.Done()
		}(internalChans[i], consumerChans[i])
	}
	for rowIdx := 0; rowIdx < b.N; rowIdx++ {
		r := row{val: rowIdx}
		i := rand.Intn(k)
		if withSema {
			semaphore <- struct{}{}
		}
		internalChans[i] <- r
		if withSema {
			<-semaphore
		}
	}
	for i := range internalChans {
		close(internalChans[i])
	}
	wg.Wait()
}

func BenchmarkOption1(b *testing.B) {
	// XXX option 1 deadlocks with semaphore!
	for _, sema := range []bool{false} {
		name := "Semaphore"
		if !sema {
			name = "NoSemaphore"
		}
		b.Run(name, func(b *testing.B) {
			for _, k := range []int{2, 4, 16, 64} {
				b.Run(strconv.Itoa(k), func(b *testing.B) {
					runOption1(b, k, sema)
				})
			}
		})
	}
}

// This implements option 2 from the RFC in
// https://github.com/cockroachdb/cockroach/pull/17105: the main routine buffers
// rows directly; each of the k goroutines listens on a condition variable.
func runOption2a(b *testing.B, k int, withSema bool) {
	// The semaphore makes sure the main routine blocks if all the consumers are
	// blocked.
	semaphore := make(chan struct{}, k)

	wg, consumerChans := setupConsumers(b, k)

	type consumerState struct {
		consumerChan chan row
		done         bool
		buf          []row
		mutex        sync.Mutex
		cond         *sync.Cond
	}

	consumers := make([]consumerState, k)
	for i := range consumers {
		consumers[i].consumerChan = consumerChans[i]
		consumers[i].cond = sync.NewCond(&consumers[i].mutex)

		wg.Add(1)
		go func(c *consumerState) {
		Outer:
			for {
				c.mutex.Lock()
				for len(c.buf) == 0 {
					if c.done {
						c.mutex.Unlock()
						break Outer
					}
					c.cond.Wait()
				}
				row := c.buf[0]
				c.buf = c.buf[1:]
				c.mutex.Unlock()
				if withSema {
					// Acquire the semaphore while we are blocked sending.
					semaphore <- struct{}{}
				}
				c.consumerChan <- row
				if withSema {
					<-semaphore
				}
			}

			close(c.consumerChan)
			wg.Done()
		}(&consumers[i])
	}
	for rowIdx := 0; rowIdx < b.N; rowIdx++ {
		r := row{val: rowIdx}
		i := rand.Intn(k)
		c := &consumers[i]
		if withSema {
			semaphore <- struct{}{}
		}
		c.mutex.Lock()
		c.buf = append(c.buf, r)
		c.mutex.Unlock()
		if withSema {
			<-semaphore
		}
		c.cond.Signal()
	}
	for i := range consumers {
		c := &consumers[i]
		c.mutex.Lock()
		c.done = true
		c.mutex.Unlock()
		c.cond.Signal()
	}
	wg.Wait()
}

func BenchmarkOption2a(b *testing.B) {
	for _, sema := range []bool{false, true} {
		name := "Semaphore"
		if !sema {
			name = "NoSemaphore"
		}
		b.Run(name, func(b *testing.B) {
			for _, k := range []int{2, 4, 16, 64} {
				b.Run(strconv.Itoa(k), func(b *testing.B) {
					runOption2a(b, k, sema)
				})
			}
		})
	}
}

// This implements option 2 from the RFC in
// https://github.com/cockroachdb/cockroach/pull/17105: the main routine buffers
// rows directly; each of the k goroutines listens on channel (instead of a
// condition variable in option 2a).
func runOption2b(b *testing.B, k int, withSema bool) {
	// The semaphore makes sure the main routine blocks if all the consumers are
	// blocked.
	semaphore := make(chan struct{}, k)

	wg, consumerChans := setupConsumers(b, k)

	type consumerState struct {
		consumerChan chan row
		done         bool
		buf          []row
		mutex        sync.Mutex
		wakeup       chan struct{}
	}

	consumers := make([]consumerState, k)
	for i := range consumers {
		consumers[i].consumerChan = consumerChans[i]
		consumers[i].wakeup = make(chan struct{}, 1)

		wg.Add(1)
		go func(c *consumerState) {
		Outer:
			for {
				c.mutex.Lock()
				for len(c.buf) == 0 {
					if c.done {
						c.mutex.Unlock()
						break Outer
					}
					c.mutex.Unlock()
					<-c.wakeup
					c.mutex.Lock()
				}
				row := c.buf[0]
				c.buf = c.buf[1:]
				c.mutex.Unlock()
				if withSema {
					// Acquire the semaphore while we are blocked sending.
					semaphore <- struct{}{}
				}
				c.consumerChan <- row
				if withSema {
					<-semaphore
				}
			}

			close(c.consumerChan)
			wg.Done()
		}(&consumers[i])
	}
	for rowIdx := 0; rowIdx < b.N; rowIdx++ {
		r := row{val: rowIdx}
		i := rand.Intn(k)
		c := &consumers[i]
		if withSema {
			semaphore <- struct{}{}
		}
		c.mutex.Lock()
		c.buf = append(c.buf, r)
		c.mutex.Unlock()
		if withSema {
			<-semaphore
		}
		select {
		case c.wakeup <- struct{}{}:
		default:
		}
	}
	for i := range consumers {
		c := &consumers[i]
		c.mutex.Lock()
		c.done = true
		c.mutex.Unlock()
		select {
		case c.wakeup <- struct{}{}:
		default:
		}
	}
	wg.Wait()
}

func BenchmarkOption2b(b *testing.B) {
	for _, sema := range []bool{false, true} {
		name := "Semaphore"
		if !sema {
			name = "NoSemaphore"
		}
		b.Run(name, func(b *testing.B) {
			for _, k := range []int{2, 4, 16, 64} {
				b.Run(strconv.Itoa(k), func(b *testing.B) {
					runOption2b(b, k, sema)
				})
			}
		})
	}
}
