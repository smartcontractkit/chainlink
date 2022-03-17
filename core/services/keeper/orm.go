package keeper

import (
	"database/sql"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pg"
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
INSERT INTO keeper_registries (job_id, keeper_index, contract_address, from_address, check_gas, block_count_per_turn, num_keepers) VALUES (
:job_id, :keeper_index, :contract_address, :from_address, :check_gas, :block_count_per_turn, :num_keepers
) ON CONFLICT (job_id) DO UPDATE SET
	keeper_index = :keeper_index,
	check_gas = :check_gas,
	block_count_per_turn = :block_count_per_turn,
	num_keepers = :num_keepers
RETURNING *
`
	err := korm.q.GetNamed(stmt, registry, registry)
	return errors.Wrap(err, "failed to upsert registry")
}

// UpsertUpkeep upserts upkeep by the given input
func (korm ORM) UpsertUpkeep(registration *UpkeepRegistration) error {
	stmt := `
INSERT INTO upkeep_registrations (registry_id, execute_gas, check_data, upkeep_id, last_run_block_height) VALUES (
:registry_id, :execute_gas, :check_data, :upkeep_id, :last_run_block_height
) ON CONFLICT (registry_id, upkeep_id) DO UPDATE SET
	execute_gas = :execute_gas,
	check_data = :check_data
RETURNING *
`
	err := korm.q.GetNamed(stmt, registration, registration)
	return errors.Wrap(err, "failed to upsert upkeep")
}

// BatchDeleteUpkeepsForJob deletes all upkeeps by the given IDs for the job with the given ID
func (korm ORM) BatchDeleteUpkeepsForJob(jobID int32, upkeepIDs []int64) (int64, error) {
	res, err := korm.q.Exec(`
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

//EligibleUpkeepsForRegistry fetches eligible upkeeps for processing
//The query checks the following conditions
// - checks the registry address is correct and the registry has some keepers associated
// -- is it my turn AND my keeper was not the last perform for this upkeep OR my keeper was the last before BUT it is past the grace period
// -- OR is it my buddy's turn AND they were the last keeper to do the perform for this upkeep
func (korm ORM) EligibleUpkeepsForRegistry(registryAddress ethkey.EIP55Address, head *types.Head, gracePeriod int64) (upkeeps []UpkeepRegistration, err error) {
	registry, err := korm.RegistryByContractAddress(registryAddress)
	if err != nil {
		return nil, errors.Wrap(err, "EligibleUpkeepsForRegistry failed to get a registry by address")
	}
	blockNumber := head.Number
	binaryHash := binaryOfFirstHashInTurn(blockNumber, registry, head)

	stmt := `
SELECT upkeep_registrations.* FROM upkeep_registrations
INNER JOIN keeper_registries ON keeper_registries.id = upkeep_registrations.registry_id
WHERE
	keeper_registries.contract_address = $1 AND
	keeper_registries.num_keepers > 0 AND 
    ((
                keeper_registries.keeper_index = ((CAST(upkeep_registrations.upkeep_id AS bit(32)) #
                                                   CAST($4 AS bit(32)))::bigint % keeper_registries.num_keepers)
            AND
                (
				upkeep_registrations.last_keeper_index IS DISTINCT FROM keeper_registries.keeper_index
				OR
				(upkeep_registrations.last_keeper_index IS NOT DISTINCT FROM keeper_registries.keeper_index AND upkeep_registrations.last_run_block_height + $2 < $3)
				)
        )
   OR
    (
                    (keeper_registries.keeper_index + 1) % keeper_registries.num_keepers =
                    ((CAST(upkeep_registrations.upkeep_id AS bit(32)) #
                      CAST($4 AS bit(32)))::bigint % keeper_registries.num_keepers)
            AND
                    upkeep_registrations.last_keeper_index IS NOT DISTINCT FROM (keeper_registries.keeper_index + 1) % keeper_registries.num_keepers
        ))
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

// binaryOfFirstHashInTurn first calculates the first potential head for a turn. It then gets the hash for that head and converts it to binary
func binaryOfFirstHashInTurn(blockNumber int64, registry Registry, head *types.Head) string {
	firstHeadInTurn := blockNumber - (blockNumber % int64(registry.BlockCountPerTurn))
	hashAtHeight := head.HashAtHeight(firstHeadInTurn)
	bigInt := new(big.Int)
	bigInt.SetString(hashAtHeight.Hex(), 0)
	binaryString := fmt.Sprintf("%b", bigInt)
	return binaryString
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
	err = korm.q.Get(&nextID, `
SELECT coalesce(max(upkeep_id), -1) + 1
FROM upkeep_registrations
WHERE registry_id = $1
`, regID)
	return nextID, errors.Wrap(err, "LowestUnsyncedID failed")
}

//SetLastRunInfoForUpkeepOnJob sets the last run block height and the associated keeper index only if the new block height is greater than the previous.
func (korm ORM) SetLastRunInfoForUpkeepOnJob(jobID int32, upkeepID, height int64, fromIndex sql.NullInt64, qopts ...pg.QOpt) error {
	_, err := korm.q.WithOpts(qopts...).Exec(`
UPDATE upkeep_registrations
SET last_run_block_height = $1,
    last_keeper_index = $4
WHERE upkeep_id = $2 AND
registry_id = (SELECT id FROM keeper_registries WHERE job_id = $3) AND
last_run_block_height < $1`, height, upkeepID, jobID, fromIndex)
	return errors.Wrap(err, "SetLastRunInfoForUpkeepOnJob failed")
}
