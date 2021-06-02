package fluxmonitorv2

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"gorm.io/gorm"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

// ORM defines an interface for database commands related to Flux Monitor v2
type ORM interface {
	MostRecentFluxMonitorRoundID(aggregator common.Address) (uint32, error)
	DeleteFluxMonitorRoundsBackThrough(aggregator common.Address, roundID uint32) error
	FindOrCreateFluxMonitorRoundStats(aggregator common.Address, roundID uint32) (FluxMonitorRoundStatsV2, error)
	UpdateFluxMonitorRoundStats(db *gorm.DB, aggregator common.Address, roundID uint32, runID int64) error
	CreateEthTransaction(db *gorm.DB, fromAddress, toAddress common.Address, payload []byte, gasLimit uint64, maxUnconfirmedTransactions uint64) error
}

type orm struct {
	db *gorm.DB
}

// NewORM initializes a new ORM
func NewORM(db *gorm.DB) *orm {
	return &orm{db}
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
func (o *orm) FindOrCreateFluxMonitorRoundStats(aggregator common.Address, roundID uint32) (FluxMonitorRoundStatsV2, error) {
	var stats FluxMonitorRoundStatsV2
	err := o.db.FirstOrCreate(&stats,
		FluxMonitorRoundStatsV2{Aggregator: aggregator, RoundID: roundID},
	).Error

	return stats, err
}

// UpdateFluxMonitorRoundStats trys to create a RoundStat record for the given oracle
// at the given round. If one already exists, it increments the num_submissions column.
func (o *orm) UpdateFluxMonitorRoundStats(db *gorm.DB, aggregator common.Address, roundID uint32, runID int64) error {
	err := db.Exec(`
        INSERT INTO flux_monitor_round_stats_v2 (
            aggregator, round_id, pipeline_run_id, num_new_round_logs, num_submissions
        ) VALUES (
            ?, ?, ?, 0, 1
        ) ON CONFLICT (aggregator, round_id)
        DO UPDATE SET
					num_submissions = flux_monitor_round_stats_v2.num_submissions + 1,
					pipeline_run_id = EXCLUDED.pipeline_run_id
    `, aggregator, roundID, runID).Error
	return errors.Wrapf(err, "Failed to insert round stats for roundID=%v, runID=%v", roundID, runID)
}

// CountFluxMonitorRoundStats counts the total number of records
func (o *orm) CountFluxMonitorRoundStats() (int, error) {
	var count int64
	err := o.db.Table("flux_monitor_round_stats_v2").Count(&count).Error

	return int(count), err
}

// CreateEthTransaction creates an ethereum transaction for the BPTXM to pick up
func (o *orm) CreateEthTransaction(
	db *gorm.DB,
	fromAddress common.Address,
	toAddress common.Address,
	payload []byte,
	gasLimit uint64,
	maxUnconfirmedTransactions uint64,
) (err error) {
	_, err = bulletprooftxmanager.CreateEthTransaction(db, fromAddress, toAddress, payload, gasLimit, maxUnconfirmedTransactions)
	return errors.Wrap(err, "Skipped Flux Monitor submission")
}
