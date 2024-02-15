package factory

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	mocks2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
)

func TestOffRamp(t *testing.T) {
	for _, versionStr := range []string{ccipdata.V1_0_0, ccipdata.V1_2_0} {
		lggr := logger.TestLogger(t)
		addr := cciptypes.Address(utils.RandomAddress().String())
		lp := mocks2.NewLogPoller(t)

		expFilterNames := []string{
			logpoller.FilterName(v1_0_0.EXEC_EXECUTION_STATE_CHANGES, addr),
			logpoller.FilterName(v1_0_0.EXEC_TOKEN_POOL_ADDED, addr),
			logpoller.FilterName(v1_0_0.EXEC_TOKEN_POOL_REMOVED, addr),
		}
		versionFinder := newMockVersionFinder(ccipconfig.EVM2EVMOffRamp, *semver.MustParse(versionStr), nil)

		lp.On("RegisterFilter", mock.Anything).Return(nil).Times(len(expFilterNames))
		_, err := NewOffRampReader(lggr, versionFinder, addr, nil, lp, nil, true)
		assert.NoError(t, err)

		for _, f := range expFilterNames {
			lp.On("UnregisterFilter", f).Return(nil)
		}
		err = CloseOffRampReader(lggr, versionFinder, addr, nil, lp, nil)
		assert.NoError(t, err)
	}
}
