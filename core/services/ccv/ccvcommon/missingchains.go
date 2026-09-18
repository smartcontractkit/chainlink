package ccvcommon

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

// DefaultMissingChainsReportInterval is the time between two reports of the same missing chain.
const DefaultMissingChainsReportInterval = 5 * time.Minute

const (
	keyChainSelector = "chain_selector"
	keyJobName       = "job_name"
)

// missingChainsMetrics holds the metric instruments of the MissingChainsMonitor.
type missingChainsMetrics struct {
	missingChain metric.Int64Gauge
}

func newMissingChainsMetrics() (*missingChainsMetrics, error) {
	missingChain, err := beholder.GetMeter().Int64Gauge("ccv_missing_chain")
	if err != nil {
		return nil, fmt.Errorf("failed to create ccv_missing_chain gauge: %w", err)
	}
	return &missingChainsMetrics{missingChain: missingChain}, nil
}

// MissingChainsMonitor reports, at a regular interval, each chain that the job configuration names
// but that has no chain object on this node. The report is a critical log line and a metric with a
// chain_selector label.
//
// The monitor stays healthy. A missing chain is a configuration alarm, not a reason to restart the
// job.
type MissingChainsMonitor struct {
	services.StateMachine
	lggr     logger.SugaredLogger
	name     string
	metrics  *missingChainsMetrics
	missing  []protocol.ChainSelector
	interval time.Duration
	stopCh   services.StopChan
	wg       sync.WaitGroup
}

// NewMissingChainsMonitor makes a monitor for the given missing chain selectors. It returns a nil
// monitor when no chain is missing, so that a correctly configured node starts no extra goroutine.
// The caller must check for nil before it appends the monitor to a service list.
func NewMissingChainsMonitor(
	lggr logger.Logger,
	name string,
	missing []protocol.ChainSelector,
	interval time.Duration,
) (*MissingChainsMonitor, error) {
	if len(missing) == 0 {
		return nil, nil
	}
	if interval <= 0 {
		interval = DefaultMissingChainsReportInterval
	}
	m, err := newMissingChainsMetrics()
	if err != nil {
		return nil, err
	}
	return &MissingChainsMonitor{
		lggr:     logger.Sugared(lggr).Named("MissingChainsMonitor"),
		name:     name,
		metrics:  m,
		missing:  missing,
		interval: interval,
		stopCh:   make(services.StopChan),
	}, nil
}

func (m *MissingChainsMonitor) Start(ctx context.Context) error {
	return m.StartOnce(m.Name(), func() error {
		m.report(ctx)
		m.wg.Go(m.reportLoop)
		return nil
	})
}

func (m *MissingChainsMonitor) Close() error {
	return m.StopOnce(m.Name(), func() error {
		close(m.stopCh)
		m.wg.Wait()
		return nil
	})
}

func (m *MissingChainsMonitor) reportLoop() {
	ctx, cancel := m.stopCh.NewCtx()
	defer cancel()

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.report(ctx)
		}
	}
}

func (m *MissingChainsMonitor) report(ctx context.Context) {
	for _, sel := range m.missing {
		m.lggr.Criticalw(
			"chain in the job configuration has no chain object on this node; the job continues without it",
			"jobName", m.name,
			"chainSelector", uint64(sel),
		)
		m.metrics.missingChain.Record(ctx, 1, metric.WithAttributes(
			attribute.String(keyChainSelector, strconv.FormatUint(uint64(sel), 10)),
			attribute.String(keyJobName, m.name),
		))
	}
}

func (m *MissingChainsMonitor) HealthReport() map[string]error {
	return map[string]error{m.Name(): m.Ready()}
}

func (m *MissingChainsMonitor) Name() string {
	return m.lggr.Name()
}
