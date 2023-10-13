package ccipdata

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestPriceRegistryFilters(t *testing.T) {
	cl := mocks.NewClient(t)

	assertFilterRegistration(t, new(lpmocks.LogPoller), func(lp *lpmocks.LogPoller, addr common.Address) Closer {
		c, err := NewPriceRegistryV1_0_0(logger.TestLogger(t), addr, lp, cl)
		require.NoError(t, err)
		return c
	}, 3)
}
