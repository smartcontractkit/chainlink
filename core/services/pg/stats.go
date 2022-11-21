package pg

import (
	"database/sql"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
