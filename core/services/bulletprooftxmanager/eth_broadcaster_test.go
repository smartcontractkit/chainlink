package bulletprooftxmanager_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	gasmocks "github.com/smartcontractkit/chainlink/core/services/gas/mocks"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestEthBroadcaster_ProcessUnstartedEthTxs_Success(t *testing.T) {
	db := pgtest.NewGormDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	cfg := configtest.NewTestGeneralConfig(t)
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState})

	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	timeNow := time.Now()

	encodedPayload := []byte{1, 2, 3}
	value := assets.NewEthValue(142)
	gasLimit := uint64(242)

	t.Run("no eth_txes at all", func(t *testing.T) {
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))
	})

	t.Run("eth_txes exist for a different from address", func(t *testing.T) {
		_, otherAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

		etx := bulletprooftxmanager.EthTx{
			FromAddress:    otherAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          bulletprooftxmanager.EthTxUnstarted,
		}
		require.NoError(t, db.Save(&etx).Error)

		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))
	})

	t.Run("existing eth_txes with broadcast_at or error", func(t *testing.T) {
		nonce := int64(342)
		errStr := "some error"

		etxUnconfirmed := bulletprooftxmanager.EthTx{
			Nonce:          &nonce,
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			BroadcastAt:    &timeNow,
			Error:          null.String{},
			State:          bulletprooftxmanager.EthTxUnconfirmed,
		}
		etxWithError := bulletprooftxmanager.EthTx{
			Nonce:          nil,
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			Error:          null.StringFrom(errStr),
			State:          bulletprooftxmanager.EthTxFatalError,
		}

		require.NoError(t, db.Save(&etxUnconfirmed).Error)
		require.NoError(t, db.Save(&etxWithError).Error)

		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))
	})

	t.Run("sends 3 EthTxs in order with higher value last, and lower values starting from the earliest", func(t *testing.T) {
		// Higher value
		expensiveEthTx := bulletprooftxmanager.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 42, 0},
			Value:          assets.NewEthValue(242),
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(0, 0),
			State:          bulletprooftxmanager.EthTxUnstarted,
		}
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(2) && tx.Value().Cmp(big.NewInt(242)) == 0
		})).Return(nil).Once()

		// Earlier
		tr := int32(99)
		b, err := json.Marshal(bulletprooftxmanager.EthTxMeta{JobID: tr})
		require.NoError(t, err)
		meta := datatypes.JSON(b)
		earlierEthTx := bulletprooftxmanager.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 42, 0},
			Value:          value,
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(0, 1),
			State:          bulletprooftxmanager.EthTxUnstarted,
			Meta:           &meta,
		}
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			if tx.Nonce() != uint64(0) {
				return false
			}
			require.Equal(t, evmcfg.ChainID(), tx.ChainId())
			require.Equal(t, gasLimit, tx.Gas())
			require.Equal(t, evmcfg.EvmGasPriceDefault(), tx.GasPrice())
			require.Equal(t, toAddress, *tx.To())
			require.Equal(t, value.ToInt().String(), tx.Value().String())
			require.Equal(t, earlierEthTx.EncodedPayload, tx.Data())
			return true
		})).Return(nil).Once()

		// Later
		laterEthTx := bulletprooftxmanager.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 42, 1},
			Value:          value,
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(1, 0),
			State:          bulletprooftxmanager.EthTxUnstarted,
		}
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			if tx.Nonce() != uint64(1) {
				return false
			}
			require.Equal(t, evmcfg.ChainID(), tx.ChainId())
			require.Equal(t, gasLimit, tx.Gas())
			require.Equal(t, evmcfg.EvmGasPriceDefault(), tx.GasPrice())
			require.Equal(t, toAddress, *tx.To())
			require.Equal(t, value.ToInt().String(), tx.Value().String())
			require.Equal(t, laterEthTx.EncodedPayload, tx.Data())
			return true
		})).Return(nil).Once()

		// Insertion order deliberately reversed to test ordering
		require.NoError(t, db.Save(&expensiveEthTx).Error)
		require.NoError(t, db.Save(&laterEthTx).Error)
		require.NoError(t, db.Save(&earlierEthTx).Error)

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check earlierEthTx and it's attempt
		// This was the earlier one sent so it has the lower nonce
		earlierTransaction, err := cltest.FindEthTxWithAttempts(db, earlierEthTx.ID)
		require.NoError(t, err)
		assert.False(t, earlierTransaction.Error.Valid)
		require.NotNil(t, earlierTransaction.FromAddress)
		assert.Equal(t, fromAddress, earlierTransaction.FromAddress)
		require.NotNil(t, earlierTransaction.Nonce)
		assert.Equal(t, int64(0), *earlierTransaction.Nonce)
		assert.NotNil(t, earlierTransaction.BroadcastAt)
		assert.Len(t, earlierTransaction.EthTxAttempts, 1)
		var m bulletprooftxmanager.EthTxMeta
		err = json.Unmarshal(*earlierEthTx.Meta, &m)
		require.NoError(t, err)
		assert.Equal(t, tr, m.JobID)

		attempt := earlierTransaction.EthTxAttempts[0]

		assert.Equal(t, earlierTransaction.ID, attempt.EthTxID)
		assert.NotNil(t, attempt.GasPrice)
		assert.Nil(t, attempt.GasTipCap)
		assert.Nil(t, attempt.GasFeeCap)
		assert.Equal(t, evmcfg.EvmGasPriceDefault().String(), attempt.GasPrice.String())

		_, err = attempt.GetSignedTx()
		require.NoError(t, err)
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt.State)
		require.Len(t, attempt.EthReceipts, 0)

		// Check laterEthTx and it's attempt
		// This was the later one sent so it has the higher nonce
		laterTransaction, err := cltest.FindEthTxWithAttempts(db, laterEthTx.ID)
		require.NoError(t, err)
		assert.False(t, earlierTransaction.Error.Valid)
		require.NotNil(t, laterTransaction.FromAddress)
		assert.Equal(t, fromAddress, laterTransaction.FromAddress)
		require.NotNil(t, laterTransaction.Nonce)
		assert.Equal(t, int64(1), *laterTransaction.Nonce)
		assert.NotNil(t, laterTransaction.BroadcastAt)
		assert.Len(t, laterTransaction.EthTxAttempts, 1)

		attempt = laterTransaction.EthTxAttempts[0]

		assert.Equal(t, laterTransaction.ID, attempt.EthTxID)
		assert.Equal(t, evmcfg.EvmGasPriceDefault().String(), attempt.GasPrice.String())

		_, err = attempt.GetSignedTx()
		require.NoError(t, err)
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt.State)
		require.Len(t, attempt.EthReceipts, 0)

		ethClient.AssertExpectations(t)
	})

	t.Run("sends transactions with type 0x2 in EIP-1559 mode", func(t *testing.T) {
		cfg.Overrides.GlobalEvmEIP1559DynamicFees = null.BoolFrom(true)
		rnd := int64(1000000000 + rand.Intn(5000))
		cfg.Overrides.GlobalEvmGasTipCapDefault = big.NewInt(rnd)
		cfg.Overrides.GlobalEvmMaxGasPriceWei = big.NewInt(rnd + 1)

		eipTxWithoutAl := bulletprooftxmanager.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 0, 0},
			Value:          assets.NewEthValue(142),
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(0, 0),
			State:          bulletprooftxmanager.EthTxUnstarted,
		}
		eipTxWithAl := bulletprooftxmanager.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 42, 0},
			Value:          assets.NewEthValue(242),
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(0, 1),
			State:          bulletprooftxmanager.EthTxUnstarted,
			AccessList:     bulletprooftxmanager.NullableEIP2930AccessListFrom(gethTypes.AccessList{gethTypes.AccessTuple{Address: cltest.NewAddress(), StorageKeys: []gethCommon.Hash{utils.NewHash()}}}),
		}
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(3) && tx.Value().Cmp(big.NewInt(142)) == 0
		})).Return(nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(4) && tx.Value().Cmp(big.NewInt(242)) == 0
		})).Return(nil).Once()

		require.NoError(t, db.Save(&eipTxWithAl).Error)
		require.NoError(t, db.Save(&eipTxWithoutAl).Error)

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check eipTxWithAl and it's attempt
		// This was the earlier one sent so it has the lower nonce
		eipTxWithAl, err := cltest.FindEthTxWithAttempts(db, eipTxWithAl.ID)
		require.NoError(t, err)
		assert.False(t, eipTxWithAl.Error.Valid)
		require.NotNil(t, eipTxWithAl.FromAddress)
		assert.Equal(t, fromAddress, eipTxWithAl.FromAddress)
		require.NotNil(t, eipTxWithAl.Nonce)
		assert.Equal(t, int64(4), *eipTxWithAl.Nonce)
		assert.NotNil(t, eipTxWithAl.BroadcastAt)
		assert.True(t, eipTxWithAl.AccessList.Valid)
		assert.Len(t, eipTxWithAl.AccessList.AccessList, 1)
		assert.Len(t, eipTxWithAl.EthTxAttempts, 1)

		attempt := eipTxWithAl.EthTxAttempts[0]

		assert.Equal(t, eipTxWithAl.ID, attempt.EthTxID)
		assert.Nil(t, attempt.GasPrice)
		assert.Equal(t, rnd, attempt.GasTipCap.ToInt().Int64())
		assert.Equal(t, rnd+1, attempt.GasFeeCap.ToInt().Int64())

		_, err = attempt.GetSignedTx()
		require.NoError(t, err)
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt.State)
		require.Len(t, attempt.EthReceipts, 0)
	})

	ethClient.AssertExpectations(t)

	t.Run("transaction simulation", func(t *testing.T) {
		t.Run("when simulation succeeds, sends tx as normal", func(t *testing.T) {
			ethTx := bulletprooftxmanager.EthTx{
				FromAddress:    fromAddress,
				ToAddress:      toAddress,
				EncodedPayload: []byte{42, 0, 0},
				Value:          assets.NewEthValue(442),
				GasLimit:       gasLimit,
				CreatedAt:      time.Unix(0, 0),
				State:          bulletprooftxmanager.EthTxUnstarted,
				Simulate:       true,
			}
			ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				return tx.Nonce() == uint64(5) && tx.Value().Cmp(big.NewInt(442)) == 0
			})).Return(nil).Once()
			ethClient.On("CallContext", mock.Anything, mock.AnythingOfType("*hexutil.Bytes"), "eth_call", mock.MatchedBy(func(callarg map[string]interface{}) bool {
				if fmt.Sprintf("%s", callarg["value"]) == "0x1ba" { // 442
					assert.Equal(t, ethTx.FromAddress, callarg["from"])
					assert.Equal(t, &ethTx.ToAddress, callarg["to"])
					assert.Equal(t, hexutil.Uint64(ethTx.GasLimit), callarg["gas"])
					assert.Nil(t, callarg["gasPrice"])
					assert.Nil(t, callarg["maxFeePerGas"])
					assert.Nil(t, callarg["maxPriorityFeePerGas"])
					assert.Equal(t, (*hexutil.Big)(&ethTx.Value), callarg["value"])
					assert.Equal(t, hexutil.Bytes(ethTx.EncodedPayload), callarg["data"])
					return true
				}
				return false
			}), "latest").Return(nil).Once()

			require.NoError(t, db.Save(&ethTx).Error)

			require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

			// Check ethtx was sent
			ethTx, err := cltest.FindEthTxWithAttempts(db, ethTx.ID)
			require.NoError(t, err)
			assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, ethTx.State)

			ethClient.AssertExpectations(t)
		})
		t.Run("with unknown error, sends tx as normal", func(t *testing.T) {
			ethTx := bulletprooftxmanager.EthTx{
				FromAddress:    fromAddress,
				ToAddress:      toAddress,
				EncodedPayload: []byte{42, 0, 0},
				Value:          assets.NewEthValue(542),
				GasLimit:       gasLimit,
				CreatedAt:      time.Unix(0, 0),
				State:          bulletprooftxmanager.EthTxUnstarted,
				Simulate:       true,
			}
			ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				return tx.Nonce() == uint64(6) && tx.Value().Cmp(big.NewInt(542)) == 0
			})).Return(nil).Once()
			ethClient.On("CallContext", mock.Anything, mock.AnythingOfType("*hexutil.Bytes"), "eth_call", mock.MatchedBy(func(callarg map[string]interface{}) bool {
				return fmt.Sprintf("%s", callarg["value"]) == "0x21e" // 542
			}), "latest").Return(errors.New("this is not a revert, something unexpected went wrong")).Once()

			require.NoError(t, db.Save(&ethTx).Error)

			require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

			ethTx, err := cltest.FindEthTxWithAttempts(db, ethTx.ID)
			require.NoError(t, err)
			assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, ethTx.State)

			ethClient.AssertExpectations(t)
		})
		t.Run("on revert, marks tx as fatally errored and does not send", func(t *testing.T) {
			ethTx := bulletprooftxmanager.EthTx{
				FromAddress:    fromAddress,
				ToAddress:      toAddress,
				EncodedPayload: []byte{42, 0, 0},
				Value:          assets.NewEthValue(642),
				GasLimit:       gasLimit,
				CreatedAt:      time.Unix(0, 0),
				State:          bulletprooftxmanager.EthTxUnstarted,
				Simulate:       true,
			}

			jerr := eth.JsonError{
				Code:    42,
				Message: "oh no, it reverted",
				Data:    []byte{42, 166, 34},
			}
			ethClient.On("CallContext", mock.Anything, mock.AnythingOfType("*hexutil.Bytes"), "eth_call", mock.MatchedBy(func(callarg map[string]interface{}) bool {
				return fmt.Sprintf("%s", callarg["value"]) == "0x282" // 642
			}), "latest").Return(&jerr).Once()

			require.NoError(t, db.Save(&ethTx).Error)

			require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

			ethTx, err := cltest.FindEthTxWithAttempts(db, ethTx.ID)
			require.NoError(t, err)
			assert.Equal(t, bulletprooftxmanager.EthTxFatalError, ethTx.State)
			assert.True(t, ethTx.Error.Valid)
			assert.Equal(t, "transaction reverted during simulation: json-rpc error { Code = 42, Message = 'oh no, it reverted', Data = 'KqYi' }", ethTx.Error.String)

			ethClient.AssertExpectations(t)
		})
	})

	ethClient.AssertExpectations(t)
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_OptimisticLockingOnEthTx(t *testing.T) {
	// non-transactional DB needed because we deliberately test for FK violation
	cfg, _, db := heavyweight.FullTestDB(t, "eth_broadcaster_optimistic_locking", true, true)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	chStartEstimate := make(chan struct{})
	chBlock := make(chan struct{})

	estimator := new(gasmocks.Estimator)
	estimator.On("GetLegacyGas", mock.Anything, mock.Anything).Return(assets.GWei(32), uint64(500), nil).Run(func(_ mock.Arguments) {
		close(chStartEstimate)
		<-chBlock
	})

	eb := bulletprooftxmanager.NewEthBroadcaster(
		db,
		ethClient,
		evmcfg,
		ethKeyStore,
		&postgres.NullEventBroadcaster{},
		[]ethkey.State{keyState},
		estimator,
		nil,
		logger.TestLogger(t),
	)

	etx := bulletprooftxmanager.EthTx{
		FromAddress:    fromAddress,
		ToAddress:      cltest.NewAddress(),
		EncodedPayload: []byte{42, 42, 0},
		Value:          *assets.NewEth(0),
		GasLimit:       500000,
		State:          bulletprooftxmanager.EthTxUnstarted,
	}
	require.NoError(t, db.Save(&etx).Error)

	go func() {
		select {
		case <-chStartEstimate:
		case <-time.After(5 * time.Second):
			t.Log("timed out waiting for estimator to be called")
			return
		}

		// Simulate a "PruneQueue" call
		assert.NoError(t, db.Exec(`DELETE FROM eth_txes WHERE state = 'unstarted'`).Error)

		close(chBlock)
	}()

	err := eb.ProcessUnstartedEthTxs(context.Background(), keyState)
	require.NoError(t, err)

	estimator.AssertExpectations(t)
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_Success_WithMultiplier(t *testing.T) {
	db := pgtest.NewGormDB(t)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	cfg := cltest.NewTestGeneralConfig(t)
	cfg.Overrides.GlobalEvmGasLimitMultiplier = null.FloatFrom(1.3)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

	eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState})

	ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
		assert.Equal(t, int(1600), int(tx.Gas()))
		return true
	})).Return(nil).Once()

	tx := bulletprooftxmanager.EthTx{
		FromAddress:    fromAddress,
		ToAddress:      gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411"),
		EncodedPayload: []byte{42, 42, 0},
		Value:          assets.NewEthValue(242),
		GasLimit:       1231,
		CreatedAt:      time.Unix(0, 0),
		State:          bulletprooftxmanager.EthTxUnstarted,
	}
	require.NoError(t, db.Save(&tx).Error)

	// Do the thing
	require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))
	ethClient.AssertExpectations(t)
}

