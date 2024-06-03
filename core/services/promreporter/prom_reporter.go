package promreporter

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

//go:generate mockery --quiet --name PrometheusBackend --output ../../internal/mocks/ --case=underscore
type (
	promReporter struct {
		services.StateMachine
		ds           sqlutil.DataSource
		chains       legacyevm.LegacyChainContainer
		lggr         logger.Logger
		backend      PrometheusBackend
		newHeads     *mailbox.Mailbox[*evmtypes.Head]
		chStop       services.StopChan
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

func NewPromReporter(ds sqlutil.DataSource, chainContainer legacyevm.LegacyChainContainer, lggr logger.Logger, opts ...interface{}) *promReporter {
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
		ds:           ds,
		chains:       chainContainer,
		lggr:         lggr.Named("PromReporter"),
		backend:      backend,
		newHeads:     mailbox.NewSingle[*evmtypes.Head](),
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

func (pr *promReporter) getTxm(evmChainID *big.Int) (txmgr.TxManager, error) {
	chain, err := pr.chains.Get(evmChainID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain: %w", err)
	}
	return chain.TxManager(), nil
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
	txm, err := pr.getTxm(evmChainID)
	if err != nil {
		return fmt.Errorf("failed to get txm: %w", err)
	}

	unconfirmed, err := txm.CountTransactionsByState(ctx, txmgrcommon.TxUnconfirmed)
	if err != nil {
		return fmt.Errorf("failed to query for unconfirmed eth_tx count: %w", err)
	}
	pr.backend.SetUnconfirmedTransactions(evmChainID, int64(unconfirmed))
	return nil
}

func (pr *promReporter) reportMaxUnconfirmedAge(ctx context.Context, evmChainID *big.Int) (err error) {
	txm, err := pr.getTxm(evmChainID)
	if err != nil {
		return fmt.Errorf("failed to get txm: %w", err)
	}

	broadcastAt, err := txm.FindEarliestUnconfirmedBroadcastTime(ctx)
	if err != nil {
		return fmt.Errorf("failed to query for min broadcast time: %w", err)
	}

	var seconds float64
	if broadcastAt.Valid {
		seconds = time.Since(broadcastAt.ValueOrZero()).Seconds()
	}
	pr.backend.SetMaxUnconfirmedAge(evmChainID, seconds)
	return nil
}

func (pr *promReporter) reportMaxUnconfirmedBlocks(ctx context.Context, head *evmtypes.Head) (err error) {
	txm, err := pr.getTxm(head.EVMChainID.ToInt())
	if err != nil {
		return fmt.Errorf("failed to get txm: %w", err)
	}

	earliestUnconfirmedTxBlock, err := txm.FindEarliestUnconfirmedTxAttemptBlock(ctx)
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
	rows, err := pr.ds.QueryContext(ctx, `
SELECT pipeline_run_id FROM pipeline_task_runs WHERE finished_at IS NULL
`)
	if err != nil {
		return errors.Wrap(err, "failed to query for pipeline_run_id")
	}
	defer func() {
		err = multierr.Combine(err, rows.Close())
	}()

	pipelineTaskRunsQueued := 0
	pipelineRunsQueuedSet := make(map[int32]struct{})
	for rows.Next() {
		var pipelineRunID int32
		if err = rows.Scan(&pipelineRunID); err != nil {
			return errors.Wrap(err, "unexpected error scanning row")
		}
		pipelineTaskRunsQueued++
		pipelineRunsQueuedSet[pipelineRunID] = struct{}{}
	}
	if err = rows.Err(); err != nil {
		return err
	}
	pipelineRunsQueued := len(pipelineRunsQueuedSet)

	pr.backend.SetPipelineTaskRunsQueued(pipelineTaskRunsQueued)
	pr.backend.SetPipelineRunsQueued(pipelineRunsQueued)

	return nil
}
