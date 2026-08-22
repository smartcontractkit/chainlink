package headreporter

import (
	"context"
	"sync"
	"time"

	common "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
)

type (
	HeadReporter interface {
		ReportNewHead(ctx context.Context, head *types.Head) error
		ReportPeriodic(ctx context.Context) error
	}

	Service struct {
		services.StateMachine
		ds             sqlutil.DataSource
		lggr           common.Logger
		newHeads       *mailbox.Mailbox[*types.Head]
		chStop         services.StopChan
		wgDone         sync.WaitGroup
		reportPeriod   time.Duration
		reporters      []HeadReporter
		unsubscribeFns []func()
	}
)

func NewHeadReporterService(ds sqlutil.DataSource, lggr common.Logger, reporters ...HeadReporter) *Service {
	return &Service{
		ds:           ds,
		lggr:         common.Named(lggr, "HeadReporter"),
		newHeads:     mailbox.NewSingle[*types.Head](),
		chStop:       make(chan struct{}),
		reporters:    reporters,
		reportPeriod: 15 * time.Second,
	}
}

func (hrd *Service) Subscribe(subFn func(heads.Trackable) (types.Head, func())) {
	_, unsubscribe := subFn(hrd)
	hrd.unsubscribeFns = append(hrd.unsubscribeFns, unsubscribe)
}

func (hrd *Service) Start(context.Context) error {
	return hrd.StartOnce(hrd.Name(), func() error {
		hrd.wgDone.Add(1)
		go hrd.eventLoop()
		return nil
	})
}

func (hrd *Service) Close() error {
	return hrd.StopOnce(hrd.Name(), func() error {
		close(hrd.chStop)
		hrd.wgDone.Wait()
		return nil
	})
}

func (hrd *Service) Name() string {
	return hrd.lggr.Name()
}

func (hrd *Service) HealthReport() map[string]error {
	return map[string]error{hrd.Name(): hrd.Healthy()}
}

func (hrd *Service) OnNewLongestChain(ctx context.Context, head *types.Head) {
	hrd.newHeads.Deliver(head)
}

func (hrd *Service) eventLoop() {
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
