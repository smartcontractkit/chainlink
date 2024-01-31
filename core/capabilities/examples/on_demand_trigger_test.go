package examples

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
)

func TestOnDemandTrigger(t *testing.T) {
	r := capabilities.NewRegistry()
	tr := NewOnDemandTrigger()
	ctx := context.Background()

	err := r.Add(ctx, tr)
	require.NoError(t, err)

	trigger, err := r.GetTrigger(ctx, tr.Info().ID)
	require.NoError(t, err)

	callback := make(chan capabilities.CapabilityResponse, 10)

	req := capabilities.CapabilityRequest{
		Metadata: capabilities.Metadata{
			WorkflowID: "hello",
		},
	}
	err = trigger.RegisterTrigger(ctx, callback, req)
	require.NoError(t, err)

	er := capabilities.CapabilityResponse{
		Value: &values.String{"hello"},
	}

	err = tr.FanOutEvent(ctx, er)
	require.NoError(t, err)

	assert.Len(t, callback, 1)
	assert.Equal(t, er, <-callback)
}

func TestOnDemandTrigger_ChannelDoesntExist(t *testing.T) {
	tr := NewOnDemandTrigger()
	ctx := context.Background()

	er := capabilities.CapabilityResponse{
		Value: &values.String{"hello"},
	}
	err := tr.SendEvent(ctx, "hello", er)
	assert.ErrorContains(t, err, "no registration")
}
