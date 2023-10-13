package ccipdata_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

func TestUSDCReaderFilters(t *testing.T) {
	assertFilterRegistration(t, new(lpmocks.LogPoller), func(lp *lpmocks.LogPoller, addr common.Address) ccipdata.Closer {
		c, err := ccipdata.NewUSDCReader(logger.TestLogger(t), addr, lp)
		require.NoError(t, err)
		return c
	}, 1)
}
