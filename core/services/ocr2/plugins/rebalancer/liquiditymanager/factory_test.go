package liquiditymanager

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

func TestNewBaseLiquidityManagerFactory(t *testing.T) {
	lp1 := mocks.NewLogPoller(t)
	lp2 := mocks.NewLogPoller(t)
	lmf := NewBaseRebalancerFactory(
		WithEvmDep(models.NetworkSelector(1), lp1, nil),
		WithEvmDep(models.NetworkSelector(2), lp2, nil),
	)
	assert.Len(t, lmf.evmDeps, 2)
}
