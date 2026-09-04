package confidentialrelay

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/patrickmn/go-cache"

	confidentialrelaytypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialrelay"
)

const capabilityCallDomain = "call-capability"
const secretsGetDomain = "get-secrets"

// capExecKey is the deterministic cache key for a capability-exec request,
// built from its logical identity: the (workflow, execution, step, capability)
// tuple the relay-DON signature binds to. Avoids hashing: the fields are
// required non-empty by Validate, so a plain join is stable and debuggable.
func capExecKey(p confidentialrelaytypes.CapabilityRequestParams) string {
	return strings.Join([]string{capabilityCallDomain, p.WorkflowID, p.ExecutionID, p.ReferenceID, p.CapabilityID}, "/")
}

// secretsKey is the deterministic cache key for a secrets-get request, built
// from its logical identity: workflow, execution, callback id.
func secretsKey(p confidentialrelaytypes.SecretsRequestParams) string {
	return strings.Join([]string{secretsGetDomain, p.WorkflowID, p.ExecutionID, strconv.Itoa(int(p.CallbackID))}, "/")
}

// pendingRequest is the in-flight marker for a logical identity. The owner
// executes and later completes it; duplicates wait on done and then respond
// with the owner's published signed result (re-wrapped with their own request
// id by the caller). signed is written before done is closed, so a waiter
// that observes done closed may read signed without further synchronization.
// Closing done is a broadcast: any number of waiters may listen concurrently,
// and all are woken together.
type pendingRequest struct {
	done chan struct{}
	// signed is the owner's published result, set before done closes.
	signed any
}

// checkPendingRequest returns the pending request in flight for key, or nil
// when none was — in which case one was just registered and this caller is
// implicitly its owner: execute, then completePendingRequest. A non-nil
// return is another caller's in-flight execution: wait for it and respond
// with its result, instead of re-executing into the remote request server's
// duplicate-requester error ("request already received from peer", a
// well-formed terminal result for the enclave's retry loop). Registration is
// atomic (go-cache Add), so exactly one caller gets nil; the entry is
// TTL-bounded so a crashed owner ages out via the background cleanup. A
// non-nil error means the pending set holds an entry of an unexpected type
// (or one vanished between Add and Get): its state is unknown, so the caller
// fails loudly rather than guesses.
func (h *Handler) checkOrCreatePendingRequest(key string) (*pendingRequest, error) {
	mine := &pendingRequest{done: make(chan struct{})}
	if h.pendingRequests.Add(key, mine, cache.DefaultExpiration) == nil {
		h.lggr.Debugw("registered pending relay request, executing", "key", key)
		return nil, nil
	}
	if existing, ok := h.pendingRequests.Get(key); ok {
		if p, ok := existing.(*pendingRequest); ok {
			h.lggr.Debugw("relay request already in flight, waiting for its result", "key", key)
			return p, nil
		}
	}
	h.lggr.Errorw("pending request cache holds an entry of unexpected type", "key", key)
	return nil, fmt.Errorf("pending set holds an invalid entry for key %q", key)
}

// completePendingRequest publishes the owner's signed result and wakes any
// waiters. Called on the execution's success path only, after the memo is set,
// so requests arriving later hit the memo instead. An owner that dies without
// completing never publishes; the TTL ages its entry out and waiters fail on
// their own deadlines, which is the "unknown failure" outcome.
func (h *Handler) completePendingRequest(key string, signed any) {
	if v, ok := h.pendingRequests.Get(key); ok {
		if pr, ok := v.(*pendingRequest); ok {
			pr.signed = signed
			close(pr.done)
		}
	}
	h.pendingRequests.Delete(key)
}

// waitForPendingRequest blocks until the owner completes (returning its
// published signed result) or ctx is done. A waiter whose context expires
// drops: the gateway has already timed out its request by then, so there is
// nobody to answer.
func waitForPendingRequest(ctx context.Context, pr *pendingRequest) (any, bool) {
	select {
	case <-pr.done:
		return pr.signed, true
	case <-ctx.Done():
		return nil, false
	}
}
