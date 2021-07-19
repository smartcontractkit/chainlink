package keeper

import (
	"context"

	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewORM(db *gorm.DB, txm transmitter, config Config, strategy bulletprooftxmanager.TxStrategy) ORM {
	return ORM{
		DB:       db,
		txm:      txm,
		config:   config,
		strategy: strategy,
	}
}

type ORM struct {
	DB       *gorm.DB
	txm      transmitter
	config   Config
	strategy bulletprooftxmanager.TxStrategy
}

func (korm ORM) Registries(ctx context.Context) (registries []Registry, _ error) {
	err := korm.DB.
		WithContext(ctx).
		Find(&registries).
		Error
	return registries, err
}

func (korm ORM) RegistryForJob(ctx context.Context, jobID int32) (registry Registry, _ error) {
	err := korm.DB.
		WithContext(ctx).
		First(&registry, "job_id = ?", jobID).
		Error
	return registry, err
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
			DoUpdates: clause.AssignmentColumns([]string{"execute_gas", "check_data", "positioning_constant"}),
		}).
		Create(registration).
		Error
}

func (korm ORM) BatchDeleteUpkeepsForJob(ctx context.Context, jobID int32, upkeedIDs []int64) (int64, error) {
	exec := korm.DB.
		WithContext(ctx).Exec(
		`DELETE FROM upkeep_registrations WHERE registry_id = (
			SELECT id from keeper_registries where job_id = ?
		) AND upkeep_id IN (?)`,
		jobID,
		upkeedIDs,
	)
	return exec.RowsAffected, exec.Error
}

func (korm ORM) EligibleUpkeepsForRegistry(
	ctx context.Context,
	registryAddress ethkey.EIP55Address,
	blockNumber int64,
	gracePeriod int64,
) (upkeeps []UpkeepRegistration, _ error) {
	err := korm.DB.
		WithContext(ctx).
		Preload("Registry").
		Order("upkeep_registrations.id ASC, upkeep_registrations.upkeep_id ASC").
		Joins("INNER JOIN keeper_registries ON keeper_registries.id = upkeep_registrations.registry_id").
		Where(`
			keeper_registries.contract_address = ? AND
			keeper_registries.num_keepers > 0 AND
			(
				upkeep_registrations.last_run_block_height = 0 OR (
					upkeep_registrations.last_run_block_height + ? < ? AND
					upkeep_registrations.last_run_block_height < (? - (? % keeper_registries.block_count_per_turn))
				)
			) AND
			keeper_registries.keeper_index = (
				upkeep_registrations.positioning_constant + ((? - (? % keeper_registries.block_count_per_turn)) / keeper_registries.block_count_per_turn)
			) % keeper_registries.num_keepers
		`, registryAddress, gracePeriod, blockNumber, blockNumber, blockNumber, blockNumber, blockNumber).
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

func (korm ORM) SetLastRunHeightForUpkeepOnJob(db *gorm.DB, jobID int32, upkeepID int64, height int64) error {
	return db.
		Exec(`UPDATE upkeep_registrations
		SET last_run_block_height = ?
		WHERE upkeep_id = ? AND
		registry_id = (
			SELECT id FROM keeper_registries WHERE job_id = ?
		);`,
			height,
			upkeepID,
			jobID,
		).Error
}

func (korm ORM) CreateEthTransactionForUpkeep(tx *gorm.DB, upkeep UpkeepRegistration, payload []byte) (bulletprooftxmanager.EthTx, error) {
	from := upkeep.Registry.FromAddress.Address()
	to := upkeep.Registry.ContractAddress.Address()
	gasLimit := upkeep.ExecuteGas + korm.config.KeeperRegistryPerformGasOverhead()
	return korm.txm.CreateEthTransaction(tx, from, to, payload, gasLimit, nil, korm.strategy)
}
