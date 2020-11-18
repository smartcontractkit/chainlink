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
)

type promReporter struct {
	db *sql.DB
}

var (
	promUnconfirmedTransactions = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "unconfirmed_transactions",
		Help: "Number of currently unconfirmed transactions",
	})
)

func NewPromReporter(db *sql.DB) store.HeadTrackable {
	return &promReporter{db}
}

// Do nothing on connect, simply wait for the next head
func (pr *promReporter) Connect(*models.Head) error {
	return nil
}

func (pr *promReporter) Disconnect() {
	// pass
}

func (pr *promReporter) OnNewLongestChain(ctx context.Context, head models.Head) {
	if err := pr.reportPendingEthTxes(ctx); err != nil {
		logger.Error(err)
	}
}

func (pr *promReporter) reportPendingEthTxes(ctx context.Context) error {
	rows, err := pr.db.QueryContext(ctx, `SELECT count(*) FROM eth_txes WHERE state = 'unconfirmed'`)
	if err != nil {
		return errors.Wrap(err, "failed to query for unconfirmed eth_tx count")
	}
	defer logger.ErrorIfCalling(rows.Close)
	var unconfirmed int64
	for rows.Next() {
		if err := rows.Scan(&unconfirmed); err != nil {
			return errors.Wrap(err, "unexpected error scanning row")
		}
	}
	promUnconfirmedTransactions.Set(float64(unconfirmed))
	return nil
}
