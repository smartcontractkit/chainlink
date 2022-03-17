package client_test

import (
	"context"
	"math/big"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	promtestutil "github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/mock"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func TestPool_Dial(t *testing.T) {
	tests := []struct {
		name            string
		poolChainID     *big.Int
		nodeChainID     int64
		sendNodeChainID int64
		nodes           []chainIDResps
		sendNodes       []chainIDResp
		errStr          string
	}{
		{
			name:            "no nodes",
			poolChainID:     testutils.FixtureChainID,
			nodeChainID:     testutils.FixtureChainID.Int64(),
			sendNodeChainID: testutils.FixtureChainID.Int64(),
			nodes:           []chainIDResps{},
			sendNodes:       []chainIDResp{},
			errStr:          "no available nodes for chain 0",
		},
		{
			name:            "normal",
			poolChainID:     testutils.FixtureChainID,
			nodeChainID:     testutils.FixtureChainID.Int64(),
			sendNodeChainID: testutils.FixtureChainID.Int64(),
			nodes: []chainIDResps{
				{ws: chainIDResp{testutils.FixtureChainID.Int64(), nil}},
			},
			sendNodes: []chainIDResp{
				{testutils.FixtureChainID.Int64(), nil},
			},
		},
		{
			name:            "node has wrong chain ID compared to pool",
			poolChainID:     testutils.FixtureChainID,
			nodeChainID:     42,
			sendNodeChainID: testutils.FixtureChainID.Int64(),
			nodes: []chainIDResps{
				{ws: chainIDResp{1, nil}},
			},
			sendNodes: []chainIDResp{
				{1, nil},
			},
			errStr: "has chain ID 42 which does not match pool chain ID of 0",
		},
		{
			name:            "sendonly node has wrong chain ID compared to pool",
			poolChainID:     testutils.FixtureChainID,
			nodeChainID:     testutils.FixtureChainID.Int64(),
			sendNodeChainID: 42,
			nodes: []chainIDResps{
				{ws: chainIDResp{testutils.FixtureChainID.Int64(), nil}},
			},
			sendNodes: []chainIDResp{
				{testutils.FixtureChainID.Int64(), nil},
			},
			errStr: "has chain ID 42 which does not match pool chain ID of 0",
		},
		{
			name:            "remote RPC has wrong chain ID for primary node (ws) - no error, it will go into retry loop",
			poolChainID:     testutils.FixtureChainID,
			nodeChainID:     testutils.FixtureChainID.Int64(),
			sendNodeChainID: testutils.FixtureChainID.Int64(),
			nodes: []chainIDResps{
				{
					ws:   chainIDResp{42, nil},
					http: &chainIDResp{testutils.FixtureChainID.Int64(), nil},
				},
			},
			sendNodes: []chainIDResp{
				{testutils.FixtureChainID.Int64(), nil},
			},
		},
		{
			name:            "remote RPC has wrong chain ID for primary node (http) - no error, it will go into retry loop",
			poolChainID:     testutils.FixtureChainID,
			nodeChainID:     testutils.FixtureChainID.Int64(),
			sendNodeChainID: testutils.FixtureChainID.Int64(),
			nodes: []chainIDResps{
				{
					ws:   chainIDResp{testutils.FixtureChainID.Int64(), nil},
					http: &chainIDResp{42, nil},
				},
			},
			sendNodes: []chainIDResp{
				{testutils.FixtureChainID.Int64(), nil},
			},
		},
		{
			name:            "remote RPC has wrong chain ID for sendonly node",
			poolChainID:     testutils.FixtureChainID,
			nodeChainID:     testutils.FixtureChainID.Int64(),
			sendNodeChainID: testutils.FixtureChainID.Int64(),
			nodes: []chainIDResps{
				{ws: chainIDResp{testutils.FixtureChainID.Int64(), nil}},
			},
			sendNodes: []chainIDResp{
				{42, nil},
			},
			// TODO: Followup; sendonly nodes should not halt if they fail to
			// dail on startup; instead should go into retry loop like
			// primaries
			// See: https://app.shortcut.com/chainlinklabs/story/31338/sendonly-nodes-should-not-halt-node-boot-if-they-fail-to-dial-instead-should-have-retry-loop-like-primaries
			errStr: "sendonly rpc ChainID doesn't match local chain ID: RPC ID=42, local ID=0",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), cltest.WaitTimeout(t))
			defer cancel()

			nodes := make([]evmclient.Node, len(test.nodes))
			for i, n := range test.nodes {
				nodes[i] = n.newNode(t, test.nodeChainID)
			}
			sendNodes := make([]evmclient.SendOnlyNode, len(test.sendNodes))
			for i, n := range test.sendNodes {
				sendNodes[i] = n.newSendOnlyNode(t, test.sendNodeChainID)
			}
			p := evmclient.NewPool(logger.TestLogger(t), nodes, sendNodes, test.poolChainID)
			err := p.Dial(ctx)
			if test.errStr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.errStr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

type chainIDResp struct {
	chainID int64
	err     error
}

func (r *chainIDResp) newSendOnlyNode(t *testing.T, nodeChainID int64) evmclient.SendOnlyNode {
	httpURL := r.newHTTPServer(t)
	return evmclient.NewSendOnlyNode(logger.TestLogger(t), *httpURL, t.Name(), big.NewInt(nodeChainID))
}
func (r *chainIDResp) newHTTPServer(t *testing.T) *url.URL {
	rpcSrv := rpc.NewServer()
	t.Cleanup(rpcSrv.Stop)
	rpcSrv.RegisterName("eth", &chainIDService{*r})
	ts := httptest.NewServer(rpcSrv)
	t.Cleanup(ts.Close)

	httpURL, err := url.Parse(ts.URL)
	require.NoError(t, err)
	return httpURL
}

type chainIDResps struct {
	ws   chainIDResp
	http *chainIDResp
	id   int32
}

func (r *chainIDResps) newNode(t *testing.T, nodeChainID int64) evmclient.Node {
	ws := cltest.NewWSServer(t, big.NewInt(r.ws.chainID), func(method string, params gjson.Result) (string, string) {
		t.Errorf("Unexpected method call: %s(%s)", method, params)
		return "", ""
	})

	wsURL, err := url.Parse(ws)
	require.NoError(t, err)

	var httpURL *url.URL
	if r.http != nil {
		httpURL = r.http.newHTTPServer(t)
	}

	defer func() { r.id++ }()
	return evmclient.NewNode(evmclient.TestNodeConfig{}, logger.TestLogger(t), *wsURL, httpURL, t.Name(), r.id, big.NewInt(nodeChainID))
}

type chainIDService struct {
	chainIDResp
}

func (x *chainIDService) ChainId(ctx context.Context) (*hexutil.Big, error) {
	if x.err != nil {
		return nil, x.err
	}
	return (*hexutil.Big)(big.NewInt(x.chainID)), nil
}

func newPool(t *testing.T, nodes []evmclient.Node) *evmclient.Pool {
	return evmclient.NewPool(logger.TestLogger(t), nodes, []evmclient.SendOnlyNode{}, &cltest.FixtureChainID)
}

func TestUnit_Pool_RunLoop(t *testing.T) {
	t.Run("reports node states to prometheus", func(t *testing.T) {
		n1 := new(evmmocks.Node)
		n1.Test(t)
		n2 := new(evmmocks.Node)
		n2.Test(t)
		n3 := new(evmmocks.Node)
		n3.Test(t)
		nodes := []evmclient.Node{n1, n2, n3}

		lggr, observedLogs := logger.TestLoggerObserved(t, zap.ErrorLevel)
		p := evmclient.NewPool(lggr, nodes, []evmclient.SendOnlyNode{}, &cltest.FixtureChainID)

		n1.On("String").Maybe().Return("n1")
		n2.On("String").Maybe().Return("n2")
		n3.On("String").Maybe().Return("n3")

		n1.On("Close").Maybe()
		n2.On("Close").Maybe()
		n3.On("Close").Maybe()

		// n1 is alive
		n1.On("Start", mock.Anything).Return(nil).Once()
		n1.On("Verify", mock.Anything, &cltest.FixtureChainID).Return(nil).Once()
		n1.On("State").Return(evmclient.NodeStateAlive)
		n1.On("ChainID").Return(testutils.FixtureChainID).Once()
		// n2 is unreachable
		n2.On("Start", mock.Anything).Return(nil).Once()
		n2.On("Verify", mock.Anything, &cltest.FixtureChainID).Return(nil).Once()
		n2.On("State").Return(evmclient.NodeStateUnreachable)
		n2.On("ChainID").Return(testutils.FixtureChainID).Once()
		// n3 is out of sync
		n3.On("Start", mock.Anything).Return(nil).Once()
		n3.On("Verify", mock.Anything, &cltest.FixtureChainID).Return(nil).Once()
		n3.On("State").Return(evmclient.NodeStateOutOfSync)
		n3.On("ChainID").Return(testutils.FixtureChainID).Once()

		require.NoError(t, p.Dial(context.Background()))
		defer p.Close()

		testutils.WaitForLogMessage(t, observedLogs, "At least one EVM primary node is dead")

		testutils.AssertEventually(t, func() bool {
			totalReported := promtestutil.CollectAndCount(evmclient.PromEVMPoolRPCNodeStates)
			if totalReported < 3 {
				return false
			}
			if promtestutil.ToFloat64(evmclient.PromEVMPoolRPCNodeStates.WithLabelValues("0", "Alive")) < 1.0 {
				return false
			}
			if promtestutil.ToFloat64(evmclient.PromEVMPoolRPCNodeStates.WithLabelValues("0", "Unreachable")) < 1.0 {
				return false
			}
			if promtestutil.ToFloat64(evmclient.PromEVMPoolRPCNodeStates.WithLabelValues("0", "OutOfSync")) < 1.0 {
				return false
			}
			return true
		})
	})

}

func TestUnit_Pool_BatchCallContextAll(t *testing.T) {
	var mockNodes []*evmmocks.Node
	var mockSendonlys []*evmmocks.SendOnlyNode
	var nodes []evmclient.Node
	var sendonlys []evmclient.SendOnlyNode

	nodeCount := 2
	sendOnlyCount := 3

	b := []rpc.BatchElem{
		rpc.BatchElem{Method: "method", Args: []interface{}{1, false}},
		rpc.BatchElem{Method: "method2"},
	}

	ctx := testutils.Context(t)

	for i := 0; i < nodeCount; i++ {
		node := new(evmmocks.Node)
		node.On("State").Return(evmclient.NodeStateAlive)
		node.Test(t)
		node.On("BatchCallContext", ctx, b).Return(nil).Once()
		nodes = append(nodes, node)
		mockNodes = append(mockNodes, node)
	}
	for i := 0; i < sendOnlyCount; i++ {
		s := new(evmmocks.SendOnlyNode)
		s.Test(t)
		s.On("BatchCallContext", ctx, b).Return(nil).Once()
		sendonlys = append(sendonlys, s)
		mockSendonlys = append(mockSendonlys, s)
	}

	p := evmclient.NewPool(logger.TestLogger(t), nodes, sendonlys, &cltest.FixtureChainID)

	p.BatchCallContextAll(ctx, b)

	for _, n := range mockNodes {
		n.AssertExpectations(t)
	}
	for _, s := range mockSendonlys {
		s.AssertExpectations(t)
	}
}
