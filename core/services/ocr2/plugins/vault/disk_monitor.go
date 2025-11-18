package vault

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	defaultTickInterval = 1 * time.Minute
)

type DiskMonitor struct {
	services.Service
	eng          *services.Engine
	tickInterval time.Duration
	lggr         logger.Logger
	sizeOfDir    func() (int64, error)

	gauge gauge
}

type gauge interface {
	Record(ctx context.Context, value int64, options ...metric.RecordOption)
}

// NewDiskMonitor creates a new DiskMonitor for the given directory path.
func NewDiskMonitor(lggr logger.Logger, dirPath string) (*DiskMonitor, error) {
	gauge, err := beholder.GetMeter().Int64Gauge("platform_vault_disk_usage_bytes")
	if err != nil {
		return nil, fmt.Errorf("failed to create gauge for vault DiskMonitor: %w", err)
	}

	dm := &DiskMonitor{
		gauge:        gauge,
		tickInterval: defaultTickInterval,
		lggr:         lggr.Named("DiskMonitor").With("dirPath", dirPath),
		sizeOfDir: func() (int64, error) {
			var totalSize int64
			err := filepath.Walk(dirPath, func(_ string, info os.FileInfo, ierr error) error {
				if ierr == nil && !info.IsDir() {
					totalSize += info.Size()
				}
				return nil
			})
			return totalSize, err
		},
	}

	dm.Service, dm.eng = services.Config{
		Name:  "DiskMonitor",
		Start: dm.start,
	}.NewServiceEngine(lggr)
	return dm, nil
}

// Start begins monitoring the directory size in a background goroutine.
func (dm *DiskMonitor) start(ctx context.Context) error {
	ticker := services.TickerConfig{}.NewTicker(dm.tickInterval)
	dm.eng.GoTick(ticker, dm.emitDirSizeMetric)
	return nil
}

// emitDirSizeMetric calculates the total size of the directory in bytes and emits a beholder metric.
func (dm *DiskMonitor) emitDirSizeMetric(ctx context.Context) {
	totalSize, err := dm.sizeOfDir()
	if err != nil {
		dm.lggr.Errorw("Failed to measure vault directory size", "error", err)
		return
	}

	dm.lggr.Debugw("Emitting vault directory size metric", "sizeBytes", totalSize)
	dm.gauge.Record(ctx, totalSize)
}
