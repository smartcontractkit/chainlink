package v2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	protoevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

func TestDisallowedExecutionHelper_GetSecretsReturnsError(t *testing.T) {
	t.Parallel()

	h := NewDisallowedExecutionHelper(logger.Test(t), make(chan<- *protoevents.LogLine), &types.LocalTimeProvider{})

	_, err := h.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{})
	require.ErrorIs(t, err, ErrSecretsCallDuringSubscription)
	require.ErrorContains(t, err, "secrets calls cannot be made during trigger subscription")
}

func TestDisallowedExecutionHelper_CallCapabilityReturnsError(t *testing.T) {
	t.Parallel()

	h := NewDisallowedExecutionHelper(logger.Test(t), make(chan<- *protoevents.LogLine), &types.LocalTimeProvider{})

	_, err := h.CallCapability(context.Background(), &sdkpb.CapabilityRequest{})
	require.ErrorIs(t, err, ErrCapabilityCallDuringSubscription)
	require.ErrorContains(t, err, "capability calls cannot be made during trigger subscription")
}

func TestDisallowedExecutionHelper_OtherHelpers(t *testing.T) {
	t.Parallel()

	h := NewDisallowedExecutionHelper(logger.Test(t), make(chan<- *protoevents.LogLine), &types.LocalTimeProvider{})

	require.Empty(t, h.GetWorkflowExecutionID())
	require.Error(t, h.EmitUserMetric(context.Background(), &eventsv2.WorkflowUserMetric{}))
}
