package performance

import (
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
	"github.com/tsenart/vegeta/lib"
	"log"
)

// API: POST /v2/jobs
func BenchmarkCreateJobApiUnderLoad(b *testing.B) {
	requirements := []Requirement{
		{ requestsPerSecond: 20, averageLatencyMax: time.Millisecond * 100, percentileMax: time.Millisecond * 100 },
		{ requestsPerSecond: 40, averageLatencyMax: time.Millisecond * 200, percentileMax: time.Millisecond * 200 },
		{ requestsPerSecond: 60, averageLatencyMax: time.Millisecond * 300, percentileMax: time.Millisecond * 300 },
	}

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	targets := GetCreateJobTargets(app)

	for _, requirement := range requirements {
		rate := uint64(requirement.requestsPerSecond) // per second
		duration := 4 * time.Second
		targeter := vegeta.NewStaticTargeter(targets...)
		attacker := vegeta.NewAttacker()

		var metrics vegeta.Metrics
		for res := range attacker.Attack(targeter, rate, duration) {
			metrics.Add(res)
		}
		metrics.Close()

		log.Printf("Average Latency: %s\n", metrics.Latencies.Mean)
		log.Printf("95th percentile: %s\n", metrics.Latencies.P95)

		assert.True(
			b,
			requirement.averageLatencyMax > metrics.Latencies.Mean,
			fmt.Sprintf(
				"Average latency was higher than the requirement of %s at %d",
				requirement.averageLatencyMax,
				requirement.requestsPerSecond),
		)
		assert.True(
			b,
			requirement.percentileMax > metrics.Latencies.P95,
			fmt.Sprintf(
				"95th percentile was higher than the requirement of %s at %d",
				requirement.percentileMax,
				requirement.requestsPerSecond),
		)
	}
}