package main

var f func()

func cucu() {
	x := 5
	f = func() { x++ }
}

func main() {
	cucu()
	f()
	f()
}
