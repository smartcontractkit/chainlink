package commit

import (
	"context"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/ccipocr3/internal/libs/testhelpers"
	"github.com/smartcontractkit/ccipocr3/internal/mocks"
	"github.com/smartcontractkit/ccipocr3/internal/reader"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	libocrtypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

func TestPlugin(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Test(t)

	testCases := []struct {
		name                  string
		description           string
		nodes                 []nodeSetup
		expErr                func(*testing.T, error)
		expOutcome            cciptypes.CommitPluginOutcome
		expTransmittedReports []cciptypes.CommitPluginReport
		initialOutcome        cciptypes.CommitPluginOutcome
	}{
		{
			name:        "EmptyOutcome",
			description: "Empty observations are returned by all nodes which leads to an empty outcome.",
			nodes:       setupEmptyOutcome(ctx, t, lggr),
			expErr:      func(t *testing.T, err error) { assert.Equal(t, testhelpers.ErrEmptyOutcome, err) },
		},
		{
			name: "AllNodesReadAllChains",
			description: "Nodes observe the latest sequence numbers and new messages after those sequence numbers. " +
				"They also observe gas prices. In this setup all nodes can read all chains.",
			nodes: setupAllNodesReadAllChains(ctx, t, lggr),
			expOutcome: cciptypes.CommitPluginOutcome{
				MaxSeqNums: []cciptypes.SeqNumChain{
					{ChainSel: chainA, SeqNum: 10},
					{ChainSel: chainB, SeqNum: 20},
				},
				MerkleRoots: []cciptypes.MerkleRootChain{
					{ChainSel: chainB, MerkleRoot: cciptypes.Bytes32{}, SeqNumsRange: cciptypes.NewSeqNumRange(21, 22)},
				},
				TokenPrices: []cciptypes.TokenPrice{},
				GasPrices: []cciptypes.GasPriceChain{
					{ChainSel: chainA, GasPrice: cciptypes.NewBigIntFromInt64(1000)},
					{ChainSel: chainB, GasPrice: cciptypes.NewBigIntFromInt64(20_000)},
				},
			},
			expTransmittedReports: []cciptypes.CommitPluginReport{
				{
					MerkleRoots: []cciptypes.MerkleRootChain{
						{ChainSel: chainB, SeqNumsRange: cciptypes.NewSeqNumRange(21, 22)},
					},
					PriceUpdates: cciptypes.PriceUpdates{
						TokenPriceUpdates: []cciptypes.TokenPrice{},
						GasPriceUpdates: []cciptypes.GasPriceChain{
							{ChainSel: chainA, GasPrice: cciptypes.NewBigIntFromInt64(1000)},
							{ChainSel: chainB, GasPrice: cciptypes.NewBigIntFromInt64(20_000)},
						},
					},
				},
			},
			initialOutcome: cciptypes.CommitPluginOutcome{
				MaxSeqNums: []cciptypes.SeqNumChain{
					{ChainSel: chainA, SeqNum: 10},
					{ChainSel: chainB, SeqNum: 20},
				},
				MerkleRoots: []cciptypes.MerkleRootChain{},
				TokenPrices: []cciptypes.TokenPrice{},
				GasPrices:   []cciptypes.GasPriceChain{},
			},
		},
		{
			name:        "NodesDoNotAgreeOnMsgs",
			description: "Nodes do not agree on messages which leads to an outcome with empty merkle roots.",
			nodes:       setupNodesDoNotAgreeOnMsgs(ctx, t, lggr),
			expOutcome: cciptypes.CommitPluginOutcome{
				MaxSeqNums: []cciptypes.SeqNumChain{
					{ChainSel: chainA, SeqNum: 10},
					{ChainSel: chainB, SeqNum: 20},
				},
				MerkleRoots: []cciptypes.MerkleRootChain{},
				TokenPrices: []cciptypes.TokenPrice{},
				GasPrices: []cciptypes.GasPriceChain{
					{ChainSel: chainA, GasPrice: cciptypes.NewBigIntFromInt64(1000)},
					{ChainSel: chainB, GasPrice: cciptypes.NewBigIntFromInt64(20_000)},
				},
			},
			expTransmittedReports: []cciptypes.CommitPluginReport{
				{
					MerkleRoots: []cciptypes.MerkleRootChain{},
					PriceUpdates: cciptypes.PriceUpdates{
						TokenPriceUpdates: []cciptypes.TokenPrice{},
						GasPriceUpdates: []cciptypes.GasPriceChain{
							{ChainSel: chainA, GasPrice: cciptypes.NewBigIntFromInt64(1000)},
							{ChainSel: chainB, GasPrice: cciptypes.NewBigIntFromInt64(20_000)},
						},
					},
				},
			},
			initialOutcome: cciptypes.CommitPluginOutcome{
				MaxSeqNums: []cciptypes.SeqNumChain{
					{ChainSel: chainA, SeqNum: 10},
					{ChainSel: chainB, SeqNum: 20},
				},
				MerkleRoots: []cciptypes.MerkleRootChain{},
				TokenPrices: []cciptypes.TokenPrice{},
				GasPrices:   []cciptypes.GasPriceChain{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-")
			t.Logf(">>> [%s]\n", tc.name)
			t.Logf(">>> %s\n", tc.description)
			defer t.Log("-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-")

			nodesSetup := tc.nodes
			nodes := make([]ocr3types.ReportingPlugin[[]byte], 0, len(nodesSetup))
			for _, n := range nodesSetup {
				nodes = append(nodes, n.node)
			}

			nodeIDs := make([]commontypes.OracleID, 0, len(nodesSetup))
			for _, n := range nodesSetup {
				nodeIDs = append(nodeIDs, n.node.nodeID)
			}
			o, err := tc.initialOutcome.Encode()
			require.NoError(t, err)
			runner := testhelpers.NewOCR3Runner(nodes, nodeIDs, o)

			res, err := runner.RunRound(ctx)
			if tc.expErr != nil {
				tc.expErr(t, err)
			} else {
				assert.NoError(t, err)
			}

			if !reflect.DeepEqual(tc.expOutcome, cciptypes.CommitPluginOutcome{}) {
				outcome, err := cciptypes.DecodeCommitPluginOutcome(res.Outcome)
				assert.NoError(t, err)
				assert.Equal(t, tc.expOutcome.TokenPrices, outcome.TokenPrices)
				assert.Equal(t, tc.expOutcome.MaxSeqNums, outcome.MaxSeqNums)
				assert.Equal(t, tc.expOutcome.GasPrices, outcome.GasPrices)

				assert.Equal(t, len(tc.expOutcome.MerkleRoots), len(outcome.MerkleRoots))
				for i, exp := range tc.expOutcome.MerkleRoots {
					assert.Equal(t, exp.ChainSel, outcome.MerkleRoots[i].ChainSel)
					assert.Equal(t, exp.SeqNumsRange, outcome.MerkleRoots[i].SeqNumsRange)
				}
			}

			assert.Equal(t, len(tc.expTransmittedReports), len(res.Transmitted))
			for i, exp := range tc.expTransmittedReports {
				actual, err := nodesSetup[0].reportCodec.Decode(ctx, res.Transmitted[i].Report)
				assert.NoError(t, err)
				assert.Equal(t, exp.PriceUpdates, actual.PriceUpdates)
				assert.Equal(t, len(exp.MerkleRoots), len(actual.MerkleRoots))
				for j, expRoot := range exp.MerkleRoots {
					assert.Equal(t, expRoot.ChainSel, actual.MerkleRoots[j].ChainSel)
					assert.Equal(t, expRoot.SeqNumsRange, actual.MerkleRoots[j].SeqNumsRange)
				}
			}
		})
	}
}

func setupEmptyOutcome(ctx context.Context, t *testing.T, lggr logger.Logger) []nodeSetup {
	cfg := cciptypes.CommitPluginConfig{
		DestChain:           chainC,
		PricedTokens:        []types.Account{tokenX},
		TokenPricesObserver: false,
		NewMsgScanBatchSize: 256,
	}

	chainConfigInfos := []reader.ChainConfigInfo{
		{
			ChainSelector: chainC,
			ChainConfig: reader.HomeChainConfigMapper{
				FChain: 1,
				Readers: []libocrtypes.PeerID{
					{1}, {2}, {3},
				},
				Config: []byte{0},
			},
		},
	}

	homeChain := setupHomeChainPoller(lggr, chainConfigInfos)
	err := homeChain.Start(ctx)
	if err != nil {
		return nil
	}

	oracleIDToP2pID := GetP2pIDs(1, 2, 3)
	nodes := []nodeSetup{
		newNode(ctx, t, lggr, 1, cfg, homeChain, oracleIDToP2pID),
		newNode(ctx, t, lggr, 2, cfg, homeChain, oracleIDToP2pID),
		newNode(ctx, t, lggr, 3, cfg, homeChain, oracleIDToP2pID),
	}

	for _, n := range nodes {
		// All nodes have issue reading the latest sequence number, should lead to empty outcomes
		n.ccipReader.On(
			"NextSeqNum",
			ctx,
			mock.Anything,
		).Return([]cciptypes.SeqNum{}, nil)
	}

	err = homeChain.Close()
	if err != nil {
		return nil
	}
	return nodes
}

func setupAllNodesReadAllChains(ctx context.Context, t *testing.T, lggr logger.Logger) []nodeSetup {
	cfg := cciptypes.CommitPluginConfig{
		DestChain:           chainC,
		PricedTokens:        []types.Account{tokenX},
		TokenPricesObserver: false,
		NewMsgScanBatchSize: 256,
	}

	chainConfigInfos := []reader.ChainConfigInfo{
		{
			ChainSelector: chainA,
			ChainConfig: reader.HomeChainConfigMapper{
				FChain: 1,
				Readers: []libocrtypes.PeerID{
					{1}, {2}, {3},
				},
				Config: []byte{0},
			},
		},
		{
			ChainSelector: chainB,
			ChainConfig: reader.HomeChainConfigMapper{
				FChain: 1,
				Readers: []libocrtypes.PeerID{
					{1}, {2}, {3},
				},
				Config: []byte{0},
			},
		},
		{
			ChainSelector: chainC,
			ChainConfig: reader.HomeChainConfigMapper{
				FChain: 1,
				Readers: []libocrtypes.PeerID{
					{1}, {2}, {3},
				},
				Config: []byte{0},
			},
		},
	}

	homeChain := setupHomeChainPoller(lggr, chainConfigInfos)
	err := homeChain.Start(ctx)
	if err != nil {
		return nil
	}
	oracleIDToP2pID := GetP2pIDs(1, 2, 3)
	n1 := newNode(ctx, t, lggr, 1, cfg, homeChain, oracleIDToP2pID)
	n2 := newNode(ctx, t, lggr, 2, cfg, homeChain, oracleIDToP2pID)
	n3 := newNode(ctx, t, lggr, 3, cfg, homeChain, oracleIDToP2pID)
	nodes := []nodeSetup{n1, n2, n3}

	for _, n := range nodes {
		// then they fetch new msgs, there is nothing new on chainA
		n.ccipReader.On(
			"MsgsBetweenSeqNums",
			ctx,
			chainA,
			cciptypes.NewSeqNumRange(11, cciptypes.SeqNum(11+cfg.NewMsgScanBatchSize)),
		).Return([]cciptypes.CCIPMsg{}, nil)

		// and there are two new message on chainB
		n.ccipReader.On(
			"MsgsBetweenSeqNums",
			ctx,
			chainB,
			cciptypes.NewSeqNumRange(21, cciptypes.SeqNum(21+cfg.NewMsgScanBatchSize)),
		).Return([]cciptypes.CCIPMsg{
			{CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{MsgHash: cciptypes.Bytes32{1}, ID: "1", SourceChain: chainB, SeqNum: 21}},
			{CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{MsgHash: cciptypes.Bytes32{2}, ID: "2", SourceChain: chainB, SeqNum: 22}},
		}, nil)

		n.ccipReader.On("GasPrices", ctx, []cciptypes.ChainSelector{chainA, chainB}).
			Return([]cciptypes.BigInt{
				cciptypes.NewBigIntFromInt64(1000),
				cciptypes.NewBigIntFromInt64(20_000),
			}, nil)

		// all nodes observe the same sequence numbers 10 for chainA and 20 for chainB
		n.ccipReader.On("NextSeqNum", ctx, []cciptypes.ChainSelector{chainA, chainB}).
			Return([]cciptypes.SeqNum{10, 20}, nil)

	}

	// No need to keep it running in the background anymore for this test
	err = homeChain.Close()
	if err != nil {
		return nil
	}

	return nodes
}

func setupNodesDoNotAgreeOnMsgs(ctx context.Context, t *testing.T, lggr logger.Logger) []nodeSetup {
	cfg := cciptypes.CommitPluginConfig{
		DestChain:           chainC,
		PricedTokens:        []types.Account{tokenX},
		TokenPricesObserver: false,
		NewMsgScanBatchSize: 256,
	}

	chainConfigInfos := []reader.ChainConfigInfo{
		{
			ChainSelector: chainA,
			ChainConfig: reader.HomeChainConfigMapper{
				FChain: 1,
				Readers: []libocrtypes.PeerID{
					{1}, {2}, {3},
				},
				Config: []byte{0},
			},
		},
		{
			ChainSelector: chainB,
			ChainConfig: reader.HomeChainConfigMapper{
				FChain: 1,
				Readers: []libocrtypes.PeerID{
					{1}, {2}, {3},
				},
				Config: []byte{0},
			},
		},
		{
			ChainSelector: chainC,
			ChainConfig: reader.HomeChainConfigMapper{
				FChain: 1,
				Readers: []libocrtypes.PeerID{
					{1}, {2}, {3},
				},
				Config: []byte{0},
			},
		},
	}

	homeChain := setupHomeChainPoller(lggr, chainConfigInfos)
	err := homeChain.Start(ctx)
	if err != nil {
		return nil
	}
	oracleIDToP2pID := GetP2pIDs(1, 2, 3)
	n1 := newNode(ctx, t, lggr, 1, cfg, homeChain, oracleIDToP2pID)
	n2 := newNode(ctx, t, lggr, 2, cfg, homeChain, oracleIDToP2pID)
	n3 := newNode(ctx, t, lggr, 3, cfg, homeChain, oracleIDToP2pID)
	nodes := []nodeSetup{n1, n2, n3}

	for i, n := range nodes {
		// all nodes observe the same sequence numbers 10 for chainA and 20 for chainB
		n.ccipReader.On("NextSeqNum", ctx, []cciptypes.ChainSelector{chainA, chainB}).
			Return([]cciptypes.SeqNum{10, 20}, nil)

		// then they fetch new msgs, there is nothing new on chainA
		n.ccipReader.On(
			"MsgsBetweenSeqNums",
			ctx,
			chainA,
			cciptypes.NewSeqNumRange(11, cciptypes.SeqNum(11+cfg.NewMsgScanBatchSize)),
		).Return([]cciptypes.CCIPMsg{}, nil)

		// and there are two new message on chainB
		n.ccipReader.On(
			"MsgsBetweenSeqNums",
			ctx,
			chainB,
			cciptypes.NewSeqNumRange(
				21,
				cciptypes.SeqNum(21+cfg.NewMsgScanBatchSize),
			),
		).Return([]cciptypes.CCIPMsg{
			{CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{MsgHash: cciptypes.Bytes32{1}, ID: "1" + strconv.Itoa(i), SourceChain: chainB, SeqNum: 21 + cciptypes.SeqNum(i*10)}},
			{CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{MsgHash: cciptypes.Bytes32{2}, ID: "2" + strconv.Itoa(i), SourceChain: chainB, SeqNum: 22 + cciptypes.SeqNum(i*20)}},
		}, nil)

		n.ccipReader.On("GasPrices", ctx, []cciptypes.ChainSelector{chainA, chainB}).
			Return([]cciptypes.BigInt{
				cciptypes.NewBigIntFromInt64(1000),
				cciptypes.NewBigIntFromInt64(20_000),
			}, nil)
	}

	// No need to keep it running in the background anymore for this test
	err = homeChain.Close()
	if err != nil {
		return nil
	}

	return nodes
}

type nodeSetup struct {
	node        *Plugin
	ccipReader  *mocks.CCIPReader
	priceReader *mocks.TokenPricesReader
	reportCodec *mocks.CommitPluginJSONReportCodec
	msgHasher   *mocks.MessageHasher
}

func newNode(ctx context.Context, t *testing.T, lggr logger.Logger, id int, cfg cciptypes.CommitPluginConfig, homeChain reader.HomeChain, oracleIDToP2pID map[commontypes.OracleID]libocrtypes.PeerID) nodeSetup {
	ccipReader := mocks.NewCCIPReader()
	priceReader := mocks.NewTokenPricesReader()
	reportCodec := mocks.NewCommitPluginJSONReportCodec()
	msgHasher := mocks.NewMessageHasher()

	node1 := NewPlugin(
		context.Background(),
		commontypes.OracleID(id),
		oracleIDToP2pID,
		cfg,
		ccipReader,
		priceReader,
		reportCodec,
		msgHasher,
		lggr,
		homeChain,
	)

	return nodeSetup{
		node:        node1,
		ccipReader:  ccipReader,
		priceReader: priceReader,
		reportCodec: reportCodec,
		msgHasher:   msgHasher,
	}
}

func setupHomeChainPoller(lggr logger.Logger, chainConfigInfos []reader.ChainConfigInfo) reader.HomeChain {
	homeChainReader := mocks.NewContractReaderMock()
	homeChainReader.On(
		"GetLatestValue", mock.Anything, "CCIPCapabilityConfiguration", "getAllChainConfigs", mock.Anything, mock.Anything).Run(
		func(args mock.Arguments) {
			arg := args.Get(4).(*[]reader.ChainConfigInfo)
			*arg = chainConfigInfos
		}).Return(nil)

	homeChain := reader.NewHomeChainConfigPoller(
		homeChainReader,
		lggr,
		10*time.Millisecond, // to prevent linting error because of logging after finishing tests, we close the poller after each test, having lower polling interval make it catch up faster
	)

	return homeChain
}
func GetP2pIDs(ids ...int) map[commontypes.OracleID]libocrtypes.PeerID {
	res := make(map[commontypes.OracleID]libocrtypes.PeerID)
	for _, id := range ids {
		res[commontypes.OracleID(id)] = libocrtypes.PeerID{byte(id)}
	}
	return res
}

var (
	chainA = cciptypes.ChainSelector(1)
	chainB = cciptypes.ChainSelector(2)
	chainC = cciptypes.ChainSelector(3)

	tokenX = types.Account("tk_xxx")
)
