package testhelp

import "testing"

func ExpectTrue(t *testing.T, value bool, message string) {
	t.Helper()
	if !value {
		t.Fatalf("expected true: %s", message)
	}
}

func ExpectFalse(t *testing.T, value bool, message string) {
	t.Helper()
	if value {
		t.Fatalf("expected false: %s", message)
	}
}
