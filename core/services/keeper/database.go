package keeper

import (
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"gorm.io/gorm/clause"
)

type DB interface {
	Registries() ([]Registry, error)
	UpsertRegistry(registry Registry) error
	UpsertUpkeep(UpkeepRegistration) error
	BatchDeleteUpkeeps(registryID int32, upkeedIDs []int64) error
	DeleteRegistryByJobID(jobID int32) error
	EligibleUpkeeps(blockNumber int64) ([]UpkeepRegistration, error)
	NextUpkeepIDForRegistry(registry Registry) (int64, error)
}

func NewDBInterface(orm *orm.ORM) DB {
	return keeperDB{
		orm,
	}
}

type keeperDB struct {
	*orm.ORM
}

func (kdb keeperDB) Registries() (registries []Registry, _ error) {
	err := kdb.DB.Find(&registries).Error
	return registries, err
}

func (kdb keeperDB) UpsertRegistry(registry Registry) error {
	return kdb.DB.Save(&registry).Error
}

func (kdb keeperDB) UpsertUpkeep(registration UpkeepRegistration) error {
	return kdb.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "registry_id"}, {Name: "upkeep_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"execute_gas", "check_data"}),
		}).
		Create(&registration).
		Error
}

func (kdb keeperDB) BatchDeleteUpkeeps(registryID int32, upkeedIDs []int64) error {
	return kdb.DB.
		Where("registry_id = ? AND upkeep_id IN (?)", registryID, upkeedIDs).
		Delete(UpkeepRegistration{}).
		Error
}

func (kdb keeperDB) DeleteRegistryByJobID(jobID int32) error {
	return kdb.DB.
		Where("job_id = ?", jobID).
		Delete(Registry{}). // auto deletes upkeep registrations
		Error
}

func (kdb keeperDB) EligibleUpkeeps(blockNumber int64) (result []UpkeepRegistration, _ error) {
	turnTakingQuery := `
		keeper_registries.keeper_index =
			(
				upkeep_registrations.positioning_constant + (? / keeper_registries.block_count_per_turn)
			) % keeper_registries.num_keepers
	`
	err := kdb.DB.
		Preload("Registry").
		Joins("INNER JOIN keeper_registries ON keeper_registries.id = upkeep_registrations.registry_id").
		Where("? % keeper_registries.block_count_per_turn = 0", blockNumber).
		Where(turnTakingQuery, blockNumber).
		Find(&result).
		Error

	return result, err
}

// NextUpkeepIDForRegistry returns the largest upkeepID + 1, indicating the expected next upkeepID
// to sync from the contract
func (kdb keeperDB) NextUpkeepIDForRegistry(reg Registry) (nextID int64, err error) {
	err = kdb.DB.
		Model(&UpkeepRegistration{}).
		Where("registry_id = ?", reg.ID).
		Select("coalesce(max(upkeep_id), -1) + 1").
		Row().
		Scan(&nextID)
	return nextID, err
}
