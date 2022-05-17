package soltxm

import (
	"context"
	"crypto/rand"
	"sync"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
)

func TestPendingTxContext(t *testing.T) {
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
	_, cancel := context.WithCancel(context.Background())
	sig := solana.Signature{}
	txs := newPendingTxContext()

	err := txs.Add(sig, cancel)
	assert.NoError(t, err)

	assert.True(t, txs.Expired(sig, 0*time.Second))   // expired for 0s lifetime
	assert.False(t, txs.Expired(sig, 60*time.Second)) // not expired for 60s lifetime

	txs.Remove(sig)
	assert.True(t, txs.Expired(sig, 60*time.Second)) // no longer exists, should be expired
}
