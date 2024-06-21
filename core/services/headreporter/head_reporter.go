package headreporter

import (
	"context"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

//go:generate mockery --quiet --name PrometheusBackend --output ../../internal/mocks/ --case=underscore
type (
	Reporter interface {
		reportOnHead(ctx context.Context, head *evmtypes.Head)
		reportPeriodic(ctx context.Context)
	}

	headReporter struct {
		services.StateMachine
		ds           sqlutil.DataSource
		chains       legacyevm.LegacyChainContainer
		lggr         logger.Logger
		newHeads     *mailbox.Mailbox[*evmtypes.Head]
		chStop       services.StopChan
		wgDone       sync.WaitGroup
		reportPeriod time.Duration
		reporters    []Reporter
	}
)

var (
	name = "HeadReporter"
)

func NewHeadReporter(ds sqlutil.DataSource, chainContainer legacyevm.LegacyChainContainer, lggr logger.Logger, opts ...interface{}) *headReporter {
	chStop := make(chan struct{})
	return &headReporter{
		ds:       ds,
		chains:   chainContainer,
		lggr:     lggr.Named(name),
		newHeads: mailbox.NewSingle[*evmtypes.Head](),
		chStop:   chStop,
		reporters: []Reporter{
			NewPrometheusReporter(ds, chainContainer, lggr, opts),
			NewTelemetryReporter(chainContainer, lggr),
		},
	}
}

func (rr *headReporter) Start(context.Context) error {
	return rr.StartOnce(name, func() error {
		rr.wgDone.Add(1)
		go rr.eventLoop()
		return nil
	})
}

func (rr *headReporter) Close() error {
	return rr.StopOnce(name, func() error {
		close(rr.chStop)
		rr.wgDone.Wait()
		return nil
	})
}
func (rr *headReporter) Name() string {
	return rr.lggr.Name()
}

func (rr *headReporter) HealthReport() map[string]error {
	return map[string]error{rr.Name(): rr.Healthy()}
}

func (rr *headReporter) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
	rr.newHeads.Deliver(head)
}

func (rr *headReporter) eventLoop() {
	rr.lggr.Debug("Starting event loop")
	defer rr.wgDone.Done()
	ctx, cancel := rr.chStop.NewCtx()
	defer cancel()
	for {
		select {
		case <-rr.newHeads.Notify():
			head, exists := rr.newHeads.Retrieve()
			if !exists {
				continue
			}
			for _, reporter := range rr.reporters {
				reporter.reportOnHead(ctx, head)
			}
		case <-time.After(rr.reportPeriod):
			for _, reporter := range rr.reporters {
				reporter.reportPeriodic(ctx)
			}
		case <-rr.chStop:
			return
		}
	}
}
