package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStaticOnRamp(t *testing.T) {
	t.Parallel()

	// static test implementation is self consistent
	ctx := context.Background()
	assert.NoError(t, OnRamp.Evaluate(ctx, OnRamp))

	// error when the test implementation is evaluates something that differs from the static implementation
	botched := OnRamp
	botched.addressResponse = "not the right address"
	err := OnRamp.Evaluate(ctx, botched)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not the right address")
}
