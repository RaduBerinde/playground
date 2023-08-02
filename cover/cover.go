package cover

func Foo(a, b int) int {
	if a < b {
		if a == 1 {
			return 10
		}
		return a * b
	}
	if a == b {
		return 100
	}
	if b == 1 {
		return 1000
	}
	return a + b
}
