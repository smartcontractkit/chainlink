package targets_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

//go:generate mockery --quiet --name ChainWriter --srcpkg=github.com/smartcontractkit/chainlink-common/pkg/types --output ./mocks/ --case=underscore
//go:generate mockery --quiet --name ChainReader --srcpkg=github.com/smartcontractkit/chainlink-common/pkg/types --output ./mocks/ --case=underscore

func TestWriteTarget(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := context.Background()

	cw := mocks.NewChainWriter(t)
	cr := mocks.NewChainReader(t)

	forwarderA := testutils.NewAddress()
	forwarderAddr := forwarderA.Hex()

	writeTarget := targets.NewWriteTarget(lggr, "test-write-target@1.0.0", cr, cw, forwarderAddr)
	require.NotNil(t, writeTarget)

	config, err := values.NewMap(map[string]any{
		"Address": forwarderAddr,
	})
	require.NoError(t, err)

	validInputs, err := values.NewMap(map[string]any{
		"signed_report": map[string]any{
			"report":     []byte{1, 2, 3},
			"signatures": [][]byte{},
		},
	})
	require.NoError(t, err)

	cr.On("GetLatestValue", mock.Anything, "forwarder", "getTransmitter", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		transmitter := args.Get(4).(*common.Address)
		*transmitter = common.HexToAddress("0x0")
	}).Once()

	cw.On("SubmitTransaction", mock.Anything, "forwarder", "report", mock.Anything, mock.Anything, forwarderAddr, mock.Anything, mock.Anything).Return(nil).Once()

	t.Run("succeeds with valid report", func(t *testing.T) {
		req := capabilities.CapabilityRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID: "test-id",
			},
			Config: config,
			Inputs: validInputs,
		}

		ch, err2 := writeTarget.Execute(ctx, req)
		require.NoError(t, err2)
		response := <-ch
		require.NotNil(t, response)
	})

	t.Run("succeeds with empty report", func(t *testing.T) {
		emptyInputs, err2 := values.NewMap(map[string]any{
			"signed_report": map[string]any{
				"report": []byte{},
			},
			"signatures": [][]byte{},
		})

		require.NoError(t, err2)
		req := capabilities.CapabilityRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowExecutionID: "test-id",
			},
			Config: config,
			Inputs: emptyInputs,
		}

		ch, err2 := writeTarget.Execute(ctx, req)
		require.NoError(t, err2)
		response := <-ch
		require.Nil(t, response.Value)
	})

	t.Run("fails when ChainReader's GetLatestValue returns error", func(t *testing.T) {
		req := capabilities.CapabilityRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID: "test-id",
			},
			Config: config,
			Inputs: validInputs,
		}
		cr.On("GetLatestValue", mock.Anything, "forwarder", "getTransmitter", mock.Anything, mock.Anything).Return(errors.New("reader error"))

		_, err = writeTarget.Execute(ctx, req)
		require.Error(t, err)
	})

	t.Run("fails when ChainWriter's SubmitTransaction returns error", func(t *testing.T) {
		req := capabilities.CapabilityRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID: "test-id",
			},
			Config: config,
			Inputs: validInputs,
		}
		cw.On("SubmitTransaction", mock.Anything, "forwarder", "report", mock.Anything, mock.Anything, forwarderAddr, mock.Anything, mock.Anything).Return(errors.New("writer error"))

		_, err = writeTarget.Execute(ctx, req)
		require.Error(t, err)
	})

	t.Run("fails with invalid config", func(t *testing.T) {
		invalidConfig, err := values.NewMap(map[string]any{
			"Address": "invalid-address",
		})
		require.NoError(t, err)

		req := capabilities.CapabilityRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID: "test-id",
			},
			Config: invalidConfig,
			Inputs: validInputs,
		}
		_, err = writeTarget.Execute(ctx, req)
		require.Error(t, err)
	})
}
