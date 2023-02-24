package promreporter

import (
	"context"
	"database/sql"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v4"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --quiet --name PrometheusBackend --output ../../internal/mocks/ --case=underscore
type (
	promReporter struct {
		db           *sql.DB
		lggr         logger.Logger
		backend      PrometheusBackend
		newHeads     *utils.Mailbox[*evmtypes.Head]
		chStop       chan struct{}
		wgDone       sync.WaitGroup
		reportPeriod time.Duration

		utils.StartStopOnce
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

func NewPromReporter(db *sql.DB, lggr logger.Logger, opts ...interface{}) *promReporter {
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
		db:           db,
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
	ctx, cancel := utils.ContextFromChan(pr.chStop)
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
	var unconfirmed int64
	if err := pr.db.QueryRowContext(ctx, `SELECT count(*) FROM eth_txes WHERE state = 'unconfirmed' AND evm_chain_id = $1`, evmChainID.String()).Scan(&unconfirmed); err != nil {
		return errors.Wrap(err, "failed to query for unconfirmed eth_tx count")
	}
	pr.backend.SetUnconfirmedTransactions(evmChainID, unconfirmed)
	return nil
}

func (pr *promReporter) reportMaxUnconfirmedAge(ctx context.Context, evmChainID *big.Int) (err error) {
	var broadcastAt null.Time
	now := time.Now()
	if err := pr.db.QueryRowContext(ctx, `SELECT min(initial_broadcast_at) FROM eth_txes WHERE state = 'unconfirmed' AND evm_chain_id = $1`, evmChainID.String()).Scan(&broadcastAt); err != nil {
		return errors.Wrap(err, "failed to query for unconfirmed eth_tx count")
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
	var earliestUnconfirmedTxBlock null.Int
	err = pr.db.QueryRowContext(ctx, `
SELECT MIN(broadcast_before_block_num) FROM eth_tx_attempts
JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id
WHERE eth_txes.state = 'unconfirmed'
AND evm_chain_id = $1
AND eth_txes.state = 'unconfirmed'`, head.EVMChainID.String()).Scan(&earliestUnconfirmedTxBlock)
	if err != nil {
		return errors.Wrap(err, "failed to query for min broadcast_before_block_num")
	}
	var blocksUnconfirmed int64
	if !earliestUnconfirmedTxBlock.IsZero() {
		blocksUnconfirmed = head.Number - earliestUnconfirmedTxBlock.ValueOrZero()
	}
	pr.backend.SetMaxUnconfirmedBlocks(head.EVMChainID.ToInt(), blocksUnconfirmed)
	return nil
}

func (pr *promReporter) reportPipelineRunStats(ctx context.Context) (err error) {
	rows, err := pr.db.QueryContext(ctx, `
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
