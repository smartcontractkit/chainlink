package pg

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/stretchr/testify/assert"
)

// testDbStater is a simple test wrapper for statFn
type testDbStater struct {
	cntr int64
}

func (s *testDbStater) Stats() sql.DBStats {
	s.cntr++
	return sql.DBStats{
		MaxOpenConnections: int(s.cntr),
		OpenConnections:    int(s.cntr),
		InUse:              int(s.cntr),
		Idle:               int(s.cntr),
		WaitCount:          s.cntr,
		WaitDuration:       time.Duration(s.cntr * int64(time.Second)),
		MaxIdleClosed:      s.cntr,
		MaxLifetimeClosed:  s.cntr,
	}

}

type statScenario func(*testing.T, *StatsReporter, time.Duration, int)

func TestStatReporter(t *testing.T) {
	interval := 2 * time.Millisecond
	expectedIntervals := 4

	for _, scenario := range []statScenario{
		testParentContextCanceled,
		testCollectAndStop,
		testMultiStart,
		testMultiStop} {

		d := new(testDbStater)
		resetProm(t)
		scenario(
			t,
			NewStatsReporter(d.Stats, StatsInterval(interval)),
			interval,
			expectedIntervals,
		)

		assertStats(t, expectedIntervals)

	}

}

// test appropriate handling of context cancellation
func testParentContextCanceled(t *testing.T, r *StatsReporter, interval time.Duration, n int) {

	ctx := context.Background()
	tctx, cancel := context.WithTimeout(ctx, time.Duration(n)*interval)

	r.Start(tctx)
	// wait for parent cancelation
	<-tctx.Done()
	// call cancel to statisy linter
	cancel()

}

// test normal stop
func testCollectAndStop(t *testing.T, r *StatsReporter, interval time.Duration, n int) {

	ctx := context.Background()

	r.Start(ctx)
	time.Sleep(time.Duration(n) * interval)
	r.Stop()

}

// test multiple start calls are idempotent
func testMultiStart(t *testing.T, r *StatsReporter, interval time.Duration, n int) {

	ctx := context.Background()

	ticker := time.NewTicker(time.Duration(n) * interval)
	defer ticker.Stop()

	r.Start(ctx)
	r.Start(ctx)
	<-ticker.C
	r.Stop()
}

// test multiple stop calls are idempotent
func testMultiStop(t *testing.T, r *StatsReporter, interval time.Duration, n int) {

	ctx := context.Background()

	ticker := time.NewTicker(time.Duration(n) * interval)
	defer ticker.Stop()

	r.Start(ctx)
	<-ticker.C
	r.Stop()
	r.Stop()
}

func assertStats(t *testing.T, expected int) {
	statInRange := func(stat float64) bool {
		return int(stat) > expected/2 && int(stat) <= expected
	}

	testutils.AssertEventually(t,
		func() bool { return statInRange(testutil.ToFloat64(promDBConnsInUse)) })

	testutils.AssertEventually(t,
		func() bool { return statInRange(testutil.ToFloat64(promDBConnsMax)) })

	testutils.AssertEventually(t,
		func() bool { return statInRange(testutil.ToFloat64(promDBConnsOpen)) })

	testutils.AssertEventually(t,
		func() bool { return statInRange(testutil.ToFloat64(promDBWaitCount)) })

	testutils.AssertEventually(t,
		func() bool { return statInRange(testutil.ToFloat64(promDBWaitDuration)) })
}

func resetProm(t *testing.T) {
	promDBConnsInUse.Set(0)
	assert.Equal(t, int(testutil.ToFloat64(promDBConnsInUse)), 0)

	promDBConnsMax.Set(0)
	assert.Equal(t, int(testutil.ToFloat64(promDBConnsMax)), 0)

	promDBConnsOpen.Set(0)
	assert.Equal(t, int(testutil.ToFloat64(promDBConnsOpen)), 0)

	promDBWaitCount.Set(0)
	assert.Equal(t, int(testutil.ToFloat64(promDBWaitCount)), 0)

	promDBWaitDuration.Set(0)
	assert.Equal(t, int(testutil.ToFloat64(promDBWaitDuration)), 0)
}
