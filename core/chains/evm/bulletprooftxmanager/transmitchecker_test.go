package bulletprooftxmanager_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/bulletprooftxmanager"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg/datatypes"
)

func TestFactory(t *testing.T) {
	client, _, _ := cltest.NewEthMocksWithDefaultChain(t)
	factory := &bulletprooftxmanager.CheckerFactory{Client: client}

	t.Run("no checker", func(t *testing.T) {
		c, err := factory.BuildChecker(bulletprooftxmanager.TransmitCheckerSpec{})
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.NoChecker{}, c)
	})

	t.Run("vrf checker", func(t *testing.T) {
		c, err := factory.BuildChecker(bulletprooftxmanager.TransmitCheckerSpec{
			CheckerType:           bulletprooftxmanager.TransmitCheckerTypeVRFV2,
			VRFCoordinatorAddress: cltest.NewAddress(),
		})
		require.NoError(t, err)
		require.IsType(t, bulletprooftxmanager.VRFV2Checker{}, c)
	})

	t.Run("simulate checker", func(t *testing.T) {
		c, err := factory.BuildChecker(bulletprooftxmanager.TransmitCheckerSpec{
			CheckerType: bulletprooftxmanager.TransmitCheckerTypeSimulate,
		})
		require.NoError(t, err)
		require.Equal(t, &bulletprooftxmanager.SimulateChecker{Client: client}, c)
	})

	t.Run("invalid checker type", func(t *testing.T) {
		_, err := factory.BuildChecker(bulletprooftxmanager.TransmitCheckerSpec{
			CheckerType: "invalid",
		})
		require.EqualError(t, err, "unrecognized checker type: invalid")
	})
}

func TestTransmitCheckers(t *testing.T) {
	client := cltest.NewEthClientMockWithDefaultChain(t)
	log := logger.TestLogger(t)
	ctx := context.Background()

	t.Run("no checker", func(t *testing.T) {
		checker := bulletprooftxmanager.NoChecker{}
		require.NoError(t, checker.Check(ctx, log, bulletprooftxmanager.EthTx{}, bulletprooftxmanager.EthTxAttempt{}))
	})

	t.Run("simulate", func(t *testing.T) {
		checker := bulletprooftxmanager.SimulateChecker{Client: client}

		tx := bulletprooftxmanager.EthTx{
			FromAddress:    common.HexToAddress("0xfe0629509E6CB8dfa7a99214ae58Ceb465d5b5A9"),
			ToAddress:      common.HexToAddress("0xff0Aac13eab788cb9a2D662D3FB661Aa5f58FA21"),
			EncodedPayload: []byte{42, 0, 0},
			Value:          assets.NewEthValue(642),
			GasLimit:       1e9,
			CreatedAt:      time.Unix(0, 0),
			State:          bulletprooftxmanager.EthTxUnstarted,
		}
		attempt := bulletprooftxmanager.EthTxAttempt{
			EthTx:     tx,
			Hash:      common.Hash{},
			CreatedAt: tx.CreatedAt,
			State:     bulletprooftxmanager.EthTxAttemptInProgress,
		}

		t.Run("success", func(t *testing.T) {
			client.On("CallContext", mock.Anything,
				mock.AnythingOfType("*hexutil.Bytes"), "eth_call",
				mock.MatchedBy(func(callarg map[string]interface{}) bool {
					return fmt.Sprintf("%s", callarg["value"]) == "0x282" // 642
				}), "latest").Return(nil).Once()

			require.NoError(t, checker.Check(ctx, log, tx, attempt))
			client.AssertExpectations(t)
		})

		t.Run("revert", func(t *testing.T) {
			jerr := evmclient.JsonError{
				Code:    42,
				Message: "oh no, it reverted",
				Data:    []byte{42, 166, 34},
			}
			client.On("CallContext", mock.Anything,
				mock.AnythingOfType("*hexutil.Bytes"), "eth_call",
				mock.MatchedBy(func(callarg map[string]interface{}) bool {
					return fmt.Sprintf("%s", callarg["value"]) == "0x282" // 642
				}), "latest").Return(&jerr).Once()

			err := checker.Check(ctx, log, tx, attempt)
			expErrMsg := "transaction reverted during simulation: json-rpc error { Code = 42, Message = 'oh no, it reverted', Data = 'KqYi' }"
			require.EqualError(t, err, expErrMsg)
			client.AssertExpectations(t)
		})

		t.Run("non revert error", func(t *testing.T) {
			client.On("CallContext", mock.Anything,
				mock.AnythingOfType("*hexutil.Bytes"), "eth_call",
				mock.MatchedBy(func(callarg map[string]interface{}) bool {
					return fmt.Sprintf("%s", callarg["value"]) == "0x282" // 642
				}), "latest").Return(errors.New("error!")).Once()

			// Non-revert errors are logged but should not prevent transmission, and do not need
			// to be passed to the caller
			require.NoError(t, checker.Check(ctx, log, tx, attempt))
			client.AssertExpectations(t)
		})
	})

	t.Run("VRF V2", func(t *testing.T) {

		newTx := func(t *testing.T, vrfReqID *big.Int) (bulletprooftxmanager.EthTx, bulletprooftxmanager.EthTxAttempt) {
			meta := bulletprooftxmanager.EthTxMeta{
				RequestID: common.BytesToHash(vrfReqID.Bytes()),
				MaxLink:   "1000000000000000000", // 1 LINK
				SubID:     2,
			}

			b, err := json.Marshal(meta)
			require.NoError(t, err)
			metaJson := datatypes.JSON(b)

			tx := bulletprooftxmanager.EthTx{
				FromAddress:    common.HexToAddress("0xfe0629509E6CB8dfa7a99214ae58Ceb465d5b5A9"),
				ToAddress:      common.HexToAddress("0xff0Aac13eab788cb9a2D662D3FB661Aa5f58FA21"),
				EncodedPayload: []byte{42, 0, 0},
				Value:          assets.NewEthValue(642),
				GasLimit:       1e9,
				CreatedAt:      time.Unix(0, 0),
				State:          bulletprooftxmanager.EthTxUnstarted,
				Meta:           &metaJson,
			}
			return tx, bulletprooftxmanager.EthTxAttempt{
				EthTx:     tx,
				Hash:      common.Hash{},
				CreatedAt: tx.CreatedAt,
				State:     bulletprooftxmanager.EthTxAttemptInProgress,
			}
		}

		checker := bulletprooftxmanager.VRFV2Checker{GetCommitment: func(_ *bind.CallOpts, requestID *big.Int) ([32]byte, error) {
			fmt.Printf("requestID: %v\n", requestID.String())
			if requestID.String() == "1" {
				// Request 1 is already fulfilled
				return [32]byte{}, nil
			} else if requestID.String() == "2" {
				// Request 2 errors
				return [32]byte{}, errors.New("error getting commitment")
			} else {
				// All other requests are unfulfilled
				return [32]byte{1}, nil
			}
		}}

		t.Run("already fulfilled", func(t *testing.T) {
			tx, attempt := newTx(t, big.NewInt(1))
			err := checker.Check(ctx, log, tx, attempt)
			require.Error(t, err, "request already fulfilled")
		})

		t.Run("not fulfilled", func(t *testing.T) {
			tx, attempt := newTx(t, big.NewInt(3))
			require.NoError(t, checker.Check(ctx, log, tx, attempt))
		})

		t.Run("error checking fulfillment, should transmit", func(t *testing.T) {
			tx, attempt := newTx(t, big.NewInt(2))
			require.NoError(t, checker.Check(ctx, log, tx, attempt))
		})
	})
}
