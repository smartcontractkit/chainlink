package factory

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	lpmocks "github.com/smartcontractkit/chainlink/v2/common/logpoller/mocks"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
)

func TestCommitStore(t *testing.T) {
	ctx := t.Context()
	for _, versionStr := range []string{ccipdata.V1_2_0} {
		lggr := logger.Test(t)
		addr := cciptypes.Address(utils.RandomAddress().String())
		lp := lpmocks.NewLogPoller(t)

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
