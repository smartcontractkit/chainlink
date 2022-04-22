package soltxm

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/client/mocks"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type soltxmProm struct {
	id                              string
	success, revert, fail, timedOut float64
}

func (p soltxmProm) assertEqual(t *testing.T) {
	assert.Equal(t, p.success, testutil.ToFloat64(promSolTxmSuccessfulTxs.WithLabelValues(p.id)), "mismatch: success")
	assert.Equal(t, p.revert, testutil.ToFloat64(promSolTxmRevertedTxs.WithLabelValues(p.id)), "mismatch: revert")
	assert.Equal(t, p.fail, testutil.ToFloat64(promSolTxmFailedTxs.WithLabelValues(p.id)), "mismatch: fail")
	assert.Equal(t, p.timedOut, testutil.ToFloat64(promSolTxmTimedOutTxs.WithLabelValues(p.id)), "mismatch: timedOut")
}

func (p soltxmProm) getInflight() float64 {
	return testutil.ToFloat64(promSolTxmInflightTxs.WithLabelValues(p.id))
}

func TestTxm(t *testing.T) {
	// set up configs needed in txm
	id := "mocknet"
	lggr := logger.TestLogger(t)
	cfg := config.NewConfig(db.ChainCfg{}, lggr)
	mc := new(mocks.ReaderWriter)
	txm := NewTxm(id, func() (client.ReaderWriter, error) {
		return mc, nil
	}, cfg, lggr)
	require.NoError(t, txm.Start(context.Background()))

	// tracking prom metrics
	prom := soltxmProm{id: id}

	// create random signature
	getSig := func() solana.Signature {
		sig := make([]byte, 64)
		rand.Read(sig)
		return solana.SignatureFromBytes(sig)
	}

	// create placeholder transaction
	getTx := func() *solana.Transaction {
		// create transfer tx
		key := solana.PublicKey{}
		tx, err := solana.NewTransaction(
			[]solana.Instruction{
				system.NewTransferInstruction(
					rand.Uint64(), // identifier is the transfer balance
					key,
					key,
				).Build(),
			},
			solana.Hash{},
			solana.TransactionPayer(key),
		)
		assert.NoError(t, err)
		return tx
	}

	// check if cached transaction is cleared
	empty := func() bool {
		count := txm.InflightTxs()
		assert.Equal(t, float64(count), prom.getInflight()) // validate prom metric and cache length
		return count == 0
	}

	waitFor := func(f func() bool) {
		for i := 0; i < 30; i++ {
			if f() {
				return
			}
			time.Sleep(time.Second)
		}
		assert.NoError(t, errors.New("unable to confirm tx cache is empty"))
	}

	// happy path (send => simulate success => tx: nil => tx: processed => tx: confirmed => done)
	t.Run("happyPath", func(t *testing.T) {
		sig := getSig()
		tx := getTx()
		var wg sync.WaitGroup
		wg.Add(3)

		sendCount := 0
		mc.On("SendTx", mock.Anything, tx).Run(func(mock.Arguments) {
			sendCount++
		}).Return(sig, nil)
		mc.On("SimulateTx", mock.Anything, tx, mock.Anything).Return(&rpc.SimulateTransactionResult{}, nil).Once()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Run(func(mock.Arguments) {
			wg.Done()
		}).Return([]*rpc.SignatureStatusesResult{nil}, nil).Once()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Run(func(mock.Arguments) {
			wg.Done()
		}).Return([]*rpc.SignatureStatusesResult{&rpc.SignatureStatusesResult{
			ConfirmationStatus: rpc.ConfirmationStatusProcessed,
		}}, nil).Once()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Run(func(mock.Arguments) {
			wg.Done()
		}).Return([]*rpc.SignatureStatusesResult{&rpc.SignatureStatusesResult{
			ConfirmationStatus: rpc.ConfirmationStatusConfirmed,
		}}, nil).Once()

		// send tx
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait()

		// no transactions stored cache list
		waitFor(empty)
		// transaction should be sent more than twice
		t.Logf("sendTx received %d calls", sendCount)
		assert.Greater(t, sendCount, 2)

		// panic if sendTx called after context cancelled
		mc.On("SendTx", mock.Anything, tx).Panic("SendTx should not be called anymore")

		// check prom metric
		prom.success++
		prom.assertEqual(t)
	})

	// fail on initial transmit (RPC immediate rejects)
	t.Run("fail_initialTx", func(t *testing.T) {
		tx := getTx()
		var wg sync.WaitGroup
		wg.Add(1)

		// should only be called once (tx does not start retry, confirming, or simulation)
		mc.On("SendTx", mock.Anything, tx).Run(func(mock.Arguments) {
			wg.Done()
		}).Return(solana.Signature{}, errors.New("FAIL")).Once()

		// tx should be able to queue
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait() // wait to be picked up and processed

		// no transactions stored cache list
		waitFor(empty)

		// check prom metric
		prom.fail++
		prom.assertEqual(t)
	})

	// tx fails simulation (simulation error)
	t.Run("fail_simulation", func(t *testing.T) {
		tx := getTx()
		sig := getSig()
		var wg sync.WaitGroup
		wg.Add(1)

		mc.On("SendTx", mock.Anything, tx).Return(sig, nil)
		mc.On("SimulateTx", mock.Anything, tx, mock.Anything).Run(func(mock.Arguments) {
			wg.Done()
		}).Return(&rpc.SimulateTransactionResult{
			Err: "FAIL",
		}, nil).Once()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Return([]*rpc.SignatureStatusesResult{nil}, nil).Maybe()

		// tx should be able to queue
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait()      // wait to be picked up and processed
		waitFor(empty) // tx cache cleared quickly

		// check prom metric
		prom.fail++
		prom.assertEqual(t)
	})

	// tx fails simulation (rpc error, timeout should clean up b/c sig status will be nil)
	t.Run("fail_simulation_confirmNil", func(t *testing.T) {
		tx := getTx()
		sig := getSig()
		var wg sync.WaitGroup
		wg.Add(1)

		mc.On("SendTx", mock.Anything, tx).Return(sig, nil)
		mc.On("SimulateTx", mock.Anything, tx, mock.Anything).Run(func(mock.Arguments) {
			wg.Done()
		}).Return(&rpc.SimulateTransactionResult{}, errors.New("FAIL")).Once()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Return([]*rpc.SignatureStatusesResult{nil}, nil)

		// tx should be able to queue
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait()      // wait to be picked up and processed
		waitFor(empty) // tx cache cleared after timeout

		// check prom metric
		prom.timedOut++
		prom.assertEqual(t)

		// panic if sendTx called after context cancelled
		mc.On("SendTx", mock.Anything, tx).Panic("SendTx should not be called anymore")
	})

	// tx passes sim, never passes processed (timeout should cleanup)
	t.Run("fail_confirm_processed", func(t *testing.T) {
		tx := getTx()
		sig := getSig()
		var wg sync.WaitGroup
		wg.Add(1)

		mc.On("SendTx", mock.Anything, tx).Return(sig, nil)
		mc.On("SimulateTx", mock.Anything, tx, mock.Anything).Run(func(mock.Arguments) {
			wg.Done()
		}).Return(&rpc.SimulateTransactionResult{}, nil).Once()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Return([]*rpc.SignatureStatusesResult{&rpc.SignatureStatusesResult{
			ConfirmationStatus: rpc.ConfirmationStatusProcessed,
		}}, nil)

		// tx should be able to queue
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait()      // wait to be picked up and processed
		waitFor(empty) // tx cache cleared after timeout

		// check prom metric
		prom.timedOut++
		prom.assertEqual(t)

		// panic if sendTx called after context cancelled
		mc.On("SendTx", mock.Anything, tx).Panic("SendTx should not be called anymore")
	})

	// tx passes sim, shows processed, moves to nil (timeout should cleanup)
	t.Run("fail_confirm_processedToNil", func(t *testing.T) {
		tx := getTx()
		sig := getSig()
		var wg sync.WaitGroup
		wg.Add(1)

		mc.On("SendTx", mock.Anything, tx).Return(sig, nil)
		mc.On("SimulateTx", mock.Anything, tx, mock.Anything).Run(func(mock.Arguments) {
			wg.Done()
		}).Return(&rpc.SimulateTransactionResult{}, nil).Once()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Return([]*rpc.SignatureStatusesResult{&rpc.SignatureStatusesResult{
			ConfirmationStatus: rpc.ConfirmationStatusProcessed,
		}}, nil).Twice()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Return([]*rpc.SignatureStatusesResult{nil}, nil)

		// tx should be able to queue
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait()      // wait to be picked up and processed
		waitFor(empty) // tx cache cleared after timeout

		// check prom metric
		prom.timedOut++
		prom.assertEqual(t)

		// panic if sendTx called after context cancelled
		mc.On("SendTx", mock.Anything, tx).Panic("SendTx should not be called anymore")
	})

	// tx passes sim, errors on confirm
	t.Run("fail_confirm_revert", func(t *testing.T) {
		tx := getTx()
		sig := getSig()
		var wg sync.WaitGroup
		wg.Add(1)

		mc.On("SendTx", mock.Anything, tx).Return(sig, nil)
		mc.On("SimulateTx", mock.Anything, tx, mock.Anything).Run(func(mock.Arguments) {
			wg.Done()
		}).Return(&rpc.SimulateTransactionResult{}, nil).Once()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Return([]*rpc.SignatureStatusesResult{&rpc.SignatureStatusesResult{
			ConfirmationStatus: rpc.ConfirmationStatusProcessed,
			Err:                "ERROR",
		}}, nil).Once()

		// tx should be able to queue
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait()      // wait to be picked up and processed
		waitFor(empty) // tx cache cleared after timeout

		// check prom metric
		prom.revert++
		prom.assertEqual(t)

		// panic if sendTx called after context cancelled
		mc.On("SendTx", mock.Anything, tx).Panic("SendTx should not be called anymore")
	})

	mc.AssertExpectations(t)
}
