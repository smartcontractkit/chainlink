package keeper

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const gasBuffer = int32(200_000)

func NewORM(db *gorm.DB) ORM {
	return ORM{
		DB: db,
	}
}

type ORM struct {
	DB *gorm.DB
}

func (korm ORM) Registries() (registries []Registry, _ error) {
	err := korm.DB.Find(&registries).Error
	return registries, err
}

func (korm ORM) UpsertRegistry(registry *Registry) error {
	return korm.DB.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "job_id"}},
			DoUpdates: clause.AssignmentColumns(
				[]string{"keeper_index", "check_gas", "block_count_per_turn", "num_keepers"},
			),
		}).
		Create(registry).
		Error
}

func (korm ORM) UpsertUpkeep(registration *UpkeepRegistration) error {
	return korm.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "registry_id"}, {Name: "upkeep_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"execute_gas", "check_data"}),
		}).
		Create(registration).
		Error
}

func (korm ORM) BatchDeleteUpkeeps(registryID int32, upkeedIDs []int64) error {
	return korm.DB.
		Where("registry_id = ? AND upkeep_id IN (?)", registryID, upkeedIDs).
		Delete(UpkeepRegistration{}).
		Error
}

// TODO - RYAN - raw sql
func (korm ORM) DeleteRegistryByJobID(jobID int32) error {
	return korm.DB.
		Where("job_id = ?", jobID).
		Delete(Registry{}). // auto deletes upkeep registrations
		Error
}

// TODO - RYAN - one where clause
func (korm ORM) EligibleUpkeeps(blockNumber int64) (upkeeps []UpkeepRegistration, _ error) {
	turnTakingQuery := `
		keeper_registries.keeper_index =
			(
				upkeep_registrations.positioning_constant + (? / keeper_registries.block_count_per_turn)
			) % keeper_registries.num_keepers
	`
	err := korm.DB.
		Preload("Registry").
		Joins("INNER JOIN keeper_registries ON keeper_registries.id = upkeep_registrations.registry_id").
		Where("? % keeper_registries.block_count_per_turn = 0", blockNumber).
		Where(turnTakingQuery, blockNumber).
		Find(&upkeeps).
		Error

	return upkeeps, err
}

// NextUpkeepIDForRegistry returns the largest upkeepID + 1, indicating the expected next upkeepID
// to sync from the contract
// LowestUnsyncedID
// todo - ryan - note not racy - and raw sql
func (korm ORM) NextUpkeepIDForRegistry(reg Registry) (nextID int64, err error) {
	err = korm.DB.
		Model(&UpkeepRegistration{}).
		Where("registry_id = ?", reg.ID).
		Select("coalesce(max(upkeep_id), -1) + 1").
		Row().
		Scan(&nextID)
	return nextID, err
}

func (korm ORM) CreateEthTransactionForUpkeep(upkeep UpkeepRegistration, payload []byte) error {
	sqlDB, err := korm.DB.DB()
	if err != nil {
		return err
	}
	_, err = sqlDB.Exec(
		`INSERT INTO eth_txes (from_address, to_address, encoded_payload, gas_limit, value, state, created_at)
		VALUES ($1,$2,$3,$4,0,'unstarted',NOW());`,
		upkeep.Registry.FromAddress.Address(),
		upkeep.Registry.ContractAddress.Address(),
		payload,
		upkeep.ExecuteGas+gasBuffer,
	)
	if err != nil {
		return errors.Wrap(err, "keeper failed to insert eth_tx")
	}
	return nil
}
