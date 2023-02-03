package client

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/stretchr/testify/assert"
)

func Test_verifyLoop(t *testing.T) {
	t.Run("skips check and sets to online if chain ID is zero", func(t *testing.T) {
		s := &sendOnlyNode{chainID: big.NewInt(0), log: logger.TestLogger(t)}
		s.wg.Add(1)

		// exits immediately, does not panic or hang
		s.verifyLoop()

		s.wg.Wait()

		assert.Equal(t, NodeStateAlive, s.State())
	})
}
