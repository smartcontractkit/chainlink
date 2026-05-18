package rebalalgo

import (
	"math/big"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func TestTargetMinBalancer_ComputeTransfersToBalance_arb_eth_opt_0(t *testing.T) {
	// Cases with minimums and targets across networks are equal
	testCases := []struct {
		name             string
		balances         map[models.NetworkSelector]int64
		minimums         map[models.NetworkSelector]int64
		targets          map[models.NetworkSelector]int64
		pendingTransfers []models.ProposedTransfer
		expTransfers     []models.ProposedTransfer
	}{
		{
			name:             "imbalanced",
			balances:         map[models.NetworkSelector]int64{eth: 1100, arb: 800, opt: 1100},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers: []models.ProposedTransfer{{From: eth, To: arb, Amount: ubig.New(big.NewInt(200))},
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(100))}},
		},
		{
			name:             "all above target",
			balances:         map[models.NetworkSelector]int64{eth: 1400, arb: 1000, opt: 1100},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers:     []models.ProposedTransfer{},
		},
		{
			name:             "arb below target",
			balances:         map[models.NetworkSelector]int64{eth: 1400, arb: 800, opt: 1100},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers:     []models.ProposedTransfer{{From: eth, To: arb, Amount: ubig.New(big.NewInt(200))}},
		},
		{
			name:             "opt below target",
			balances:         map[models.NetworkSelector]int64{eth: 1400, arb: 1000, opt: 900},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers:     []models.ProposedTransfer{{From: eth, To: opt, Amount: ubig.New(big.NewInt(100))}},
		},
		{
			name:             "eth below target",
			balances:         map[models.NetworkSelector]int64{eth: 900, arb: 1000, opt: 1300},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers:     []models.ProposedTransfer{{From: opt, To: eth, Amount: ubig.New(big.NewInt(100))}},
		},
		{
			name:             "both opt and arb below target",
			balances:         map[models.NetworkSelector]int64{eth: 1500, arb: 800, opt: 900},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers: []models.ProposedTransfer{{From: eth, To: opt, Amount: ubig.New(big.NewInt(100))},
				{From: eth, To: arb, Amount: ubig.New(big.NewInt(200))}},
		},
		{
			name:             "eth below target and requires two transfers to reach target",
			balances:         map[models.NetworkSelector]int64{eth: 800, arb: 1100, opt: 1150},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers: []models.ProposedTransfer{{From: opt, To: eth, Amount: ubig.New(big.NewInt(150))},
				{From: arb, To: eth, Amount: ubig.New(big.NewInt(50))}},
		},
		{
			name:             "eth below with two sources to reach target",
			balances:         map[models.NetworkSelector]int64{eth: 900, arb: 1800, opt: 2000},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers:     []models.ProposedTransfer{{From: opt, To: eth, Amount: ubig.New(big.NewInt(100))}},
		},
		{
			name:             "eth below with two sources to reach target",
			balances:         map[models.NetworkSelector]int64{eth: 700, arb: 1400, opt: 1400},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers:     []models.ProposedTransfer{{From: opt, To: eth, Amount: ubig.New(big.NewInt(300))}},
		},
		{
			name:             "arb is below target but there is an inflight transfer that causes arb to have a surplus and give back to eth",
			balances:         map[models.NetworkSelector]int64{eth: 1000, arb: 800, opt: 1000},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{{From: eth, To: arb, Amount: ubig.New(big.NewInt(250)), Status: models.TransferStatusInflight}},
			expTransfers: []models.ProposedTransfer{
				{From: arb, To: eth, Amount: ubig.New(big.NewInt(50))},
			},
		},
		{
			name:             "arb is below target but there is inflight transfer that isn't onchain yet",
			balances:         map[models.NetworkSelector]int64{eth: 1250, arb: 800, opt: 1000},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{{From: eth, To: arb, Amount: ubig.New(big.NewInt(250)), Status: models.TransferStatusProposed}},
			expTransfers:     []models.ProposedTransfer{},
		},
		{
			name:             "eth is below target there are two sources of funding but one is already inflight",
			balances:         map[models.NetworkSelector]int64{eth: 800, arb: 2000, opt: 2200},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{{From: arb, To: eth, Amount: ubig.New(big.NewInt(250))}},
			expTransfers:     []models.ProposedTransfer{},
		},
		{
			name:             "eth is below target there are two sources of funding but one is already inflight but will not cover target",
			balances:         map[models.NetworkSelector]int64{eth: 800, arb: 2000, opt: 2200},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{{From: arb, To: eth, Amount: ubig.New(big.NewInt(100))}},
			expTransfers:     []models.ProposedTransfer{{From: opt, To: eth, Amount: ubig.New(big.NewInt(100))}},
		},
		{
			name:             "eth is below target there are two sources of funding but one is already inflight that will not cover target, both sources are used",
			balances:         map[models.NetworkSelector]int64{eth: 100, arb: 1100, opt: 1200},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{{From: arb, To: eth, Amount: ubig.New(big.NewInt(100))}},
			expTransfers: []models.ProposedTransfer{
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(200))},
				{From: arb, To: eth, Amount: ubig.New(big.NewInt(100))}},
		},
		{
			name:             "arb below target but there is no single full funding to reach target",
			balances:         map[models.NetworkSelector]int64{eth: 1100, arb: 800, opt: 1050},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers: []models.ProposedTransfer{
				{From: eth, To: arb, Amount: ubig.New(big.NewInt(150))},
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(50))},
			},
		},
		{
			name:             "opt is below target and arb can fund it with a transfer to eth",
			balances:         map[models.NetworkSelector]int64{eth: 1000, arb: 1300, opt: 800},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers: []models.ProposedTransfer{
				{From: arb, To: eth, Amount: ubig.New(big.NewInt(200))}, //we send 200 to eth knowing we are sending
				{From: eth, To: opt, Amount: ubig.New(big.NewInt(200))}, //200 to opt
			},
		},
		{
			name:             "both opt and eth are below target one transfer should be made",
			balances:         map[models.NetworkSelector]int64{eth: 900, arb: 1300, opt: 900},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers: []models.ProposedTransfer{
				{From: arb, To: eth, Amount: ubig.New(big.NewInt(200))}, //we over fill with 200 to eth knowing we are sending
				{From: eth, To: opt, Amount: ubig.New(big.NewInt(100))}, //100 to opt
			},
		},
		{
			name:             "both opt and eth are below target arb cannot fully fund both",
			balances:         map[models.NetworkSelector]int64{eth: 900, arb: 1150, opt: 900},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers: []models.ProposedTransfer{
				{From: arb, To: eth, Amount: ubig.New(big.NewInt(150))}, //we over fill with 150 to eth knowing we are sending
				{From: eth, To: opt, Amount: ubig.New(big.NewInt(50))},  //50 to opt
			},
		},
		{
			name:             "arb is below target requires transfers From both eth and opt",
			balances:         map[models.NetworkSelector]int64{eth: 1100, arb: 800, opt: 1400},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers: []models.ProposedTransfer{
				{From: eth, To: arb, Amount: ubig.New(big.NewInt(200))},
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(100))},
			},
		},
		{
			name:             "arb rebalancing is disabled and eth is below target",
			balances:         map[models.NetworkSelector]int64{eth: 800, arb: 1000, opt: 2000},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 0, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers: []models.ProposedTransfer{
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(200))},
			},
		},
		{
			name:             "both arb and opt are below target and balance cannot cover both",
			balances:         map[models.NetworkSelector]int64{eth: 1200, arb: 900, opt: 800},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers: []models.ProposedTransfer{
				{From: eth, To: opt, Amount: ubig.New(big.NewInt(200))},
			},
		},
	}

	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.DebugLevel)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			overrides := make(map[models.NetworkSelector]*big.Int)
			g := graph.NewGraph()
			for net, b := range tc.balances {
				g.(graph.GraphTest).AddNetwork(net, graph.Data{
					Liquidity:        big.NewInt(b),
					NetworkSelector:  net,
					MinimumLiquidity: big.NewInt(tc.minimums[net]),
				})
				overrides[net] = big.NewInt(tc.targets[net])
			}
			assert.NoError(t, g.(graph.GraphTest).AddConnection(eth, arb))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(arb, eth))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(eth, opt))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(opt, eth))

			pluginConfigOverrides := models.PluginConfig{
				RebalancerConfig: models.RebalancerConfig{
					Type:                   "target-and-min",
					DefaultTarget:          big.NewInt(5),
					NetworkTargetOverrides: overrides,
				},
			}

			r := NewTargetMinBalancer(lggr, pluginConfigOverrides)

			unexecuted := make([]UnexecutedTransfer, 0, len(tc.pendingTransfers))
			for _, tr := range tc.pendingTransfers {
				unexecuted = append(unexecuted, models.PendingTransfer{
					Transfer: models.Transfer{
						From:   tr.From,
						To:     tr.To,
						Amount: tr.Amount,
					},
					Status: tr.Status,
				})
			}
			transfersToBalance, err := r.ComputeTransfersToBalance(g, unexecuted)
			assert.NoError(t, err)

			for _, tr := range transfersToBalance {
				t.Logf("actual transfer: %s -> %s = %s", tr.From, tr.To, tr.Amount)
			}
			sort.Sort(models.ProposedTransfers(tc.expTransfers))
			require.Len(t, transfersToBalance, len(tc.expTransfers))
			for i, tr := range tc.expTransfers {
				t.Logf("expected transfer: %s -> %s = %s", tr.From, tr.To, tr.Amount)
				assert.Equal(t, tr.From, transfersToBalance[i].From)
				assert.Equal(t, tr.To, transfersToBalance[i].To)
				assert.Equal(t, tr.Amount.Int64(), transfersToBalance[i].Amount.Int64())
			}
		})
	}
}

