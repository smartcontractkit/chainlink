package rebalancer

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	bridgemocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/bridge/mocks"
	discoverermocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/discoverer/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
	rebalancer_mocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/rebalancermocks"
)

type mockDeps struct {
	mockFactory           *rebalancer_mocks.Factory
	mockRebalancer        *rebalancer_mocks.Rebalancer
	mockDiscovererFactory *discoverermocks.Factory
	mockBridgeFactory     *bridgemocks.Factory
}

func newPlugin(t *testing.T) (*Plugin, mockDeps) {
	f := 10
	closeTimeout := 5 * time.Second
	rootNetwork := models.NetworkSelector(1)
	rootAddr := models.Address(utils.RandomAddress())

	lmFactory := rebalancer_mocks.NewFactory(t)
	rb := rebalancer_mocks.NewRebalancer(t)
	discovererFactory := discoverermocks.NewFactory(t)
	bridgeFactory := bridgemocks.NewFactory(t)
	return NewPlugin(f, closeTimeout, rootNetwork, rootAddr, lmFactory, discovererFactory, bridgeFactory, rb, logger.TestLogger(t)), mockDeps{
		mockFactory:           lmFactory,
		mockRebalancer:        rb,
		mockDiscovererFactory: discovererFactory,
		mockBridgeFactory:     bridgeFactory,
	}
}

func TestPluginQuery(t *testing.T) {
	p, _ := newPlugin(t)
	q, err := p.Query(context.Background(), ocr3types.OutcomeContext{})
	require.Empty(t, q, "query should always be empty")
	require.NoError(t, err)
}

func TestPluginObservation(t *testing.T) {
	// TODO: fix test
	// ctx := testutils.Context(t)
	// p, deps := newPlugin(t)

	// networks := p.rebalancerGraph.GetNetworks()
	// require.Len(t, networks, 0, "plugin should initially contain zero nodes in the graph")

	// obs, err := p.Observation(ctx, ocr3types.OutcomeContext{}, ocrtypes.Query{})
	// require.NoError(t, err)
	// expObs := models.NewObservation(
	// 	[]models.NetworkLiquidity{
	// 		{Network: net, Liquidity: big.NewInt(1234)},
	// 	},
	// 	[]models.PendingTransfer{},
	// ).Encode()
	// require.Equal(t, ocrtypes.Observation(expObs), obs)
}
