package chains_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/stretchr/testify/assert"
)

func Test_ChainFromID(t *testing.T) {
	t.Run("returns existing chain", func(t *testing.T) {
		c := chains.ChainFromID(big.NewInt(1))

		assert.Equal(t, big.NewInt(1), c.ID())
	})
	t.Run("falls back to generic chain if missing", func(t *testing.T) {
		c := chains.ChainFromID(big.NewInt(0))

		assert.Equal(t, big.NewInt(0), c.ID())
		assert.Equal(t, "", c.Config().LinkContractAddress)

		c2 := chains.ChainFromID(big.NewInt(0))

		assert.Equal(t, c, c2)

		c3 := chains.ChainFromID(big.NewInt(98765))

		assert.Equal(t, big.NewInt(98765), c3.ID())
		assert.Equal(t, "", c3.Config().LinkContractAddress)
	})
}