func TestTargetMinBalancer_ComputeTransfersToBalance_arb_eth_opt_pending_status_behavior(t *testing.T) {
	testCases := []struct {
		name             string
		balances         map[models.NetworkSelector]int64
		minimums         map[models.NetworkSelector]int64
		targets          map[models.NetworkSelector]int64
		pendingTransfers []models.ProposedTransfer
		expTransfers     []models.ProposedTransfer
	}{
		{
			name:     "eth is below target there are multiple inflight transfers",
			balances: map[models.NetworkSelector]int64{eth: 100, arb: 2000, opt: 2000},
			minimums: map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:  map[models.NetworkSelector]int64{eth: 2000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{
				{From: arb, To: eth, Amount: ubig.New(big.NewInt(50)), Status: models.TransferStatusInflight},
				{From: arb, To: eth, Amount: ubig.New(big.NewInt(100)), Status: models.TransferStatusInflight},
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(200)), Status: models.TransferStatusInflight},
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(200)), Status: models.TransferStatusInflight},
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(50)), Status: models.TransferStatusProposed},
			},
			expTransfers: []models.ProposedTransfer{
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(450))},
				{From: arb, To: eth, Amount: ubig.New(big.NewInt(850))},
			},
		},
		{
			name:     "eth is below target there are multiple inflight transfers but not enough to balance",
			balances: map[models.NetworkSelector]int64{eth: 100, arb: 1100, opt: 2000},
			minimums: map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:  map[models.NetworkSelector]int64{eth: 2000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{
				{From: arb, To: eth, Amount: ubig.New(big.NewInt(50)), Status: models.TransferStatusInflight},
				{From: arb, To: eth, Amount: ubig.New(big.NewInt(100)), Status: models.TransferStatusInflight},
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(200)), Status: models.TransferStatusInflight},
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(200)), Status: models.TransferStatusInflight},
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(50)), Status: models.TransferStatusProposed},
			},
			expTransfers: []models.ProposedTransfer{
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(550))},
			},
		},
		{
			name:     "eth is below target there are multiple inflight transfers surplus so can take from either",
			balances: map[models.NetworkSelector]int64{eth: 100, arb: 4000, opt: 2000},
			minimums: map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:  map[models.NetworkSelector]int64{eth: 2000, arb: 1000, opt: 1000},
			pendingTransfers: []models.ProposedTransfer{
				{From: arb, To: eth, Amount: ubig.New(big.NewInt(50)), Status: models.TransferStatusReady},
				{From: arb, To: eth, Amount: ubig.New(big.NewInt(100)), Status: models.TransferStatusReady},
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(200)), Status: models.TransferStatusReady},
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(200)), Status: models.TransferStatusReady},
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(50)), Status: models.TransferStatusProposed},
			},
			expTransfers: []models.ProposedTransfer{
				{From: arb, To: eth, Amount: ubig.New(big.NewInt(1300))},
			},
		},
	}

	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.DebugLevel)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			overrides := make(map[models.NetworkSelector]*big.Int)

			g := graph.NewGraph()
			for net, b := range tc.balances {
				g.(graph.GraphTest).AddNetwork(net, graph.Data{
					Liquidity:        big.NewInt(b),
					NetworkSelector:  net,
					MinimumLiquidity: big.NewInt(tc.minimums[net]),
				})
				overrides[net] = big.NewInt(tc.targets[net])
			}
			assert.NoError(t, g.(graph.GraphTest).AddConnection(eth, arb))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(arb, eth))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(eth, opt))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(opt, eth))

			pluginConfigOverrides := models.PluginConfig{
				RebalancerConfig: models.RebalancerConfig{
					Type:                   "target-and-min",
					DefaultTarget:          big.NewInt(5),
					NetworkTargetOverrides: overrides,
				},
			}
			r := NewTargetMinBalancer(lggr, pluginConfigOverrides)

			unexecuted := make([]UnexecutedTransfer, 0, len(tc.pendingTransfers))
			for _, tr := range tc.pendingTransfers {
				unexecuted = append(unexecuted, models.PendingTransfer{
					Transfer: models.Transfer{
						From:   tr.From,
						To:     tr.To,
						Amount: tr.Amount,
					},
					Status: tr.Status,
				})
			}
			transfersToBalance, err := r.ComputeTransfersToBalance(g, unexecuted)
			assert.NoError(t, err)

			for _, tr := range transfersToBalance {
				t.Logf("actual transfer: %s -> %s = %s", tr.From, tr.To, tr.Amount)
			}
			sort.Sort(models.ProposedTransfers(tc.expTransfers))
			require.Len(t, transfersToBalance, len(tc.expTransfers))
			for i, tr := range tc.expTransfers {
				t.Logf("expected transfer: %s -> %s = %s", tr.From, tr.To, tr.Amount)
				assert.Equal(t, tr.From, transfersToBalance[i].From)
				assert.Equal(t, tr.To, transfersToBalance[i].To)
				assert.Equal(t, tr.Amount.Int64(), transfersToBalance[i].Amount.Int64())
			}
		})
	}
}

