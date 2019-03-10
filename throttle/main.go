package main

import (
	"fmt"
	"time"
)

var m = make(map[int]int)

const inner = 400000
const outer = 100

func workload() {
	for i := 0; i < inner; i++ {
		m[i] = i
	}
}

// Sample output (for 400000 x 100):
//
// Unthrottled: 2.625012595s (sum: 2.624992595s)
// With throttling: 1m5.177411266s (sum: 6.467121915s)   <-- should be 26s

func main() {
	before := time.Now()
	var sum time.Duration
	for i := 0; i < outer; i++ {
		iterStart := time.Now()
		workload()
		iterEnd := time.Now()
		sum += iterEnd.Sub(iterStart)
	}
	fmt.Printf("Unthrottled: %s (sum: %s)\n", time.Now().Sub(before), sum)

	sum = 0
	before = time.Now()
	for i := 0; i < outer; i++ {
		iterStart := time.Now()
		workload()
		iterEnd := time.Now()
		sum += iterEnd.Sub(iterStart)
		s := iterEnd.Sub(iterStart) * 9
		time.Sleep(s)
		//fmt.Printf("Sleep(%s), actually slept for %s\n", s, time.Now().Sub(iterEnd))
	}
	fmt.Printf("With throttling: %s (sum: %s)\n", time.Now().Sub(before), sum)
}
