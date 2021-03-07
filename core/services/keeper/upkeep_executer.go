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
			executer.processActiveUpkeeps()
		}
	}
}

func (executer UpkeepExecuter) processActiveUpkeeps() {
	// Keepers could miss their turn in the turn taking algo if they are too overloaded
	// with work because processActiveUpkeeps() blocks
	logger.Debug("received new block, running checkUpkeep for keeper registrations")

	activeUpkeepss, err := executer.keeperORM.EligibleUpkeeps(executer.blockHeight.Load())
	if err != nil {
		logger.Errorf("unable to load active registrations: %v", err)
		return
	}

	for _, reg := range activeUpkeepss {
		executer.concurrentExecute(reg)
	}
}

func (executer UpkeepExecuter) concurrentExecute(upkeep UpkeepRegistration) {
	executer.executionQueue <- struct{}{}
	go executer.execute(upkeep)
}

// execute will call checkForUpkeep and, if it succeeds, triger a job on the CL node
// DEV: must perform call "manually" because abigen wrapper can only send tx
func (executer UpkeepExecuter) execute(upkeep UpkeepRegistration) {
	defer func() { <-executer.executionQueue }()

	msg, err := constructCheckUpkeepCallMsg(upkeep)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debugf("Checking upkeep on registry: %s, upkeepID %d", upkeep.Registry.ContractAddress.Hex(), upkeep.UpkeepID)

	checkUpkeepResult, err := executer.ethClient.CallContract(context.Background(), msg, nil)
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

	err = executer.keeperORM.InsertEthTXForUpkeep(upkeep, performTxData)
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