func TestTargetMinBalancer_ComputeTransfersToBalance_arb_eth_opt_base(t *testing.T) {
	testCases := []struct {
		name             string
		balances         map[models.NetworkSelector]int64
		minimums         map[models.NetworkSelector]int64
		targets          map[models.NetworkSelector]int64
		pendingTransfers []models.ProposedTransfer
		expTransfers     []models.ProposedTransfer
	}{
		{
			name:             "all above targets",
			balances:         map[models.NetworkSelector]int64{eth: 1100, arb: 1000, opt: 1100, base: 1000},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500, base: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000, base: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers:     []models.ProposedTransfer{},
		},
		{
			name:             "arb and base below target, eth and opt above target: eth tops up arb & base, opt tops up eth",
			balances:         map[models.NetworkSelector]int64{eth: 1100, arb: 900, opt: 1100, base: 900},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500, base: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000, base: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers: []models.ProposedTransfer{
				{From: eth, To: arb, Amount: ubig.New(big.NewInt(100))},
				{From: eth, To: base, Amount: ubig.New(big.NewInt(100))},
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(100))},
			},
		},
		{
			name:             "eth below target: gets funding by opt and base",
			balances:         map[models.NetworkSelector]int64{eth: 500, arb: 1000, opt: 1300, base: 1200},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500, base: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000, base: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers: []models.ProposedTransfer{
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(300))},
				{From: base, To: eth, Amount: ubig.New(big.NewInt(200))},
			},
		},
		{
			name:             "eth and arb below target: eth gets funding by opt and base, eth funds arb",
			balances:         map[models.NetworkSelector]int64{eth: 500, arb: 700, opt: 1300, base: 1500},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500, base: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000, base: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers: []models.ProposedTransfer{
				{From: base, To: eth, Amount: ubig.New(big.NewInt(500))},
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(300))},
				// we send the 300 from opt to eth but will not yet send it to arb
				// because it would make eth dip below minimum.
			},
		},
		{
			name:     "eth and arb below target with pending: eth gets funding by opt and base(inflight), eth funds arb",
			balances: map[models.NetworkSelector]int64{eth: 700, arb: 700, opt: 1300, base: 1000},
			minimums: map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500, base: 500},
			targets:  map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000, base: 1000},
			pendingTransfers: []models.ProposedTransfer{
				{From: base, To: eth, Amount: ubig.New(big.NewInt(300)), Status: models.TransferStatusNotReady},
			},
			expTransfers: []models.ProposedTransfer{
				{From: opt, To: eth, Amount: ubig.New(big.NewInt(300))},
				// can't send to arb because eth would dip below minimum
				//{From: eth, To: arb, Amount: ubig.New(big.NewInt(300))},
			},
		},
		{
			name: "opt and arb below target: eth funds opt and arb with base funds heading to eth",
			// this scenario shows that we will let eth temporarily go below target to fund opt and arb because we know we have funds coming from base
			balances:         map[models.NetworkSelector]int64{eth: 1200, arb: 800, opt: 800, base: 1200},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500, base: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000, base: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers: []models.ProposedTransfer{
				{From: eth, To: arb, Amount: ubig.New(big.NewInt(200))},
				{From: eth, To: opt, Amount: ubig.New(big.NewInt(200))},
				{From: base, To: eth, Amount: ubig.New(big.NewInt(200))},
			},
		},
		{
			name:             "all below targets",
			balances:         map[models.NetworkSelector]int64{eth: 1100, arb: 1000, opt: 1100, base: 1000},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500, base: 500},
			targets:          map[models.NetworkSelector]int64{eth: 5000, arb: 5000, opt: 5000, base: 5000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers:     []models.ProposedTransfer{},
		},
	}

	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.DebugLevel)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			overrides := make(map[models.NetworkSelector]*big.Int)
			g := graph.NewGraph()
			for net, b := range tc.balances {
				g.(graph.GraphTest).AddNetwork(net, graph.Data{
					Liquidity:        big.NewInt(b),
					NetworkSelector:  net,
					MinimumLiquidity: big.NewInt(tc.minimums[net]),
				})
				overrides[net] = big.NewInt(tc.targets[net])
			}
			assert.NoError(t, g.(graph.GraphTest).AddConnection(eth, arb))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(arb, eth))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(eth, opt))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(opt, eth))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(eth, base))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(base, eth))

			pluginConfigOverrides := models.PluginConfig{
				RebalancerConfig: models.RebalancerConfig{
					Type:                   "target-and-min",
					DefaultTarget:          big.NewInt(5),
					NetworkTargetOverrides: overrides,
				},
			}
			r := NewTargetMinBalancer(lggr, pluginConfigOverrides)

			unexecuted := make([]UnexecutedTransfer, 0, len(tc.pendingTransfers))
			for _, tr := range tc.pendingTransfers {
				unexecuted = append(unexecuted, models.PendingTransfer{
					Transfer: models.Transfer{
						From:   tr.From,
						To:     tr.To,
						Amount: tr.Amount,
					},
					Status: tr.Status,
				})
			}
			transfersToBalance, err := r.ComputeTransfersToBalance(g, unexecuted)
			assert.NoError(t, err)

			for _, tr := range transfersToBalance {
				t.Logf("actual transfer: %s -> %s = %s", tr.From, tr.To, tr.Amount)
			}
			sort.Sort(models.ProposedTransfers(tc.expTransfers))
			require.Len(t, transfersToBalance, len(tc.expTransfers))
			for i, tr := range tc.expTransfers {
				t.Logf("expected transfer: %s -> %s = %s", tr.From, tr.To, tr.Amount)
				assert.Equal(t, tr.From, transfersToBalance[i].From)
				assert.Equal(t, tr.To, transfersToBalance[i].To)
				assert.Equal(t, tr.Amount.Int64(), transfersToBalance[i].Amount.Int64())
			}
		})
	}
}

