package v2

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/diskmonitor"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// NewWorkflowModuleCacheDiskMonitor measures cacheDir on tickInterval and exports disk usage
// to both Beholder telemetry and the node /metrics endpoint.
func NewWorkflowModuleCacheDiskMonitor(
	lggr logger.Logger,
	cacheDir string,
	tickInterval time.Duration,
) (*diskmonitor.DiskMonitor, error) {
	return diskmonitor.NewDiskMonitor(
		lggr,
		cacheDir,
		GaugeWorkflowModuleCacheDiskUsageBytes,
		tickInterval,
		diskmonitor.WithPrometheusGauge(promModuleCacheDiskUsageBytes),
	)
}
