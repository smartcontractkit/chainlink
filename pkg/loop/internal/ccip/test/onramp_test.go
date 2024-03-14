package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStaticOnRamp(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	assert.NoError(t, OnRamp.Evaluate(ctx, OnRamp))
}
