package soltxm

import (
	"crypto/rand"
	"encoding/binary"
	"math"
	"sync"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
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

func XXXGetFeePrice(t *testing.T, tx *solana.Transaction) uint64 {
	require.True(t, len(tx.Message.Instructions) > 0, "not enough instructions")
	require.Equal(t, 9, len(tx.Message.Instructions[0].Data), "fee instruction should be first") // 1 byte function selector, 8 byte little-endian encoded uint64

	return binary.LittleEndian.Uint64([]byte(tx.Message.Instructions[0].Data)[1:])
}

// Test doubling progression of fee
func TestPendingTx_FeeBumping(t *testing.T) {
	tx := PendingTx{baseTx: &solana.Transaction{}}
	n := 10
	init := true

	for i := 0; i < n; i++ {
		// initial tx should use the default price arg
		txWithFee, fee, err := tx.SetComputeUnitPrice(0, 0, 10_000)
		require.NoError(t, err)

		if init {
			v := uint64(0)
			assert.Equal(t, v, fee)
			assert.Equal(t, v, XXXGetFeePrice(t, txWithFee))
			init = false
		} else {
			v := uint64(math.Pow(2, float64(i-1)))
			assert.Equal(t, v, fee)
			assert.Equal(t, v, XXXGetFeePrice(t, txWithFee))
		}

		// if tx has been broadcast should begin X^2 increases
		tx.broadcast = true
		tx.currentFee = fee // track current fee
	}
}

// Test combination of inputs for robustness
func FuzzPendingTx_SetComputeUnitPrice(f *testing.F) {
	f.Add(uint64(2), uint64(0), uint64(0), uint64(0), true)
	f.Add(uint64(0), uint64(0), uint64(0), uint64(10), false)
	f.Add(uint64(1), uint64(100), uint64(0), uint64(1000), true)
	f.Add(uint64(0), uint64(0), uint64(10), uint64(0), false)
	f.Add(uint64(10), uint64(0), uint64(0), uint64(10), true)
	f.Add(uint64(100), uint64(0), uint64(0), uint64(0), false)

	f.Fuzz(func(t *testing.T, init, base, min, max uint64, broadcast bool) {
		tx := PendingTx{
			baseTx:     &solana.Transaction{},
			currentFee: init,
			broadcast:  broadcast,
		}

		txWithFee, fee, err := tx.SetComputeUnitPrice(base, min, max)

		// if parameters are out of bounds, should error
		if base < min || base > max {
			assert.Error(t, err)
			return
		}

		assert.Equal(t, fee, XXXGetFeePrice(t, txWithFee), "tx fee + output fee should match")
		assert.True(t, fee >= min, "fee should be bounded by minimum")
		assert.True(t, fee <= max, "fee should be bounded by maximum")

		if !broadcast {
			assert.Equal(t, base, fee)
		} else {
			if init == 0 {
				assert.True(t, 1 == fee || fee == max || fee == min, "if starting at 0 => doubling = 1, bounded by min & max")
			} else {
				assert.True(t, 2*init == fee || fee == max || fee == min, "double initial fee or bounded by min & max")
			}
		}
	})
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
			assert.NoError(t, txs.Add(idTemp, sigTemp, 0))

			// validate get method
			txGet, exists := txs.GetBySignature(sigTemp)
			assert.True(t, exists)
			assert.Equal(t, idTemp, txGet.id)
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

			_, exists := txs.GetBySignature(list[i])
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
			err := txs.Add(id, sigTemp, 0)
			assert.NoError(t, err)

			// validate get method
			txGet, exists := txs.GetBySignature(sigTemp)
			assert.True(t, exists)
			assert.Equal(t, id, txGet.id)
			assert.Equal(t, i+1, len(txGet.signatures))
			assert.Equal(t, sigTemp, txGet.signatures[len(txGet.signatures)-1])
		}

		// return list of signatures
		list := txs.ListSignatures()
		assert.Equal(t, n, len(list))

		// clear transaction by completing 1 signature
		txs.OnSuccess(list[0])
		assert.Equal(t, 0, len(txs.ListSignatures()))
		for i := 0; i < len(list); i++ {
			_, exists := txs.GetBySignature(list[i])
			assert.False(t, exists)
		}
	})

	t.Run("duplicateSignatures", func(t *testing.T) {
		txs := newPendingTxMemory()
		id := txs.New(PendingTx{})

		// duplicate for same tx
		sig := XXXNewSignature(t)
		assert.NoError(t, txs.Add(id, sig, 0))
		assert.Error(t, txs.Add(id, sig, 0))

		// duplicate for different txs
		assert.Error(t, txs.Add(txs.New(PendingTx{}), sig, 0))
	})

	t.Run("zeroID_zeroSignature", func(t *testing.T) {
		txs := newPendingTxMemory()
		id := txs.New(PendingTx{})
		assert.True(t, id != uuid.Nil)

		assert.Error(t, txs.Add(id, solana.Signature{}, 0))
		assert.Error(t, txs.Add(uuid.Nil, XXXNewSignature(t), 0))
	})
}

func TestPendingTxMemory_race(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		txCtx := newPendingTxMemory()
		id := txCtx.New(PendingTx{})
		sig := XXXNewSignature(t)
		var wg sync.WaitGroup
		wg.Add(2)
		var err [2]error

		go func() {
			err[0] = txCtx.Add(id, sig, 0)
			wg.Done()
		}()
		go func() {
			err[1] = txCtx.Add(id, sig, 0)
			wg.Done()
		}()

		wg.Wait()
		assert.True(t, (err[0] != nil && err[1] == nil) || (err[0] == nil && err[1] != nil), "one and only one 'add' should have errored")
	})

	t.Run("remove", func(t *testing.T) {
		txCtx := newPendingTxMemory()
		id := txCtx.New(PendingTx{})
		var wg sync.WaitGroup
		wg.Add(2)
		var err [2]error

		go func() {
			err[0] = txCtx.Remove(id)
			wg.Done()
		}()
		go func() {
			err[1] = txCtx.Remove(id)
			wg.Done()
		}()

		wg.Wait()
		assert.True(t, (err[0] != nil && err[1] == nil) || (err[0] == nil && err[1] != nil), "one and only one 'add' should have errored")
	})
}
