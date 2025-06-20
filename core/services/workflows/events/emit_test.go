package events_test

import (
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder/beholdertest"
	pb "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
)

func TestEmitTimestampNano(t *testing.T) {
	// t.Parallel() // TODO: the beholder tester uses t.SetEnv and cannot use t.Parallel
	triggerID := "trigger_" + uuid.NewString()
	executionID := "execution_" + uuid.NewString()
	capabilityID := "capability_" + uuid.NewString()
	stepRef := "step"
	beholderObserver := beholdertest.NewObserver(t)
	labels := map[string]string{
		platform.KeyWorkflowOwner: "owner",
	}

	// basic regex for RFC3339Nano using ISO 8601 or tz offset format
	timeMatcher := regexp.MustCompile(`[0-9\-]{10}T[0-9:]{8}\.[0-9Z\-:\+]+`)

	t.Run(events.WorkflowExecutionStarted, func(t *testing.T) {
		require.NoError(t, events.EmitExecutionStartedEvent(t.Context(), labels, triggerID, executionID))
		require.Len(t, labels, 1)

		msgs := beholderObserver.Messages(t, "beholder_entity", "workflows.v1."+events.WorkflowExecutionStarted)
		require.Len(t, msgs, 1)

		var expected pb.WorkflowExecutionStarted

		require.NoError(t, proto.Unmarshal(msgs[0].Body, &expected))
		assert.True(t, timeMatcher.MatchString(expected.Timestamp), expected.Timestamp)
	})

	t.Run(events.WorkflowExecutionFinished, func(t *testing.T) {
		require.NoError(t, events.EmitExecutionFinishedEvent(t.Context(), labels, "status", executionID))
		require.Len(t, labels, 1)

		msgs := beholderObserver.Messages(t, "beholder_entity", "workflows.v1."+events.WorkflowExecutionFinished)
		require.Len(t, msgs, 1)

		var expected pb.WorkflowExecutionFinished

		require.NoError(t, proto.Unmarshal(msgs[0].Body, &expected))
		assert.True(t, timeMatcher.MatchString(expected.Timestamp), expected.Timestamp)
	})

	t.Run(events.CapabilityExecutionStarted, func(t *testing.T) {
		require.NoError(t, events.EmitCapabilityStartedEvent(t.Context(), labels, executionID, capabilityID, stepRef))
		require.Len(t, labels, 1)

		msgs := beholderObserver.Messages(t, "beholder_entity", "workflows.v1."+events.CapabilityExecutionStarted)
		require.Len(t, msgs, 1)

		var expected pb.CapabilityExecutionStarted

		require.NoError(t, proto.Unmarshal(msgs[0].Body, &expected))
		assert.True(t, timeMatcher.MatchString(expected.Timestamp), expected.Timestamp)
	})

	t.Run(events.CapabilityExecutionFinished, func(t *testing.T) {
		require.NoError(t, events.EmitCapabilityFinishedEvent(t.Context(), labels, executionID, capabilityID, stepRef, "status"))
		require.Len(t, labels, 1)

		msgs := beholderObserver.Messages(t, "beholder_entity", "workflows.v1."+events.CapabilityExecutionFinished)
		require.Len(t, msgs, 1)

		var expected pb.CapabilityExecutionFinished

		require.NoError(t, proto.Unmarshal(msgs[0].Body, &expected))
		assert.True(t, timeMatcher.MatchString(expected.Timestamp), expected.Timestamp)
	})
}
