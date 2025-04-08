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
	mocks2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

func TestPriceRegistry(t *testing.T) {
	ctx := testutils.Context(t)

	for _, versionStr := range []string{ccipdata.V1_2_0} {
		lggr := logger.Test(t)
		addr := cciptypes.Address(utils.RandomAddress().String())
		lp := mocks2.NewLogPoller(t)

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
