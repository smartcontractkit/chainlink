package keeper

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"go.uber.org/atomic"
)

const (
	checkUpkeep        = "checkUpkeep"
	performUpkeep      = "performUpkeep"
	executionQueueSize = 10
)

func NewUpkeepExecuter(
	keeperORM KeeperORM,
	ethClient eth.Client,
) UpkeepExecuter {
	return UpkeepExecuter{
		blockHeight:    atomic.NewInt64(0),
		ethClient:      ethClient,
		keeperORM:      keeperORM,
		isRunning:      atomic.NewBool(false),
		executionQueue: make(chan struct{}, executionQueueSize),
		chDone:         make(chan struct{}),
		chSignalRun:    make(chan struct{}, 1),
	}
}

// UpkeepExecuter fulfills HeadTrackable interface
var _ store.HeadTrackable = UpkeepExecuter{}

type UpkeepExecuter struct {
	blockHeight *atomic.Int64
	ethClient   eth.Client
	keeperORM   KeeperORM
	isRunning   *atomic.Bool

	executionQueue chan struct{}
	chDone         chan struct{}
	chSignalRun    chan struct{}
}

func (executer UpkeepExecuter) Start() error {
	if executer.isRunning.Load() {
		return errors.New("already started")
	}
	executer.isRunning.Store(true)
	go executer.run()
	return nil
}

func (executer UpkeepExecuter) Close() error {
	close(executer.chDone)
	return nil
}

func (executer UpkeepExecuter) Connect(head *models.Head) error {
	return nil
}

func (executer UpkeepExecuter) Disconnect() {}

func (executer UpkeepExecuter) OnNewLongestChain(ctx context.Context, head models.Head) {
	executer.blockHeight.Store(head.Number)
	// avoid blocking if signal already in buffer
	select {
	case executer.chSignalRun <- struct{}{}:
	default:
	}
}

func (executer UpkeepExecuter) run() {
	for {
		select {
		case <-executer.chDone:
			return
		case <-executer.chSignalRun:
			executer.processActiveRegistrations()
		}
	}
}

func (executer UpkeepExecuter) processActiveRegistrations() {
	// Keepers could miss their turn in the turn taking algo if they are too overloaded
	// with work because processActiveRegistrations() blocks - this could be parallelized
	// but will need a cap
	logger.Debug("received new block, running checkUpkeep for keeper registrations")

	activeRegistrations, err := executer.keeperORM.EligibleUpkeeps(executer.blockHeight.Load())
	if err != nil {
		logger.Errorf("unable to load active registrations: %v", err)
		return
	}

	for _, reg := range activeRegistrations {
		executer.concurrentExecute(reg)
	}
}

func (executer UpkeepExecuter) concurrentExecute(registration UpkeepRegistration) {
	executer.executionQueue <- struct{}{}
	go executer.execute(registration)
}

// execute will call checkForUpkeep and, if it succeeds, triger a job on the CL node
func (executer UpkeepExecuter) execute(registration UpkeepRegistration) {
	// pop queue when done executing
	defer func() {
		<-executer.executionQueue
	}()

	checkPayload, err := RegistryABI.Pack(
		checkUpkeep,
		big.NewInt(int64(registration.UpkeepID)),
		registration.Registry.FromAddress.Address(),
	)
	if err != nil {
		logger.Error(err)
		return
	}

	to := registration.Registry.ContractAddress.Address()
	msg := ethereum.CallMsg{
		From: utils.ZeroAddress,
		To:   &to,
		Gas:  uint64(registration.Registry.CheckGas),
		Data: checkPayload,
	}

	logger.Debugf("Checking upkeep on registry: %s, upkeepID %d", registration.Registry.ContractAddress.Hex(), registration.UpkeepID)

	result, err := executer.ethClient.CallContract(context.Background(), msg, nil)
	if err != nil {
		logger.Debugf("checkUpkeep failed on registry: %s, upkeepID %d", registration.Registry.ContractAddress.Hex(), registration.UpkeepID)
		return
	}

	performTxData, err := constructPerformUpkeepTxData(result, registration.UpkeepID)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debugf("Performing upkeep on registry: %s, upkeepID %d", registration.Registry.ContractAddress.Hex(), registration.UpkeepID)

	err = executer.keeperORM.InsertEthTXForUpkeep(registration, performTxData)
	if err != nil {
		logger.Error(err)
	}
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
