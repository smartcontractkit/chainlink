package keeper

import (
	"context"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	gasBuffer         = int32(200_000)
	maxUnconfirmedTXs = uint64(5)
)

func NewORM(db *gorm.DB) ORM {
	return ORM{
		DB: db,
	}
}

type ORM struct {
	DB *gorm.DB
}

func (korm ORM) Registries(ctx context.Context) (registries []Registry, _ error) {
	err := korm.DB.
		WithContext(ctx).
		Find(&registries).
		Error
	return registries, err
}

func (korm ORM) UpsertRegistry(ctx context.Context, registry *Registry) error {
	return korm.DB.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "job_id"}},
			DoUpdates: clause.AssignmentColumns(
				[]string{"keeper_index", "check_gas", "block_count_per_turn", "num_keepers"},
			),
		}).
		Create(registry).
		Error
}

func (korm ORM) UpsertUpkeep(ctx context.Context, registration *UpkeepRegistration) error {
	return korm.DB.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "registry_id"}, {Name: "upkeep_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"execute_gas", "check_data"}),
		}).
		Create(registration).
		Error
}

func (korm ORM) BatchDeleteUpkeeps(ctx context.Context, registryID int32, upkeedIDs []int64) error {
	return korm.DB.
		WithContext(ctx).
		Where("registry_id = ? AND upkeep_id IN (?)", registryID, upkeedIDs).
		Delete(UpkeepRegistration{}).
		Error
}

func (korm ORM) DeleteRegistryByJobID(ctx context.Context, jobID int32) error {
	return korm.DB.
		WithContext(ctx).
		Where("job_id = ?", jobID).
		Delete(Registry{}). // auto deletes upkeep registrations
		Error
}

func (korm ORM) EligibleUpkeeps(ctx context.Context, blockNumber int64) (upkeeps []UpkeepRegistration, _ error) {
	err := korm.DB.
		WithContext(ctx).
		Preload("Registry").
		Joins("INNER JOIN keeper_registries ON keeper_registries.id = upkeep_registrations.registry_id").
		Where(`
			? % keeper_registries.block_count_per_turn = 0 AND
			keeper_registries.keeper_index =
			(
				upkeep_registrations.positioning_constant + (? / keeper_registries.block_count_per_turn)
			) % keeper_registries.num_keepers
		`, blockNumber, blockNumber).
		Find(&upkeeps).
		Error

	return upkeeps, err
}

// LowestUnsyncedID returns the largest upkeepID + 1, indicating the expected next upkeepID
// to sync from the contract
func (korm ORM) LowestUnsyncedID(ctx context.Context, reg Registry) (nextID int64, err error) {
	err = korm.DB.
		WithContext(ctx).
		Model(&UpkeepRegistration{}).
		Where("registry_id = ?", reg.ID).
		Select("coalesce(max(upkeep_id), -1) + 1").
		Row().
		Scan(&nextID)
	return nextID, err
}

func (korm ORM) CreateEthTransactionForUpkeep(ctx context.Context, upkeep UpkeepRegistration, payload []byte) error {
	sqlDB, err := korm.DB.DB()
	if err != nil {
		return err
	}

	from := upkeep.Registry.FromAddress.Address()
	err = utils.CheckOKToTransmit(ctx, sqlDB, from, maxUnconfirmedTXs)
	if err != nil {
		return errors.Wrap(err, "transmitter#CreateEthTransaction")
	}

	value := 0
	res, err := sqlDB.ExecContext(ctx, `
		INSERT INTO eth_txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at)
		SELECT $1,$2,$3,$4,$5,'unstarted',NOW()
		WHERE NOT EXISTS (
			SELECT 1 FROM eth_tx_attempts
			JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id
			WHERE eth_txes.from_address = $1
				AND eth_txes.state = 'unconfirmed'
				AND eth_tx_attempts.state = 'insufficient_eth'
		);`,
		from,
		upkeep.Registry.ContractAddress.Address(),
		payload,
		value,
		upkeep.ExecuteGas+gasBuffer,
	)

	if err != nil {
		return errors.Wrap(err, "keeper failed to insert eth_tx")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "transmitter failed to get RowsAffected on eth_tx insert")
	}
	if rowsAffected == 0 {
		err := errors.Errorf("Skipped upkeep because wallet is out of eth: %s", from.Hex())
		logger.Warnw(err.Error())
		return err
	}

	return nil
}
