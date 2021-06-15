package keeper

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"sync"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// syncUpkeepQueueSize represents the max number of upkeeps that can be synced in parallel
const syncUpkeepQueueSize = 10

func (rs *RegistrySynchronizer) fullSync() {
	contractAddress := rs.job.KeeperSpec.ContractAddress
	logger.Debugf("fullSyncing registry %s", contractAddress.Hex())

	var err error
	defer func() {
		logger.ErrorIf(err, fmt.Sprintf("unable to fullSync registry %s", contractAddress.Hex()))
	}()

	registry, err := rs.syncRegistry()
	if err != nil {
		return
	}
	if err = rs.addNewUpkeeps(registry); err != nil {
		return
	}
	if err = rs.deleteCanceledUpkeeps(registry); err != nil {
		return
	}
}

func (rs *RegistrySynchronizer) syncRegistry() (Registry, error) {
	registry, err := rs.newRegistryFromChain()
	if err != nil {
		return Registry{}, err
	}
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	if err = rs.orm.UpsertRegistry(ctx, &registry); err != nil {
		return Registry{}, err
	}
	return registry, err
}

func (rs *RegistrySynchronizer) addNewUpkeeps(reg Registry) error {
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	nextUpkeepID, err := rs.orm.LowestUnsyncedID(ctx, reg)
	if err != nil {
		return errors.Wrap(err, "RegistrySynchronizer: unable to find next ID for registry")
	}

	countOnContractBig, err := rs.contract.GetUpkeepCount(nil)
	if err != nil {
		return errors.Wrapf(err, "RegistrySynchronizer: unable to get upkeep count")
	}
	countOnContract := countOnContractBig.Int64()

	if nextUpkeepID > countOnContract {
		return errors.New("RegistrySynchronizer: invariant, contract should always have at least as many upkeeps as DB")
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
	err := rs.syncUpkeep(registry, upkeepID)
	if err != nil {
		logger.ErrorIf(err, fmt.Sprintf("unable to sync upkeep #%d on registry %s", upkeepID, registry.ContractAddress.Hex()))
	}
}

func (rs *RegistrySynchronizer) syncUpkeep(registry Registry, upkeepID int64) error {
	upkeepConfig, err := rs.contract.GetUpkeep(nil, big.NewInt(int64(upkeepID)))
	if err != nil {
		return err
	}
	positioningConstant, err := CalcPositioningConstant(upkeepID, registry.ContractAddress)
	if err != nil {
		return err
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
	return rs.orm.UpsertUpkeep(ctx, &newUpkeep)
}

func (rs *RegistrySynchronizer) deleteCanceledUpkeeps(reg Registry) error {
	canceledBigs, err := rs.contract.GetCanceledUpkeepList(nil)
	if err != nil {
		return err
	}
	canceled := make([]int64, len(canceledBigs))
	for idx, upkeepID := range canceledBigs {
		canceled[idx] = upkeepID.Int64()
	}
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	_, err = rs.orm.BatchDeleteUpkeepsForJob(ctx, rs.job.ID, canceled)
	return err
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
		return Registry{}, err
	}
	keeperAddresses, err := rs.contract.GetKeeperList(nil)
	if err != nil {
		return Registry{}, err
	}
	keeperIndex := int32(-1)
	for idx, address := range keeperAddresses {
		if address == fromAddress.Address() {
			keeperIndex = int32(idx)
		}
	}
	if keeperIndex == -1 {
		logger.Warnf("unable to find %s in keeper list on registry %s", fromAddress.Hex(), contractAddress.Hex())
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

// the positioning constant is fixed because upkeepID and registryAddress are immutable
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
