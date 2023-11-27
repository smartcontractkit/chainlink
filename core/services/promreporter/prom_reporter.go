package promreporter

import (
	"context"
	"fmt"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"go.uber.org/multierr"
)

//go:generate mockery --quiet --name PrometheusBackend --output ../../internal/mocks/ --case=underscore
type (
	promReporter struct {
		services.StateMachine
		txStore      txmgr.TxStore
		lggr         logger.Logger
		backend      PrometheusBackend
		newHeads     *utils.Mailbox[*evmtypes.Head]
		chStop       utils.StopChan
		wgDone       sync.WaitGroup
		reportPeriod time.Duration
	}

	PrometheusBackend interface {
		SetUnconfirmedTransactions(*big.Int, int64)
		SetMaxUnconfirmedAge(*big.Int, float64)
		SetMaxUnconfirmedBlocks(*big.Int, int64)
		SetPipelineRunsQueued(n int)
		SetPipelineTaskRunsQueued(n int)
	}

	defaultBackend struct{}
)

var (
	promUnconfirmedTransactions = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "unconfirmed_transactions",
		Help: "Number of currently unconfirmed transactions",
	}, []string{"evmChainID"})
	promMaxUnconfirmedAge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "max_unconfirmed_tx_age",
		Help: "The length of time the oldest unconfirmed transaction has been in that state (in seconds). Will be 0 if there are no unconfirmed transactions.",
	}, []string{"evmChainID"})
	promMaxUnconfirmedBlocks = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "max_unconfirmed_blocks",
		Help: "The max number of blocks any currently unconfirmed transaction has been unconfirmed for",
	}, []string{"evmChainID"})
	promPipelineRunsQueued = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pipeline_runs_queued",
		Help: "The total number of pipeline runs that are awaiting execution",
	})
	promPipelineTaskRunsQueued = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pipeline_task_runs_queued",
		Help: "The total number of pipeline task runs that are awaiting execution",
	})
)

func (defaultBackend) SetUnconfirmedTransactions(evmChainID *big.Int, n int64) {
	promUnconfirmedTransactions.WithLabelValues(evmChainID.String()).Set(float64(n))
}

func (defaultBackend) SetMaxUnconfirmedAge(evmChainID *big.Int, s float64) {
	promMaxUnconfirmedAge.WithLabelValues(evmChainID.String()).Set(s)
}

func (defaultBackend) SetMaxUnconfirmedBlocks(evmChainID *big.Int, n int64) {
	promMaxUnconfirmedBlocks.WithLabelValues(evmChainID.String()).Set(float64(n))
}

func (defaultBackend) SetPipelineRunsQueued(n int) {
	promPipelineTaskRunsQueued.Set(float64(n))
}

func (defaultBackend) SetPipelineTaskRunsQueued(n int) {
	promPipelineRunsQueued.Set(float64(n))
}

func NewPromReporter(txStore txmgr.TxStore, lggr logger.Logger, opts ...interface{}) *promReporter {
	var backend PrometheusBackend = defaultBackend{}
	period := 15 * time.Second
	for _, opt := range opts {
		switch v := opt.(type) {
		case time.Duration:
			period = v
		case PrometheusBackend:
			backend = v
		}
	}

	chStop := make(chan struct{})
	return &promReporter{
		txStore:      txStore,
		lggr:         lggr.Named("PromReporter"),
		backend:      backend,
		newHeads:     utils.NewSingleMailbox[*evmtypes.Head](),
		chStop:       chStop,
		reportPeriod: period,
	}
}

// Start starts PromReporter.
func (pr *promReporter) Start(context.Context) error {
	return pr.StartOnce("PromReporter", func() error {
		pr.wgDone.Add(1)
		go pr.eventLoop()
		return nil
	})
}

func (pr *promReporter) Close() error {
	return pr.StopOnce("PromReporter", func() error {
		close(pr.chStop)
		pr.wgDone.Wait()
		return nil
	})
}
func (pr *promReporter) Name() string {
	return pr.lggr.Name()
}

