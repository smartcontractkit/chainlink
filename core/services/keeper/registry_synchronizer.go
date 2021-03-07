package keeper

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_contract"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"go.uber.org/atomic"
)

const syncUpkeepQueueSize = 10

func NewRegistrySynchronizer(
	job job.Job,
	contract *keeper_registry_contract.KeeperRegistryContract,
	keeperORM KeeperORM,
	syncInterval time.Duration,
) RegistrySynchronizer {
	return RegistrySynchronizer{
		contract:  contract,
		keeperORM: keeperORM,
		job:       job,
		interval:  syncInterval,
		isRunning: atomic.NewBool(false),
		chDone:    make(chan struct{}),
	}
}

// RegistrySynchronizer conforms to the job.Service interface
var _ job.Service = RegistrySynchronizer{}

type RegistrySynchronizer struct {
	contract  *keeper_registry_contract.KeeperRegistryContract
	interval  time.Duration
	isRunning *atomic.Bool
	job       job.Job
	keeperORM KeeperORM

	chDone chan struct{}
}

func (rs RegistrySynchronizer) Start() error {
	if rs.isRunning.Load() {
		return errors.New("already started")
	}
	rs.isRunning.Store(true)
	go rs.run()
	return nil
}

func (rs RegistrySynchronizer) Close() error {
	close(rs.chDone)
	return nil
}

func (rs RegistrySynchronizer) run() {
	ticker := time.NewTicker(rs.interval)
	defer ticker.Stop()

	for {
		select {
		case <-rs.chDone:
			return
		case <-ticker.C:
			rs.syncRegistry()
		}
	}
}

func (rs RegistrySynchronizer) syncRegistry() {
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

func (rs RegistrySynchronizer) addNewUpkeeps(reg Registry) error {
	nextUpkeepID, err := rs.keeperORM.NextUpkeepIDForRegistry(reg)
	if err != nil {
		return err
	}

	countOnContractBig, err := rs.contract.GetUpkeepCount(nil)
	if err != nil {
		return err
	}
	countOnContract := countOnContractBig.Int64()

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

func (rs RegistrySynchronizer) deleteCanceledUpkeeps(reg Registry) error {
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

func (rs RegistrySynchronizer) syncUpkeep(registry Registry, upkeepID int64, doneCallback func()) {
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
		logger.Errorf("unable to sync upkeep #%d on registry %s, err: %v", upkeepID, registry.ContractAddress.Hex(), err)
	}
}

func (rs RegistrySynchronizer) newSyncedRegistry(job job.Job) (Registry, error) {
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
	found := false
	var keeperIndex int32
	for idx, address := range keeperAddresses {
		if address == fromAddress.Address() {
			keeperIndex = int32(idx)
			found = true
		}
	}
	if !found {
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
