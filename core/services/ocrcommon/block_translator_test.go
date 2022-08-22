package ocrcommon_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"

	"github.com/stretchr/testify/assert"
)

func Test_BlockTranslator(t *testing.T) {
	t.Parallel()

	ethClient := evmtest.NewEthClientMock(t)
	ctx := testutils.Context(t)
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
}
