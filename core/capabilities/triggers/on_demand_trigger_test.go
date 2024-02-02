package triggers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

const testID = "test-id-1"

func TestOnDemand(t *testing.T) {
	r := capabilities.NewRegistry()
	tr := NewOnDemand()
	ctx := testutils.Context(t)

	err := r.Add(ctx, tr)
	require.NoError(t, err)

	trigger, err := r.GetTrigger(ctx, tr.Info().ID)
	require.NoError(t, err)

	callback := make(chan capabilities.CapabilityResponse, 10)

	req := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowExecutionID: testID,
		},
	}
	err = trigger.RegisterTrigger(ctx, callback, req)
	require.NoError(t, err)

	er := capabilities.CapabilityResponse{
		Value: &values.String{Underlying: testID},
	}

	err = tr.FanOutEvent(ctx, er)
	require.NoError(t, err)

	assert.Len(t, callback, 1)
	assert.Equal(t, er, <-callback)
}

func TestOnDemand_ChannelDoesntExist(t *testing.T) {
	tr := NewOnDemand()
	ctx := testutils.Context(t)

	er := capabilities.CapabilityResponse{
		Value: &values.String{Underlying: testID},
	}
	err := tr.SendEvent(ctx, testID, er)
	assert.ErrorContains(t, err, "no registration")
}

func TestOnDemand_(t *testing.T) {
	tr := NewOnDemand()
	ctx := testutils.Context(t)

	req := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID: "hello",
		},
	}
	callback := make(chan capabilities.CapabilityResponse, 10)

	err := tr.RegisterTrigger(ctx, callback, req)
	require.NoError(t, err)

	er := capabilities.CapabilityResponse{
		Value: &values.String{Underlying: testID},
	}
	err = tr.SendEvent(ctx, "hello", er)
	require.NoError(t, err)

	assert.Len(t, callback, 1)
	assert.Equal(t, er, <-callback)
}
