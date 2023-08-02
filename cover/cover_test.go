package cover

import "testing"

func TestFoo(t *testing.T) {
	testCases := []struct{ a, b, exp int }{
		{a: 10, b: 100, exp: 1000},
		{a: 1, b: 100, exp: 10},
		{a: 5, b: 5, exp: 100},
		{a: 50, b: 5, exp: 55},
	}

	for _, tc := range testCases {
		if res := Foo(tc.a, tc.b); res != tc.exp {
			t.Errorf("%d %d: expected %d, got %d", tc.a, tc.b, tc.exp, res)
		}
	}
}
