package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStaticOffRamp(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	assert.NoError(t, OffRamp.Evaluate(ctx, OffRamp))
}
