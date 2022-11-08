package soltxm

import (
	"context"
	"crypto/rand"
	"sync"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
)

func TestPendingTxContext(t *testing.T) {
	// setup
	var wg sync.WaitGroup
	ctx := testutils.Context(t)
	newProcess := func(i int) (solana.Signature, context.CancelFunc) {
		// make random signature
		sig := make([]byte, 64)
		_, err := rand.Read(sig)
		require.NoError(t, err)

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
	txs := newPendingTxContext()
	n := 5
	for i := 0; i < n; i++ {
		sig, cancel := newProcess(i)
		err := txs.Add(sig, cancel)
		assert.NoError(t, err)
	}

	// return list of signatures
	list := txs.ListAll()
	assert.Equal(t, n, len(list))

	// stop all sub processes
	for i := 0; i < len(list); i++ {
		txs.Remove(list[i])
		assert.Equal(t, n-i-1, len(txs.ListAll()))
	}
	wg.Wait()
}

func TestPendingTxContext_expired(t *testing.T) {
	_, cancel := context.WithCancel(testutils.Context(t))
	sig := solana.Signature{}
	txs := newPendingTxContext()

	err := txs.Add(sig, cancel)
	assert.NoError(t, err)

	assert.True(t, txs.Expired(sig, 0*time.Second))   // expired for 0s lifetime
	assert.False(t, txs.Expired(sig, 60*time.Second)) // not expired for 60s lifetime

	txs.Remove(sig)
	assert.True(t, txs.Expired(sig, 60*time.Second)) // no longer exists, should be expired
}

func TestPendingTxContext_race(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		txCtx := newPendingTxContext()
		var wg sync.WaitGroup
		wg.Add(2)
		var err [2]error

		go func() {
			err[0] = txCtx.Add(solana.Signature{}, func() {})
			wg.Done()
		}()
		go func() {
			err[1] = txCtx.Add(solana.Signature{}, func() {})
			wg.Done()
		}()

		wg.Wait()
		assert.True(t, (err[0] != nil && err[1] == nil) || (err[0] == nil && err[1] != nil), "one and only one 'add' should have errored")
	})

	t.Run("remove", func(t *testing.T) {
		txCtx := newPendingTxContext()
		require.NoError(t, txCtx.Add(solana.Signature{}, func() {}))
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			assert.NotPanics(t, func() { txCtx.Remove(solana.Signature{}) })
			wg.Done()
		}()
		go func() {
			assert.NotPanics(t, func() { txCtx.Remove(solana.Signature{}) })
			wg.Done()
		}()

		wg.Wait()
	})
}
