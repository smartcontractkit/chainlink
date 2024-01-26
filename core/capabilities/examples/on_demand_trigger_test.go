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
	tr := NewOnDemandTrigger()
	ctx := context.Background()

	callback := make(chan capabilities.CapabilityResponse, 10)
	m, err := values.NewMap(map[string]any{"weid": "hello"})
	require.NoError(t, err)

	err = tr.RegisterTrigger(ctx, callback, m)
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
