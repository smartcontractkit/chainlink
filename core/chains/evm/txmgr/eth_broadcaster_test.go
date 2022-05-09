package txmgr_test

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

	"github.com/smartcontractkit/chainlink/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	gasmocks "github.com/smartcontractkit/chainlink/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pg/datatypes"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestEthBroadcaster_ProcessUnstartedEthTxs_Success(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	checkerFactory := &txmgr.CheckerFactory{Client: ethClient}

	eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState}, checkerFactory)

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

		etx := txmgr.EthTx{
			FromAddress:    otherAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          txmgr.EthTxUnstarted,
		}
		require.NoError(t, borm.InsertEthTx(&etx))

		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))
	})

	t.Run("existing eth_txes with broadcast_at or error", func(t *testing.T) {
		nonce := int64(342)
		errStr := "some error"

		etxUnconfirmed := txmgr.EthTx{
			Nonce:              &nonce,
			FromAddress:        fromAddress,
			ToAddress:          toAddress,
			EncodedPayload:     encodedPayload,
			Value:              value,
			GasLimit:           gasLimit,
			BroadcastAt:        &timeNow,
			InitialBroadcastAt: &timeNow,
			Error:              null.String{},
			State:              txmgr.EthTxUnconfirmed,
		}
		etxWithError := txmgr.EthTx{
			Nonce:          nil,
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			Error:          null.StringFrom(errStr),
			State:          txmgr.EthTxFatalError,
		}

		require.NoError(t, borm.InsertEthTx(&etxUnconfirmed))
		require.NoError(t, borm.InsertEthTx(&etxWithError))

		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))
	})

	t.Run("sends 3 EthTxs in order with higher value last, and lower values starting from the earliest", func(t *testing.T) {
		// Higher value
		expensiveEthTx := txmgr.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 42, 0},
			Value:          assets.NewEthValue(242),
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(0, 0),
			State:          txmgr.EthTxUnstarted,
		}
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(2) && tx.Value().Cmp(big.NewInt(242)) == 0
		})).Return(nil).Once()

		// Earlier
		tr := int32(99)
		b, err := json.Marshal(txmgr.EthTxMeta{JobID: tr})
		require.NoError(t, err)
		meta := datatypes.JSON(b)
		earlierEthTx := txmgr.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 42, 0},
			Value:          value,
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(0, 1),
			State:          txmgr.EthTxUnstarted,
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
		laterEthTx := txmgr.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 42, 1},
			Value:          value,
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(1, 0),
			State:          txmgr.EthTxUnstarted,
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
		require.NoError(t, borm.InsertEthTx(&expensiveEthTx))
		require.NoError(t, borm.InsertEthTx(&laterEthTx))
		require.NoError(t, borm.InsertEthTx(&earlierEthTx))

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check earlierEthTx and it's attempt
		// This was the earlier one sent so it has the lower nonce
		earlierTransaction, err := borm.FindEthTxWithAttempts(earlierEthTx.ID)
		require.NoError(t, err)
		assert.False(t, earlierTransaction.Error.Valid)
		require.NotNil(t, earlierTransaction.FromAddress)
		assert.Equal(t, fromAddress, earlierTransaction.FromAddress)
		require.NotNil(t, earlierTransaction.Nonce)
		assert.Equal(t, int64(0), *earlierTransaction.Nonce)
		assert.NotNil(t, earlierTransaction.BroadcastAt)
		assert.NotNil(t, earlierTransaction.InitialBroadcastAt)
		assert.Len(t, earlierTransaction.EthTxAttempts, 1)
		var m txmgr.EthTxMeta
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
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, attempt.State)
		require.Len(t, attempt.EthReceipts, 0)

		// Check laterEthTx and it's attempt
		// This was the later one sent so it has the higher nonce
		laterTransaction, err := borm.FindEthTxWithAttempts(laterEthTx.ID)
		require.NoError(t, err)
		assert.False(t, earlierTransaction.Error.Valid)
		require.NotNil(t, laterTransaction.FromAddress)
		assert.Equal(t, fromAddress, laterTransaction.FromAddress)
		require.NotNil(t, laterTransaction.Nonce)
		assert.Equal(t, int64(1), *laterTransaction.Nonce)
		assert.NotNil(t, laterTransaction.BroadcastAt)
		assert.NotNil(t, earlierTransaction.InitialBroadcastAt)
		assert.Len(t, laterTransaction.EthTxAttempts, 1)

		attempt = laterTransaction.EthTxAttempts[0]

		assert.Equal(t, laterTransaction.ID, attempt.EthTxID)
		assert.Equal(t, evmcfg.EvmGasPriceDefault().String(), attempt.GasPrice.String())

		_, err = attempt.GetSignedTx()
		require.NoError(t, err)
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, attempt.State)
		require.Len(t, attempt.EthReceipts, 0)

		ethClient.AssertExpectations(t)
	})

	t.Run("sends transactions with type 0x2 in EIP-1559 mode", func(t *testing.T) {
		cfg.Overrides.GlobalEvmEIP1559DynamicFees = null.BoolFrom(true)
		rnd := int64(1000000000 + rand.Intn(5000))
		cfg.Overrides.GlobalEvmGasTipCapDefault = big.NewInt(rnd)
		cfg.Overrides.GlobalEvmGasFeeCapDefault = big.NewInt(rnd + 1)
		cfg.Overrides.GlobalEvmMaxGasPriceWei = big.NewInt(rnd + 2)

		eipTxWithoutAl := txmgr.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 0, 0},
			Value:          assets.NewEthValue(142),
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(0, 0),
			State:          txmgr.EthTxUnstarted,
		}
		eipTxWithAl := txmgr.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 42, 0},
			Value:          assets.NewEthValue(242),
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(0, 1),
			State:          txmgr.EthTxUnstarted,
			AccessList:     txmgr.NullableEIP2930AccessListFrom(gethTypes.AccessList{gethTypes.AccessTuple{Address: testutils.NewAddress(), StorageKeys: []gethCommon.Hash{utils.NewHash()}}}),
		}
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(3) && tx.Value().Cmp(big.NewInt(142)) == 0
		})).Return(nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(4) && tx.Value().Cmp(big.NewInt(242)) == 0
		})).Return(nil).Once()

		require.NoError(t, borm.InsertEthTx(&eipTxWithAl))
		require.NoError(t, borm.InsertEthTx(&eipTxWithoutAl))

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check eipTxWithAl and it's attempt
		// This was the earlier one sent so it has the lower nonce
		eipTxWithAl, err := borm.FindEthTxWithAttempts(eipTxWithAl.ID)
		require.NoError(t, err)
		assert.False(t, eipTxWithAl.Error.Valid)
		require.NotNil(t, eipTxWithAl.FromAddress)
		assert.Equal(t, fromAddress, eipTxWithAl.FromAddress)
		require.NotNil(t, eipTxWithAl.Nonce)
		assert.Equal(t, int64(4), *eipTxWithAl.Nonce)
		assert.NotNil(t, eipTxWithAl.BroadcastAt)
		assert.NotNil(t, eipTxWithAl.InitialBroadcastAt)
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
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, attempt.State)
		require.Len(t, attempt.EthReceipts, 0)
	})

	ethClient.AssertExpectations(t)

	t.Run("transaction simulation", func(t *testing.T) {
		t.Run("when simulation succeeds, sends tx as normal", func(t *testing.T) {
			ethTx := txmgr.EthTx{
				FromAddress:    fromAddress,
				ToAddress:      toAddress,
				EncodedPayload: []byte{42, 0, 0},
				Value:          assets.NewEthValue(442),
				GasLimit:       gasLimit,
				CreatedAt:      time.Unix(0, 0),
				State:          txmgr.EthTxUnstarted,
				TransmitChecker: checkerToJson(t, txmgr.TransmitCheckerSpec{
					CheckerType: txmgr.TransmitCheckerTypeSimulate,
				}),
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

			require.NoError(t, borm.InsertEthTx(&ethTx))

			require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

			// Check ethtx was sent
			ethTx, err := borm.FindEthTxWithAttempts(ethTx.ID)
			require.NoError(t, err)
			assert.Equal(t, txmgr.EthTxUnconfirmed, ethTx.State)

			ethClient.AssertExpectations(t)
		})
		t.Run("with unknown error, sends tx as normal", func(t *testing.T) {
			ethTx := txmgr.EthTx{
				FromAddress:    fromAddress,
				ToAddress:      toAddress,
				EncodedPayload: []byte{42, 0, 0},
				Value:          assets.NewEthValue(542),
				GasLimit:       gasLimit,
				CreatedAt:      time.Unix(0, 0),
				State:          txmgr.EthTxUnstarted,
				TransmitChecker: checkerToJson(t, txmgr.TransmitCheckerSpec{
					CheckerType: txmgr.TransmitCheckerTypeSimulate,
				}),
			}
			ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				return tx.Nonce() == uint64(6) && tx.Value().Cmp(big.NewInt(542)) == 0
			})).Return(nil).Once()
			ethClient.On("CallContext", mock.Anything, mock.AnythingOfType("*hexutil.Bytes"), "eth_call", mock.MatchedBy(func(callarg map[string]interface{}) bool {
				return fmt.Sprintf("%s", callarg["value"]) == "0x21e" // 542
			}), "latest").Return(errors.New("this is not a revert, something unexpected went wrong")).Once()

			require.NoError(t, borm.InsertEthTx(&ethTx))

			require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

			ethTx, err := borm.FindEthTxWithAttempts(ethTx.ID)
			require.NoError(t, err)
			assert.Equal(t, txmgr.EthTxUnconfirmed, ethTx.State)

			ethClient.AssertExpectations(t)
		})
		t.Run("on revert, marks tx as fatally errored and does not send", func(t *testing.T) {
			ethTx := txmgr.EthTx{
				FromAddress:    fromAddress,
				ToAddress:      toAddress,
				EncodedPayload: []byte{42, 0, 0},
				Value:          assets.NewEthValue(642),
				GasLimit:       gasLimit,
				CreatedAt:      time.Unix(0, 0),
				State:          txmgr.EthTxUnstarted,
				TransmitChecker: checkerToJson(t, txmgr.TransmitCheckerSpec{
					CheckerType: txmgr.TransmitCheckerTypeSimulate,
				}),
			}

			jerr := evmclient.JsonError{
				Code:    42,
				Message: "oh no, it reverted",
				Data:    []byte{42, 166, 34},
			}
			ethClient.On("CallContext", mock.Anything, mock.AnythingOfType("*hexutil.Bytes"), "eth_call", mock.MatchedBy(func(callarg map[string]interface{}) bool {
				return fmt.Sprintf("%s", callarg["value"]) == "0x282" // 642
			}), "latest").Return(&jerr).Once()

			require.NoError(t, borm.InsertEthTx(&ethTx))

			require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

			ethTx, err := borm.FindEthTxWithAttempts(ethTx.ID)
			require.NoError(t, err)
			assert.Equal(t, txmgr.EthTxFatalError, ethTx.State)
			assert.True(t, ethTx.Error.Valid)
			assert.Equal(t, "transaction reverted during simulation: json-rpc error { Code = 42, Message = 'oh no, it reverted', Data = 'KqYi' }", ethTx.Error.String)

			ethClient.AssertExpectations(t)
		})
	})

	ethClient.AssertExpectations(t)
}

