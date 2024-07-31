package targets_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	coreMocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func TestWriteTarget(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := context.Background()

	cw := mocks.NewChainWriter(t)
	cr := mocks.NewChainReader(t)

	forwarderA := testutils.NewAddress()
	forwarderAddr := forwarderA.Hex()

	writeTarget := targets.NewWriteTarget(lggr, "test-write-target@1.0.0", cr, cw, forwarderAddr, nil, 0)
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

	cr.On("Bind", mock.Anything, []types.BoundContract{
		{
			Address: forwarderAddr,
			Name:    "forwarder",
		},
	}).Return(nil)

	cr.On("GetLatestValue", mock.Anything, "forwarder", "getTransmitter", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		transmitter := args.Get(5).(*common.Address)
		*transmitter = common.HexToAddress("0x0")
	}).Once()

	cw.On("SubmitTransaction", mock.Anything, "forwarder", "report", mock.Anything, mock.Anything, forwarderAddr, mock.Anything, mock.Anything).Return(nil).Once()

	t.Run("succeeds with valid report", func(t *testing.T) {
		req := capabilities.CapabilityRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID:          "test-id",
				WorkflowExecutionID: "dd3709ac7d8dd6fa4fae0fb87b73f318a4da2526c123e159b72435e3b2fe8751",
				WorkflowDonID:       1,
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
				WorkflowID:          "test-id",
				WorkflowExecutionID: "dd3709ac7d8dd6fa4fae0fb87b73f318a4da2526c123e159b72435e3b2fe8751",
				WorkflowDonID:       1,
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
				WorkflowID:          "test-id",
				WorkflowExecutionID: "dd3709ac7d8dd6fa4fae0fb87b73f318a4da2526c123e159b72435e3b2fe8751",
				WorkflowDonID:       1,
			},
			Config: config,
			Inputs: validInputs,
		}
		cr.On("GetLatestValue", mock.Anything, "forwarder", "getTransmitter", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("reader error"))

		_, err = writeTarget.Execute(ctx, req)
		require.Error(t, err)
	})

	t.Run("fails when ChainWriter's SubmitTransaction returns error", func(t *testing.T) {
		req := capabilities.CapabilityRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID:          "test-id",
				WorkflowExecutionID: "dd3709ac7d8dd6fa4fae0fb87b73f318a4da2526c123e159b72435e3b2fe8751",
				WorkflowDonID:       1,
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
				WorkflowID:          "test-id",
				WorkflowExecutionID: "dd3709ac7d8dd6fa4fae0fb87b73f318a4da2526c123e159b72435e3b2fe8751",
				WorkflowDonID:       1,
			},
			Config: invalidConfig,
			Inputs: validInputs,
		}
		_, err = writeTarget.Execute(ctx, req)
		require.Error(t, err)
	})
}
func TestResolveLocalNodeInfo(t *testing.T) {
	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.DebugLevel)

	cw := mocks.NewChainWriter(t)
	cr := mocks.NewChainReader(t)
	forwarderA := testutils.NewAddress()
	forwarderAddr := forwarderA.Hex()

	registry := coreMocks.NewCapabilitiesRegistry(t)
	writeTarget := targets.NewWriteTarget(lggr, "test-write-target@1.0.0", cr, cw, forwarderAddr, registry, 10)
	require.NotNil(t, writeTarget)

	peerIDLogs := observedLogs.FilterFieldKey("peerID")
	require.Equal(t, peerIDLogs.Len(), 0, "we should not have any peerID sugared logs before the registry returns anything")
	workflowDONIDLogs := observedLogs.FilterFieldKey("workflowDONID")
	require.Equal(t, workflowDONIDLogs.Len(), 0, "we should not have any workflowDONID sugared logs before the registry returns anything")
	configVersionLogs := observedLogs.FilterFieldKey("workflowDONConfigVersion")
	require.Equal(t, configVersionLogs.Len(), 0, "we should not have any configVersion sugared logs before the registry returns anything")

	var pid p2ptypes.PeerID
	err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)
	registry.On("GetLocalNode", mock.Anything).Return(capabilities.Node{
		PeerID: &pid,
		WorkflowDON: capabilities.DON{
			ID:            1,
			ConfigVersion: 2,
		},
	}, nil)

	time.Sleep(50 * time.Millisecond)
	peerIDLogs = observedLogs.FilterFieldKey("peerID")
	require.Greater(t, peerIDLogs.Len(), 0, "we should have some peerID sugared logs after the registry returns something")
	workflowDONIDLogs = observedLogs.FilterFieldKey("workflowDONID")
	require.Greater(t, workflowDONIDLogs.Len(), 0, "we should have some workflowDONID sugared logs after the registry returns something")
	configVersionLogs = observedLogs.FilterFieldKey("workflowDONConfigVersion")
	require.Greater(t, configVersionLogs.Len(), 0, "we should have some configVersion sugared logs after the registry returns something")
}
