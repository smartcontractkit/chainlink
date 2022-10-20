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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
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

var (
	// {"InstructionError":[0,{"Custom":6003}]}
	instructionErr = map[string][]interface{}{
		"InstructionError": {
			0, map[string]int{"Custom": 6003},
		},
	}

	confirmedStatus = &rpc.SignatureStatusesResult{
		ConfirmationStatus: rpc.ConfirmationStatusConfirmed,
	}
)

type soltxmProm struct {
	id                                      string
	success, error, revert, reject, invalid float64
}

func (p soltxmProm) assertEqual(t *testing.T) {
	assert.Equal(t, p.success, testutil.ToFloat64(promSolTxmSuccessTxs.WithLabelValues(p.id)), "mismatch: success")
	assert.Equal(t, p.error, testutil.ToFloat64(promSolTxmErrorTxs.WithLabelValues(p.id)), "mismatch: error")
	assert.Equal(t, p.revert, testutil.ToFloat64(promSolTxmRevertTxs.WithLabelValues(p.id)), "mismatch: revert")
	assert.Equal(t, p.reject, testutil.ToFloat64(promSolTxmRejectTxs.WithLabelValues(p.id)), "mismatch: reject")
	assert.Equal(t, p.invalid, testutil.ToFloat64(promSolTxmInvalidBlockhash.WithLabelValues(p.id)), "mismatch: invalid")
}

func (p soltxmProm) getInflight() float64 {
	return testutil.ToFloat64(promSolTxmPendingTxs.WithLabelValues(p.id))
}

