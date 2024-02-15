package consensus

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const workflowTestID = "consensus-workflow-test-id-1"
const workflowExecutionTestID = "consensus-workflow-execution-test-id-1"

func TestOCR3Capability(t *testing.T) {
	n := time.Now()
	fc := clockwork.NewFakeClockAt(n)

	ctx := tests.Context(t)
	s := newStore(1*time.Second, fc)
	cp := newCapability(s, fc)
	require.NoError(t, cp.Start(ctx))

	callback := make(chan capabilities.CapabilityResponse, 10)
	config, err := values.NewMap(map[string]any{"aggregation_method": "data_feeds_2_0"})
	require.NoError(t, err)

	registerReq := capabilities.RegisterToWorkflowRequest{
		Metadata: capabilities.RegistrationMetadata{
			WorkflowID: workflowTestID,
		},
		Config: config,
	}

	err = cp.RegisterToWorkflow(ctx, registerReq)
	require.NoError(t, err)

	ethUsdValue, err := decimal.NewFromString("1.123456")

	require.NoError(t, err)

	obs := map[string]any{"ETH_USD": ethUsdValue}
	inputs, err := values.NewMap(map[string]any{"observations": obs})
	require.NoError(t, err)

	executeReq := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          workflowTestID,
			WorkflowExecutionID: workflowExecutionTestID,
		},
		Config: config,
		Inputs: inputs,
	}
	err = cp.Execute(ctx, callback, executeReq)
	require.NoError(t, err)

	obsv, err := values.NewMap(obs)
	require.NoError(t, err)

	// Mock the oracle returning a response
	err = cp.response(ctx, response{
		Value:     obsv,
		RequestID: workflowExecutionTestID,
	})
	require.NoError(t, err)

	expectedCapabilityResponse := capabilities.CapabilityResponse{
		Value: obsv,
	}
	assert.Len(t, callback, 1)
	assert.Equal(t, expectedCapabilityResponse, <-callback)
}

func TestOCR3Capability_Eviction(t *testing.T) {
	n := time.Now()
	fc := clockwork.NewFakeClockAt(n)

	ctx := tests.Context(t)
	rea := time.Second
	s := newStore(rea, fc)
	cp := newCapability(s, fc)
	require.NoError(t, cp.Start(ctx))

	config, err := values.NewMap(map[string]any{"aggregation_method": "data_feeds_2_0"})
	require.NoError(t, err)

	ethUsdValue, err := decimal.NewFromString("1.123456")
	require.NoError(t, err)
	inputs, err := values.NewMap(map[string]any{"observations": map[string]any{"ETH_USD": ethUsdValue}})
	require.NoError(t, err)

	rid := uuid.New().String()
	executeReq := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          workflowTestID,
			WorkflowExecutionID: rid,
		},
		Config: config,
		Inputs: inputs,
	}
	callback := make(chan capabilities.CapabilityResponse, 10)
	err = cp.Execute(ctx, callback, executeReq)
	require.NoError(t, err)

	fc.Advance(1 * time.Hour)
	resp := <-callback
	assert.ErrorContains(t, resp.Err, "timeout exceeded: could not process request before expiry")

	_, err = s.get(ctx, rid)
	assert.ErrorContains(t, err, "not found")
}
