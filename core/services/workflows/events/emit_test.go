package events_test

import (
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	pb "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
)

func TestEmitTimestampNano(t *testing.T) {
	// t.Parallel() // TODO: the beholder tester uses t.SetEnv and cannot use t.Parallel
	triggerID := "trigger_" + uuid.NewString()
	executionID := "execution_" + uuid.NewString()
	capabilityID := "capability_" + uuid.NewString()
	stepRef := "step"
	beholder := tests.Beholder(t)
	labeler := custmsg.NewLabeler()

	// basic regex for RFC3339Nano using ISO 8601 or tz offset format
	timeMatcher := regexp.MustCompile(`[0-9\-]{10}T[0-9:]{8}\.[0-9Z\-:\+]+`)

	t.Run(events.WorkflowExecutionStarted, func(t *testing.T) {
		require.NoError(t, events.EmitExecutionStartedEvent(t.Context(), labeler, triggerID, executionID))

		msgs := beholder.Messages(t, "beholder_entity", "workflows.v1."+events.WorkflowExecutionStarted)
		require.Len(t, msgs, 1)

		var expected pb.WorkflowExecutionStarted

		require.NoError(t, proto.Unmarshal(msgs[0].Body, &expected))
		assert.True(t, timeMatcher.MatchString(expected.Timestamp), expected.Timestamp)
	})

	t.Run(events.WorkflowExecutionFinished, func(t *testing.T) {
		require.NoError(t, events.EmitExecutionFinishedEvent(t.Context(), labeler, triggerID, executionID))

		msgs := beholder.Messages(t, "beholder_entity", "workflows.v1."+events.WorkflowExecutionFinished)
		require.Len(t, msgs, 1)

		var expected pb.WorkflowExecutionFinished

		require.NoError(t, proto.Unmarshal(msgs[0].Body, &expected))
		assert.True(t, timeMatcher.MatchString(expected.Timestamp), expected.Timestamp)
	})

	t.Run(events.CapabilityExecutionStarted, func(t *testing.T) {
		require.NoError(t, events.EmitCapabilityStartedEvent(t.Context(), labeler, executionID, capabilityID, stepRef))

		msgs := beholder.Messages(t, "beholder_entity", "workflows.v1."+events.CapabilityExecutionStarted)
		require.Len(t, msgs, 1)

		var expected pb.CapabilityExecutionStarted

		require.NoError(t, proto.Unmarshal(msgs[0].Body, &expected))
		assert.True(t, timeMatcher.MatchString(expected.Timestamp), expected.Timestamp)
	})

	t.Run(events.CapabilityExecutionFinished, func(t *testing.T) {
		require.NoError(t, events.EmitCapabilityFinishedEvent(t.Context(), labeler, executionID, capabilityID, stepRef, "status"))

		msgs := beholder.Messages(t, "beholder_entity", "workflows.v1."+events.CapabilityExecutionFinished)
		require.Len(t, msgs, 1)

		var expected pb.CapabilityExecutionFinished

		require.NoError(t, proto.Unmarshal(msgs[0].Body, &expected))
		assert.True(t, timeMatcher.MatchString(expected.Timestamp), expected.Timestamp)
	})
}
