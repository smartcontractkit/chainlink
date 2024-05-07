package client_test

import (
	"context"
	"math/big"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	promtestutil "github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

type poolConfig struct {
	selectionMode       string
	noNewHeadsThreshold time.Duration
	leaseDuration       time.Duration
}

func (c poolConfig) NodeSelectionMode() string {
	return c.selectionMode
}

func (c poolConfig) NodeNoNewHeadsThreshold() time.Duration {
	return c.noNewHeadsThreshold
}

func (c poolConfig) LeaseDuration() time.Duration {
	return c.leaseDuration
}

var defaultConfig evmclient.PoolConfig = &poolConfig{
	selectionMode:       evmclient.NodeSelectionMode_RoundRobin,
	noNewHeadsThreshold: 0,
	leaseDuration:       time.Second * 0,
}

func TestPool_Dial(t *testing.T) {
	t.Parallel()

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
			name:            "remote RPC has wrong chain ID for sendonly node - no error, it will go into retry loop",
			poolChainID:     testutils.FixtureChainID,
			nodeChainID:     testutils.FixtureChainID.Int64(),
			sendNodeChainID: testutils.FixtureChainID.Int64(),
			nodes: []chainIDResps{
				{ws: chainIDResp{testutils.FixtureChainID.Int64(), nil}},
			},
			sendNodes: []chainIDResp{
				{42, nil},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx := testutils.Context(t)

			nodes := make([]evmclient.Node, len(test.nodes))
			for i, n := range test.nodes {
				nodes[i] = n.newNode(t, test.nodeChainID)
			}
			sendNodes := make([]evmclient.SendOnlyNode, len(test.sendNodes))
			for i, n := range test.sendNodes {
				sendNodes[i] = n.newSendOnlyNode(t, test.sendNodeChainID)
			}
			p := evmclient.NewPool(logger.Test(t), defaultConfig.NodeSelectionMode(), defaultConfig.LeaseDuration(), time.Second*0, nodes, sendNodes, test.poolChainID, "")
			err := p.Dial(ctx)
			if err == nil {
				t.Cleanup(func() { assert.NoError(t, p.Close()) })
			}
			assert.True(t, p.ChainType().IsValid())
			assert.False(t, p.ChainType().IsL2())
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
	return evmclient.NewSendOnlyNode(logger.Test(t), *httpURL, t.Name(), big.NewInt(nodeChainID))
}

func (r *chainIDResp) newHTTPServer(t *testing.T) *url.URL {
	rpcSrv := rpc.NewServer()
	t.Cleanup(rpcSrv.Stop)
	err := rpcSrv.RegisterName("eth", &chainIDService{*r})
	require.NoError(t, err)
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
	ws := testutils.NewWSServer(t, big.NewInt(r.ws.chainID), func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
		switch method {
		case "eth_subscribe":
			resp.Result = `"0x00"`
			resp.Notify = headResult
			return
		case "eth_unsubscribe":
			resp.Result = "true"
			return
		}
		t.Errorf("Unexpected method call: %s(%s)", method, params)
		return
	}).WSURL().String()

	wsURL, err := url.Parse(ws)
	require.NoError(t, err)

	var httpURL *url.URL
	if r.http != nil {
		httpURL = r.http.newHTTPServer(t)
	}

	defer func() { r.id++ }()
	return evmclient.NewNode(evmclient.TestNodePoolConfig{}, time.Second*0, logger.Test(t), *wsURL, httpURL, t.Name(), r.id, big.NewInt(nodeChainID), 0)
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

func TestUnit_Pool_RunLoop(t *testing.T) {
	t.Parallel()

	n1 := evmmocks.NewNode(t)
	n2 := evmmocks.NewNode(t)
	n3 := evmmocks.NewNode(t)
	nodes := []evmclient.Node{n1, n2, n3}

	lggr, observedLogs := logger.TestObserved(t, zap.ErrorLevel)
	p := evmclient.NewPool(lggr, defaultConfig.NodeSelectionMode(), defaultConfig.LeaseDuration(), time.Second*0, nodes, []evmclient.SendOnlyNode{}, &cltest.FixtureChainID, "")

	n1.On("String").Maybe().Return("n1")
	n2.On("String").Maybe().Return("n2")
	n3.On("String").Maybe().Return("n3")

	n1.On("Close").Maybe().Return(nil)
	n2.On("Close").Maybe().Return(nil)
	n3.On("Close").Maybe().Return(nil)

	// n1 is alive
	n1.On("Start", mock.Anything).Return(nil).Once()
	n1.On("State").Return(evmclient.NodeStateAlive)
	n1.On("ChainID").Return(testutils.FixtureChainID).Once()
	// n2 is unreachable
	n2.On("Start", mock.Anything).Return(nil).Once()
	n2.On("State").Return(evmclient.NodeStateUnreachable)
	n2.On("ChainID").Return(testutils.FixtureChainID).Once()
	// n3 is out of sync
	n3.On("Start", mock.Anything).Return(nil).Once()
	n3.On("State").Return(evmclient.NodeStateOutOfSync)
	n3.On("ChainID").Return(testutils.FixtureChainID).Once()

	require.NoError(t, p.Dial(testutils.Context(t)))
	t.Cleanup(func() { assert.NoError(t, p.Close()) })

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
}

func TestUnit_Pool_BatchCallContextAll(t *testing.T) {
	t.Parallel()

	var nodes []evmclient.Node
	var sendonlys []evmclient.SendOnlyNode

	nodeCount := 2
	sendOnlyCount := 3

	b := []rpc.BatchElem{
		{Method: "method", Args: []interface{}{1, false}},
		{Method: "method2"},
	}

	ctx := testutils.Context(t)

	for i := 0; i < nodeCount; i++ {
		node := evmmocks.NewNode(t)
		node.On("State").Return(evmclient.NodeStateAlive).Maybe()
		node.On("BatchCallContext", ctx, b).Return(nil).Once()
		nodes = append(nodes, node)
	}
	for i := 0; i < sendOnlyCount; i++ {
		s := evmmocks.NewSendOnlyNode(t)
		s.On("BatchCallContext", ctx, b).Return(nil).Once()
		sendonlys = append(sendonlys, s)
	}

	p := evmclient.NewPool(logger.Test(t), defaultConfig.NodeSelectionMode(), defaultConfig.LeaseDuration(), time.Second*0, nodes, sendonlys, &cltest.FixtureChainID, "")

	assert.True(t, p.ChainType().IsValid())
	assert.False(t, p.ChainType().IsL2())
	require.NoError(t, p.BatchCallContextAll(ctx, b))
}

func TestUnit_Pool_LeaseDuration(t *testing.T) {
	t.Parallel()

	n1 := evmmocks.NewNode(t)
	n2 := evmmocks.NewNode(t)
	nodes := []evmclient.Node{n1, n2}
	type nodeStateSwitch struct {
		isAlive bool
		mu      sync.RWMutex
	}

	nodeSwitch := nodeStateSwitch{
		isAlive: true,
		mu:      sync.RWMutex{},
	}

	n1.On("String").Maybe().Return("n1")
	n2.On("String").Maybe().Return("n2")
	n1.On("Close").Maybe().Return(nil)
	n2.On("Close").Maybe().Return(nil)
	n2.On("UnsubscribeAllExceptAliveLoop").Return()
	n2.On("SubscribersCount").Return(int32(2))

	n1.On("Start", mock.Anything).Return(nil).Once()
	n1.On("State").Return(func() evmclient.NodeState {
		nodeSwitch.mu.RLock()
		defer nodeSwitch.mu.RUnlock()
		if nodeSwitch.isAlive {
			return evmclient.NodeStateAlive
		}
		return evmclient.NodeStateOutOfSync
	})
	n1.On("Order").Return(int32(1))
	n1.On("ChainID").Return(testutils.FixtureChainID).Once()

	n2.On("Start", mock.Anything).Return(nil).Once()
	n2.On("State").Return(evmclient.NodeStateAlive)
	n2.On("Order").Return(int32(2))
	n2.On("ChainID").Return(testutils.FixtureChainID).Once()

	lggr, observedLogs := logger.TestObserved(t, zap.InfoLevel)
	p := evmclient.NewPool(lggr, "PriorityLevel", time.Second*2, time.Second*0, nodes, []evmclient.SendOnlyNode{}, &cltest.FixtureChainID, "")
	require.NoError(t, p.Dial(testutils.Context(t)))
	t.Cleanup(func() { assert.NoError(t, p.Close()) })

	testutils.WaitForLogMessage(t, observedLogs, "The pool will switch to best node every 2s")
	nodeSwitch.mu.Lock()
	nodeSwitch.isAlive = false
	nodeSwitch.mu.Unlock()
	testutils.WaitForLogMessage(t, observedLogs, "At least one EVM primary node is dead")
	nodeSwitch.mu.Lock()
	nodeSwitch.isAlive = true
	nodeSwitch.mu.Unlock()
	testutils.WaitForLogMessage(t, observedLogs, `Switching to best node from "n2" to "n1"`)
}
