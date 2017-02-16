// Results:
//   BenchmarkChan1-4 	 3000000	       532 ns/op
//   BenchmarkChan2-4 	 1000000	      1122 ns/op
//   BenchmarkChan1a-4	 3000000	       599 ns/op
//   BenchmarkChan2a-4	 2000000	       848 ns/op
//
// Conclusion: single channel with slice+error is better than separate data and
// err channels.

package main

import (
	"testing"
	"time"
)

const chanBuf = 10

// Variant 1: single channel with multi-purpose message
type chanMsg struct {
	data []int
	err  error
}

var sum int

func receiverChan1(c <-chan chanMsg) {
	for m := range c {
		for x := range m.data {
			sum += x
		}
		if m.err != nil {
			return
		}
	}
}

func BenchmarkChan1(b *testing.B) {
	c := make(chan chanMsg, chanBuf)
	go receiverChan1(c)
	d := make([]int, 1)
	for i := 0; i < b.N; i++ {
		d[0] = i
		c <- chanMsg{d, nil}
	}
	close(c)
}

// Variant 1a: single channel with multi-purpose message, plus ticker channel

func receiverChan1a(c <-chan chanMsg, t <-chan time.Time) {
	for {
		select {
		case m := <-c:
			for x := range m.data {
				sum += x
			}
			if m.err != nil {
				return
			}
		case <-t:
		}
	}
}

func BenchmarkChan1a(b *testing.B) {
	t := time.Tick(100 * time.Microsecond)
	c := make(chan chanMsg, chanBuf)
	go receiverChan1a(c, t)
	d := make([]int, 1)
	for i := 0; i < b.N; i++ {
		d[0] = i
		c <- chanMsg{d, nil}
	}
	close(c)
}

// Variant 2: two channels

func receiverChan2(dataChan <-chan []int, errChan <-chan error) {
	for {
		select {
		case d := <-dataChan:
			for x := range d {
				sum += x
			}
		case <-errChan:
			return
		}
	}
}

func BenchmarkChan2(b *testing.B) {
	dataChan := make(chan []int, chanBuf)
	errChan := make(chan error)
	go receiverChan2(dataChan, errChan)
	d := make([]int, 1)
	for i := 0; i < b.N; i++ {
		d[0] = i
		dataChan <- d
	}
	close(errChan)
}

// Variant 2a: two channels, plus ticker channel

func receiverChan2a(dataChan <-chan []int, errChan <-chan error, t <-chan time.Time) {
	for {
		select {
		case d := <-dataChan:
			for x := range d {
				sum += x
			}
		case <-errChan:
			return
		case <-t:
		}
	}
}

func BenchmarkChan2a(b *testing.B) {
	t := time.Tick(100 * time.Microsecond)
	dataChan := make(chan []int, chanBuf)
	errChan := make(chan error)
	go receiverChan2a(dataChan, errChan, t)
	d := make([]int, 1)
	for i := 0; i < b.N; i++ {
		d[0] = i
		dataChan <- d
	}
	close(errChan)
}

type chanMeta struct {
	ranges []int
	err    error
}

type chanMetaMsg struct {
	data []int
	meta chanMeta
}

// Variant 3: single channel with larger metadata

func receiverChan3(c <-chan chanMetaMsg) {
	for m := range c {
		for x := range m.data {
			sum += x
		}
		if m.meta.err != nil {
			return
		}
	}
}

func BenchmarkChan3(b *testing.B) {
	c := make(chan chanMetaMsg, chanBuf)
	go receiverChan3(c)
	d := make([]int, 1)
	for i := 0; i < b.N; i++ {
		d[0] = i
		c <- chanMetaMsg{d, chanMeta{}}
	}
	close(c)
}

// Variant 3a: single channel with larger metadata + ticker channel

func receiverChan3a(c <-chan chanMetaMsg, t <-chan time.Time) {
	for {
		select {
		case m := <-c:
			for x := range m.data {
				sum += x
			}
			if m.meta.err != nil {
				return
			}
		case <-t:
		}
	}
}

func BenchmarkChan3a(b *testing.B) {
	t := time.Tick(100 * time.Microsecond)
	c := make(chan chanMetaMsg, chanBuf)
	go receiverChan3a(c, t)
	d := make([]int, 1)
	for i := 0; i < b.N; i++ {
		d[0] = i
		c <- chanMetaMsg{d, chanMeta{}}
	}
	close(c)
}
