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

func TestCommitStore(t *testing.T) {
	for _, versionStr := range []string{ccipdata.V1_0_0, ccipdata.V1_2_0} {
		lggr := logger.TestLogger(t)
		addr := cciptypes.Address(utils.RandomAddress().String())
		lp := mocks2.NewLogPoller(t)

		feeEstimatorConfig := ccipdatamocks.NewFeeEstimatorConfigReader(t)

		lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
		versionFinder := newMockVersionFinder(ccipconfig.CommitStore, *semver.MustParse(versionStr), nil)
		_, err := NewCommitStoreReader(lggr, versionFinder, addr, nil, lp, feeEstimatorConfig)
		assert.NoError(t, err)

		expFilterName := logpoller.FilterName(v1_0_0.EXEC_REPORT_ACCEPTS, addr)
		lp.On("UnregisterFilter", mock.Anything, expFilterName).Return(nil)
		err = CloseCommitStoreReader(lggr, versionFinder, addr, nil, lp, feeEstimatorConfig)
		assert.NoError(t, err)
	}
}
