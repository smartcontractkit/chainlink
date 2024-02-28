package rebalancer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/commontypes"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/rebalancer_report_encoder"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/bridge"
	bridgemocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/bridge/mocks"
	discoverermocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/discoverer/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/graph"
	liquiditymanagermocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditymanager/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquidityrebalancer"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
	rebalancermocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/rebalancermocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/testhelpers"
)

func TestPlugin_Query(t *testing.T) {
	t.Run("query should always be empty", func(t *testing.T) {
		p := newPluginWithMocksAndDefaults(t)
		q, err := p.plugin.Query(context.Background(), ocr3types.OutcomeContext{})
		assert.NoError(t, err)
		assert.Empty(t, q)

		q, err = p.plugin.Query(context.Background(), ocr3types.OutcomeContext{SeqNr: 1234})
		assert.NoError(t, err)
		assert.Empty(t, q)
	})
}

func TestPlugin_Observation(t *testing.T) {
	ctx := testutils.Context(t)

	testCases := []struct {
		name            string
		seqNr           uint64
		observedGraph   func(t *testing.T) (graph.Graph, error)
		previousOutcome models.Outcome
		bridges         map[[2]models.NetworkSelector]func(t *testing.T) (bridge.Bridge, error)
		expObservation  models.Observation
		expErr          func(t *testing.T, err error)
	}{
		{
			name:  "no neighboring networks",
			seqNr: 1,
			observedGraph: func(t *testing.T) (graph.Graph, error) {
				return graph.NewGraph(), nil
			},
			previousOutcome: models.Outcome{},
			bridges:         nil,
			expObservation: models.NewObservation(
				[]models.NetworkLiquidity{},
				[]models.Transfer{},
				[]models.PendingTransfer{},
				[]models.Edge{},
				[]models.ConfigDigestWithMeta{},
			),
			expErr: nil,
		},
		{
			name:  "two networks that generate a full report",
			seqNr: 1,
			observedGraph: func(t *testing.T) (graph.Graph, error) {
				g := graph.NewGraph()
				g.AddNetwork(networkA, graph.Data{
					Liquidity:         big.NewInt(1000),
					TokenAddress:      tokenX,
					RebalancerAddress: rebalancerA,
					NetworkSelector:   networkA,
					ConfigDigest:      cfgDigest1,
				})
				g.AddNetwork(networkB, graph.Data{
					Liquidity:         big.NewInt(2000),
					TokenAddress:      tokenY,
					RebalancerAddress: rebalancerB,
					NetworkSelector:   networkB,
					ConfigDigest:      cfgDigest2,
				})
				assert.NoError(t, g.AddConnection(networkA, networkB))
				assert.NoError(t, g.AddConnection(networkB, networkA))
				return g, nil
			},
			previousOutcome: models.Outcome{
				ProposedTransfers: []models.ProposedTransfer{
					{
						From:   networkB,
						To:     networkA,
						Amount: ubig.New(big.NewInt(123)),
					},
				},
			},
			bridges: map[[2]models.NetworkSelector]func(t *testing.T) (bridge.Bridge, error){
				{networkA, networkB}: func(t *testing.T) (bridge.Bridge, error) {
					b := bridgemocks.NewBridge(t)
					b.On("GetTransfers", ctx, tokenX, tokenY).Return([]models.PendingTransfer{
						{
							Transfer: models.NewTransfer(networkB, networkA, big.NewInt(200), time.Time{}, []byte("abc")),
							Status:   models.TransferStatusReady,
							ID:       "some-id",
						},
					}, nil)
					return b, nil
				},
				{networkB, networkA}: func(t *testing.T) (bridge.Bridge, error) {
					b := bridgemocks.NewBridge(t)
					b.On("GetTransfers", ctx, tokenY, tokenX).Return(nil, nil)

					// call bridge for the proposed transfer of the previous outcome
					b.On("GetBridgePayloadAndFee", ctx, models.Transfer{
						From:               networkB,
						To:                 networkA,
						Amount:             ubig.New(big.NewInt(123)),
						Sender:             rebalancerB,
						Receiver:           rebalancerA,
						LocalTokenAddress:  tokenY,
						RemoteTokenAddress: tokenX,
					}).Return([]byte("bridge-payload"), big.NewInt(4040), nil)
					return b, nil
				},
			},
			expObservation: models.NewObservation(
				[]models.NetworkLiquidity{
					{Network: networkA, Liquidity: ubig.New(big.NewInt(1000))},
					{Network: networkB, Liquidity: ubig.New(big.NewInt(2000))},
				},
				[]models.Transfer{
					{
						From:               networkB,
						To:                 networkA,
						Amount:             ubig.New(big.NewInt(123)),
						Sender:             rebalancerB,
						Receiver:           rebalancerA,
						LocalTokenAddress:  tokenY,
						RemoteTokenAddress: tokenX,
						BridgeData:         []byte("bridge-payload"),
						NativeBridgeFee:    ubig.New(big.NewInt(4040)),
					},
				},
				[]models.PendingTransfer{
					{
						Transfer: models.NewTransfer(networkB, networkA, big.NewInt(200), time.Time{}, []byte("abc")),
						Status:   models.TransferStatusReady,
						ID:       "some-id",
					},
				},
				[]models.Edge{
					models.NewEdge(networkA, networkB),
					models.NewEdge(networkB, networkA),
				},
				[]models.ConfigDigestWithMeta{
					{Digest: cfgDigest1, NetworkSel: networkA},
					{Digest: cfgDigest2, NetworkSel: networkB},
				},
			),
			expErr: nil,
		},
		{
			name:  "observation returned an error",
			seqNr: 1,
			observedGraph: func(t *testing.T) (graph.Graph, error) {
				return nil, errSomethingWentWrong
			},
			previousOutcome: models.Outcome{},
			bridges:         map[[2]models.NetworkSelector]func(t *testing.T) (bridge.Bridge, error){},
			expErr: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, errSomethingWentWrong))
				assert.False(t, err.Error() == errSomethingWentWrong.Error()) // error should be wrapped
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := newPluginWithMocksAndDefaults(t)

			// syncGraph
			mockDiscoverer := discoverermocks.NewDiscoverer(t)
			p.discovererFactory.
				On("NewDiscoverer", p.plugin.rootNetwork, p.plugin.rootAddress).
				Return(mockDiscoverer, nil)
			g, err := tc.observedGraph(t)
			mockDiscoverer.
				On("Discover", ctx).
				Return(g, err)

			// loadPendingTransfers && resolveProposedTransfers
			for sourceDest, bridgeFn := range tc.bridges {
				br, err2 := bridgeFn(t)
				p.bridgeFactory.
					On("NewBridge", sourceDest[0], sourceDest[1]).
					Return(br, err2)
			}

			// run the observation
			obs, err := p.plugin.Observation(ctx, ocr3types.OutcomeContext{
				SeqNr:           tc.seqNr,
				PreviousOutcome: tc.previousOutcome.Encode(),
			}, ocrtypes.Query{})

			if tc.expErr != nil {
				tc.expErr(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, string(tc.expObservation.Encode()), string(obs))
		})
	}
}

