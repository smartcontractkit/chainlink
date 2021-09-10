package ocrcommon_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"

	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/stretchr/testify/assert"
)

func Test_BlockTranslator(t *testing.T) {
	ethClient := cltest.NewEthClientMock(t)
	ctx := context.Background()

	t.Run("for L1 chains, returns the block changed argument", func(t *testing.T) {
		chain := evmtest.ChainEthMainnet()

		bt := ocrcommon.NewBlockTranslator(chain, ethClient)

		from, to := bt.NumberToQueryRange(ctx, 42)

		assert.Equal(t, big.NewInt(42), from)
		assert.Equal(t, big.NewInt(42), to)
	})

	t.Run("for optimism, returns an initial block number and nil", func(t *testing.T) {
		bt := ocrcommon.NewBlockTranslator(evmtest.ChainOptimismMainnet(), ethClient)
		from, to := bt.NumberToQueryRange(ctx, 42)
		assert.Equal(t, big.NewInt(0), from)
		assert.Equal(t, (*big.Int)(nil), to)

		bt = ocrcommon.NewBlockTranslator(evmtest.ChainOptimismKovan(), ethClient)
		from, to = bt.NumberToQueryRange(ctx, 42)
		assert.Equal(t, big.NewInt(0), from)
		assert.Equal(t, (*big.Int)(nil), to)
	})

	t.Run("for arbitrum, returns the ArbitrumBlockTranslator", func(t *testing.T) {
		bt := ocrcommon.NewBlockTranslator(evmtest.ChainArbitrumMainnet(), ethClient)
		assert.IsType(t, &ocrcommon.ArbitrumBlockTranslator{}, bt)

		bt = ocrcommon.NewBlockTranslator(evmtest.ChainArbitrumRinkeby(), ethClient)
		assert.IsType(t, &ocrcommon.ArbitrumBlockTranslator{}, bt)
	})

	ethClient.AssertExpectations(t)
}
