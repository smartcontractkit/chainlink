package logger

import "testing"

func TestWeakRegistryCloseIdempotent(t *testing.T) {
	r := NewWeakRegistry[int]()

	r.Close()
	assertNotPanics(t, func() { r.Close() })
}

func assertNotPanics(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if p := recover(); p != nil {
			t.Fatalf("unexpected panic: %v", p)
		}
	}()
	fn()
}
