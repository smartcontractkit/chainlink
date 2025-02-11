package keeper

import (
	"context"
	"encoding/binary"
	"math"
	"sync"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-integrations/evm/types"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
)

func (rs *RegistrySynchronizer) fullSync(ctx context.Context) {
	rs.logger.Debugf("fullSyncing registry %s", rs.job.KeeperSpec.ContractAddress.Hex())

	registry, err := rs.syncRegistry(ctx)
	if err != nil {
		rs.logger.Error(errors.Wrap(err, "failed to sync registry during fullSyncing registry"))
		return
	}

	if err := rs.fullSyncUpkeeps(ctx, registry); err != nil {
		rs.logger.Error(errors.Wrap(err, "failed to sync upkeeps during fullSyncing registry"))
		return
	}
	rs.logger.Debugf("fullSyncing registry successful %s", rs.job.KeeperSpec.ContractAddress.Hex())
}

func (rs *RegistrySynchronizer) syncRegistry(ctx context.Context) (Registry, error) {
	registry, err := rs.newRegistryFromChain(ctx)
	if err != nil {
		return Registry{}, errors.Wrap(err, "failed to get new registry from chain")
	}

	err = rs.orm.UpsertRegistry(ctx, &registry)
	if err != nil {
		return Registry{}, errors.Wrap(err, "failed to upsert registry")
	}

	return registry, nil
}

func (rs *RegistrySynchronizer) fullSyncUpkeeps(ctx context.Context, reg Registry) error {
	activeUpkeepIDs, err := rs.registryWrapper.GetActiveUpkeepIDs(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "unable to get active upkeep IDs")
	}

	existingUpkeepIDs, err := rs.orm.AllUpkeepIDsForRegistry(ctx, reg.ID)
	if err != nil {
		return errors.Wrap(err, "unable to fetch existing upkeep IDs from DB")
	}

	activeSet := make(map[string]bool)
	allActiveUpkeeps := make([]big.Big, 0)
	for _, upkeepID := range activeUpkeepIDs {
		activeSet[upkeepID.String()] = true
		allActiveUpkeeps = append(allActiveUpkeeps, *big.New(upkeepID))
	}
	rs.batchSyncUpkeepsOnRegistry(ctx, reg, allActiveUpkeeps)

	// All upkeeps in existingUpkeepIDs, not in activeUpkeepIDs should be deleted
	canceled := make([]big.Big, 0)
	for _, upkeepID := range existingUpkeepIDs {
		if _, found := activeSet[upkeepID.ToInt().String()]; !found {
			canceled = append(canceled, upkeepID)
		}
	}
	if _, err := rs.orm.BatchDeleteUpkeepsForJob(ctx, rs.job.ID, canceled); err != nil {
		return errors.Wrap(err, "failed to batch delete upkeeps from job")
	}
	return nil
}

// batchSyncUpkeepsOnRegistry syncs <syncUpkeepQueueSize> upkeeps at a time in parallel
// for all the IDs within newUpkeeps slice
func (rs *RegistrySynchronizer) batchSyncUpkeepsOnRegistry(ctx context.Context, reg Registry, newUpkeeps []big.Big) {
	wg := sync.WaitGroup{}
	chSyncUpkeepQueue := make(chan struct{}, rs.syncUpkeepQueueSize)

	done := func() { <-chSyncUpkeepQueue; wg.Done() }
	for i := range newUpkeeps {
		select {
		case <-ctx.Done():
		case chSyncUpkeepQueue <- struct{}{}:
			wg.Add(1)
			go rs.syncUpkeepWithCallback(ctx, &rs.registryWrapper, reg, &newUpkeeps[i], done)
		}
	}

	wg.Wait()
}