func TestPlugin_ValidateObservation(t *testing.T) {
	testCases := []struct {
		name   string
		obs    ocrtypes.Observation
		expErr func(t *testing.T, err error)
	}{
		{
			name: "some random bytes",
			obs:  ocrtypes.Observation("abc"),
			expErr: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "empty is ok",
			obs:  ocrtypes.Observation("{}"),
		},
		{
			name: "some observation",
			obs: models.NewObservation(
				[]models.NetworkLiquidity{},
				[]models.Transfer{{}, {}},
				[]models.PendingTransfer{},
				[]models.Edge{},
				[]models.ConfigDigestWithMeta{},
			).Encode(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := newPluginWithMocksAndDefaults(t)
			ao := ocrtypes.AttributedObservation{
				Observation: tc.obs,
				Observer:    commontypes.OracleID(uint8(rand.Intn(10))), // ignored by the plugin
			}
			err := p.plugin.ValidateObservation(ocr3types.OutcomeContext{}, ocrtypes.Query{}, ao)
			if tc.expErr != nil {
				tc.expErr(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestPlugin_ObservationQuorum(t *testing.T) {
	p := newPluginWithMocksAndDefaults(t)
	res, err := p.plugin.ObservationQuorum(ocr3types.OutcomeContext{}, ocrtypes.Query{})
	assert.NoError(t, err)
	assert.Equal(t, ocr3types.QuorumTwoFPlusOne, res)
}

func TestPlugin_Outcome(t *testing.T) {
	testCases := []struct {
		name            string
		observations    []models.Observation
		f               int
		bridges         map[[2]models.NetworkSelector]func(t *testing.T) (*bridgemocks.Bridge, error)
		expectedOutcome models.Outcome
		expErr          func(t *testing.T, err error)
	}{
		{
			name:         "zero observations",
			observations: []models.Observation{},
			f:            4,
			bridges:      nil,
			expectedOutcome: models.Outcome{
				ProposedTransfers: []models.ProposedTransfer{},
				ResolvedTransfers: []models.Transfer{},
				PendingTransfers:  []models.PendingTransfer{},
				ConfigDigests:     []models.ConfigDigestWithMeta{},
			},
			expErr: nil,
		},
		{
			name: "new proposed transfers to reach balance",
			observations: slicesRepeat[models.Observation](models.Observation{
				LiquidityPerChain: []models.NetworkLiquidity{
					{Network: networkA, Liquidity: ubig.New(big.NewInt(1000))},
					{Network: networkB, Liquidity: ubig.New(big.NewInt(2000))},
				},
				Edges: []models.Edge{
					{Source: networkA, Dest: networkB},
					{Source: networkB, Dest: networkA},
				},
				ConfigDigests: []models.ConfigDigestWithMeta{
					{NetworkSel: networkA, Digest: cfgDigest1},
					{NetworkSel: networkB, Digest: cfgDigest2},
				},
			}, 5),
			f:       2,
			bridges: nil,
			expectedOutcome: models.Outcome{
				ProposedTransfers: []models.ProposedTransfer{
					{From: networkA, To: networkB, Amount: ubig.New(big.NewInt(1000))},
				},
				ResolvedTransfers: []models.Transfer{},
				PendingTransfers:  []models.PendingTransfer{},
				ConfigDigests: []models.ConfigDigestWithMeta{
					{NetworkSel: networkA, Digest: cfgDigest1},
					{NetworkSel: networkB, Digest: cfgDigest2},
				},
			},
			expErr: nil,
		},
		{
			name: "not enough observations to reach consensus",
			observations: slicesRepeat[models.Observation](models.Observation{
				LiquidityPerChain: []models.NetworkLiquidity{
					{Network: networkA, Liquidity: ubig.New(big.NewInt(1000))},
					{Network: networkB, Liquidity: ubig.New(big.NewInt(2000))},
				},
				Edges: []models.Edge{
					{Source: networkA, Dest: networkB},
					{Source: networkB, Dest: networkA},
				},
				ConfigDigests: []models.ConfigDigestWithMeta{
					{NetworkSel: networkA, Digest: cfgDigest1},
					{NetworkSel: networkB, Digest: cfgDigest2},
				},
			}, 2),
			f:       10,
			bridges: nil,
			expectedOutcome: models.Outcome{
				ProposedTransfers: []models.ProposedTransfer{},
				ResolvedTransfers: []models.Transfer{},
				PendingTransfers:  []models.PendingTransfer{},
				ConfigDigests:     []models.ConfigDigestWithMeta{},
			},
			expErr: nil,
		},
		{
			name: "different nodes see different liquidity, median is selected",
			observations: []models.Observation{
				{
					LiquidityPerChain: []models.NetworkLiquidity{
						{Network: networkA, Liquidity: ubig.New(big.NewInt(1000))},
						{Network: networkB, Liquidity: ubig.New(big.NewInt(2000))},
					},
					Edges: []models.Edge{
						{Source: networkA, Dest: networkB},
						{Source: networkB, Dest: networkA},
					},
					ConfigDigests: []models.ConfigDigestWithMeta{
						{NetworkSel: networkA, Digest: cfgDigest1},
						{NetworkSel: networkB, Digest: cfgDigest2},
					},
				},
				{
					LiquidityPerChain: []models.NetworkLiquidity{
						{Network: networkA, Liquidity: ubig.New(big.NewInt(1100))},
						{Network: networkB, Liquidity: ubig.New(big.NewInt(2100))},
					},
					Edges: []models.Edge{
						{Source: networkA, Dest: networkB},
						{Source: networkB, Dest: networkA},
					},
					ConfigDigests: []models.ConfigDigestWithMeta{
						{NetworkSel: networkA, Digest: cfgDigest1},
						{NetworkSel: networkB, Digest: cfgDigest2},
					},
				},
				{
					LiquidityPerChain: []models.NetworkLiquidity{
						{Network: networkA, Liquidity: ubig.New(big.NewInt(1100))},
						{Network: networkB, Liquidity: ubig.New(big.NewInt(2100))},
					},
					Edges: []models.Edge{
						{Source: networkA, Dest: networkB},
						{Source: networkB, Dest: networkA},
					},
					ConfigDigests: []models.ConfigDigestWithMeta{
						{NetworkSel: networkA, Digest: cfgDigest1},
						{NetworkSel: networkB, Digest: cfgDigest2},
					},
				},
				{
					LiquidityPerChain: []models.NetworkLiquidity{
						{Network: networkA, Liquidity: ubig.New(big.NewInt(1100))},
						{Network: networkB, Liquidity: ubig.New(big.NewInt(2100))},
					},
					Edges: []models.Edge{
						{Source: networkA, Dest: networkB},
						{Source: networkB, Dest: networkA},
					},
					ConfigDigests: []models.ConfigDigestWithMeta{
						{NetworkSel: networkA, Digest: cfgDigest1},
						{NetworkSel: networkB, Digest: cfgDigest2},
					},
				},
			},
			f:       1,
			bridges: nil,
			expectedOutcome: models.Outcome{
				ProposedTransfers: []models.ProposedTransfer{
					{From: networkA, To: networkB, Amount: ubig.New(big.NewInt(1100))},
				},
				ResolvedTransfers: []models.Transfer{},
				PendingTransfers:  []models.PendingTransfer{},
				ConfigDigests: []models.ConfigDigestWithMeta{
					{NetworkSel: networkA, Digest: cfgDigest1},
					{NetworkSel: networkB, Digest: cfgDigest2},
				},
			},
			expErr: nil,
		},
		{
			name: "there is a pending transfer we should not get a new proposed transfer",
			observations: slicesRepeat[models.Observation](models.Observation{
				LiquidityPerChain: []models.NetworkLiquidity{
					{Network: networkA, Liquidity: ubig.New(big.NewInt(1000))},
					{Network: networkB, Liquidity: ubig.New(big.NewInt(2000))},
				},
				Edges: []models.Edge{
					{Source: networkA, Dest: networkB},
					{Source: networkB, Dest: networkA},
				},
				ConfigDigests: []models.ConfigDigestWithMeta{
					{NetworkSel: networkA, Digest: cfgDigest1},
					{NetworkSel: networkB, Digest: cfgDigest2},
				},
				PendingTransfers: []models.PendingTransfer{
					{
						Transfer: models.NewTransfer(networkA, networkB, big.NewInt(1000), date2010, []byte("abc")),
						Status:   models.TransferStatusReady,
						ID:       "some-transfer-id",
					},
				},
				ResolvedTransfers: []models.Transfer{
					models.NewTransfer(networkB, networkA, big.NewInt(234), date2011, []byte("ba-resolved")),
				},
			}, 5),
			f: 2,
			bridges: map[[2]models.NetworkSelector]func(t *testing.T) (*bridgemocks.Bridge, error){
				{networkB, networkA}: func(t *testing.T) (*bridgemocks.Bridge, error) {
					br := bridgemocks.NewBridge(t)
					br.On("QuorumizedBridgePayload", slicesRepeat([]byte("ba-resolved"), 5), 2).
						Return([]byte("quorum-ba-resolved"), nil)
					return br, nil
				},
			},
			expectedOutcome: models.Outcome{
				ProposedTransfers: []models.ProposedTransfer{},
				ResolvedTransfers: []models.Transfer{
					models.NewTransfer(networkB, networkA, big.NewInt(234), date2011, []byte("quorum-ba-resolved")),
				},
				PendingTransfers: []models.PendingTransfer{
					{
						Transfer: models.NewTransfer(networkA, networkB, big.NewInt(1000), date2010, []byte("abc")),
						Status:   models.TransferStatusReady,
						ID:       "some-transfer-id",
					},
				},
				ConfigDigests: []models.ConfigDigestWithMeta{
					{NetworkSel: networkA, Digest: cfgDigest1},
					{NetworkSel: networkB, Digest: cfgDigest2},
				},
			},
			expErr: nil,
		},
	}

	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.ErrorLevel)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := newPluginWithMocks(t, tc.f, time.Minute, networkA, rebalancerA, lggr)

			for sourceDest, bridgeFn := range tc.bridges {
				br, err := bridgeFn(t)
				p.bridgeFactory.
					On("NewBridge", sourceDest[0], sourceDest[1]).
					Return(br, err)
			}

			attributedObservations := make([]ocrtypes.AttributedObservation, 0, len(tc.observations))
			for _, o := range tc.observations {
				attributedObservations = append(attributedObservations, ocrtypes.AttributedObservation{
					Observation: o.Encode(),
					Observer:    commontypes.OracleID(uint8(rand.Intn(10))),
				})
			}

			outc, err := p.plugin.Outcome(ocr3types.OutcomeContext{}, ocrtypes.Query{}, attributedObservations)
			if tc.expErr != nil {
				tc.expErr(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, string(tc.expectedOutcome.Encode()), string(outc))
		})
	}
}

func TestPlugin_Reports(t *testing.T) {
	testCases := []struct {
		name              string
		outcome           models.Outcome
		expReports        []models.Report
		seqNr             uint64
		expErr            func(*testing.T, error)
		rebalancerAddress map[models.NetworkSelector]models.Address
	}{
		{
			name:       "empty outcome",
			outcome:    models.Outcome{},
			expReports: nil,
			seqNr:      0,
			expErr:     nil,
		},
		{
			name: "newly proposed transfers are ignored until they get resolved in the next round",
			outcome: models.NewOutcome(
				[]models.ProposedTransfer{
					{From: networkA, To: networkB, Amount: ubig.New(big.NewInt(1000))},
				},
				nil,
				nil,
				[]models.ConfigDigestWithMeta{{Digest: cfgDigest1, NetworkSel: networkA}, {Digest: cfgDigest2, NetworkSel: networkB}}),
			expReports: []models.Report{},
			rebalancerAddress: map[models.NetworkSelector]models.Address{
				networkA: rebalancerA,
				networkB: rebalancerB,
			},
			seqNr:  2,
			expErr: nil,
		},
		{
			name: "resolved and pending transfers are included in the reports",
			outcome: models.NewOutcome(
				nil,
				[]models.Transfer{
					{From: networkA, To: networkB, Amount: ubig.New(big.NewInt(1000)), BridgeData: []byte("ab")},
					{From: networkB, To: networkC, Amount: ubig.New(big.NewInt(2000)), BridgeData: []byte("bc")},
				},
				[]models.PendingTransfer{
					{
						Transfer: models.Transfer{From: networkC, To: networkD, Amount: ubig.New(big.NewInt(3000)), BridgeData: []byte("cd")},
						Status:   models.TransferStatusFinalized,
						ID:       "ab",
					},
					{
						Transfer: models.Transfer{From: networkC, To: networkD, Amount: ubig.New(big.NewInt(4000)), BridgeData: []byte("cd2")},
						Status:   models.TransferStatusNotReady, // ignored
						ID:       "ab",
					},
				},
				[]models.ConfigDigestWithMeta{
					{Digest: cfgDigest1, NetworkSel: networkA},
					{Digest: cfgDigest2, NetworkSel: networkB},
					{Digest: cfgDigest3, NetworkSel: networkC},
					{Digest: cfgDigest4, NetworkSel: networkD},
				},
			),
			expReports: []models.Report{
				{
					Transfers: []models.Transfer{
						{From: networkA, To: networkB, Amount: ubig.New(big.NewInt(1000)), BridgeData: []byte("ab")},
					},
					LiquidityManagerAddress: rebalancerA,
					NetworkID:               networkA,
					ConfigDigest:            cfgDigest1,
				},
				{
					Transfers: []models.Transfer{
						{From: networkB, To: networkC, Amount: ubig.New(big.NewInt(2000)), BridgeData: []byte("bc")},
					},
					LiquidityManagerAddress: rebalancerB,
					NetworkID:               networkB,
					ConfigDigest:            cfgDigest2,
				},
				{
					Transfers: []models.Transfer{
						{From: networkC, To: networkD, Amount: ubig.New(big.NewInt(3000)), BridgeData: []byte("cd")},
					},
					LiquidityManagerAddress: rebalancerD,
					NetworkID:               networkD,
					ConfigDigest:            cfgDigest4,
				},
			},
			rebalancerAddress: map[models.NetworkSelector]models.Address{
				networkA: rebalancerA,
				networkB: rebalancerB,
				networkC: rebalancerC,
				networkD: rebalancerD,
			},
			seqNr:  2,
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := newPluginWithMocksAndDefaults(t)
			for net, addr := range tc.rebalancerAddress {
				p.plugin.rebalancerGraph.AddNetwork(net, graph.Data{RebalancerAddress: addr, NetworkSelector: net})
			}

			reports, err := p.plugin.Reports(tc.seqNr, tc.outcome.Encode())
			if tc.expErr != nil {
				tc.expErr(t, err)
				return
			}
			jsonEncoder := NewJsonReportCodec()
			assert.NoError(t, err)
			assert.Equal(t, len(tc.expReports), len(reports))
			for i := range tc.expReports {
				encodedReport, err := jsonEncoder.Encode(tc.expReports[i])
				assert.NoError(t, err)
				assert.Equal(t, string(encodedReport), string(reports[i].Report))
			}
		})
	}
}

func TestPlugin_ShouldAcceptAttestedReport(t *testing.T) {
	testCases := []struct {
		name   string
		seqNr  uint64
		report ocr3types.ReportWithInfo[models.Report]
		expRes bool
		expErr bool
	}{
		{
			name:  "some valid report",
			seqNr: 123,
			report: ocr3types.ReportWithInfo[models.Report]{
				Report: []byte(`{"transfers": [], "networkID": 123}`),
			},
			expRes: true,
			expErr: false,
		},
		{
			name:  "some invalid report",
			seqNr: 123,
			report: ocr3types.ReportWithInfo[models.Report]{
				Report: []byte(`this cannot be decoded`),
			},
			expRes: false,
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := newPluginWithMocksAndDefaults(t)
			res, err := p.plugin.ShouldAcceptAttestedReport(context.Background(), tc.seqNr, tc.report)
			if tc.expErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expRes, res)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expRes, res)
		})
	}
}

