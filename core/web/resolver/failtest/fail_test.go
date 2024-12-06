package failtest

import "testing"

func TestImmediateFail(t *testing.T) {
	t.Fatalf("This test fails immediately.")
}
