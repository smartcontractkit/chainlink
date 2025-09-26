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
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"

	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
)

func TestEmit(t *testing.T) {
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
		require.NoError(t, events.EmitExecutionFinishedEvent(t.Context(), labels, "status", executionID, nil))
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
		require.NoError(t, events.EmitCapabilityFinishedEvent(t.Context(), labels, executionID, capabilityID, stepRef, "status", nil))
		require.Len(t, labels, 1)

		msgs := beholderObserver.Messages(t, "beholder_entity", "workflows.v1."+events.CapabilityExecutionFinished)
		require.Len(t, msgs, 1)

		var expected pb.CapabilityExecutionFinished

		require.NoError(t, proto.Unmarshal(msgs[0].Body, &expected))
		assert.True(t, timeMatcher.MatchString(expected.Timestamp), expected.Timestamp)
	})

	t.Run(events.UserLogs, func(t *testing.T) {
		logLines := []*pb.LogLine{
			{
				NodeTimestamp: "2024-01-01T00:00:00Z",
				Message:       "Test log message",
			},
			{
				NodeTimestamp: "2024-01-01T00:01:00Z",
				Message:       "Second log message",
			},
		}
		require.NoError(t, events.EmitUserLogs(t.Context(), labels, logLines, executionID))
		require.Len(t, labels, 1)

		// Verify v1 event
		v1Msgs := beholderObserver.Messages(t, "beholder_entity", "workflows.v1."+events.UserLogs)
		require.Len(t, v1Msgs, 1)

		var receivedV1 pb.UserLogs
		require.NoError(t, proto.Unmarshal(v1Msgs[0].Body, &receivedV1))
		assert.Equal(t, logLines[0].Message, receivedV1.LogLines[0].Message)
		assert.Equal(t, logLines[1].Message, receivedV1.LogLines[1].Message)

		// Verify v2 events are emitted (one per log line)
		v2Msgs := beholderObserver.Messages(t, "beholder_entity", "workflows.v2."+events.WorkflowUserLog)
		require.Len(t, v2Msgs, 2)

		var msg1 eventsv2.WorkflowUserLog
		require.NoError(t, proto.Unmarshal(v2Msgs[0].Body, &msg1))
		assert.Equal(t, executionID, msg1.WorkflowExecutionID)
		assert.Equal(t, logLines[0].NodeTimestamp, msg1.Timestamp)
		assert.Equal(t, logLines[0].Message, msg1.Msg)
		assert.NotNil(t, msg1.CreInfo)
		assert.NotNil(t, msg1.Workflow)

		var msg2 eventsv2.WorkflowUserLog
		require.NoError(t, proto.Unmarshal(v2Msgs[1].Body, &msg2))
		assert.Equal(t, executionID, msg2.WorkflowExecutionID)
		assert.Equal(t, logLines[1].NodeTimestamp, msg2.Timestamp)
		assert.Equal(t, logLines[1].Message, msg2.Msg)
		assert.NotNil(t, msg2.CreInfo)
		assert.NotNil(t, msg2.Workflow)
		// Labels not utilized, left unchecked
	})
}
