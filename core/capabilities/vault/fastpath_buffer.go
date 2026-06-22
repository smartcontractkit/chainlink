package vault

import (
	"context"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

const fastPathResponseFormat = "fastpath/v1"

// FastPathResponseFormat is the vaulttypes.Response.Format value for unsigned fast-path reads.
const FastPathResponseFormat = fastPathResponseFormat

// FastPathRequest is a pending GetSecrets request awaiting per-node share generation.
type FastPathRequest struct {
	ID      string
	Request *vaultcommon.GetSecretsRequest
	Expiry  time.Time
}

// FastPathSource is implemented by the fast-path buffer and consumed by the OCR plugin.
type FastPathSource interface {
	Drain() []FastPathRequest
	Complete(id string, resp *vaulttypes.Response)
	ExpireOlderThan(t time.Time)
}

type fastPathEntry struct {
	request      *vaultcommon.GetSecretsRequest
	requestID    string
	expiry       time.Time
	responseChan chan *vaulttypes.Response
}

type FastPathBuffer struct {
	mu           sync.Mutex
	pending      map[string]*fastPathEntry
	inFlight     map[string]*fastPathEntry
	clock        clockwork.Clock
	expiresAfter time.Duration
	lggr         logger.Logger
}

func newFastPathBuffer(lggr logger.Logger, clock clockwork.Clock, expiresAfter time.Duration) *FastPathBuffer {
	return &FastPathBuffer{
		pending:      make(map[string]*fastPathEntry),
		inFlight:     make(map[string]*fastPathEntry),
		clock:        clock,
		expiresAfter: expiresAfter,
		lggr:         logger.Named(lggr, "VaultFastPathBuffer"),
	}
}

// NewFastPathBuffer constructs a fast-path request buffer shared by the vault capability and OCR plugin.
func NewFastPathBuffer(lggr logger.Logger, clock clockwork.Clock, expiresAfter time.Duration) *FastPathBuffer {
	return newFastPathBuffer(lggr, clock, expiresAfter)
}

// Submit enqueues a fast-path GetSecrets request and returns its response channel.
func (b *FastPathBuffer) Submit(id string, req *vaultcommon.GetSecretsRequest) <-chan *vaulttypes.Response {
	respCh := make(chan *vaulttypes.Response, 1)
	entry := &fastPathEntry{
		request:      req,
		requestID:    id,
		expiry:       b.clock.Now().Add(b.expiresAfter),
		responseChan: respCh,
	}

	b.mu.Lock()
	b.pending[id] = entry
	b.mu.Unlock()

	return respCh
}

func (b *FastPathBuffer) Drain() []FastPathRequest {
	b.mu.Lock()
	defer b.mu.Unlock()

	out := make([]FastPathRequest, 0, len(b.pending))
	for id, entry := range b.pending {
		out = append(out, FastPathRequest{
			ID:      id,
			Request: entry.request,
			Expiry:  entry.expiry,
		})
		b.inFlight[id] = entry
		delete(b.pending, id)
	}
	return out
}

func (b *FastPathBuffer) Complete(id string, resp *vaulttypes.Response) {
	b.mu.Lock()
	entry, ok := b.inFlight[id]
	if ok {
		delete(b.inFlight, id)
	}
	b.mu.Unlock()

	if !ok {
		b.lggr.Debugw("fast-path complete for unknown request", "requestID", id)
		return
	}

	select {
	case entry.responseChan <- resp:
	default:
		b.lggr.Warnw("fast-path response channel full or closed", "requestID", id)
	}
}

func (b *FastPathBuffer) ExpireOlderThan(t time.Time) {
	b.mu.Lock()
	expiredPending := make([]*fastPathEntry, 0)
	for id, entry := range b.pending {
		if !entry.expiry.After(t) {
			expiredPending = append(expiredPending, entry)
			delete(b.pending, id)
		}
	}
	expiredInFlight := make([]*fastPathEntry, 0)
	for id, entry := range b.inFlight {
		if !entry.expiry.After(t) {
			expiredInFlight = append(expiredInFlight, entry)
			delete(b.inFlight, id)
		}
	}
	b.mu.Unlock()

	timeoutResp := &vaulttypes.Response{Error: context.DeadlineExceeded.Error()}
	for _, entry := range append(expiredPending, expiredInFlight...) {
		b.lggr.Debugw("fast-path request expired", "requestID", entry.requestID)
		resp := *timeoutResp
		resp.ID = entry.requestID
		select {
		case entry.responseChan <- &resp:
		default:
		}
	}
}

func (b *FastPathBuffer) drainAllWithError(errMsg string) {
	b.mu.Lock()
	entries := make([]*fastPathEntry, 0, len(b.pending)+len(b.inFlight))
	for _, entry := range b.pending {
		entries = append(entries, entry)
	}
	for _, entry := range b.inFlight {
		entries = append(entries, entry)
	}
	b.pending = make(map[string]*fastPathEntry)
	b.inFlight = make(map[string]*fastPathEntry)
	b.mu.Unlock()

	for _, entry := range entries {
		select {
		case entry.responseChan <- &vaulttypes.Response{ID: entry.requestID, Error: errMsg}:
		default:
		}
	}
}
