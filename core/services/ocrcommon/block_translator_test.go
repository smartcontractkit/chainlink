package ocrcommon_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/stretchr/testify/assert"
)

func Test_BlockTranslator(t *testing.T) {
	ethClient := cltest.NewEthClientMock(t)
	ctx := context.Background()
	lggr := logger.TestLogger(t)

	t.Run("for L1 chains, returns the block changed argument", func(t *testing.T) {
		bt := ocrcommon.NewBlockTranslator(evmtest.ChainEthMainnet(t), ethClient, lggr)

		from, to := bt.NumberToQueryRange(ctx, 42)

		assert.Equal(t, big.NewInt(42), from)
		assert.Equal(t, big.NewInt(42), to)
	})

	t.Run("for optimism, uses the default translator", func(t *testing.T) {
		bt := ocrcommon.NewBlockTranslator(evmtest.ChainOptimismMainnet(t), ethClient, lggr)
		from, to := bt.NumberToQueryRange(ctx, 42)
		assert.Equal(t, big.NewInt(42), from)
		assert.Equal(t, big.NewInt(42), to)

		bt = ocrcommon.NewBlockTranslator(evmtest.ChainOptimismKovan(t), ethClient, lggr)
		from, to = bt.NumberToQueryRange(ctx, 42)
		assert.Equal(t, big.NewInt(42), from)
		assert.Equal(t, big.NewInt(42), to)
	})

	t.Run("for arbitrum, returns the ArbitrumBlockTranslator", func(t *testing.T) {
		bt := ocrcommon.NewBlockTranslator(evmtest.ChainArbitrumMainnet(t), ethClient, lggr)
		assert.IsType(t, &ocrcommon.ArbitrumBlockTranslator{}, bt)

		bt = ocrcommon.NewBlockTranslator(evmtest.ChainArbitrumRinkeby(t), ethClient, lggr)
		assert.IsType(t, &ocrcommon.ArbitrumBlockTranslator{}, bt)
	})

	ethClient.AssertExpectations(t)
}
