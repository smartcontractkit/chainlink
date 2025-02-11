package factory

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	mocks2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
)

func TestOffRamp(t *testing.T) {
	ctx := tests.Context(t)
	for _, versionStr := range []string{ccipdata.V1_2_0} {
		lggr := logger.Test(t)
		addr := cciptypes.Address(utils.RandomAddress().String())
		lp := mocks2.NewLogPoller(t)

		feeEstimatorConfig := ccipdatamocks.NewFeeEstimatorConfigReader(t)

		expFilterNames := []string{
			logpoller.FilterName(v1_2_0.ExecExecutionStateChanges, addr),
			logpoller.FilterName(v1_2_0.ExecTokenPoolAdded, addr),
			logpoller.FilterName(v1_2_0.ExecTokenPoolRemoved, addr),
		}
		versionFinder := newMockVersionFinder(ccipconfig.EVM2EVMOffRamp, *semver.MustParse(versionStr), nil)

		lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil).Times(len(expFilterNames))
		_, err := NewOffRampReader(ctx, lggr, versionFinder, addr, nil, lp, nil, nil, true, feeEstimatorConfig)
		assert.NoError(t, err)

		for _, f := range expFilterNames {
			lp.On("UnregisterFilter", mock.Anything, f).Return(nil)
		}
		err = CloseOffRampReader(ctx, lggr, versionFinder, addr, nil, lp, nil, nil, feeEstimatorConfig)
		assert.NoError(t, err)
	}
}