func TestPlugin_ShouldTransmitAcceptedReport(t *testing.T) {
	testCases := []struct {
		name         string
		report       models.Report
		reportSeqNr  uint64
		onChainSeqNr uint64
		expRes       bool
		expErr       bool
	}{
		{
			name: "a valid report that should be transmitted",
			report: models.Report{
				Transfers: []models.Transfer{
					{From: networkA, To: networkB},
				},
				LiquidityManagerAddress: rebalancerA,
				NetworkID:               networkA,
				ConfigDigest:            cfgDigest1,
			},
			reportSeqNr:  11,
			onChainSeqNr: 10,
			expRes:       true,
			expErr:       false,
		},
		{
			name: "report will not get transmitted since the seq num matches the on chain",
			report: models.Report{
				Transfers: []models.Transfer{
					{From: networkA, To: networkB},
				},
				LiquidityManagerAddress: rebalancerA,
				NetworkID:               networkA,
				ConfigDigest:            cfgDigest1,
			},
			reportSeqNr:  11,
			onChainSeqNr: 11,
			expRes:       false,
			expErr:       false,
		},
	}

	ctx := testutils.Context(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := newPluginWithMocksAndDefaults(t)
			rb := liquiditymanagermocks.NewRebalancer(t)

			p.lmFactory.
				On("NewRebalancer", tc.report.NetworkID, tc.report.LiquidityManagerAddress).
				Return(rb, nil)

			rb.
				On("GetLatestSequenceNumber", ctx).
				Return(tc.onChainSeqNr, nil)

			encodedReport, err := p.plugin.reportCodec.Encode(tc.report)
			assert.NoError(t, err)

			res, err := p.plugin.ShouldTransmitAcceptedReport(ctx, tc.reportSeqNr, ocr3types.ReportWithInfo[models.Report]{
				Report: encodedReport,
				Info:   tc.report,
			})

			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expRes, res)
		})
	}
}

