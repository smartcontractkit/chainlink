package services

import (
	"context"
	"database/sql"

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
)

func (defaultBackend) SetUnconfirmedTransactions(n int64) {
	promUnconfirmedTransactions.Set(float64(n))
}

func (defaultBackend) SetMaxUnconfirmedBlocks(n int64) {
	promMaxUnconfirmedBlocks.Set(float64(n))
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
	err := multierr.Combine(
		errors.Wrap(pr.reportPendingEthTxes(ctx), "reportPendingEthTxes failed"),
		errors.Wrap(pr.reportMaxUnconfirmedBlocks(ctx, head), "reportMaxUnconfirmedBlocks failed"),
	)

	if err != nil {
		logger.Errorw("Error reporting prometheus metrics", "err", err)
	}
}

func (pr *promReporter) reportPendingEthTxes(ctx context.Context) (err error) {
	rows, err := pr.db.QueryContext(ctx, `SELECT count(*) FROM eth_txes WHERE state = 'unconfirmed'`)
	if err != nil {
		return errors.Wrap(err, "failed to query for unconfirmed eth_tx count")
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
	rows, err := pr.db.QueryContext(ctx, `
SELECT MIN(broadcast_before_block_num) FROM eth_tx_attempts
JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id
AND eth_txes.state = 'unconfirmed'`)
	if err != nil {
		return errors.Wrap(err, "failed to query for min broadcast_before_block_num")
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