func (rs *RegistrySynchronizer) syncUpkeepWithCallback(ctx context.Context, getter upkeepGetter, registry Registry, upkeepID *big.Big, doneCallback func()) {
	defer doneCallback()

	if err := rs.syncUpkeep(ctx, getter, registry, upkeepID); err != nil {
		rs.logger.With("err", err.Error()).With(
			"upkeepID", NewUpkeepIdentifier(upkeepID).String(),
			"registryContract", registry.ContractAddress.Hex(),
		).Error("unable to sync upkeep on registry")
	}
}

func (rs *RegistrySynchronizer) syncUpkeep(ctx context.Context, getter upkeepGetter, registry Registry, upkeepID *big.Big) error {
	upkeep, err := getter.GetUpkeep(nil, upkeepID.ToInt())
	if err != nil {
		return errors.Wrap(err, "failed to get upkeep config")
	}

	if upkeep.ExecuteGas <= uint32(0) {
		return errors.Errorf("execute gas is zero for upkeep %s", NewUpkeepIdentifier(upkeepID).String())
	}

	positioningConstant, err := CalcPositioningConstant(upkeepID, registry.ContractAddress)
	if err != nil {
		return errors.Wrap(err, "failed to calc positioning constant")
	}
	newUpkeep := UpkeepRegistration{
		CheckData:           upkeep.CheckData,
		ExecuteGas:          upkeep.ExecuteGas,
		RegistryID:          registry.ID,
		PositioningConstant: positioningConstant,
		UpkeepID:            upkeepID,
	}
	if err := rs.orm.UpsertUpkeep(ctx, &newUpkeep); err != nil {
		return errors.Wrap(err, "failed to upsert upkeep")
	}

	if err := rs.orm.UpdateUpkeepLastKeeperIndex(ctx, rs.job.ID, upkeepID, types.EIP55AddressFromAddress(upkeep.LastKeeper)); err != nil {
		return errors.Wrap(err, "failed to update upkeep last keeper index")
	}

	return nil
}

// newRegistryFromChain returns a Registry struct with fields synched from those on chain
func (rs *RegistrySynchronizer) newRegistryFromChain(ctx context.Context) (Registry, error) {
	fromAddress := rs.effectiveKeeperAddress
	contractAddress := rs.job.KeeperSpec.ContractAddress

	registryConfig, err := rs.registryWrapper.GetConfig(nil)
	if err != nil {
		rs.jrm.TryRecordError(ctx, rs.job.ID, err.Error())
		return Registry{}, errors.Wrap(err, "failed to get contract config")
	}

	keeperIndex := int32(-1)
	keeperMap := map[types.EIP55Address]int32{}
	for idx, address := range registryConfig.KeeperAddresses {
		keeperMap[types.EIP55AddressFromAddress(address)] = int32(idx)
		if address == fromAddress {
			keeperIndex = int32(idx)
		}
	}
	if keeperIndex == -1 {
		rs.logger.Warnf("unable to find %s in keeper list on registry %s", fromAddress.Hex(), contractAddress.Hex())
	}

	return Registry{
		BlockCountPerTurn: registryConfig.BlockCountPerTurn,
		CheckGas:          registryConfig.CheckGas,
		ContractAddress:   contractAddress,
		FromAddress:       rs.job.KeeperSpec.FromAddress,
		JobID:             rs.job.ID,
		KeeperIndex:       keeperIndex,
		NumKeepers:        int32(len(registryConfig.KeeperAddresses)),
		KeeperIndexMap:    keeperMap,
	}, nil
}

// CalcPositioningConstant calculates a positioning constant.
// The positioning constant is fixed because upkeepID and registryAddress are immutable
func CalcPositioningConstant(upkeepID *big.Big, registryAddress types.EIP55Address) (int32, error) {
	upkeepBytes := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(upkeepBytes, upkeepID.Mod(big.NewI(math.MaxInt64)).Int64())
	bytesToHash := utils.ConcatBytes(upkeepBytes, registryAddress.Bytes())
	checksum, err := utils.Keccak256(bytesToHash)
	if err != nil {
		return 0, err
	}
	constant := binary.BigEndian.Uint16(checksum[:2])
	return int32(constant), nil
}
