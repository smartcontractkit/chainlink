package offchainreporting_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/stretchr/testify/assert"
)

func Test_BlockTranslator(t *testing.T) {
	t.Run("for L1 chains, returns the block changed argument", func(t *testing.T) {
		chain := chains.ChainFromID(big.NewInt(1))

		bt := offchainreporting.NewBlockTranslator(chain)
		bt.Start()
		defer bt.Close()

		from, to := bt.NumberToQueryRange(42)

		assert.Equal(t, big.NewInt(42), from)
		assert.Equal(t, big.NewInt(42), to)
	})

	t.Run("for L2 chains, returns an initial block number and nil", func(t *testing.T) {
		bt := offchainreporting.NewBlockTranslator(chains.ArbitrumMainnet)
		bt.Start()
		defer bt.Close()
		from, to := bt.NumberToQueryRange(42)
		assert.Equal(t, big.NewInt(12525700), from)
		assert.Equal(t, (*big.Int)(nil), to)

		bt = offchainreporting.NewBlockTranslator(chains.ArbitrumRinkeby)
		bt.Start()
		defer bt.Close()
		from, to = bt.NumberToQueryRange(42)
		assert.Equal(t, big.NewInt(8700589), from)
		assert.Equal(t, (*big.Int)(nil), to)

		bt = offchainreporting.NewBlockTranslator(chains.OptimismMainnet)
		bt.Start()
		defer bt.Close()
		from, to = bt.NumberToQueryRange(42)
		assert.Equal(t, big.NewInt(0), from)
		assert.Equal(t, (*big.Int)(nil), to)

		bt = offchainreporting.NewBlockTranslator(chains.OptimismKovan)
		bt.Start()
		defer bt.Close()
		from, to = bt.NumberToQueryRange(42)
		assert.Equal(t, big.NewInt(0), from)
		assert.Equal(t, (*big.Int)(nil), to)
	})
}
