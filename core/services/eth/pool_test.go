package eth_test

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

func Test_Pool_AddNode(t *testing.T) {
	l := logger.TestLogger(t)
	pool := eth.NewPool(l, []eth.Node{}, []eth.SendOnlyNode{}, big.NewInt(32))

	resp := `
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": "0x20"
}
`
	_, wsUrl, wsCleanup := cltest.NewWSServer(resp, func(data []byte) {
		req := cltest.ParseJSON(t, bytes.NewReader(data))
		require.Equal(t, "eth_chainId", req.Get("method").String())
	})
	t.Cleanup(wsCleanup)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		defer r.Body.Close()
		req := cltest.ParseJSON(t, bytes.NewReader(data))
		require.Equal(t, "eth_chainId", req.Get("method").String())
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, resp)
	}))
	t.Cleanup(s.Close)

	t.Run("errors with unstarted pool", func(t *testing.T) {
		n := eth.NewNode(l, *cltest.MustURL(t, wsUrl), nil, "test primary node", big.NewInt(32))
		err := pool.AddNode(context.Background(), n)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot add node; pool is not started")

		s := eth.NewSendOnlyNode(l, *cltest.MustURL(t, s.URL), "test sendonly node", big.NewInt(32))
		err = pool.AddSendOnlyNode(context.Background(), s)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot add send only node; pool is not started")
	})

	eth.StartPool(pool)

	t.Run("adding node with nil chain ID is not allowed", func(t *testing.T) {
		n := eth.NewNode(l, *cltest.MustURL(t, wsUrl), nil, "test primary node", nil)
		err := pool.AddNode(context.Background(), n)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot add node with chain ID <nil> to pool with chain ID 32")

		s := eth.NewSendOnlyNode(l, *cltest.MustURL(t, s.URL), "test send only node", nil)
		err = pool.AddSendOnlyNode(context.Background(), s)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot add send only node with chain ID <nil> to pool with chain ID 32")
	})

	t.Run("adding node with wrong chain ID is not allowed", func(t *testing.T) {
		n := eth.NewNode(l, *cltest.MustURL(t, wsUrl), nil, "test primary node", big.NewInt(42))
		err := pool.AddNode(context.Background(), n)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot add node with chain ID 42 to pool with chain ID 32")

		s := eth.NewSendOnlyNode(l, *cltest.MustURL(t, s.URL), "test send only node", big.NewInt(42))
		err = pool.AddSendOnlyNode(context.Background(), s)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot add send only node with chain ID 42 to pool with chain ID 32")
	})

	t.Run("adds nodes with correct chain ID", func(t *testing.T) {
		n := eth.NewNode(l, *cltest.MustURL(t, wsUrl), nil, "test primary node", big.NewInt(32))

		err := pool.AddNode(context.Background(), n)
		require.NoError(t, err)

		nodes := eth.Nodes(pool)
		assert.Len(t, nodes, 1)

		s := eth.NewSendOnlyNode(l, *cltest.MustURL(t, s.URL), "test send only node", big.NewInt(32))

		err = pool.AddSendOnlyNode(context.Background(), s)
		require.NoError(t, err)

		sendonlys := eth.SendOnlyNodes(pool)
		assert.Len(t, sendonlys, 1)
	})
}

func Test_Pool_RemoveNode(t *testing.T) {
	t.Fatal("todo")
}
