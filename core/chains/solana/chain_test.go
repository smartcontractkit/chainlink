package solana

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

const TestSolanaGenesisHashTemplate = `{"jsonrpc":"2.0","result":"%s","id":1}`

func TestSolanaChain_GetClient(t *testing.T) {
	checkOnce := map[string]struct{}{}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out := fmt.Sprintf(TestSolanaGenesisHashTemplate, client.MainnetGenesisHash) // mainnet genesis hash

		if !strings.Contains(r.URL.Path, "/mismatch") {
			// devnet gensis hash
			out = fmt.Sprintf(TestSolanaGenesisHashTemplate, client.DevnetGenesisHash)

			// clients with correct chainID should request chainID only once
			if _, exists := checkOnce[r.URL.Path]; exists {
				assert.NoError(t, errors.Errorf("rpc has been called once already for successful client '%s'", r.URL.Path))
			}
			checkOnce[r.URL.Path] = struct{}{}
		}

		_, err := w.Write([]byte(out))
		require.NoError(t, err)
	}))
	defer mockServer.Close()

	solORM := &mockORM{}
	lggr := logger.TestLogger(t)
	testChain := chain{
		id:          "devnet",
		orm:         solORM,
		cfg:         config.NewConfig(db.ChainCfg{}, lggr),
		lggr:        logger.TestLogger(t),
		clientCache: map[string]*verifiedCachedClient{},
	}

	// random nodes (happy path, all valid)
	solORM.nodesForChain = []db.Node{
		{
			SolanaChainID: "devnet",
			SolanaURL:     mockServer.URL + "/1",
		},
		{
			SolanaChainID: "devnet",
			SolanaURL:     mockServer.URL + "/2",
		},
	}
	_, err := testChain.getClient()
	assert.NoError(t, err)

	// random nodes (happy path, 1 valid + multiple invalid)
	solORM.nodesForChain = []db.Node{
		{
			SolanaChainID: "devnet",
			SolanaURL:     mockServer.URL + "/1",
		},
		{
			SolanaChainID: "devnet",
			SolanaURL:     mockServer.URL + "/mismatch/1",
		},
		{
			SolanaChainID: "devnet",
			SolanaURL:     mockServer.URL + "/mismatch/2",
		},
		{
			SolanaChainID: "devnet",
			SolanaURL:     mockServer.URL + "/mismatch/3",
		},
		{
			SolanaChainID: "devnet",
			SolanaURL:     mockServer.URL + "/mismatch/4",
		},
	}
	_, err = testChain.getClient()
	assert.NoError(t, err)

	// empty nodes response
	solORM.nodesForChain = nil
	_, err = testChain.getClient()
	assert.Error(t, err)

	// no valid nodes to select from
	solORM.nodesForChain = []db.Node{
		{
			SolanaChainID: "devnet",
			SolanaURL:     mockServer.URL + "/mismatch/1",
		},
		{
			SolanaChainID: "devnet",
			SolanaURL:     mockServer.URL + "/mismatch/2",
		},
	}
	_, err = testChain.getClient()
	assert.NoError(t, err)
}

func TestSolanaChain_VerifiedClient(t *testing.T) {
	called := false
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out := `{ "jsonrpc": "2.0", "result": 1234, "id": 1 }` // getSlot response

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		// handle getGenesisHash request
		if strings.Contains(string(body), "getGenesisHash") {
			// should only be called once, chainID will be cached in chain
			// allowing `mismatch` to be ignored, since invalid nodes will try to verify the chain ID
			// if it is not verified
			if !strings.Contains(r.URL.Path, "/mismatch") && called {
				assert.NoError(t, errors.New("rpc has been called once already"))
			}
			// devnet genesis hash
			out = fmt.Sprintf(TestSolanaGenesisHashTemplate, client.DevnetGenesisHash)
		}
		_, err = w.Write([]byte(out))
		require.NoError(t, err)
		called = true
	}))
	defer mockServer.Close()

	lggr := logger.TestLogger(t)
	testChain := chain{
		cfg:         config.NewConfig(db.ChainCfg{}, lggr),
		lggr:        logger.TestLogger(t),
		clientCache: map[string]*verifiedCachedClient{},
	}
	node := db.Node{SolanaURL: mockServer.URL}

	// happy path
	testChain.id = "devnet"
	_, err := testChain.verifiedClient(node)
	assert.NoError(t, err)

	// retrieve cached client and retrieve slot height
	c, err := testChain.verifiedClient(node)
	assert.NoError(t, err)
	slot, err := c.SlotHeight()
	assert.NoError(t, err)
	assert.Equal(t, uint64(1234), slot)

	node.SolanaURL = mockServer.URL + "/mismatch"
	testChain.id = "incorrect"
	c, err = testChain.verifiedClient(node)
	assert.NoError(t, err)
	_, err = c.ChainID()
	// expect error from id mismatch (even if using a cached client) when performing RPC calls
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("client returned mismatched chain id (expected: %s, got: %s): %s", "incorrect", "devnet", node.SolanaURL), err.Error())
}

