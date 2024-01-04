package rebalancer

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditygraph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
	mocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/rebalancermocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type mockDeps struct {
	mockFactory    *mocks.Factory
	mockRebalancer *mocks.Rebalancer
}

func newPlugin(t *testing.T) (*Plugin, mockDeps) {
	f := 10
	closeTimeout := 5 * time.Second
	rootNetwork := models.NetworkID(1)
	rootAddr := models.Address(utils.RandomAddress())

	lmGraph := liquiditygraph.NewGraph()
	lmFactory := mocks.NewFactory(t)
	rb := mocks.NewRebalancer(t)
	return NewPlugin(f, closeTimeout, rootNetwork, rootAddr, lmFactory, lmGraph, rb), mockDeps{
		mockFactory:    lmFactory,
		mockRebalancer: rb,
	}
}

func TestPluginQuery(t *testing.T) {
	p, _ := newPlugin(t)
	q, err := p.Query(context.Background(), ocr3types.OutcomeContext{})
	assert.Empty(t, q, "query should always be empty")
	assert.NoError(t, err)
}

func TestPluginObservation(t *testing.T) {
	ctx := testutils.Context(t)
	p, deps := newPlugin(t)

	lms := p.liquidityManagers.GetAll()
	assert.Len(t, lms, 1, "plugin should initially contain one lm")
	net := maps.Keys(lms)[0]
	addr := maps.Values(lms)[0]

	mockLM := mocks.NewLiquidityManager(t)
	deps.mockFactory.On("NewLiquidityManager", net, addr).Return(mockLM, nil)

	mockLM.On("GetLiquidityManagers", ctx).Return(map[models.NetworkID]models.Address{}, nil)
	mockLM.On("GetBalance", ctx).Return(big.NewInt(1234), nil)
	mockLM.On("GetPendingTransfers", ctx).Return([]models.PendingTransfer{}, nil)

	obs, err := p.Observation(ctx, ocr3types.OutcomeContext{}, ocrtypes.Query{})
	assert.NoError(t, err)
	expObs := models.NewObservation(
		[]models.NetworkLiquidity{
			{Network: net, Liquidity: big.NewInt(1234)},
		},
		[]models.PendingTransfer{},
	).Encode()
	assert.Equal(t, ocrtypes.Observation(expObs), obs)
}
