package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const (
	checkUpkeep          = "checkUpkeep"
	performUpkeep        = "performUpkeep"
	executionQueueSize   = 10
	queuedEthTransaction = "successfully queued performUpkeep eth transaction"
)

// UpkeepExecutor fulfills Service and HeadBroadcastable interfaces
var _ job.Service = (*UpkeepExecutor)(nil)
var _ services.HeadBroadcastable = (*UpkeepExecutor)(nil)

type UpkeepExecutor struct {
	chStop          chan struct{}
	ethClient       eth.Client
	executionQueue  chan struct{}
	headBroadcaster *services.HeadBroadcaster
	job             job.Job
	mailbox         *utils.Mailbox
	maxGracePeriod  int64
	orm             ORM
	pr              pipeline.Runner
	wgDone          sync.WaitGroup
	utils.StartStopOnce
}

func NewUpkeepExecutor(
	job job.Job,
	db *gorm.DB,
	pr pipeline.Runner,
	ethClient eth.Client,
	headBroadcaster *services.HeadBroadcaster,
	maxGracePeriod int64,
) *UpkeepExecutor {
	return &UpkeepExecutor{
		chStop:          make(chan struct{}),
		ethClient:       ethClient,
		executionQueue:  make(chan struct{}, executionQueueSize),
		headBroadcaster: headBroadcaster,
		job:             job,
		mailbox:         utils.NewMailbox(1),
		maxGracePeriod:  maxGracePeriod,
		orm:             NewORM(db),
		pr:              pr,
		wgDone:          sync.WaitGroup{},
		StartStopOnce:   utils.StartStopOnce{},
	}
}

func (executor *UpkeepExecutor) Start() error {
	return executor.StartOnce("UpkeepExecutor", func() error {
		executor.wgDone.Add(2)
		go executor.run()
		unsubscribe := executor.headBroadcaster.Subscribe(executor)
		go func() {
			defer unsubscribe()
			defer executor.wgDone.Done()
			<-executor.chStop
		}()
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

func (executor *UpkeepExecutor) OnNewLongestChain(ctx context.Context, head models.Head) {
	executor.mailbox.Deliver(head)
}

func (executor *UpkeepExecutor) run() {
	defer executor.wgDone.Done()
	for {
		select {
		case <-executor.chStop:
			return
		case <-executor.mailbox.Notify():
			executor.processActiveUpkeeps()
		}
	}
}

func (executor *UpkeepExecutor) processActiveUpkeeps() {
	// Keepers could miss their turn in the turn taking algo if they are too overloaded
	// with work because processActiveUpkeeps() blocks
	head, ok := executor.mailbox.Retrieve().(models.Head)
	if !ok {
		logger.Errorf("expected `models.Head`, got %T", head)
		return
	}

	logger.Debugw("UpkeepExecutor: checking active upkeeps", "blockheight", head.Number, "jobID", executor.job.ID)

	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	activeUpkeeps, err := executor.orm.EligibleUpkeeps(ctx, head.Number, executor.maxGracePeriod)
	if err != nil {
		logger.Errorf("unable to load active registrations: %v", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(activeUpkeeps))
	done := func() { <-executor.executionQueue; wg.Done() }
	for _, reg := range activeUpkeeps {
		executor.executionQueue <- struct{}{}
		go executor.execute(reg, head.Number, done)
	}

	wg.Wait()
}

// execute will call checkForUpkeep and, if it succeeds, trigger a job on the CL node
// DEV: must perform contract call "manually" because abigen wrapper can only send tx
func (executor *UpkeepExecutor) execute(upkeep UpkeepRegistration, headNumber int64, done func()) {
	defer done()
	start := time.Now()
	logArgs := []interface{}{
		"jobID", executor.job.ID,
		"blockNum", headNumber,
		"registryAddress", upkeep.Registry.ContractAddress.Hex(),
		"upkeepID", upkeep.UpkeepID,
	}

	msg, err := constructCheckUpkeepCallMsg(upkeep)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debugw("UpkeepExecutor: checking upkeep", logArgs...)

	ctxService, cancel := utils.ContextFromChan(executor.chStop)
	defer cancel()

	checkUpkeepResult, err := executor.ethClient.CallContract(ctxService, msg, nil)
	if err != nil {
		logger.Debugw(fmt.Sprintf("UpkeepExecutor: checkUpkeep failed: %v", err), logArgs...)
		return
	}

	performTxData, err := constructPerformUpkeepTxData(checkUpkeepResult, upkeep.UpkeepID)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debugw("UpkeepExecutor: performing upkeep", logArgs...)

	ctxQuery, _ := postgres.DefaultQueryCtx()
	ctxCombined, cancel := utils.CombinedContext(executor.chStop, ctxQuery)
	defer cancel()

	etx, err := executor.orm.CreateEthTransactionForUpkeep(ctxCombined, upkeep, performTxData)
	if err != nil {
		logger.Error(err)
	}

	// Save a run indicating we performed an upkeep.
	f := time.Now()
	_, err = executor.pr.InsertFinishedRun(ctxCombined, pipeline.Run{
		PipelineSpecID: executor.job.PipelineSpecID,
		Meta: pipeline.JSONSerializable{
			Val: map[string]interface{}{"eth_tx_id": etx.ID},
		},
		Errors:     pipeline.RunErrors{null.String{}},
		Outputs:    pipeline.JSONSerializable{Val: fmt.Sprintf("queued tx from %v to %v txdata %v", etx.FromAddress, etx.ToAddress, hex.EncodeToString(etx.EncodedPayload))},
		CreatedAt:  start,
		FinishedAt: &f,
	}, nil, false)
	if err != nil {
		logger.Error(err)
	}

	ctxQuery, cancel = postgres.DefaultQueryCtx()
	defer cancel()
	ctxCombined, cancel = utils.CombinedContext(executor.chStop, ctxQuery)
	defer cancel()
	// DEV: this is the block that initiated the run, not the block height when broadcast nor the block
	// that the tx gets confirmed in. This is fine because this grace period is just used as a fallback
	// in case we miss the UpkeepPerformed log or the tx errors. It does not need to be exact.
	err = executor.orm.SetLastRunHeightForUpkeepOnJob(ctxCombined, executor.job.ID, upkeep.UpkeepID, headNumber)
	if err != nil {
		logger.Errorw("UpkeepExecutor: unable to setLastRunHeightForUpkeep for upkeep", logArgs...)
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
