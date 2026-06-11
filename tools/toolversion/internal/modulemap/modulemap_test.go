package modulemap

import "testing"

func TestModulePath(t *testing.T) {
	t.Parallel()
	mod, err := ModulePath("mockery")
	if err != nil {
		t.Fatal(err)
	}
	if mod != "github.com/vektra/mockery/v2" {
		t.Fatalf("got %q", mod)
	}
	_, err = ModulePath("protoc")
	if err == nil {
		t.Fatal("expected error for protoc")
	}
}
