package ccipexec

import (
	"sync"
	"time"

	"github.com/pkg/errors"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// InflightInternalExecutionReport serves the same purpose as InflightCommitReport
// see the comment on that struct for context.
type InflightInternalExecutionReport struct {
	createdAt time.Time
	messages  []cciptypes.EVM2EVMMessage
}

// inflightExecReportsContainer holds existing inflight reports.
// it provides a thread-safe access as it is called from multiple goroutines,
// e.g. reporting and transmission protocols.
type inflightExecReportsContainer struct {
	locker  sync.RWMutex
	reports []InflightInternalExecutionReport

	cacheExpiry time.Duration
}

func newInflightExecReportsContainer(inflightCacheExpiry time.Duration) *inflightExecReportsContainer {
	return &inflightExecReportsContainer{
		locker:      sync.RWMutex{},
		reports:     make([]InflightInternalExecutionReport, 0),
		cacheExpiry: inflightCacheExpiry,
	}
}

func (container *inflightExecReportsContainer) getAll() []InflightInternalExecutionReport {
	container.locker.RLock()
	defer container.locker.RUnlock()

	reports := make([]InflightInternalExecutionReport, len(container.reports))
	copy(reports, container.reports)

	return reports
}

func (container *inflightExecReportsContainer) expire(lggr logger.Logger) {
	container.locker.Lock()
	defer container.locker.Unlock()
	// Reap old inflight txs and check if any messages in the report are inflight.
	var stillInFlight []InflightInternalExecutionReport
	for _, report := range container.reports {
		if time.Since(report.createdAt) > container.cacheExpiry {
			// Happy path: inflight report was successfully transmitted onchain, we remove it from inflight and onchain state reflects inflight.
			// Sad path: inflight report reverts onchain, we remove it from inflight, onchain state does not reflect the change so we retry.
			lggr.Infow("Inflight report expired", "messages", report.messages)
		} else {
			stillInFlight = append(stillInFlight, report)
		}
	}
	container.reports = stillInFlight
}

func (container *inflightExecReportsContainer) add(lggr logger.Logger, messages []cciptypes.EVM2EVMMessage) error {
	container.locker.Lock()
	defer container.locker.Unlock()

	for _, report := range container.reports {
		if (len(report.messages) > 0) && (report.messages[0].SequenceNumber == messages[0].SequenceNumber) {
			return errors.Errorf("report is already in flight")
		}
	}

	// Otherwise not already in flight, add it.
	lggr.Info("Inflight report added")
	container.reports = append(container.reports, InflightInternalExecutionReport{
		createdAt: time.Now(),
		messages:  messages,
	})
	return nil
}
