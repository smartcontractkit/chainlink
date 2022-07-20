package fluxmonitorv2

import (
	"database/sql"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/sqlx"
)

type transmitter interface {
	CreateEthTransaction(newTx txmgr.NewTx, qopts ...pg.QOpt) (etx txmgr.EthTx, err error)
}

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

// ORM defines an interface for database commands related to Flux Monitor v2
type ORM interface {
	MostRecentFluxMonitorRoundID(aggregator common.Address) (uint32, error)
	DeleteFluxMonitorRoundsBackThrough(aggregator common.Address, roundID uint32) error
	FindOrCreateFluxMonitorRoundStats(aggregator common.Address, roundID uint32, newRoundLogs uint) (FluxMonitorRoundStatsV2, error)
	UpdateFluxMonitorRoundStats(aggregator common.Address, roundID uint32, runID int64, newRoundLogsAddition uint, qopts ...pg.QOpt) error
	CreateEthTransaction(fromAddress, toAddress common.Address, payload []byte, gasLimit uint64, qopts ...pg.QOpt) error
	CountFluxMonitorRoundStats() (count int, err error)
}

type orm struct {
	q        pg.Q
	txm      transmitter
	strategy txmgr.TxStrategy
	checker  txmgr.TransmitCheckerSpec
	logger   logger.Logger
}

// NewORM initializes a new ORM
func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig, txm transmitter, strategy txmgr.TxStrategy, checker txmgr.TransmitCheckerSpec) ORM {
	namedLogger := lggr.Named("FluxMonitorORM")
	q := pg.NewQ(db, namedLogger, cfg)
	return &orm{
		q,
		txm,
		strategy,
		checker,
		namedLogger,
	}
}

// MostRecentFluxMonitorRoundID finds roundID of the most recent round that the
// provided oracle address submitted to
func (o *orm) MostRecentFluxMonitorRoundID(aggregator common.Address) (uint32, error) {
	var stats FluxMonitorRoundStatsV2
	err := o.q.Get(&stats, `SELECT * FROM flux_monitor_round_stats_v2 WHERE aggregator = $1 ORDER BY round_id DESC LIMIT 1`, aggregator)
	return stats.RoundID, errors.Wrap(err, "MostRecentFluxMonitorRoundID failed")
}

// DeleteFluxMonitorRoundsBackThrough deletes all the RoundStat records for a
// given oracle address starting from the most recent round back through the
// given round
func (o *orm) DeleteFluxMonitorRoundsBackThrough(aggregator common.Address, roundID uint32) error {
	_, err := o.q.Exec(`
        DELETE FROM flux_monitor_round_stats_v2
        WHERE aggregator = $1
          AND round_id >= $2
    `, aggregator, roundID)
	return errors.Wrap(err, "DeleteFluxMonitorRoundsBackThrough failed")
}

// FindOrCreateFluxMonitorRoundStats find the round stats record for a given
// oracle on a given round, or creates it if no record exists
func (o *orm) FindOrCreateFluxMonitorRoundStats(aggregator common.Address, roundID uint32, newRoundLogs uint) (stats FluxMonitorRoundStatsV2, err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		err = tx.Get(&stats,
			`INSERT INTO flux_monitor_round_stats_v2 (aggregator, round_id, num_new_round_logs, num_submissions) VALUES ($1, $2, $3, 0)
		ON CONFLICT (aggregator, round_id) DO NOTHING`,
			aggregator, roundID, newRoundLogs)
		if errors.Is(err, sql.ErrNoRows) {
			err = tx.Get(&stats, `SELECT * FROM flux_monitor_round_stats_v2 WHERE aggregator=$1 AND round_id=$2`, aggregator, roundID)
		}
		return err
	})

	return stats, errors.Wrap(err, "FindOrCreateFluxMonitorRoundStats failed")
}

// UpdateFluxMonitorRoundStats trys to create a RoundStat record for the given oracle
// at the given round. If one already exists, it increments the num_submissions column.
func (o *orm) UpdateFluxMonitorRoundStats(aggregator common.Address, roundID uint32, runID int64, newRoundLogsAddition uint, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	err := q.ExecQ(`
        INSERT INTO flux_monitor_round_stats_v2 (
            aggregator, round_id, pipeline_run_id, num_new_round_logs, num_submissions
        ) VALUES (
            $1, $2, $3, $4, 1
        ) ON CONFLICT (aggregator, round_id)
        DO UPDATE SET
          num_new_round_logs = flux_monitor_round_stats_v2.num_new_round_logs + $5,
					num_submissions    = flux_monitor_round_stats_v2.num_submissions + 1,
					pipeline_run_id    = EXCLUDED.pipeline_run_id
    `, aggregator, roundID, runID, newRoundLogsAddition, newRoundLogsAddition)
	return errors.Wrapf(err, "Failed to insert round stats for roundID=%v, runID=%v, newRoundLogsAddition=%v", roundID, runID, newRoundLogsAddition)
}

// CountFluxMonitorRoundStats counts the total number of records
func (o *orm) CountFluxMonitorRoundStats() (count int, err error) {
	err = o.q.Get(&count, `SELECT count(*) FROM flux_monitor_round_stats_v2`)
	return count, errors.Wrap(err, "CountFluxMonitorRoundStats failed")
}

// CreateEthTransaction creates an ethereum transaction for the Txm to pick up
func (o *orm) CreateEthTransaction(
	fromAddress common.Address,
	toAddress common.Address,
	payload []byte,
	gasLimit uint64,
	qopts ...pg.QOpt,
) (err error) {
	_, err = o.txm.CreateEthTransaction(txmgr.NewTx{
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: payload,
		GasLimit:       gasLimit,
		Strategy:       o.strategy,
		Checker:        o.checker,
	}, qopts...)
	return errors.Wrap(err, "Skipped Flux Monitor submission")
}
