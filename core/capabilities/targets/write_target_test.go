package targets_test

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type testHarness struct {
	lggr          logger.Logger
	cw            *mocks.ContractWriter
	cr            *mocks.ContractValueGetter
	config        *values.Map
	validInputs   *values.Map
	validMetadata capabilities.RequestMetadata
	writeTarget   *targets.WriteTarget
	forwarderAddr string
	binding       types.BoundContract
	gasLimit      uint64
}

func setup(t *testing.T) testHarness {
	lggr := logger.TestLogger(t)

	cw := mocks.NewContractWriter(t)
	cr := mocks.NewContractValueGetter(t)

	forwarderA := testutils.NewAddress()
	forwarderAddr := forwarderA.Hex()

	var txGasLimit uint64 = 400_000
	writeTarget := targets.NewWriteTarget(lggr, "test-write-target@1.0.0", cr, cw, forwarderAddr, txGasLimit)
	require.NotNil(t, writeTarget)

	config, err := values.NewMap(map[string]any{
		"Address": forwarderAddr,
	})
	require.NoError(t, err)

	reportID := [2]byte{0x00, 0x01}
	var workflowName [10]byte
	copy(workflowName[:], []byte("name"))
	workflowOwnerString := "219BFD3D78fbb740c614432975CBE829E26C490e"
	workflowOwner := common.HexToAddress(workflowOwnerString)
	reportMetadata := targets.ReportV1Metadata{
		Version:             1,
		WorkflowExecutionID: [32]byte{},
		Timestamp:           0,
		DonID:               0,
		DonConfigVersion:    0,
		WorkflowCID:         [32]byte{},
		WorkflowName:        workflowName,
		WorkflowOwner:       workflowOwner,
		ReportID:            reportID,
	}

	reportMetadataBytes, err := reportMetadata.Encode()
	require.NoError(t, err)

	validInputs, err := values.NewMap(map[string]any{
		"signed_report": map[string]any{
			"report":     reportMetadataBytes,
			"signatures": [][]byte{},
			"context":    []byte{4, 5},
			"id":         reportID[:],
		},
	})
	require.NoError(t, err)

	validMetadata := capabilities.RequestMetadata{
		WorkflowID:          hex.EncodeToString(reportMetadata.WorkflowCID[:]),
		WorkflowOwner:       workflowOwnerString,
		WorkflowName:        hex.EncodeToString(reportMetadata.WorkflowName[:]),
		WorkflowExecutionID: hex.EncodeToString(reportMetadata.WorkflowExecutionID[:]),
	}

	binding := types.BoundContract{
		Address: forwarderAddr,
		Name:    "forwarder",
	}

	cr.On("Bind", mock.Anything, []types.BoundContract{binding}).Return(nil)

	return testHarness{
		lggr:          lggr,
		cw:            cw,
		cr:            cr,
		config:        config,
		validInputs:   validInputs,
		validMetadata: validMetadata,
		writeTarget:   writeTarget,
		forwarderAddr: forwarderAddr,
		binding:       binding,
		gasLimit:      txGasLimit - targets.ForwarderContractLogicGasCost,
	}
}
func TestWriteTarget(t *testing.T) {
	th := setup(t)
	th.cr.EXPECT().GetLatestValue(mock.Anything, th.binding.ReadIdentifier("getTransmissionInfo"), mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(_ context.Context, _ string, _ primitives.ConfidenceLevel, _, retVal any) {
		transmissionInfo := retVal.(*targets.TransmissionInfo)
		*transmissionInfo = targets.TransmissionInfo{
			GasLimit:        big.NewInt(0),
			InvalidReceiver: false,
			State:           0,
			Success:         false,
			TransmissionID:  [32]byte{},
			Transmitter:     common.HexToAddress("0x0"),
		}
	})
	t.Run("succeeds with valid report", func(t *testing.T) {
		req := capabilities.CapabilityRequest{
			Metadata: th.validMetadata,
			Config:   th.config,
			Inputs:   th.validInputs,
		}
		th.cw.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(types.Finalized, nil).Once()
		th.cw.On("SubmitTransaction", mock.Anything, "forwarder", "report", mock.Anything, mock.Anything, th.forwarderAddr, mock.Anything, mock.Anything).Return(nil).Once()

		ctx := testutils.Context(t)
		response, err2 := th.writeTarget.Execute(ctx, req)
		require.NoError(t, err2)
		require.NotNil(t, response)
	})

	t.Run("fails when ChainWriter's SubmitTransaction returns error", func(t *testing.T) {
		req := capabilities.CapabilityRequest{
			Metadata: th.validMetadata,
			Config:   th.config,
			Inputs:   th.validInputs,
		}
		th.cw.On("SubmitTransaction", mock.Anything, "forwarder", "report", mock.Anything, mock.Anything, th.forwarderAddr, mock.Anything, mock.Anything).Return(errors.New("writer error"))

		ctx := testutils.Context(t)
		_, err := th.writeTarget.Execute(ctx, req)
		require.Error(t, err)
	})

	t.Run("passes gas limit set on config to the chain writer", func(t *testing.T) {
		configGasLimit, err2 := values.NewMap(map[string]any{
			"Address":  th.forwarderAddr,
			"GasLimit": 500000,
		})
		require.NoError(t, err2)
		req := capabilities.CapabilityRequest{
			Metadata: th.validMetadata,
			Config:   configGasLimit,
			Inputs:   th.validInputs,
		}

		meta := types.TxMeta{WorkflowExecutionID: &req.Metadata.WorkflowExecutionID, GasLimit: big.NewInt(500000)}
		th.cw.On("SubmitTransaction", mock.Anything, "forwarder", "report", mock.Anything, mock.Anything, th.forwarderAddr, &meta, mock.Anything).Return(types.ErrSettingTransactionGasLimitNotSupported)

		ctx := testutils.Context(t)
		_, err2 = th.writeTarget.Execute(ctx, req)
		require.Error(t, err2)
	})

	t.Run("retries without gas limit when ChainWriter's SubmitTransaction returns error due to gas limit not supported", func(t *testing.T) {
		configGasLimit, err2 := values.NewMap(map[string]any{
			"Address":  th.forwarderAddr,
			"GasLimit": 500000,
		})
		require.NoError(t, err2)
		req := capabilities.CapabilityRequest{
			Metadata: th.validMetadata,
			Config:   configGasLimit,
			Inputs:   th.validInputs,
		}

		meta := types.TxMeta{WorkflowExecutionID: &req.Metadata.WorkflowExecutionID, GasLimit: big.NewInt(500000)}
		th.cw.On("SubmitTransaction", mock.Anything, "forwarder", "report", mock.Anything, mock.Anything, th.forwarderAddr, &meta, mock.Anything).Return(types.ErrSettingTransactionGasLimitNotSupported)
		meta = types.TxMeta{WorkflowExecutionID: &req.Metadata.WorkflowExecutionID}
		th.cw.On("SubmitTransaction", mock.Anything, "forwarder", "report", mock.Anything, mock.Anything, th.forwarderAddr, &meta, mock.Anything).Return(nil)

		configGasLimit, err := values.NewMap(map[string]any{
			"Address": th.forwarderAddr,
		})
		require.NoError(t, err)
		req = capabilities.CapabilityRequest{
			Metadata: th.validMetadata,
			Config:   configGasLimit,
			Inputs:   th.validInputs,
		}

		ctx := testutils.Context(t)
		_, err2 = th.writeTarget.Execute(ctx, req)
		require.Error(t, err2)
	})

	t.Run("fails when ChainReader's GetLatestValue returns error", func(t *testing.T) {
		req := capabilities.CapabilityRequest{
			Metadata: th.validMetadata,
			Config:   th.config,
			Inputs:   th.validInputs,
		}
		th.cr.EXPECT().GetLatestValue(mock.Anything, th.binding.ReadIdentifier("getTransmissionInfo"), mock.Anything, mock.Anything, mock.Anything).Return(errors.New("reader error"))

		ctx := testutils.Context(t)
		_, err := th.writeTarget.Execute(ctx, req)
		require.Error(t, err)
	})

	t.Run("fails with invalid config", func(t *testing.T) {
		invalidConfig, err2 := values.NewMap(map[string]any{
			"Address": "invalid-address",
		})
		require.NoError(t, err2)

		req := capabilities.CapabilityRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID: "test-id",
			},
			Config: invalidConfig,
			Inputs: th.validInputs,
		}
		ctx := testutils.Context(t)
		_, err2 = th.writeTarget.Execute(ctx, req)
		require.Error(t, err2)
	})

	t.Run("fails with nil config", func(t *testing.T) {
		req := capabilities.CapabilityRequest{
			Metadata: th.validMetadata,
			Config:   nil,
			Inputs:   th.validInputs,
		}
		ctx := testutils.Context(t)
		_, err2 := th.writeTarget.Execute(ctx, req)
		require.Error(t, err2)
	})

	t.Run("fails with nil inputs", func(t *testing.T) {
		req := capabilities.CapabilityRequest{
			Metadata: th.validMetadata,
			Config:   th.config,
			Inputs:   nil,
		}
		ctx := testutils.Context(t)
		_, err2 := th.writeTarget.Execute(ctx, req)
		require.Error(t, err2)
	})
}

