package factory

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	ccipdata2 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/ccipdata"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/version"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	lpmocks "github.com/smartcontractkit/chainlink/v2/common/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

func TestPriceRegistry(t *testing.T) {
	ctx := testutils.Context(t)

	for _, versionStr := range []string{ccipdata2.V1_2_0} {
		lggr := logger.Test(t)
		addr := cciptypes.Address(utils.RandomAddress().String())
		lp := lpmocks.NewLogPoller(t)

		expFilterNames := []string{
			logpoller.FilterName(ccipdata.COMMIT_PRICE_UPDATES, addr),
			logpoller.FilterName(ccipdata.FEE_TOKEN_ADDED, addr),
			logpoller.FilterName(ccipdata.FEE_TOKEN_REMOVED, addr),
		}
		versionFinder := newMockVersionFinder(ccipconfig.PriceRegistry, *semver.MustParse(versionStr), nil)

		lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil).Times(len(expFilterNames))
		_, err := NewPriceRegistryReader(ctx, lggr, versionFinder, addr, lp, nil)
		assert.NoError(t, err)

		for _, f := range expFilterNames {
			lp.On("UnregisterFilter", mock.Anything, f).Return(nil)
		}
		err = ClosePriceRegistryReader(ctx, lggr, versionFinder, addr, lp, nil)
		assert.NoError(t, err)
	}
}
