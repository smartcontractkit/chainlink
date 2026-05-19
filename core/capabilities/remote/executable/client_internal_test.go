package executable

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable/request"
)

// TestClient_deleteRequest_RemovesEntryFromMap is a focused unit test for the helper added
// alongside the request-map leak fix. Execute defers c.deleteRequest(req.ID()) so the entry
// added by storeRequest is removed regardless of how Execute exits (success, response error,
// ctx cancellation). Without that cleanup, repeated Execute calls with the same workflow
// execution ID + reference ID failed in storeRequest with "request for ID ... already exists"
// until the expireRequests ticker reaped the stale entry.
func TestClient_deleteRequest_RemovesEntryFromMap(t *testing.T) {
	c := &client{
		requestIDToCallerRequest: map[string]*request.ClientRequest{},
		mutex:                    sync.Mutex{},
	}

	// Seed the map directly with a non-nil placeholder; deleteRequest only deletes by key, so the
	// value's internal state is irrelevant for this test.
	c.requestIDToCallerRequest["Execute:wfExecID:refID"] = &request.ClientRequest{}
	assert.Len(t, c.requestIDToCallerRequest, 1)

	c.deleteRequest("Execute:wfExecID:refID")
	assert.Empty(t, c.requestIDToCallerRequest, "deleteRequest should remove the entry")

	// Idempotent: deleting an already-absent key is a no-op (matters because expireRequests can
	// race with the Execute defer and reap the entry first).
	c.deleteRequest("Execute:wfExecID:refID")
	assert.Empty(t, c.requestIDToCallerRequest)
}