func TestEthBroadcaster_AssignsNonceOnStart(t *testing.T) {
	var err error
	db := pgtest.NewGormDB(t)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	k1, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, true)
	k2, dummyAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, false)
	keyStates := []ethkey.State{k1, k2}

	cfg := cltest.NewTestGeneralConfig(t)
	cfg.Overrides.GlobalEvmNonceAutoSync = null.BoolFrom(true)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ethNodeNonce := uint64(22)

	t.Run("when eth node returns error", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, keyStates)

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(account gethCommon.Address) bool {
			return account.Hex() == dummyAddress.Hex()
		})).Return(uint64(0), nil).Once()
		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(account gethCommon.Address) bool {
			return account.Hex() == fromAddress.Hex()
		})).Return(ethNodeNonce, errors.New("something exploded")).Once()

		err = eb.Start()
		require.Error(t, err)
		defer eb.Close()
		require.Contains(t, err.Error(), "something exploded")

		// dummy address got updated
		var n int
		err := db.Raw(`SELECT next_nonce FROM eth_key_states WHERE address = ?`, dummyAddress).Scan(&n).Error
		require.NoError(t, err)
		require.Equal(t, 0, n)

		// real address did not update (it errored)
		err = db.Raw(`SELECT next_nonce FROM eth_key_states WHERE address = ?`, fromAddress).Scan(&n).Error
		require.NoError(t, err)
		require.Equal(t, 0, n)

		ethClient.AssertExpectations(t)
	})

	t.Run("when eth node returns nonce", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, keyStates)

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(account gethCommon.Address) bool {
			return account.Hex() == dummyAddress.Hex()
		})).Return(uint64(0), nil).Once()
		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(account gethCommon.Address) bool {
			return account.Hex() == fromAddress.Hex()
		})).Return(ethNodeNonce, nil).Once()

		require.NoError(t, eb.Start())
		defer eb.Close()

		// Check keyState to make sure it has correct nonce assigned
		var states []ethkey.State
		err := db.Order("created_at asc").Find(&states).Error
		require.NoError(t, err)
		state := states[0]

		assert.NotNil(t, state.NextNonce)
		assert.Equal(t, int64(ethNodeNonce), state.NextNonce)

		// The dummy key did not get updated
		state2 := states[1]
		assert.Equal(t, dummyAddress.Hex(), state2.Address.Hex())
		assert.Equal(t, 0, int(state2.NextNonce))

		ethClient.AssertExpectations(t)
	})
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_ResumingFromCrash(t *testing.T) {
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	value := assets.NewEthValue(142)
	gasLimit := uint64(242)
	encodedPayload := []byte{0, 1}
	nextNonce := int64(916714082576372851)
	firstNonce := nextNonce
	secondNonce := nextNonce + 1
	cfg := cltest.NewTestGeneralConfig(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	t.Run("cannot be more than one transaction per address in an unfinished state", func(t *testing.T) {
		db := pgtest.NewGormDB(t)

		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, nextNonce)

		firstInProgress := bulletprooftxmanager.EthTx{
			FromAddress:    fromAddress,
			Nonce:          &firstNonce,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			BroadcastAt:    nil,
			Error:          null.String{},
			State:          bulletprooftxmanager.EthTxInProgress,
		}

		secondInProgress := bulletprooftxmanager.EthTx{
			FromAddress:    fromAddress,
			Nonce:          &secondNonce,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			BroadcastAt:    nil,
			Error:          null.String{},
			State:          bulletprooftxmanager.EthTxInProgress,
		}

		require.NoError(t, db.Create(&firstInProgress).Error)
		err := db.Create(&secondInProgress).Error
		require.Error(t, err)
		assert.EqualError(t, err, "ERROR: duplicate key value violates unique constraint \"idx_only_one_in_progress_tx_per_account_id_per_evm_chain_id\" (SQLSTATE 23505)")
	})

	t.Run("previous run assigned nonce but never broadcast", func(t *testing.T) {
		db := pgtest.NewGormDB(t)

		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, nextNonce)

		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState})

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so eth_key_states.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, db, firstNonce, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check it was saved correctly with its attempt
		etx, err := cltest.FindEthTxWithAttempts(db, inProgressEthTx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.EthTxAttempts, 1)
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, etx.EthTxAttempts[0].State)

		ethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and broadcast but it fatally errored before we could save", func(t *testing.T) {
		db := pgtest.NewGormDB(t)

		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, nextNonce)

		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState})

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, db, firstNonce, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		})).Return(errors.New("exceeds block gas limit")).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check it was saved correctly with its attempt
		etx, err := cltest.FindEthTxWithAttempts(db, inProgressEthTx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.True(t, etx.Error.Valid)
		assert.Equal(t, "exceeds block gas limit", etx.Error.String)
		assert.Len(t, etx.EthTxAttempts, 0)

		ethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and broadcast and is now in mempool", func(t *testing.T) {
		db := pgtest.NewGormDB(t)

		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, nextNonce)

		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState})

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, db, firstNonce, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		})).Return(errors.New("known transaction: a1313bd99a81fb4d8ad1d2e90b67c6b3fa77545c990d6251444b83b70b6f8980")).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check it was saved correctly with its attempt
		etx, err := cltest.FindEthTxWithAttempts(db, inProgressEthTx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.EthTxAttempts, 1)

		ethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and broadcast and now the transaction has been confirmed", func(t *testing.T) {
		db := pgtest.NewGormDB(t)

		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, nextNonce)

		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState})

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, db, firstNonce, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		})).Return(errors.New("nonce too low")).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check it was saved correctly with its attempt
		etx, err := cltest.FindEthTxWithAttempts(db, inProgressEthTx.ID)
		require.NoError(t, err)

		require.NotNil(t, etx.BroadcastAt)
		assert.Equal(t, *etx.BroadcastAt, etx.CreatedAt)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.EthTxAttempts, 1)

		ethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and then failed to reach node for some reason and node is still down", func(t *testing.T) {
		failedToReachNodeError := context.DeadlineExceeded
		db := pgtest.NewGormDB(t)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, nextNonce)

		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState})

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, db, firstNonce, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		})).Return(failedToReachNodeError).Once()

		// Do the thing
		err := eb.ProcessUnstartedEthTxs(context.Background(), keyState)
		require.Error(t, err)
		assert.Contains(t, err.Error(), failedToReachNodeError.Error())

		// Check it was left in the unfinished state
		etx, err := cltest.FindEthTxWithAttempts(db, inProgressEthTx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.Equal(t, nextNonce, *etx.Nonce)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.EthTxAttempts, 1)

		ethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and broadcast transaction then crashed and rebooted with a different configured gas price", func(t *testing.T) {
		db := pgtest.NewGormDB(t)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, nextNonce)

		cfg := cltest.NewTestGeneralConfig(t)
		// Configured gas price changed
		cfg.Overrides.GlobalEvmGasPriceDefault = big.NewInt(500000000000)
		evmcfg := evmtest.NewChainScopedConfig(t, cfg)

		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState})

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, db, firstNonce, fromAddress)
		require.Len(t, inProgressEthTx.EthTxAttempts, 1)
		attempt := inProgressEthTx.EthTxAttempts[0]

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			// Ensure that the gas price is the same as the original attempt
			s, e := attempt.GetSignedTx()
			require.NoError(t, e)
			return tx.Nonce() == uint64(firstNonce) && tx.GasPrice().Int64() == s.GasPrice().Int64()
		})).Return(errors.New("known transaction: a1313bd99a81fb4d8ad1d2e90b67c6b3fa77545c990d6251444b83b70b6f8980")).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check it was saved correctly with its attempt
		etx, err := cltest.FindEthTxWithAttempts(db, inProgressEthTx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt = etx.EthTxAttempts[0]
		s, err := attempt.GetSignedTx()
		require.NoError(t, err)
		assert.Equal(t, int64(342), s.GasPrice().Int64())
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt.State)

		ethClient.AssertExpectations(t)
	})
}

