package rebalalgo

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

var eth = models.NetworkSelector(5009297550715157269)
var opt = models.NetworkSelector(3734403246176062136)
var arb = models.NetworkSelector(4949039107694359620)
var base = models.NetworkSelector(15971525489660198786)
var celo = models.NetworkSelector(1346049177634351622)

func TestTargetBalanceRebalancer_ComputeTransfersToBalance_arb_eth_opt(t *testing.T) {
	type transfer struct {
		from   models.NetworkSelector
		to     models.NetworkSelector
		am     int64
		status models.TransferStatus
	}

	testCases := []struct {
		name             string
		balances         map[models.NetworkSelector]int64
		minimums         map[models.NetworkSelector]int64
		pendingTransfers []transfer
		expTransfers     []transfer
	}{
		{
			name:             "all above target",
			balances:         map[models.NetworkSelector]int64{eth: 1400, arb: 1000, opt: 1100},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{},
			expTransfers:     []transfer{},
		},
		{
			name:             "arb below target",
			balances:         map[models.NetworkSelector]int64{eth: 1400, arb: 800, opt: 1100},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{},
			expTransfers:     []transfer{{from: eth, to: arb, am: 200}},
		},
		{
			name:             "opt below target",
			balances:         map[models.NetworkSelector]int64{eth: 1400, arb: 1000, opt: 900},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{},
			expTransfers:     []transfer{{from: eth, to: opt, am: 100}},
		},
		{
			name:             "eth below target",
			balances:         map[models.NetworkSelector]int64{eth: 900, arb: 1000, opt: 1300},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{},
			expTransfers:     []transfer{{from: opt, to: eth, am: 100}},
		},
		{
			name:             "both opt and arb below target",
			balances:         map[models.NetworkSelector]int64{eth: 1500, arb: 800, opt: 900},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{},
			expTransfers:     []transfer{{from: eth, to: opt, am: 100}, {from: eth, to: arb, am: 200}},
		},
		{
			name:             "eth below target and requires two transfers to reach target",
			balances:         map[models.NetworkSelector]int64{eth: 800, arb: 1100, opt: 1150},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{},
			expTransfers:     []transfer{{from: opt, to: eth, am: 150}, {from: arb, to: eth, am: 50}},
		},
		{
			name:             "eth below with two sources to reach target the highest one is selected",
			balances:         map[models.NetworkSelector]int64{eth: 900, arb: 2000, opt: 1800},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{},
			expTransfers:     []transfer{{from: arb, to: eth, am: 100}},
		},
		{
			name:             "eth below with two sources to reach target the highest one is selected - reversed",
			balances:         map[models.NetworkSelector]int64{eth: 900, arb: 1800, opt: 2000},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{},
			expTransfers:     []transfer{{from: opt, to: eth, am: 100}},
		},
		{
			name:             "eth below with two sources to reach target should be deterministic",
			balances:         map[models.NetworkSelector]int64{eth: 700, arb: 1400, opt: 1400},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{},
			expTransfers:     []transfer{{from: opt, to: eth, am: 300}},
		},
		{
			name:             "arb is below target but there is an inflight transfer",
			balances:         map[models.NetworkSelector]int64{eth: 1000, arb: 800, opt: 1000},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{{from: eth, to: arb, am: 250}},
			expTransfers:     []transfer{},
		},
		{
			name: "arb is below target but there is inflight transfer that isn't onchain yet",
			// since it's not on-chain yet the source balance should be expected to be lower than the current value
			balances:         map[models.NetworkSelector]int64{eth: 1250, arb: 800, opt: 1000},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{{from: eth, to: arb, am: 250, status: models.TransferStatusProposed}},
			expTransfers:     []transfer{},
		},
		{
			name:             "eth is below target there are two sources of funding but one is already inflight",
			balances:         map[models.NetworkSelector]int64{eth: 800, arb: 2000, opt: 2200},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{{from: arb, to: eth, am: 250}},
			expTransfers:     []transfer{},
		},
		{
			name:             "eth is below target there are two sources of funding but one is already inflight but will not cover target",
			balances:         map[models.NetworkSelector]int64{eth: 800, arb: 2000, opt: 2200},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{{from: arb, to: eth, am: 100}},
			expTransfers:     []transfer{{from: opt, to: eth, am: 100}},
		},
		{
			name:             "eth is below target there are two sources of funding but one is already inflight that will not cover target, both sources are used",
			balances:         map[models.NetworkSelector]int64{eth: 100, arb: 1100, opt: 1200},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{{from: arb, to: eth, am: 100}},
			expTransfers:     []transfer{{from: opt, to: eth, am: 200}, {from: arb, to: eth, am: 100}},
		},
		{
			name:     "eth is below target there are multiple inflight transfers",
			balances: map[models.NetworkSelector]int64{eth: 100, arb: 2000, opt: 2000},
			minimums: map[models.NetworkSelector]int64{eth: 2000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{
				{from: arb, to: eth, am: 50},
				{from: arb, to: eth, am: 100},
				{from: opt, to: eth, am: 50, status: models.TransferStatusProposed},
				{from: opt, to: eth, am: 200},
				{from: opt, to: eth, am: 200},
			},
			expTransfers: []transfer{
				{from: opt, to: eth, am: 300},
				{from: arb, to: eth, am: 1000},
			},
		},
		{
			name:     "eth is below target there are multiple inflight transfers 2",
			balances: map[models.NetworkSelector]int64{eth: 100, arb: 1100, opt: 2000},
			minimums: map[models.NetworkSelector]int64{eth: 2000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{
				{from: arb, to: eth, am: 50},
				{from: arb, to: eth, am: 100},
				{from: opt, to: eth, am: 50, status: models.TransferStatusProposed},
				{from: opt, to: eth, am: 200},
				{from: opt, to: eth, am: 200},
			},
			expTransfers: []transfer{
				{from: opt, to: eth, am: 950},
				{from: arb, to: eth, am: 100},
			},
		},
		{
			name:     "eth is below target there are multiple inflight transfers 3",
			balances: map[models.NetworkSelector]int64{eth: 100, arb: 4000, opt: 2000},
			minimums: map[models.NetworkSelector]int64{eth: 2000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{
				{from: arb, to: eth, am: 50},
				{from: arb, to: eth, am: 100},
				{from: opt, to: eth, am: 50, status: models.TransferStatusProposed},
				{from: opt, to: eth, am: 200},
				{from: opt, to: eth, am: 200},
			},
			expTransfers: []transfer{
				{from: arb, to: eth, am: 1300},
			},
		},
		{
			name:             "arb below target but there is no full funding to reach target",
			balances:         map[models.NetworkSelector]int64{eth: 1100, arb: 800, opt: 1050},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{},
			expTransfers: []transfer{
				{from: opt, to: eth, am: 50}, // 2hop transfer
				{from: eth, to: arb, am: 100},
			}, // transfer is made but without reaching target
		},
		{
			name:             "opt is below target and arb can fund it with a transfer to eth",
			balances:         map[models.NetworkSelector]int64{eth: 1000, arb: 1300, opt: 800},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{},
			expTransfers: []transfer{
				{from: arb, to: eth, am: 200},
			},
		},
		{
			name:             "both opt and eth are below target one transfer should be made",
			balances:         map[models.NetworkSelector]int64{eth: 900, arb: 1300, opt: 900},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{},
			expTransfers: []transfer{
				{from: arb, to: eth, am: 200},
			},
		},
		{
			name:             "both opt and eth are below target arb cannot fully fund both",
			balances:         map[models.NetworkSelector]int64{eth: 900, arb: 1150, opt: 900},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{},
			expTransfers: []transfer{
				{from: arb, to: eth, am: 150},
			},
		},
		{
			name:             "arb is below target requires transfers from both eth and opt",
			balances:         map[models.NetworkSelector]int64{eth: 1100, arb: 800, opt: 1400},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{},
			expTransfers: []transfer{
				{from: opt, to: eth, am: 100},
				{from: eth, to: arb, am: 100},
			},
		},
		{
			name:             "arb rebalancing is disabled and eth is below target",
			balances:         map[models.NetworkSelector]int64{eth: 800, arb: 1000, opt: 2000},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 0, opt: 1000},
			pendingTransfers: []transfer{},
			expTransfers: []transfer{
				{from: opt, to: eth, am: 200},
			},
		},
		{
			name: "both arb and opt are below target and balance cannot cover both",
			// in this case owner should reconsider minimums or increase liquidity
			balances:         map[models.NetworkSelector]int64{eth: 1200, arb: 900, opt: 800},
			minimums:         map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []transfer{},
			expTransfers: []transfer{
				{from: eth, to: opt, am: 200},
			},
		},
	}

	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.DebugLevel)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := graph.NewGraph()
			for net, b := range tc.balances {
				g.(graph.GraphTest).AddNetwork(net, graph.Data{
					Liquidity:        big.NewInt(b),
					NetworkSelector:  net,
					MinimumLiquidity: big.NewInt(tc.minimums[net]),
					TargetLiquidity:  big.NewInt(tc.minimums[net]),
				})
			}
			assert.NoError(t, g.(graph.GraphTest).AddConnection(eth, arb))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(arb, eth))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(eth, opt))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(opt, eth))

			r := NewMinLiquidityRebalancer(lggr)

			unexecuted := make([]UnexecutedTransfer, 0, len(tc.pendingTransfers))
			for _, tr := range tc.pendingTransfers {
				unexecuted = append(unexecuted, models.PendingTransfer{
					Transfer: models.Transfer{
						From:   tr.from,
						To:     tr.to,
						Amount: ubig.New(big.NewInt(tr.am)),
					},
					Status: tr.status,
				})
			}
			transfersToBalance, err := r.ComputeTransfersToBalance(g, unexecuted)
			assert.NoError(t, err)

			for _, tr := range transfersToBalance {
				t.Logf("actual transfer: %v->%v %s", tr.From, tr.To, tr.Amount)
			}

			assert.Len(t, transfersToBalance, len(tc.expTransfers))
			for i, tr := range tc.expTransfers {
				t.Logf("expected transfer: %v->%v %d", tr.from, tr.to, tr.am)
				assert.Equal(t, tr.from, transfersToBalance[i].From)
				assert.Equal(t, tr.to, transfersToBalance[i].To)
				assert.Equal(t, tr.am, transfersToBalance[i].Amount.Int64())
			}
		})
	}
}
