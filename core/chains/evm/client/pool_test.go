package client_test

import (
	"context"
	"math/big"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmclientmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func TestPool_Dial(t *testing.T) {
	tests := []struct {
		name      string
		presetID  *big.Int
		nodes     []chainIDResps
		sendNodes []chainIDResp
		wantErr   bool
		errStr    string
	}{
		{
			name:      "no nodes",
			presetID:  &cltest.FixtureChainID,
			nodes:     []chainIDResps{},
			sendNodes: []chainIDResp{},
			wantErr:   true,
			errStr:    "no available nodes for chain 0",
		},
		{
			name:     "normal",
			presetID: &cltest.FixtureChainID,
			nodes: []chainIDResps{
				{ws: chainIDResp{1, nil}},
			},
			sendNodes: []chainIDResp{
				{1, nil},
			},
		},
		{
			name:     "normal preset",
			presetID: big.NewInt(1),
			nodes: []chainIDResps{
				{ws: chainIDResp{1, nil}},
			},
			sendNodes: []chainIDResp{
				{1, nil},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), cltest.WaitTimeout(t))
			defer cancel()

			nodes := make([]evmclient.Node, len(test.nodes))
			for i, n := range test.nodes {
				nodes[i] = n.newNode(t)
			}
			sendNodes := make([]evmclient.SendOnlyNode, len(test.sendNodes))
			for i, n := range test.sendNodes {
				sendNodes[i] = n.newSendOnlyNode(t)
			}
			p := evmclient.NewPool(logger.TestLogger(t), nodes, sendNodes, test.presetID)
			err := p.Dial(ctx)
			if test.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.errStr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPool_Dial_Errors(t *testing.T) {
	t.Run("starts and kicks off retry loop even if dial errors", func(t *testing.T) {
		node := new(evmclientmocks.Node)
		node.On("String").Return("node").Maybe()
		node.On("Close").Maybe()
		node.Test(t)
		nodes := []evmclient.Node{node}
		p := newPool(t, nodes)

		node.On("Dial", mock.Anything).Return(errors.New("error"))

		err := p.Dial(context.Background())
		require.NoError(t, err)

		p.Close()

		node.AssertExpectations(t)
	})

	t.Run("starts and kicks off retry loop even on verification errors", func(t *testing.T) {
		node := new(evmclientmocks.Node)
		node.On("String").Return("node").Maybe()
		node.On("Close").Maybe()
		node.Test(t)
		nodes := []evmclient.Node{node}
		p := newPool(t, nodes)

		node.On("Dial", mock.Anything).Return(nil)
		node.On("Verify", mock.Anything, &cltest.FixtureChainID).Return(errors.New("error"))

		err := p.Dial(context.Background())
		require.NoError(t, err)

		p.Close()

		node.AssertExpectations(t)
	})
}

type chainIDResp struct {
	chainID int64
	err     error
}

func (r *chainIDResp) newSendOnlyNode(t *testing.T) evmclient.SendOnlyNode {
	httpURL := r.newHTTPServer(t)
	return evmclient.NewSendOnlyNode(logger.TestLogger(t), *httpURL, t.Name())
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
}

func (r *chainIDResps) newNode(t *testing.T) evmclient.Node {
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

	return evmclient.NewNode(logger.TestLogger(t), *wsURL, httpURL, t.Name())
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

func TestPool_RunLoop(t *testing.T) {
	t.Run("with several nodes and different types of errors", func(t *testing.T) {
		n1 := new(evmclientmocks.Node)
		n1.Test(t)
		n2 := new(evmclientmocks.Node)
		n2.Test(t)
		n3 := new(evmclientmocks.Node)
		n3.Test(t)
		nodes := []evmclient.Node{n1, n2, n3}
		p := newPool(t, nodes)

		n1.On("String").Maybe().Return("n1")
		n2.On("String").Maybe().Return("n2")
		n3.On("String").Maybe().Return("n3")

		n1.On("Close").Maybe()
		n2.On("Close").Maybe()
		n3.On("Close").Maybe()

		wait := make(chan struct{})
		// n1 succeeds
		n1.On("Dial", mock.Anything).Return(nil).Once()
		n1.On("Verify", mock.Anything, &cltest.FixtureChainID).Return(nil).Once()
		n1.On("State").Return(evmclient.NodeStateAlive)
		// n2 fails once then succeeds in runloop
		n2.On("Dial", mock.Anything).Return(errors.New("first error")).Once()
		n2.On("State").Return(evmclient.NodeStateDead)
		// n3 succeeds dial then fails verification
		n3.On("Dial", mock.Anything).Return(nil).Once()
		n3.On("State").Return(evmclient.NodeStateDialed)
		n3.On("Verify", mock.Anything, &cltest.FixtureChainID).Return(errors.New("Verify error")).Once()
		n3.On("Verify", mock.Anything, &cltest.FixtureChainID).Once().Return(nil).Run(func(_ mock.Arguments) {
			close(wait)
		})

		// Handle spurious extra calls after
		n2.On("Dial", mock.Anything).Maybe().Return(nil)
		n3.On("Verify", mock.Anything, mock.Anything).Maybe().Return(nil)

		require.NoError(t, p.Dial(context.Background()))

		select {
		case <-wait:
		case <-time.After(cltest.WaitTimeout(t)):
			t.Fatal("timed out waiting for Dial call")
		}
		p.Close()

		n1.AssertExpectations(t)
		n2.AssertExpectations(t)
	})

}
