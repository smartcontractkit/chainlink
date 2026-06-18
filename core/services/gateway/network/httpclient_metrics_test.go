package network

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewHTTPClientMetrics(t *testing.T) {
	m, err := newHTTPClientMetrics()
	require.NoError(t, err)
	require.NotNil(t, m)
}

func TestHTTPClientMetrics_RecordPhase(t *testing.T) {
	m, err := newHTTPClientMetrics()
	require.NoError(t, err)

	ctx := context.Background()
	for _, phase := range []TracePhase{
		PhaseGetConn,
		PhaseDNSLookup,
		PhaseTCPConnect,
		PhaseTLSHandshake,
		PhaseWroteRequest,
		PhaseTimeToFirstByte,
		PhaseTotal,
	} {
		m.recordPhase(ctx, "GET", phase, 25*time.Millisecond)
	}
}

func TestHTTPClientMetrics_RecordTotal(t *testing.T) {
	m, err := newHTTPClientMetrics()
	require.NoError(t, err)

	ctx := context.Background()
	m.recordTotal(ctx, "GET", 200, true, false, 25*time.Millisecond)
	m.recordTotal(ctx, "POST", 500, true, true, 100*time.Millisecond)
	m.recordTotal(ctx, "GET", 0, false, false, 5*time.Millisecond)
}

func TestHTTPClientMetricViews(t *testing.T) {
	views := HTTPClientMetricViews()
	require.Len(t, views, 1)
}
