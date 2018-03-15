package performance_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"fmt"
	"time"
	"github.com/tsenart/vegeta/lib"
	"log"
	"github.com/smartcontractkit/chainlink/internal/performance"
)


type Requirement struct {
	requestsPerSecond int
	averageLatencyMax, percentileMax, runDurationMax time.Duration
}

// API: POST /v2/jobs
func BenchmarkCreateJobApiUnderLoad(b *testing.B) {
	requirements := []Requirement{
		{ requestsPerSecond: 20, averageLatencyMax: time.Millisecond * 100, percentileMax: time.Millisecond * 100 },
		{ requestsPerSecond: 40, averageLatencyMax: time.Millisecond * 200, percentileMax: time.Millisecond * 200 },
		{ requestsPerSecond: 60, averageLatencyMax: time.Millisecond * 300, percentileMax: time.Millisecond * 300 },
	}

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	targets := performance.GetCreateJobTargets(app)

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

// API: GET /v2/jobs/:id
func BenchmarkViewJobApiUnderLoad(b *testing.B) {
	requirements := []Requirement{
		{ requestsPerSecond: 20, averageLatencyMax: time.Millisecond * 100, percentileMax: time.Millisecond * 100 },
		{ requestsPerSecond: 40, averageLatencyMax: time.Millisecond * 200, percentileMax: time.Millisecond * 200 },
		{ requestsPerSecond: 60, averageLatencyMax: time.Millisecond * 300, percentileMax: time.Millisecond * 300 },
	}

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	createJobTargets := performance.GetCreateJobTargets(app)

	for _, requirement := range requirements {

		rate := uint64(requirement.requestsPerSecond)
		duration := 4 * time.Second

		targeter := vegeta.NewStaticTargeter(createJobTargets...)
		attacker := vegeta.NewAttacker()

		for range attacker.Attack(targeter, rate, duration){}

		viewJobTargets := performance.GetViewJobTargets(app)
		targeter = vegeta.NewStaticTargeter(viewJobTargets...)
		attacker = vegeta.NewAttacker()

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

// API: POST /v2/jobs/:id/runs
func BenchmarkCreateJobRunsApiUnderLoad(b *testing.B) {
	requirements := []Requirement{
		{
			requestsPerSecond: 20,
			averageLatencyMax: time.Millisecond * 100,
			percentileMax: time.Millisecond * 100,
			runDurationMax: time.Millisecond * 100,
		}, {
			requestsPerSecond: 40,
			averageLatencyMax: time.Millisecond * 200,
			percentileMax: time.Millisecond * 100,
			runDurationMax: time.Millisecond * 200,
		}, {
			requestsPerSecond: 60,
			averageLatencyMax: time.Millisecond * 300,
			percentileMax: time.Millisecond * 100,
			runDurationMax: time.Millisecond * 300,
		},
	}

	for _, requirement := range requirements {
		app, cleanup := cltest.NewApplicationWithKeyStore()

		rate := uint64(requirement.requestsPerSecond)
		duration := 4 * time.Second

		createJobTargets := performance.GetCreateJobTargets(app)
		targeter := vegeta.NewStaticTargeter(createJobTargets...)
		attacker := vegeta.NewAttacker()

		for range attacker.Attack(targeter, rate, duration){}

		jobRunTargets := performance.GetJobRunTargets(app)
		targeter = vegeta.NewStaticTargeter(jobRunTargets...)
		attacker = vegeta.NewAttacker()

		var metrics vegeta.Metrics
		for res := range attacker.Attack(targeter, rate, duration) {
			metrics.Add(res)
		}
		metrics.Close()

		log.Printf("Average Latency: %s\n", metrics.Latencies.Mean)
		log.Printf("95th percentile: %s\n", metrics.Latencies.P95)

		averageLatency := performance.CalculateAverageJobRunLatency(app)
		log.Printf("Time to Complete Runs after API Calls: %s\n", averageLatency)

		assert.True(
			b,
			requirement.runDurationMax > averageLatency,
			fmt.Sprintf(
				"Time for running tasks to complete was higher than the requirement of %s at %d",
				requirement.runDurationMax,
				requirement.requestsPerSecond),
		)
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

		cleanup()
	}
}

// API: GET /v2/jobs/:id/runs
func BenchmarkViewJobRunsApiUnderLoad(b *testing.B) {
	requirements := []Requirement{
		{ requestsPerSecond: 20, averageLatencyMax: time.Millisecond * 100, percentileMax: time.Millisecond * 100 },
		{ requestsPerSecond: 40, averageLatencyMax: time.Millisecond * 200, percentileMax: time.Millisecond * 200 },
		{ requestsPerSecond: 60, averageLatencyMax: time.Millisecond * 300, percentileMax: time.Millisecond * 300 },
	}

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	createJobTargets := performance.GetCreateJobTargets(app)

	for _, requirement := range requirements {
		rate := uint64(requirement.requestsPerSecond) // per second
		duration := 4 * time.Second

		targeter := vegeta.NewStaticTargeter(createJobTargets...)
		attacker := vegeta.NewAttacker()

		for range attacker.Attack(targeter, rate, duration){}

		jobRunTargets := performance.GetJobRunTargets(app)
		targeter = vegeta.NewStaticTargeter(jobRunTargets...)
		attacker = vegeta.NewAttacker()

		for range attacker.Attack(targeter, rate, duration){}
		performance.WaitForJobRunsToComplete(app)

		viewJobRunTargets := performance.GetViewJobRunTargets(app)
		targeter = vegeta.NewStaticTargeter(viewJobRunTargets...)
		attacker = vegeta.NewAttacker()

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