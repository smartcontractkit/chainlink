package soltxm

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/client/mocks"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/solkey"
	keyMocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
)

type soltxmProm struct {
	id                                                        string
	success, error, revert, reject, drop, simRevert, simOther float64
}

func (p soltxmProm) assertEqual(t *testing.T) {
	assert.Equal(t, p.success, testutil.ToFloat64(promSolTxmSuccessTxs.WithLabelValues(p.id)), "mismatch: success")
	assert.Equal(t, p.error, testutil.ToFloat64(promSolTxmErrorTxs.WithLabelValues(p.id)), "mismatch: error")
	assert.Equal(t, p.revert, testutil.ToFloat64(promSolTxmRevertTxs.WithLabelValues(p.id)), "mismatch: revert")
	assert.Equal(t, p.reject, testutil.ToFloat64(promSolTxmRejectTxs.WithLabelValues(p.id)), "mismatch: reject")
	assert.Equal(t, p.drop, testutil.ToFloat64(promSolTxmDropTxs.WithLabelValues(p.id)), "mismatch: drop")
	assert.Equal(t, p.simRevert, testutil.ToFloat64(promSolTxmSimRevertTxs.WithLabelValues(p.id)), "mismatch: simRevert")
	assert.Equal(t, p.simOther, testutil.ToFloat64(promSolTxmSimOtherTxs.WithLabelValues(p.id)), "mismatch: simOther")
}

func (p soltxmProm) getInflight() float64 {
	return testutil.ToFloat64(promSolTxmPendingTxs.WithLabelValues(p.id))
}

// create placeholder transaction
func getTx(t *testing.T, pubkey solana.PublicKey) *solana.Transaction {
	// create transfer tx
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				rand.Uint64(), // identifier is the transfer balance
				pubkey,
				pubkey,
			).Build(),
		},
		solana.Hash{},
		solana.TransactionPayer(pubkey),
	)
	assert.NoError(t, err)
	return tx
}

func newReaderWriterMock(t *testing.T) *mocks.ReaderWriter {
	m := new(mocks.ReaderWriter)
	m.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}

