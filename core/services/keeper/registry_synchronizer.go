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
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

const syncUpkeepQueueSize = 10

func NewRegistrySynchronizer(
	job job.Job,
	contract *keeper_registry_wrapper.KeeperRegistry,
	db *gorm.DB,
	syncInterval time.Duration,
) *RegistrySynchronizer {
	return &RegistrySynchronizer{
		contract:      contract,
		interval:      syncInterval,
		job:           job,
		orm:           NewORM(db),
		StartStopOnce: utils.StartStopOnce{},
		wgDone:        sync.WaitGroup{},
		chStop:        make(chan struct{}),
	}
}

// RegistrySynchronizer conforms to the job.Service interface
var _ job.Service = (*RegistrySynchronizer)(nil)

type RegistrySynchronizer struct {
	contract *keeper_registry_wrapper.KeeperRegistry
	interval time.Duration
	job      job.Job
	orm      ORM
	wgDone   sync.WaitGroup
	chStop   chan struct{}
	utils.StartStopOnce
}

func (rs *RegistrySynchronizer) Start() error {
	return rs.StartOnce("RegistrySynchronizer", func() error {
		go rs.run()
		return nil
	})
}

func (rs *RegistrySynchronizer) Close() error {
	if !rs.OkayToStop() {
		return errors.New("RegistrySynchronizer is already stopped")
	}
	close(rs.chStop)
	rs.wgDone.Wait()
	return nil
}

func (rs *RegistrySynchronizer) run() {
	rs.wgDone.Add(1)
	ticker := time.NewTicker(rs.interval)
	defer rs.wgDone.Done()
	defer ticker.Stop()

	for {
		select {
		case <-rs.chStop:
			return
		case <-ticker.C:
			rs.syncRegistry()
		}
	}
}

func (rs *RegistrySynchronizer) syncRegistry() {
	contractAddress := rs.job.KeeperSpec.ContractAddress
	logger.Debugf("syncing registry %s", contractAddress.Hex())

	var err error
	defer func() {
		logger.ErrorIf(err, fmt.Sprintf("unable to sync registry %s", contractAddress.Hex()))
	}()

	registry, err := rs.newSyncedRegistry(rs.job)
	if err != nil {
		return
	}
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	if err = rs.orm.UpsertRegistry(ctx, &registry); err != nil {
		return
	}
	if err = rs.addNewUpkeeps(registry); err != nil {
		return
	}
	if err = rs.deleteCanceledUpkeeps(registry); err != nil {
		return
	}
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

	// batch sync registries
	wg := sync.WaitGroup{}
	wg.Add(int(countOnContract - nextUpkeepID))
	chSyncUpkeepQueue := make(chan struct{}, syncUpkeepQueueSize)

	done := func() {
		select {
		case <-rs.chStop:
		case <-chSyncUpkeepQueue:
		}
		wg.Done()
	}

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
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	return rs.orm.BatchDeleteUpkeeps(ctx, reg.ID, canceled)
}

func (rs *RegistrySynchronizer) syncUpkeep(registry Registry, upkeepID int64, doneCallback func()) {
	defer doneCallback()

	var err error
	defer func() {
		logger.ErrorIf(err, fmt.Sprintf("unable to sync upkeep #%d on registry %s", upkeepID, registry.ContractAddress.Hex()))
	}()

	upkeepConfig, err := rs.contract.GetUpkeep(nil, big.NewInt(int64(upkeepID)))
	if err != nil {
		return
	}
	positioningConstant, err := calcPositioningConstant(upkeepID, registry.ContractAddress, registry.NumKeepers)
	if err != nil {
		return
	}
	newUpkeep := UpkeepRegistration{
		CheckData:           upkeepConfig.CheckData,
		ExecuteGas:          int32(upkeepConfig.ExecuteGas),
		RegistryID:          registry.ID,
		PositioningConstant: positioningConstant,
		UpkeepID:            upkeepID,
	}
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	err = rs.orm.UpsertUpkeep(ctx, &newUpkeep)
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
