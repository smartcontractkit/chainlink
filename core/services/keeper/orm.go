package keeper

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
)

// ORM implements ORM layer using PostgreSQL
type ORM struct {
	DB       *gorm.DB
	txm      transmitter
	config   Config
	strategy bulletprooftxmanager.TxStrategy
}

// NewORM is the constructor of postgresORM
func NewORM(db *gorm.DB, txm transmitter, config Config, strategy bulletprooftxmanager.TxStrategy) ORM {
	return ORM{
		DB:       db,
		txm:      txm,
		config:   config,
		strategy: strategy,
	}
}

// WithTransaction wraps the given handler into transaction context
func (korm ORM) WithTransaction(cb func(ctx context.Context) error) error {
	return postgres.NewGormTransactionManager(korm.DB).Transact(cb)
}

func (korm ORM) Registries(ctx context.Context) (registries []Registry, _ error) {
	err := korm.getDB(ctx).
		Find(&registries).
		Error
	return registries, err
}

func (korm ORM) RegistryForJob(ctx context.Context, jobID int32) (registry Registry, _ error) {
	err := korm.getDB(ctx).
		First(&registry, "job_id = ?", jobID).
		Error
	return registry, err
}

func (korm ORM) UpsertRegistry(ctx context.Context, registry *Registry) error {
	return korm.getDB(ctx).
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
	return korm.getDB(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "registry_id"}, {Name: "upkeep_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"execute_gas", "check_data", "positioning_constant"}),
		}).
		Create(registration).
		Error
}

func (korm ORM) BatchDeleteUpkeepsForJob(ctx context.Context, jobID int32, upkeedIDs []int64) (int64, error) {
	exec := korm.getDB(ctx).
		Exec(
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
	blockNumber, gracePeriod int64,
) (upkeeps []UpkeepRegistration, _ error) {
	err := korm.getDB(ctx).
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
func (korm ORM) LowestUnsyncedID(ctx context.Context, regID int32) (nextID int64, err error) {
	err = korm.getDB(ctx).
		Model(&UpkeepRegistration{}).
		Where("registry_id = ?", regID).
		Select("coalesce(max(upkeep_id), -1) + 1").
		Row().
		Scan(&nextID)
	return nextID, err
}

func (korm ORM) SetLastRunHeightForUpkeepOnJob(ctx context.Context, jobID int32, upkeepID, height int64) error {
	return korm.getDB(ctx).
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

func (korm ORM) CreateEthTransactionForUpkeep(ctx context.Context, upkeep UpkeepRegistration, payload []byte) (bulletprooftxmanager.EthTx, error) {
	from := upkeep.Registry.FromAddress.Address()
	to := upkeep.Registry.ContractAddress.Address()
	gasLimit := upkeep.ExecuteGas + korm.config.KeeperRegistryPerformGasOverhead()
	return korm.txm.CreateEthTransaction(korm.getDB(ctx), bulletprooftxmanager.NewTx{
		FromAddress:    from,
		ToAddress:      to,
		EncodedPayload: payload,
		GasLimit:       gasLimit,
		Meta:           nil,
		Strategy:       korm.strategy,
	})
}

func (korm ORM) getDB(ctx context.Context) *gorm.DB {
	return postgres.TxFromContext(ctx, korm.DB).WithContext(ctx)
}
