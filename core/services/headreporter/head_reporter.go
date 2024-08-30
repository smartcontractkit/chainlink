package headreporter

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type (
	HeadReporter interface {
		ReportNewHead(ctx context.Context, head *evmtypes.Head) error
		ReportPeriodic(ctx context.Context) error
	}

	HeadReporterService struct {
		services.StateMachine
		ds             sqlutil.DataSource
		lggr           logger.Logger
		newHeads       *mailbox.Mailbox[*evmtypes.Head]
		chStop         services.StopChan
		wgDone         sync.WaitGroup
		reportPeriod   time.Duration
		reporters      []HeadReporter
		unsubscribeFns []func()
	}
)

func NewHeadReporterService(ds sqlutil.DataSource, lggr logger.Logger, reporters ...HeadReporter) *HeadReporterService {
	return &HeadReporterService{
		ds:           ds,
		lggr:         lggr.Named("HeadReporter"),
		newHeads:     mailbox.NewSingle[*evmtypes.Head](),
		chStop:       make(chan struct{}),
		reporters:    reporters,
		reportPeriod: 15 * time.Second,
	}
}

func (hrd *HeadReporterService) Subscribe(subFn func(types.HeadTrackable) (evmtypes.Head, func())) {
	_, unsubscribe := subFn(hrd)
	hrd.unsubscribeFns = append(hrd.unsubscribeFns, unsubscribe)
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
	after := time.After(hrd.reportPeriod)
	for {
		select {
		case <-hrd.newHeads.Notify():
			head, exists := hrd.newHeads.Retrieve()
			if !exists {
				continue
			}
			for _, reporter := range hrd.reporters {
				err := reporter.ReportNewHead(ctx, head)
				if err != nil && ctx.Err() == nil {
					hrd.lggr.Errorw("Error reporting new head", "err", err)
				}
			}
		case <-after:
			for _, reporter := range hrd.reporters {
				err := reporter.ReportPeriodic(ctx)
				if err != nil && ctx.Err() == nil {
					hrd.lggr.Errorw("Error in periodic report", "err", err)
				}
			}
			after = time.After(hrd.reportPeriod)
		case <-hrd.chStop:
			return
		}
	}
}
