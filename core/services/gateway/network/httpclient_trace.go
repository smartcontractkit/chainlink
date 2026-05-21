package network

import (
	"context"
	"crypto/tls"
	"net/http/httptrace"
	"sync/atomic"
	"time"
)

// requestTrace holds out-of-band signals captured during a single HTTP request
// so the caller can include them in the final request metric.
type requestTrace struct {
	connReused atomic.Bool
}

// newClientTrace returns an httptrace.ClientTrace that records per-phase
// durations on m, along with a requestTrace exposing connection reuse.

// Metrics recorded from paired start/end callbacks (GetConn, DNS, TCP connect,
// TLS handshake) are phase-local durations. WroteRequest and
// GotFirstResponseByte record cumulative elapsed time since requestStart.
//
// httptrace may invoke callbacks concurrently so each phase keeps its own atomic start timestamp.
func newClientTrace(ctx context.Context, method string, requestStart time.Time, m *httpClientMetrics) (*httptrace.ClientTrace, *requestTrace) {
	rt := &requestTrace{}

	var (
		getConnStart atomic.Pointer[time.Time]
		dnsStart     atomic.Pointer[time.Time]
		connectStart atomic.Pointer[time.Time]
		tlsStart     atomic.Pointer[time.Time]
	)

	storeNow := func(p *atomic.Pointer[time.Time]) {
		now := time.Now()
		p.Store(&now)
	}

	recordPhase := func(start *atomic.Pointer[time.Time], phase TracePhase) {
		s := start.Load()
		if s == nil {
			return
		}
		m.recordPhase(ctx, method, phase, time.Since(*s))
	}

	return &httptrace.ClientTrace{
		GetConn: func(string) {
			storeNow(&getConnStart)
		},
		GotConn: func(info httptrace.GotConnInfo) {
			recordPhase(&getConnStart, PhaseGetConn)
			if info.Reused {
				rt.connReused.Store(true)
			}
		},
		DNSStart: func(httptrace.DNSStartInfo) {
			storeNow(&dnsStart)
		},
		DNSDone: func(httptrace.DNSDoneInfo) {
			recordPhase(&dnsStart, PhaseDNSLookup)
		},
		ConnectStart: func(string, string) {
			storeNow(&connectStart)
		},
		ConnectDone: func(string, string, error) {
			recordPhase(&connectStart, PhaseTCPConnect)
		},
		TLSHandshakeStart: func() {
			storeNow(&tlsStart)
		},
		TLSHandshakeDone: func(tls.ConnectionState, error) {
			recordPhase(&tlsStart, PhaseTLSHandshake)
		},
		WroteRequest: func(httptrace.WroteRequestInfo) {
			m.recordPhase(ctx, method, PhaseWroteRequest, time.Since(requestStart))
		},
		GotFirstResponseByte: func() {
			m.recordPhase(ctx, method, PhaseTimeToFirstByte, time.Since(requestStart))
		},
	}, rt
}
