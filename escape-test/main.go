package main

import (
	"fmt"
	"unsafe"
)

var f func()

func cucu() {
	x := 5
	f = func() { x++ }
}

func main() {
	cucu()
	f()
	f()
	a := 5
	x := unsafe.Pointer(&a)
	fmt.Printf("%d\n", x)

	b := [2]int{1, 2}
	c := b[:1]
	c[1] = 3
	y := unsafe.Pointer(&c)
	fmt.Printf("%d\n", b, y)
}