func TestPlugin_Close(t *testing.T) {
	p := newPluginWithMocksAndDefaults(t)

	g := graph.NewGraph()
	g.AddNetwork(networkA, graph.Data{RebalancerAddress: rebalancerA})
	g.AddNetwork(networkB, graph.Data{RebalancerAddress: rebalancerB})
	g.AddNetwork(networkC, graph.Data{RebalancerAddress: rebalancerC})
	p.plugin.rebalancerGraph = g

	rbA := liquiditymanagermocks.NewRebalancer(t)
	rbB := liquiditymanagermocks.NewRebalancer(t)
	rbC := liquiditymanagermocks.NewRebalancer(t)

	p.lmFactory.On("GetRebalancer", networkA, rebalancerA).Return(rbA, errSomethingWentWrong) //  networkA errors while getting the rebalancer
	p.lmFactory.On("GetRebalancer", networkB, rebalancerB).Return(rbB, nil)
	p.lmFactory.On("GetRebalancer", networkC, rebalancerC).Return(rbC, nil)

	rbB.On("Close", mock.Anything).Return(errSomethingWentWrong) // networkB errors while closing
	rbC.On("Close", mock.Anything).Return(nil)                   // networkC is still closed

	err := p.plugin.Close()
	assert.Error(t, err)
	assert.Equal(t, "get rebalancer (1, 0x000000000000000000000000000000000000000A): "+
		"some error that indicates something went wrong; "+
		"close rebalancer (2, 0x000000000000000000000000000000000000000b): "+
		"some error that indicates something went wrong", err.Error())
}

