package keeper

import (
	"encoding/binary"
	"math"
	"sync"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func (rs *RegistrySynchronizer) fullSync() {
	rs.logger.Debugf("fullSyncing registry %s", rs.job.KeeperSpec.ContractAddress.Hex())

	registry, err := rs.syncRegistry()
	if err != nil {
		rs.logger.Error(errors.Wrap(err, "failed to sync registry during fullSyncing registry"))
		return
	}

	if err := rs.fullSyncUpkeeps(registry); err != nil {
		rs.logger.Error(errors.Wrap(err, "failed to sync upkeeps during fullSyncing registry"))
		return
	}
	rs.logger.Debugf("fullSyncing registry successful %s", rs.job.KeeperSpec.ContractAddress.Hex())
}

func (rs *RegistrySynchronizer) syncRegistry() (Registry, error) {
	registry, err := rs.newRegistryFromChain()
	if err != nil {
		return Registry{}, errors.Wrap(err, "failed to get new registry from chain")
	}

	err = rs.orm.UpsertRegistry(&registry)
	if err != nil {
		return Registry{}, errors.Wrap(err, "failed to upsert registry")
	}

	return registry, nil
}

func (rs *RegistrySynchronizer) fullSyncUpkeeps(reg Registry) error {
	activeUpkeepIDs, err := rs.registryWrapper.GetActiveUpkeepIDs(nil)
	if err != nil {
		return errors.Wrap(err, "unable to get active upkeep IDs")
	}

	existingUpkeepIDs, err := rs.orm.AllUpkeepIDsForRegistry(reg.ID)
	if err != nil {
		return errors.Wrap(err, "unable to fetch existing upkeep IDs from DB")
	}

	activeSet := make(map[string]bool)
	allActiveUpkeeps := make([]utils.Big, 0)
	for _, upkeepID := range activeUpkeepIDs {
		activeSet[upkeepID.String()] = true
		allActiveUpkeeps = append(allActiveUpkeeps, *utils.NewBig(upkeepID))
	}
	rs.batchSyncUpkeepsOnRegistry(reg, allActiveUpkeeps)

	// All upkeeps in existingUpkeepIDs, not in activeUpkeepIDs should be deleted
	canceled := make([]utils.Big, 0)
	for _, upkeepID := range existingUpkeepIDs {
		if _, found := activeSet[upkeepID.ToInt().String()]; !found {
			canceled = append(canceled, upkeepID)
		}
	}
	if _, err := rs.orm.BatchDeleteUpkeepsForJob(rs.job.ID, canceled); err != nil {
		return errors.Wrap(err, "failed to batch delete upkeeps from job")
	}
	return nil
}

// batchSyncUpkeepsOnRegistry syncs <syncUpkeepQueueSize> upkeeps at a time in parallel
// for all the IDs within newUpkeeps slice
func (rs *RegistrySynchronizer) batchSyncUpkeepsOnRegistry(reg Registry, newUpkeeps []utils.Big) {
	wg := sync.WaitGroup{}
	wg.Add(len(newUpkeeps))
	chSyncUpkeepQueue := make(chan struct{}, rs.syncUpkeepQueueSize)

	done := func() { <-chSyncUpkeepQueue; wg.Done() }
	for i := range newUpkeeps {
		select {
		case <-rs.chStop:
			return
		case chSyncUpkeepQueue <- struct{}{}:
			go rs.syncUpkeepWithCallback(&rs.registryWrapper, reg, &newUpkeeps[i], done)
		}
	}

	wg.Wait()
}

func (rs *RegistrySynchronizer) syncUpkeepWithCallback(getter upkeepGetter, registry Registry, upkeepID *utils.Big, doneCallback func()) {
	defer doneCallback()

	if err := rs.syncUpkeep(getter, registry, upkeepID); err != nil {
		rs.logger.With("error", err.Error()).With(
			"upkeepID", NewUpkeepIdentifier(upkeepID).String(),
			"registryContract", registry.ContractAddress.Hex(),
		).Error("unable to sync upkeep on registry")
	}
}

func (rs *RegistrySynchronizer) syncUpkeep(getter upkeepGetter, registry Registry, upkeepID *utils.Big) error {
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
	if err := rs.orm.UpsertUpkeep(&newUpkeep); err != nil {
		return errors.Wrap(err, "failed to upsert upkeep")
	}

	if err := rs.orm.UpdateUpkeepLastKeeperIndex(rs.job.ID, upkeepID, ethkey.EIP55AddressFromAddress(upkeep.LastKeeper)); err != nil {
		return errors.Wrap(err, "failed to update upkeep last keeper index")
	}

	return nil
}

// newRegistryFromChain returns a Registry struct with fields synched from those on chain
func (rs *RegistrySynchronizer) newRegistryFromChain() (Registry, error) {
	fromAddress := rs.effectiveKeeperAddress
	contractAddress := rs.job.KeeperSpec.ContractAddress

	registryConfig, err := rs.registryWrapper.GetConfig(nil)
	if err != nil {
		rs.jrm.TryRecordError(rs.job.ID, err.Error())
		return Registry{}, errors.Wrap(err, "failed to get contract config")
	}

	keeperIndex := int32(-1)
	keeperMap := map[ethkey.EIP55Address]int32{}
	for idx, address := range registryConfig.KeeperAddresses {
		keeperMap[ethkey.EIP55AddressFromAddress(address)] = int32(idx)
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
func CalcPositioningConstant(upkeepID *utils.Big, registryAddress ethkey.EIP55Address) (int32, error) {
	upkeepBytes := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(upkeepBytes, upkeepID.Mod(math.MaxInt64).Int64())
	bytesToHash := utils.ConcatBytes(upkeepBytes, registryAddress.Bytes())
	checksum, err := utils.Keccak256(bytesToHash)
	if err != nil {
		return 0, err
	}
	constant := binary.BigEndian.Uint16(checksum[:2])
	return int32(constant), nil
}
