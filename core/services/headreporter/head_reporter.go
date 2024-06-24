package headreporter

import (
	"context"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

//go:generate mockery --quiet --name HeadReporter --output ../../internal/mocks/ --case=underscore
type (
	HeadReporter interface {
		ReportNewHead(ctx context.Context, head *evmtypes.Head)
		ReportPeriodic(ctx context.Context)
	}

	HeadReporterService struct {
		services.StateMachine
		ds           sqlutil.DataSource
		chains       legacyevm.LegacyChainContainer
		lggr         logger.Logger
		newHeads     *mailbox.Mailbox[*evmtypes.Head]
		chStop       services.StopChan
		wgDone       sync.WaitGroup
		reportPeriod time.Duration
		reporters    []HeadReporter
	}
)

func NewHeadReporterService(config config.HeadReport, ds sqlutil.DataSource, chainContainer legacyevm.LegacyChainContainer, lggr logger.Logger, monitoringEndpointGen telemetry.MonitoringEndpointGenerator, opts ...interface{}) *HeadReporterService {
	reporters := make([]HeadReporter, 2)
	reporters = append(reporters, NewPrometheusReporter(ds, chainContainer, lggr, opts))
	if config.TelemetryEnabled() {
		reporters = append(reporters, NewTelemetryReporter(chainContainer, lggr, monitoringEndpointGen))
	}
	return NewHeadReporterServiceWithReporters(ds, chainContainer, lggr, reporters, opts)
}

func NewHeadReporterServiceWithReporters(ds sqlutil.DataSource, chainContainer legacyevm.LegacyChainContainer, lggr logger.Logger, reporters []HeadReporter, opts ...interface{}) *HeadReporterService {
	reportPeriod := 30 * time.Second
	for _, opt := range opts {
		switch v := opt.(type) {
		case time.Duration:
			reportPeriod = v
		}
	}
	chStop := make(chan struct{})
	return &HeadReporterService{
		ds:           ds,
		chains:       chainContainer,
		lggr:         lggr.Named("HeadReporterService"),
		newHeads:     mailbox.NewSingle[*evmtypes.Head](),
		chStop:       chStop,
		wgDone:       sync.WaitGroup{},
		reportPeriod: reportPeriod,
		reporters:    reporters,
	}
}

func (hrd *HeadReporterService) Start(context.Context) error {
	return hrd.StartOnce(hrd.Name(), func() error {
		hrd.wgDone.Add(1)
		go hrd.eventLoop()
		return nil
	})
}

func (hrd *HeadReporterService) Close() error {
	return hrd.StopOnce(hrd.Name(), func() error {
		close(hrd.chStop)
		hrd.wgDone.Wait()
		return nil
	})
}

func (hrd *HeadReporterService) Name() string {
	return hrd.lggr.Name()
}

func (hrd *HeadReporterService) HealthReport() map[string]error {
	return map[string]error{hrd.Name(): hrd.Healthy()}
}

func (hrd *HeadReporterService) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
	hrd.newHeads.Deliver(head)
}

func (hrd *HeadReporterService) eventLoop() {
	hrd.lggr.Debug("Starting event loop")
	defer hrd.wgDone.Done()
	ctx, cancel := hrd.chStop.NewCtx()
	defer cancel()
	for {
		select {
		case <-hrd.newHeads.Notify():
			head, exists := hrd.newHeads.Retrieve()
			if !exists {
				continue
			}
			for _, reporter := range hrd.reporters {
				reporter.ReportNewHead(ctx, head)
			}
		case <-time.After(hrd.reportPeriod):
			for _, reporter := range hrd.reporters {
				reporter.ReportPeriodic(ctx)
			}
		case <-hrd.chStop:
			return
		}
	}
}
