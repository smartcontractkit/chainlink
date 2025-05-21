package keeper

import (
	"context"
	"math/rand"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// ORM implements ORM layer using PostgreSQL
type ORM struct {
	ds     sqlutil.DataSource
	logger logger.Logger
}

// NewORM is the constructor of postgresORM
func NewORM(ds sqlutil.DataSource, lggr logger.Logger) *ORM {
	lggr = lggr.Named("KeeperORM")
	return &ORM{
		ds:     ds,
		logger: lggr,
	}
}

func (o *ORM) DataSource() sqlutil.DataSource {
	return o.ds
}

// Registries returns all registries
func (o *ORM) Registries(ctx context.Context) ([]Registry, error) {
	var registries []Registry
	err := o.ds.SelectContext(ctx, &registries, `SELECT * FROM keeper_registries ORDER BY id ASC`)
	return registries, errors.Wrap(err, "failed to get registries")
}

// RegistryByContractAddress returns a single registry based on provided address
func (o *ORM) RegistryByContractAddress(ctx context.Context, registryAddress types.EIP55Address) (Registry, error) {
	var registry Registry
	err := o.ds.GetContext(ctx, &registry, `SELECT * FROM keeper_registries WHERE keeper_registries.contract_address = $1`, registryAddress)
	return registry, errors.Wrap(err, "failed to get registry")
}

// RegistryForJob returns a specific registry for a job with the given ID
func (o *ORM) RegistryForJob(ctx context.Context, jobID int32) (Registry, error) {
	var registry Registry
	err := o.ds.GetContext(ctx, &registry, `SELECT * FROM keeper_registries WHERE job_id = $1 LIMIT 1`, jobID)
	return registry, errors.Wrapf(err, "failed to get registry with job_id %d", jobID)
}

// UpsertRegistry upserts registry by the given input
func (o *ORM) UpsertRegistry(ctx context.Context, registry *Registry) error {
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
	query, args, err := o.ds.BindNamed(stmt, registry)
	if err != nil {
		return errors.Wrap(err, "failed to upsert registry")
	}
	err = o.ds.GetContext(ctx, registry, query, args...)
	return errors.Wrap(err, "failed to upsert registry")
}

// UpsertUpkeep upserts upkeep by the given input
func (o *ORM) UpsertUpkeep(ctx context.Context, registration *UpkeepRegistration) error {
	stmt := `
INSERT INTO upkeep_registrations (registry_id, execute_gas, check_data, upkeep_id, positioning_constant, last_run_block_height) VALUES (
:registry_id, :execute_gas, :check_data, :upkeep_id, :positioning_constant, :last_run_block_height
) ON CONFLICT (registry_id, upkeep_id) DO UPDATE SET
	execute_gas = :execute_gas,
	check_data = :check_data,
	positioning_constant = :positioning_constant
RETURNING *
`
	query, args, err := o.ds.BindNamed(stmt, registration)
	if err != nil {
		return errors.Wrap(err, "failed to upsert upkeep")
	}
	err = o.ds.GetContext(ctx, registration, query, args...)
	return errors.Wrap(err, "failed to upsert upkeep")
}

// UpdateUpkeepLastKeeperIndex updates the last keeper index for an upkeep
func (o *ORM) UpdateUpkeepLastKeeperIndex(ctx context.Context, jobID int32, upkeepID *big.Big, fromAddress types.EIP55Address) error {
	_, err := o.ds.ExecContext(ctx, `
	UPDATE upkeep_registrations
	SET
		last_keeper_index = CAST((SELECT keeper_index_map -> $3 FROM keeper_registries WHERE job_id = $1) AS int)
	WHERE upkeep_id = $2 AND
	registry_id = (SELECT id FROM keeper_registries WHERE job_id = $1)`,
		jobID, upkeepID, fromAddress.Hex())
	return errors.Wrap(err, "UpdateUpkeepLastKeeperIndex failed")
}

// BatchDeleteUpkeepsForJob deletes all upkeeps by the given IDs for the job with the given ID
func (o *ORM) BatchDeleteUpkeepsForJob(ctx context.Context, jobID int32, upkeepIDs []big.Big) (int64, error) {
	strIds := []string{}
	for _, upkeepID := range upkeepIDs {
		strIds = append(strIds, upkeepID.String())
	}
	res, err := o.ds.ExecContext(ctx, `
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

// EligibleUpkeepsForRegistry fetches eligible upkeeps for processing
// The query checks the following conditions
// - checks the registry address is correct and the registry has some keepers associated
// -- is it my turn AND my keeper was not the last perform for this upkeep OR my keeper was the last before BUT it is past the grace period
// -- OR is it my buddy's turn AND they were the last keeper to do the perform for this upkeep
// DEV: note we cast upkeep_id and binaryHash as 32 bits, even though both are 256 bit numbers when performing XOR. This is enough information
// to distribute the upkeeps over the keepers so long as num keepers < 4294967296
func (o *ORM) EligibleUpkeepsForRegistry(ctx context.Context, registryAddress types.EIP55Address, blockNumber int64, gracePeriod int64, binaryHash string) (upkeeps []UpkeepRegistration, err error) {
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
						-- last keeper != me
						upkeep_registrations.last_keeper_index IS DISTINCT FROM keeper_registries.keeper_index
						OR
						-- last keeper == me AND its past the grace period
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
				-- last keeper == my buddy
				upkeep_registrations.last_keeper_index IS NOT DISTINCT FROM
					(keeper_registries.keeper_index + 1) % keeper_registries.num_keepers
				-- buddy system only active if we have multiple keeper nodes
				AND keeper_registries.num_keepers > 1
			)
		)
`
	if err = o.ds.SelectContext(ctx, &upkeeps, stmt, registryAddress, gracePeriod, blockNumber, binaryHash); err != nil {
		return upkeeps, errors.Wrap(err, "EligibleUpkeepsForRegistry failed to get upkeep_registrations")
	}
	if err = o.loadUpkeepsRegistry(ctx, upkeeps); err != nil {
		return upkeeps, errors.Wrap(err, "EligibleUpkeepsForRegistry failed to load Registry on upkeeps")
	}

	rand.Shuffle(len(upkeeps), func(i, j int) {
		upkeeps[i], upkeeps[j] = upkeeps[j], upkeeps[i]
	})

	return upkeeps, err
}

func (o *ORM) loadUpkeepsRegistry(ctx context.Context, upkeeps []UpkeepRegistration) error {
	registryIDM := make(map[int64]*Registry)
	var registryIDs []int64
	for _, upkeep := range upkeeps {
		if _, exists := registryIDM[upkeep.RegistryID]; !exists {
			registryIDM[upkeep.RegistryID] = new(Registry)
			registryIDs = append(registryIDs, upkeep.RegistryID)
		}
	}
	var registries []*Registry
	err := o.ds.SelectContext(ctx, &registries, `SELECT * FROM keeper_registries WHERE id = ANY($1)`, pq.Array(registryIDs))
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

func (o *ORM) AllUpkeepIDsForRegistry(ctx context.Context, regID int64) (upkeeps []big.Big, err error) {
	err = o.ds.SelectContext(ctx, &upkeeps, `
SELECT upkeep_id
FROM upkeep_registrations
WHERE registry_id = $1
`, regID)
	return upkeeps, errors.Wrap(err, "allUpkeepIDs failed")
}

// SetLastRunInfoForUpkeepOnJob sets the last run block height and the associated keeper index only if the new block height is greater than the previous.
func (o *ORM) SetLastRunInfoForUpkeepOnJob(ctx context.Context, jobID int32, upkeepID *big.Big, height int64, fromAddress types.EIP55Address) (int64, error) {
	res, err := o.ds.ExecContext(ctx, `
	UPDATE upkeep_registrations
	SET last_run_block_height = $1,
		last_keeper_index = CAST((SELECT keeper_index_map -> $4 FROM keeper_registries WHERE job_id = $3) AS int)
	WHERE upkeep_id = $2 AND
	registry_id = (SELECT id FROM keeper_registries WHERE job_id = $3) AND
	last_run_block_height <= $1`, height, upkeepID, jobID, fromAddress.Hex())

	if err != nil {
		return 0, errors.Wrap(err, "SetLastRunInfoForUpkeepOnJob failed")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "SetLastRunInfoForUpkeepOnJob failed to get RowsAffected")
	}
	return rowsAffected, nil
}
