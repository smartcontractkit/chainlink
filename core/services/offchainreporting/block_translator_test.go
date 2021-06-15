package offchainreporting_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/stretchr/testify/assert"
)

func Test_BlockTranslator(t *testing.T) {
	t.Run("for non-arbitrum chains, returns the block changed argument", func(t *testing.T) {
		chain := chains.ChainFromID(big.NewInt(1))

		bt := offchainreporting.NewBlockTranslator(chain)
		bt.Start()
		defer bt.Close()

		from, to := bt.NumberToQueryRange(42)

		assert.Equal(t, big.NewInt(42), from)
		assert.Equal(t, big.NewInt(42), to)
	})

	t.Run("for arbitrum chains, returns 0 and nil", func(t *testing.T) {
		chain := chains.ChainFromID(big.NewInt(42161))

		bt := offchainreporting.NewBlockTranslator(chain)
		bt.Start()
		defer bt.Close()

		from, to := bt.NumberToQueryRange(42)

		assert.Equal(t, big.NewInt(0), from)
		assert.Equal(t, (*big.Int)(nil), to)
	})
}