func TestEthBroadcaster_TransmitChecking(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	checkerFactory := &testCheckerFactory{}

	eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState}, checkerFactory)

	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	gasLimit := uint64(242)

	t.Run("when transmit checking times out, sends tx as normal", func(t *testing.T) {
		// Checker will return a canceled error
		checkerFactory.err = context.Canceled

		ethTx := txmgr.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 0, 0},
			Value:          assets.NewEthValue(442),
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(0, 0),
			State:          txmgr.EthTxUnstarted,
			TransmitChecker: checkerToJson(t, txmgr.TransmitCheckerSpec{
				CheckerType: txmgr.TransmitCheckerTypeSimulate,
			}),
		}
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == 0 && tx.Value().Cmp(big.NewInt(442)) == 0
		})).Return(nil).Once()

		require.NoError(t, borm.InsertEthTx(&ethTx))
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check ethtx was sent
		ethTx, err := borm.FindEthTxWithAttempts(ethTx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgr.EthTxUnconfirmed, ethTx.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("when transmit checking succeeds, sends tx as normal", func(t *testing.T) {
		// Checker will return no error
		checkerFactory.err = nil

		ethTx := txmgr.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 0, 0},
			Value:          assets.NewEthValue(442),
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(0, 0),
			State:          txmgr.EthTxUnstarted,
			TransmitChecker: checkerToJson(t, txmgr.TransmitCheckerSpec{
				CheckerType: txmgr.TransmitCheckerTypeSimulate,
			}),
		}
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == 1 && tx.Value().Cmp(big.NewInt(442)) == 0
		})).Return(nil).Once()

		require.NoError(t, borm.InsertEthTx(&ethTx))
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check ethtx was sent
		ethTx, err := borm.FindEthTxWithAttempts(ethTx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgr.EthTxUnconfirmed, ethTx.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("when transmit errors, fatally error transaction", func(t *testing.T) {
		// Checker will return a fatal error
		checkerFactory.err = errors.New("fatal checker error")

		ethTx := txmgr.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 0, 0},
			Value:          assets.NewEthValue(442),
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(0, 0),
			State:          txmgr.EthTxUnstarted,
			TransmitChecker: checkerToJson(t, txmgr.TransmitCheckerSpec{
				CheckerType: txmgr.TransmitCheckerTypeSimulate,
			}),
		}

		require.NoError(t, borm.InsertEthTx(&ethTx))
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check ethtx was sent
		ethTx, err := borm.FindEthTxWithAttempts(ethTx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgr.EthTxFatalError, ethTx.State)
		assert.True(t, ethTx.Error.Valid)
		assert.Equal(t, "fatal checker error", ethTx.Error.String)

		ethClient.AssertExpectations(t)
	})
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_OptimisticLockingOnEthTx(t *testing.T) {
	// non-transactional DB needed because we deliberately test for FK violation
	cfg, db := heavyweight.FullTestDB(t, "eth_broadcaster_optimistic_locking")
	borm := cltest.NewTxmORM(t, db, cfg)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	chStartEstimate := make(chan struct{})
	chBlock := make(chan struct{})

	estimator := new(gasmocks.Estimator)
	estimator.On("GetLegacyGas", mock.Anything, mock.Anything).Return(assets.GWei(32), uint64(500), nil).Run(func(_ mock.Arguments) {
		close(chStartEstimate)
		<-chBlock
	})

	eb := txmgr.NewEthBroadcaster(
		db,
		ethClient,
		evmcfg,
		ethKeyStore,
		&pg.NullEventBroadcaster{},
		[]ethkey.State{keyState},
		estimator,
		nil,
		logger.TestLogger(t),
		&testCheckerFactory{},
	)

	etx := txmgr.EthTx{
		FromAddress:    fromAddress,
		ToAddress:      testutils.NewAddress(),
		EncodedPayload: []byte{42, 42, 0},
		Value:          *assets.NewEth(0),
		GasLimit:       500000,
		State:          txmgr.EthTxUnstarted,
	}
	require.NoError(t, borm.InsertEthTx(&etx))

	go func() {
		select {
		case <-chStartEstimate:
		case <-time.After(5 * time.Second):
			t.Log("timed out waiting for estimator to be called")
			return
		}

		// Simulate a "PruneQueue" call
		assert.NoError(t, utils.JustError(db.Exec(`DELETE FROM eth_txes WHERE state = 'unstarted'`)))

		close(chBlock)
	}()

	err := eb.ProcessUnstartedEthTxs(context.Background(), keyState)
	require.NoError(t, err)

	estimator.AssertExpectations(t)
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_Success_WithMultiplier(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	cfg.Overrides.GlobalEvmGasLimitMultiplier = null.FloatFrom(1.3)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

	eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState}, &testCheckerFactory{})

	ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
		assert.Equal(t, int(1600), int(tx.Gas()))
		return true
	})).Return(nil).Once()

	tx := txmgr.EthTx{
		FromAddress:    fromAddress,
		ToAddress:      gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411"),
		EncodedPayload: []byte{42, 42, 0},
		Value:          assets.NewEthValue(242),
		GasLimit:       1231,
		CreatedAt:      time.Unix(0, 0),
		State:          txmgr.EthTxUnstarted,
	}
	require.NoError(t, borm.InsertEthTx(&tx))

	// Do the thing
	require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))
	ethClient.AssertExpectations(t)
}