func TestTxm(t *testing.T) {
	// set up configs needed in txm
	id := "mocknet"
	lggr := logger.TestLogger(t)
	cfg := config.NewConfig(db.ChainCfg{}, lggr)
	mc := newReaderWriterMock(t)

	// mock solana keystore
	key, err := solkey.New()
	pubkey := key.PublicKey()

	require.NoError(t, err)
	mkey := keyMocks.NewSolana(t)
	mkey.On("Get", key.ID()).Return(key, nil)

	txm := NewTxm(id, func() (client.ReaderWriter, error) {
		return mc, nil
	}, cfg, mkey, lggr)
	require.NoError(t, txm.Start(testutils.Context(t)))

	// tracking prom metrics
	prom := soltxmProm{id: id}

	// create random signature
	getSig := func() solana.Signature {
		sig := make([]byte, 64)
		rand.Read(sig)
		return solana.SignatureFromBytes(sig)
	}

	// check if cached transaction is cleared
	empty := func() bool {
		count := txm.InflightTxs()
		assert.Equal(t, float64(count), prom.getInflight()) // validate prom metric and txs length
		return count == 0
	}

	// adjust wait time based on config
	waitDuration := cfg.TxConfirmTimeout()
	waitFor := func(f func() bool) {
		for i := 0; i < int(waitDuration.Seconds()*1.5); i++ {
			if f() {
				return
			}
			time.Sleep(time.Second)
		}
		assert.NoError(t, errors.New("unable to confirm inflight txs is empty"))
	}

	// happy path (send => simulate success => tx: nil => tx: processed => tx: confirmed => done)
	t.Run("happyPath", func(t *testing.T) {
		sig := getSig()
		tx := getTx(t, pubkey)
		var wg sync.WaitGroup
		wg.Add(3)

		sendCount := 0
		var countRW sync.RWMutex
		mc.On("SendTx", mock.Anything, tx).Run(func(mock.Arguments) {
			countRW.Lock()
			sendCount++
			countRW.Unlock()
		}).After(500*time.Millisecond).Return(sig, nil)
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

		// no transactions stored inflight txs list
		waitFor(empty)
		// transaction should be sent more than twice
		countRW.RLock()
		t.Logf("sendTx received %d calls", sendCount)
		assert.Greater(t, sendCount, 2)
		countRW.RUnlock()

		// panic if sendTx called after context cancelled
		mc.On("SendTx", mock.Anything, tx).Panic("SendTx should not be called anymore")

		// check prom metric
		prom.success++
		prom.assertEqual(t)
	})

	// fail on initial transmit (RPC immediate rejects)
	t.Run("fail_initialTx", func(t *testing.T) {
		tx := getTx(t, pubkey)
		var wg sync.WaitGroup
		wg.Add(1)

		// should only be called once (tx does not start retry, confirming, or simulation)
		mc.On("SendTx", mock.Anything, tx).Run(func(mock.Arguments) {
			wg.Done()
		}).Return(solana.Signature{}, errors.New("FAIL")).Once()

		// tx should be able to queue
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait() // wait to be picked up and processed

		// no transactions stored inflight txs list
		waitFor(empty)

		// check prom metric
		prom.error++
		prom.reject++
		prom.assertEqual(t)
	})

	// tx fails simulation (simulation error)
	t.Run("fail_simulation", func(t *testing.T) {
		tx := getTx(t, pubkey)
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
		waitFor(empty) // txs cleared quickly

		// check prom metric
		prom.error++
		prom.simOther++
		prom.assertEqual(t)
	})

	// tx fails simulation (rpc error, timeout should clean up b/c sig status will be nil)
	t.Run("fail_simulation_confirmNil", func(t *testing.T) {
		tx := getTx(t, pubkey)
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
		waitFor(empty) // txs cleared after timeout

		// check prom metric
		prom.error++
		prom.drop++
		prom.assertEqual(t)

		// panic if sendTx called after context cancelled
		mc.On("SendTx", mock.Anything, tx).Panic("SendTx should not be called anymore")
	})

	// tx fails simulation with an InstructionError (indicates reverted execution)
	// manager should cancel sending retry immediately + increment reverted prom metric
	t.Run("fail_simulation_instructionError", func(t *testing.T) {
		tx := getTx(t, pubkey)
		sig := getSig()
		var wg sync.WaitGroup
		wg.Add(1)

		// {"InstructionError":[0,{"Custom":6003}]}
		tempErr := map[string][]interface{}{
			"InstructionError": []interface{}{
				0, map[string]int{"Custom": 6003},
			},
		}
		mc.On("SendTx", mock.Anything, tx).Return(sig, nil)
		mc.On("SimulateTx", mock.Anything, tx, mock.Anything).Run(func(mock.Arguments) {
			wg.Done()
		}).Return(&rpc.SimulateTransactionResult{
			Err: tempErr,
		}, nil).Once()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Return([]*rpc.SignatureStatusesResult{nil}, nil).Maybe()

		// tx should be able to queue
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait()      // wait to be picked up and processed
		waitFor(empty) // txs cleared after timeout

		// check prom metric
		prom.error++
		prom.simRevert++
		prom.assertEqual(t)

		// panic if sendTx called after context cancelled
		mc.On("SendTx", mock.Anything, tx).Panic("SendTx should not be called anymore")
	})

	// tx fails simulation with BlockHashNotFound error
	// txm should continue to confirm tx (in this case it will succeed)
	t.Run("fail_simulation_blockhashNotFound", func(t *testing.T) {
		tx := getTx(t, pubkey)
		sig := getSig()
		var wg sync.WaitGroup
		wg.Add(2)

		mc.On("SendTx", mock.Anything, tx).Return(sig, nil)
		mc.On("SimulateTx", mock.Anything, tx, mock.Anything).Run(func(mock.Arguments) {
			wg.Done()
		}).Return(&rpc.SimulateTransactionResult{
			Err: "BlockhashNotFound",
		}, nil).Once()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Run(func(mock.Arguments) {
			wg.Done()
		}).Return([]*rpc.SignatureStatusesResult{&rpc.SignatureStatusesResult{
			ConfirmationStatus: rpc.ConfirmationStatusConfirmed,
		}}, nil).Once()
		// tx should be able to queue
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait()      // wait to be picked up and processed
		waitFor(empty) // txs cleared after timeout

		// check prom metric
		prom.success++
		prom.assertEqual(t)

		// panic if sendTx called after context cancelled
		mc.On("SendTx", mock.Anything, tx).Panic("SendTx should not be called anymore")
	})

	// tx fails simulation with AlreadyProcessed error
	// txm should continue to confirm tx (in this case it will revert)
	t.Run("fail_simulation_alreadyProcessed", func(t *testing.T) {
		tx := getTx(t, pubkey)
		sig := getSig()
		var wg sync.WaitGroup
		wg.Add(2)

		mc.On("SendTx", mock.Anything, tx).Return(sig, nil)
		mc.On("SimulateTx", mock.Anything, tx, mock.Anything).Run(func(mock.Arguments) {
			wg.Done()
		}).Return(&rpc.SimulateTransactionResult{
			Err: "AlreadyProcessed",
		}, nil).Once()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Run(func(mock.Arguments) {
			wg.Done()
		}).Return([]*rpc.SignatureStatusesResult{&rpc.SignatureStatusesResult{
			Err:                "ERROR",
			ConfirmationStatus: rpc.ConfirmationStatusConfirmed,
		}}, nil).Once()
		// tx should be able to queue
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait()      // wait to be picked up and processed
		waitFor(empty) // txs cleared after timeout

		// check prom metric
		prom.revert++
		prom.error++
		prom.assertEqual(t)

		// panic if sendTx called after context cancelled
		mc.On("SendTx", mock.Anything, tx).Panic("SendTx should not be called anymore")
	})

	// tx passes sim, never passes processed (timeout should cleanup)
	t.Run("fail_confirm_processed", func(t *testing.T) {
		tx := getTx(t, pubkey)
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
		waitFor(empty) // inflight txs cleared after timeout

		// check prom metric
		prom.error++
		prom.drop++
		prom.assertEqual(t)

		// panic if sendTx called after context cancelled
		mc.On("SendTx", mock.Anything, tx).Panic("SendTx should not be called anymore")
	})

	// tx passes sim, shows processed, moves to nil (timeout should cleanup)
	t.Run("fail_confirm_processedToNil", func(t *testing.T) {
		tx := getTx(t, pubkey)
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
		waitFor(empty) // inflight txs cleared after timeout

		// check prom metric
		prom.error++
		prom.drop++
		prom.assertEqual(t)

		// panic if sendTx called after context cancelled
		mc.On("SendTx", mock.Anything, tx).Panic("SendTx should not be called anymore")
	})

	// tx passes sim, errors on confirm
	t.Run("fail_confirm_revert", func(t *testing.T) {
		tx := getTx(t, pubkey)
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
		waitFor(empty) // inflight txs cleared after timeout

		// check prom metric
		prom.error++
		prom.revert++
		prom.assertEqual(t)

		// panic if sendTx called after context cancelled
		mc.On("SendTx", mock.Anything, tx).Panic("SendTx should not be called anymore")
	})
}

