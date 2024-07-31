package evmliquiditymanager

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func TestNewBaseLiquidityManagerFactory(t *testing.T) {
	lggr := logger.TestLogger(t)
	net1 := models.NetworkSelector(1)
	addr1 := models.Address(utils.RandomAddress())
	net2 := models.NetworkSelector(2)

	t.Run("base constructor test", func(t *testing.T) {
		client1 := mocks.NewClient(t)
		client2 := mocks.NewClient(t)
		lmf := NewBaseLiquidityManagerFactory(lggr, WithEvmDep(net1, client1), WithEvmDep(net2, client2))
		assert.Len(t, lmf.evmDeps, 2)
	})

	t.Run("wrong cached type", func(t *testing.T) {
		lmf := NewBaseLiquidityManagerFactory(lggr)
		lmf.cache.Store(lmf.cacheKey(net1, addr1), 1234)
		_, err := lmf.GetLiquidityManager(net1, addr1)
		assert.Equal(t, ErrInternalCacheIssue, err)
	})

	t.Run("get from cache", func(t *testing.T) {
		lmf := NewBaseLiquidityManagerFactory(lggr)
		evmRb := &EvmLiquidityManager{}
		var rb LiquidityManager = evmRb
		lmf.cache.Store(lmf.cacheKey(net1, addr1), rb)
		_, err := lmf.GetLiquidityManager(net1, addr1)
		assert.NoError(t, err)
	})

	t.Run("cache key", func(t *testing.T) {
		lmf := NewBaseLiquidityManagerFactory(lggr)
		net1 := models.NetworkSelector(1)
		addr1 := models.Address(common.HexToAddress("0x000000000000000000000000000000000000dEaD"))
		assert.Equal(t, "rebalancer-1-0x000000000000000000000000000000000000dEaD", lmf.cacheKey(net1, addr1))
	})
}
