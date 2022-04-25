package keeper

import (
	"encoding/binary"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func (rs *RegistrySynchronizer) fullSync() {
	contractAddress := rs.job.KeeperSpec.ContractAddress
	rs.logger.Debugf("fullSyncing registry %s", contractAddress.Hex())

	registry, err := rs.syncRegistry()
	if err != nil {
		rs.logger.With("error", err).Error("failed to sync registry during fullSyncing registry")
		return
	}

	if err := rs.fullSyncUpkeeps(registry); err != nil {
		rs.logger.With("error", err).Error("failed to sync upkeeps during fullSyncing registry")
		return
	}
}

func (rs *RegistrySynchronizer) syncRegistry() (Registry, error) {
	registry, err := rs.newRegistryFromChain()
	if err != nil {
		return Registry{}, errors.Wrap(err, "failed to get new registry from chain")
	}

	if err := rs.orm.UpsertRegistry(&registry); err != nil {
		return Registry{}, errors.Wrap(err, "failed to upsert registry")
	}

	return registry, nil
}

func (rs *RegistrySynchronizer) fullSyncUpkeeps(reg Registry) error {
	switch rs.version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		return rs.fullSyncUpkeeps1_1(reg)
	case RegistryVersion_1_2:
		return rs.fullSyncUpkeeps1_2(reg)
	}
	return nil
}

func (rs *RegistrySynchronizer) fullSyncUpkeeps1_1(reg Registry) error {
	// Add new upkeeps which are not in DB
	nextUpkeepID, err := rs.orm.LowestUnsyncedID(reg.ID)
	if err != nil {
		return errors.Wrap(err, "unable to find next ID for registry")
	}

	countOnContractBig, err := rs.contract1_1.GetUpkeepCount(nil)
	if err != nil {
		return errors.Wrapf(err, "unable to get upkeep count")
	}
	countOnContract := countOnContractBig.Int64()

	if nextUpkeepID > countOnContract {
		return errors.New("invariant, contract should always have at least as many upkeeps as DB")
	}
	// newUpkeeps is the range nextUpkeepID, nextUpkeepID + 1 , ... , countOnContract-1
	newUpkeeps := make([]*big.Int, countOnContract-nextUpkeepID)
	for i := range newUpkeeps {
		newUpkeeps[i] = big.NewInt(nextUpkeepID + int64(i))
	}
	rs.batchSyncUpkeepsOnRegistry(reg, newUpkeeps)

	// Delete upkeeps which have been cancelled
	canceledBigs, err := rs.contract1_1.GetCanceledUpkeepList(nil)
	if err != nil {
		return errors.Wrap(err, "failed to get canceled upkeep list")
	}
	canceled := make([]int64, len(canceledBigs))
	for idx, upkeepID := range canceledBigs {
		canceled[idx] = upkeepID.Int64()
	}
	if _, err := rs.orm.BatchDeleteUpkeepsForJob(rs.job.ID, canceled); err != nil {
		return errors.Wrap(err, "failed to batch delete upkeeps from job")
	}
	return nil
}

func (rs *RegistrySynchronizer) fullSyncUpkeeps1_2(reg Registry) error {
	startIndex := big.NewInt(0)
	maxCount := big.NewInt(0)
	// TODO (sc-37024): Get active upkeep IDs from contract in batches
	activeUpkeepIDs, err := rs.contract1_2.GetActiveUpkeepIDs(nil, startIndex, maxCount)
	if err != nil {
		return errors.Wrapf(err, "unable to get active upkeep IDs")
	}
	existingUpkeepIDs, err := rs.orm.AllUpkeepIDsForRegistry(reg.ID)
	if err != nil {
		return errors.Wrap(err, "unable to fetch existing upkeep IDs from DB")
	}

	existingSet := make(map[string]bool)
	activeSet := make(map[string]bool)
	// New upkeeps are all elements in activeUpkeepIDs which are not in existingUpkeepIDs
	newUpkeeps := make([]*big.Int, 0)
	for _, upkeepID := range existingUpkeepIDs {
		existingSet[upkeepID.ToInt().String()] = true
	}
	for _, upkeepID := range activeUpkeepIDs {
		activeSet[upkeepID.String()] = true
		if _, found := existingSet[upkeepID.String()]; !found {
			newUpkeeps = append(newUpkeeps, upkeepID)
		}
	}
	rs.batchSyncUpkeepsOnRegistry(reg, newUpkeeps)

	// All upkeeps in existingUpkeepIDs, not in activeUpkeepIDs should be deleted
	canceled := make([]int64, 0)
	for _, upkeepID := range existingUpkeepIDs {
		if _, found := activeSet[upkeepID.ToInt().String()]; !found {
			canceled = append(canceled, upkeepID.ToInt().Int64())
		}
	}
	if _, err := rs.orm.BatchDeleteUpkeepsForJob(rs.job.ID, canceled); err != nil {
		return errors.Wrap(err, "failed to batch delete upkeeps from job")
	}
	return nil
}

