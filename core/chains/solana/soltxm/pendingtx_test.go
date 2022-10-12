package soltxm

import (
	"crypto/rand"
	"fmt"
	"sync"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testing only
func XXXNewSignature(t *testing.T) solana.Signature {
	// make random signature
	sig := make([]byte, 64)
	_, err := rand.Read(sig)
	require.NoError(t, err)

	return solana.SignatureFromBytes(sig)
}

func TestPendingTxMemory(t *testing.T) {
	t.Run("happyPath", func(t *testing.T) {
		// init inflight txs map + store some signatures and cancelFunc
		txs := newPendingTxMemory()
		n := 5
		for i := 0; i < n; i++ {
			// 1 tx, 1 signature
			idTemp := txs.New(PendingTx{})
			sigTemp := XXXNewSignature(t)
			assert.NoError(t, txs.Add(idTemp, sigTemp))

			// validate get method
			idGet, txGet, exists := txs.Get(sigTemp)
			assert.True(t, exists)
			assert.Equal(t, idTemp, idGet)
			assert.Equal(t, 1, len(txGet.signatures))
			assert.Equal(t, sigTemp, txGet.signatures[0])
		}

		// return list of signatures
		list := txs.ListSignatures()
		assert.Equal(t, n, len(list))

		// stop all sub processes
		for i := 0; i < len(list); i++ {
			txs.OnSuccess(list[i])
			assert.Equal(t, n-i-1, len(txs.ListSignatures()))

			_, _, exists := txs.Get(list[i])
			assert.False(t, exists)
		}
	})

	t.Run("oneTxManySig", func(t *testing.T) {
		// init inflight txs map + store some signatures and cancelFunc
		txs := newPendingTxMemory()
		n := 5
		var tx0 PendingTx
		id := txs.New(tx0) // store 1 tx
		for i := 0; i < n; i++ {
			// 1 tx, many signatures
			sigTemp := XXXNewSignature(t)
			err := txs.Add(id, sigTemp)
			assert.NoError(t, err)

			// validate get method
			idGet, txGet, exists := txs.Get(sigTemp)
			assert.True(t, exists)
			assert.Equal(t, id, idGet)
			assert.Equal(t, i+1, len(txGet.signatures))
			assert.Equal(t, sigTemp, txGet.signatures[len(txGet.signatures)-1])
		}

		// return list of signatures
		list := txs.ListSignatures()
		assert.Equal(t, n, len(list))

		// stop all sub processes by completing 1 signature
		txs.OnSuccess(list[0])
		assert.Equal(t, 0, len(txs.ListSignatures()))
		for i := 0; i < len(list); i++ {
			_, _, exists := txs.Get(list[i])
			assert.False(t, exists)
		}
	})

	t.Run("duplicateSignatures", func(t *testing.T) {
		// TODO
		// duplicate for same tx
		// duplicate for different txs
	})

	t.Run("stateMachine", func(t *testing.T) {
		// TODO
	})
}

func TestPendingTxMemory_race(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		txCtx := newPendingTxMemory()
		id := txCtx.New(PendingTx{})
		var wg sync.WaitGroup
		wg.Add(2)
		var err [2]error

		go func() {
			err[0] = txCtx.Add(id, solana.Signature{})
			wg.Done()
		}()
		go func() {
			err[1] = txCtx.Add(id, solana.Signature{})
			wg.Done()
		}()

		wg.Wait()
		fmt.Println(err)
		assert.True(t, (err[0] != nil && err[1] == nil) || (err[0] == nil && err[1] != nil), "one and only one 'add' should have errored")
	})

	t.Run("remove", func(t *testing.T) {
		txCtx := newPendingTxMemory()
		id := txCtx.New(PendingTx{})
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			assert.NotPanics(t, func() { txCtx.Remove(id) })
			wg.Done()
		}()
		go func() {
			assert.NotPanics(t, func() { txCtx.Remove(id) })
			wg.Done()
		}()

		wg.Wait()
	})
}
