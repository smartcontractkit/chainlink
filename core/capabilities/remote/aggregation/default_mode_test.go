package aggregation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
)

const (
	testPayloadSize = 2000
	testDonSize     = 7
	testMinResp     = uint32(3)
)

type aggregationScenario struct {
	name          string
	payloads      [][]byte
	minResponses  uint32
	expectPayload []byte
	expectError   bool
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

func Benchmark_AggregateModeRaw(b *testing.B) {
	for _, sc := range aggregationScenarios(b) {
		if sc.expectError {
			// not interested in breaking scenarios for benchmarks, only valid ones
			continue
		}
		b.Run(sc.name, func(b *testing.B) {
			for b.Loop() {
				_, _ = AggregateModeRaw(sc.payloads, sc.minResponses)
			}
		})
	}
}

func testPayload(tb testing.TB, data string, eventID string) []byte {
	tb.Helper()
	val, err := values.NewMap(map[string]any{"event": data})
	require.NoError(tb, err)
	resp := commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{ID: eventID, Outputs: val},
	}
	m, err := pb.MarshalTriggerResponse(resp)
	require.NoError(tb, err)
	return m
}

func testDataString(size int, seed int) string {
	buf := make([]byte, size)
	for i := range buf {
		buf[i] = 'A' + byte((seed+i)%26)
	}
	return string(buf)
}

// aggregationScenarios returns the shared set of test cases used by both correctness tests and benchmarks.
// If scenarios are added or modified, regenerate the benchmark baseline:
//
//	cd ./core/capabilities/remote/aggregation/
//	go test -bench=Benchmark_AggregateModeRaw -benchmem -count=6 . > testdata/aggregate_mode_raw_benchmark.txt
func aggregationScenarios(tb testing.TB) []aggregationScenario {
	tb.Helper()

	pA := testPayload(tb, testDataString(testPayloadSize, 0), "event-A")
	pB := testPayload(tb, testDataString(testPayloadSize, 13), "event-B")
	pC := testPayload(tb, testDataString(testPayloadSize, 7), "event-C")

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
		},
		{
			name:          "3 identical payloads - returns A",
			payloads:      [][]byte{pA, pA, pA},
			minResponses:  testMinResp,
			expectPayload: pA,
		},
		{
			name: fmt.Sprintf("%d payloads: 3×A + %d different - returns A", testDonSize, testDonSize-2),
			payloads: func() [][]byte {
				p := make([][]byte, testDonSize)
				p[0] = pA
				p[1] = pA
				p[2] = pA
				for i := 3; i < testDonSize; i++ {
					p[i] = testPayload(tb, testDataString(testPayloadSize, i*7), fmt.Sprintf("diff-%d", i))
				}
				return p
			}(),
			minResponses:  testMinResp,
			expectPayload: pA,
		},
		{
			name:          fmt.Sprintf("%d payloads: %d×A + %d×B - returns B (mode)", testDonSize, testDonSize/2, testDonSize-testDonSize/2),
			payloads:      nMixed(testDonSize/2, testDonSize-testDonSize/2),
			minResponses:  testMinResp,
			expectPayload: pB,
		},
		{
			name:          fmt.Sprintf("%d payloads: %d×A + %d×B - returns A (mode)", testDonSize, testDonSize-testDonSize/2, testDonSize/2),
			payloads:      nMixed(testDonSize-testDonSize/2, testDonSize/2),
			minResponses:  testMinResp,
			expectPayload: pA,
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
		},
		{
			name:          fmt.Sprintf("%d payloads: [A A B B C C C] - returns C (mode)", testDonSize),
			payloads:      [][]byte{pA, pA, pB, pB, pC, pC, pC},
			minResponses:  testMinResp,
			expectPayload: pC,
		},
		{
			name:          fmt.Sprintf("%d payloads: [B B A A C C C] - returns C (mode)", testDonSize),
			payloads:      [][]byte{pB, pB, pA, pA, pC, pC, pC},
			minResponses:  testMinResp,
			expectPayload: pC,
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