func (pr *promReporter) HealthReport() map[string]error {
	return map[string]error{pr.Name(): pr.Healthy()}
}

func (pr *promReporter) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
	pr.newHeads.Deliver(head)
}

func (pr *promReporter) eventLoop() {
	pr.lggr.Debug("Starting event loop")
	defer pr.wgDone.Done()
	ctx, cancel := pr.chStop.NewCtx()
	defer cancel()
	for {
		select {
		case <-pr.newHeads.Notify():
			head, exists := pr.newHeads.Retrieve()
			if !exists {
				continue
			}
			pr.reportHeadMetrics(ctx, head)
		case <-time.After(pr.reportPeriod):
			if err := errors.Wrap(pr.reportPipelineRunStats(ctx), "reportPipelineRunStats failed"); err != nil {
				pr.lggr.Errorw("Error reporting prometheus metrics", "err", err)
			}

		case <-pr.chStop:
			return
		}
	}
}

func (pr *promReporter) reportHeadMetrics(ctx context.Context, head *evmtypes.Head) {
	evmChainID := head.EVMChainID.ToInt()
	err := multierr.Combine(
		errors.Wrap(pr.reportPendingEthTxes(ctx, evmChainID), "reportPendingEthTxes failed"),
		errors.Wrap(pr.reportMaxUnconfirmedAge(ctx, evmChainID), "reportMaxUnconfirmedAge failed"),
		errors.Wrap(pr.reportMaxUnconfirmedBlocks(ctx, head), "reportMaxUnconfirmedBlocks failed"),
	)

	if err != nil && ctx.Err() == nil {
		pr.lggr.Errorw("Error reporting prometheus metrics", "err", err)
	}
}

func (pr *promReporter) reportPendingEthTxes(ctx context.Context, evmChainID *big.Int) (err error) {
	unconfirmed, err := pr.txStore.CountAllUnconfirmedTransactions(ctx, evmChainID)
	if err != nil {
		return fmt.Errorf("failed to query for unconfirmed eth_tx count: %w", err)
	}
	pr.backend.SetUnconfirmedTransactions(evmChainID, int64(unconfirmed))
	return nil
}

func (pr *promReporter) reportMaxUnconfirmedAge(ctx context.Context, evmChainID *big.Int) (err error) {
	now := time.Now()
	broadcastAt, err := pr.txStore.FindMinUnconfirmedBroadcastTime(ctx, evmChainID)
	if err != nil {
		return fmt.Errorf("failed to query for min broadcast time: %w", err)
	}

	var seconds float64
	if broadcastAt.Valid {
		nanos := now.Sub(broadcastAt.ValueOrZero())
		seconds = float64(nanos) / 1000000000
	}
	pr.backend.SetMaxUnconfirmedAge(evmChainID, seconds)
	return nil
}

func (pr *promReporter) reportMaxUnconfirmedBlocks(ctx context.Context, head *evmtypes.Head) (err error) {
	earliestUnconfirmedTxBlock, err := pr.txStore.FindEarliestUnconfirmedTxBlock(ctx, head.EVMChainID.ToInt())
	if err != nil {
		return fmt.Errorf("failed to query for earliest unconfirmed tx block: %w", err)
	}

	var blocksUnconfirmed int64
	if !earliestUnconfirmedTxBlock.IsZero() {
		blocksUnconfirmed = head.Number - earliestUnconfirmedTxBlock.ValueOrZero()
	}
	pr.backend.SetMaxUnconfirmedBlocks(head.EVMChainID.ToInt(), blocksUnconfirmed)
	return nil
}

func (pr *promReporter) reportPipelineRunStats(ctx context.Context) (err error) {
	pipelineTaskRunsQueued, pipelineRunsQueued, err := pr.txStore.GetPipelineRunStats(ctx)
	if err != nil {
		return fmt.Errorf("failed to query for pipeline run stats: %w", err)
	}

	pr.backend.SetPipelineTaskRunsQueued(pipelineTaskRunsQueued)
	pr.backend.SetPipelineRunsQueued(pipelineRunsQueued)

	return nil
}
