package performance_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/internal/performance"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
	"github.com/tsenart/vegeta/lib"
)

type Requirement struct {
	requestsPerSecond int
	averageLatencyMax time.Duration
	percentileMax     time.Duration
	runDurationMax    time.Duration
}

var (
	r            Requirement
	requirements = []Requirement{
		{20, time.Millisecond * 100, time.Millisecond * 100, time.Millisecond * 100},
		{40, time.Millisecond * 200, time.Millisecond * 200, time.Millisecond * 200},
		{60, time.Millisecond * 300, time.Millisecond * 300, time.Millisecond * 300},
	}

	app     *cltest.TestApplication
	cleanup func()
)

func Setup() func(b *testing.B) {
	 return func(b *testing.B) {
		for _, req := range requirements {
			app, cleanup = cltest.NewApplicationWithKeyStore()
			r = req
			b.Run("CreateJobs", CreateJobApiUnderLoad)
			b.Run("ViewJobs", ViewJobApiUnderLoad)
			b.Run("CreateJobRuns", CreateJobRunsApiUnderLoad)
			b.Run("ViewJobRuns", ViewJobRunsApiUnderLoad)
			cleanup()
		}
	}
}

func BenchmarkNodeUnderLoad(b *testing.B) {
	tests := Setup()
	tests(b)
}

func CreateJobApiUnderLoad(b *testing.B) {
	targets := performance.GetPostJobTargets(app)

	rate := uint64(r.requestsPerSecond) // per second
	duration := 4 * time.Second
	targeter := vegeta.NewStaticTargeter(targets...)
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration) {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("Average Latency: %s\n", metrics.Latencies.Mean)
	fmt.Printf("95th percentile: %s\n", metrics.Latencies.P95)

	assert.Condition(
		b,
		func() bool { return r.averageLatencyMax > metrics.Latencies.Mean },
		fmt.Sprintf(
			"Average latency was higher than the requirement of %s at %d",
			r.averageLatencyMax,
			r.requestsPerSecond),
	)
	assert.Condition(
		b,
		func() bool { return r.percentileMax > metrics.Latencies.P95 },
		fmt.Sprintf(
			"95th percentile was higher than the requirement of %s at %d",
			r.percentileMax,
			r.requestsPerSecond),
	)
}

func ViewJobApiUnderLoad(b *testing.B) {
	targets := performance.GetViewJobTargets(app)

	rate := uint64(r.requestsPerSecond) // per second
	duration := 4 * time.Second
	targeter := vegeta.NewStaticTargeter(targets...)
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration) {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("Average Latency: %s\n", metrics.Latencies.Mean)
	fmt.Printf("95th percentile: %s\n", metrics.Latencies.P95)

	assert.Condition(
		b,
		func() bool { return r.averageLatencyMax > metrics.Latencies.Mean },
		fmt.Sprintf(
			"Average latency was higher than the requirement of %s at %d",
			r.averageLatencyMax,
			r.requestsPerSecond),
	)
	assert.Condition(
		b,
		func() bool { return r.percentileMax > metrics.Latencies.P95 },
		fmt.Sprintf(
			"95th percentile was higher than the requirement of %s at %d",
			r.percentileMax,
			r.requestsPerSecond),
	)
}

func CreateJobRunsApiUnderLoad(b *testing.B) {
	targets := performance.GetJobRunTargets(app)

	rate := uint64(r.requestsPerSecond) // per second
	duration := 4 * time.Second
	targeter := vegeta.NewStaticTargeter(targets...)
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration) {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("Average Latency: %s\n", metrics.Latencies.Mean)
	fmt.Printf("95th percentile: %s\n", metrics.Latencies.P95)

	averageLatency := performance.CalculateAverageJobRunLatency(app)
	fmt.Printf("Time to Complete Runs after API Calls: %s\n", averageLatency)

	assert.Condition(
		b,
		func() bool { return r.runDurationMax > averageLatency },
		fmt.Sprintf(
			"Time for running tasks to complete was higher than the requirement of %s at %d",
			r.runDurationMax,
			r.requestsPerSecond),
	)
	assert.Condition(
		b,
		func() bool { return r.averageLatencyMax > metrics.Latencies.Mean },
		fmt.Sprintf(
			"Average latency was higher than the requirement of %s at %d",
			r.averageLatencyMax,
			r.requestsPerSecond),
	)
	assert.Condition(
		b,
		func() bool { return r.percentileMax > metrics.Latencies.P95 },
		fmt.Sprintf(
			"95th percentile was higher than the requirement of %s at %d",
			r.percentileMax,
			r.requestsPerSecond),
	)
}

func ViewJobRunsApiUnderLoad(b *testing.B) {
	targets := performance.GetViewJobRunTargets(app)

	rate := uint64(r.requestsPerSecond) // per second
	duration := 4 * time.Second
	targeter := vegeta.NewStaticTargeter(targets...)
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration) {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("Average Latency: %s\n", metrics.Latencies.Mean)
	fmt.Printf("95th percentile: %s\n", metrics.Latencies.P95)

	assert.Condition(
		b,
		func() bool { return r.averageLatencyMax > metrics.Latencies.Mean },
		fmt.Sprintf(
			"Average latency was higher than the requirement of %s at %d",
			r.averageLatencyMax,
			r.requestsPerSecond),
	)
	assert.Condition(
		b,
		func() bool { return r.percentileMax > metrics.Latencies.P95 },
		fmt.Sprintf(
			"95th percentile was higher than the requirement of %s at %d",
			r.percentileMax,
			r.requestsPerSecond),
	)
}