package executable

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable/request"
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
	// for the active map (delete is a no-op), and the completedAt stamp is refreshed — fine,
	// since the goal is "did this complete recently".
	c.markCompleted("Execute:wfExecID:refID")
	assert.Empty(t, c.requestIDToCallerRequest)
	assert.Len(t, c.completedRequestIDs, 1)
}

// TestClient_Receive_RecognizesRecentlyCompletedRequests asserts the data-path that Receive
// uses to distinguish a late post-quorum response from a truly unknown message: after
// markCompleted, the active map no longer has the request but completedRequestIDs does. This
// is the regression guard for Copilot's review note on the original PR — every non-quorum
// signer's late response otherwise logged "received response for unknown message ID" at Warn,
// generating N-F-1 warnings per Execute call.
func TestClient_Receive_RecognizesRecentlyCompletedRequests(t *testing.T) {
	c := &client{
		requestIDToCallerRequest: map[string]*request.ClientRequest{},
		completedRequestIDs:      map[string]time.Time{},
		mutex:                    sync.Mutex{},
	}

	const messageID = "Execute:wfExecID:refID"
	c.requestIDToCallerRequest[messageID] = &request.ClientRequest{}

	// Simulate Execute returning: active entry removed, completedRequestIDs populated.
	c.markCompleted(messageID)

	// Receive's decision: if requestIDToCallerRequest has no entry, check completedRequestIDs
	// before logging Warn. After markCompleted, the late-response branch must fire.
	_, isActive := c.requestIDToCallerRequest[messageID]
	assert.False(t, isActive, "request should no longer be in the active map")
	_, recentlyCompleted := c.completedRequestIDs[messageID]
	assert.True(t, recentlyCompleted, "request should be recognized as recently completed")
}
