package solana

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	"github.com/smartcontractkit/chainlink/core/chains/solana/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSolanaChain_GetClient(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// devnet gensis hash
		out := `{"jsonrpc":"2.0","result":"EtWTRABZaYq6iMfeYKouRu166VU2xqa1wcaWoxPkrZBG","id":1}`
		if strings.Contains(r.URL.Path, "/mismatch") {
			out = `{"jsonrpc":"2.0","result":"5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d","id":1}`
		}

		_, err := w.Write([]byte(out))
		require.NoError(t, err)
	}))
	defer mockServer.Close()

	solORM := new(mocks.ORM)
	lggr := logger.TestLogger(t)
	testChain := chain{
		id:            "devnet",
		orm:           solORM,
		cfg:           config.NewConfig(db.ChainCfg{}, lggr),
		lggr:          logger.TestLogger(t),
		clientChainID: map[string]string{},
		chainIDLock:   &sync.RWMutex{},
	}

	// named nodes response (happy path)
	solORM.On("NodeNamed", mock.Anything).Return(db.Node{
		SolanaChainID: "devnet",
		SolanaURL:     mockServer.URL + "/0",
	}, nil).Once()
	_, err := testChain.getClient("namedNode")
	assert.NoError(t, err)

  // random nodes (happy path, all valid)
  solORM.On("NodesForChain", mock.Anything, mock.Anything, mock.Anything).Return([]db.Node{
		db.Node{
			SolanaChainID: "devnet",
			SolanaURL:     mockServer.URL + "/1",
		},
		db.Node{
			SolanaChainID: "devnet",
			SolanaURL:     mockServer.URL + "/2",
		},
	}, 2, nil).Once()
	_, err = testChain.getClient("")
	assert.NoError(t, err)

  // random nodes (happy path, 1 valid + multiple invalid)
  solORM.On("NodesForChain", mock.Anything, mock.Anything, mock.Anything).Return([]db.Node{
    db.Node{
      SolanaChainID: "devnet",
      SolanaURL:     mockServer.URL + "/A",
    },
    db.Node{
      SolanaChainID: "devnet",
      SolanaURL:     mockServer.URL + "/mismatch/A",
    },
    db.Node{
      SolanaChainID: "devnet",
      SolanaURL:     mockServer.URL + "/mismatch/B",
    },
    db.Node{
      SolanaChainID: "devnet",
      SolanaURL:     mockServer.URL + "/mismatch/C",
    },
    db.Node{
      SolanaChainID: "devnet",
      SolanaURL:     mockServer.URL + "/mismatch/D",
    },
  }, 2, nil).Once()
  _, err = testChain.getClient("")
  assert.NoError(t, err)

	// empty nodes response
	solORM.On("NodesForChain", mock.Anything, mock.Anything, mock.Anything).Return([]db.Node{}, 0, nil).Once()
	_, err = testChain.getClient("")
	assert.Error(t, err)

	// named nodes response wrong genesis hash
	solORM.On("NodeNamed", mock.Anything).Return(db.Node{
		SolanaChainID: "devnet",
		SolanaURL:     mockServer.URL + "/mismatch/0",
	}, nil).Once()
	_, err = testChain.getClient("namedNode")
	assert.Error(t, err)

	// no valid nodes to select from
	solORM.On("NodesForChain", mock.Anything, mock.Anything, mock.Anything).Return([]db.Node{
		db.Node{
			SolanaChainID: "devnet",
			SolanaURL:     mockServer.URL + "/mismatch/1",
		},
		db.Node{
			SolanaChainID: "devnet",
			SolanaURL:     mockServer.URL + "/mismatch/2",
		},
	}, 2, nil).Once()
	_, err = testChain.getClient("")
	assert.Error(t, err)
}

func TestSolanaChain_VerifiedClient(t *testing.T) {
	called := false
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// should only be called once, chainID will be cached in chain
		if called {
			assert.NoError(t, errors.New("rpc has been called once already"))
		}

		// devnet gensis hash
		out := `{"jsonrpc":"2.0","result":"EtWTRABZaYq6iMfeYKouRu166VU2xqa1wcaWoxPkrZBG","id":1}`
		_, err := w.Write([]byte(out))
		require.NoError(t, err)
		called = true
	}))
	defer mockServer.Close()

	lggr := logger.TestLogger(t)
	testChain := chain{
		cfg:           config.NewConfig(db.ChainCfg{}, lggr),
		lggr:          logger.TestLogger(t),
		clientChainID: map[string]string{},
		chainIDLock:   &sync.RWMutex{},
	}
	node := db.Node{SolanaURL: mockServer.URL}

	// happy path
	testChain.id = "devnet"
	_, err := testChain.verifiedClient(node)
	assert.NoError(t, err)
	assert.Equal(t, testChain.id, testChain.clientChainID[node.SolanaURL])

	// expect error from id mismatch
	testChain.id = "incorrect"
	_, err = testChain.verifiedClient(node)
	assert.Error(t, err)
}
