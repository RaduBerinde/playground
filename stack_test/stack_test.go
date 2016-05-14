package main

import (
	"bytes"
	"fmt"
	"testing"
)

type Flags struct {
	numBranches, maxDepth      int
	a2, a3, a4, a5, a6, a7, a8 int
}

func Func1(buf *bytes.Buffer, depth int, f *Flags) {
	fmt.Fprintf(buf, "%d", depth+f.a2+f.a3+f.a4+f.a5+f.a6+f.a7+f.a8)
	if depth < f.maxDepth {
		for i := 0; i < f.numBranches; i++ {
			Func1(buf, depth+1, f)
		}
	}
}

func Func2(buf *bytes.Buffer, depth int, f Flags) {
	fmt.Fprintf(buf, "%d", depth+f.a2+f.a3+f.a4+f.a5+f.a6+f.a7+f.a8)
	if depth < f.maxDepth {
		for i := 0; i < f.numBranches; i++ {
			Func2(buf, depth+1, f)
		}
	}
}

func BenchmarkFunc1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := Flags{numBranches: 2, maxDepth: 5, a2: i}
		var buf bytes.Buffer
		Func1(&buf, 0, &f)
	}
}

func BenchmarkFunc2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := Flags{numBranches: 2, maxDepth: 5, a2: i}
		var buf bytes.Buffer
		Func2(&buf, 0, f)
	}
}
