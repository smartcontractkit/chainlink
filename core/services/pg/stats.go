package pg

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const dbStatsInternal = 10 * time.Second

var (
	promDBConnsMax = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_conns_max",
		Help: "Maximum number of open connections to the database.",
	})
	promDBConnsOpen = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_conns_open",
		Help: "The number of established connections both in use and idle.",
	})
	promDBConnsInUse = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_conns_used",
		Help: "The number of connections currently in use.",
	})
	promDBWaitCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_wait_count",
		Help: "The total number of connections waited for.",
	})
	promDBWaitDuration = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_wait_time_seconds",
		Help: "The total time blocked waiting for a new connection.",
	})
)

func publishStats(stats sql.DBStats) {
	promDBConnsMax.Set(float64(stats.MaxOpenConnections))
	promDBConnsOpen.Set(float64(stats.OpenConnections))
	promDBConnsInUse.Set(float64(stats.InUse))

	promDBWaitCount.Set(float64(stats.WaitCount))
	promDBWaitDuration.Set(stats.WaitDuration.Seconds())
}

type StatsReporterOpt func(*StatsReporter)

func StatsInterval(d time.Duration) StatsReporterOpt {
	return func(r *StatsReporter) {
		r.interval = d
	}
}

func StatsCustomReporterFn(fn ReportFn) StatsReporterOpt {
	return func(r *StatsReporter) {
		r.reportFn = fn
	}
}

type (
	StatFn   func() sql.DBStats
	ReportFn func(sql.DBStats)
)

type StatsReporter struct {
	statFn   StatFn
	reportFn ReportFn
	interval time.Duration
	cancel   context.CancelFunc
	lggr     logger.Logger
	once     sync.Once
	wg       sync.WaitGroup
}

func NewStatsReporter(fn StatFn, lggr logger.Logger, opts ...StatsReporterOpt) *StatsReporter {
	r := &StatsReporter{
		statFn:   fn,
		reportFn: publishStats,
		interval: dbStatsInternal,
		lggr:     lggr.Named("StatsReporter"),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *StatsReporter) Start(ctx context.Context) {
	startOnce := func() {
		r.wg.Add(1)
		r.lggr.Debug("Starting DB stat reporter")
		rctx, cancelFunc := context.WithCancel(ctx)
		r.cancel = cancelFunc
		go r.loop(rctx)
	}

	r.once.Do(startOnce)
}

// Stop stops all resources owned by the reporter and waits
// for all of them to be done
func (r *StatsReporter) Stop() {
	if r.cancel != nil {
		r.lggr.Debug("Stopping DB stat reporter")
		r.cancel()
		r.cancel = nil
		r.wg.Wait()
	}
}

func (r *StatsReporter) loop(ctx context.Context) {
	defer r.wg.Done()

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	r.reportFn(r.statFn())
	for {
		select {
		case <-ticker.C:
			r.reportFn(r.statFn())
		case <-ctx.Done():
			r.lggr.Debug("stat reporter loop received done. stopping...")
			return
		}
	}
}
