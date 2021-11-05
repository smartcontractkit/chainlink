package fluxmonitorv2

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"gorm.io/gorm"
)

type transmitter interface {
	CreateEthTransaction(newTx bulletprooftxmanager.NewTx, qopts ...postgres.QOpt) (etx bulletprooftxmanager.EthTx, err error)
}

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

// ORM defines an interface for database commands related to Flux Monitor v2
type ORM interface {
	MostRecentFluxMonitorRoundID(aggregator common.Address) (uint32, error)
	DeleteFluxMonitorRoundsBackThrough(aggregator common.Address, roundID uint32) error
	FindOrCreateFluxMonitorRoundStats(aggregator common.Address, roundID uint32, newRoundLogs uint) (FluxMonitorRoundStatsV2, error)
	UpdateFluxMonitorRoundStats(aggregator common.Address, roundID uint32, runID int64, newRoundLogsAddition uint, qopts ...postgres.QOpt) error
	CreateEthTransaction(fromAddress, toAddress common.Address, payload []byte, gasLimit uint64) error
}

type orm struct {
	db       *gorm.DB
	txm      transmitter
	strategy bulletprooftxmanager.TxStrategy
}

// NewORM initializes a new ORM
func NewORM(db *gorm.DB, txm transmitter, strategy bulletprooftxmanager.TxStrategy) *orm {
	return &orm{db, txm, strategy}
}

// MostRecentFluxMonitorRoundID finds roundID of the most recent round that the
// provided oracle address submitted to
func (o *orm) MostRecentFluxMonitorRoundID(aggregator common.Address) (uint32, error) {
	var stats FluxMonitorRoundStatsV2
	err := o.db.
		Order("round_id DESC").
		First(&stats, "aggregator = ?", aggregator).
		Error
	if err != nil {
		return 0, err
	}

	return stats.RoundID, nil
}

// DeleteFluxMonitorRoundsBackThrough deletes all the RoundStat records for a
// given oracle address starting from the most recent round back through the
// given round
func (o *orm) DeleteFluxMonitorRoundsBackThrough(aggregator common.Address, roundID uint32) error {
	return o.db.Exec(`
        DELETE FROM flux_monitor_round_stats_v2
        WHERE aggregator = ?
          AND round_id >= ?
    `, aggregator, roundID).Error
}

// FindOrCreateFluxMonitorRoundStats find the round stats record for a given
// oracle on a given round, or creates it if no record exists
func (o *orm) FindOrCreateFluxMonitorRoundStats(aggregator common.Address, roundID uint32, newRoundLogs uint) (FluxMonitorRoundStatsV2, error) {

	// new potential entry to be inserted
	var stats FluxMonitorRoundStatsV2
	stats.Aggregator = aggregator
	stats.RoundID = roundID
	stats.NumNewRoundLogs = uint64(newRoundLogs)

	err := o.db.FirstOrCreate(&stats,
		// conditions for finding the existing one
		FluxMonitorRoundStatsV2{Aggregator: aggregator, RoundID: roundID},
	).Error

	return stats, err
}

// UpdateFluxMonitorRoundStats trys to create a RoundStat record for the given oracle
// at the given round. If one already exists, it increments the num_submissions column.
func (o *orm) UpdateFluxMonitorRoundStats(aggregator common.Address, roundID uint32, runID int64, newRoundLogsAddition uint, qopts ...postgres.QOpt) error {
	q := postgres.NewQ(postgres.UnwrapGormDB(o.db), qopts...)
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
func (o *orm) CountFluxMonitorRoundStats() (int, error) {
	var count int64
	err := o.db.Table("flux_monitor_round_stats_v2").Count(&count).Error

	return int(count), err
}

// CreateEthTransaction creates an ethereum transaction for the BPTXM to pick up
func (o *orm) CreateEthTransaction(
	fromAddress common.Address,
	toAddress common.Address,
	payload []byte,
	gasLimit uint64,
) (err error) {
	_, err = o.txm.CreateEthTransaction(bulletprooftxmanager.NewTx{
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: payload,
		GasLimit:       gasLimit,
		Meta:           nil,
		Strategy:       o.strategy,
	})
	return errors.Wrap(err, "Skipped Flux Monitor submission")
}
