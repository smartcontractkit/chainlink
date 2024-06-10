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
	"github.com/ethereum/go-ethereum/rpc"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	v1 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
)

func TestFactory(t *testing.T) {
	client := testutils.NewEthClientMockWithDefaultChain(t)
	factory := &txmgr.CheckerFactory{Client: client}

	t.Run("no checker", func(t *testing.T) {
		c, err := factory.BuildChecker(txmgr.TransmitCheckerSpec{})
		require.NoError(t, err)
		require.Equal(t, txmgr.NoChecker, c)
	})

	t.Run("vrf v1 checker", func(t *testing.T) {
		c, err := factory.BuildChecker(txmgr.TransmitCheckerSpec{
			CheckerType:           txmgr.TransmitCheckerTypeVRFV1,
			VRFCoordinatorAddress: testutils.NewAddressPtr(),
		})
		require.NoError(t, err)
		require.IsType(t, &txmgr.VRFV1Checker{}, c)
	})

	t.Run("vrf v2 checker", func(t *testing.T) {
		c, err := factory.BuildChecker(txmgr.TransmitCheckerSpec{
			CheckerType:           txmgr.TransmitCheckerTypeVRFV2,
			VRFCoordinatorAddress: testutils.NewAddressPtr(),
			VRFRequestBlockNumber: big.NewInt(1),
		})
		require.NoError(t, err)
		require.IsType(t, &txmgr.VRFV2Checker{}, c)

		// request block number not provided should error out.
		c, err = factory.BuildChecker(txmgr.TransmitCheckerSpec{
			CheckerType:           txmgr.TransmitCheckerTypeVRFV2,
			VRFCoordinatorAddress: testutils.NewAddressPtr(),
		})
		require.Error(t, err)
		require.Nil(t, c)
	})

	t.Run("vrf v2 plus checker", func(t *testing.T) {
		c, err := factory.BuildChecker(txmgr.TransmitCheckerSpec{
			CheckerType:           txmgr.TransmitCheckerTypeVRFV2Plus,
			VRFCoordinatorAddress: testutils.NewAddressPtr(),
			VRFRequestBlockNumber: big.NewInt(1),
		})
		require.NoError(t, err)
		require.IsType(t, &txmgr.VRFV2Checker{}, c)

		// request block number not provided should error out.
		c, err = factory.BuildChecker(txmgr.TransmitCheckerSpec{
			CheckerType:           txmgr.TransmitCheckerTypeVRFV2Plus,
			VRFCoordinatorAddress: testutils.NewAddressPtr(),
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
	client := testutils.NewEthClientMockWithDefaultChain(t)
	log := logger.Sugared(logger.Test(t))
	ctx := tests.Context(t)

	t.Run("no checker", func(t *testing.T) {
		checker := txmgr.NoChecker
		require.NoError(t, checker.Check(ctx, log, txmgr.Tx{}, txmgr.TxAttempt{}))
	})

	t.Run("simulate", func(t *testing.T) {
		checker := txmgr.SimulateChecker{Client: client}

		tx := txmgr.Tx{
			FromAddress:    common.HexToAddress("0xfe0629509E6CB8dfa7a99214ae58Ceb465d5b5A9"),
			ToAddress:      common.HexToAddress("0xff0Aac13eab788cb9a2D662D3FB661Aa5f58FA21"),
			EncodedPayload: []byte{42, 0, 0},
			Value:          big.Int(assets.NewEthValue(642)),
			FeeLimit:       1e9,
			CreatedAt:      time.Unix(0, 0),
			State:          txmgrcommon.TxUnstarted,
		}
		attempt := txmgr.TxAttempt{
			Tx:        tx,
			Hash:      common.Hash{},
			CreatedAt: tx.CreatedAt,
			State:     txmgrtypes.TxAttemptInProgress,
		}

		t.Run("success", func(t *testing.T) {
			client.On("CallContext", mock.Anything,
				mock.AnythingOfType("*hexutil.Bytes"), "eth_call",
				mock.MatchedBy(func(callarg map[string]interface{}) bool {
					return fmt.Sprintf("%s", callarg["value"]) == "0x282" // 642
				}), "latest").Return(nil).Once()

			require.NoError(t, checker.Check(ctx, log, tx, attempt))
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
		})

		t.Run("non revert error", func(t *testing.T) {
			client.On("CallContext", mock.Anything,
				mock.AnythingOfType("*hexutil.Bytes"), "eth_call",
				mock.MatchedBy(func(callarg map[string]interface{}) bool {
					return fmt.Sprintf("%s", callarg["value"]) == "0x282" // 642
				}), "latest").Return(pkgerrors.New("error")).Once()

			// Non-revert errors are logged but should not prevent transmission, and do not need
			// to be passed to the caller
			require.NoError(t, checker.Check(ctx, log, tx, attempt))
		})
	})

	t.Run("VRF V1", func(t *testing.T) {
		testDefaultSubID := uint64(2)
		testDefaultMaxLink := "1000000000000000000"

		txRequest := func(t *testing.T, vrfReqID [32]byte, nilTxHash bool) (txmgr.Tx, txmgr.TxAttempt) {
			h := common.BytesToHash(vrfReqID[:])
			txHash := common.Hash{}
			meta := txmgr.TxMeta{
				RequestID:     &h,
				MaxLink:       &testDefaultMaxLink, // 1 LINK
				SubID:         &testDefaultSubID,
				RequestTxHash: &txHash,
			}

			if nilTxHash {
				meta.RequestTxHash = nil
			}

			b, err := json.Marshal(meta)
			require.NoError(t, err)
			metaJson := sqlutil.JSON(b)

			tx := txmgr.Tx{
				FromAddress:    common.HexToAddress("0xfe0629509E6CB8dfa7a99214ae58Ceb465d5b5A9"),
				ToAddress:      common.HexToAddress("0xff0Aac13eab788cb9a2D662D3FB661Aa5f58FA21"),
				EncodedPayload: []byte{42, 0, 0},
				Value:          big.Int(assets.NewEthValue(642)),
				FeeLimit:       1e9,
				CreatedAt:      time.Unix(0, 0),
				State:          txmgrcommon.TxUnstarted,
				Meta:           &metaJson,
			}
			return tx, txmgr.TxAttempt{
				Tx:        tx,
				Hash:      common.Hash{},
				CreatedAt: tx.CreatedAt,
				State:     txmgrtypes.TxAttemptInProgress,
			}
		}

		r1 := [32]byte{1}
		r2 := [32]byte{2}
		r3 := [32]byte{3}

		checker := txmgr.VRFV1Checker{
			Callbacks: func(opts *bind.CallOpts, reqID [32]byte) (v1.Callbacks, error) {
				if opts.BlockNumber.Cmp(big.NewInt(6)) != 0 {
					// Ensure correct logic is applied to get callbacks.
					return v1.Callbacks{}, pkgerrors.New("error getting callback")
				}
				if reqID == r1 {
					// Request 1 is already fulfilled
					return v1.Callbacks{
						SeedAndBlockNum: [32]byte{},
					}, nil
				} else if reqID == r2 {
					// Request 2 errors
					return v1.Callbacks{}, pkgerrors.New("error getting commitment")
				}
				return v1.Callbacks{
					SeedAndBlockNum: [32]byte{1},
				}, nil
			},
			Client: client,
		}

		mockBatch := client.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 && b[0].Method == "eth_getBlockByNumber" && b[1].Method == "eth_getTransactionReceipt"
		})).Return(nil).Run(func(args mock.Arguments) {
			batch := args.Get(1).([]rpc.BatchElem)

			// Return block 10 for eth_getBlockByNumber
			mostRecentHead := batch[0].Result.(*evmtypes.Head)
			mostRecentHead.Number = 10

			// Return block 6 for eth_getTransactionReceipt
			requestTransactionReceipt := batch[1].Result.(*types.Receipt)
			requestTransactionReceipt.BlockNumber = big.NewInt(6)
		})

		t.Run("already fulfilled", func(t *testing.T) {
			tx, attempt := txRequest(t, r1, false)
			err := checker.Check(ctx, log, tx, attempt)
			require.Error(t, err, "request already fulfilled")
		})

		t.Run("nil RequestTxHash", func(t *testing.T) {
			tx, attempt := txRequest(t, r1, true)
			err := checker.Check(ctx, log, tx, attempt)
			require.NoError(t, err)
		})

		t.Run("not fulfilled", func(t *testing.T) {
			tx, attempt := txRequest(t, r3, false)
			require.NoError(t, checker.Check(ctx, log, tx, attempt))
		})

		t.Run("error checking fulfillment, should transmit", func(t *testing.T) {
			tx, attempt := txRequest(t, r2, false)
			require.NoError(t, checker.Check(ctx, log, tx, attempt))
		})

		t.Run("failure fetching tx receipt and block head", func(t *testing.T) {
			tx, attempt := txRequest(t, r1, false)
			mockBatch.Return(pkgerrors.New("could not fetch"))
			err := checker.Check(ctx, log, tx, attempt)
			require.NoError(t, err)
		})
	})

	t.Run("VRF V2", func(t *testing.T) {
		testDefaultSubID := uint64(2)
		testDefaultMaxLink := "1000000000000000000"

		txRequest := func(t *testing.T, vrfReqID *big.Int) (txmgr.Tx, txmgr.TxAttempt) {
			h := common.BytesToHash(vrfReqID.Bytes())
			meta := txmgr.TxMeta{
				RequestID: &h,
				MaxLink:   &testDefaultMaxLink, // 1 LINK
				SubID:     &testDefaultSubID,
			}

			b, err := json.Marshal(meta)
			require.NoError(t, err)
			metaJson := sqlutil.JSON(b)

			tx := txmgr.Tx{
				FromAddress:    common.HexToAddress("0xfe0629509E6CB8dfa7a99214ae58Ceb465d5b5A9"),
				ToAddress:      common.HexToAddress("0xff0Aac13eab788cb9a2D662D3FB661Aa5f58FA21"),
				EncodedPayload: []byte{42, 0, 0},
				Value:          big.Int(assets.NewEthValue(642)),
				FeeLimit:       1e9,
				CreatedAt:      time.Unix(0, 0),
				State:          txmgrcommon.TxUnstarted,
				Meta:           &metaJson,
			}
			return tx, txmgr.TxAttempt{
				Tx:        tx,
				Hash:      common.Hash{},
				CreatedAt: tx.CreatedAt,
				State:     txmgrtypes.TxAttemptInProgress,
			}
		}

		checker := txmgr.VRFV2Checker{
			GetCommitment: func(_ *bind.CallOpts, requestID *big.Int) ([32]byte, error) {
				if requestID.String() == "1" {
					// Request 1 is already fulfilled
					return [32]byte{}, nil
				} else if requestID.String() == "2" {
					// Request 2 errors
					return [32]byte{}, pkgerrors.New("error getting commitment")
				}
				// All other requests are unfulfilled
				return [32]byte{1}, nil
			},
			HeadByNumber: func(ctx context.Context, n *big.Int) (*evmtypes.Head, error) {
				return &evmtypes.Head{
					Number: 1,
				}, nil
			},
			RequestBlockNumber: big.NewInt(1),
		}

		t.Run("already fulfilled", func(t *testing.T) {
			tx, attempt := txRequest(t, big.NewInt(1))
			err := checker.Check(ctx, log, tx, attempt)
			require.Error(t, err, "request already fulfilled")
		})

		t.Run("not fulfilled", func(t *testing.T) {
			tx, attempt := txRequest(t, big.NewInt(3))
			require.NoError(t, checker.Check(ctx, log, tx, attempt))
		})

		t.Run("error checking fulfillment, should transmit", func(t *testing.T) {
			tx, attempt := txRequest(t, big.NewInt(2))
			require.NoError(t, checker.Check(ctx, log, tx, attempt))
		})

		t.Run("can't get header", func(t *testing.T) {
			checker.HeadByNumber = func(ctx context.Context, n *big.Int) (*evmtypes.Head, error) {
				return nil, pkgerrors.New("can't get head")
			}
			tx, attempt := txRequest(t, big.NewInt(3))
			require.NoError(t, checker.Check(ctx, log, tx, attempt))
		})

		t.Run("nil request block number", func(t *testing.T) {
			checker.HeadByNumber = func(ctx context.Context, n *big.Int) (*evmtypes.Head, error) {
				return &evmtypes.Head{
					Number: 1,
				}, nil
			}
			checker.RequestBlockNumber = nil
			tx, attempt := txRequest(t, big.NewInt(4))
			require.NoError(t, checker.Check(ctx, log, tx, attempt))
		})
	})
}
