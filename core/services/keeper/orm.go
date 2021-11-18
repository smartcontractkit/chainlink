package keeper

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/sqlx"
)

// ORM implements ORM layer using PostgreSQL
type ORM struct {
	DB       *sqlx.DB
	txm      transmitter
	config   Config
	strategy bulletprooftxmanager.TxStrategy
	logger   logger.Logger
}

// NewORM is the constructor of postgresORM
func NewORM(db *sqlx.DB, lggr logger.Logger, txm transmitter, config Config, strategy bulletprooftxmanager.TxStrategy) ORM {
	lggr = lggr.Named("KeeperORM")
	return ORM{
		DB:       db,
		txm:      txm,
		config:   config,
		strategy: strategy,
		logger:   lggr,
	}
}

// Registries returns all registries
func (korm ORM) Registries() ([]Registry, error) {
	var registries []Registry
	err := pg.NewQ(korm.DB).Select(&registries, `SELECT * FROM keeper_registries ORDER BY id ASC`)
	return registries, errors.Wrap(err, "failed to get registries")
}

// RegistryForJob returns a specific registry for a job with the given ID
func (korm ORM) RegistryForJob(jobID int32) (Registry, error) {
	var registry Registry
	err := pg.NewQ(korm.DB).Get(&registry, `SELECT * FROM keeper_registries WHERE job_id = $1 LIMIT 1`, jobID)
	return registry, errors.Wrapf(err, "failed to get registry with job_id %d", jobID)
}

// UpsertRegistry upserts registry by the given input
func (korm ORM) UpsertRegistry(registry *Registry) error {
	stmt := `
INSERT INTO keeper_registries (job_id, keeper_index, contract_address, from_address, check_gas, block_count_per_turn, num_keepers) VALUES (
:job_id, :keeper_index, :contract_address, :from_address, :check_gas, :block_count_per_turn, :num_keepers
) ON CONFLICT (job_id) DO UPDATE SET
	keeper_index = :keeper_index,
	check_gas = :check_gas,
	block_count_per_turn = :block_count_per_turn,
	num_keepers = :num_keepers
RETURNING *
`
	err := pg.NewQ(korm.DB).GetNamed(stmt, registry, registry)
	return errors.Wrap(err, "failed to upsert registry")
}

// UpsertUpkeep upserts upkeep by the given input
func (korm ORM) UpsertUpkeep(registration *UpkeepRegistration) error {
	stmt := `
INSERT INTO upkeep_registrations (registry_id, execute_gas, check_data, upkeep_id, positioning_constant, last_run_block_height) VALUES (
:registry_id, :execute_gas, :check_data, :upkeep_id, :positioning_constant, :last_run_block_height
) ON CONFLICT (registry_id, upkeep_id) DO UPDATE SET
	execute_gas = :execute_gas,
	check_data = :check_data,
	positioning_constant = :positioning_constant
RETURNING *
`
	err := pg.NewQ(korm.DB).GetNamed(stmt, registration, registration)
	return errors.Wrap(err, "failed to upsert upkeep")
}

// BatchDeleteUpkeepsForJob deletes all upkeeps by the given IDs for the job with the given ID
func (korm ORM) BatchDeleteUpkeepsForJob(jobID int32, upkeepIDs []int64) (int64, error) {
	res, err := pg.NewQ(korm.DB).Exec(`
DELETE FROM upkeep_registrations WHERE registry_id IN (
	SELECT id FROM keeper_registries WHERE job_id = $1
) AND upkeep_id = ANY($2)
`, jobID, upkeepIDs)
	if err != nil {
		return 0, errors.Wrap(err, "BatchDeleteUpkeepsForJob failed to delete")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "BatchDeleteUpkeepsForJob failed to get RowsAffected")
	}
	return rowsAffected, nil
}

func (korm ORM) EligibleUpkeepsForRegistry(
	registryAddress ethkey.EIP55Address,
	blockNumber, gracePeriod int64,
) (upkeeps []UpkeepRegistration, err error) {
	err = pg.NewQ(korm.DB).Transaction(korm.logger, func(tx pg.Queryer) error {
		stmt := `
SELECT upkeep_registrations.* FROM upkeep_registrations
INNER JOIN keeper_registries ON keeper_registries.id = upkeep_registrations.registry_id
WHERE
	keeper_registries.contract_address = $1 AND
	keeper_registries.num_keepers > 0 AND
	(
		upkeep_registrations.last_run_block_height = 0 OR (
			upkeep_registrations.last_run_block_height + $2 < $3 AND
			upkeep_registrations.last_run_block_height < ($3 - ($3 % keeper_registries.block_count_per_turn))
		)
	) AND
	keeper_registries.keeper_index = (
		upkeep_registrations.positioning_constant + (($3 - ($3 % keeper_registries.block_count_per_turn)) / keeper_registries.block_count_per_turn)
	) % keeper_registries.num_keepers
ORDER BY upkeep_registrations.id ASC, upkeep_registrations.upkeep_id ASC
`
		if err = tx.Select(&upkeeps, stmt, registryAddress, gracePeriod, blockNumber); err != nil {
			return errors.Wrap(err, "EligibleUpkeepsForRegistry failed to get upkeep_registrations")
		}
		if err = loadUpkeepsRegistry(tx, upkeeps); err != nil {
			return errors.Wrap(err, "EligibleUpkeepsForRegistry failed to load Registry on upkeeps")
		}
		return nil
	}, pg.OptReadOnlyTx())

	return upkeeps, err
}

func loadUpkeepsRegistry(q pg.Queryer, upkeeps []UpkeepRegistration) error {
	registryIDM := make(map[int64]*Registry)
	var registryIDs []int64
	for _, upkeep := range upkeeps {
		if _, exists := registryIDM[upkeep.RegistryID]; !exists {
			registryIDM[upkeep.RegistryID] = new(Registry)
			registryIDs = append(registryIDs, upkeep.RegistryID)
		}
	}
	var registries []*Registry
	err := q.Select(&registries, `SELECT * FROM keeper_registries WHERE id = ANY($1)`, pq.Array(registryIDs))
	if err != nil {
		return errors.Wrap(err, "loadUpkeepsRegistry failed")
	}
	for _, reg := range registries {
		registryIDM[reg.ID] = reg
	}
	for i, upkeep := range upkeeps {
		upkeeps[i].Registry = *registryIDM[upkeep.RegistryID]
	}
	return nil
}

// LowestUnsyncedID returns the largest upkeepID + 1, indicating the expected next upkeepID
// to sync from the contract
func (korm ORM) LowestUnsyncedID(regID int64) (nextID int64, err error) {
	err = pg.NewQ(korm.DB).Get(&nextID, `
SELECT coalesce(max(upkeep_id), -1) + 1
FROM upkeep_registrations
WHERE registry_id = $1
`, regID)
	return nextID, errors.Wrap(err, "LowestUnsyncedID failed")
}

func (korm ORM) SetLastRunHeightForUpkeepOnJob(jobID int32, upkeepID, height int64, qopts ...pg.QOpt) error {
	_, err := pg.NewQ(korm.DB, qopts...).Exec(`
UPDATE upkeep_registrations
SET last_run_block_height = $1
WHERE upkeep_id = $2 AND
registry_id = (
	SELECT id FROM keeper_registries WHERE job_id = $3
)`, height, upkeepID, jobID)
	return errors.Wrap(err, "SetLastRunHeightForUpkeepOnJob failed")
}
