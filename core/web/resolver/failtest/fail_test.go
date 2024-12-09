package failtest

import "testing"

func TestImmediateFail(t *testing.T) {
	t.Fatalf("This test fails immediately.")
}

func TestPassing(t *testing.T) {
	t.Log("This test passes.")
}