func TestSolanaChain_VerifiedClient_ParallelClients(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out := fmt.Sprintf(TestSolanaGenesisHashTemplate, client.DevnetGenesisHash)
		_, err := w.Write([]byte(out))
		require.NoError(t, err)
	}))
	defer mockServer.Close()

	lggr := logger.TestLogger(t)
	testChain := chain{
		id:          "devnet",
		cfg:         config.NewConfig(db.ChainCfg{}, lggr),
		lggr:        logger.TestLogger(t),
		clientCache: map[string]*verifiedCachedClient{},
	}
	node := db.Node{SolanaURL: mockServer.URL}

	var wg sync.WaitGroup
	wg.Add(2)

	var client0 client.ReaderWriter
	var client1 client.ReaderWriter
	var err0 error
	var err1 error

	// call verifiedClient in parallel
	go func() {
		client0, err0 = testChain.verifiedClient(node)
		assert.NoError(t, err0)
		wg.Done()
	}()
	go func() {
		client1, err1 = testChain.verifiedClient(node)
		assert.NoError(t, err1)
		wg.Done()
	}()

	wg.Wait()

	// check if pointers are all the same
	assert.Equal(t, testChain.clientCache[mockServer.URL], client0)
	assert.Equal(t, testChain.clientCache[mockServer.URL], client1)
}

var _ ORM = &mockORM{}

type mockORM struct {
	nodesForChain []db.Node
}

func (m *mockORM) GetChainsByIDs(ids []string) (chains []DBChain, err error) {
	panic("implement me")
}

func (m *mockORM) GetNodesByChainIDs(chainIDs []string, qopts ...pg.QOpt) (nodes []db.Node, err error) {
	panic("implement me")
}

func (m *mockORM) StoreString(chainID string, key, val string) error {
	panic("implement me")
}

func (m *mockORM) Clear(chainID string, key string) error {
	panic("implement me")
}

func (m *mockORM) NodesForChain(chainID string, offset, limit int, qopts ...pg.QOpt) (nodes []db.Node, count int, err error) {
	return m.nodesForChain, len(m.nodesForChain), nil
}

func (m *mockORM) Chain(s string, opt ...pg.QOpt) (DBChain, error) { panic("unimplemented") }

func (m *mockORM) Chains(offset, limit int, qopts ...pg.QOpt) ([]DBChain, int, error) {
	panic("unimplemented")
}

func (m *mockORM) CreateChain(id string, config *db.ChainCfg, qopts ...pg.QOpt) (DBChain, error) {
	panic("unimplemented")
}

func (m *mockORM) UpdateChain(id string, enabled bool, config *db.ChainCfg, qopts ...pg.QOpt) (DBChain, error) {
	panic("unimplemented")
}

func (m *mockORM) DeleteChain(id string, qopts ...pg.QOpt) error { panic("unimplemented") }

func (m *mockORM) EnabledChains(opt ...pg.QOpt) ([]DBChain, error) { panic("unimplemented") }

func (m *mockORM) CreateNode(node db.Node, opt ...pg.QOpt) (db.Node, error) {
	panic("unimplemented")
}

func (m *mockORM) DeleteNode(i int32, opt ...pg.QOpt) error { panic("unimplemented") }

func (m *mockORM) NodeNamed(s string, opt ...pg.QOpt) (db.Node, error) { panic("unimplemented") }

func (m *mockORM) Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []db.Node, count int, err error) {
	panic("unimplemented")
}

func (m *mockORM) SetupNodes([]db.Node, []string) error    { panic("unimplemented") }
func (m *mockORM) EnsureChains([]string, ...pg.QOpt) error { panic("unimplemented") }
