package offchainreporting_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/stretchr/testify/assert"
)

func Test_BlockTranslator(t *testing.T) {
	ethClient := cltest.NewEthClientMock(t)
	ctx := context.Background()

	t.Run("for L1 chains, returns the block changed argument", func(t *testing.T) {
		bt := offchainreporting.NewBlockTranslator(evmtest.ChainEthMainnet(t), ethClient, logger.Default)

		from, to := bt.NumberToQueryRange(ctx, 42)

		assert.Equal(t, big.NewInt(42), from)
		assert.Equal(t, big.NewInt(42), to)
	})

	t.Run("for optimism, returns an initial block number and nil", func(t *testing.T) {
		bt := offchainreporting.NewBlockTranslator(evmtest.ChainOptimismMainnet(t), ethClient, logger.Default)
		from, to := bt.NumberToQueryRange(ctx, 42)
		assert.Equal(t, big.NewInt(0), from)
		assert.Equal(t, (*big.Int)(nil), to)

		bt = offchainreporting.NewBlockTranslator(evmtest.ChainOptimismKovan(t), ethClient, logger.Default)
		from, to = bt.NumberToQueryRange(ctx, 42)
		assert.Equal(t, big.NewInt(0), from)
		assert.Equal(t, (*big.Int)(nil), to)
	})

	t.Run("for arbitrum, returns the ArbitrumBlockTranslator", func(t *testing.T) {
		bt := offchainreporting.NewBlockTranslator(evmtest.ChainArbitrumMainnet(t), ethClient, logger.Default)
		assert.IsType(t, &offchainreporting.ArbitrumBlockTranslator{}, bt)

		bt = offchainreporting.NewBlockTranslator(evmtest.ChainArbitrumRinkeby(t), ethClient, logger.Default)
		assert.IsType(t, &offchainreporting.ArbitrumBlockTranslator{}, bt)
	})

	ethClient.AssertExpectations(t)
}
