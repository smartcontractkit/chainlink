package eth_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/tidwall/gjson"

	eth "github.com/smartcontractkit/chainlink/core/chains/evm/eth"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NodeWrapError(t *testing.T) {
	t.Run("handles nil errors", func(t *testing.T) {
		err := eth.Wrap(nil, "foo")
		assert.NoError(t, err)
	})

	t.Run("adds extra info to context deadline exceeded errors", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()

		err := ctx.Err()

		err = eth.Wrap(err, "foo")

		assert.EqualError(t, err, "foo call failed: remote eth node timed out: context deadline exceeded")
	})
}

func Test_NodeStateTransitions(t *testing.T) {
	nInvalid := eth.NewNode(logger.TestLogger(t), *cltest.MustParseURL(t, "ws://example.invalid"), nil, "test node")
	wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (string, string) {
		return "", ""
	})

	nValid := eth.NewNode(logger.TestLogger(t), *cltest.MustParseURL(t, wsURL), nil, "test node")

	assert.Equal(t, eth.NodeStateUndialed, nInvalid.State())
	assert.Equal(t, eth.NodeStateUndialed, nValid.State())

	var err error

	t.Run("Verify before Dial", func(t *testing.T) {
		err = nInvalid.Verify(context.Background(), &cltest.FixtureChainID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot verify undialed node")
		err = nValid.Verify(context.Background(), &cltest.FixtureChainID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot verify undialed node")
	})

	assert.Equal(t, eth.NodeStateUndialed, nInvalid.State())
	assert.Equal(t, eth.NodeStateUndialed, nValid.State())

	t.Run("Dial state changes", func(t *testing.T) {
		err = nInvalid.Dial(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error while dialing websocket")

		assert.Equal(t, eth.NodeStateDead, nInvalid.State())

		// make sure that verifying dead node doesn't crash
		err = nInvalid.Verify(context.Background(), &cltest.FixtureChainID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot verify dead node")

		assert.Equal(t, eth.NodeStateDead, nInvalid.State())

		err = nValid.Dial(context.Background())
		require.NoError(t, err)

		assert.Equal(t, eth.NodeStateDialed, nValid.State())
	})

	t.Run("Verify after dial", func(t *testing.T) {
		err = nValid.Verify(context.Background(), big.NewInt(99))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "websocket rpc ChainID doesn't match local chain ID: RPC ID=0, local ID=99")

		assert.Equal(t, eth.NodeStateInvalidChainID, nValid.State())

		err = nValid.Verify(context.Background(), &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Equal(t, eth.NodeStateAlive, nValid.State())
	})

	t.Run("Close state changes", func(t *testing.T) {
		nInvalid.Close()
		assert.Equal(t, eth.NodeStateClosed, nInvalid.State())
		nValid.Close()
		assert.Equal(t, eth.NodeStateClosed, nValid.State())
	})
}
