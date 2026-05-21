package executable

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
)

// TestClient_markCompleted_MovesEntryFromActiveToCompleted is a focused unit test for the
// helper added alongside the request-map leak fix. Execute defers c.markCompleted(req.ID())
// so the entry added by storeRequest is removed from the active map regardless of how Execute
// exits (success, response error, ctx cancellation), and the requestID is recorded in
// completedRequestIDs so Receive can distinguish expected-but-late post-quorum responses from
// truly unknown messages.
func TestClient_markCompleted_MovesEntryFromActiveToCompleted(t *testing.T) {
	c := &client{
		requestIDToCallerRequest: map[string]*request.ClientRequest{},
		completedRequestIDs:      map[string]time.Time{},
		mutex:                    sync.Mutex{},
	}

	// Seed the active map directly with a non-nil placeholder; markCompleted only deletes by key,
	// so the value's internal state is irrelevant for this test.
	c.requestIDToCallerRequest["Execute:wfExecID:refID"] = &request.ClientRequest{}
	assert.Len(t, c.requestIDToCallerRequest, 1)
	assert.Empty(t, c.completedRequestIDs)

	before := time.Now()
	c.markCompleted("Execute:wfExecID:refID")
	after := time.Now()

	assert.Empty(t, c.requestIDToCallerRequest, "markCompleted should remove the active entry")
	assert.Len(t, c.completedRequestIDs, 1, "markCompleted should record the completion")
	completedAt, ok := c.completedRequestIDs["Execute:wfExecID:refID"]
	assert.True(t, ok)
	assert.True(t, !completedAt.Before(before) && !completedAt.After(after),
		"completion timestamp should be between before and after the markCompleted call")

	// Idempotent for the active map: calling markCompleted on an already-absent key is a no-op
	// for the active map (delete is a no-op), and the completedAt stamp is refreshed, fine,
	// since the goal is "did this complete recently".
	c.markCompleted("Execute:wfExecID:refID")
	assert.Empty(t, c.requestIDToCallerRequest)
	assert.Len(t, c.completedRequestIDs, 1)
}

// TestClient_Receive_LateResponseAfterMarkCompleted_DoesNotWarn invokes Receive with a
// message for a requestID whose Execute has already returned, and asserts that no Warn line
// is emitted. This is the behavioral regression guard for Copilot's review note: every late
// post-quorum response previously logged "received response for unknown message ID" at Warn,
// generating N-F-1 warnings per Execute call in a healthy DON.
func TestClient_Receive_LateResponseAfterMarkCompleted_DoesNotWarn(t *testing.T) {
	lggr, observed := logger.TestObserved(t, zapcore.DebugLevel)
	c := &client{
		lggr:                     lggr,
		requestIDToCallerRequest: map[string]*request.ClientRequest{},
		completedRequestIDs:      map[string]time.Time{},
		mutex:                    sync.Mutex{},
	}

	const messageID = "Execute:wfExecID:refID"

	// Simulate the state right after Execute returned: completedRequestIDs has the ID.
	c.markCompleted(messageID)

	// A late response from a non-quorum signer arrives at Receive.
	c.Receive(context.Background(), &types.MessageBody{
		MessageId: []byte(messageID),
		Sender:    []byte("late-signer-peer"),
	})

	warns := observed.FilterLevelExact(zapcore.WarnLevel).All()
	for _, entry := range warns {
		assert.NotContains(t, entry.Message, "unknown message ID",
			"late response after Execute should not log 'unknown message ID' at Warn")
	}

	debugs := observed.FilterMessageSnippet("late response after Execute returned").All()
	assert.NotEmpty(t, debugs, "late response should log at Debug")
}

// TestClient_Receive_UnknownMessageID_StillWarns asserts the other half of the contract: a
// message ID that was never registered (and is not in completedRequestIDs) still produces a
// Warn. This guards against an over-correction where the late-response Debug branch swallows
// genuinely unknown messages.
func TestClient_Receive_UnknownMessageID_StillWarns(t *testing.T) {
	lggr, observed := logger.TestObserved(t, zapcore.DebugLevel)
	c := &client{
		lggr:                     lggr,
		requestIDToCallerRequest: map[string]*request.ClientRequest{},
		completedRequestIDs:      map[string]time.Time{},
		mutex:                    sync.Mutex{},
	}

	c.Receive(context.Background(), &types.MessageBody{
		MessageId: []byte("Execute:never-registered:refID"),
		Sender:    []byte("rogue-peer"),
	})

	warns := observed.FilterMessageSnippet("unknown message ID").All()
	assert.NotEmpty(t, warns, "truly unknown messageID should still surface as Warn")
}
