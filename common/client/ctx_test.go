package client

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestContext(t *testing.T) {
	ctx := tests.Context(t)
	assert.False(t, CtxIsHeathCheckRequest(ctx), "expected false for test context")
	ctx = CtxAddHealthCheckFlag(ctx)
	assert.True(t, CtxIsHeathCheckRequest(ctx), "expected context to contain the healthcheck flag")
}