func TestTxm_Enqueue(t *testing.T) {
	// set up configs needed in txm
	lggr := logger.TestLogger(t)
	cfg := config.NewConfig(db.ChainCfg{}, lggr)
	mc := newReaderWriterMock(t)

	// mock solana keystore
	key, err := solkey.New()
	pubkey := key.PublicKey()

	require.NoError(t, err)
	mkey := keyMocks.NewSolana(t)
	mkey.On("Get", key.ID()).Return(key, nil)
	zerokey := solana.PublicKey{}
	mkey.On("Get", zerokey.String()).Return(solkey.Key{}, keystore.KeyNotFoundError{ID: zerokey.String(), KeyType: "Solana"})

	txm := NewTxm("enqueue_test", func() (client.ReaderWriter, error) {
		return mc, nil
	}, cfg, mkey, lggr)

	txs := []struct {
		name string
		tx   *solana.Transaction
		fail bool
	}{
		{"success", getTx(t, pubkey), false},
		{"invalid_key", getTx(t, zerokey), true},
		{"nil_pointer", nil, true},
		{"empty_tx", &solana.Transaction{}, true},
	}

	for _, run := range txs {
		t.Run(run.name, func(t *testing.T) {
			if !run.fail {
				assert.NoError(t, txm.Enqueue(run.name, run.tx))
				return
			}
			assert.Error(t, txm.Enqueue(run.name, run.tx))
		})
	}
}
