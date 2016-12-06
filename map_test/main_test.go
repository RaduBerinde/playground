package main

import (
	"fmt"
	"testing"
)

func BenchmarkSlice(b *testing.B) {
	for _, n := range []int{1, 2, 3, 4, 8, 16, 32} {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m := make([]int, n)
				for j := 0; j < n; j++ {
					m[j] = j
				}
			}
		})
	}
}

func BenchmarkMap(b *testing.B) {
	for _, n := range []int{1, 2, 3, 4, 8, 16, 32} {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m := make(map[int]int)
				for j := 0; j < n; j++ {
					m[j] = j
				}
			}
		})
	}
}
