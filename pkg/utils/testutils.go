package utils

import (
	"context"
	"testing"
)

func Context(t *testing.T) context.Context {
	ctx := context.Background()
	var cancel func()

	if d, ok := t.Deadline(); ok {
		ctx, cancel = context.WithDeadline(ctx, d)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	t.Cleanup(cancel)
	return ctx
}
