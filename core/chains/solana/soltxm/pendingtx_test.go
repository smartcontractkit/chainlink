package soltxm

import (
	"context"
	"crypto/rand"
	"sync"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
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
		processCtx, cancel := context.WithCancel(ctx)
		wg.Add(1)
		go func() {
			<-processCtx.Done()
			wg.Done()
		}()
		return solana.SignatureFromBytes(sig), cancel
	}

	// init inflight txs map + store some signatures and cancelFunc
	txs := newPendingTxContext()
	ids := map[solana.Signature]uuid.UUID{}
	n := 5
	for i := 0; i < n; i++ {
		sig, cancel := newProcess(i)
		id, err := txs.New(sig, cancel)
		assert.NoError(t, err)
		ids[sig] = id
	}

	// cannot add signature for non existent ID
	require.Error(t, txs.Add(uuid.New(), solana.Signature{}))

	// return list of signatures
	list := txs.ListAll()
	assert.Equal(t, n, len(list))

	// stop all sub processes
	for i := 0; i < len(list); i++ {
		id := txs.Remove(list[i])
		assert.Equal(t, n-i-1, len(txs.ListAll()))
		assert.Equal(t, ids[list[i]], id)

		// second remove should not return valid id - already removed
		assert.Equal(t, uuid.Nil, txs.Remove(list[i]))
	}
	wg.Wait()
}

func TestPendingTxContext_expired(t *testing.T) {
	_, cancel := context.WithCancel(testutils.Context(t))
	sig := solana.Signature{}
	txs := newPendingTxContext()

	id, err := txs.New(sig, cancel)
	assert.NoError(t, err)

	assert.True(t, txs.Expired(sig, 0*time.Second))   // expired for 0s lifetime
	assert.False(t, txs.Expired(sig, 60*time.Second)) // not expired for 60s lifetime

	assert.Equal(t, id, txs.Remove(sig))
	assert.False(t, txs.Expired(sig, 60*time.Second)) // no longer exists, should return false
}

func TestPendingTxContext_race(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		txCtx := newPendingTxContext()
		var wg sync.WaitGroup
		wg.Add(2)
		var err [2]error

		go func() {
			_, err[0] = txCtx.New(solana.Signature{}, func() {})
			wg.Done()
		}()
		go func() {
			_, err[1] = txCtx.New(solana.Signature{}, func() {})
			wg.Done()
		}()

		wg.Wait()
		assert.True(t, (err[0] != nil && err[1] == nil) || (err[0] == nil && err[1] != nil), "one and only one 'add' should have errored")
	})

	t.Run("add", func(t *testing.T) {
		txCtx := newPendingTxContext()
		id, createErr := txCtx.New(solana.Signature{}, func() {})
		require.NoError(t, createErr)
		var wg sync.WaitGroup
		wg.Add(2)
		var err [2]error

		go func() {
			err[0] = txCtx.Add(id, solana.Signature{1})
			wg.Done()
		}()
		go func() {
			err[1] = txCtx.Add(id, solana.Signature{1})
			wg.Done()
		}()

		wg.Wait()
		assert.True(t, (err[0] != nil && err[1] == nil) || (err[0] == nil && err[1] != nil), "one and only one 'add' should have errored")
	})

	t.Run("remove", func(t *testing.T) {
		txCtx := newPendingTxContext()
		_, err := txCtx.New(solana.Signature{}, func() {})
		require.NoError(t, err)
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
