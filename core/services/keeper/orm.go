package keeper

import (
	"math/rand"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ORM implements ORM layer using PostgreSQL
type ORM struct {
	q        pg.Q
	config   Config
	strategy txmgr.TxStrategy
	logger   logger.Logger
}

// NewORM is the constructor of postgresORM
func NewORM(db *sqlx.DB, lggr logger.Logger, config Config, strategy txmgr.TxStrategy) ORM {
	lggr = lggr.Named("KeeperORM")
	return ORM{
		q:        pg.NewQ(db, lggr, config),
		config:   config,
		strategy: strategy,
		logger:   lggr,
	}
}

func (korm ORM) Q() pg.Q {
	return korm.q
}

// Registries returns all registries
func (korm ORM) Registries() ([]Registry, error) {
	var registries []Registry
	err := korm.q.Select(&registries, `SELECT * FROM keeper_registries ORDER BY id ASC`)
	return registries, errors.Wrap(err, "failed to get registries")
}

// RegistryByContractAddress returns a single registry based on provided address
func (korm ORM) RegistryByContractAddress(registryAddress ethkey.EIP55Address) (Registry, error) {
	var registry Registry
	err := korm.q.Get(&registry, `SELECT * FROM keeper_registries WHERE keeper_registries.contract_address = $1`, registryAddress)
	return registry, errors.Wrap(err, "failed to get registry")
}

// RegistryForJob returns a specific registry for a job with the given ID
func (korm ORM) RegistryForJob(jobID int32) (Registry, error) {
	var registry Registry
	err := korm.q.Get(&registry, `SELECT * FROM keeper_registries WHERE job_id = $1 LIMIT 1`, jobID)
	return registry, errors.Wrapf(err, "failed to get registry with job_id %d", jobID)
}

// UpsertRegistry upserts registry by the given input
func (korm ORM) UpsertRegistry(registry *Registry) error {
	stmt := `
INSERT INTO keeper_registries (job_id, keeper_index, contract_address, from_address, check_gas, block_count_per_turn, num_keepers, keeper_index_map) VALUES (
:job_id, :keeper_index, :contract_address, :from_address, :check_gas, :block_count_per_turn, :num_keepers, :keeper_index_map
) ON CONFLICT (job_id) DO UPDATE SET
	keeper_index = :keeper_index,
	check_gas = :check_gas,
	block_count_per_turn = :block_count_per_turn,
	num_keepers = :num_keepers,
	keeper_index_map = :keeper_index_map
RETURNING *
`
	err := korm.q.GetNamed(stmt, registry, registry)
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
	err := korm.q.GetNamed(stmt, registration, registration)
	return errors.Wrap(err, "failed to upsert upkeep")
}

// UpdateUpkeepLastKeeperIndex updates the last keeper index for an upkeep
func (korm ORM) UpdateUpkeepLastKeeperIndex(jobID int32, upkeepID *utils.Big, fromAddress ethkey.EIP55Address) error {
	_, err := korm.q.Exec(`
	UPDATE upkeep_registrations
	SET
		last_keeper_index = CAST((SELECT keeper_index_map -> $3 FROM keeper_registries WHERE job_id = $1) as int)
	WHERE upkeep_id = $2 AND
	registry_id = (SELECT id FROM keeper_registries WHERE job_id = $1)`,
		jobID, upkeepID, fromAddress.Hex())
	return errors.Wrap(err, "UpdateUpkeepLastKeeperIndex failed")
}

// BatchDeleteUpkeepsForJob deletes all upkeeps by the given IDs for the job with the given ID
func (korm ORM) BatchDeleteUpkeepsForJob(jobID int32, upkeepIDs []utils.Big) (int64, error) {
	strIds := []string{}
	for _, upkeepID := range upkeepIDs {
		strIds = append(strIds, upkeepID.String())
	}
	res, err := korm.q.Exec(`
DELETE FROM upkeep_registrations WHERE registry_id IN (
	SELECT id FROM keeper_registries WHERE job_id = $1
) AND upkeep_id = ANY($2)
`, jobID, strIds)
	if err != nil {
		return 0, errors.Wrap(err, "BatchDeleteUpkeepsForJob failed to delete")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "BatchDeleteUpkeepsForJob failed to get RowsAffected")
	}
	return rowsAffected, nil
}

