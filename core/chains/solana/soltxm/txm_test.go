package soltxm_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	solanaGo "github.com/gagliardetto/solana-go"
	solanaClient "github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	"github.com/smartcontractkit/chainlink/core/chains/solana/soltxm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTxm_HappyPath(t *testing.T) {
	received := make(chan struct{})

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out := fmt.Sprintf(`{"jsonrpc":"2.0","result":"5vwvmiymecC16SA6nt1YNNNrMEvMX4YpCJ6Et7N3M9bM3YFVNQJsThHnZPsA8AwnmTJ9LR2u8CUsVXkQ5vnMkKEo","id":1}`)
		_, err := w.Write([]byte(out))
		require.NoError(t, err)
		close(received)
	}))
	defer mockServer.Close()

	// set up txm
	lggr := logger.TestLogger(t)
	cfg := config.NewConfig(db.ChainCfg{}, lggr)
	client, err := solanaClient.NewClient(mockServer.URL, cfg, 2*time.Second, lggr)
	require.NoError(t, err)
	getClient := func() (solanaClient.ReaderWriter, error) {
		return client, nil
	}
	txm := soltxm.NewTxm(getClient, cfg, lggr)

	// start
	require.NoError(t, txm.Start(context.Background()))

	// already started
	assert.Error(t, txm.Start(context.Background()))

	// submit tx
	tx := solanaGo.Transaction{}
	tx.Signatures = append(tx.Signatures, solanaGo.Signature{})
	tx.Message.Header.NumRequiredSignatures = 1
	assert.NoError(t, txm.Enqueue("testAccount", &tx))
	<-received
}