func TestTargetMinBalancer_ComputeTransfersToBalance_islands_in_graph(t *testing.T) {
	// these test have 4 networks in a spoke graph with an island node (celo) that does not have connections to the rest of the graph
	testCases := []struct {
		name             string
		balances         map[models.NetworkSelector]int64
		minimums         map[models.NetworkSelector]int64
		targets          map[models.NetworkSelector]int64
		pendingTransfers []models.ProposedTransfer
		expTransfers     []models.ProposedTransfer
	}{
		{
			name:             "all above targets",
			balances:         map[models.NetworkSelector]int64{eth: 1100, arb: 1000, opt: 1100, base: 1000, celo: 1000},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500, base: 500, celo: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000, base: 1000, celo: 1000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers:     []models.ProposedTransfer{},
		},
		{
			name: "eth and arb below: inflight transfer from eth to celo",
			// because celo is not connected to anything then nothing is done.
			balances: map[models.NetworkSelector]int64{eth: 700, arb: 900, opt: 1000, base: 1000, celo: 1000},
			minimums: map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500, base: 500, celo: 500},
			targets:  map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000, base: 1000, celo: 0},
			pendingTransfers: []models.ProposedTransfer{
				{From: eth, To: celo, Amount: ubig.New(big.NewInt(300)), Status: models.TransferStatusNotReady},
			},
			expTransfers: []models.ProposedTransfer{},
		},
		{
			name:             "celo stole all our liquidity, so we can't transfer anywhere cause everyone is below target",
			balances:         map[models.NetworkSelector]int64{eth: 700, arb: 900, opt: 800, base: 900, celo: 1700},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500, base: 500, celo: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000, base: 1000, celo: 0},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers:     []models.ProposedTransfer{},
		},
		{
			name: "celo stole some of our liquidity: base sends surplus to eth",
			// base sends it surplus to eth but nothing else can happen because eth is below target
			balances:         map[models.NetworkSelector]int64{eth: 700, arb: 900, opt: 700, base: 1100, celo: 1600},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500, base: 500, celo: 500},
			targets:          map[models.NetworkSelector]int64{eth: 1000, arb: 1000, opt: 1000, base: 1000, celo: 0},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers: []models.ProposedTransfer{
				{From: base, To: eth, Amount: ubig.New(big.NewInt(100))},
			},
		},
		{
			name:             "all below targets",
			balances:         map[models.NetworkSelector]int64{eth: 1100, arb: 1000, opt: 1100, base: 1000, celo: 1000},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500, base: 500, celo: 500},
			targets:          map[models.NetworkSelector]int64{eth: 5000, arb: 5000, opt: 5000, base: 5000, celo: 5000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers:     []models.ProposedTransfer{},
		},
	}

	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.DebugLevel)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			overrides := make(map[models.NetworkSelector]*big.Int)

			g := graph.NewGraph()
			for net, b := range tc.balances {
				g.(graph.GraphTest).AddNetwork(net, graph.Data{
					Liquidity:        big.NewInt(b),
					NetworkSelector:  net,
					MinimumLiquidity: big.NewInt(tc.minimums[net]),
				})
				overrides[net] = big.NewInt(tc.targets[net])
			}
			assert.NoError(t, g.(graph.GraphTest).AddConnection(eth, arb))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(arb, eth))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(eth, opt))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(opt, eth))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(eth, base))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(base, eth))

			pluginConfigOverrides := models.PluginConfig{
				RebalancerConfig: models.RebalancerConfig{
					Type:                   "target-and-min",
					DefaultTarget:          big.NewInt(5),
					NetworkTargetOverrides: overrides,
				},
			}
			r := NewTargetMinBalancer(lggr, pluginConfigOverrides)

			unexecuted := make([]UnexecutedTransfer, 0, len(tc.pendingTransfers))
			for _, tr := range tc.pendingTransfers {
				unexecuted = append(unexecuted, models.PendingTransfer{
					Transfer: models.Transfer{
						From:   tr.From,
						To:     tr.To,
						Amount: tr.Amount,
					},
					Status: tr.Status,
				})
			}
			transfersToBalance, err := r.ComputeTransfersToBalance(g, unexecuted)
			assert.NoError(t, err)

			for _, tr := range transfersToBalance {
				t.Logf("actual transfer: %s -> %s = %s", tr.From, tr.To, tr.Amount)
			}
			sort.Sort(models.ProposedTransfers(tc.expTransfers))
			require.Len(t, transfersToBalance, len(tc.expTransfers))
			for i, tr := range tc.expTransfers {
				t.Logf("expected transfer: %s -> %s = %s", tr.From, tr.To, tr.Amount)
				assert.Equal(t, tr.From, transfersToBalance[i].From)
				assert.Equal(t, tr.To, transfersToBalance[i].To)
				assert.Equal(t, tr.Amount.Int64(), transfersToBalance[i].Amount.Int64())
			}
		})
	}
}

