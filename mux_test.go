package cmdmux

import (
	"testing"
)

func TestDefaultTokenizer(t *testing.T) {

	input := "hello this is a test"
	expected := []string{"hello", "this", "is", "a", "test"}

	output := DefaultTokenizer(input)
	if !pathMatch(expected, output) {
		t.Fatalf("token paths did not match. Expected: %v Actual: %v", expected, output)
	}
}

func pathMatch(l, r []string) bool {
	if len(l) != len(r) {
		return false
	}
	for i, v := range l {
		if v != r[i] {
			return false
		}
	}
	return true
}
