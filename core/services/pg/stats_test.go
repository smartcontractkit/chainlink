package pg

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/logger"
)

// testDbStater is a simple test wrapper for statFn
type testDbStater struct {
	t         *testing.T
	name      string
	cntr      int64
	testGauge prometheus.Gauge
}

func newtestDbStater(t *testing.T, name string) *testDbStater {
	return &testDbStater{
		t:    t,
		name: name,
		testGauge: promauto.NewGauge(prometheus.GaugeOpts{
			Name: strings.ReplaceAll(name, " ", "_"),
		}),
	}
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

func (s *testDbStater) report(stats sql.DBStats) {
	s.t.Logf("reporting stats +%v", stats)
	s.testGauge.Set(float64(stats.MaxOpenConnections))
}

func (s *testDbStater) checkReport() {
	assert.Greater(s.t, testutil.ToFloat64(s.testGauge), float64(0), s.name)
}

type statScenario struct {
	name   string
	testFn func(*testing.T, *StatsReporter, time.Duration, int)
}

func TestStatReporter(t *testing.T) {
	interval := 2 * time.Millisecond
	expectedIntervals := 4

	lggr := logger.TestLogger(t)

	for _, scenario := range []statScenario{
		{name: "parent_ctx_canceled", testFn: testParentContextCanceled},
		{name: "normal_collect_and_stop", testFn: testCollectAndStop},
		{name: "mutli_start", testFn: testMultiStart},
		{name: "multi_stop", testFn: testMultiStop},
	} {

		d := newtestDbStater(t, scenario.name)
		reporter := NewStatsReporter(d.Stats,
			lggr,
			StatsInterval(interval),
			StatsCustomReporterFn(d.report),
		)

		scenario.testFn(
			t,
			reporter,
			interval,
			expectedIntervals,
		)

		d.checkReport()

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