func TestEthBroadcaster_AssignsNonceOnStart(t *testing.T) {
	var err error
	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	k1, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, true)
	k2, dummyAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, false)
	keyStates := []ethkey.State{k1, k2}

	cfg.Overrides.GlobalEvmNonceAutoSync = null.BoolFrom(true)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ethNodeNonce := uint64(22)

	t.Run("when eth node returns error", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, keyStates, &testCheckerFactory{})

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(account gethCommon.Address) bool {
			return account.Hex() == dummyAddress.Hex()
		})).Return(uint64(0), nil).Once()
		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(account gethCommon.Address) bool {
			return account.Hex() == fromAddress.Hex()
		})).Return(ethNodeNonce, errors.New("something exploded")).Once()

		err = eb.Start(testutils.Context(t))
		require.Error(t, err)
		defer eb.Close()
		require.Contains(t, err.Error(), "something exploded")

		// dummy address got updated
		var n int
		err := db.Get(&n, `SELECT next_nonce FROM eth_key_states WHERE address = $1`, dummyAddress)
		require.NoError(t, err)
		require.Equal(t, 0, n)

		// real address did not update (it errored)
		err = db.Get(&n, `SELECT next_nonce FROM eth_key_states WHERE address = $1`, fromAddress)
		require.NoError(t, err)
		require.Equal(t, 0, n)

		ethClient.AssertExpectations(t)
	})

	t.Run("when eth node returns nonce", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, keyStates, &testCheckerFactory{})

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(account gethCommon.Address) bool {
			return account.Hex() == dummyAddress.Hex()
		})).Return(uint64(0), nil).Once()
		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(account gethCommon.Address) bool {
			return account.Hex() == fromAddress.Hex()
		})).Return(ethNodeNonce, nil).Once()

		require.NoError(t, eb.Start(testutils.Context(t)))
		defer eb.Close()

		// Check keyState to make sure it has correct nonce assigned
		var states []ethkey.State
		err := db.Select(&states, `SELECT * FROM eth_key_states ORDER BY created_at ASC, id ASC`)
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
		db := pgtest.NewSqlxDB(t)
		borm := cltest.NewTxmORM(t, db, cfg)

		ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, nextNonce)

		firstInProgress := txmgr.EthTx{
			FromAddress:    fromAddress,
			Nonce:          &firstNonce,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			Error:          null.String{},
			State:          txmgr.EthTxInProgress,
		}

		secondInProgress := txmgr.EthTx{
			FromAddress:    fromAddress,
			Nonce:          &secondNonce,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			Error:          null.String{},
			State:          txmgr.EthTxInProgress,
		}

		require.NoError(t, borm.InsertEthTx(&firstInProgress))
		err := borm.InsertEthTx(&secondInProgress)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ERROR: duplicate key value violates unique constraint \"idx_only_one_in_progress_tx_per_account_id_per_evm_chain_id\" (SQLSTATE 23505)")
	})

	t.Run("previous run assigned nonce but never broadcast", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		borm := cltest.NewTxmORM(t, db, cfg)

		ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
		keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, nextNonce)

		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState}, &testCheckerFactory{})

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so eth_key_states.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, borm, firstNonce, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check it was saved correctly with its attempt
		etx, err := borm.FindEthTxWithAttempts(inProgressEthTx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.InitialBroadcastAt)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.EthTxAttempts, 1)
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, etx.EthTxAttempts[0].State)

		ethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and broadcast but it fatally errored before we could save", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		borm := cltest.NewTxmORM(t, db, cfg)

		ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
		keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, nextNonce)

		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState}, &testCheckerFactory{})

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, borm, firstNonce, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		})).Return(errors.New("exceeds block gas limit")).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check it was saved correctly with its attempt
		etx, err := borm.FindEthTxWithAttempts(inProgressEthTx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.Nil(t, etx.InitialBroadcastAt)
		assert.True(t, etx.Error.Valid)
		assert.Equal(t, "exceeds block gas limit", etx.Error.String)
		assert.Len(t, etx.EthTxAttempts, 0)

		ethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and broadcast and is now in mempool", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		borm := cltest.NewTxmORM(t, db, cfg)

		ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
		keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, nextNonce)

		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState}, &testCheckerFactory{})

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, borm, firstNonce, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		})).Return(errors.New("known transaction: a1313bd99a81fb4d8ad1d2e90b67c6b3fa77545c990d6251444b83b70b6f8980")).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check it was saved correctly with its attempt
		etx, err := borm.FindEthTxWithAttempts(inProgressEthTx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.InitialBroadcastAt)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.EthTxAttempts, 1)

		ethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and broadcast and now the transaction has been confirmed", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		borm := cltest.NewTxmORM(t, db, cfg)

		ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
		keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, nextNonce)

		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState}, &testCheckerFactory{})

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, borm, firstNonce, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		})).Return(errors.New("nonce too low")).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check it was saved correctly with its attempt
		etx, err := borm.FindEthTxWithAttempts(inProgressEthTx.ID)
		require.NoError(t, err)

		require.NotNil(t, etx.BroadcastAt)
		assert.Equal(t, *etx.BroadcastAt, etx.CreatedAt)
		assert.NotNil(t, etx.InitialBroadcastAt)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.EthTxAttempts, 1)

		ethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and then failed to reach node for some reason and node is still down", func(t *testing.T) {
		failedToReachNodeError := context.DeadlineExceeded
		db := pgtest.NewSqlxDB(t)
		borm := cltest.NewTxmORM(t, db, cfg)

		ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
		keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, nextNonce)

		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState}, &testCheckerFactory{})

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, borm, firstNonce, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		})).Return(failedToReachNodeError).Once()

		// Do the thing
		err := eb.ProcessUnstartedEthTxs(context.Background(), keyState)
		require.Error(t, err)
		assert.Contains(t, err.Error(), failedToReachNodeError.Error())

		// Check it was left in the unfinished state
		etx, err := borm.FindEthTxWithAttempts(inProgressEthTx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.Nil(t, etx.InitialBroadcastAt)
		assert.Equal(t, nextNonce, *etx.Nonce)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.EthTxAttempts, 1)

		ethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and broadcast transaction then crashed and rebooted with a different configured gas price", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		borm := cltest.NewTxmORM(t, db, cfg)

		ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
		keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, nextNonce)

		cfg := cltest.NewTestGeneralConfig(t)
		// Configured gas price changed
		cfg.Overrides.GlobalEvmGasPriceDefault = big.NewInt(500000000000)
		evmcfg := evmtest.NewChainScopedConfig(t, cfg)

		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState}, &testCheckerFactory{})

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, borm, firstNonce, fromAddress)
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
		etx, err := borm.FindEthTxWithAttempts(inProgressEthTx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.InitialBroadcastAt)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt = etx.EthTxAttempts[0]
		s, err := attempt.GetSignedTx()
		require.NoError(t, err)
		assert.Equal(t, int64(342), s.GasPrice().Int64())
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, attempt.State)

		ethClient.AssertExpectations(t)
	})
}