func getLocalNextNonce(t *testing.T, db *gorm.DB, fromAddress gethCommon.Address) uint64 {
	n, err := bulletprooftxmanager.GetNextNonce(db, fromAddress, &cltest.FixtureChainID)
	require.NoError(t, err)
	require.NotNil(t, n)
	return uint64(n)
}

// Note that all of these tests share the same database, and ordering matters.
// This in order to more deeply test ProcessUnstartedEthTxs over
// multiple runs with previous errors in the database.
func TestEthBroadcaster_ProcessUnstartedEthTxs_Errors(t *testing.T) {
	var err error
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	value := assets.NewEthValue(142)
	gasLimit := uint64(242)
	encodedPayload := []byte{0, 1}

	db := pgtest.NewGormDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	cfg := cltest.NewTestGeneralConfig(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

	eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState})

	require.NoError(t, db.Exec(`SET CONSTRAINTS pipeline_runs_pipeline_spec_id_fkey DEFERRED`).Error)

	t.Run("if external wallet sent a transaction from the account and now the nonce is one higher than it should be and we got replacement underpriced then we assume a previous transaction of ours was the one that succeeded, and hand off to EthConfirmer", func(t *testing.T) {
		etx := bulletprooftxmanager.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          bulletprooftxmanager.EthTxUnstarted,
		}
		require.NoError(t, db.Save(&etx).Error)

		// First send, replacement underpriced
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(0)
		})).Return(errors.New("replacement transaction underpriced")).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		ethClient.AssertExpectations(t)

		// Check that the transaction was saved correctly with its attempt
		// We assume success and hand off to eth confirmer to eventually mark it as failed
		var latestID int64
		var etx1 bulletprooftxmanager.EthTx
		require.NoError(t, db.Raw("SELECT max(id) FROM eth_txes").Row().Scan(&latestID))
		etx1, err = cltest.FindEthTxWithAttempts(db, latestID)
		require.NoError(t, err)
		require.NotNil(t, etx1.BroadcastAt)
		assert.NotEqual(t, etx1.CreatedAt, *etx1.BroadcastAt)
		require.NotNil(t, etx1.Nonce)
		assert.Equal(t, int64(0), *etx1.Nonce)
		assert.False(t, etx1.Error.Valid)
		assert.Len(t, etx1.EthTxAttempts, 1)

		// Check that the local nonce was incremented by one
		var finalNextNonce int64
		finalNextNonce, err = bulletprooftxmanager.GetNextNonce(db, fromAddress, &cltest.FixtureChainID)
		require.NoError(t, err)
		require.NotNil(t, finalNextNonce)
		require.Equal(t, int64(1), finalNextNonce)
	})

	t.Run("geth client returns an error in the fatal errors category", func(t *testing.T) {
		fatalErrorExample := "exceeds block gas limit"
		localNextNonce := getLocalNextNonce(t, db, fromAddress)

		t.Run("without callback", func(t *testing.T) {
			etx := bulletprooftxmanager.EthTx{
				FromAddress:    fromAddress,
				ToAddress:      toAddress,
				EncodedPayload: encodedPayload,
				Value:          value,
				GasLimit:       gasLimit,
				State:          bulletprooftxmanager.EthTxUnstarted,
			}
			require.NoError(t, db.Save(&etx).Error)

			ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				return tx.Nonce() == localNextNonce
			})).Return(errors.New(fatalErrorExample)).Once()

			require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

			// Check it was saved correctly with its attempt
			etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
			require.NoError(t, err)

			assert.Nil(t, etx.BroadcastAt)
			require.Nil(t, etx.Nonce)
			assert.True(t, etx.Error.Valid)
			assert.Contains(t, etx.Error.String, "exceeds block gas limit")
			assert.Len(t, etx.EthTxAttempts, 0)

			// Check that the key had its nonce reset
			var state ethkey.State
			require.NoError(t, db.First(&state).Error)
			// Saved NextNonce must be the same as before because this transaction
			// was not accepted by the eth node and never can be
			require.NotNil(t, state.NextNonce)
			require.Equal(t, int64(localNextNonce), state.NextNonce)
		})

		t.Run("with callback", func(t *testing.T) {
			run := cltest.MustInsertPipelineRun(t, db)
			tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)
			etx := bulletprooftxmanager.EthTx{
				FromAddress:       fromAddress,
				ToAddress:         toAddress,
				EncodedPayload:    encodedPayload,
				Value:             value,
				GasLimit:          gasLimit,
				State:             bulletprooftxmanager.EthTxUnstarted,
				PipelineTaskRunID: uuid.NullUUID{UUID: tr.ID, Valid: true},
			}

			t.Run("with erroring callback bails out", func(t *testing.T) {
				require.NoError(t, db.Save(&etx).Error)
				fn := func(id uuid.UUID, result interface{}, err error) error {
					return errors.New("something exploded in the callback")
				}

				bulletprooftxmanager.SetResumeCallbackOnEthBroadcaster(fn, eb)

				ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
					return tx.Nonce() == localNextNonce
				})).Return(errors.New(fatalErrorExample)).Once()

				err := eb.ProcessUnstartedEthTxs(context.Background(), keyState)
				require.Error(t, err)
				require.Contains(t, err.Error(), "something exploded in the callback")
			})

			t.Run("calls resume with error", func(t *testing.T) {
				fn := func(id uuid.UUID, result interface{}, err error) error {
					require.Equal(t, id, tr.ID)
					require.Nil(t, result)
					require.Error(t, err)
					require.Contains(t, err.Error(), "fatal error while sending transaction: exceeds block gas limit")
					return nil
				}

				bulletprooftxmanager.SetResumeCallbackOnEthBroadcaster(fn, eb)

				ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
					return tx.Nonce() == localNextNonce
				})).Return(errors.New(fatalErrorExample)).Once()

				require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))
			})

		})

		ethClient.AssertExpectations(t)
	})

	bulletprooftxmanager.SetResumeCallbackOnEthBroadcaster(nil, eb)

	t.Run("geth client fails with error indicating that the transaction was too expensive", func(t *testing.T) {
		tooExpensiveError := "tx fee (1.10 ether) exceeds the configured cap (1.00 ether)"
		localNextNonce := getLocalNextNonce(t, db, fromAddress)

		etx := bulletprooftxmanager.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          bulletprooftxmanager.EthTxUnstarted,
		}
		require.NoError(t, db.Save(&etx).Error)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New(tooExpensiveError)).Once()

		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check it was saved with no attempt and a fatal error
		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		require.Nil(t, etx.Nonce)
		assert.True(t, etx.Error.Valid)
		assert.Contains(t, etx.Error.String, "tx fee (1.10 ether) exceeds the configured cap (1.00 ether)")
		assert.Len(t, etx.EthTxAttempts, 0)

		// Check that the key had its nonce reset
		var state ethkey.State
		require.NoError(t, db.First(&state).Error)
		// Saved NextNonce must be the same as before because this transaction
		// was not accepted by the eth node and never can be
		require.NotNil(t, state.NextNonce)
		require.Equal(t, int64(localNextNonce), state.NextNonce)

		ethClient.AssertExpectations(t)
	})

	t.Run("eth client call fails with an unexpected random error", func(t *testing.T) {
		retryableErrorExample := "geth shit the bed again"
		localNextNonce := getLocalNextNonce(t, db, fromAddress)

		etx := bulletprooftxmanager.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          bulletprooftxmanager.EthTxUnstarted,
		}
		require.NoError(t, db.Save(&etx).Error)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New(retryableErrorExample)).Once()

		// Do the thing
		err = eb.ProcessUnstartedEthTxs(context.Background(), keyState)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("error while sending transaction %v: %s", etx.ID, retryableErrorExample))

		// Check it was saved correctly with its attempt
		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		require.NotNil(t, etx.Nonce)
		assert.False(t, etx.Error.Valid)
		assert.Equal(t, bulletprooftxmanager.EthTxInProgress, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt := etx.EthTxAttempts[0]
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptInProgress, attempt.State)

		ethClient.AssertExpectations(t)

		// Now on the second run, it is successful
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(nil).Once()

		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check it was saved correctly with its attempt
		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		require.NotNil(t, etx.Nonce)
		assert.False(t, etx.Error.Valid)
		assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt = etx.EthTxAttempts[0]
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("eth node returns underpriced transaction", func(t *testing.T) {
		// This happens if a transaction's gas price is below the minimum
		// configured for the transaction pool.
		// This is a configuration error by the node operator, since it means they set the base gas level too low.
		underpricedError := "transaction underpriced"
		localNextNonce := getLocalNextNonce(t, db, fromAddress)

		etx := bulletprooftxmanager.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          bulletprooftxmanager.EthTxUnstarted,
		}
		require.NoError(t, db.Save(&etx).Error)

		// First was underpriced
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasPrice().Cmp(evmcfg.EvmGasPriceDefault()) == 0
		})).Return(errors.New(underpricedError)).Once()

		// Second with gas bump was still underpriced
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasPrice().Cmp(big.NewInt(25000000000)) == 0
		})).Return(errors.New(underpricedError)).Once()

		// Third succeeded
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasPrice().Cmp(big.NewInt(30000000000)) == 0
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		ethClient.AssertExpectations(t)

		// Check it was saved correctly with its attempt
		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		require.NotNil(t, etx.Nonce)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt := etx.EthTxAttempts[0]
		assert.Equal(t, big.NewInt(30000000000).String(), attempt.GasPrice.String())
	})

	etxUnfinished := bulletprooftxmanager.EthTx{
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: encodedPayload,
		Value:          value,
		GasLimit:       gasLimit,
		State:          bulletprooftxmanager.EthTxUnstarted,
	}
	require.NoError(t, db.Save(&etxUnfinished).Error)

	t.Run("failed to reach node for some reason", func(t *testing.T) {
		failedToReachNodeError := context.DeadlineExceeded
		localNextNonce := getLocalNextNonce(t, db, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(failedToReachNodeError).Once()

		// Do the thing
		err = eb.ProcessUnstartedEthTxs(context.Background(), keyState)
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("error while sending transaction %v: context deadline exceeded", etxUnfinished.ID))

		// Check it was left in the unfinished state
		etx, err := cltest.FindEthTxWithAttempts(db, etxUnfinished.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.Nonce)
		assert.False(t, etx.Error.Valid)
		assert.Equal(t, bulletprooftxmanager.EthTxInProgress, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptInProgress, etx.EthTxAttempts[0].State)

		ethClient.AssertExpectations(t)
	})

	t.Run("eth node returns temporarily underpriced transaction", func(t *testing.T) {
		// This happens if parity is rejecting transactions that are not priced high enough to even get into the mempool at all
		// It should pretend it was accepted into the mempool and hand off to ethConfirmer to bump gas as normal
		temporarilyUnderpricedError := "There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee."
		localNextNonce := getLocalNextNonce(t, db, fromAddress)

		// Re-use the previously unfinished transaction, no need to insert new

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New(temporarilyUnderpricedError)).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check it was saved correctly with its attempt
		etx, err := cltest.FindEthTxWithAttempts(db, etxUnfinished.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		require.NotNil(t, etx.Nonce)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt := etx.EthTxAttempts[0]
		assert.Equal(t, big.NewInt(20000000000).String(), attempt.GasPrice.String())

		ethClient.AssertExpectations(t)
	})

	t.Run("eth node returns underpriced transaction and bumping gas doesn't increase it", func(t *testing.T) {
		// This happens if a transaction's gas price is below the minimum
		// configured for the transaction pool.
		// This is a configuration error by the node operator, since it means they set the base gas level too low.
		underpricedError := "transaction underpriced"
		localNextNonce := getLocalNextNonce(t, db, fromAddress)
		// In this scenario the node operator REALLY fucked up and set the bump
		// to zero (even though that should not be possible due to config
		// validation)
		cfg.Overrides.GlobalEvmGasBumpWei = big.NewInt(0)
		cfg.Overrides.GlobalEvmGasBumpPercent = null.IntFrom(0)

		etx := bulletprooftxmanager.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          bulletprooftxmanager.EthTxUnstarted,
		}
		require.NoError(t, db.Save(&etx).Error)

		// First was underpriced
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasPrice().Cmp(evmcfg.EvmGasPriceDefault()) == 0
		})).Return(errors.New(underpricedError)).Once()

		// Do the thing
		err := eb.ProcessUnstartedEthTxs(context.Background(), keyState)
		require.Error(t, err)
		require.Contains(t, err.Error(), "bumped gas price of 20000000000 is equal to original gas price of 20000000000. ACTION REQUIRED: This is a configuration error, you must increase either ETH_GAS_BUMP_PERCENT or ETH_GAS_BUMP_WEI")

		// TEARDOWN: Clear out the unsent tx before the next test
		require.NoError(t, db.Exec(`DELETE FROM eth_txes WHERE nonce = ?`, localNextNonce).Error)

		ethClient.AssertExpectations(t)
	})

	t.Run("eth tx is left in progress if eth node returns insufficient eth", func(t *testing.T) {
		insufficientEthError := "insufficient funds for transfer"
		localNextNonce := getLocalNextNonce(t, db, fromAddress)
		etx := bulletprooftxmanager.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          bulletprooftxmanager.EthTxUnstarted,
		}
		require.NoError(t, db.Save(&etx).Error)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New(insufficientEthError)).Once()

		err := eb.ProcessUnstartedEthTxs(context.Background(), keyState)
		require.EqualError(t, err, "processUnstartedEthTxs failed: insufficient funds for transfer")

		// Check it was saved correctly with its attempt
		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		require.NotNil(t, etx.Nonce)
		assert.False(t, etx.Error.Valid)
		assert.Equal(t, bulletprooftxmanager.EthTxInProgress, etx.State)
		require.Len(t, etx.EthTxAttempts, 1)
		attempt := etx.EthTxAttempts[0]
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptInProgress, attempt.State)
		assert.Nil(t, attempt.BroadcastBeforeBlockNum)

		ethClient.AssertExpectations(t)
	})

	require.NoError(t, db.Exec(`DELETE FROM eth_txes`).Error)
	cfg.Overrides.GlobalEvmEIP1559DynamicFees = null.BoolFrom(true)

	t.Run("eth node returns underpriced transaction for EIP-1559 tx, should return error", func(t *testing.T) {
		// Experimentally this error is not actually possible; eth nodes will accept literally any price for EIP-1559 transactions
		underpricedError := "transaction underpriced"
		localNextNonce := getLocalNextNonce(t, db, fromAddress)

		etx := bulletprooftxmanager.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          bulletprooftxmanager.EthTxUnstarted,
		}
		require.NoError(t, db.Save(&etx).Error)

		// First was underpriced
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New(underpricedError)).Once()

		err := eb.ProcessUnstartedEthTxs(context.Background(), keyState)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "bumping gas on initial send is not supported for EIP-1559 transactions")

		ethClient.AssertExpectations(t)
	})
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_KeystoreErrors(t *testing.T) {
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	value := assets.NewEthValue(142)
	gasLimit := uint64(242)
	encodedPayload := []byte{0, 1}
	localNonce := 0

	db := pgtest.NewGormDB(t)

	realKeystore := cltest.NewKeyStore(t, db)
	keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, realKeystore.Eth())

	cfg := cltest.NewTestGeneralConfig(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

	kst := new(ksmocks.Eth)
	eb := cltest.NewEthBroadcaster(t, db, ethClient, kst, evmcfg, []ethkey.State{keyState})

	t.Run("tx signing fails", func(t *testing.T) {
		etx := bulletprooftxmanager.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          bulletprooftxmanager.EthTxUnstarted,
		}
		require.NoError(t, db.Save(&etx).Error)

		tx := *gethTypes.NewTx(&gethTypes.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.AnythingOfType("*types.Transaction"),
			mock.MatchedBy(func(chainID *big.Int) bool {
				return chainID.Cmp(evmcfg.ChainID()) == 0
			})).Return(&tx, errors.New("could not sign transaction")).Once()
		kst.On("GetState", fromAddress.Hex()).Return(ethkey.State{}, nil)

		// Do the thing
		err := eb.ProcessUnstartedEthTxs(context.Background(), keyState)
		require.Error(t, err)
		require.Contains(t, err.Error(), "could not sign transaction")

		// Check that the transaction is left in unstarted state
		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)

		assert.Equal(t, bulletprooftxmanager.EthTxUnstarted, etx.State)
		assert.Len(t, etx.EthTxAttempts, 0)

		// Check that the key did not have its nonce incremented
		var keyState ethkey.State
		require.NoError(t, db.First(&keyState).Error)
		require.NotNil(t, keyState.NextNonce)
		require.Equal(t, int64(localNonce), keyState.NextNonce)

		kst.AssertExpectations(t)
	})

	// Should have done nothing
	ethClient.AssertExpectations(t)
}

