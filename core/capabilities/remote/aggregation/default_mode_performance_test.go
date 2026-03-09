package aggregation

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
)

func perfPayload(t *testing.T, data string, eventID string) []byte {
	t.Helper()
	val, err := values.NewMap(map[string]any{"event": data})
	require.NoError(t, err)
	resp := commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{ID: eventID, Outputs: val},
	}
	m, err := pb.MarshalTriggerResponse(resp)
	require.NoError(t, err)
	return m
}

func perfDataString(size int, seed int) string {
	buf := make([]byte, size)
	for i := range buf {
		buf[i] = 'A' + byte((seed+i)%26)
	}
	return string(buf)
}

// TestAggregateModeRawFast_Performance asserts that AggregateModeRawFast stays
// within expected latency budgets. If a future change causes a regression,
// these thresholds will catch it.
func TestAggregateModeRawFast_Performance(t *testing.T) {
	const (
		payloadSize  = 2000
		donSize      = 7
		minResponses = 2
		iterations   = 100_000
	)

	payloadA := perfPayload(t, perfDataString(payloadSize, 0), "event-A")
	payloadB := perfPayload(t, perfDataString(payloadSize, 13), "event-B")

	scenarios := []struct {
		name     string
		build    func() [][]byte
		maxAvgUs float64 // max allowed average latency in microseconds
	}{
		{
			name: fmt.Sprintf("%d identical payloads (production fast path)", donSize),
			build: func() [][]byte {
				p := make([][]byte, donSize)
				for i := range p {
					p[i] = payloadA
				}
				return p
			},
			maxAvgUs: 0.5, // fast path: bytes.Equal only, sub-microsecond
		},
		{
			name: "2 identical payloads (minimum fast path)",
			build: func() [][]byte {
				return [][]byte{payloadA, payloadA}
			},
			maxAvgUs: 0.5,
		},
		{
			name: fmt.Sprintf("%d payloads: 2 identical + %d different (all different between each other) (hash fallback)", donSize, donSize-2),
			build: func() [][]byte {
				p := make([][]byte, donSize)
				p[0] = payloadA
				p[1] = payloadA
				for i := 2; i < donSize; i++ {
					p[i] = perfPayload(t, perfDataString(payloadSize, i*7), fmt.Sprintf("diff-%d", i))
				}
				return p
			},
			maxAvgUs: 10.0, // hash fallback with donSize SHA256 operations
		},
		{
			name: fmt.Sprintf("%d payloads: %d×A + %d×B, mode is B (hash fallback)", donSize, donSize/2, donSize-donSize/2),
			build: func() [][]byte {
				p := make([][]byte, donSize)
				for i := 0; i < donSize/2; i++ {
					p[i] = payloadA
				}
				for i := donSize / 2; i < donSize; i++ {
					p[i] = payloadB
				}
				return p
			},
			maxAvgUs: 10.0,
		},
		{
			name: fmt.Sprintf("%d payloads: %d×A + %d×B, mode is A (hash fallback)", donSize, donSize-donSize/2, donSize/2),
			build: func() [][]byte {
				p := make([][]byte, donSize)
				for i := 0; i < donSize-donSize/2; i++ {
					p[i] = payloadA
				}
				for i := donSize - donSize/2; i < donSize; i++ {
					p[i] = payloadB
				}
				return p
			},
			maxAvgUs: 10.0,
		},
		{
			name: fmt.Sprintf("%d payloads: %d×A + 1×B, mode is A (hash fallback)", donSize, donSize-1),
			build: func() [][]byte {
				p := make([][]byte, donSize)
				for i := 0; i < donSize-1; i++ {
					p[i] = payloadA
				}
				p[donSize-1] = payloadB
				return p
			},
			maxAvgUs: 10.0,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			payloads := sc.build()

			// Correctness: cross-check against original
			resultOld, err := AggregateModeRaw(payloads, minResponses)
			require.NoError(t, err)
			result, err := AggregateModeRawFast(payloads, minResponses)
			require.NoError(t, err)
			require.Equal(t, resultOld, result, "AggregateModeRawFast must match AggregateModeRaw")

			// Warm up
			for range 5000 {
				_, _ = AggregateModeRawFast(payloads, minResponses)
			}

			// Measure
			start := time.Now()
			for range iterations {
				_, _ = AggregateModeRawFast(payloads, minResponses)
			}
			elapsed := time.Since(start)
			avgUs := float64(elapsed.Nanoseconds()) / float64(iterations) / 1e3

			t.Logf("  Payloads: %d × %d bytes | Iterations: %d", len(payloads), len(payloads[0]), iterations)
			t.Logf("  Avg latency: %.3f µs (threshold: %.1f µs)", avgUs, sc.maxAvgUs)

			assert.LessOrEqual(t, avgUs, sc.maxAvgUs,
				"AggregateModeRawFast exceeded latency budget: got %.3f µs, max %.1f µs", avgUs, sc.maxAvgUs)
		})
	}
}

