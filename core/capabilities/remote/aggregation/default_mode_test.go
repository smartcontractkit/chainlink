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

const (
	testPayloadSize = 2000
	testDonSize     = 7
	testMinResp     = uint32(2)
)

type aggregationScenario struct {
	name          string
	payloads      [][]byte
	minResponses  uint32
	expectPayload []byte
	expectError   bool
	maxAvgUs      float64 // latency budget for performance tests (0 = skip perf)
}

var (
	triggerEvent1 = map[string]any{"event": "triggerEvent1"}
	triggerEvent2 = map[string]any{"event": "triggerEvent2"}
)

func TestDefaultModeAggregator_Aggregate(t *testing.T) {
	val, err := values.NewMap(triggerEvent1)
	require.NoError(t, err)
	capResponse1 := commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{
			Outputs: val,
		},
		Err: nil,
	}
	marshaled1, err := pb.MarshalTriggerResponse(capResponse1)
	require.NoError(t, err)

	val2, err := values.NewMap(triggerEvent2)
	require.NoError(t, err)
	capResponse2 := commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{
			Outputs: val2,
		},
		Err: nil,
	}
	marshaled2, err := pb.MarshalTriggerResponse(capResponse2)
	require.NoError(t, err)

	agg := NewDefaultModeAggregator(2)
	_, err = agg.Aggregate("", [][]byte{marshaled1})
	require.Error(t, err)

	_, err = agg.Aggregate("", [][]byte{marshaled1, marshaled2})
	require.Error(t, err)

	res, err := agg.Aggregate("", [][]byte{marshaled1, marshaled2, marshaled1})
	require.NoError(t, err)
	require.Equal(t, res, capResponse1)
}

func TestAggregateModeRaw_Correctness(t *testing.T) {
	for _, sc := range aggregationScenarios(t) {
		t.Run(sc.name, func(t *testing.T) {
			result, err := AggregateModeRaw(sc.payloads, sc.minResponses)
			if sc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, sc.expectPayload, result)
		})
	}
}

func TestAggregateModeRaw_Performance(t *testing.T) {
	const iterations = 100_000

	for _, sc := range aggregationScenarios(t) {
		if sc.expectError || sc.maxAvgUs == 0 {
			continue
		}
		t.Run(sc.name, func(t *testing.T) {
			for range 5000 {
				_, _ = AggregateModeRaw(sc.payloads, sc.minResponses)
			}

			start := time.Now()
			for range iterations {
				_, _ = AggregateModeRaw(sc.payloads, sc.minResponses)
			}
			elapsed := time.Since(start)
			avgUs := float64(elapsed.Nanoseconds()) / float64(iterations) / 1e3

			t.Logf("  Payloads: %d × %d bytes | Iterations: %d", len(sc.payloads), len(sc.payloads[0]), iterations)
			t.Logf("  Avg latency: %.3f µs (threshold: %.1f µs)", avgUs, sc.maxAvgUs)

			assert.LessOrEqual(t, avgUs, sc.maxAvgUs,
				"AggregateModeRaw exceeded latency budget: got %.3f µs, max %.1f µs", avgUs, sc.maxAvgUs)
		})
	}
}

func testPayload(t *testing.T, data string, eventID string) []byte {
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

func testDataString(size int, seed int) string {
	buf := make([]byte, size)
	for i := range buf {
		buf[i] = 'A' + byte((seed+i)%26)
	}
	return string(buf)
}

func aggregationScenarios(t *testing.T) []aggregationScenario {
	t.Helper()

	pA := testPayload(t, testDataString(testPayloadSize, 0), "event-A")
	pB := testPayload(t, testDataString(testPayloadSize, 13), "event-B")
	pC := testPayload(t, testDataString(testPayloadSize, 7), "event-C")

	nAllTheSame := func() [][]byte {
		p := make([][]byte, testDonSize)
		for i := range p {
			p[i] = pA
		}
		return p
	}
	nMixed := func(countA, countB int) [][]byte {
		p := make([][]byte, countA+countB)
		for i := 0; i < countA; i++ {
			p[i] = pA
		}
		for i := countA; i < countA+countB; i++ {
			p[i] = pB
		}
		return p
	}

	return []aggregationScenario{
		{
			name:          fmt.Sprintf("%d identical payloads - returns A", testDonSize),
			payloads:      nAllTheSame(),
			minResponses:  testMinResp,
			expectPayload: pA,
			maxAvgUs:      0.5,
		},
		{
			name:          "2 identical payloads - returns A",
			payloads:      [][]byte{pA, pA},
			minResponses:  testMinResp,
			expectPayload: pA,
			maxAvgUs:      0.5,
		},
		{
			name: fmt.Sprintf("%d payloads: 2×A + %d different - returns A", testDonSize, testDonSize-2),
			payloads: func() [][]byte {
				p := make([][]byte, testDonSize)
				p[0] = pA
				p[1] = pA
				for i := 2; i < testDonSize; i++ {
					p[i] = testPayload(t, testDataString(testPayloadSize, i*7), fmt.Sprintf("diff-%d", i))
				}
				return p
			}(),
			minResponses:  testMinResp,
			expectPayload: pA,
			maxAvgUs:      10.0,
		},
		{
			name:          fmt.Sprintf("%d payloads: %d×A + %d×B - returns B (mode)", testDonSize, testDonSize/2, testDonSize-testDonSize/2),
			payloads:      nMixed(testDonSize/2, testDonSize-testDonSize/2),
			minResponses:  testMinResp,
			expectPayload: pB,
			maxAvgUs:      10.0,
		},
		{
			name:          fmt.Sprintf("%d payloads: %d×A + %d×B - returns A (mode)", testDonSize, testDonSize-testDonSize/2, testDonSize/2),
			payloads:      nMixed(testDonSize-testDonSize/2, testDonSize/2),
			minResponses:  testMinResp,
			expectPayload: pA,
			maxAvgUs:      10.0,
		},
		{
			name: fmt.Sprintf("%d payloads: %d×A + 1×B - returns A (mode)", testDonSize, testDonSize-1),
			payloads: func() [][]byte {
				p := make([][]byte, testDonSize)
				for i := 0; i < testDonSize-1; i++ {
					p[i] = pA
				}
				p[testDonSize-1] = pB
				return p
			}(),
			minResponses:  testMinResp,
			expectPayload: pA,
			maxAvgUs:      10.0,
		},
		{
			name:          fmt.Sprintf("%d payloads: [A A B B C C C] - returns C (mode)", testDonSize),
			payloads:      [][]byte{pA, pA, pB, pB, pC, pC, pC},
			minResponses:  testMinResp,
			expectPayload: pC,
			maxAvgUs:      10.0,
		},
		{
			name:          fmt.Sprintf("%d payloads: [B B A A C C C] - returns C (mode)", testDonSize),
			payloads:      [][]byte{pB, pB, pA, pA, pC, pC, pC},
			minResponses:  testMinResp,
			expectPayload: pC,
			maxAvgUs:      10.0,
		},
		{
			name:         "1 payload with minResponses=2 - error",
			payloads:     [][]byte{pA},
			minResponses: 2,
			expectError:  true,
		},
		{
			name:         "2 different payloads with minResponses=2 - error",
			payloads:     [][]byte{pA, pB},
			minResponses: 2,
			expectError:  true,
		},
	}
}