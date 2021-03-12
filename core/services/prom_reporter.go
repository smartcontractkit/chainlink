package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v4"
)

//go:generate mockery --name PrometheusBackend --output ../internal/mocks/ --case=underscore
type (
	promReporter struct {
		db      *sql.DB
		backend PrometheusBackend
	}

	PrometheusBackend interface {
		SetUnconfirmedTransactions(int64)
		SetMaxUnconfirmedBlocks(int64)
		SetPipelineRunsQueued(n int)
		SetPipelineTaskRunsQueued(n int)
	}

	defaultBackend struct{}
)

var (
	promUnconfirmedTransactions = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "unconfirmed_transactions",
		Help: "Number of currently unconfirmed transactions",
	})
	promMaxUnconfirmedBlocks = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "max_unconfirmed_blocks",
		Help: "The max number of blocks any currently unconfirmed transaction has been unconfirmed for",
	})
	promPipelineRunsQueued = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pipeline_runs_queued",
		Help: "The total number of pipeline runs that are awaiting execution",
	})
	promPipelineTaskRunsQueued = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pipeline_task_runs_queued",
		Help: "The total number of pipeline task runs that are awaiting execution",
	})
)

func (defaultBackend) SetUnconfirmedTransactions(n int64) {
	promUnconfirmedTransactions.Set(float64(n))
}

func (defaultBackend) SetMaxUnconfirmedBlocks(n int64) {
	promMaxUnconfirmedBlocks.Set(float64(n))
}

func (defaultBackend) SetPipelineRunsQueued(n int) {
	promPipelineTaskRunsQueued.Set(float64(n))
}

func (defaultBackend) SetPipelineTaskRunsQueued(n int) {
	promPipelineRunsQueued.Set(float64(n))
}

func NewPromReporter(db *sql.DB, opts ...PrometheusBackend) store.HeadTrackable {
	var backend PrometheusBackend
	if len(opts) > 0 {
		backend = opts[0]
	} else {
		backend = defaultBackend{}
	}

	return &promReporter{db, backend}
}

// Do nothing on connect, simply wait for the next head
func (pr *promReporter) Connect(*models.Head) error {
	return nil
}

func (pr *promReporter) Disconnect() {
	// pass
}

func (pr *promReporter) OnNewLongestChain(ctx context.Context, head models.Head) {
	start := time.Now()
	err := multierr.Combine(
		errors.Wrap(pr.reportPendingEthTxes(ctx), "reportPendingEthTxes failed"),
		errors.Wrap(pr.reportMaxUnconfirmedBlocks(ctx, head), "reportMaxUnconfirmedBlocks failed"),
		errors.Wrap(pr.reportPipelineRunStats(ctx), "reportPipelineRunStats failed"),
	)

	if err != nil {
		deadline, _ := ctx.Deadline()
		took := time.Now().Sub(start) * time.Millisecond
		available := deadline.Sub(start) * time.Millisecond
		logger.Errorw(fmt.Sprintf("Error reporting prometheus metrics. Took: %d ms, time available when starting: %d ms", took, available), "err", err)
	}
}

func (pr *promReporter) reportPendingEthTxes(ctx context.Context) (err error) {
	start := time.Now()
	rows, err := pr.db.QueryContext(ctx, `SELECT count(*) FROM eth_txes WHERE state = 'unconfirmed'`)
	if err != nil {
		deadline, _ := ctx.Deadline()
		return errors.Wrap(err,
			fmt.Sprintf("failed to query for unconfirmed eth_tx count. Query took: %d ms. Time left when starting: %d ms",
				time.Now().Sub(start) * time.Millisecond,
				deadline.Sub(start) * time.Millisecond,
			),
		)
	}

	defer func() {
		err = multierr.Combine(err, rows.Close())
	}()

	var unconfirmed int64
	for rows.Next() {
		if err := rows.Scan(&unconfirmed); err != nil {
			return errors.Wrap(err, "unexpected error scanning row")
		}
	}
	pr.backend.SetUnconfirmedTransactions(unconfirmed)
	return nil
}

func (pr *promReporter) reportMaxUnconfirmedBlocks(ctx context.Context, head models.Head) (err error) {
	start := time.Now()
	rows, err := pr.db.QueryContext(ctx, `
SELECT MIN(broadcast_before_block_num) FROM eth_tx_attempts
JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id
AND eth_txes.state = 'unconfirmed'`)
	if err != nil {
		deadline, _ := ctx.Deadline()
		return errors.Wrap(err,
			fmt.Sprintf("failed to query for min broadcast_before_block_num. Query took: %d ms. Time left when starting: %d ms",
				time.Now().Sub(start) * time.Millisecond,
				deadline.Sub(start) * time.Millisecond,
			),
		)
	}
	defer func() {
		err = multierr.Combine(err, rows.Close())
	}()

	var earliestUnconfirmedTxBlock null.Int
	for rows.Next() {
		if err := rows.Scan(&earliestUnconfirmedTxBlock); err != nil {
			return errors.Wrap(err, "unexpected error scanning row")
		}
	}
	var blocksUnconfirmed int64
	if !earliestUnconfirmedTxBlock.IsZero() {
		blocksUnconfirmed = head.Number - earliestUnconfirmedTxBlock.ValueOrZero()
	}
	pr.backend.SetMaxUnconfirmedBlocks(blocksUnconfirmed)
	return nil
}

func (pr *promReporter) reportPipelineRunStats(ctx context.Context) (err error) {
	start := time.Now()
	rows, err := pr.db.QueryContext(ctx, `
SELECT pipeline_run_id FROM pipeline_task_runs WHERE finished_at IS NULL
`)
	if err != nil {
		deadline, _ := ctx.Deadline()
		return errors.Wrap(err,
			fmt.Sprintf("failed to query for pipeline_run_id. Query took: %d ms. Time left when starting: %d ms",
				time.Now().Sub(start) * time.Millisecond,
				deadline.Sub(start) * time.Millisecond,
			),
		)
	}
	defer func() {
		err = multierr.Combine(err, rows.Close())
	}()

	pipelineTaskRunsQueued := 0
	pipelineRunsQueuedSet := make(map[int32]struct{})
	for rows.Next() {
		var pipelineRunID int32
		if err := rows.Scan(&pipelineRunID); err != nil {
			return errors.Wrap(err, "unexpected error scanning row")
		}
		pipelineTaskRunsQueued++
		pipelineRunsQueuedSet[pipelineRunID] = struct{}{}
	}
	pipelineRunsQueued := len(pipelineRunsQueuedSet)

	pr.backend.SetPipelineTaskRunsQueued(pipelineTaskRunsQueued)
	pr.backend.SetPipelineRunsQueued(pipelineRunsQueued)

	return nil
}
