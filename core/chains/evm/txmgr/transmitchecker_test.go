package txmgr_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	v1 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg/datatypes"
)

func TestFactory(t *testing.T) {
	client := cltest.NewEthMocksWithDefaultChain(t)
	factory := &txmgr.CheckerFactory{Client: client}

	t.Run("no checker", func(t *testing.T) {
		c, err := factory.BuildChecker(txmgr.TransmitCheckerSpec{})
		require.NoError(t, err)
		require.Equal(t, txmgr.NoChecker, c)
	})

	t.Run("vrf v1 checker", func(t *testing.T) {
		c, err := factory.BuildChecker(txmgr.TransmitCheckerSpec{
			CheckerType:           txmgr.TransmitCheckerTypeVRFV1,
			VRFCoordinatorAddress: testutils.NewAddress(),
		})
		require.NoError(t, err)
		require.IsType(t, &txmgr.VRFV1Checker{}, c)
	})

	t.Run("vrf v2 checker", func(t *testing.T) {
		c, err := factory.BuildChecker(txmgr.TransmitCheckerSpec{
			CheckerType:           txmgr.TransmitCheckerTypeVRFV2,
			VRFCoordinatorAddress: testutils.NewAddress(),
			VRFRequestBlockNumber: big.NewInt(1),
		})
		require.NoError(t, err)
		require.IsType(t, &txmgr.VRFV2Checker{}, c)

		// request block number not provided should error out.
		c, err = factory.BuildChecker(txmgr.TransmitCheckerSpec{
			CheckerType:           txmgr.TransmitCheckerTypeVRFV2,
			VRFCoordinatorAddress: testutils.NewAddress(),
		})
		require.Error(t, err)
		require.Nil(t, c)
	})

	t.Run("simulate checker", func(t *testing.T) {
		c, err := factory.BuildChecker(txmgr.TransmitCheckerSpec{
			CheckerType: txmgr.TransmitCheckerTypeSimulate,
		})
		require.NoError(t, err)
		require.Equal(t, &txmgr.SimulateChecker{Client: client}, c)
	})

	t.Run("invalid checker type", func(t *testing.T) {
		_, err := factory.BuildChecker(txmgr.TransmitCheckerSpec{
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
		checker := txmgr.NoChecker
		require.NoError(t, checker.Check(ctx, log, txmgr.EthTx{}, txmgr.EthTxAttempt{}))
	})

	t.Run("simulate", func(t *testing.T) {
		checker := txmgr.SimulateChecker{Client: client}

		tx := txmgr.EthTx{
			FromAddress:    common.HexToAddress("0xfe0629509E6CB8dfa7a99214ae58Ceb465d5b5A9"),
			ToAddress:      common.HexToAddress("0xff0Aac13eab788cb9a2D662D3FB661Aa5f58FA21"),
			EncodedPayload: []byte{42, 0, 0},
			Value:          assets.NewEthValue(642),
			GasLimit:       1e9,
			CreatedAt:      time.Unix(0, 0),
			State:          txmgr.EthTxUnstarted,
		}
		attempt := txmgr.EthTxAttempt{
			EthTx:     tx,
			Hash:      common.Hash{},
			CreatedAt: tx.CreatedAt,
			State:     txmgr.EthTxAttemptInProgress,
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

	t.Run("VRF V1", func(t *testing.T) {
		testDefaultSubID := uint64(2)
		testDefaultMaxLink := "1000000000000000000"

		newTx := func(t *testing.T, vrfReqID [32]byte) (txmgr.EthTx, txmgr.EthTxAttempt) {
			meta := txmgr.EthTxMeta{
				RequestID: common.BytesToHash(vrfReqID[:]),
				MaxLink:   &testDefaultMaxLink, // 1 LINK
				SubID:     &testDefaultSubID,
			}

			b, err := json.Marshal(meta)
			require.NoError(t, err)
			metaJson := datatypes.JSON(b)

			tx := txmgr.EthTx{
				FromAddress:    common.HexToAddress("0xfe0629509E6CB8dfa7a99214ae58Ceb465d5b5A9"),
				ToAddress:      common.HexToAddress("0xff0Aac13eab788cb9a2D662D3FB661Aa5f58FA21"),
				EncodedPayload: []byte{42, 0, 0},
				Value:          assets.NewEthValue(642),
				GasLimit:       1e9,
				CreatedAt:      time.Unix(0, 0),
				State:          txmgr.EthTxUnstarted,
				Meta:           &metaJson,
			}
			return tx, txmgr.EthTxAttempt{
				EthTx:     tx,
				Hash:      common.Hash{},
				CreatedAt: tx.CreatedAt,
				State:     txmgr.EthTxAttemptInProgress,
			}
		}

		r1 := [32]byte{1}
		r2 := [32]byte{2}
		r3 := [32]byte{3}

		checker := txmgr.VRFV1Checker{Callbacks: func(opts *bind.CallOpts, reqID [32]byte) (v1.Callbacks, error) {
			if reqID == r1 {
				// Request 1 is already fulfilled
				return v1.Callbacks{
					SeedAndBlockNum: [32]byte{},
				}, nil
			} else if reqID == r2 {
				// Request 2 errors
				return v1.Callbacks{}, errors.New("error getting commitment")
			} else {
				return v1.Callbacks{
					SeedAndBlockNum: [32]byte{1},
				}, nil
			}
		}}

		t.Run("already fulfilled", func(t *testing.T) {
			tx, attempt := newTx(t, r1)
			err := checker.Check(ctx, log, tx, attempt)
			require.Error(t, err, "request already fulfilled")
		})

		t.Run("not fulfilled", func(t *testing.T) {
			tx, attempt := newTx(t, r3)
			require.NoError(t, checker.Check(ctx, log, tx, attempt))
		})

		t.Run("error checking fulfillment, should transmit", func(t *testing.T) {
			tx, attempt := newTx(t, r2)
			require.NoError(t, checker.Check(ctx, log, tx, attempt))
		})
	})

	t.Run("VRF V2", func(t *testing.T) {
		testDefaultSubID := uint64(2)
		testDefaultMaxLink := "1000000000000000000"

		newTx := func(t *testing.T, vrfReqID *big.Int) (txmgr.EthTx, txmgr.EthTxAttempt) {
			meta := txmgr.EthTxMeta{
				RequestID: common.BytesToHash(vrfReqID.Bytes()),
				MaxLink:   &testDefaultMaxLink, // 1 LINK
				SubID:     &testDefaultSubID,
			}

			b, err := json.Marshal(meta)
			require.NoError(t, err)
			metaJson := datatypes.JSON(b)

			tx := txmgr.EthTx{
				FromAddress:    common.HexToAddress("0xfe0629509E6CB8dfa7a99214ae58Ceb465d5b5A9"),
				ToAddress:      common.HexToAddress("0xff0Aac13eab788cb9a2D662D3FB661Aa5f58FA21"),
				EncodedPayload: []byte{42, 0, 0},
				Value:          assets.NewEthValue(642),
				GasLimit:       1e9,
				CreatedAt:      time.Unix(0, 0),
				State:          txmgr.EthTxUnstarted,
				Meta:           &metaJson,
			}
			return tx, txmgr.EthTxAttempt{
				EthTx:     tx,
				Hash:      common.Hash{},
				CreatedAt: tx.CreatedAt,
				State:     txmgr.EthTxAttemptInProgress,
			}
		}

		checker := txmgr.VRFV2Checker{
			GetCommitment: func(_ *bind.CallOpts, requestID *big.Int) ([32]byte, error) {
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
			},
			HeaderByNumber: func(ctx context.Context, n *big.Int) (*types.Header, error) {
				return &types.Header{
					Number: big.NewInt(1),
				}, nil
			},
			RequestBlockNumber: big.NewInt(1),
		}

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

		t.Run("can't get header", func(t *testing.T) {
			checker.HeaderByNumber = func(ctx context.Context, n *big.Int) (*types.Header, error) {
				return nil, errors.New("can't get head")
			}
			tx, attempt := newTx(t, big.NewInt(3))
			require.NoError(t, checker.Check(ctx, log, tx, attempt))
		})

		t.Run("nil request block number", func(t *testing.T) {
			checker.HeaderByNumber = func(ctx context.Context, n *big.Int) (*types.Header, error) {
				return &types.Header{
					Number: big.NewInt(1),
				}, nil
			}
			checker.RequestBlockNumber = nil
			tx, attempt := newTx(t, big.NewInt(4))
			require.NoError(t, checker.Check(ctx, log, tx, attempt))
		})
	})
}
