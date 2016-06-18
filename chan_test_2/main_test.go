package main

import "testing"

const chanBuf = 10

type chanMsg struct {
	data []int
	err  error
}

var sum int

func push(c chan<- chanMsg) {
	c <- chanMsg{data: []int{1, 2, 3}}
	close(c)
}

func BenchmarkFoo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := make(chan chanMsg, 16)
		go push(c)
	loop:
		for {
			select {
			case msg, ok := <-c:
				if !ok {
					break loop
				}
				if msg.data[0] != 1 {
					panic(msg.data)
				}
			}
		}
	}
}
