package eth_test

import (
	"context"
	"errors"
	"math/big"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/logger"
	. "github.com/smartcontractkit/chainlink/core/services/eth"
)

func TestPool_Dial(t *testing.T) {
	tests := []struct {
		name        string
		presetID    *big.Int
		nodes       []chainIDResps
		sendNodes   []chainIDResp
		expectErr   bool
		multiErrCnt int
	}{
		{
			name: "normal",
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
		{
			name:      "wrong id",
			nodes:     []chainIDResps{{ws: chainIDResp{1, nil}}},
			sendNodes: []chainIDResp{{2, nil}},
			expectErr: true,
		},
		{
			name:      "wrong id preset",
			presetID:  big.NewInt(1),
			nodes:     []chainIDResps{{ws: chainIDResp{1, nil}}},
			sendNodes: []chainIDResp{{2, nil}},
			expectErr: true,
		},
		{
			name:     "wrong id preset multiple",
			presetID: big.NewInt(1),
			nodes: []chainIDResps{
				{ws: chainIDResp{1, nil}, http: &chainIDResp{2, nil}},
				{ws: chainIDResp{3, nil}, http: &chainIDResp{1, nil}},
			},
			sendNodes: []chainIDResp{
				{2, nil},
				{6, nil},
			},
			expectErr:   true,
			multiErrCnt: 4,
		},
		{
			name:      "error",
			nodes:     []chainIDResps{{ws: chainIDResp{1, nil}}},
			sendNodes: []chainIDResp{{-1, errors.New("fake")}},
			expectErr: true,
		},
		{
			name:      "error preset",
			presetID:  big.NewInt(1),
			nodes:     []chainIDResps{{ws: chainIDResp{1, nil}}},
			sendNodes: []chainIDResp{{-1, errors.New("fake")}},
			expectErr: true,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), cltest.DefaultWaitTimeout)
			defer cancel()

			nodes := make([]Node, len(test.nodes))
			for i, n := range test.nodes {
				nodes[i] = n.newNode(t)
			}
			sendNodes := make([]SendOnlyNode, len(test.sendNodes))
			for i, n := range test.sendNodes {
				sendNodes[i] = n.newSendOnlyNode(t)
			}
			p := NewPool(logger.TestLogger(t), nodes, sendNodes, test.presetID)
			if err := p.Dial(ctx); err != nil {
				if test.expectErr {
					if test.multiErrCnt > 0 {
						assert.Equal(t, test.multiErrCnt, len(multierr.Errors(err)))
					}
				} else {
					t.Error(err)
				}
			} else if test.expectErr {
				t.Error("expected error")
			}
		})
	}
}

type chainIDResp struct {
	chainID int64
	err     error
}

func (r *chainIDResp) newSendOnlyNode(t *testing.T) SendOnlyNode {
	httpURL := r.newHTTPServer(t)
	return NewSendOnlyNode(logger.TestLogger(t), *httpURL, t.Name())
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

func (r *chainIDResps) newNode(t *testing.T) Node {
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

	return NewNode(logger.TestLogger(t), *wsURL, httpURL, t.Name())
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
