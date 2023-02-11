package pg

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/logger"
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

type StatFn func() sql.DBStats

type StatsReporter struct {
	statFn   StatFn
	interval time.Duration
	cancel   context.CancelFunc
	lggr     logger.Logger
	once     sync.Once
}

func NewStatsReporter(fn StatFn, lggr logger.Logger, opts ...StatsReporterOpt) *StatsReporter {
	r := &StatsReporter{
		statFn:   fn,
		interval: dbStatsInternal,
		lggr:     lggr.Named("stat-reporter"),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *StatsReporter) Start(ctx context.Context) {
	run := func() {
		r.lggr.Debugf("Starting DB stat reporter")
		rctx, cancelFunc := context.WithCancel(ctx)
		r.cancel = cancelFunc
		go r.loop(rctx)
	}

	r.once.Do(run)
}

func (r *StatsReporter) Stop() {

	if r.cancel != nil {
		r.lggr.Debugf("Stopping DB stat reporter")
		r.cancel()
		r.cancel = nil
	}

}

func (r *StatsReporter) loop(ctx context.Context) {

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.lggr.Debugf("publishing stats")
			publishStats(r.statFn())
		case <-ctx.Done():
			r.lggr.Debugf("stat reporter loop received done. stopping...")
			return
		}
	}
}
