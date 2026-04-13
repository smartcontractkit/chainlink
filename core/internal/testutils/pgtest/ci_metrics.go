package pgtest

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

// EnvCIMetrics enables stderr metrics for CI (periodic samples + peak milestones).
// Set CHAINLINK_PGTEST_CI_METRICS=true in the workflow; search job logs for "pgtest_ci_metrics".
const EnvCIMetrics = "CHAINLINK_PGTEST_CI_METRICS"

var lastPeakReported atomic.Int64

func init() {
	if os.Getenv(EnvCIMetrics) != "true" {
		return
	}
	go ciMetricsTicker()
}

func ciMetricsTicker() {
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()
	for range tick.C {
		fmt.Fprintf(os.Stderr, "%s pgtest_ci_metrics sample peak_concurrent_txdb=%d current_open=%d\n",
			time.Now().Format(time.RFC3339), peakTxdbSessions.Load(), openTxdbSessions.Load())
	}
}

// maybeReportPeakCIMetrics logs when concurrent txdb sessions reach a new high (per OS process / test package).
func maybeReportPeakCIMetrics() {
	if os.Getenv(EnvCIMetrics) != "true" {
		return
	}
	p := peakTxdbSessions.Load()
	o := openTxdbSessions.Load()
	if p < 10 {
		return
	}
	prev := lastPeakReported.Load()
	if prev > 0 && p-prev < 25 {
		return
	}
	if !lastPeakReported.CompareAndSwap(prev, p) {
		return
	}
	fmt.Fprintf(os.Stderr, "%s pgtest_ci_metrics peak_concurrent_txdb=%d current_open=%d (since_last_report=%d)\n",
		time.Now().Format(time.RFC3339), p, o, p-prev)
}
