// Output:
// BenchmarkReflectValue-4	100000000	        20.2 ns/op
// BenchmarkEquality-4    	200000000	         9.67 ns/op

package main

import (
	"reflect"
	"testing"
)

type SomeInterface interface {
	Cool() int
}

type SomeType struct {
	x int
}

func (t *SomeType) Cool() int { return t.x }

func BenchmarkReflectValue(b *testing.B) {
	var arr [1000]SomeInterface

	for i := 0; i < 1000; i++ {
		arr[i] = &SomeType{i}
	}

	test := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if reflect.ValueOf(arr[i%1000]) == reflect.ValueOf(arr[(i+1)%1000]) {
			test++
		}
	}
	if test != 0 {
		b.Errorf("Equality!")
	}
}

func BenchmarkEquality(b *testing.B) {
	var arr [1000]SomeInterface

	for i := 0; i < 1000; i++ {
		arr[i] = &SomeType{i}
	}

	test := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if arr[i%1000] == arr[(i+1)%1000] {
			test++
		}
	}
	if test != 0 {
		b.Errorf("Equality!")
	}
}