// NewEligibleUpkeepsForRegistry fetches eligible upkeeps for processing
//The query checks the following conditions
// - checks the registry address is correct and the registry has some keepers associated
// -- is it my turn AND my keeper was not the last perform for this upkeep OR my keeper was the last before BUT it is past the grace period
// -- OR is it my buddy's turn AND they were the last keeper to do the perform for this upkeep
// DEV: note we cast upkeep_id and binaryHash as 32 bits, even though both are 256 bit numbers when performing XOR. This is enough information
// to disribute the upkeeps over the keepers so long as num keepers < 4294967296
func (korm ORM) NewEligibleUpkeepsForRegistry(registryAddress ethkey.EIP55Address, blockNumber int64, gracePeriod int64, binaryHash string) (upkeeps []UpkeepRegistration, err error) {
	stmt := `
SELECT upkeep_registrations.*
FROM upkeep_registrations
  INNER JOIN keeper_registries ON keeper_registries.id = upkeep_registrations.registry_id,
  LATERAL ABS(
		(least_significant(uint256_to_bit(upkeep_registrations.upkeep_id), 32) # least_significant($4, 32))::bigint
	) AS turn
WHERE keeper_registries.contract_address = $1
  AND keeper_registries.num_keepers > 0
  AND
		(
			(
				-- my turn
				keeper_registries.keeper_index = turn % keeper_registries.num_keepers
				AND
					(
						upkeep_registrations.last_keeper_index IS DISTINCT FROM keeper_registries.keeper_index
						OR
						(upkeep_registrations.last_keeper_index IS NOT DISTINCT FROM
							keeper_registries.keeper_index AND
							upkeep_registrations.last_run_block_height + $2 < $3)
					)
			)
			OR
			(
				-- my buddy's turn
				(keeper_registries.keeper_index + 1) % keeper_registries.num_keepers =
					turn % keeper_registries.num_keepers
				AND
				upkeep_registrations.last_keeper_index IS NOT DISTINCT FROM
					(keeper_registries.keeper_index + 1) % keeper_registries.num_keepers
			)
		)
`
	if err = korm.q.Select(&upkeeps, stmt, registryAddress, gracePeriod, blockNumber, binaryHash); err != nil {
		return upkeeps, errors.Wrap(err, "EligibleUpkeepsForRegistry failed to get upkeep_registrations")
	}
	if err = loadUpkeepsRegistry(korm.q, upkeeps); err != nil {
		return upkeeps, errors.Wrap(err, "EligibleUpkeepsForRegistry failed to load Registry on upkeeps")
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(upkeeps), func(i, j int) {
		upkeeps[i], upkeeps[j] = upkeeps[j], upkeeps[i]
	})

	return upkeeps, err
}

func (korm ORM) EligibleUpkeepsForRegistry(registryAddress ethkey.EIP55Address, blockNumber, gracePeriod int64) (upkeeps []UpkeepRegistration, err error) {
	err = korm.q.Transaction(func(tx pg.Queryer) error {
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

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(upkeeps), func(i, j int) {
		upkeeps[i], upkeeps[j] = upkeeps[j], upkeeps[i]
	})

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

func (korm ORM) AllUpkeepIDsForRegistry(regID int64) (upkeeps []utils.Big, err error) {
	err = korm.q.Select(&upkeeps, `
SELECT upkeep_id
FROM upkeep_registrations
WHERE registry_id = $1
`, regID)
	return upkeeps, errors.Wrap(err, "allUpkeepIDs failed")
}

//SetLastRunInfoForUpkeepOnJob sets the last run block height and the associated keeper index only if the new block height is greater than the previous.
func (korm ORM) SetLastRunInfoForUpkeepOnJob(jobID int32, upkeepID *utils.Big, height int64, fromAddress ethkey.EIP55Address, qopts ...pg.QOpt) error {
	_, err := korm.q.WithOpts(qopts...).Exec(`
	UPDATE upkeep_registrations
	SET last_run_block_height = $1,
		last_keeper_index = CAST((SELECT keeper_index_map -> $4 FROM keeper_registries WHERE job_id = $3) as int)
	WHERE upkeep_id = $2 AND
	registry_id = (SELECT id FROM keeper_registries WHERE job_id = $3) AND
	last_run_block_height < $1`, height, upkeepID, jobID, fromAddress.Hex())
	return errors.Wrap(err, "SetLastRunInfoForUpkeepOnJob failed")
}
