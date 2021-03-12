package keeper

import (
	"context"
	"errors"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"go.uber.org/atomic"
	"gorm.io/gorm"
)

const (
	checkUpkeep        = "checkUpkeep"
	performUpkeep      = "performUpkeep"
	executionQueueSize = 10
)

func NewUpkeepExecutor(
	db *gorm.DB,
	ethClient eth.Client,
) *UpkeepExecutor {
	return &UpkeepExecutor{
		blockHeight:    atomic.NewInt64(0),
		wgDone:         sync.WaitGroup{},
		ethClient:      ethClient,
		orm:            NewORM(db),
		executionQueue: make(chan struct{}, executionQueueSize),
		chStop:         make(chan struct{}),
		chSignalRun:    make(chan struct{}, 1),
		StartStopOnce:  utils.StartStopOnce{},
	}
}

// UpkeepExecutor fulfills HeadTrackable interface
var _ store.HeadTrackable = (*UpkeepExecutor)(nil)

type UpkeepExecutor struct {
	blockHeight *atomic.Int64
	wgDone      sync.WaitGroup
	ethClient   eth.Client
	orm         ORM

	executionQueue chan struct{}
	chStop         chan struct{}
	chSignalRun    chan struct{}

	utils.StartStopOnce
}

func (executor *UpkeepExecutor) Start() error {
	return executor.StartOnce("UpkeepExecutor", func() error {
		go executor.run()
		return nil
	})
}

func (executor *UpkeepExecutor) Close() error {
	if !executor.OkayToStop() {
		return errors.New("UpkeepExecutor is already stopped")
	}
	close(executor.chStop)
	executor.wgDone.Wait()
	return nil
}

func (executor *UpkeepExecutor) Connect(head *models.Head) error {
	return nil
}

func (executor *UpkeepExecutor) Disconnect() {}

func (executor *UpkeepExecutor) OnNewLongestChain(ctx context.Context, head models.Head) {
	executor.blockHeight.Store(head.Number)
	// avoid blocking if signal already in buffer
	select {
	case executor.chSignalRun <- struct{}{}:
	default:
	}
}

func (executor *UpkeepExecutor) run() {
	executor.wgDone.Add(1)
	defer executor.wgDone.Done()
	for {
		select {
		case <-executor.chStop:
			return
		case <-executor.chSignalRun:
			executor.processActiveUpkeeps()
		}
	}
}

func (executor *UpkeepExecutor) processActiveUpkeeps() {
	// Keepers could miss their turn in the turn taking algo if they are too overloaded
	// with work because processActiveUpkeeps() blocks
	blockHeight := executor.blockHeight.Load()
	logger.Debug("received new block, running checkUpkeep for keeper registrations", "blockheight", blockHeight)

	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	activeUpkeeps, err := executor.orm.EligibleUpkeeps(ctx, blockHeight)
	if err != nil {
		logger.Errorf("unable to load active registrations: %v", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(activeUpkeeps))
	done := func() { <-executor.executionQueue; wg.Done() }
	for _, reg := range activeUpkeeps {
		executor.executionQueue <- struct{}{}
		go executor.execute(reg, done)
	}

	wg.Wait()
}

// execute will call checkForUpkeep and, if it succeeds, trigger a job on the CL node
// DEV: must perform contract call "manually" because abigen wrapper can only send tx
func (executor *UpkeepExecutor) execute(upkeep UpkeepRegistration, done func()) {
	defer done()

	msg, err := constructCheckUpkeepCallMsg(upkeep)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debugf("Checking upkeep on registry: %s, upkeepID %d", upkeep.Registry.ContractAddress.Hex(), upkeep.UpkeepID)

	ctxService, cancel := utils.ContextFromChan(executor.chStop)
	defer cancel()

	checkUpkeepResult, err := executor.ethClient.CallContract(ctxService, msg, nil)
	if err != nil {
		logger.Debugf("checkUpkeep failed on registry: %s, upkeepID %d", upkeep.Registry.ContractAddress.Hex(), upkeep.UpkeepID)
		return
	}

	performTxData, err := constructPerformUpkeepTxData(checkUpkeepResult, upkeep.UpkeepID)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debugf("Performing upkeep on registry: %s, upkeepID %d", upkeep.Registry.ContractAddress.Hex(), upkeep.UpkeepID)

	ctxQuery, _ := postgres.DefaultQueryCtx()
	ctxCombined, cancel := utils.CombinedContext(executor.chStop, ctxQuery)
	defer cancel()

	err = executor.orm.CreateEthTransactionForUpkeep(ctxCombined, upkeep, performTxData)
	if err != nil {
		logger.Error(err)
	}
}

func constructCheckUpkeepCallMsg(upkeep UpkeepRegistration) (ethereum.CallMsg, error) {
	checkPayload, err := RegistryABI.Pack(
		checkUpkeep,
		big.NewInt(int64(upkeep.UpkeepID)),
		upkeep.Registry.FromAddress.Address(),
	)
	if err != nil {
		return ethereum.CallMsg{}, err
	}

	to := upkeep.Registry.ContractAddress.Address()
	msg := ethereum.CallMsg{
		From: utils.ZeroAddress,
		To:   &to,
		Gas:  uint64(upkeep.Registry.CheckGas),
		Data: checkPayload,
	}

	return msg, nil
}

func constructPerformUpkeepTxData(checkUpkeepResult []byte, upkeepID int64) ([]byte, error) {
	unpackedResult, err := RegistryABI.Unpack(checkUpkeep, checkUpkeepResult)
	if err != nil {
		return nil, err
	}

	performData, ok := unpackedResult[0].([]byte)
	if !ok {
		return nil, errors.New("checkupkeep payload not as expected")
	}

	performTxData, err := RegistryABI.Pack(
		performUpkeep,
		big.NewInt(upkeepID),
		performData,
	)
	if err != nil {
		return nil, err
	}

	return performTxData, nil
}
