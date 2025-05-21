package fluxmonitorv2

import (
	"context"
	"database/sql"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type transmitter interface {
	CreateTransaction(ctx context.Context, txRequest txmgr.TxRequest) (tx txmgr.Tx, err error)
}

// ORM defines an interface for database commands related to Flux Monitor v2
type ORM interface {
	MostRecentFluxMonitorRoundID(ctx context.Context, aggregator common.Address) (uint32, error)
	DeleteFluxMonitorRoundsBackThrough(ctx context.Context, aggregator common.Address, roundID uint32) error
	FindOrCreateFluxMonitorRoundStats(ctx context.Context, aggregator common.Address, roundID uint32, newRoundLogs uint) (FluxMonitorRoundStatsV2, error)
	UpdateFluxMonitorRoundStats(ctx context.Context, aggregator common.Address, roundID uint32, runID int64, newRoundLogsAddition uint) error
	CreateEthTransaction(ctx context.Context, fromAddress, toAddress common.Address, payload []byte, gasLimit uint64, idempotencyKey *string) error
	CountFluxMonitorRoundStats(ctx context.Context) (count int, err error)

	WithDataSource(sqlutil.DataSource) ORM
}

type orm struct {
	ds       sqlutil.DataSource
	txm      transmitter
	strategy types.TxStrategy
	checker  txmgr.TransmitCheckerSpec
	logger   logger.Logger
}

func (o *orm) WithDataSource(ds sqlutil.DataSource) ORM { return o.withDataSource(ds) }

func (o *orm) withDataSource(ds sqlutil.DataSource) *orm {
	return &orm{ds, o.txm, o.strategy, o.checker, o.logger}
}

// NewORM initializes a new ORM
func NewORM(ds sqlutil.DataSource, lggr logger.Logger, txm transmitter, strategy types.TxStrategy, checker txmgr.TransmitCheckerSpec) ORM {
	namedLogger := lggr.Named("FluxMonitorORM")
	return &orm{ds, txm, strategy, checker, namedLogger}
}

// MostRecentFluxMonitorRoundID finds roundID of the most recent round that the
// provided oracle address submitted to
func (o *orm) MostRecentFluxMonitorRoundID(ctx context.Context, aggregator common.Address) (uint32, error) {
	var stats FluxMonitorRoundStatsV2
	err := o.ds.GetContext(ctx, &stats, `SELECT * FROM flux_monitor_round_stats_v2 WHERE aggregator = $1 ORDER BY round_id DESC LIMIT 1`, aggregator)
	return stats.RoundID, errors.Wrap(err, "MostRecentFluxMonitorRoundID failed")
}

// DeleteFluxMonitorRoundsBackThrough deletes all the RoundStat records for a
// given oracle address starting from the most recent round back through the
// given round
func (o *orm) DeleteFluxMonitorRoundsBackThrough(ctx context.Context, aggregator common.Address, roundID uint32) error {
	_, err := o.ds.ExecContext(ctx, `
        DELETE FROM flux_monitor_round_stats_v2
        WHERE aggregator = $1
          AND round_id >= $2
    `, aggregator, roundID)
	return errors.Wrap(err, "DeleteFluxMonitorRoundsBackThrough failed")
}

// FindOrCreateFluxMonitorRoundStats find the round stats record for a given
// oracle on a given round, or creates it if no record exists
func (o *orm) FindOrCreateFluxMonitorRoundStats(ctx context.Context, aggregator common.Address, roundID uint32, newRoundLogs uint) (stats FluxMonitorRoundStatsV2, err error) {
	err = sqlutil.Transact(ctx, o.withDataSource, o.ds, nil, func(tx *orm) error {
		err = tx.ds.GetContext(ctx, &stats,
			`INSERT INTO flux_monitor_round_stats_v2 (aggregator, round_id, num_new_round_logs, num_submissions) VALUES ($1, $2, $3, 0)
		ON CONFLICT (aggregator, round_id) DO NOTHING`,
			aggregator, roundID, newRoundLogs)
		if errors.Is(err, sql.ErrNoRows) {
			err = tx.ds.GetContext(ctx, &stats, `SELECT * FROM flux_monitor_round_stats_v2 WHERE aggregator=$1 AND round_id=$2`, aggregator, roundID)
		}
		return err
	})

	return stats, errors.Wrap(err, "FindOrCreateFluxMonitorRoundStats failed")
}

// UpdateFluxMonitorRoundStats trys to create a RoundStat record for the given oracle
// at the given round. If one already exists, it increments the num_submissions column.
func (o *orm) UpdateFluxMonitorRoundStats(ctx context.Context, aggregator common.Address, roundID uint32, runID int64, newRoundLogsAddition uint) error {
	_, err := o.ds.ExecContext(ctx, `
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
func (o *orm) CountFluxMonitorRoundStats(ctx context.Context) (count int, err error) {
	err = o.ds.GetContext(ctx, &count, `SELECT count(*) FROM flux_monitor_round_stats_v2`)
	return count, errors.Wrap(err, "CountFluxMonitorRoundStats failed")
}

// CreateEthTransaction creates an ethereum transaction for the Txm to pick up
func (o *orm) CreateEthTransaction(
	ctx context.Context,
	fromAddress common.Address,
	toAddress common.Address,
	payload []byte,
	gasLimit uint64,
	idempotencyKey *string,
) (err error) {
	_, err = o.txm.CreateTransaction(ctx, txmgr.TxRequest{
		IdempotencyKey: idempotencyKey,
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: payload,
		FeeLimit:       gasLimit,
		Strategy:       o.strategy,
		Checker:        o.checker,
	})
	return errors.Wrap(err, "Skipped Flux Monitor submission")
}
