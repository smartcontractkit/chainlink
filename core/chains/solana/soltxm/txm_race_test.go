package soltxm

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"testing"
	"time"

	solanaGo "github.com/gagliardetto/solana-go"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/solana/fees"
	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils"

	solanaClient "github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/solkey"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func NewTestTx() (tx solanaGo.Transaction) {
	tx.Message.AccountKeys = append(tx.Message.AccountKeys, solanaGo.PublicKey{})
	return tx
}

// Test race condition for saving + reading signatures when bumping fees
// A slow RPC can cause the tx (before bump) to be processed after the bumped tx
// The bumped tx will cause the storage check to move on to the next tx signature even with a inflight "non-bumped" tx
func TestTxm_SendWithRetry_Race(t *testing.T) {
	// test config
	txRetryDuration := 2 * time.Second

	// mocks init
	lggr, observer := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	cfg := config.NewConfig(db.ChainCfg{
		ComputeUnitPriceMax: null.IntFrom(10),
		ComputeUnitPriceMin: null.IntFrom(0),
		FeeBumpPeriod:       utils.MustNewDuration(txRetryDuration / 6),
	}, lggr)
	ks := ksmocks.NewSolana(t)
	fee, err := fees.NewFixedPriceEstimator(cfg)
	require.NoError(t, err)

	// keystore mock
	key, err := solkey.New()
	require.NoError(t, err)
	ks.On("Get", mock.Anything).Return(key, nil)

	// assemble minimal tx for testing retry
	tx := NewTestTx()

	testRunner := func(t *testing.T, client solanaClient.ReaderWriter) {
		getClient := func() (solanaClient.ReaderWriter, error) {
			return client, nil
		}

		// build minimal txm
		txm := NewTxm("retry_race", getClient, cfg, ks, lggr)
		txm.fee = fee

		_, _, _, err := txm.sendWithRetry(
			context.Background(),
			tx,
			txRetryDuration,
		)
		require.NoError(t, err)

		time.Sleep(txRetryDuration / 4 * 5)                                     // wait 1.25x longer of tx life to capture all logs
		assert.Equal(t, observer.FilterLevelExact(zapcore.ErrorLevel).Len(), 0) // assert no error logs
		lastLog := observer.All()[len(observer.All())-1]
		assert.Contains(t, lastLog.Message, "stopped tx retry") // assert that all retry goroutines exit successfully
	}

	t.Run("delay in rebroadcasting tx", func(t *testing.T) {
		client := newReaderWriterMock(t)
		// client mock
		txs := map[string]solanaGo.Signature{}
		var lock sync.RWMutex
		client.On("SendTx", mock.Anything, mock.Anything).Return(
			// build new sig if tx is different
			func(_ context.Context, tx *solanaGo.Transaction) solanaGo.Signature {
				strTx := tx.String()

				// if exists, slow down client response to trigger race
				lock.RLock()
				val, exists := txs[strTx]
				lock.RUnlock()
				if exists {
					time.Sleep(txRetryDuration / 3)
					return val
				}

				lock.Lock()
				defer lock.Unlock()
				// recheck existence
				val, exists = txs[strTx]
				if exists {
					return val
				}
				sig := make([]byte, 16)
				rand.Read(sig)
				txs[strTx] = solanaGo.SignatureFromBytes(sig)

				return txs[strTx]
			},
			nil,
		)
		testRunner(t, client)
	})

	t.Run("delay in broadcasting new tx", func(t *testing.T) {
		client := newReaderWriterMock(t)
		// client mock
		txs := map[string]solanaGo.Signature{}
		var lock sync.RWMutex
		client.On("SendTx", mock.Anything, mock.Anything).Return(
			// build new sig if tx is different
			func(_ context.Context, tx *solanaGo.Transaction) solanaGo.Signature {
				strTx := tx.String()

				lock.Lock()
				// check existence
				val, exists := txs[strTx]
				if exists {
					lock.Unlock()
					return val
				}
				sig := make([]byte, 16)
				rand.Read(sig)
				txs[strTx] = solanaGo.SignatureFromBytes(sig)
				lock.Unlock()

				// don't lock on delay
				// delay every new bumping tx
				time.Sleep(txRetryDuration / 3)

				lock.RLock()
				defer lock.RUnlock()
				return txs[strTx]
			},
			nil,
		)
		testRunner(t, client)
	})

	t.Run("overlapping bumping tx", func(t *testing.T) {
		client := newReaderWriterMock(t)
		// client mock
		txs := map[string]solanaGo.Signature{}
		var lock sync.RWMutex
		client.On("SendTx", mock.Anything, mock.Anything).Return(
			// build new sig if tx is different
			func(_ context.Context, tx *solanaGo.Transaction) solanaGo.Signature {
				strTx := tx.String()

				lock.Lock()
				// recheck existence
				val, exists := txs[strTx]
				if exists {
					lock.Unlock()
					return val
				}
				sig := make([]byte, 16)
				rand.Read(sig)
				txs[strTx] = solanaGo.SignatureFromBytes(sig)

				triggerDelay := len(txs) == 2
				lock.Unlock()

				// don't lock on delay
				// only delay on the first bump tx
				// ------------------------------
				// init tx - no delay
				// rebroadcast - no delay (tx + sig already exists, does not reach this point)
				// first bump tx - DELAY
				// rebroadcast bump tx - no delay (tx + sig already exists, does not reach this point)
				// second bump tx - no delay
				// etc
				if triggerDelay {
					time.Sleep(txRetryDuration * 2 / 3)
				}

				lock.RLock()
				defer lock.RUnlock()
				return txs[strTx]
			},
			nil,
		)
		testRunner(t, client)
	})

	t.Run("bumping tx errors and ctx cleans up waitgroup blocks", func(t *testing.T) {
		client := newReaderWriterMock(t)
		// client mock - first tx is always successful
		tx0 := NewTestTx()
		require.NoError(t, fees.SetComputeUnitPrice(&tx0, 0))
		client.On("SendTx", mock.Anything, mock.MatchedBy(func(tx *solanaGo.Transaction) bool {
			return bytes.Equal(tx.Message.Instructions[0].Data, tx0.Message.Instructions[0].Data)
		})).Return(solanaGo.Signature{1}, nil)

		// init bump tx fails, rebroadcast is successful
		tx1 := NewTestTx()
		require.NoError(t, fees.SetComputeUnitPrice(&tx1, 1))
		client.On("SendTx", mock.Anything, mock.MatchedBy(func(tx *solanaGo.Transaction) bool {
			return bytes.Equal(tx.Message.Instructions[0].Data, tx1.Message.Instructions[0].Data)
		})).Return(solanaGo.Signature{}, fmt.Errorf("BUMP FAILED")).Once()
		client.On("SendTx", mock.Anything, mock.MatchedBy(func(tx *solanaGo.Transaction) bool {
			return bytes.Equal(tx.Message.Instructions[0].Data, tx1.Message.Instructions[0].Data)
		})).Return(solanaGo.Signature{2}, nil)

		// init bump tx success, rebroadcast fails
		tx2 := NewTestTx()
		require.NoError(t, fees.SetComputeUnitPrice(&tx2, 2))
		client.On("SendTx", mock.Anything, mock.MatchedBy(func(tx *solanaGo.Transaction) bool {
			return bytes.Equal(tx.Message.Instructions[0].Data, tx2.Message.Instructions[0].Data)
		})).Return(solanaGo.Signature{3}, nil).Once()
		client.On("SendTx", mock.Anything, mock.MatchedBy(func(tx *solanaGo.Transaction) bool {
			return bytes.Equal(tx.Message.Instructions[0].Data, tx2.Message.Instructions[0].Data)
		})).Return(solanaGo.Signature{}, fmt.Errorf("REBROADCAST FAILED"))

		// always successful
		tx3 := NewTestTx()
		require.NoError(t, fees.SetComputeUnitPrice(&tx3, 4))
		client.On("SendTx", mock.Anything, mock.MatchedBy(func(tx *solanaGo.Transaction) bool {
			return bytes.Equal(tx.Message.Instructions[0].Data, tx3.Message.Instructions[0].Data)
		})).Return(solanaGo.Signature{4}, nil)

		testRunner(t, client)
	})
}