func getLocalNextNonce(t *testing.T, q pg.Q, fromAddress gethCommon.Address) uint64 {
	n, err := txmgr.GetNextNonce(q, fromAddress, &cltest.FixtureChainID)
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

	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	q := pg.NewQ(db, logger.TestLogger(t), cfg)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

	eb := cltest.NewEthBroadcaster(t, db, ethClient, ethKeyStore, evmcfg, []ethkey.State{keyState}, &testCheckerFactory{})

	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS pipeline_runs_pipeline_spec_id_fkey DEFERRED`)))

	t.Run("if external wallet sent a transaction from the account and now the nonce is one higher than it should be and we got replacement underpriced then we assume a previous transaction of ours was the one that succeeded, and hand off to EthConfirmer", func(t *testing.T) {
		etx := txmgr.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          txmgr.EthTxUnstarted,
		}
		require.NoError(t, borm.InsertEthTx(&etx))

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
		var etx1 txmgr.EthTx
		require.NoError(t, db.Get(&latestID, "SELECT max(id) FROM eth_txes"))
		etx1, err = borm.FindEthTxWithAttempts(latestID)
		require.NoError(t, err)
		require.NotNil(t, etx1.BroadcastAt)
		assert.NotEqual(t, etx1.CreatedAt, *etx1.BroadcastAt)
		assert.NotNil(t, etx1.InitialBroadcastAt)
		require.NotNil(t, etx1.Nonce)
		assert.Equal(t, int64(0), *etx1.Nonce)
		assert.False(t, etx1.Error.Valid)
		assert.Len(t, etx1.EthTxAttempts, 1)

		// Check that the local nonce was incremented by one
		var finalNextNonce int64
		finalNextNonce, err = txmgr.GetNextNonce(q, fromAddress, &cltest.FixtureChainID)
		require.NoError(t, err)
		require.NotNil(t, finalNextNonce)
		require.Equal(t, int64(1), finalNextNonce)
	})

	t.Run("geth Client returns an error in the fatal errors category", func(t *testing.T) {
		fatalErrorExample := "exceeds block gas limit"
		localNextNonce := getLocalNextNonce(t, q, fromAddress)

		t.Run("without callback", func(t *testing.T) {
			etx := txmgr.EthTx{
				FromAddress:    fromAddress,
				ToAddress:      toAddress,
				EncodedPayload: encodedPayload,
				Value:          value,
				GasLimit:       gasLimit,
				State:          txmgr.EthTxUnstarted,
			}
			require.NoError(t, borm.InsertEthTx(&etx))

			ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				return tx.Nonce() == localNextNonce
			})).Return(errors.New(fatalErrorExample)).Once()

			require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

			// Check it was saved correctly with its attempt
			etx, err = borm.FindEthTxWithAttempts(etx.ID)
			require.NoError(t, err)

			assert.Nil(t, etx.BroadcastAt)
			assert.Nil(t, etx.InitialBroadcastAt)
			require.Nil(t, etx.Nonce)
			assert.True(t, etx.Error.Valid)
			assert.Contains(t, etx.Error.String, "exceeds block gas limit")
			assert.Len(t, etx.EthTxAttempts, 0)

			// Check that the key had its nonce reset
			var state ethkey.State
			require.NoError(t, db.Get(&state, `SELECT * FROM eth_key_states`))
			// Saved NextNonce must be the same as before because this transaction
			// was not accepted by the eth node and never can be
			require.NotNil(t, state.NextNonce)
			require.Equal(t, int64(localNextNonce), state.NextNonce)
		})

		t.Run("with callback", func(t *testing.T) {
			run := cltest.MustInsertPipelineRun(t, db)
			tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)
			etx := txmgr.EthTx{
				FromAddress:       fromAddress,
				ToAddress:         toAddress,
				EncodedPayload:    encodedPayload,
				Value:             value,
				GasLimit:          gasLimit,
				State:             txmgr.EthTxUnstarted,
				PipelineTaskRunID: uuid.NullUUID{UUID: tr.ID, Valid: true},
			}

			t.Run("with erroring callback bails out", func(t *testing.T) {
				require.NoError(t, borm.InsertEthTx(&etx))
				fn := func(id uuid.UUID, result interface{}, err error) error {
					return errors.New("something exploded in the callback")
				}

				txmgr.SetResumeCallbackOnEthBroadcaster(fn, eb)

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

				txmgr.SetResumeCallbackOnEthBroadcaster(fn, eb)

				ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
					return tx.Nonce() == localNextNonce
				})).Return(errors.New(fatalErrorExample)).Once()

				require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))
			})

		})

		ethClient.AssertExpectations(t)
	})

	txmgr.SetResumeCallbackOnEthBroadcaster(nil, eb)

	t.Run("geth Client fails with error indicating that the transaction was too expensive", func(t *testing.T) {
		tooExpensiveError := "tx fee (1.10 ether) exceeds the configured cap (1.00 ether)"
		localNextNonce := getLocalNextNonce(t, q, fromAddress)

		etx := txmgr.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          txmgr.EthTxUnstarted,
		}
		require.NoError(t, borm.InsertEthTx(&etx))

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New(tooExpensiveError)).Once()

		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check it was saved with no attempt and a fatal error
		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.Nil(t, etx.InitialBroadcastAt)
		require.Nil(t, etx.Nonce)
		assert.True(t, etx.Error.Valid)
		assert.Contains(t, etx.Error.String, "tx fee (1.10 ether) exceeds the configured cap (1.00 ether)")
		assert.Len(t, etx.EthTxAttempts, 0)

		// Check that the key had its nonce reset
		var state ethkey.State
		require.NoError(t, db.Get(&state, `SELECT * FROM eth_key_states`))
		// Saved NextNonce must be the same as before because this transaction
		// was not accepted by the eth node and never can be
		require.NotNil(t, state.NextNonce)
		require.Equal(t, int64(localNextNonce), state.NextNonce)

		ethClient.AssertExpectations(t)
	})

	t.Run("eth Client call fails with an unexpected random error", func(t *testing.T) {
		retryableErrorExample := "geth shit the bed again"
		localNextNonce := getLocalNextNonce(t, q, fromAddress)

		etx := txmgr.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          txmgr.EthTxUnstarted,
		}
		require.NoError(t, borm.InsertEthTx(&etx))

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New(retryableErrorExample)).Once()

		// Do the thing
		err = eb.ProcessUnstartedEthTxs(context.Background(), keyState)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("error while sending transaction %v: %s", etx.ID, retryableErrorExample))

		// Check it was saved correctly with its attempt
		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.Nil(t, etx.InitialBroadcastAt)
		require.NotNil(t, etx.Nonce)
		assert.False(t, etx.Error.Valid)
		assert.Equal(t, txmgr.EthTxInProgress, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt := etx.EthTxAttempts[0]
		assert.Equal(t, txmgr.EthTxAttemptInProgress, attempt.State)

		ethClient.AssertExpectations(t)

		// Now on the second run, it is successful
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(nil).Once()

		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check it was saved correctly with its attempt
		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.InitialBroadcastAt)
		require.NotNil(t, etx.Nonce)
		assert.False(t, etx.Error.Valid)
		assert.Equal(t, txmgr.EthTxUnconfirmed, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt = etx.EthTxAttempts[0]
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, attempt.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("eth node returns underpriced transaction", func(t *testing.T) {
		// This happens if a transaction's gas price is below the minimum
		// configured for the transaction pool.
		// This is a configuration error by the node operator, since it means they set the base gas level too low.
		underpricedError := "transaction underpriced"
		localNextNonce := getLocalNextNonce(t, q, fromAddress)

		etx := txmgr.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          txmgr.EthTxUnstarted,
		}
		require.NoError(t, borm.InsertEthTx(&etx))

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
		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.InitialBroadcastAt)
		require.NotNil(t, etx.Nonce)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt := etx.EthTxAttempts[0]
		assert.Equal(t, big.NewInt(30000000000).String(), attempt.GasPrice.String())
	})

	etxUnfinished := txmgr.EthTx{
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: encodedPayload,
		Value:          value,
		GasLimit:       gasLimit,
		State:          txmgr.EthTxUnstarted,
	}
	require.NoError(t, borm.InsertEthTx(&etxUnfinished))

	t.Run("failed to reach node for some reason", func(t *testing.T) {
		failedToReachNodeError := context.DeadlineExceeded
		localNextNonce := getLocalNextNonce(t, q, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(failedToReachNodeError).Once()

		// Do the thing
		err = eb.ProcessUnstartedEthTxs(context.Background(), keyState)
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("error while sending transaction %v: context deadline exceeded", etxUnfinished.ID))

		// Check it was left in the unfinished state
		etx, err := borm.FindEthTxWithAttempts(etxUnfinished.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.Nil(t, etx.InitialBroadcastAt)
		assert.NotNil(t, etx.Nonce)
		assert.False(t, etx.Error.Valid)
		assert.Equal(t, txmgr.EthTxInProgress, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		assert.Equal(t, txmgr.EthTxAttemptInProgress, etx.EthTxAttempts[0].State)

		ethClient.AssertExpectations(t)
	})

	t.Run("eth node returns temporarily underpriced transaction", func(t *testing.T) {
		// This happens if parity is rejecting transactions that are not priced high enough to even get into the mempool at all
		// It should pretend it was accepted into the mempool and hand off to ethConfirmer to bump gas as normal
		temporarilyUnderpricedError := "There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee."
		localNextNonce := getLocalNextNonce(t, q, fromAddress)

		// Re-use the previously unfinished transaction, no need to insert new

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New(temporarilyUnderpricedError)).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(context.Background(), keyState))

		// Check it was saved correctly with its attempt
		etx, err := borm.FindEthTxWithAttempts(etxUnfinished.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.InitialBroadcastAt)
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
		localNextNonce := getLocalNextNonce(t, q, fromAddress)
		// In this scenario the node operator REALLY fucked up and set the bump
		// to zero (even though that should not be possible due to config
		// validation)
		oldGasBumpWei := evmcfg.EvmGasBumpWei()
		oldGasBumpPercent := evmcfg.EvmGasBumpPercent()
		cfg.Overrides.GlobalEvmGasBumpWei = big.NewInt(0)
		cfg.Overrides.GlobalEvmGasBumpPercent = null.IntFrom(0)
		defer func() {
			cfg.Overrides.GlobalEvmGasBumpWei = oldGasBumpWei
			cfg.Overrides.GlobalEvmGasBumpPercent = null.IntFrom(int64(oldGasBumpPercent))
		}()

		etx := txmgr.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          txmgr.EthTxUnstarted,
		}
		require.NoError(t, borm.InsertEthTx(&etx))

		// First was underpriced
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasPrice().Cmp(evmcfg.EvmGasPriceDefault()) == 0
		})).Return(errors.New(underpricedError)).Once()

		// Do the thing
		err := eb.ProcessUnstartedEthTxs(context.Background(), keyState)
		require.Error(t, err)
		require.Contains(t, err.Error(), "bumped gas price of 20000000000 is equal to original gas price of 20000000000. ACTION REQUIRED: This is a configuration error, you must increase either ETH_GAS_BUMP_PERCENT or ETH_GAS_BUMP_WEI")

		// TEARDOWN: Clear out the unsent tx before the next test
		pgtest.MustExec(t, db, `DELETE FROM eth_txes WHERE nonce = $1`, localNextNonce)

		ethClient.AssertExpectations(t)
	})

	t.Run("eth tx is left in progress if eth node returns insufficient eth", func(t *testing.T) {
		insufficientEthError := "insufficient funds for transfer"
		localNextNonce := getLocalNextNonce(t, q, fromAddress)
		etx := txmgr.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          txmgr.EthTxUnstarted,
		}
		require.NoError(t, borm.InsertEthTx(&etx))

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New(insufficientEthError)).Once()

		err := eb.ProcessUnstartedEthTxs(context.Background(), keyState)
		require.EqualError(t, err, "processUnstartedEthTxs failed: insufficient funds for transfer")

		// Check it was saved correctly with its attempt
		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.Nil(t, etx.InitialBroadcastAt)
		require.NotNil(t, etx.Nonce)
		assert.False(t, etx.Error.Valid)
		assert.Equal(t, txmgr.EthTxInProgress, etx.State)
		require.Len(t, etx.EthTxAttempts, 1)
		attempt := etx.EthTxAttempts[0]
		assert.Equal(t, txmgr.EthTxAttemptInProgress, attempt.State)
		assert.Nil(t, attempt.BroadcastBeforeBlockNum)

		ethClient.AssertExpectations(t)
	})

	pgtest.MustExec(t, db, `DELETE FROM eth_txes`)
	cfg.Overrides.GlobalEvmEIP1559DynamicFees = null.BoolFrom(true)

	t.Run("eth node returns underpriced transaction and bumping gas doesn't increase it in EIP-1559 mode", func(t *testing.T) {
		// This happens if a transaction's gas price is below the minimum
		// configured for the transaction pool.
		// This is a configuration error by the node operator, since it means they set the base gas level too low.

		// In this scenario the node operator REALLY fucked up and set the bump
		// to zero (even though that should not be possible due to config
		// validation)
		oldGasBumpWei := evmcfg.EvmGasBumpWei()
		oldGasBumpPercent := evmcfg.EvmGasBumpPercent()
		cfg.Overrides.GlobalEvmGasBumpWei = big.NewInt(0)
		cfg.Overrides.GlobalEvmGasBumpPercent = null.IntFrom(0)
		defer func() {
			cfg.Overrides.GlobalEvmGasBumpWei = oldGasBumpWei
			cfg.Overrides.GlobalEvmGasBumpPercent = null.IntFrom(int64(oldGasBumpPercent))
		}()

		etx := txmgr.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          txmgr.EthTxUnstarted,
		}
		require.NoError(t, borm.InsertEthTx(&etx))

		underpricedError := "transaction underpriced"
		localNextNonce := getLocalNextNonce(t, q, fromAddress)
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasTipCap().Cmp(big.NewInt(1)) == 0
		})).Return(errors.New(underpricedError)).Once()

		// Check gas tip cap verification
		err := eb.ProcessUnstartedEthTxs(context.Background(), keyState)
		require.Error(t, err)
		require.Contains(t, err.Error(), "bumped gas tip cap of 1 is less than or equal to original gas tip cap of 1")

		pgtest.MustExec(t, db, `DELETE FROM eth_txes`)
	})

	t.Run("eth node returns underpriced transaction in EIP-1559 mode, bumps until inclusion", func(t *testing.T) {
		// This happens if a transaction's gas price is below the minimum
		// configured for the transaction pool.
		// This is a configuration error by the node operator, since it means they set the base gas level too low.
		underpricedError := "transaction underpriced"
		localNextNonce := getLocalNextNonce(t, q, fromAddress)

		etx := txmgr.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          txmgr.EthTxUnstarted,
		}
		require.NoError(t, borm.InsertEthTx(&etx))

		// Check gas tip cap verification
		cfg.Overrides.GlobalEvmGasTipCapDefault = big.NewInt(0)
		err := eb.ProcessUnstartedEthTxs(context.Background(), keyState)
		require.Error(t, err)
		require.Contains(t, err.Error(), "specified gas tip cap of 0 is below min configured gas tip of 1 for key")

		gasTipCapDefault := big.NewInt(42)
		cfg.Overrides.GlobalEvmGasTipCapDefault = gasTipCapDefault
		// Second was underpriced but above minimum
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasTipCap().Cmp(gasTipCapDefault) == 0
		})).Return(errors.New(underpricedError)).Once()
		// Resend at the bumped price
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasTipCap().Cmp(big.NewInt(0).Add(gasTipCapDefault, evmcfg.EvmGasBumpWei())) == 0
		})).Return(errors.New(underpricedError)).Once()
		// Final bump succeeds
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasTipCap().Cmp(big.NewInt(0).Add(gasTipCapDefault, big.NewInt(0).Mul(evmcfg.EvmGasBumpWei(), big.NewInt(2)))) == 0
		})).Return(nil).Once()

		err = eb.ProcessUnstartedEthTxs(context.Background(), keyState)
		require.NoError(t, err)

		// TEARDOWN: Clear out the unsent tx before the next test
		pgtest.MustExec(t, db, `DELETE FROM eth_txes WHERE nonce = $1`, localNextNonce)

		ethClient.AssertExpectations(t)
	})

}

func TestEthBroadcaster_ProcessUnstartedEthTxs_KeystoreErrors(t *testing.T) {
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	value := assets.NewEthValue(142)
	gasLimit := uint64(242)
	encodedPayload := []byte{0, 1}
	localNonce := 0

	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)

	realKeystore := cltest.NewKeyStore(t, db, cfg)
	keyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, realKeystore.Eth())

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

	kst := new(ksmocks.Eth)
	eb := cltest.NewEthBroadcaster(t, db, ethClient, kst, evmcfg, []ethkey.State{keyState}, &testCheckerFactory{})

	t.Run("tx signing fails", func(t *testing.T) {
		etx := txmgr.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          txmgr.EthTxUnstarted,
		}
		require.NoError(t, borm.InsertEthTx(&etx))

		tx := *gethTypes.NewTx(&gethTypes.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.AnythingOfType("*types.Transaction"),
			mock.MatchedBy(func(chainID *big.Int) bool {
				return chainID.Cmp(evmcfg.ChainID()) == 0
			})).Return(&tx, errors.New("could not sign transaction")).Once()

		// Do the thing
		err := eb.ProcessUnstartedEthTxs(context.Background(), keyState)
		require.Error(t, err)
		require.Contains(t, err.Error(), "could not sign transaction")

		// Check that the transaction is left in unstarted state
		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgr.EthTxUnstarted, etx.State)
		assert.Len(t, etx.EthTxAttempts, 0)

		// Check that the key did not have its nonce incremented
		var keyState ethkey.State
		require.NoError(t, db.Get(&keyState, `SELECT * FROM eth_key_states`))
		require.NotNil(t, keyState.NextNonce)
		require.Equal(t, int64(localNonce), keyState.NextNonce)

		kst.AssertExpectations(t)
	})

	// Should have done nothing
	ethClient.AssertExpectations(t)
}

func TestEthBroadcaster_GetNextNonce(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	keyState, _ := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	q := pg.NewQ(db, logger.TestLogger(t), cfg)
	nonce, err := txmgr.GetNextNonce(q, keyState.Address.Address(), &cltest.FixtureChainID)
	assert.NoError(t, err)
	require.NotNil(t, nonce)
	assert.Equal(t, int64(0), nonce)
}

func TestEthBroadcaster_IncrementNextNonce(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	keyState, _ := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	// Cannot increment if supplied nonce doesn't match existing
	require.Error(t, txmgr.IncrementNextNonce(db, keyState.Address.Address(), &cltest.FixtureChainID, int64(42)))

	require.NoError(t, txmgr.IncrementNextNonce(db, keyState.Address.Address(), &cltest.FixtureChainID, int64(0)))

	// Nonce bumped to 1
	require.NoError(t, db.Get(&keyState, `SELECT * FROM eth_key_states LIMIT 1`))
	require.NotNil(t, keyState.NextNonce)
	require.Equal(t, int64(1), keyState.NextNonce)
}

func TestEthBroadcaster_Trigger(t *testing.T) {
	t.Parallel()

	// Simple sanity check to make sure it doesn't block
	db := pgtest.NewSqlxDB(t)

	cfg := cltest.NewTestGeneralConfig(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	eb := cltest.NewEthBroadcaster(t, db, cltest.NewEthClientMockWithDefaultChain(t), ethKeyStore, evmcfg, []ethkey.State{}, &testCheckerFactory{})

	eb.Trigger(testutils.NewAddress())
	eb.Trigger(testutils.NewAddress())
	eb.Trigger(testutils.NewAddress())
}

func TestEthBroadcaster_EthTxInsertEventCausesTriggerToFire(t *testing.T) {
	// NOTE: Testing triggers requires committing transactions and does not work with transactional tests
	cfg, db := heavyweight.FullTestDB(t, "eth_tx_triggers")
	borm := cltest.NewTxmORM(t, db, cfg)

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	eventBroadcaster := cltest.NewEventBroadcaster(t, evmcfg.DatabaseURL())
	require.NoError(t, eventBroadcaster.Start(testutils.Context(t)))
	t.Cleanup(func() { require.NoError(t, eventBroadcaster.Close()) })

	ethTxInsertListener, err := eventBroadcaster.Subscribe(pg.ChannelInsertOnEthTx, "")
	require.NoError(t, err)

	// Give it some time to start listening
	time.Sleep(100 * time.Millisecond)

	mustInsertUnstartedEthTx(t, borm, fromAddress)
	gomega.NewWithT(t).Eventually(ethTxInsertListener.Events()).Should(gomega.Receive())
}

func checkerToJson(t *testing.T, checker txmgr.TransmitCheckerSpec) *datatypes.JSON {
	b, err := json.Marshal(checker)
	require.NoError(t, err)
	j := datatypes.JSON(b)
	return &j
}

type testCheckerFactory struct {
	err error
}

func (t *testCheckerFactory) BuildChecker(spec txmgr.TransmitCheckerSpec) (txmgr.TransmitChecker, error) {
	return &testChecker{t.err}, nil
}

type testChecker struct {
	err error
}

func (t *testChecker) Check(
	_ context.Context,
	_ logger.Logger,
	_ txmgr.EthTx,
	_ txmgr.EthTxAttempt,
) error {
	return t.err
}
