package client_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

type evmRPC = commonclient.RPC[
	*big.Int,
	evmtypes.Nonce,
	common.Address,
	common.Hash,
	*types.Transaction,
	common.Hash,
	types.Log,
	ethereum.FilterQuery,
	*evmtypes.Receipt,
	*assets.Wei,
	*evmtypes.Head,
]

func Test_NodeWrapError(t *testing.T) {
	t.Parallel()

	t.Run("handles nil errors", func(t *testing.T) {
		err := evmclient.Wrap(nil, "foo")
		assert.NoError(t, err)
	})

	t.Run("adds extra info to context deadline exceeded errors", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(testutils.Context(t), 0)
		defer cancel()

		err := ctx.Err()

		err = evmclient.Wrap(err, "foo")

		assert.EqualError(t, err, "foo call failed: remote eth node timed out: context deadline exceeded")
	})
}
