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
)

func TestOnRamp(t *testing.T) {
	ctx := t.Context()
	for _, versionStr := range []string{ccipdata.V1_2_0, ccipdata.V1_5_0} {
		lggr := logger.Test(t)
		addr := cciptypes.Address(utils.RandomAddress().String())
		lp := lpmocks.NewLogPoller(t)

		sourceSelector := uint64(1000)
		destSelector := uint64(2000)

		expFilterNames := []string{
			logpoller.FilterName(ccipdata.COMMIT_CCIP_SENDS, addr),
			logpoller.FilterName(ccipdata.CONFIG_CHANGED, addr),
		}
		versionFinder := newMockVersionFinder(ccipconfig.EVM2EVMOnRamp, *semver.MustParse(versionStr), nil)

		lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil).Times(len(expFilterNames))
		_, err := NewOnRampReader(ctx, lggr, versionFinder, sourceSelector, destSelector, addr, lp, nil)
		assert.NoError(t, err)

		for _, f := range expFilterNames {
			lp.On("UnregisterFilter", mock.Anything, f).Return(nil)
		}
		err = CloseOnRampReader(ctx, lggr, versionFinder, sourceSelector, destSelector, addr, lp, nil)
		assert.NoError(t, err)
	}
}
