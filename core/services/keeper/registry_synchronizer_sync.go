package keeper

import (
	"encoding/binary"
	"math/big"
	"sync"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// syncUpkeepQueueSize represents the max number of upkeeps that can be synced in parallel
const syncUpkeepQueueSize = 10

func (rs *RegistrySynchronizer) fullSync() {
	contractAddress := rs.job.KeeperSpec.ContractAddress
	rs.logger.Debugf("fullSyncing registry %s", contractAddress.Hex())

	registry, err := rs.syncRegistry()
	if err != nil {
		rs.logger.With("error", err).Error("failed to sync registry during fullSyncing registry")
		return
	}
	if err := rs.addNewUpkeeps(registry); err != nil {
		rs.logger.With("error", err).Error("failed to add new upkeeps during fullSyncing registry")
		return
	}
	if err := rs.deleteCanceledUpkeeps(); err != nil {
		rs.logger.With("error", err).Error("failed to delete canceled upkeeps during fullSyncing registry")
		return
	}
}

func (rs *RegistrySynchronizer) syncRegistry() (Registry, error) {
	registry, err := rs.newRegistryFromChain()
	if err != nil {
		return Registry{}, errors.Wrap(err, "failed to get new registry from chain")
	}

	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	if err := rs.orm.UpsertRegistry(ctx, &registry); err != nil {
		return Registry{}, errors.Wrap(err, "failed to upsert registry")
	}

	return registry, nil
}

func (rs *RegistrySynchronizer) addNewUpkeeps(reg Registry) error {
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	nextUpkeepID, err := rs.orm.LowestUnsyncedID(ctx, reg.ID)
	if err != nil {
		return errors.Wrap(err, "unable to find next ID for registry")
	}

	countOnContractBig, err := rs.contract.GetUpkeepCount(nil)
	if err != nil {
		return errors.Wrapf(err, "unable to get upkeep count")
	}
	countOnContract := countOnContractBig.Int64()

	if nextUpkeepID > countOnContract {
		return errors.New("invariant, contract should always have at least as many upkeeps as DB")
	}

	rs.batchSyncUpkeepsOnRegistry(reg, nextUpkeepID, countOnContract)
	return nil
}

// batchSyncUpkeepsOnRegistry syncs <syncUpkeepQueueSize> upkeeps at a time in parallel
// starting at upkeep ID <start> and up to (but not including) <end>
func (rs *RegistrySynchronizer) batchSyncUpkeepsOnRegistry(reg Registry, start, end int64) {
	wg := sync.WaitGroup{}
	wg.Add(int(end - start))
	chSyncUpkeepQueue := make(chan struct{}, syncUpkeepQueueSize)

	done := func() { <-chSyncUpkeepQueue; wg.Done() }
	for upkeepID := start; upkeepID < end; upkeepID++ {
		select {
		case <-rs.chStop:
			return
		case chSyncUpkeepQueue <- struct{}{}:
			go rs.syncUpkeepWithCallback(reg, upkeepID, done)
		}
	}

	wg.Wait()
}

func (rs *RegistrySynchronizer) syncUpkeepWithCallback(registry Registry, upkeepID int64, doneCallback func()) {
	defer doneCallback()

	if err := rs.syncUpkeep(registry, upkeepID); err != nil {
		rs.logger.With("error", err).With(
			"upkeepID", upkeepID,
			"registryContract", registry.ContractAddress.Hex(),
		).Error("unable to sync upkeep on registry")
	}
}

func (rs *RegistrySynchronizer) syncUpkeep(registry Registry, upkeepID int64) error {
	upkeepConfig, err := rs.contract.GetUpkeep(nil, big.NewInt(upkeepID))
	if err != nil {
		return errors.Wrap(err, "failed to get upkeep config")
	}
	positioningConstant, err := CalcPositioningConstant(upkeepID, registry.ContractAddress)
	if err != nil {
		return errors.Wrap(err, "failed to calc positioning constant")
	}
	newUpkeep := UpkeepRegistration{
		CheckData:           upkeepConfig.CheckData,
		ExecuteGas:          uint64(upkeepConfig.ExecuteGas),
		RegistryID:          registry.ID,
		PositioningConstant: positioningConstant,
		UpkeepID:            upkeepID,
	}
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	if err := rs.orm.UpsertUpkeep(ctx, &newUpkeep); err != nil {
		return errors.Wrap(err, "failed to upsert upkeep")
	}

	return nil
}

func (rs *RegistrySynchronizer) deleteCanceledUpkeeps() error {
	canceledBigs, err := rs.contract.GetCanceledUpkeepList(nil)
	if err != nil {
		return errors.Wrap(err, "failed to get canceled upkeep list")
	}
	canceled := make([]int64, len(canceledBigs))
	for idx, upkeepID := range canceledBigs {
		canceled[idx] = upkeepID.Int64()
	}
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	if _, err := rs.orm.BatchDeleteUpkeepsForJob(ctx, rs.job.ID, canceled); err != nil {
		return errors.Wrap(err, "failed to batch delete upkeeps from job")
	}

	return nil
}

// newRegistryFromChain returns a Registry stuct with fields synched from those on chain
func (rs *RegistrySynchronizer) newRegistryFromChain() (Registry, error) {
	fromAddress := rs.job.KeeperSpec.FromAddress
	contractAddress := rs.job.KeeperSpec.ContractAddress
	config, err := rs.contract.GetConfig(nil)
	if err != nil {
		ctx, cancel := postgres.DefaultQueryCtx()
		defer cancel()
		rs.jrm.RecordError(ctx, rs.job.ID, err.Error())
		return Registry{}, errors.Wrap(err, "failed to get contract config")
	}
	keeperAddresses, err := rs.contract.GetKeeperList(nil)
	if err != nil {
		return Registry{}, errors.Wrap(err, "failed to get keeper list")
	}
	keeperIndex := int32(-1)
	for idx, address := range keeperAddresses {
		if address == fromAddress.Address() {
			keeperIndex = int32(idx)
		}
	}
	if keeperIndex == -1 {
		rs.logger.Warnf("unable to find %s in keeper list on registry %s", fromAddress.Hex(), contractAddress.Hex())
	}

	return Registry{
		BlockCountPerTurn: int32(config.BlockCountPerTurn.Int64()),
		CheckGas:          int32(config.CheckGasLimit),
		ContractAddress:   contractAddress,
		FromAddress:       fromAddress,
		JobID:             rs.job.ID,
		KeeperIndex:       keeperIndex,
		NumKeepers:        int32(len(keeperAddresses)),
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