func TestPlugin_E2EWithMocks(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.ErrorLevel)

	testCases := []testCase{
		twoNodesFourRounds(t),
	}

	for _, tc := range testCases {
		tc.validate(t)
		t.Run(tc.name, func(t *testing.T) {
			// init the nodes and the ocr3 runner
			nodes := make([]node, tc.numNodes)
			plugins := make([]ocr3types.ReportingPlugin[models.Report], tc.numNodes)
			for i := range nodes {
				nodes[i] = newNode(t, lggr, tc.f)
				plugins[i] = nodes[i].plugin
			}
			ocr3Runner := testhelpers.NewOCR3Runner[models.Report](plugins)

			for numRound, round := range tc.rounds {
				for i, n := range nodes {
					t.Logf(">>> running round: %d", numRound)
					// the node will first discover the graph, let's mock the observed graph
					discoverer := discoverermocks.NewDiscoverer(t)
					n.discovererFactory.
						On("NewDiscoverer", n.plugin.rootNetwork, n.plugin.rootAddress).
						Return(discoverer, nil).Maybe()
					discoverer.
						On("Discover", mock.Anything).
						Return(round.discoveredGraphPerNode[i](), nil).Maybe()

					// the node will now try to load the pending transfers of all the available bridges
					// let's mock the pending transfers
					observedGraph := round.discoveredGraphPerNode[i]()
					edges, err := observedGraph.GetEdges()
					assert.NoError(t, err)
					for _, edge := range edges {
						br, ok := n.bridges[[2]models.NetworkSelector{edge.Source, edge.Dest}]
						assert.True(t, ok, "the test case is wrong, bridge is not defined %d->%d", edge.Source, edge.Dest)
						n.bridgeFactory.On("NewBridge", edge.Source, edge.Dest).Return(br, nil).Maybe()

						pendingTransfers := make([]models.PendingTransfer, 0)
						for _, tr := range round.pendingTransfersPerNode[i] {
							if tr.From == edge.Source && tr.To == edge.Dest {
								pendingTransfers = append(pendingTransfers, tr)
							}
						}
						br.
							On("GetTransfers", mock.Anything, mock.Anything, mock.Anything).
							Return(pendingTransfers, nil).Maybe()

						br.
							On("GetBridgePayloadAndFee", mock.Anything, mock.Anything).
							Return(nil, big.NewInt(10), nil).Maybe()

						br.
							On("QuorumizedBridgePayload", mock.Anything, mock.Anything).
							Return(nil, nil).Maybe()
					}

					for net, seqNum := range round.seqNumPerRebalancer {
						rb, exists := n.rebalancers[net]
						assert.True(t, exists, "test case is wrong, seq num of rebalancer is not defined")
						rb.On("GetLatestSequenceNumber", mock.Anything).Return(seqNum, nil).Maybe()
						n.rbFactory.On("NewRebalancer", net, mock.Anything).Return(rb, nil).Maybe()
					}
				}

				transmitted, notAccepted, notTransmitted, outcome, err := ocr3Runner.RunRound(ctx)
				if round.expErr {
					assert.Error(t, err)
					continue
				}
				assertOutcomeEqual(t, round.expOutcome, outcome)
				assertReportsSlicesEqual(t, round.expTransmitted, transmitted)
				assertReportsSlicesEqual(t, round.expNotAccepted, notAccepted)
				assertReportsSlicesEqual(t, round.expNotTransmitted, notTransmitted)
			}
		})
	}
}

