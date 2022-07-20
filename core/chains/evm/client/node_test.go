package client_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
)

func Test_NodeWrapError(t *testing.T) {
	t.Run("handles nil errors", func(t *testing.T) {
		err := evmclient.Wrap(nil, "foo")
		assert.NoError(t, err)
	})

	t.Run("adds extra info to context deadline exceeded errors", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()

		err := ctx.Err()

		err = evmclient.Wrap(err, "foo")

		assert.EqualError(t, err, "foo call failed: remote eth node timed out: context deadline exceeded")
	})
}
