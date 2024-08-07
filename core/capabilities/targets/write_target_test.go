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
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestWriteTarget(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := context.Background()

	cw := mocks.NewChainWriter(t)
	cr := mocks.NewChainReader(t)

	forwarderA := testutils.NewAddress()
	forwarderAddr := forwarderA.Hex()

	writeTarget := targets.NewWriteTarget(lggr, "test-write-target@1.0.0", cr, cw, forwarderAddr, 400_000)
	require.NotNil(t, writeTarget)

	config, err := values.NewMap(map[string]any{
		"Address": forwarderAddr,
	})
	require.NoError(t, err)

	// type ReportV1Metadata struct {
	// 	Version             uint8
	// 	WorkflowExecutionID [32]byte
	// 	Timestamp           uint32
	// 	DonID               uint32
	// 	DonConfigVersion    uint32
	// 	WorkflowCID         [32]byte
	// 	WorkflowName        [10]byte
	// 	WorkflowOwner       [20]byte
	// 	ReportID            [2]byte
	// }

	reportMetadata := targets.ReportV1Metadata{
		Version:             1,
		WorkflowExecutionID: [32]byte{1, 2, 3},
		Timestamp:           0,
		DonID:               0,
		DonConfigVersion:    0,
		WorkflowCID:         [32]byte{1, 2, 3},
		WorkflowName:        [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		WorkflowOwner:       [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		ReportID:            [2]byte{1, 2},
	}

	reportMetadataBytes, err := reportMetadata.Bytes()
	require.NoError(t, err)

	validInputs, err := values.NewMap(map[string]any{
		"signed_report": map[string]any{
			"report":     reportMetadataBytes,
			"signatures": [][]byte{},
		},
	})
	require.NoError(t, err)

	validMetadata := capabilities.RequestMetadata{
		WorkflowID:          hex.EncodeToString(reportMetadata.WorkflowCID[:]),
		WorkflowOwner:       hex.EncodeToString(reportMetadata.WorkflowOwner[:]),
		WorkflowName:        hex.EncodeToString(reportMetadata.WorkflowName[:]),
		WorkflowExecutionID: hex.EncodeToString(reportMetadata.WorkflowExecutionID[:]),
	}

	cr.On("Bind", mock.Anything, []types.BoundContract{
		{
			Address: forwarderAddr,
			Name:    "forwarder",
		},
	}).Return(nil)

	cr.On("GetLatestValue", mock.Anything, "forwarder", "getTransmissionInfo", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		transmissionInfo := args.Get(5).(*targets.TransmissionInfo)
		*transmissionInfo = targets.TransmissionInfo{
			GasLimit:        big.NewInt(0),
			InvalidReceiver: false,
			State:           0,
			Success:         false,
			TransmissionId:  [32]byte{},
			Transmitter:     common.HexToAddress("0x0"),
		}
	}).Once()

	cw.On("SubmitTransaction", mock.Anything, "forwarder", "report", mock.Anything, mock.Anything, forwarderAddr, mock.Anything, mock.Anything).Return(nil).Once()

	t.Run("succeeds with valid report", func(t *testing.T) {
		req := capabilities.CapabilityRequest{
			Metadata: validMetadata,
			Config:   config,
			Inputs:   validInputs,
		}

		ch, err2 := writeTarget.Execute(ctx, req)
		require.NoError(t, err2)
		response := <-ch
		require.NotNil(t, response)
	})

	t.Run("fails when ChainReader's GetLatestValue returns error", func(t *testing.T) {
		req := capabilities.CapabilityRequest{
			Metadata: validMetadata,
			Config:   config,
			Inputs:   validInputs,
		}
		cr.On("GetLatestValue", mock.Anything, "forwarder", "getTransmissionInfo", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("reader error"))

		_, err = writeTarget.Execute(ctx, req)
		require.Error(t, err)
	})

	t.Run("fails when ChainWriter's SubmitTransaction returns error", func(t *testing.T) {
		req := capabilities.CapabilityRequest{
			Metadata: validMetadata,
			Config:   config,
			Inputs:   validInputs,
		}
		cw.On("SubmitTransaction", mock.Anything, "forwarder", "report", mock.Anything, mock.Anything, forwarderAddr, mock.Anything, mock.Anything).Return(errors.New("writer error"))

		_, err = writeTarget.Execute(ctx, req)
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
			Inputs: validInputs,
		}
		_, err2 = writeTarget.Execute(ctx, req)
		require.Error(t, err2)
	})

	t.Run("fails with nil config", func(t *testing.T) {
		req := capabilities.CapabilityRequest{
			Metadata: validMetadata,
			Config:   nil,
			Inputs:   validInputs,
		}
		_, err2 := writeTarget.Execute(ctx, req)
		require.Error(t, err2)
	})

	t.Run("fails with nil inputs", func(t *testing.T) {
		req := capabilities.CapabilityRequest{
			Metadata: validMetadata,
			Config:   config,
			Inputs:   nil,
		}
		_, err2 := writeTarget.Execute(ctx, req)
		require.Error(t, err2)
	})
}