func twoNodesFourRounds(t *testing.T) testCase {
	g := graph.NewGraph()
	g.AddNetwork(networkA, graph.Data{
		Liquidity:         big.NewInt(1000),
		TokenAddress:      tokenX,
		RebalancerAddress: rebalancerA,
		XChainRebalancers: nil,
		NetworkSelector:   networkA,
		ConfigDigest:      cfgDigest1,
	})
	g.AddNetwork(networkB, graph.Data{
		Liquidity:         big.NewInt(2000),
		TokenAddress:      tokenX,
		RebalancerAddress: rebalancerB,
		XChainRebalancers: nil,
		NetworkSelector:   networkB,
		ConfigDigest:      cfgDigest2,
	})
	assert.NoError(t, g.AddConnection(networkA, networkB))
	assert.NoError(t, g.AddConnection(networkB, networkA))

	return testCase{
		name:     "two nodes four rounds nothing inflight",
		numNodes: 2,
		f:        1,
		rounds: []roundData{
			{ // round 1 - new transfers to reach balance in the outcome
				discoveredGraphPerNode: []func() graph.Graph{
					func() graph.Graph { return g },
					func() graph.Graph { return g },
				},
				pendingTransfersPerNode: [][]models.PendingTransfer{{}, {}},
				expTransmitted:          []ocr3types.ReportWithInfo[models.Report]{},
				expNotTransmitted:       []ocr3types.ReportWithInfo[models.Report]{},
				expNotAccepted:          []ocr3types.ReportWithInfo[models.Report]{},
				expOutcome: models.NewOutcome(
					[]models.ProposedTransfer{
						{From: networkA, To: networkB, Amount: ubig.New(big.NewInt(1000))},
					},
					nil,
					nil,
					[]models.ConfigDigestWithMeta{{Digest: cfgDigest1, NetworkSel: networkA}, {Digest: cfgDigest2, NetworkSel: networkB}}),
				seqNumPerRebalancer: map[models.NetworkSelector]uint64{
					networkA: 1,
					networkB: 2,
				},
			},
			{ // round 2 - the transfers of the previous outcome are included in the report
				discoveredGraphPerNode: []func() graph.Graph{
					func() graph.Graph { return g },
					func() graph.Graph { return g },
				},
				pendingTransfersPerNode: [][]models.PendingTransfer{{}, {}},
				expTransmitted: []ocr3types.ReportWithInfo[models.Report]{
					{
						Info: models.Report{
							Transfers:               []models.Transfer{models.NewTransfer(networkA, networkB, big.NewInt(1000), time.Time{}, nil)},
							LiquidityManagerAddress: rebalancerA,
							NetworkID:               networkA,
							ConfigDigest:            cfgDigest1,
						},
					},
				},
				expNotTransmitted: []ocr3types.ReportWithInfo[models.Report]{},
				expNotAccepted:    []ocr3types.ReportWithInfo[models.Report]{},
				expOutcome: models.NewOutcome(
					nil,
					[]models.Transfer{{From: networkA, To: networkB, Amount: ubig.New(big.NewInt(1000))}},
					nil,
					[]models.ConfigDigestWithMeta{{Digest: cfgDigest1, NetworkSel: networkA}, {Digest: cfgDigest2, NetworkSel: networkB}}),
				seqNumPerRebalancer: map[models.NetworkSelector]uint64{
					networkA: 1,
					networkB: 2,
				},
			},
			{ // round 3 - nothing new
				discoveredGraphPerNode: []func() graph.Graph{
					func() graph.Graph { return g },
					func() graph.Graph { return g },
				},
				pendingTransfersPerNode: [][]models.PendingTransfer{{}, {}},
				expTransmitted:          []ocr3types.ReportWithInfo[models.Report]{},
				expNotTransmitted:       []ocr3types.ReportWithInfo[models.Report]{},
				expNotAccepted:          []ocr3types.ReportWithInfo[models.Report]{},
				expOutcome: models.NewOutcome(
					[]models.ProposedTransfer{
						// TODO: this slice should be empty for this test to pass.
						// right now the plugin will propose the same transfer again because it's missing an inflight cache.
						{From: networkA, To: networkB, Amount: ubig.New(big.NewInt(1000))},
					},
					nil,
					nil,
					[]models.ConfigDigestWithMeta{{Digest: cfgDigest1, NetworkSel: networkA}, {Digest: cfgDigest2, NetworkSel: networkB}}),
				seqNumPerRebalancer: map[models.NetworkSelector]uint64{
					networkA: 2,
					networkB: 3,
				},
			},
		},
	}
}

