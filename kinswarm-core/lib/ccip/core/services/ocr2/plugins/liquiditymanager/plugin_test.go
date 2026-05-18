package liquiditymanager

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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/libocr/commontypes"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/report_encoder"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge"
	bridgemocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/mocks"
	liquiditymanagermocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/chain/evm/mocks"
	discoverermocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/discoverer/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/rebalalgo"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/testhelpers"
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
		inflight        func(t *testing.T, p *pluginWithMocks)
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
				[]models.Transfer{},
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
				a := graph.Data{
					Liquidity:               big.NewInt(1000),
					TokenAddress:            tokenX,
					LiquidityManagerAddress: rebalancerA,
					NetworkSelector:         networkA,
					ConfigDigest:            cfgDigest1,
				}
				b := graph.Data{
					Liquidity:               big.NewInt(2000),
					TokenAddress:            tokenY,
					LiquidityManagerAddress: rebalancerB,
					NetworkSelector:         networkB,
					ConfigDigest:            cfgDigest2,
				}
				assert.NoError(t, g.Add(a, b))
				assert.NoError(t, g.Add(b, a))
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
				[]models.Transfer{},
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
		{
			name:  "inflight transfers are correctly expired when they become pending",
			seqNr: 10,
			observedGraph: func(t *testing.T) (graph.Graph, error) {
				g := graph.NewGraph()
				a := graph.Data{
					Liquidity:               big.NewInt(1000),
					TokenAddress:            tokenX,
					LiquidityManagerAddress: rebalancerA,
					NetworkSelector:         networkA,
					ConfigDigest:            cfgDigest1,
				}
				b := graph.Data{
					Liquidity:               big.NewInt(2000),
					TokenAddress:            tokenY,
					LiquidityManagerAddress: rebalancerB,
					NetworkSelector:         networkB,
					ConfigDigest:            cfgDigest2,
				}
				assert.NoError(t, g.Add(a, b))
				assert.NoError(t, g.Add(b, a))
				return g, nil
			},
			previousOutcome: models.Outcome{},
			bridges: map[[2]models.NetworkSelector]func(t *testing.T) (bridge.Bridge, error){
				{networkA, networkB}: func(t *testing.T) (bridge.Bridge, error) {
					b := bridgemocks.NewBridge(t)
					b.On("GetTransfers", ctx, tokenX, tokenY).Return([]models.PendingTransfer{}, nil)
					return b, nil
				},
				{networkB, networkA}: func(t *testing.T) (bridge.Bridge, error) {
					b := bridgemocks.NewBridge(t)
					b.On("GetTransfers", ctx, tokenY, tokenX).Return([]models.PendingTransfer{
						{
							Transfer: models.NewTransfer(networkB, networkA, big.NewInt(200), time.Time{}, []byte("abc")),
							Status:   models.TransferStatusReady,
							ID:       "some-id",
						},
					}, nil)
					return b, nil
				},
			},
			inflight: func(t *testing.T, p *pluginWithMocks) {
				// A -> B transfer will be returned as pending,
				// but was previously inflight.
				p.plugin.inflight.Add(models.Transfer{
					From:   networkB,
					To:     networkA,
					Amount: ubig.New(big.NewInt(200)),
				})
			},
			expObservation: models.NewObservation(
				[]models.NetworkLiquidity{
					{Network: networkA, Liquidity: ubig.New(big.NewInt(1000))},
					{Network: networkB, Liquidity: ubig.New(big.NewInt(2000))},
				},
				[]models.Transfer{},
				[]models.PendingTransfer{
					{
						Transfer: models.NewTransfer(networkB, networkA, big.NewInt(200), time.Time{}, []byte("abc")),
						Status:   models.TransferStatusReady,
						ID:       "some-id",
					},
				},
				[]models.Transfer{},
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := newPluginWithMocksAndDefaults(t)

			// syncGraph
			mockDiscoverer := discoverermocks.NewDiscoverer(t)
			p.discovererFactory.
				On("NewDiscoverer", mock.Anything, mock.Anything).
				Return(mockDiscoverer, nil).Maybe()
			g, err := tc.observedGraph(t)
			mockDiscoverer.
				On("Discover", ctx).
				Return(g, err)
			mockDiscoverer.On("DiscoverBalances", ctx, g).Return(nil).Maybe()
			p.plugin.discoverer = mockDiscoverer

			// loadPendingTransfers && resolveProposedTransfers
			for sourceDest, bridgeFn := range tc.bridges {
				br, err2 := bridgeFn(t)
				p.bridgeFactory.
					On("NewBridge", ctx, sourceDest[0], sourceDest[1]).
					Return(br, err2)
			}

			prevObs, err := tc.previousOutcome.Encode()
			assert.NoError(t, err)
			// run the observation
			obs, err := p.plugin.Observation(ctx, ocr3types.OutcomeContext{
				SeqNr:           tc.seqNr,
				PreviousOutcome: prevObs,
			}, ocrtypes.Query{})

			if tc.expErr != nil {
				tc.expErr(t, err)
				return
			}
			assert.NoError(t, err)
			o, err := tc.expObservation.Encode()
			assert.NoError(t, err)
			assert.Equal(t, string(o), string(obs))
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
			expErr: func(t *testing.T, err error) {
				require.Error(t, err)
			},
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
			expErr: func(t *testing.T, err error) {
				require.Error(t, err)
			},
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
		{
			name: "there is an inflight transfer we should not get a new proposed transfer",
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
				PendingTransfers: []models.PendingTransfer{},
				InflightTransfers: []models.Transfer{
					models.NewTransfer(networkA, networkB, big.NewInt(1000), date2010, []byte("abc")),
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
				PendingTransfers: []models.PendingTransfer{},
				ConfigDigests: []models.ConfigDigestWithMeta{
					{NetworkSel: networkA, Digest: cfgDigest1},
					{NetworkSel: networkB, Digest: cfgDigest2},
				},
			},
			expErr: nil,
		},
		{
			name: "inflight transfers and pending transfers",
			observations: slicesRepeat[models.Observation](models.Observation{
				LiquidityPerChain: []models.NetworkLiquidity{
					{Network: networkA, Liquidity: ubig.New(big.NewInt(1000))},
					{Network: networkB, Liquidity: ubig.New(big.NewInt(2000))},
					{Network: networkC, Liquidity: ubig.New(big.NewInt(3000))},
					{Network: networkD, Liquidity: ubig.New(big.NewInt(4000))},
				},
				Edges: []models.Edge{
					{Source: networkA, Dest: networkB},
					{Source: networkB, Dest: networkA},
					{Source: networkB, Dest: networkC},
					{Source: networkC, Dest: networkA},
					{Source: networkD, Dest: networkC},
					{Source: networkC, Dest: networkD},
				},
				ConfigDigests: []models.ConfigDigestWithMeta{
					{NetworkSel: networkA, Digest: cfgDigest1},
					{NetworkSel: networkB, Digest: cfgDigest2},
					{NetworkSel: networkC, Digest: cfgDigest3},
					{NetworkSel: networkD, Digest: cfgDigest4},
				},
				PendingTransfers: []models.PendingTransfer{
					{
						Transfer: models.NewTransfer(networkA, networkB, big.NewInt(1000), date2010, []byte("abc")), // pending from A -> B
						Status:   models.TransferStatusReady,
					},
				},
				InflightTransfers: []models.Transfer{
					models.NewTransfer(networkB, networkC, big.NewInt(1000), date2010, []byte("abc")), // inflight from B -> C
				},
				ResolvedTransfers: []models.Transfer{
					models.NewTransfer(networkC, networkA, big.NewInt(234), date2011, []byte("ba-resolved")),
				},
			}, 5),
			f: 2,
			bridges: map[[2]models.NetworkSelector]func(t *testing.T) (*bridgemocks.Bridge, error){
				{networkC, networkA}: func(t *testing.T) (*bridgemocks.Bridge, error) {
					br := bridgemocks.NewBridge(t)
					br.On("QuorumizedBridgePayload", slicesRepeat([]byte("ba-resolved"), 5), 2).
						Return([]byte("quorum-ba-resolved"), nil)
					return br, nil
				},
			},
			expectedOutcome: models.Outcome{
				ProposedTransfers: []models.ProposedTransfer{
					{
						From:   networkC,
						To:     networkD,
						Amount: ubig.New(big.NewInt(2766)),
					},
				},
				ResolvedTransfers: []models.Transfer{
					models.NewTransfer(networkC, networkA, big.NewInt(234), date2011, []byte("quorum-ba-resolved")),
				},
				PendingTransfers: []models.PendingTransfer{
					{
						Transfer: models.NewTransfer(networkA, networkB, big.NewInt(1000), date2010, []byte("abc")),
						Status:   models.TransferStatusReady,
					},
				},
				ConfigDigests: []models.ConfigDigestWithMeta{
					{NetworkSel: networkA, Digest: cfgDigest1},
					{NetworkSel: networkB, Digest: cfgDigest2},
					{NetworkSel: networkC, Digest: cfgDigest3},
					{NetworkSel: networkD, Digest: cfgDigest4},
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
					On("GetBridge", sourceDest[0], sourceDest[1]).
					Return(br, err)
			}

			attributedObservations := make([]ocrtypes.AttributedObservation, 0, len(tc.observations))
			for _, o := range tc.observations {
				obs, _ := o.Encode()
				attributedObservations = append(attributedObservations, ocrtypes.AttributedObservation{
					Observation: obs,
					Observer:    commontypes.OracleID(uint8(rand.Intn(10))),
				})
			}

			outc, err := p.plugin.Outcome(ocr3types.OutcomeContext{}, ocrtypes.Query{}, attributedObservations)
			if tc.expErr != nil {
				tc.expErr(t, err)
				return
			}
			assert.NoError(t, err)
			expectedOutcome, err := tc.expectedOutcome.Encode()
			assert.NoError(t, err)
			assert.Equal(t, string(expectedOutcome), string(outc))
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
				p.plugin.liquidityGraph.(graph.GraphTest).AddNetwork(net, graph.Data{LiquidityManagerAddress: addr, NetworkSelector: net})
			}
			outcome, err := tc.outcome.Encode()
			assert.NoError(t, err)
			reports, err := p.plugin.Reports(tc.seqNr, outcome)
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
		name       string
		seqNr      uint64
		genReport  func(t *testing.T, p *pluginWithMocks) ocr3types.ReportWithInfo[models.Report]
		before     func(t *testing.T, p *pluginWithMocks)
		assertions func(t *testing.T, p *pluginWithMocks)
		expRes     bool
		expErr     bool
	}{
		{
			name:  "empty invalid report",
			seqNr: 123,
			genReport: func(t *testing.T, p *pluginWithMocks) ocr3types.ReportWithInfo[models.Report] {
				return ocr3types.ReportWithInfo[models.Report]{
					Report: []byte(`{"transfers": [], "networkID": 123}`),
				}
			},
			expRes: false,
			expErr: false,
		},
		{
			name:  "some invalid report",
			seqNr: 123,
			genReport: func(t *testing.T, p *pluginWithMocks) ocr3types.ReportWithInfo[models.Report] {
				return ocr3types.ReportWithInfo[models.Report]{
					Report: []byte(`this cannot be decoded`),
				}
			},
			expRes: false,
			expErr: true,
		},
		{
			name:  "valid report",
			seqNr: 123,
			genReport: func(t *testing.T, p *pluginWithMocks) ocr3types.ReportWithInfo[models.Report] {
				rep := models.Report{
					Transfers: []models.Transfer{
						{From: networkA, To: networkB, Amount: ubig.New(big.NewInt(1000))},
					},
					LiquidityManagerAddress: rebalancerA,
					NetworkID:               networkA,
				}
				encodedReport, err := p.plugin.reportCodec.Encode(rep)
				require.NoError(t, err)
				return ocr3types.ReportWithInfo[models.Report]{
					Report: encodedReport,
					Info:   rep,
				}
			},
			before: func(t *testing.T, p *pluginWithMocks) {
				// onchain sequence number < report sequence number
				// enough balance onchain
				mockRebalancer := liquiditymanagermocks.NewLiquidityManager(t)
				mockRebalancer.On("GetLatestSequenceNumber", mock.Anything).Return(uint64(122), nil)
				mockRebalancer.On("GetBalance", mock.Anything).Return(big.NewInt(1000), nil)
				p.lmFactory.On("NewLiquidityManager", networkA, rebalancerA).Return(mockRebalancer, nil).Once()
			},
			assertions: func(t *testing.T, p *pluginWithMocks) {
				p.lmFactory.AssertExpectations(t)
				inflight := p.plugin.inflight.GetAll()
				require.Len(t, inflight, 1)
			},
			expRes: true,
			expErr: false,
		},
		{
			name:  "stale report",
			seqNr: 123,
			genReport: func(t *testing.T, p *pluginWithMocks) ocr3types.ReportWithInfo[models.Report] {
				rep := models.Report{
					Transfers: []models.Transfer{
						{From: networkA, To: networkB, Amount: ubig.New(big.NewInt(1000))},
					},
					LiquidityManagerAddress: rebalancerA,
					NetworkID:               networkA,
				}
				encodedReport, err := p.plugin.reportCodec.Encode(rep)
				require.NoError(t, err)
				return ocr3types.ReportWithInfo[models.Report]{
					Report: encodedReport,
					Info:   rep,
				}
			},
			before: func(t *testing.T, p *pluginWithMocks) {
				// onchain sequence number == report sequence number
				mockRebalancer := liquiditymanagermocks.NewLiquidityManager(t)
				mockRebalancer.On("GetLatestSequenceNumber", mock.Anything).Return(uint64(123), nil)
				p.lmFactory.On("NewLiquidityManager", networkA, rebalancerA).Return(mockRebalancer, nil).Once()
			},
			assertions: func(t *testing.T, p *pluginWithMocks) {
				p.lmFactory.AssertExpectations(t)
				inflight := p.plugin.inflight.GetAll()
				require.Len(t, inflight, 0)
			},
			expRes: false,
			expErr: false,
		},
		{
			name:  "not enough onchain balance",
			seqNr: 123,
			genReport: func(t *testing.T, p *pluginWithMocks) ocr3types.ReportWithInfo[models.Report] {
				rep := models.Report{
					Transfers: []models.Transfer{
						{From: networkA, To: networkB, Amount: ubig.New(big.NewInt(1000))},
					},
					LiquidityManagerAddress: rebalancerA,
					NetworkID:               networkA,
				}
				encodedReport, err := p.plugin.reportCodec.Encode(rep)
				require.NoError(t, err)
				return ocr3types.ReportWithInfo[models.Report]{
					Report: encodedReport,
					Info:   rep,
				}
			},
			before: func(t *testing.T, p *pluginWithMocks) {
				// onchain sequence number < report sequence number
				// not enough balance onchain
				mockRebalancer := liquiditymanagermocks.NewLiquidityManager(t)
				mockRebalancer.On("GetLatestSequenceNumber", mock.Anything).Return(uint64(122), nil)
				mockRebalancer.On("GetBalance", mock.Anything).Return(big.NewInt(900), nil)
				p.lmFactory.On("NewLiquidityManager", networkA, rebalancerA).Return(mockRebalancer, nil).Once()
			},
			assertions: func(t *testing.T, p *pluginWithMocks) {
				p.lmFactory.AssertExpectations(t)
				inflight := p.plugin.inflight.GetAll()
				require.Len(t, inflight, 0)
			},
			expRes: false,
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := newPluginWithMocksAndDefaults(t)
			if tc.before != nil {
				tc.before(t, p)
			}
			if tc.assertions != nil {
				defer tc.assertions(t, p)
			}
			res, err := p.plugin.ShouldAcceptAttestedReport(context.Background(), tc.seqNr, tc.genReport(t, p))
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
		name           string
		report         models.Report
		reportSeqNr    uint64
		onChainSeqNr   uint64
		onchainBalance *big.Int
		expRes         bool
		expErr         bool
	}{
		{
			name: "a valid report that should be transmitted",
			report: models.Report{
				Transfers: []models.Transfer{
					{From: networkA, To: networkB, Amount: ubig.New(big.NewInt(1000))},
				},
				LiquidityManagerAddress: rebalancerA,
				NetworkID:               networkA,
				ConfigDigest:            cfgDigest1,
			},
			reportSeqNr:    11,
			onChainSeqNr:   10,
			onchainBalance: big.NewInt(1000),
			expRes:         true,
			expErr:         false,
		},
		{
			name: "report will not get transmitted since the seq num matches the on chain",
			report: models.Report{
				Transfers: []models.Transfer{
					{From: networkA, To: networkB, Amount: ubig.New(big.NewInt(1000))},
				},
				LiquidityManagerAddress: rebalancerA,
				NetworkID:               networkA,
				ConfigDigest:            cfgDigest1,
			},
			reportSeqNr:    11,
			onChainSeqNr:   11,
			onchainBalance: big.NewInt(1000),
			expRes:         false,
			expErr:         false,
		},
		{
			name: "report will not get transmitted since the on chain balance is not enough",
			report: models.Report{
				Transfers: []models.Transfer{
					{From: networkA, To: networkB, Amount: ubig.New(big.NewInt(1000))},
				},
				LiquidityManagerAddress: rebalancerA,
				NetworkID:               networkA,
				ConfigDigest:            cfgDigest1,
			},
			reportSeqNr:    11,
			onChainSeqNr:   10,
			onchainBalance: big.NewInt(900),
			expRes:         false,
			expErr:         false,
		},
	}

	ctx := testutils.Context(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := newPluginWithMocksAndDefaults(t)
			rb := liquiditymanagermocks.NewLiquidityManager(t)

			p.lmFactory.
				On("NewLiquidityManager", tc.report.NetworkID, tc.report.LiquidityManagerAddress).
				Return(rb, nil)

			rb.
				On("GetLatestSequenceNumber", ctx).
				Return(tc.onChainSeqNr, nil)

			// will only get called if onchain sequence number is less than the report sequence number
			rb.
				On("GetBalance", mock.Anything).
				Return(tc.onchainBalance, nil).
				Maybe()

			defer p.lmFactory.AssertExpectations(t)
			defer rb.AssertExpectations(t)

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
	g.(graph.GraphTest).AddNetwork(networkA, graph.Data{LiquidityManagerAddress: rebalancerA})
	g.(graph.GraphTest).AddNetwork(networkB, graph.Data{LiquidityManagerAddress: rebalancerB})
	g.(graph.GraphTest).AddNetwork(networkC, graph.Data{LiquidityManagerAddress: rebalancerC})
	p.plugin.liquidityGraph = g

	rbA := liquiditymanagermocks.NewLiquidityManager(t)
	rbB := liquiditymanagermocks.NewLiquidityManager(t)
	rbC := liquiditymanagermocks.NewLiquidityManager(t)

	p.lmFactory.On("GetLiquidityManager", networkA, rebalancerA).Return(rbA, errSomethingWentWrong) //  networkA errors while getting the rebalancer
	p.lmFactory.On("GetLiquidityManager", networkB, rebalancerB).Return(rbB, nil)
	p.lmFactory.On("GetLiquidityManager", networkC, rebalancerC).Return(rbC, nil)

	rbB.On("Close", mock.Anything).Return(errSomethingWentWrong) // networkB errors while closing
	rbC.On("Close", mock.Anything).Return(nil)                   // networkC is still closed

	err := p.plugin.Close()
	assert.Error(t, err)
	assert.Equal(t, "get liquidityManager (1, 0x000000000000000000000000000000000000000A): "+
		"some error that indicates something went wrong; "+
		"close liquidityManager (2, 0x000000000000000000000000000000000000000b): "+
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
					n.resetMocks(t)

					// the node will first discover the graph, let's mock the observed graph
					g := round.discoveredGraphPerNode[i]()
					discoverer := discoverermocks.NewDiscoverer(t)
					discoverer.
						On("Discover", mock.Anything).
						Return(g, nil).Maybe()
					discoverer.On("DiscoverBalances", mock.Anything, mock.Anything).Return(nil).Maybe()
					n.plugin.discoverer = discoverer
					n.plugin.liquidityGraph = g

					// the node will now try to load the pending transfers of all the available bridges
					// let's mock the pending transfers
					observedGraph := round.discoveredGraphPerNode[i]()
					edges, err := observedGraph.GetEdges()
					require.NoError(t, err)
					for _, edge := range edges {
						br, ok := n.bridges[[2]models.NetworkSelector{edge.Source, edge.Dest}]
						require.True(t, ok, "the test case is wrong, bridge is not defined %d->%d", edge.Source, edge.Dest)
						n.bridgeFactory.On("NewBridge", mock.Anything /* cancelContext */, edge.Source, edge.Dest).Return(br, nil).Maybe()
						n.bridgeFactory.On("GetBridge", edge.Source, edge.Dest).Return(br, nil).Maybe()

						pendingTransfers := make([]models.PendingTransfer, 0)
						for _, tr := range round.pendingTransfersPerNode[i] {
							if tr.From == edge.Source && tr.To == edge.Dest {
								pendingTransfers = append(pendingTransfers, tr)
							}
						}

						localToken, err := round.discoveredGraphPerNode[i]().GetTokenAddress(edge.Source)
						require.NoError(t, err)
						remoteToken, err := round.discoveredGraphPerNode[i]().GetTokenAddress(edge.Dest)
						require.NoError(t, err)

						br.
							On("GetTransfers", mock.Anything, localToken, remoteToken).
							Return(pendingTransfers, nil).Maybe()

						br.
							On("GetBridgePayloadAndFee", mock.Anything, mock.Anything).
							Return(nil, nativeBridgeFee, nil).Maybe()

						br.
							On("QuorumizedBridgePayload", mock.Anything, mock.Anything).
							Return(nil, nil).Maybe()
					}

					for net, data := range round.dataPerRebalancer {
						rb, exists := n.rebalancers[net]
						require.True(t, exists, "test case is wrong, seq num of rebalancer is not defined")
						rb.On("GetLatestSequenceNumber", mock.Anything).Return(data.seqNr, nil).Maybe()
						rb.On("GetBalance", mock.Anything).Return(func(context.Context) (*big.Int, error) {
							return new(big.Int).Set(data.liquidity), nil
						}).Maybe()
						n.rbFactory.On("NewLiquidityManager", net, mock.Anything).Return(rb, nil).Maybe()
					}
				}

				t.Logf(">>> running round: %d", numRound+1)
				roundResult, err := ocr3Runner.RunRound(ctx)
				if round.expErr {
					require.Error(t, err)
					continue
				}

				inflights := make([][]models.Transfer, 0, len(nodes))
				for _, n := range nodes {
					all := n.plugin.inflight.GetAll()
					inflights = append(inflights, all)
				}

				assertOutcomeEqual(t, round.expOutcome, roundResult.Outcome)
				assertReportsSlicesEqual(t, round.expTransmitted, roundResult.Transmitted)
				assertReportsSlicesEqual(t, round.expNotAccepted, roundResult.NotAccepted)
				assertReportsSlicesEqual(t, round.expNotTransmitted, roundResult.NotTransmitted)
				require.Equal(t, round.inflightPerNode, inflights)
			}
		})
	}
}

