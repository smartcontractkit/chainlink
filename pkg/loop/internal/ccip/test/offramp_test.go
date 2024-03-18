package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStaticOffRamp(t *testing.T) {
	t.Parallel()

	// static test implementation is self consistent
	ctx := context.Background()
	assert.NoError(t, OffRamp.Evaluate(ctx, OffRamp))

	// error when the test implementation is evaluates something that differs from the static implementation
	botched := OffRamp
	botched.addressResponse = "oops"
	err := OffRamp.Evaluate(ctx, botched)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "oops")
}
