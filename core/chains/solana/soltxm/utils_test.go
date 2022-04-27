package soltxm

import (
	"context"
	"crypto/rand"
	"sync"
	"testing"

	"github.com/gagliardetto/solana-go"
	solanaClient "github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/stretchr/testify/assert"
)

func TestTxProcesses(t *testing.T) {
	// setup
	var wg sync.WaitGroup
	ctx := context.Background()
	newProcess := func(i int) (solana.Signature, context.CancelFunc) {
		// make random signature
		sig := make([]byte, 64)
		rand.Read(sig)

		// start subprocess to wait for context
		ctx, cancel := context.WithCancel(ctx)
		wg.Add(1)
		go func() {
			<-ctx.Done()
			wg.Done()
		}()
		return solana.SignatureFromBytes(sig), cancel
	}

	// init inflight txs map + store some signatures and cancelFunc
	txs := NewTxProcesses("test")
	n := 5
	for i := 0; i < n; i++ {
		sig, cancel := newProcess(i)
		err := txs.Insert(sig, cancel)
		assert.NoError(t, err)
	}

	// return list of signatures
	list := txs.List()
	assert.Equal(t, n, len(list))

	// stop all sub processes
	for i := 0; i < len(list); i++ {
		txs.Cancel(list[i])
		assert.Equal(t, n-i-1, len(txs.List()))
	}
	wg.Wait()
}

func TestValidClient(t *testing.T) {
	called := false
	tc := func() (solanaClient.ReaderWriter, error) {
		return &solanaClient.Client{}, nil
	}
	callOnce := func() (solanaClient.ReaderWriter, error) {
		assert.False(t, called)
		called = true
		return tc()
	}

	// Get should only request a client once, use cached afterward
	t.Run("get", func(t *testing.T) {
		c := NewValidClient(callOnce)
		rw, err := c.Get()
		assert.NoError(t, err)
		assert.NotNil(t, rw)
		assert.NotNil(t, c.client)

		// used cached client
		rw, err = c.Get()
		assert.NoError(t, err)
		assert.NotNil(t, rw)
	})

	// Clear removes the cached client
	t.Run("clear", func(t *testing.T) {
		c := NewValidClient(tc)
		c.Get()
		assert.NotNil(t, c.client)

		c.Clear()
		assert.Nil(t, c.client)
	})

	// Race checks a race condition of Getting and Clearing a new client
	t.Run("race", func(t *testing.T) {
		c := NewValidClient(tc)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			rw, err := c.Get()
			assert.NoError(t, err)
			assert.NotNil(t, rw)
			wg.Done()
		}()
		go func() {
			c.Clear()
			wg.Done()
		}()
		wg.Wait()
	})
}