func assertReportsSlicesEqual(t *testing.T, r1, r2 []ocr3types.ReportWithInfo[models.Report]) {
	assert.Equal(t, len(r1), len(r2))
	for i := range r1 {
		assertReportsEqual(t, r1[i], r2[i])
	}
}

func assertReportsEqual(t *testing.T, r1, r2 ocr3types.ReportWithInfo[models.Report]) {
	assertTransfersEqual(t, r1.Info.Transfers, r2.Info.Transfers)
	assert.Equal(t, r1.Info.NetworkID, r2.Info.NetworkID)
	assert.Equal(t, r1.Info.LiquidityManagerAddress, r2.Info.LiquidityManagerAddress)
	assert.Equal(t, r1.Info.ConfigDigest.Hex(), r2.Info.ConfigDigest.Hex())
}

func assertTransfersEqual(t *testing.T, a, b []models.Transfer) {
	assert.Equal(t, len(a), len(b))
	for i := range a {
		assert.Equal(t, a[i].From, b[i].From)
		assert.Equal(t, a[i].To, b[i].To)
		assert.Equal(t, a[i].Amount, b[i].Amount)
	}
}

func assertPendingTransfersEqual(t *testing.T, a, b []models.PendingTransfer) {
	assert.Equal(t, len(a), len(b))
	for i := range a {
		assert.Equal(t, a[i].From, b[i].From)
		assert.Equal(t, a[i].To, b[i].To)
		assert.Equal(t, a[i].Amount, b[i].Amount)
	}
}

func assertProposedTransfersEqual(t *testing.T, a, b []models.ProposedTransfer) {
	assert.Equal(t, len(a), len(b))
	for i := range a {
		assert.Equal(t, a[i].From, b[i].From)
		assert.Equal(t, a[i].To, b[i].To)
		assert.Equal(t, a[i].Amount, b[i].Amount)
	}
}

func assertOutcomeEqual(t *testing.T, exp models.Outcome, got []byte) {
	decodedOutcome := models.Outcome{}
	err := json.Unmarshal(got, &decodedOutcome)
	assert.NoError(t, err)

	assert.Equal(t, exp.ConfigDigests, decodedOutcome.ConfigDigests)
	assertTransfersEqual(t, exp.ResolvedTransfers, decodedOutcome.ResolvedTransfers)
	assertPendingTransfersEqual(t, exp.PendingTransfers, decodedOutcome.PendingTransfers)
	assertProposedTransfersEqual(t, exp.ProposedTransfers, decodedOutcome.ProposedTransfers)
}

type testCase struct {
	name     string
	numNodes int
	f        int
	rounds   []roundData
}

func (tc *testCase) validate(t *testing.T) {
	assert.Positive(t, len(tc.rounds))
	assert.Positive(t, tc.numNodes)
	assert.NotEmpty(t, tc.name)

	for _, r := range tc.rounds {
		assert.Equal(t, len(r.discoveredGraphPerNode), tc.numNodes, "you should define discovered graph per node")
		assert.Equal(t, len(r.pendingTransfersPerNode), tc.numNodes, "you should define pending transfers per node")
		assert.Positive(t, len(r.seqNumPerRebalancer), "you should define the seq nums of the rebalancers")
	}
}

type roundData struct {
	discoveredGraphPerNode  []func() graph.Graph
	pendingTransfersPerNode [][]models.PendingTransfer
	seqNumPerRebalancer     map[models.NetworkSelector]uint64
	expOutcome              models.Outcome

	expTransmitted    []ocr3types.ReportWithInfo[models.Report]
	expNotAccepted    []ocr3types.ReportWithInfo[models.Report]
	expNotTransmitted []ocr3types.ReportWithInfo[models.Report]
	expErr            bool
}

type node struct {
	plugin            *Plugin
	rbFactory         *rebalancermocks.Factory
	discovererFactory *discoverermocks.Factory
	bridgeFactory     *bridgemocks.Factory
	rebalancers       map[models.NetworkSelector]*liquiditymanagermocks.Rebalancer
	bridges           map[[2]models.NetworkSelector]*bridgemocks.Bridge
}

func newNode(t *testing.T, lggr logger.Logger, f int) node {
	lmFactory := rebalancermocks.NewFactory(t)
	discovererFactory := discoverermocks.NewFactory(t)
	bridgeFactory := bridgemocks.NewFactory(t)
	rebalancerAlg := liquidityrebalancer.NewPingPong()

	node1 := NewPlugin(
		f,
		time.Minute,
		networkA,
		models.Address(utils.RandomAddress()),
		lmFactory,
		discovererFactory,
		bridgeFactory,
		rebalancerAlg,
		NewJsonReportCodec(),
		lggr,
	)

	bridgeMocks := make(map[[2]models.NetworkSelector]*bridgemocks.Bridge)
	for _, b := range bridges {
		bridgeMocks[b] = bridgemocks.NewBridge(t)
	}

	return node{
		plugin:            node1,
		rbFactory:         lmFactory,
		discovererFactory: discovererFactory,
		bridgeFactory:     bridgeFactory,
		bridges:           bridgeMocks,
		rebalancers: map[models.NetworkSelector]*liquiditymanagermocks.Rebalancer{
			networkA: liquiditymanagermocks.NewRebalancer(t),
			networkB: liquiditymanagermocks.NewRebalancer(t),
			networkC: liquiditymanagermocks.NewRebalancer(t),
			networkD: liquiditymanagermocks.NewRebalancer(t), // todo: loop
		},
	}
}

