package pg

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

// testDbStater is a simple test implementation of the DBStater interface
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

type statScneario func(*testing.T, *StatsReporter, time.Duration, int)

func TestStatReporter(t *testing.T) {
	interval := 100 * time.Microsecond
	expectedIntervals := 7

	for _, scenario := range []statScneario{
		testParentContextCanceled,
		testCollectAndStop,
		testMultiStart,
		testMultiStop} {
		scenario(
			t,
			NewStatsReporter(new(testDbStater), StatsInterval(interval)),
			interval,
			expectedIntervals,
		)

		// if shutdown of the reporter was broken this sleep would afford time
		// for unwanted, post-stop, metrics
		time.Sleep(2 * interval)
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

	r.Start(ctx)
	time.Sleep(interval)
	r.Start(ctx)
	time.Sleep(time.Duration(n-1) * interval)
	r.Stop()
}

// test multiple stop calls are idempotent
func testMultiStop(t *testing.T, r *StatsReporter, interval time.Duration, n int) {

	ctx := context.Background()

	r.Start(ctx)
	time.Sleep(time.Duration(n) * interval)
	r.Stop()
	r.Stop()
}

func assertStats(t *testing.T, expected int) {
	// use in delta because counters inside  go routines are inherently fuzzy
	assert.InDelta(t, expected, testutil.ToFloat64(promDBConnsInUse), 1)
	assert.InDelta(t, expected, testutil.ToFloat64(promDBConnsMax), 1)
	assert.InDelta(t, expected, testutil.ToFloat64(promDBConnsOpen), 1)
	assert.InDelta(t, expected, testutil.ToFloat64(promDBWaitCount), 1)
	assert.InDelta(t, expected, testutil.ToFloat64(promDBWaitDuration), 1)
}
