package environment

import (
	"context"
	"testing"
)

func TestSetupOutputCloseIsIdempotent(t *testing.T) {
	out := &SetupOutput{}

	if err := out.Close(context.Background()); err != nil {
		t.Fatalf("expected first close to succeed: %v", err)
	}
	if err := out.Close(context.Background()); err != nil {
		t.Fatalf("expected second close to succeed: %v", err)
	}
}