type pluginWithMocks struct {
	plugin            *Plugin
	lmFactory         *rebalancermocks.Factory
	discovererFactory *discoverermocks.Factory
	bridgeFactory     *bridgemocks.Factory
	rebalancerAlg     *liquidityrebalancer.PingPong
}

func newPluginWithMocksAndDefaults(t *testing.T) *pluginWithMocks {
	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.ErrorLevel)
	return newPluginWithMocks(t, 4, 4*time.Second, networkA, rebalancerA, lggr)
}

func newPluginWithMocks(
	t *testing.T,
	f int,
	closePluginTimeout time.Duration,
	rootNetwork models.NetworkSelector,
	rootAddress models.Address,
	lggr logger.Logger,
) *pluginWithMocks {
	lmFactory := rebalancermocks.NewFactory(t)
	discovererFactory := discoverermocks.NewFactory(t)
	bridgeFactory := bridgemocks.NewFactory(t)
	rebalancerAlg := liquidityrebalancer.NewPingPong()
	return &pluginWithMocks{
		plugin: NewPlugin(
			f,
			closePluginTimeout,
			rootNetwork,
			rootAddress,
			lmFactory,
			discovererFactory,
			bridgeFactory,
			rebalancerAlg,
			NewJsonReportCodec(),
			lggr,
		),
		lmFactory:         lmFactory,
		discovererFactory: discovererFactory,
		bridgeFactory:     bridgeFactory,
		rebalancerAlg:     rebalancerAlg,
	}
}

func slicesRepeat[T any](val T, times int) []T {
	r := make([]T, times)
	for i := range r {
		r[i] = val
	}
	return r
}

// test helper variables

var (
	networkA = models.NetworkSelector(1)
	networkB = models.NetworkSelector(2)
	networkC = models.NetworkSelector(3)
	networkD = models.NetworkSelector(4)

	rebalancerA = models.Address(common.HexToAddress("0xa"))
	rebalancerB = models.Address(common.HexToAddress("0xb"))
	rebalancerC = models.Address(common.HexToAddress("0xc"))
	rebalancerD = models.Address(common.HexToAddress("0xd"))

	tokenX = models.Address(common.HexToAddress("0x1"))
	tokenY = models.Address(common.HexToAddress("0x2"))

	bridgeAB = [2]models.NetworkSelector{networkA, networkB}
	bridgeBA = [2]models.NetworkSelector{networkB, networkA}

	bridges = [][2]models.NetworkSelector{
		bridgeAB,
		bridgeBA,
	}

	cfgDigest1 = models.ConfigDigest{ConfigDigest: ocrtypes.ConfigDigest([32]byte{1})}
	cfgDigest2 = models.ConfigDigest{ConfigDigest: ocrtypes.ConfigDigest([32]byte{2})}
	cfgDigest3 = models.ConfigDigest{ConfigDigest: ocrtypes.ConfigDigest([32]byte{3})}
	cfgDigest4 = models.ConfigDigest{ConfigDigest: ocrtypes.ConfigDigest([32]byte{4})}

	errSomethingWentWrong = errors.New("some error that indicates something went wrong")

	date2010 = time.Date(2010, 5, 6, 12, 4, 4, 0, time.UTC)
	date2011 = time.Date(2011, 5, 6, 12, 4, 4, 0, time.UTC)
)

// JsonReportCodec is used in tests
type JsonReportCodec struct{}

func NewJsonReportCodec() JsonReportCodec {
	return JsonReportCodec{}
}

func (j JsonReportCodec) Encode(report models.Report) ([]byte, error) {
	instructions, err := report.ToLiquidityInstructions()
	if err != nil {
		return nil, fmt.Errorf("converting to liquidity instructions: %w", err)
	}
	return json.Marshal(instructions)
}

func (j JsonReportCodec) Decode(networkID models.NetworkSelector, rebalancerAddress models.Address, binaryReport []byte) (models.Report, rebalancer_report_encoder.IRebalancerLiquidityInstructions, error) {
	var instructions rebalancer_report_encoder.IRebalancerLiquidityInstructions
	err := json.Unmarshal(binaryReport, &instructions)
	if err != nil {
		return models.Report{}, rebalancer_report_encoder.IRebalancerLiquidityInstructions{}, err
	}
	var r models.Report
	for _, sendInstruction := range instructions.SendLiquidityParams {
		r.Transfers = append(r.Transfers, models.Transfer{
			From: networkID,
			To:   models.NetworkSelector(sendInstruction.RemoteChainSelector),
			Amount: ubig.New(
				sendInstruction.Amount,
			),
		})
	}
	for _, receiveInstruction := range instructions.ReceiveLiquidityParams {
		r.Transfers = append(r.Transfers, models.Transfer{
			From: models.NetworkSelector(receiveInstruction.RemoteChainSelector),
			To:   networkID,
			Amount: ubig.New(
				receiveInstruction.Amount,
			),
		})
	}
	r.LiquidityManagerAddress = rebalancerAddress
	return r, instructions, err
}
