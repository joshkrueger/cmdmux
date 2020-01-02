package cmdmux

import (
	"testing"
)

func TestDefaultTokenizer(t *testing.T) {

	input := " hello this is a test"
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

func TestArgs(t *testing.T) {
	args := NewArgs()

	args.Set("foo", "bar")
	args.SetMeta("foo", "not-bar")

	if args.Get("foo") != "bar" {
		t.Error("Args didn't return expected value")
	}

	md := make(Metadata)
	md["foo"] = "not-bar"
	md["bar"] = "not-foo"

	args.SetMetadata(md)

	if args.GetMeta("foo") != "not-bar" {
		t.Error("Args metadata didn't return expected value")
	}

	if args.GetMeta("bar") != "not-foo" {
		t.Error("Args metadata didn't return expected value")
	}
}