// TestAggregateModeRaw_Performance benchmarks the original AggregateModeRaw
// using the same scenarios as TestAggregateModeRawFast_Performance.
func TestAggregateModeRaw_Performance(t *testing.T) {
	const (
		payloadSize  = 2000
		donSize      = 7
		minResponses = 2
		iterations   = 100_000
	)

	payloadA := perfPayload(t, perfDataString(payloadSize, 0), "event-A")
	payloadB := perfPayload(t, perfDataString(payloadSize, 13), "event-B")

	scenarios := []struct {
		name  string
		build func() [][]byte
	}{
		{
			name: fmt.Sprintf("%d identical payloads (production case)", donSize),
			build: func() [][]byte {
				p := make([][]byte, donSize)
				for i := range p {
					p[i] = payloadA
				}
				return p
			},
		},
		{
			name: "2 identical payloads (minimum)",
			build: func() [][]byte {
				return [][]byte{payloadA, payloadA}
			},
		},
		{
			name: fmt.Sprintf("%d payloads: 2 identical + %d different (all different between each other) (hash fallback)", donSize, donSize-2),
			build: func() [][]byte {
				p := make([][]byte, donSize)
				p[0] = payloadA
				p[1] = payloadA
				for i := 2; i < donSize; i++ {
					p[i] = perfPayload(t, perfDataString(payloadSize, i*7), fmt.Sprintf("diff-%d", i))
				}
				return p
			},
		},
		{
			name: fmt.Sprintf("%d payloads: %d×A + %d×B, mode is B (hash fallback)", donSize, donSize/2, donSize-donSize/2),
			build: func() [][]byte {
				p := make([][]byte, donSize)
				for i := 0; i < donSize/2; i++ {
					p[i] = payloadA
				}
				for i := donSize / 2; i < donSize; i++ {
					p[i] = payloadB
				}
				return p
			},
		},
		{
			name: fmt.Sprintf("%d payloads: %d×A + %d×B, mode is A (hash fallback)", donSize, donSize-donSize/2, donSize/2),
			build: func() [][]byte {
				p := make([][]byte, donSize)
				for i := 0; i < donSize-donSize/2; i++ {
					p[i] = payloadA
				}
				for i := donSize - donSize/2; i < donSize; i++ {
					p[i] = payloadB
				}
				return p
			},
		},
		{
			name: fmt.Sprintf("%d payloads: %d×A + 1×B, mode is A (hash fallback)", donSize, donSize-1),
			build: func() [][]byte {
				p := make([][]byte, donSize)
				for i := 0; i < donSize-1; i++ {
					p[i] = payloadA
				}
				p[donSize-1] = payloadB
				return p
			},
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			payloads := sc.build()

			result, err := AggregateModeRaw(payloads, minResponses)
			require.NoError(t, err)
			require.NotNil(t, result)

			// Warm up
			for range 5000 {
				_, _ = AggregateModeRaw(payloads, minResponses)
			}

			// Measure
			start := time.Now()
			for range iterations {
				_, _ = AggregateModeRaw(payloads, minResponses)
			}
			elapsed := time.Since(start)
			avgUs := float64(elapsed.Nanoseconds()) / float64(iterations) / 1e3

			t.Logf("  Payloads: %d × %d bytes | Iterations: %d", len(payloads), len(payloads[0]), iterations)
			t.Logf("  Avg latency: %.3f µs", avgUs)
		})
	}
}