func TestWriteTarget_ValidateRequest(t *testing.T) {
	th := setup(t)
	tests := []struct {
		name          string
		modifyRequest func(*capabilities.CapabilityRequest)
		expectedError string
	}{
		{
			name: "non-matching WorkflowOwner",
			modifyRequest: func(req *capabilities.CapabilityRequest) {
				req.Metadata.WorkflowOwner = "nonmatchingowner"
			},
			expectedError: "WorkflowOwner in the report does not match WorkflowOwner in the request metadata",
		},
		{
			name: "non-matching WorkflowName",
			modifyRequest: func(req *capabilities.CapabilityRequest) {
				req.Metadata.WorkflowName = hex.EncodeToString([]byte("nonmatchingname"))
			},
			expectedError: "WorkflowName in the report does not match WorkflowName in the request metadata",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := capabilities.CapabilityRequest{
				Metadata: th.validMetadata,
				Config:   th.config,
				Inputs:   th.validInputs,
			}
			tt.modifyRequest(&req)

			ctx := testutils.Context(t)
			_, err := th.writeTarget.Execute(ctx, req)
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestWriteTarget_UnconfirmedTransaction(t *testing.T) {
	t.Run("succeeds when transaction is unconfirmed but transmission succeeded ", func(t *testing.T) {
		th := setup(t)
		th.cw.On("SubmitTransaction", mock.Anything, "forwarder", "report", mock.Anything, mock.Anything, th.forwarderAddr, mock.Anything, mock.Anything).Return(nil).Once()
		callCount := 0
		th.cr.On("GetLatestValue", mock.Anything, th.binding.ReadIdentifier("getTransmissionInfo"), mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			transmissionInfo := args.Get(4).(*targets.TransmissionInfo)
			if callCount == 0 {
				*transmissionInfo = targets.TransmissionInfo{
					GasLimit:        big.NewInt(0),
					InvalidReceiver: false,
					State:           targets.TransmissionStateNotAttempted,
					Success:         false,
					TransmissionID:  [32]byte{},
					Transmitter:     common.HexToAddress("0x0"),
				}
			} else {
				*transmissionInfo = targets.TransmissionInfo{
					GasLimit:        big.NewInt(0),
					InvalidReceiver: false,
					State:           targets.TransmissionStateSucceeded,
					Success:         true,
					TransmissionID:  [32]byte{},
					Transmitter:     common.HexToAddress("0x0"),
				}
			}
			callCount++
		})
		req := capabilities.CapabilityRequest{
			Metadata: th.validMetadata,
			Config:   th.config,
			Inputs:   th.validInputs,
		}

		th.cw.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(types.Unconfirmed, nil).Once()

		ctx := testutils.Context(t)
		response, err2 := th.writeTarget.Execute(ctx, req)
		require.NoError(t, err2)
		require.NotNil(t, response)
	})

	t.Run("getTransmissionInfo twice, transaction written to the forwarder, but failed to be written to the consumer contract because of revert in receiver", func(t *testing.T) {
		th := setup(t)
		th.cw.On("SubmitTransaction", mock.Anything, "forwarder", "report", mock.Anything, mock.Anything, th.forwarderAddr, mock.Anything, mock.Anything).Return(nil).Once()
		callCount := 0
		th.cr.On("GetLatestValue", mock.Anything, th.binding.ReadIdentifier("getTransmissionInfo"), mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			transmissionInfo := args.Get(4).(*targets.TransmissionInfo)
			if callCount == 0 {
				*transmissionInfo = targets.TransmissionInfo{
					GasLimit:        big.NewInt(0),
					InvalidReceiver: false,
					State:           targets.TransmissionStateNotAttempted,
					Success:         false,
					TransmissionID:  [32]byte{},
					Transmitter:     common.HexToAddress("0x0"),
				}
			} else {
				*transmissionInfo = targets.TransmissionInfo{
					GasLimit:        big.NewInt(0),
					InvalidReceiver: false,
					State:           targets.TransmissionStateFailed,
					Success:         false,
					TransmissionID:  [32]byte{},
					Transmitter:     common.HexToAddress("0x0"),
				}
			}
			callCount++
		})
		req := capabilities.CapabilityRequest{
			Metadata: th.validMetadata,
			Config:   th.config,
			Inputs:   th.validInputs,
		}

		th.cw.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(types.Unconfirmed, nil).Once()
		ctx := testutils.Context(t)
		_, err2 := th.writeTarget.Execute(ctx, req)
		require.Error(t, err2)
		require.ErrorIs(t, err2, targets.ErrTxFailed)
	})

	t.Run("getTransmissionInfo twice, transaction written to the forwarder, but failed to be written to the consumer contract because of invalid receiver", func(t *testing.T) {
		th := setup(t)
		th.cw.On("SubmitTransaction", mock.Anything, "forwarder", "report", mock.Anything, mock.Anything, th.forwarderAddr, mock.Anything, mock.Anything).Return(nil).Once()
		callCount := 0
		th.cr.On("GetLatestValue", mock.Anything, th.binding.ReadIdentifier("getTransmissionInfo"), mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			transmissionInfo := args.Get(4).(*targets.TransmissionInfo)
			if callCount == 0 {
				*transmissionInfo = targets.TransmissionInfo{
					GasLimit:        big.NewInt(0),
					InvalidReceiver: false,
					State:           targets.TransmissionStateNotAttempted,
					Success:         false,
					TransmissionID:  [32]byte{},
					Transmitter:     common.HexToAddress("0x0"),
				}
			} else {
				*transmissionInfo = targets.TransmissionInfo{
					GasLimit:        big.NewInt(0),
					InvalidReceiver: true,
					State:           targets.TransmissionStateInvalidReceiver,
					Success:         false,
					TransmissionID:  [32]byte{},
					Transmitter:     common.HexToAddress("0x0"),
				}
			}
			callCount++
		})
		req := capabilities.CapabilityRequest{
			Metadata: th.validMetadata,
			Config:   th.config,
			Inputs:   th.validInputs,
		}

		th.cw.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(types.Unconfirmed, nil).Once()
		ctx := testutils.Context(t)
		_, err2 := th.writeTarget.Execute(ctx, req)
		require.Error(t, err2)
		require.ErrorIs(t, err2, targets.ErrTxFailed)
	})

	t.Run("getTransmissionInfo once, transaction written to the forwarder, but failed to be written to the consumer contract because of invalid receiver", func(t *testing.T) {
		th := setup(t)
		th.cr.On("GetLatestValue", mock.Anything, th.binding.ReadIdentifier("getTransmissionInfo"), mock.Anything, mock.Anything, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
			transmissionInfo := args.Get(4).(*targets.TransmissionInfo)
			*transmissionInfo = targets.TransmissionInfo{
				GasLimit:        big.NewInt(0),
				InvalidReceiver: true,
				State:           targets.TransmissionStateInvalidReceiver,
				Success:         false,
				TransmissionID:  [32]byte{},
				Transmitter:     common.HexToAddress("0x0"),
			}
		})
		req := capabilities.CapabilityRequest{
			Metadata: th.validMetadata,
			Config:   th.config,
			Inputs:   th.validInputs,
		}

		ctx := testutils.Context(t)
		_, err2 := th.writeTarget.Execute(ctx, req)
		require.Error(t, err2)
		require.ErrorIs(t, err2, targets.ErrTxFailed)
	})

	t.Run("getTransmissionInfo once, transaction written to the forwarder, but failed to be written to the consumer contract because of revert in receiver", func(t *testing.T) {
		th := setup(t)
		th.cr.On("GetLatestValue", mock.Anything, th.binding.ReadIdentifier("getTransmissionInfo"), mock.Anything, mock.Anything, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
			transmissionInfo := args.Get(4).(*targets.TransmissionInfo)
			*transmissionInfo = targets.TransmissionInfo{
				GasLimit:        big.NewInt(0).SetUint64(th.gasLimit + 1), // has sufficient gas
				InvalidReceiver: false,
				State:           targets.TransmissionStateFailed,
				Success:         false,
				TransmissionID:  [32]byte{},
				Transmitter:     common.HexToAddress("0x0"),
			}
		})
		req := capabilities.CapabilityRequest{
			Metadata: th.validMetadata,
			Config:   th.config,
			Inputs:   th.validInputs,
		}

		ctx := testutils.Context(t)
		_, err2 := th.writeTarget.Execute(ctx, req)
		require.Error(t, err2)
		require.ErrorIs(t, err2, targets.ErrTxFailed)
	})

	t.Run("getTransmissionInfo twice, first time receiver reverted because of insufficient gas, second time succeeded", func(t *testing.T) {
		th := setup(t)
		th.cw.On("SubmitTransaction", mock.Anything, "forwarder", "report", mock.Anything, mock.Anything, th.forwarderAddr, mock.Anything, mock.Anything).Return(nil).Once()
		callCount := 0
		th.cr.On("GetLatestValue", mock.Anything, th.binding.ReadIdentifier("getTransmissionInfo"), mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			transmissionInfo := args.Get(4).(*targets.TransmissionInfo)
			if callCount == 0 {
				*transmissionInfo = targets.TransmissionInfo{
					GasLimit:        big.NewInt(0).SetUint64(th.gasLimit - 1), // has insufficient gas
					InvalidReceiver: false,
					State:           targets.TransmissionStateFailed,
					Success:         false,
					TransmissionID:  [32]byte{},
					Transmitter:     common.HexToAddress("0x0"),
				}
			} else {
				*transmissionInfo = targets.TransmissionInfo{
					GasLimit:        big.NewInt(0).SetUint64(th.gasLimit + 1), // has sufficient gas
					InvalidReceiver: false,
					State:           targets.TransmissionStateSucceeded,
					Success:         true,
					TransmissionID:  [32]byte{},
					Transmitter:     common.HexToAddress("0x0"),
				}
			}
			callCount++
		})
		req := capabilities.CapabilityRequest{
			Metadata: th.validMetadata,
			Config:   th.config,
			Inputs:   th.validInputs,
		}

		th.cw.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(types.Unconfirmed, nil).Once()

		ctx := testutils.Context(t)
		response, err2 := th.writeTarget.Execute(ctx, req)
		require.NoError(t, err2)
		require.NotNil(t, response)
	})
}
