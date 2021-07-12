package offchainreporting_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/stretchr/testify/assert"
)

func Test_BlockTranslator(t *testing.T) {
	ethClient := new(mocks.Client)
	ctx := context.Background()

	t.Run("for L1 chains, returns the block changed argument", func(t *testing.T) {
		chain := chains.ChainFromID(big.NewInt(1))

		bt := offchainreporting.NewBlockTranslator(chain, ethClient)

		from, to := bt.NumberToQueryRange(ctx, 42)

		assert.Equal(t, big.NewInt(42), from)
		assert.Equal(t, big.NewInt(42), to)
	})

	t.Run("for optimism, returns an initial block number and nil", func(t *testing.T) {
		bt := offchainreporting.NewBlockTranslator(chains.OptimismMainnet, ethClient)
		from, to := bt.NumberToQueryRange(ctx, 42)
		assert.Equal(t, big.NewInt(0), from)
		assert.Equal(t, (*big.Int)(nil), to)

		bt = offchainreporting.NewBlockTranslator(chains.OptimismKovan, ethClient)
		from, to = bt.NumberToQueryRange(ctx, 42)
		assert.Equal(t, big.NewInt(0), from)
		assert.Equal(t, (*big.Int)(nil), to)
	})

	t.Run("for arbitrum, returns the ArbitrumBlockTranslator", func(t *testing.T) {
		bt := offchainreporting.NewBlockTranslator(chains.ArbitrumMainnet, ethClient)
		assert.IsType(t, &offchainreporting.ArbitrumBlockTranslator{}, bt)

		bt = offchainreporting.NewBlockTranslator(chains.ArbitrumRinkeby, ethClient)
		assert.IsType(t, &offchainreporting.ArbitrumBlockTranslator{}, bt)
	})

	ethClient.AssertExpectations(t)
}
