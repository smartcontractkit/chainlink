package cre

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	crecrypto "github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
)

func TestAptosAccountForNode_UsesMetadataKeyWithoutCallingNodeAPI(t *testing.T) {
	t.Parallel()

	var hits atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		http.NotFound(w, r)
	}))
	t.Cleanup(server.Close)

	expected, err := crecrypto.NormalizeAptosAccount("0x1")
	require.NoError(t, err)

	node := &Node{
		Name: "node-1",
		Keys: &secrets.NodeKeys{
			Aptos: &crecrypto.AptosKey{Account: expected},
		},
		Clients: NodeClients{
			RestClient: &clclient.ChainlinkClient{APIClient: resty.New().SetBaseURL(server.URL)},
		},
	}

	account, err := aptosAccountForNode(context.Background(), node)
	require.NoError(t, err)
	require.Equal(t, expected, account)
	require.Zero(t, hits.Load(), "node API must not be called when metadata already has the Aptos key")
}

func TestAptosAccountForNode_FallsBackToNodeAPIAndCachesKey(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/keys/aptos" {
			t.Errorf("unexpected path: got %q want %q", r.URL.Path, "/v2/keys/aptos")
			http.Error(w, fmt.Sprintf("unexpected path %q", r.URL.Path), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"data":[{"attributes":{"account":"0x1","publicKey":"0xabc123"}}]}`))
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	node := &Node{
		Name: "node-1",
		Keys: &secrets.NodeKeys{},
		Clients: NodeClients{
			RestClient: &clclient.ChainlinkClient{APIClient: resty.New().SetBaseURL(server.URL)},
		},
	}

	account, err := aptosAccountForNode(context.Background(), node)
	require.NoError(t, err)

	expected, err := crecrypto.NormalizeAptosAccount("0x1")
	require.NoError(t, err)
	require.Equal(t, expected, account)
	require.NotNil(t, node.Keys.Aptos)
	require.Equal(t, expected, node.Keys.Aptos.Account)
	require.Equal(t, "0xabc123", node.Keys.Aptos.PublicKey)
}
