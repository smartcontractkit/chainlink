package factory

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	mocks2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
)

func TestCommitStore(t *testing.T) {
	ctx := tests.Context(t)
	for _, versionStr := range []string{ccipdata.V1_2_0} {
		lggr := logger.Test(t)
		addr := cciptypes.Address(utils.RandomAddress().String())
		lp := mocks2.NewLogPoller(t)

		feeEstimatorConfig := ccipdatamocks.NewFeeEstimatorConfigReader(t)

		lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
		versionFinder := newMockVersionFinder(ccipconfig.CommitStore, *semver.MustParse(versionStr), nil)
		_, err := NewCommitStoreReader(ctx, lggr, versionFinder, addr, nil, lp, feeEstimatorConfig)
		assert.NoError(t, err)

		expFilterName := logpoller.FilterName(v1_2_0.ExecReportAccepts, addr)
		lp.On("UnregisterFilter", mock.Anything, expFilterName).Return(nil)
		err = CloseCommitStoreReader(ctx, lggr, versionFinder, addr, nil, lp, feeEstimatorConfig)
		assert.NoError(t, err)
	}
}
