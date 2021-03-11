package keeper

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const syncUpkeepQueueSize = 10

func NewRegistrySynchronizer(
	job job.Job,
	contract *keeper_registry_wrapper.KeeperRegistry,
	keeperORM KeeperORM,
	syncInterval time.Duration,
) *RegistrySynchronizer {
	return &RegistrySynchronizer{
		contract:      contract,
		doneWG:        sync.WaitGroup{},
		keeperORM:     keeperORM,
		job:           job,
		interval:      syncInterval,
		chDone:        make(chan struct{}),
		StartStopOnce: utils.StartStopOnce{},
	}
}

// RegistrySynchronizer conforms to the job.Service interface
var _ job.Service = &RegistrySynchronizer{}

type RegistrySynchronizer struct {
	contract  *keeper_registry_wrapper.KeeperRegistry
	doneWG    sync.WaitGroup
	interval  time.Duration
	job       job.Job
	keeperORM KeeperORM

	chDone chan struct{}

	utils.StartStopOnce
}

func (rs *RegistrySynchronizer) Start() error {
	if !rs.OkayToStart() {
		return errors.New("RegistrySynchronizer is already started")
	}
	go rs.run()
	return nil
}

func (rs *RegistrySynchronizer) Close() error {
	if !rs.OkayToStop() {
		return errors.New("RegistrySynchronizer is already stopped")
	}
	close(rs.chDone)
	rs.doneWG.Wait()
	return nil
}

func (rs *RegistrySynchronizer) run() {
	rs.doneWG.Add(1)
	ticker := time.NewTicker(rs.interval)
	defer ticker.Stop()

	for {
		select {
		case <-rs.chDone:
			rs.doneWG.Done()
			return
		case <-ticker.C:
			rs.syncRegistry()
		}
	}
}

func (rs *RegistrySynchronizer) syncRegistry() {
	contractAddress := rs.job.KeeperSpec.ContractAddress
	logger.Debugf("syncing registry %s", contractAddress.Hex())

	err := func() error {
		registry, err := rs.newSyncedRegistry(rs.job)
		if err != nil {
			return err
		}
		if err = rs.keeperORM.UpsertRegistry(&registry); err != nil {
			return err
		}
		if err = rs.addNewUpkeeps(registry); err != nil {
			return err
		}
		if err = rs.deleteCanceledUpkeeps(registry); err != nil {
			return err
		}
		return nil
	}()

	if err != nil {
		logger.Errorf("unable to sync registry %s, err: %v", contractAddress.Hex(), err)
	}
}

func (rs *RegistrySynchronizer) addNewUpkeeps(reg Registry) error {
	nextUpkeepID, err := rs.keeperORM.NextUpkeepIDForRegistry(reg)
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

	wg := sync.WaitGroup{}
	wg.Add(int(countOnContract - nextUpkeepID))

	// batch sync registries
	chSyncUpkeepQueue := make(chan struct{}, syncUpkeepQueueSize)
	done := func() { <-chSyncUpkeepQueue; wg.Done() }
	for upkeepID := nextUpkeepID; upkeepID < countOnContract; upkeepID++ {
		chSyncUpkeepQueue <- struct{}{}
		go rs.syncUpkeep(reg, upkeepID, done)
	}

	wg.Wait()
	return nil
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
	return rs.keeperORM.BatchDeleteUpkeeps(reg.ID, canceled)
}

func (rs *RegistrySynchronizer) syncUpkeep(registry Registry, upkeepID int64, doneCallback func()) {
	defer doneCallback()

	err := func() error {
		upkeepConfig, err := rs.contract.GetUpkeep(nil, big.NewInt(int64(upkeepID)))
		if err != nil {
			return err
		}
		positioningConstant, err := calcPositioningConstant(upkeepID, registry.ContractAddress, registry.NumKeepers)
		if err != nil {
			return fmt.Errorf("unable to calculate positioning constant: %v", err)
		}
		newUpkeep := UpkeepRegistration{
			CheckData:           upkeepConfig.CheckData,
			ExecuteGas:          int32(upkeepConfig.ExecuteGas),
			RegistryID:          registry.ID,
			PositioningConstant: positioningConstant,
			UpkeepID:            upkeepID,
		}

		return rs.keeperORM.UpsertUpkeep(&newUpkeep)
	}()

	if err != nil {
		logger.Errorw(
			fmt.Sprintf("unable to sync upkeep #%d on registry %s", upkeepID, registry.ContractAddress.Hex()),
			"error",
			err,
		)
	}
}

func (rs *RegistrySynchronizer) newSyncedRegistry(job job.Job) (Registry, error) {
	fromAddress := job.KeeperSpec.FromAddress
	contractAddress := job.KeeperSpec.ContractAddress
	config, err := rs.contract.GetConfig(nil)
	if err != nil {
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
		return Registry{}, fmt.Errorf("unable to find %s in keeper list on registry %s", fromAddress.Hex(), contractAddress.Hex())
	}

	return Registry{
		BlockCountPerTurn: int32(config.BlockCountPerTurn.Int64()),
		CheckGas:          int32(config.CheckGasLimit),
		ContractAddress:   contractAddress,
		FromAddress:       fromAddress,
		JobID:             job.ID,
		KeeperIndex:       keeperIndex,
		NumKeepers:        int32(len(keeperAddresses)),
	}, nil
}

func calcPositioningConstant(upkeepID int64, registryAddress models.EIP55Address, numKeepers int32) (int32, error) {
	if numKeepers == 0 {
		return 0, errors.New("cannot calc positioning constant with 0 keepers")
	}

	upkeepBytes := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(upkeepBytes, upkeepID)
	bytesToHash := utils.ConcatBytes(upkeepBytes, registryAddress.Bytes())
	hash, err := utils.Keccak256(bytesToHash)
	if err != nil {
		return 0, err
	}
	hashUint := big.NewInt(0).SetBytes(hash)
	constant := big.NewInt(0).Mod(hashUint, big.NewInt(int64(numKeepers)))

	return int32(constant.Int64()), nil
}
