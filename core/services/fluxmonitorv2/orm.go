package fluxmonitorv2

import (
	"context"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

// ORM defines an interface for database commands related to Flux Monitor v2
type ORM interface {
	MostRecentFluxMonitorRoundID(aggregator common.Address) (uint32, error)
	DeleteFluxMonitorRoundsBackThrough(aggregator common.Address, roundID uint32) error
	FindOrCreateFluxMonitorRoundStats(aggregator common.Address, roundID uint32) (FluxMonitorRoundStatsV2, error)
	UpdateFluxMonitorRoundStats(aggregator common.Address, roundID uint32, runID int64) error
	CreateEthTransaction(fromAddress, toAddress common.Address, payload []byte, gasLimit uint64, maxUnconfirmedTransactions uint64) error
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
func (o *orm) UpdateFluxMonitorRoundStats(aggregator common.Address, roundID uint32, runID int64) error {
	return o.db.Exec(`
        INSERT INTO flux_monitor_round_stats_v2 (
            aggregator, round_id, pipeline_run_id, num_new_round_logs, num_submissions
        ) VALUES (
            ?, ?, ?, 0, 1
        ) ON CONFLICT (aggregator, round_id)
        DO UPDATE SET
					num_submissions = flux_monitor_round_stats_v2.num_submissions + 1,
					pipeline_run_id = EXCLUDED.pipeline_run_id
    `, aggregator, roundID, runID).Error
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
	maxUnconfirmedTransactions uint64,
) error {
	db, err := o.db.DB()
	if err != nil {
		return errors.Wrap(err, "orm#CreateEthTransaction")
	}

	err = utils.CheckOKToTransmit(context.Background(), db, fromAddress, maxUnconfirmedTransactions)
	if err != nil {
		return errors.Wrap(err, "orm#CreateEthTransaction")
	}

	value := 0

	dbtx := o.db.Exec(`
INSERT INTO eth_txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at)
SELECT ?,?,?,?,?,'unstarted',NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM eth_tx_attempts
	JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id
	WHERE eth_txes.from_address = ?
		AND eth_txes.state = 'unconfirmed'
		AND eth_tx_attempts.state = 'insufficient_eth'
);
`, fromAddress, toAddress, payload, value, gasLimit, fromAddress)
	if dbtx.Error != nil {
		return errors.Wrap(dbtx.Error, "failed to insert eth_tx")
	}
	if dbtx.RowsAffected == 0 {
		// Unsure why this would be an wallet out of eth error
		// TODO - What is this error message
		err := errors.Errorf("Skipped Flux Monitor submission because wallet is out of eth: %s", fromAddress.Hex())
		logger.Warnw(err.Error(),
			"fromAddress", fromAddress,
			"toAddress", toAddress,
			"payload", "0x"+hex.EncodeToString(payload),
			"value", value,
			"gasLimit", gasLimit,
		)

		return err
	}

	return nil
}
