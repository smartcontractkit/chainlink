package soltxm

import (
	"context"
	"crypto/rand"
	"sync"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
)

func TestTxCache(t *testing.T) {
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

	// init cache + store some signatures and cancelFunc
	cache := NewTxCache()
	n := 5
	for i := 0; i < n; i++ {
		sig, cancel := newProcess(i)
		err := cache.Insert(sig, cancel)
		assert.NoError(t, err)
	}

	// return list of signatures
	list := cache.List()
	assert.Equal(t, n, len(list))

	// stop all sub processes
	for i := 0; i < len(list); i++ {
		cache.Cancel(list[i])
		assert.Equal(t, n-i-1, len(cache.List()))
	}
	wg.Wait()
}
