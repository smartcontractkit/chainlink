package txm

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	gasmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func TestBroadcaster_Lifecycle(t *testing.T) {
	lggr := logger.Test(t)
	chainID := big.NewInt(0)

	cfg, db := heavyweight.FullTestDBV2(t, nil)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	bcfg := BroadcasterConfig{
		FallbackPollInterval: evmcfg.Database().Listener().FallbackPollInterval(),
		MaxInFlight:          evmcfg.EVM().Transactions().MaxInFlight(),
		NonceAutoSync:        false,
	}

	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	client := evmtest.NewEthClientMockWithDefaultChain(t)
	ks := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

	cltest.MustInsertRandomKeyReturningState(t, ks)
	estimator := gasmocks.NewEvmFeeEstimator(t)
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*chainID, evmcfg.EVM().GasEstimator(), ks, estimator)
	ss := NewSequenceSyncer(lggr, txStore, client)
	b := NewBroadcaster(txBuilder, lggr, txStore, client, bcfg, ks, ss)

	// Can't close an unstarted instance
	err := b.Close()
	require.Error(t, err)
	ctx := testutils.Context(t)

	client.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(2), nil).Once()
	// Can start a new instance
	// TODO: fix Start failing on DB
	err = b.Start(ctx)
	require.NoError(t, err)

	// Can successfully close once
	err = b.Close()
	require.NoError(t, err)

	// Can't start more than once (Broadcaster uses services.StateMachine)
	err = b.Start(ctx)
	require.Error(t, err)
	// Can't close more than once (Broadcaster uses services.StateMachine)
	err = b.Close()
	require.Error(t, err)
}

