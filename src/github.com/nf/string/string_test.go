package string

import "testing"

func Test(t *testing.T) {
	var tests = []struct {
		s, want string
	}{
		{"",""},
		{"hello", "hello"},
	}
	for _, c := range tests {
		got := Print()
		if got != c.want{
			t.Error("Error")
		}
	}
}