func TestEthBroadcaster_GetNextNonce(t *testing.T) {
	db := pgtest.NewGormDB(t)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	keyState, _ := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	nonce, err := bulletprooftxmanager.GetNextNonce(db, keyState.Address.Address(), &cltest.FixtureChainID)
	assert.NoError(t, err)
	require.NotNil(t, nonce)
	assert.Equal(t, int64(0), nonce)
}

func TestEthBroadcaster_IncrementNextNonce(t *testing.T) {
	db := pgtest.NewGormDB(t)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	keyState, _ := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	// Cannot increment if supplied nonce doesn't match existing
	require.Error(t, bulletprooftxmanager.IncrementNextNonce(db, keyState.Address.Address(), &cltest.FixtureChainID, int64(42)))

	require.NoError(t, bulletprooftxmanager.IncrementNextNonce(db, keyState.Address.Address(), &cltest.FixtureChainID, int64(0)))

	// Nonce bumped to 1
	require.NoError(t, db.First(&keyState).Error)
	require.NotNil(t, keyState.NextNonce)
	require.Equal(t, int64(1), keyState.NextNonce)
}

func TestEthBroadcaster_Trigger(t *testing.T) {
	t.Parallel()

	// Simple sanity check to make sure it doesn't block
	db := pgtest.NewGormDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	eb := cltest.NewEthBroadcaster(t, db, cltest.NewEthClientMockWithDefaultChain(t), ethKeyStore, evmcfg, []ethkey.State{})

	eb.Trigger(cltest.NewAddress())
	eb.Trigger(cltest.NewAddress())
	eb.Trigger(cltest.NewAddress())
}

func TestEthBroadcaster_EthTxInsertEventCausesTriggerToFire(t *testing.T) {
	// NOTE: Testing triggers requires committing transactions and does not work with transactional tests
	cfg, _, db := heavyweight.FullTestDB(t, "eth_tx_triggers", true, true)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	eventBroadcaster := postgres.NewEventBroadcaster(evmcfg.DatabaseURL(), 0, 0, logger.TestLogger(t))
	require.NoError(t, eventBroadcaster.Start())
	t.Cleanup(func() { require.NoError(t, eventBroadcaster.Close()) })

	ethTxInsertListener, err := eventBroadcaster.Subscribe(postgres.ChannelInsertOnEthTx, "")
	require.NoError(t, err)

	// Give it some time to start listening
	time.Sleep(100 * time.Millisecond)

	mustInsertUnstartedEthTx(t, db, fromAddress)
	gomega.NewGomegaWithT(t).Eventually(ethTxInsertListener.Events()).Should(gomega.Receive())
}