func (p *soltxmProm) Reset() {
	promSolTxmSuccessTxs.Reset()
	promSolTxmPendingTxs.Reset()
	promSolTxmErrorTxs.Reset()
	promSolTxmRevertTxs.Reset()
	promSolTxmRejectTxs.Reset()
	promSolTxmInvalidBlockhash.Reset()
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
	require.NoError(t, err)
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
	cfg := config.NewConfig(db.ChainCfg{ // reduce time for faster test execution 17-18s => 7s
		ConfirmPollPeriod:       utils.MustNewDuration(5 * time.Millisecond),
		TxConfirmTimeout:        utils.MustNewDuration(10 * time.Millisecond),
		DefaultComputeUnitPrice: null.IntFrom(100),
	}, lggr)
	mTx := mock.AnythingOfType("*solana.Transaction")
	mSig := mock.AnythingOfType("[]solana.Signature")

	// mock solana keystore
	key, err := solkey.New()
	pubkey := key.PublicKey()

	require.NoError(t, err)
	mkey := keyMocks.NewSolana(t)
	mkey.On("Get", key.ID()).Return(key, nil)

	// tracking prom metrics
	prom := soltxmProm{id: id}
	prom.Reset()          // clear previous state
	t.Cleanup(prom.Reset) // clean up existing state

	// create new to limit the scope of each txm test
	initTxm := func() (*Txm, *mocks.ReaderWriter, func() bool) {
		mc := newReaderWriterMock(t)

		txm := NewTxm(id, func() (client.ReaderWriter, error) {
			return mc, nil
		}, cfg, mkey, lggr)
		require.NoError(t, txm.Start(testutils.Context(t)))

		// provide function to check if cached transaction is cleared
		empty := func() bool {
			idCount, sigCount := txm.InflightTxs()
			assert.Equal(t, float64(sigCount), prom.getInflight()) // validate prom metric and txs length
			t.Logf("tx count: IDs - %d, sigs - %d", idCount, sigCount)
			return idCount == 0 && sigCount == 0
		}

		return txm, mc, empty
	}

	// create random signature
	getSig := func() solana.Signature {
		sig := make([]byte, 64)
		rand.Read(sig)
		return solana.SignatureFromBytes(sig)
	}

	// poll for max 30s before quiting (all txs should complete)
	waitFor := func(f func() bool) {
		for i := 0; i < 30; i++ {
			if f() {
				return
			}
			time.Sleep(time.Second)
		}
		assert.NoError(t, errors.New("unable to confirm inflight txs is empty"))
	}

	// happy path (send => simulate success => tx: nil => tx: processed => tx: confirmed => done)
	t.Run("success_happyPath", func(t *testing.T) {
		txm, mc, empty := initTxm()
		sig := getSig()
		tx := getTx(t, pubkey)
		var wg sync.WaitGroup
		wg.Add(3)

		mc.On("SendTx", mock.Anything, mTx).Return(sig, nil).Once()
		mc.On("IsBlockhashValid", mock.Anything, solana.Hash{}).Return(true, nil).Once()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Run(func(mock.Arguments) {
			wg.Done()
		}).Return([]*rpc.SignatureStatusesResult{nil}, nil).Once()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Run(func(mock.Arguments) {
			wg.Done()
		}).Return([]*rpc.SignatureStatusesResult{{
			ConfirmationStatus: rpc.ConfirmationStatusProcessed,
		}}, nil).Once()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Run(func(mock.Arguments) {
			wg.Done()
		}).Return([]*rpc.SignatureStatusesResult{confirmedStatus}, nil).Once()

		// send tx
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait()

		// no transactions stored inflight txs list
		waitFor(empty)

		// check prom metric
		prom.success++
		prom.assertEqual(t)
	})

	// txm should handle RPC rejection properly
	// RPC rejection should not trigger a fee bump
	// fail on initial transmit (RPC immediate rejects)
	t.Run("success_rpcFail_retry", func(t *testing.T) {
		txm, mc, empty := initTxm()
		tx := getTx(t, pubkey)
		sig := getSig()
		var wg sync.WaitGroup
		wg.Add(3)
		var fee uint64

		// immediate fail initial send (should retry without fee bump)
		mc.On("SendTx", mock.Anything, mTx).Run(func(args mock.Arguments) {
			rawTx := args.Get(1).(*solana.Transaction)
			fee = XXXGetFeePrice(t, rawTx)
			wg.Done()
		}).Return(solana.Signature{}, errors.New("FAIL")).Once()
		mc.On("SendTx", mock.Anything, mTx).Run(func(args mock.Arguments) {
			rawTx := args.Get(1).(*solana.Transaction)
			val := XXXGetFeePrice(t, rawTx)
			assert.Equal(t, fee, val)
			wg.Done()
		}).Return(sig, nil).Once()
		mc.On("IsBlockhashValid", mock.Anything, solana.Hash{}).Return(true, nil).Twice()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Run(func(mock.Arguments) {
			wg.Done()
		}).Return([]*rpc.SignatureStatusesResult{confirmedStatus}, nil).Once()

		// tx should be able to queue
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait() // wait to be picked up and processed

		// no transactions stored inflight txs list
		waitFor(empty)

		// check prom metric
		prom.success++
		prom.reject++
		prom.error++
		prom.assertEqual(t)
	})

	// initial tx is sent, tx was dropped
	// second tx sent with higher fee, tx was dropped
	// third tx sent with higher fee, included
	// validate that fees are increasing properly
	// txm should complete all pending signatures for the 1 base tx
	t.Run("success_initialDropped_bumpFeeConfirmed", func(t *testing.T) {
		txm, mc, empty := initTxm()
		tx := getTx(t, pubkey)
		sigs := []solana.Signature{getSig(), getSig(), getSig()}
		fees := []uint64{100, 200, 400} // base fee set to 100, follow progression
		var wg sync.WaitGroup
		wg.Add(3)

		// 3 retried txs
		i := 0
		mc.On("SendTx", mock.Anything, mTx).Run(func(args mock.Arguments) {
			wg.Done()
			rawTx := args.Get(1).(*solana.Transaction)
			val := XXXGetFeePrice(t, rawTx)
			assert.Equal(t, fees[i], val)
		}).Return(
			func(_ context.Context, _ *solana.Transaction) solana.Signature {
				defer func() { i++ }()
				return sigs[i]
			},
			func(_ context.Context, _ *solana.Transaction) error {
				return nil
			},
		).Times(3)
		mc.On("IsBlockhashValid", mock.Anything, solana.Hash{}).Return(true, nil).Times(3)
		mc.On("SignatureStatuses", mock.Anything, mSig).Return(
			func(_ context.Context, sigs []solana.Signature) (out []*rpc.SignatureStatusesResult) {
				for i := 0; i < len(sigs); i++ {
					out = append(out, nil)
				}

				if len(sigs) >= 3 {
					out[2] = confirmedStatus
				}

				return
			},
			func(_ context.Context, _ []solana.Signature) error {
				return nil
			},
		)

		// tx should be able to queue
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait()      // wait to be picked up and processed
		waitFor(empty) // txs cleared after timeout

		// check prom metric
		prom.success++
		prom.assertEqual(t)
	})

	// initial tx is not confirmed in time, is eventually confirmed
	// second tx sent, also executed on chain but reverted
	t.Run("success_initialSuccess_retryFail", func(t *testing.T) {
		txm, mc, empty := initTxm()
		tx := getTx(t, pubkey)
		sig := getSig()
		sigRetry := getSig()
		var wg sync.WaitGroup
		wg.Add(2)

		mc.On("SendTx", mock.Anything, mTx).Return(sig, nil).Run(
			func(_ mock.Arguments) {
				wg.Done()
			},
		).Once()
		mc.On("SendTx", mock.Anything, mTx).Return(sigRetry, nil).Run(
			func(_ mock.Arguments) {
				wg.Done()
			},
		).Once()
		mc.On("IsBlockhashValid", mock.Anything, solana.Hash{}).Return(true, nil).Twice()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Return(
			[]*rpc.SignatureStatusesResult{nil}, nil)
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig, sigRetry}).Return(
			[]*rpc.SignatureStatusesResult{confirmedStatus, {Err: instructionErr}}, nil).Maybe()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sigRetry, sig}).Return(
			[]*rpc.SignatureStatusesResult{{Err: instructionErr}, confirmedStatus}, nil).Maybe()

		// tx should be able to queue
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait()      // wait to be picked up and processed
		waitFor(empty) // inflight txs cleared after timeout

		// check prom metric
		prom.success++
		prom.assertEqual(t)
	})

	// tx fails with an InstructionError (indicates reverted execution)
	t.Run("fail_onchainRevert", func(t *testing.T) {
		txm, mc, empty := initTxm()
		tx := getTx(t, pubkey)
		sig := getSig()
		var wg sync.WaitGroup
		wg.Add(2)

		mc.On("SendTx", mock.Anything, mTx).Run(func(_ mock.Arguments) {
			wg.Done()
		}).Return(sig, nil).Once()
		mc.On("IsBlockhashValid", mock.Anything, solana.Hash{}).Return(true, nil).Once()
		mc.On("SignatureStatuses", mock.Anything, mSig).Run(func(_ mock.Arguments) {
			wg.Done()
		}).Return(
			[]*rpc.SignatureStatusesResult{
				{
					Err: instructionErr,
				},
			}, nil).Once()

		// tx should be able to queue
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait()      // wait to be picked up and processed
		waitFor(empty) // txs cleared after timeout

		// check prom metric
		prom.error++
		prom.revert++
		prom.assertEqual(t)
	})

	// tx shows processed, moves to nil, txm should rebroadcast
	t.Run("success_processedToNil_rebroadcast", func(t *testing.T) {
		txm, mc, empty := initTxm()
		tx := getTx(t, pubkey)
		sig := getSig()
		sigRetry := getSig()
		var wg sync.WaitGroup
		wg.Add(2)

		mc.On("SendTx", mock.Anything, mTx).Return(sig, nil).Run(
			func(_ mock.Arguments) {
				wg.Done()
			},
		).Once()
		mc.On("SendTx", mock.Anything, mTx).Return(sigRetry, nil).Run(
			func(_ mock.Arguments) {
				wg.Done()
			},
		).Once()
		mc.On("IsBlockhashValid", mock.Anything, solana.Hash{}).Return(true, nil).Twice()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Return([]*rpc.SignatureStatusesResult{{
			ConfirmationStatus: rpc.ConfirmationStatusProcessed,
		}}, nil).Twice()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig}).Return(
			[]*rpc.SignatureStatusesResult{nil}, nil) // drop tx (reorg)
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sig, sigRetry}).Return(
			[]*rpc.SignatureStatusesResult{nil, confirmedStatus}, nil).Maybe()
		mc.On("SignatureStatuses", mock.Anything, []solana.Signature{sigRetry, sig}).Return(
			[]*rpc.SignatureStatusesResult{confirmedStatus, nil}, nil).Maybe()

		// tx should be able to queue
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait()      // wait to be picked up and processed
		waitFor(empty) // inflight txs cleared after timeout

		// check prom metric
		prom.success++
		prom.assertEqual(t)
	})

	// tx fails with an invalid blockhash
	t.Run("fail_invalidBlockhash", func(t *testing.T) {
		txm, mc, empty := initTxm()
		tx := getTx(t, pubkey)
		var wg sync.WaitGroup
		wg.Add(1)

		mc.On("IsBlockhashValid", mock.Anything, solana.Hash{}).Run(func(_ mock.Arguments) {
			wg.Done()
		}).Return(false, nil).Once()

		// tx should be able to queue
		assert.NoError(t, txm.Enqueue(t.Name(), tx))
		wg.Wait()      // wait to be picked up and processed
		waitFor(empty) // txs cleared after timeout

		// check prom metric
		prom.error++
		prom.invalid++
		prom.assertEqual(t)
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

// test mismatched signatures + responses in confirmer
func TestTxm_Confirmer(t *testing.T) {
	lggr, logs := logger.TestLoggerObserved(t, zapcore.ErrorLevel)
	cfg := config.NewConfig(db.ChainCfg{}, lggr)
	mc := newReaderWriterMock(t)
	var wg sync.WaitGroup
	wg.Add(4)

	txm := NewTxm("enqueue_test", func() (client.ReaderWriter, error) {
		return mc, nil
	}, cfg, keyMocks.NewSolana(t), lggr)

	id := txm.txs.New(PendingTx{})
	assert.NoError(t, txm.txs.Add(id, XXXNewSignature(t), 0))

	i := 0
	mc.On("SignatureStatuses", mock.Anything, mock.AnythingOfType("[]solana.Signature")).Run(
		func(args mock.Arguments) { wg.Done() },
	).Return(
		func(_ context.Context, _ []solana.Signature) (out []*rpc.SignatureStatusesResult) {
			for j := 0; j < 2-i; j++ {
				out = append(out, nil)
			}

			i++
			return out
		}, nil).Times(3) // return 2, 1, 0 results
	mc.On("SignatureStatuses", mock.Anything, mock.AnythingOfType("[]solana.Signature")).Run(
		func(args mock.Arguments) { wg.Done() },
	).Return([]*rpc.SignatureStatusesResult{confirmedStatus}, nil).Once()

	go txm.confirm(context.Background()) // start only confirmer
	wg.Wait()
	assert.Equal(t, 1, len(logs.All())) // extra results and equal results do not trigger error
	assert.Equal(t, "mismatch requested signatures and responses length: 1 > 0", logs.All()[0].Entry.Message)
}