func TestBroadcaster_ProcessUnstartedTxs_InProgress(t *testing.T) {
	lggr := logger.Test(t)
	chainID := big.NewInt(0)

	cfg, db := heavyweight.FullTestDBV2(t, nil)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	bcfg := BroadcasterConfig{
		FallbackPollInterval: evmcfg.Database().Listener().FallbackPollInterval(),
		MaxInFlight:          evmcfg.EVM().Transactions().MaxInFlight(),
		NonceAutoSync:        false,
	}
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	client := evmtest.NewEthClientMockWithDefaultChain(t)
	keyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

	_, addr1 := cltest.MustInsertRandomKey(t, keyStore)
	_, addr2 := cltest.MustInsertRandomKey(t, keyStore)

	estimator := gasmocks.NewEvmFeeEstimator(t)
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*chainID, evmcfg.EVM().GasEstimator(), keyStore, estimator)
	ss := NewSequenceSyncer(lggr, txStore, client)
	b := NewBroadcaster(txBuilder, lggr, txStore, client, bcfg, keyStore, ss)

	ctx := testutils.Context(t)
	encodedPayload := []byte{1, 2, 3}
	value := big.Int(assets.NewEthValue(142))
	gasLimit := uint32(242)
	t.Run("no txs at all", func(t *testing.T) {
		err := b.ProcessUnstartedTxs(ctx, addr1)
		require.NoError(t, err)
	})

	t.Run("returns error if handleAnyInProgressTx fails", func(t *testing.T) {
		temptCtx := context.Background()
		ctxWithCancel, cancel := context.WithCancel(temptCtx)
		cancel()
		err := b.ProcessUnstartedTxs(ctxWithCancel, addr1)
		require.Error(t, err)
	})

	t.Run("handles in progress tx successfully", func(t *testing.T) {
		nonce := evmtypes.Nonce(2)

		txInProgress := txmgr.Tx{
			FromAddress:    addr1,
			Sequence:       &nonce,
			ToAddress:      addr2,
			EncodedPayload: encodedPayload,
			Value:          value,
			FeeLimit:       gasLimit,
			Error:          null.String{},
			State:          TxInProgress,
		}

		err := txStore.InsertTx(&txInProgress)
		require.NoError(t, err)

		estimator.On("GetFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gas.EvmFee{Legacy: assets.GWei(32)}, uint32(500), nil).Once()
		client.On("SendTransaction", mock.Anything, mock.Anything).Return(nil).Once()
		ss.LoadNextSequenceMap(ctx, []common.Address{addr1})
		seq, err := ss.GetNextSequence(ctx, addr1)
		require.NoError(t, err)
		require.Equal(t, evmtypes.Nonce(3), seq)

		err = b.ProcessUnstartedTxs(ctx, addr1)
		require.NoError(t, err)
		// Nonce should have been incremented after successful broadcast
		seq, err = ss.GetNextSequence(ctx, addr1)
		require.NoError(t, err)
		require.Equal(t, evmtypes.Nonce(4), seq)
	})

	t.Run("returns error if handleInProgressTx fails", func(t *testing.T) {
		nonce := evmtypes.Nonce(3)

		txInProgress := txmgr.Tx{
			FromAddress:    addr1,
			Sequence:       &nonce,
			ToAddress:      addr2,
			EncodedPayload: encodedPayload,
			Value:          value,
			FeeLimit:       gasLimit,
			Error:          null.String{},
			State:          TxInProgress,
		}

		err := txStore.InsertTx(&txInProgress)
		require.NoError(t, err)

		estimator.On("GetFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gas.EvmFee{Legacy: assets.GWei(32)}, uint32(500), nil).Once()
		client.On("SendTransaction", mock.Anything, mock.Anything).Return(errors.New("RPC error")).Once()
		// Pending nonce wasn't updated, meaning tx wasn't successful.
		client.On("PendingNonceAt", mock.Anything, addr1).Return(uint64(3), nil).Once()

		err = b.ProcessUnstartedTxs(ctx, addr1)
		require.Error(t, err)
	})

}
func TestBroadcaster_ProcessUnstartedTxs_Unstarted(t *testing.T) {
	lggr := logger.Test(t)
	chainID := big.NewInt(0)

	cfg, db := heavyweight.FullTestDBV2(t, nil)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	client := evmtest.NewEthClientMockWithDefaultChain(t)
	keyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

	_, addr1 := cltest.MustInsertRandomKey(t, keyStore)
	_, addr2 := cltest.MustInsertRandomKey(t, keyStore)

	estimator := gasmocks.NewEvmFeeEstimator(t)
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*chainID, evmcfg.EVM().GasEstimator(), keyStore, estimator)
	ss := NewSequenceSyncer(lggr, txStore, client)

	ctx := testutils.Context(t)
	encodedPayload := []byte{1, 2, 3}
	value := big.Int(assets.NewEthValue(142))
	gasLimit := uint32(242)
	timeNow := time.Now()
	nonce := evmtypes.Nonce(0)
	t.Run("skips check if MaxInFlight is 0", func(t *testing.T) {
		bcfg := BroadcasterConfig{
			FallbackPollInterval: evmcfg.Database().Listener().FallbackPollInterval(),
			MaxInFlight:          0,
			NonceAutoSync:        false,
		}
		b := NewBroadcaster(txBuilder, lggr, txStore, client, bcfg, keyStore, ss)
		err := b.ProcessUnstartedTxs(ctx, utils.RandomAddress())
		require.NoError(t, err)
	})

	t.Run("picks up a new unstarted tx if in flight txs are less than threshold", func(t *testing.T) {
		bcfg := BroadcasterConfig{
			FallbackPollInterval: evmcfg.Database().Listener().FallbackPollInterval(),
			MaxInFlight:          3,
			NonceAutoSync:        false,
		}

		txUnconfirmed := txmgr.Tx{
			Sequence:           &nonce,
			FromAddress:        addr1,
			ToAddress:          addr2,
			EncodedPayload:     encodedPayload,
			Value:              value,
			FeeLimit:           gasLimit,
			BroadcastAt:        &timeNow,
			InitialBroadcastAt: &timeNow,
			Error:              null.String{},
			State:              TxUnconfirmed,
		}

		txUnstarted := txmgr.Tx{
			FromAddress:    addr1,
			ToAddress:      addr2,
			EncodedPayload: encodedPayload,
			Value:          value,
			FeeLimit:       gasLimit,
			State:          TxUnstarted,
		}

		require.NoError(t, txStore.InsertTx(&txUnconfirmed))
		require.NoError(t, txStore.InsertTx(&txUnstarted))

		ss.LoadNextSequenceMap(ctx, []common.Address{addr1})
		seq, err := ss.GetNextSequence(ctx, addr1)
		require.NoError(t, err)
		require.Equal(t, evmtypes.Nonce(nonce+1), seq)
		estimator.On("GetFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gas.EvmFee{Legacy: assets.GWei(32)}, uint32(500), nil).Once()
		client.On("SendTransaction", mock.Anything, mock.Anything).Return(nil).Once()

		b := NewBroadcaster(txBuilder, lggr, txStore, client, bcfg, keyStore, ss)
		require.NoError(t, b.ProcessUnstartedTxs(ctx, addr1))
		// Nonce should have been incremented after successful broadcast
		seq, err = ss.GetNextSequence(ctx, addr1)
		require.NoError(t, err)
		require.Equal(t, evmtypes.Nonce(2), seq)
	})

	t.Run("picks up a new unstarted tx and returns error if tx fails", func(t *testing.T) {
		bcfg := BroadcasterConfig{
			FallbackPollInterval: evmcfg.Database().Listener().FallbackPollInterval(),
			MaxInFlight:          evmcfg.EVM().Transactions().MaxInFlight(),
			NonceAutoSync:        false,
		}

		txUnstarted := txmgr.Tx{
			FromAddress:    addr1,
			ToAddress:      addr2,
			EncodedPayload: encodedPayload,
			Value:          value,
			FeeLimit:       gasLimit,
			State:          TxUnstarted,
		}

		db := pgtest.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db, cfg.Database())
		ss := NewSequenceSyncer(lggr, txStore, client)
		require.NoError(t, txStore.InsertTx(&txUnstarted))

		client.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil).Once()
		ss.LoadNextSequenceMap(ctx, []common.Address{addr1})
		seq, err := ss.GetNextSequence(ctx, addr1)
		require.NoError(t, err)
		require.Equal(t, evmtypes.Nonce(0), seq)

		estimator.On("GetFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gas.EvmFee{Legacy: assets.GWei(32)}, uint32(500), nil).Once()
		client.On("SendTransaction", mock.Anything, mock.Anything).Return(errors.New("RPC error")).Once()
		client.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil).Once()

		b := NewBroadcaster(txBuilder, lggr, txStore, client, bcfg, keyStore, ss)
		require.Error(t, b.ProcessUnstartedTxs(ctx, addr1))
	})

	t.Run("marks unstarted tx unconfirmed if tx fails but on-chain pending nonce increases", func(t *testing.T) {
		bcfg := BroadcasterConfig{
			FallbackPollInterval: evmcfg.Database().Listener().FallbackPollInterval(),
			MaxInFlight:          evmcfg.EVM().Transactions().MaxInFlight(),
			NonceAutoSync:        false,
		}

		txUnstarted := txmgr.Tx{
			FromAddress:    addr1,
			ToAddress:      addr2,
			EncodedPayload: encodedPayload,
			Value:          value,
			FeeLimit:       gasLimit,
			State:          TxUnstarted,
		}

		db := pgtest.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db, cfg.Database())
		ss := NewSequenceSyncer(lggr, txStore, client)
		require.NoError(t, txStore.InsertTx(&txUnstarted))

		client.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil).Once()
		ss.LoadNextSequenceMap(ctx, []common.Address{addr1})
		seq, err := ss.GetNextSequence(ctx, addr1)
		require.NoError(t, err)
		require.Equal(t, evmtypes.Nonce(0), seq)

		estimator.On("GetFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gas.EvmFee{Legacy: assets.GWei(32)}, uint32(500), nil).Once()
		client.On("SendTransaction", mock.Anything, mock.Anything).Return(errors.New("RPC error")).Once()
		client.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil).Once()

		b := NewBroadcaster(txBuilder, lggr, txStore, client, bcfg, keyStore, ss)
		require.NoError(t, b.ProcessUnstartedTxs(ctx, addr1))
	})
}
