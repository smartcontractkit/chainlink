package factory

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	mocks2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
)

func TestOffRamp(t *testing.T) {
	for _, versionStr := range []string{ccipdata.V1_0_0, ccipdata.V1_2_0} {
		lggr := logger.TestLogger(t)
		addr := cciptypes.Address(utils.RandomAddress().String())
		lp := mocks2.NewLogPoller(t)

		feeEstimatorConfig := ccipdatamocks.NewFeeEstimatorConfigReader(t)

		expFilterNames := []string{
			logpoller.FilterName(v1_0_0.EXEC_EXECUTION_STATE_CHANGES, addr),
			logpoller.FilterName(v1_0_0.EXEC_TOKEN_POOL_ADDED, addr),
			logpoller.FilterName(v1_0_0.EXEC_TOKEN_POOL_REMOVED, addr),
		}
		versionFinder := newMockVersionFinder(ccipconfig.EVM2EVMOffRamp, *semver.MustParse(versionStr), nil)

		lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil).Times(len(expFilterNames))
		_, err := NewOffRampReader(lggr, versionFinder, addr, nil, lp, nil, nil, true, feeEstimatorConfig)
		assert.NoError(t, err)

		for _, f := range expFilterNames {
			lp.On("UnregisterFilter", mock.Anything, f).Return(nil)
		}
		err = CloseOffRampReader(lggr, versionFinder, addr, nil, lp, nil, nil, feeEstimatorConfig)
		assert.NoError(t, err)
	}
}