func TestTargetMinBalancer_ComputeTransfersToBalance_no_tiny_transfers(t *testing.T) {
	testCases := []struct {
		name             string
		balances         map[models.NetworkSelector]int64
		minimums         map[models.NetworkSelector]int64
		targets          map[models.NetworkSelector]int64
		pendingTransfers []models.ProposedTransfer
		expTransfers     []models.ProposedTransfer
	}{
		{
			name:             "arb and opt below but not more than 5%, so even tho eth has funds we wait",
			balances:         map[models.NetworkSelector]int64{eth: 2100, arb: 1950, opt: 1950},
			minimums:         map[models.NetworkSelector]int64{eth: 500, arb: 500, opt: 500},
			targets:          map[models.NetworkSelector]int64{eth: 2000, arb: 2000, opt: 2000},
			pendingTransfers: []models.ProposedTransfer{},
			expTransfers:     []models.ProposedTransfer{},
		},
	}

	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.DebugLevel)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			overrides := make(map[models.NetworkSelector]*big.Int)

			g := graph.NewGraph()
			for net, b := range tc.balances {
				g.(graph.GraphTest).AddNetwork(net, graph.Data{
					Liquidity:        big.NewInt(b),
					NetworkSelector:  net,
					MinimumLiquidity: big.NewInt(tc.minimums[net]),
				})
				overrides[net] = big.NewInt(tc.targets[net])
			}
			assert.NoError(t, g.(graph.GraphTest).AddConnection(eth, arb))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(arb, eth))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(eth, opt))
			assert.NoError(t, g.(graph.GraphTest).AddConnection(opt, eth))

			pluginConfigOverrides := models.PluginConfig{
				RebalancerConfig: models.RebalancerConfig{
					Type:                   "target-and-min",
					DefaultTarget:          big.NewInt(5),
					NetworkTargetOverrides: overrides,
				},
			}
			r := NewTargetMinBalancer(lggr, pluginConfigOverrides)

			unexecuted := make([]UnexecutedTransfer, 0, len(tc.pendingTransfers))
			for _, tr := range tc.pendingTransfers {
				unexecuted = append(unexecuted, models.PendingTransfer{
					Transfer: models.Transfer{
						From:   tr.From,
						To:     tr.To,
						Amount: tr.Amount,
					},
					Status: tr.Status,
				})
			}
			transfersToBalance, err := r.ComputeTransfersToBalance(g, unexecuted)
			assert.NoError(t, err)

			for _, tr := range transfersToBalance {
				t.Logf("actual transfer: %s -> %s = %s", tr.From, tr.To, tr.Amount)
			}
			sort.Sort(models.ProposedTransfers(tc.expTransfers))
			require.Len(t, transfersToBalance, len(tc.expTransfers))
			for i, tr := range tc.expTransfers {
				t.Logf("expected transfer: %s -> %s = %s", tr.From, tr.To, tr.Amount)
				assert.Equal(t, tr.From, transfersToBalance[i].From)
				assert.Equal(t, tr.To, transfersToBalance[i].To)
				assert.Equal(t, tr.Amount.Int64(), transfersToBalance[i].Amount.Int64())
			}
		})
	}
}