// batchSyncUpkeepsOnRegistry syncs <syncUpkeepQueueSize> upkeeps at a time in parallel
// for all the IDs within newUpkeeps slice
func (rs *RegistrySynchronizer) batchSyncUpkeepsOnRegistry(reg Registry, newUpkeeps []*big.Int) {
	wg := sync.WaitGroup{}
	wg.Add(len(newUpkeeps))
	chSyncUpkeepQueue := make(chan struct{}, rs.syncUpkeepQueueSize)

	done := func() { <-chSyncUpkeepQueue; wg.Done() }
	for _, upkeepID := range newUpkeeps {
		select {
		case <-rs.chStop:
			return
		case chSyncUpkeepQueue <- struct{}{}:
			go rs.syncUpkeepWithCallback(reg, upkeepID, done)
		}
	}

	wg.Wait()
}

func (rs *RegistrySynchronizer) syncUpkeepWithCallback(registry Registry, upkeepID *big.Int, doneCallback func()) {
	defer doneCallback()

	if err := rs.syncUpkeep(registry, upkeepID); err != nil {
		rs.logger.With("error", err).With(
			"upkeepID", upkeepID,
			"registryContract", registry.ContractAddress.Hex(),
		).Error("unable to sync upkeep on registry")
	}
}

func (rs *RegistrySynchronizer) syncUpkeep(registry Registry, upkeepID *big.Int) error {
	var checkData []byte
	var executeGas uint64
	switch rs.version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		upkeepConfig, err := rs.contract1_1.GetUpkeep(nil, upkeepID)
		if err != nil {
			return errors.Wrap(err, "failed to get upkeep config")
		}
		checkData = upkeepConfig.CheckData
		executeGas = uint64(upkeepConfig.ExecuteGas)
	case RegistryVersion_1_2:
		upkeepConfig, err := rs.contract1_2.GetUpkeep(nil, upkeepID)
		if err != nil {
			return errors.Wrap(err, "failed to get upkeep config")
		}
		checkData = upkeepConfig.CheckData
		executeGas = uint64(upkeepConfig.ExecuteGas)
	}

	newUpkeep := UpkeepRegistration{
		CheckData:  checkData,
		ExecuteGas: executeGas,
		RegistryID: registry.ID,
		UpkeepID:   upkeepID.Int64(),
	}
	if err := rs.orm.UpsertUpkeep(&newUpkeep); err != nil {
		return errors.Wrap(err, "failed to upsert upkeep")
	}

	return nil
}

// newRegistryFromChain returns a Registry struct with fields synched from those on chain
func (rs *RegistrySynchronizer) newRegistryFromChain() (Registry, error) {
	fromAddress := rs.job.KeeperSpec.FromAddress
	contractAddress := rs.job.KeeperSpec.ContractAddress

	var blockCountPerTurn int32
	var checkGas int32
	var keeperAddresses []common.Address

	switch rs.version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		config, err := rs.contract1_1.GetConfig(nil)
		if err != nil {
			rs.jrm.TryRecordError(rs.job.ID, err.Error())
			return Registry{}, errors.Wrap(err, "failed to get contract config")
		}
		keeperAddresses, err = rs.contract1_1.GetKeeperList(nil)
		if err != nil {
			return Registry{}, errors.Wrap(err, "failed to get keeper list")
		}
		blockCountPerTurn = int32(config.BlockCountPerTurn.Int64())
		checkGas = int32(config.CheckGasLimit)
	case RegistryVersion_1_2:
		state, err := rs.contract1_2.GetState(nil)
		if err != nil {
			rs.jrm.TryRecordError(rs.job.ID, err.Error())
			return Registry{}, errors.Wrap(err, "failed to get contract state")
		}
		keeperAddresses = state.Keepers
		blockCountPerTurn = int32(state.Config.BlockCountPerTurn.Int64())
		checkGas = int32(state.Config.CheckGasLimit)
	}

	keeperIndex := int32(-1)
	keeperMap := map[ethkey.EIP55Address]int32{}
	for idx, address := range keeperAddresses {
		keeperMap[ethkey.EIP55AddressFromAddress(address)] = int32(idx)
		if address == fromAddress.Address() {
			keeperIndex = int32(idx)
		}
	}
	if keeperIndex == -1 {
		rs.logger.Warnf("unable to find %s in keeper list on registry %s", fromAddress.Hex(), contractAddress.Hex())
	}

	return Registry{
		BlockCountPerTurn: blockCountPerTurn,
		CheckGas:          checkGas,
		ContractAddress:   contractAddress,
		FromAddress:       fromAddress,
		JobID:             rs.job.ID,
		KeeperIndex:       keeperIndex,
		NumKeepers:        int32(len(keeperAddresses)),
		KeeperIndexMap:    keeperMap,
	}, nil
}

// CalcPositioningConstant calculates a positioning constant.
// The positioning constant is fixed because upkeepID and registryAddress are immutable
func CalcPositioningConstant(upkeepID int64, registryAddress ethkey.EIP55Address) (int32, error) {
	upkeepBytes := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(upkeepBytes, upkeepID)
	bytesToHash := utils.ConcatBytes(upkeepBytes, registryAddress.Bytes())
	checksum, err := utils.Keccak256(bytesToHash)
	if err != nil {
		return 0, err
	}
	constant := binary.BigEndian.Uint16(checksum[:2])
	return int32(constant), nil
}