func twoNodesFourRounds(t *testing.T) testCase {
	g := graph.NewGraph()
	a := graph.Data{
		Liquidity:               big.NewInt(1000),
		TokenAddress:            tokenX,
		LiquidityManagerAddress: rebalancerA,
		XChainLiquidityManagers: nil,
		NetworkSelector:         networkA,
		ConfigDigest:            cfgDigest1,
	}
	b := graph.Data{
		Liquidity:               big.NewInt(2000),
		TokenAddress:            tokenY,
		LiquidityManagerAddress: rebalancerB,
		XChainLiquidityManagers: nil,
		NetworkSelector:         networkB,
		ConfigDigest:            cfgDigest2,
	}
	require.NoError(t, g.Add(a, b))
	require.NoError(t, g.Add(b, a))

	return testCase{
		name:     "four nodes four rounds",
		numNodes: 4,
		f:        1,
		rounds: []roundData{
			{
				// round 1 - new transfers to reach balance are generated in the outcome.
				discoveredGraphPerNode: []func() graph.Graph{
					func() graph.Graph { return g },
					func() graph.Graph { return g },
					func() graph.Graph { return g },
					func() graph.Graph { return g },
				},
				pendingTransfersPerNode: [][]models.PendingTransfer{{}, {}, {}, {}},
				inflightPerNode:         [][]models.Transfer{{}, {}, {}, {}},
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
				dataPerRebalancer: map[models.NetworkSelector]perRebalancerData{
					networkA: {
						seqNr:     1,
						liquidity: big.NewInt(1000),
					},
					networkB: {
						seqNr:     2,
						liquidity: big.NewInt(2000),
					},
				},
			},
			{
				// round 2 - the transfers of the previous outcome are included in the report.
				// they are also marked as in-flight.
				discoveredGraphPerNode: []func() graph.Graph{
					func() graph.Graph { return g },
					func() graph.Graph { return g },
					func() graph.Graph { return g },
					func() graph.Graph { return g },
				},
				pendingTransfersPerNode: [][]models.PendingTransfer{{}, {}, {}, {}},
				inflightPerNode: [][]models.Transfer{
					{models.Transfer{From: networkA, To: networkB, Amount: ubig.NewI(1000), LocalTokenAddress: tokenX,
						RemoteTokenAddress: tokenY, Sender: rebalancerA, Receiver: rebalancerB, BridgeData: hexutil.Bytes{}, NativeBridgeFee: ubig.New(nativeBridgeFee)}},
					{models.Transfer{From: networkA, To: networkB, Amount: ubig.NewI(1000), LocalTokenAddress: tokenX,
						RemoteTokenAddress: tokenY, Sender: rebalancerA, Receiver: rebalancerB, BridgeData: hexutil.Bytes{}, NativeBridgeFee: ubig.New(nativeBridgeFee)}},
					{models.Transfer{From: networkA, To: networkB, Amount: ubig.NewI(1000), LocalTokenAddress: tokenX,
						RemoteTokenAddress: tokenY, Sender: rebalancerA, Receiver: rebalancerB, BridgeData: hexutil.Bytes{}, NativeBridgeFee: ubig.New(nativeBridgeFee)}},
					{models.Transfer{From: networkA, To: networkB, Amount: ubig.NewI(1000), LocalTokenAddress: tokenX,
						RemoteTokenAddress: tokenY, Sender: rebalancerA, Receiver: rebalancerB, BridgeData: hexutil.Bytes{}, NativeBridgeFee: ubig.New(nativeBridgeFee)}},
				},
				expTransmitted: []ocr3types.ReportWithInfo[models.Report]{
					{
						Info: models.Report{
							Transfers:               []models.Transfer{{From: networkA, To: networkB, Amount: ubig.NewI(1000), NativeBridgeFee: ubig.New(nativeBridgeFee), LocalTokenAddress: tokenX, RemoteTokenAddress: tokenY, Sender: rebalancerA, Receiver: rebalancerB}},
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
					[]models.Transfer{{From: networkA, To: networkB, Amount: ubig.NewI(1000), LocalTokenAddress: tokenX, RemoteTokenAddress: tokenY, Sender: rebalancerA, Receiver: rebalancerB, NativeBridgeFee: ubig.New(nativeBridgeFee)}},
					nil,
					[]models.ConfigDigestWithMeta{{Digest: cfgDigest1, NetworkSel: networkA}, {Digest: cfgDigest2, NetworkSel: networkB}}),
				dataPerRebalancer: map[models.NetworkSelector]perRebalancerData{
					networkA: {
						seqNr:     1,
						liquidity: big.NewInt(1000),
					},
					networkB: {
						seqNr:     2,
						liquidity: big.NewInt(2000),
					},
				},
			},
			{
				// round 3 - the transfer is in flight.
				// no new transfers should be generated.
				discoveredGraphPerNode: []func() graph.Graph{
					func() graph.Graph { return g },
					func() graph.Graph { return g },
					func() graph.Graph { return g },
					func() graph.Graph { return g },
				},
				pendingTransfersPerNode: [][]models.PendingTransfer{{}, {}, {}, {}},
				inflightPerNode: [][]models.Transfer{
					{models.Transfer{From: networkA, To: networkB, Amount: ubig.NewI(1000), LocalTokenAddress: tokenX,
						RemoteTokenAddress: tokenY, Sender: rebalancerA, Receiver: rebalancerB, BridgeData: hexutil.Bytes{}, NativeBridgeFee: ubig.New(nativeBridgeFee)}},
					{models.Transfer{From: networkA, To: networkB, Amount: ubig.NewI(1000), LocalTokenAddress: tokenX,
						RemoteTokenAddress: tokenY, Sender: rebalancerA, Receiver: rebalancerB, BridgeData: hexutil.Bytes{}, NativeBridgeFee: ubig.New(nativeBridgeFee)}},
					{models.Transfer{From: networkA, To: networkB, Amount: ubig.NewI(1000), LocalTokenAddress: tokenX,
						RemoteTokenAddress: tokenY, Sender: rebalancerA, Receiver: rebalancerB, BridgeData: hexutil.Bytes{}, NativeBridgeFee: ubig.New(nativeBridgeFee)}},
					{models.Transfer{From: networkA, To: networkB, Amount: ubig.NewI(1000), LocalTokenAddress: tokenX,
						RemoteTokenAddress: tokenY, Sender: rebalancerA, Receiver: rebalancerB, BridgeData: hexutil.Bytes{}, NativeBridgeFee: ubig.New(nativeBridgeFee)}},
				},
				expTransmitted:    []ocr3types.ReportWithInfo[models.Report]{},
				expNotTransmitted: []ocr3types.ReportWithInfo[models.Report]{},
				expNotAccepted:    []ocr3types.ReportWithInfo[models.Report]{},
				expOutcome: models.NewOutcome(
					[]models.ProposedTransfer{},
					nil,
					nil,
					[]models.ConfigDigestWithMeta{{Digest: cfgDigest1, NetworkSel: networkA}, {Digest: cfgDigest2, NetworkSel: networkB}}),
				dataPerRebalancer: map[models.NetworkSelector]perRebalancerData{
					networkA: {
						seqNr:     1,
						liquidity: big.NewInt(1000),
					},
					networkB: {
						seqNr:     2,
						liquidity: big.NewInt(2000),
					},
				},
			},
			{
				// round 4 - the transfer becomes pending, and should no longer
				// be in flight. no new transfers should be generated still.
				discoveredGraphPerNode: []func() graph.Graph{
					func() graph.Graph { return g },
					func() graph.Graph { return g },
					func() graph.Graph { return g },
					func() graph.Graph { return g },
				},
				pendingTransfersPerNode: [][]models.PendingTransfer{{
					{
						Transfer: models.Transfer{
							From:               networkA,
							To:                 networkB,
							Amount:             ubig.NewI(1000),
							NativeBridgeFee:    ubig.New(nativeBridgeFee),
							LocalTokenAddress:  tokenX,
							RemoteTokenAddress: tokenY,
							Stage:              1,
						},
						Status: models.TransferStatusNotReady,
					},
				}, {
					{
						Transfer: models.Transfer{
							From:               networkA,
							To:                 networkB,
							Amount:             ubig.NewI(1000),
							NativeBridgeFee:    ubig.New(nativeBridgeFee),
							LocalTokenAddress:  tokenX,
							RemoteTokenAddress: tokenY,
							Stage:              1,
						},
						Status: models.TransferStatusNotReady,
					},
				}, {
					{
						Transfer: models.Transfer{
							From:               networkA,
							To:                 networkB,
							Amount:             ubig.NewI(1000),
							NativeBridgeFee:    ubig.New(nativeBridgeFee),
							LocalTokenAddress:  tokenX,
							RemoteTokenAddress: tokenY,
							Stage:              1,
						},
						Status: models.TransferStatusNotReady,
					},
				}, {
					{
						Transfer: models.Transfer{
							From:               networkA,
							To:                 networkB,
							Amount:             ubig.NewI(1000),
							NativeBridgeFee:    ubig.New(nativeBridgeFee),
							LocalTokenAddress:  tokenX,
							RemoteTokenAddress: tokenY,
							Stage:              1,
						},
						Status: models.TransferStatusNotReady,
					},
				}},
				inflightPerNode:   [][]models.Transfer{{}, {}, {}, {}}, // no longer inflight
				expTransmitted:    []ocr3types.ReportWithInfo[models.Report]{},
				expNotTransmitted: []ocr3types.ReportWithInfo[models.Report]{},
				expNotAccepted:    []ocr3types.ReportWithInfo[models.Report]{},
				expOutcome: models.NewOutcome(
					[]models.ProposedTransfer{},
					nil,
					[]models.PendingTransfer{
						{
							Transfer: models.Transfer{
								From:               networkA,
								To:                 networkB,
								Amount:             ubig.NewI(1000),
								NativeBridgeFee:    ubig.New(nativeBridgeFee),
								LocalTokenAddress:  tokenX,
								RemoteTokenAddress: tokenY,
							},
							Status: models.TransferStatusNotReady,
						},
					},
					[]models.ConfigDigestWithMeta{{Digest: cfgDigest1, NetworkSel: networkA}, {Digest: cfgDigest2, NetworkSel: networkB}}),
				dataPerRebalancer: map[models.NetworkSelector]perRebalancerData{
					networkA: {
						seqNr:     2,             // report posted on networkA, sequence number is incremented.
						liquidity: big.NewInt(0), // liquidity updated on A, and is pending to B.
					},
					networkB: {
						seqNr:     2,
						liquidity: big.NewInt(2000),
					},
				},
			},
		},
	}
}

func assertReportsSlicesEqual(t *testing.T, r1, r2 []ocr3types.ReportWithInfo[models.Report]) {
	require.Equal(t, len(r1), len(r2))
	for i := range r1 {
		assertReportsEqual(t, r1[i], r2[i])
	}
}

func assertReportsEqual(t *testing.T, r1, r2 ocr3types.ReportWithInfo[models.Report]) {
	assertTransfersEqual(t, r1.Info.Transfers, r2.Info.Transfers)
	require.Equal(t, r1.Info.NetworkID, r2.Info.NetworkID)
	require.Equal(t, r1.Info.LiquidityManagerAddress, r2.Info.LiquidityManagerAddress)
	require.Equal(t, r1.Info.ConfigDigest.Hex(), r2.Info.ConfigDigest.Hex())
}

func assertTransfersEqual(t *testing.T, a, b []models.Transfer) {
	require.Equal(t, len(a), len(b))
	for i := range a {
		require.Equal(t, a[i].From, b[i].From)
		require.Equal(t, a[i].To, b[i].To)
		require.Equal(t, a[i].Amount, b[i].Amount)
	}
}

func assertPendingTransfersEqual(t *testing.T, a, b []models.PendingTransfer) {
	require.Equal(t, len(a), len(b))
	for i := range a {
		require.Equal(t, a[i].From, b[i].From)
		require.Equal(t, a[i].To, b[i].To)
		require.Equal(t, a[i].Amount, b[i].Amount)
	}
}

func assertProposedTransfersEqual(t *testing.T, a, b []models.ProposedTransfer) {
	require.Equal(t, len(a), len(b))
	for i := range a {
		require.Equal(t, a[i].From, b[i].From)
		require.Equal(t, a[i].To, b[i].To)
		require.Equal(t, a[i].Amount, b[i].Amount)
	}
}

func assertOutcomeEqual(t *testing.T, exp models.Outcome, got []byte) {
	decodedOutcome, err := models.DecodeOutcome(got)
	require.NoError(t, err)
	require.Equal(t, exp.ConfigDigests, decodedOutcome.ConfigDigests)
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
	require.Positive(t, len(tc.rounds))
	require.Positive(t, tc.numNodes)
	require.NotEmpty(t, tc.name)

	for _, r := range tc.rounds {
		require.Equal(t, len(r.discoveredGraphPerNode), tc.numNodes, "you should define discovered graph per node")
		require.Equal(t, len(r.pendingTransfersPerNode), tc.numNodes, "you should define pending transfers per node")
		require.Positive(t, len(r.dataPerRebalancer), "you should define the seq nums of the rebalancers")
		require.Positive(t, len(r.dataPerRebalancer), "you should define the data of the rebalancers")
	}
}

type perRebalancerData struct {
	seqNr     uint64
	liquidity *big.Int
}

type roundData struct {
	discoveredGraphPerNode  []func() graph.Graph
	pendingTransfersPerNode [][]models.PendingTransfer
	dataPerRebalancer       map[models.NetworkSelector]perRebalancerData
	expOutcome              models.Outcome

	inflightPerNode   [][]models.Transfer
	expTransmitted    []ocr3types.ReportWithInfo[models.Report]
	expNotAccepted    []ocr3types.ReportWithInfo[models.Report]
	expNotTransmitted []ocr3types.ReportWithInfo[models.Report]
	expErr            bool
}

type node struct {
	plugin            *Plugin
	rbFactory         *mocks.Factory
	discovererFactory *discoverermocks.Factory
	bridgeFactory     *bridgemocks.Factory
	rebalancers       map[models.NetworkSelector]*liquiditymanagermocks.LiquidityManager
	bridges           map[[2]models.NetworkSelector]*bridgemocks.Bridge
}

func (n *node) resetMocks(t *testing.T) {
	lmFactory := mocks.NewFactory(t)
	discovererMock := discoverermocks.NewDiscoverer(t)
	discovererMock.On("DiscoverBalances", mock.Anything, mock.Anything).Return(nil).Maybe()
	bridgeFactory := bridgemocks.NewFactory(t)
	bridgeMocks := make(map[[2]models.NetworkSelector]*bridgemocks.Bridge)
	for _, b := range bridges {
		bridgeMocks[b] = bridgemocks.NewBridge(t)
	}

	n.bridgeFactory = bridgeFactory
	n.rbFactory = lmFactory
	n.bridges = bridgeMocks

	n.plugin.bridgeFactory = bridgeFactory
	n.plugin.discoverer = discovererMock
	n.plugin.liquidityManagerFactory = lmFactory
}

func newNode(t *testing.T, lggr logger.Logger, f int) node {
	lmFactory := mocks.NewFactory(t)
	discovererFactory := discoverermocks.NewFactory(t)
	discovererMock := discoverermocks.NewDiscoverer(t)
	discovererMock.On("DiscoverBalances", mock.Anything, mock.Anything).Return(nil).Maybe()
	// g := graph.NewGraph()
	// discovererMock.On("Discover", mock.Anything).Return(g, nil).Maybe()
	discovererFactory.On("NewDiscoverer", mock.Anything, mock.Anything).Return(discovererMock, nil).Maybe()
	bridgeFactory := bridgemocks.NewFactory(t)
	rebalancerAlg := rebalalgo.NewPingPong()

	node1 := NewPlugin(
		f,
		time.Minute,
		networkA,
		models.Address(utils.RandomAddress()),
		lmFactory,
		discovererMock,
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
		rebalancers: map[models.NetworkSelector]*liquiditymanagermocks.LiquidityManager{
			networkA: liquiditymanagermocks.NewLiquidityManager(t),
			networkB: liquiditymanagermocks.NewLiquidityManager(t),
			networkC: liquiditymanagermocks.NewLiquidityManager(t),
			networkD: liquiditymanagermocks.NewLiquidityManager(t), // todo: loop
		},
	}
}

type pluginWithMocks struct {
	plugin            *Plugin
	lmFactory         *mocks.Factory
	discovererFactory *discoverermocks.Factory
	bridgeFactory     *bridgemocks.Factory
	rebalancerAlg     *rebalalgo.PingPong
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
	lmFactory := mocks.NewFactory(t)
	discovererFactory := discoverermocks.NewFactory(t)
	discovererMock := discoverermocks.NewDiscoverer(t)
	discovererFactory.On("NewDiscoverer", mock.Anything, mock.Anything).Return(discovererMock, nil).Maybe()
	bridgeFactory := bridgemocks.NewFactory(t)
	rebalancerAlg := rebalalgo.NewPingPong()
	return &pluginWithMocks{
		plugin: NewPlugin(
			f,
			closePluginTimeout,
			rootNetwork,
			rootAddress,
			lmFactory,
			discovererMock,
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

	nativeBridgeFee = big.NewInt(10)
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

func (j JsonReportCodec) Decode(networkID models.NetworkSelector, rebalancerAddress models.Address, binaryReport []byte) (models.Report, report_encoder.ILiquidityManagerLiquidityInstructions, error) {
	var instructions report_encoder.ILiquidityManagerLiquidityInstructions
	err := json.Unmarshal(binaryReport, &instructions)
	if err != nil {
		return models.Report{}, report_encoder.ILiquidityManagerLiquidityInstructions{}, err
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
