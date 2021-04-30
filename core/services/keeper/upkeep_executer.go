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
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const (
	checkUpkeep          = "checkUpkeep"
	performUpkeep        = "performUpkeep"
	executionQueueSize   = 10
	queuedEthTransaction = "successfully queued performUpkeep eth transaction"
)

// UpkeepExecuter fulfills Service and HeadBroadcastable interfaces
var _ job.Service = (*UpkeepExecuter)(nil)
var _ services.HeadBroadcastable = (*UpkeepExecuter)(nil)

type UpkeepExecuter struct {
	chStop            chan struct{}
	ethClient         eth.Client
	executionQueue    chan struct{}
	headBroadcaster   *services.HeadBroadcaster
	job               job.Job
	mailbox           *utils.Mailbox
	maxGracePeriod    int64
	maxUnconfirmedTXs uint64
	orm               ORM
	pr                pipeline.Runner
	wgDone            sync.WaitGroup
	utils.StartStopOnce
}

func NewUpkeepExecuter(
	job job.Job,
	db *gorm.DB,
	pr pipeline.Runner,
	ethClient eth.Client,
	headBroadcaster *services.HeadBroadcaster,
	config *orm.Config,
) *UpkeepExecuter {
	return &UpkeepExecuter{
		chStop:            make(chan struct{}),
		ethClient:         ethClient,
		executionQueue:    make(chan struct{}, executionQueueSize),
		headBroadcaster:   headBroadcaster,
		job:               job,
		mailbox:           utils.NewMailbox(1),
		maxUnconfirmedTXs: config.EthMaxUnconfirmedTransactions(),
		maxGracePeriod:    config.KeeperMaximumGracePeriod(),
		orm:               NewORM(db),
		pr:                pr,
		wgDone:            sync.WaitGroup{},
		StartStopOnce:     utils.StartStopOnce{},
	}
}

func (executer *UpkeepExecuter) Start() error {
	return executer.StartOnce("UpkeepExecuter", func() error {
		executer.wgDone.Add(2)
		go executer.run()
		unsubscribe := executer.headBroadcaster.Subscribe(executer)
		go func() {
			defer unsubscribe()
			defer executer.wgDone.Done()
			<-executer.chStop
		}()
		return nil
	})
}

func (executer *UpkeepExecuter) Close() error {
	if !executer.OkayToStop() {
		return errors.New("UpkeepExecuter is already stopped")
	}
	close(executer.chStop)
	executer.wgDone.Wait()
	return nil
}

func (executer *UpkeepExecuter) OnNewLongestChain(ctx context.Context, head models.Head) {
	executer.mailbox.Deliver(head)
}

func (executer *UpkeepExecuter) run() {
	defer executer.wgDone.Done()
	for {
		select {
		case <-executer.chStop:
			return
		case <-executer.mailbox.Notify():
			executer.processActiveUpkeeps()
		}
	}
}

func (executer *UpkeepExecuter) processActiveUpkeeps() {
	// Keepers could miss their turn in the turn taking algo if they are too overloaded
	// with work because processActiveUpkeeps() blocks
	head, ok := executer.mailbox.Retrieve().(models.Head)
	if !ok {
		logger.Errorf("expected `models.Head`, got %T", head)
		return
	}

	logger.Debugw("UpkeepExecuter: checking active upkeeps", "blockheight", head.Number, "jobID", executer.job.ID)

	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	activeUpkeeps, err := executer.orm.EligibleUpkeeps(ctx, head.Number, executer.maxGracePeriod)
	if err != nil {
		logger.Errorf("unable to load active registrations: %v", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(activeUpkeeps))
	done := func() { <-executer.executionQueue; wg.Done() }
	for _, reg := range activeUpkeeps {
		executer.executionQueue <- struct{}{}
		go executer.execute(reg, head.Number, done)
	}

	wg.Wait()
}

// execute will call checkForUpkeep and, if it succeeds, trigger a job on the CL node
// DEV: must perform contract call "manually" because abigen wrapper can only send tx
func (executer *UpkeepExecuter) execute(upkeep UpkeepRegistration, headNumber int64, done func()) {
	defer done()
	start := time.Now()
	logArgs := []interface{}{
		"jobID", executer.job.ID,
		"blockNum", headNumber,
		"registryAddress", upkeep.Registry.ContractAddress.Hex(),
		"upkeepID", upkeep.UpkeepID,
	}

	msg, err := constructCheckUpkeepCallMsg(upkeep)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debugw("UpkeepExecuter: checking upkeep", logArgs...)

	ctxService, cancel := utils.ContextFromChan(executer.chStop)
	defer cancel()

	checkUpkeepResult, err := executer.ethClient.CallContract(ctxService, msg, nil)
	if err != nil {
		logger.Debugw(fmt.Sprintf("UpkeepExecuter: checkUpkeep failed: %v", err), logArgs...)
		return
	}

	performTxData, err := constructPerformUpkeepTxData(checkUpkeepResult, upkeep.UpkeepID)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debugw("UpkeepExecuter: performing upkeep", logArgs...)

	ctxQuery, _ := postgres.DefaultQueryCtx()
	ctxCombined, cancel := utils.CombinedContext(executer.chStop, ctxQuery)
	defer cancel()

	etx, err := executer.orm.CreateEthTransactionForUpkeep(ctxCombined, upkeep, performTxData, executer.maxUnconfirmedTXs)
	if err != nil {
		logger.Error(err)
	}

	// Save a run indicating we performed an upkeep.
	f := time.Now()
	_, err = executer.pr.InsertFinishedRun(ctxCombined, pipeline.Run{
		PipelineSpecID: executer.job.PipelineSpecID,
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
	ctxCombined, cancel = utils.CombinedContext(executer.chStop, ctxQuery)
	defer cancel()
	// DEV: this is the block that initiated the run, not the block height when broadcast nor the block
	// that the tx gets confirmed in. This is fine because this grace period is just used as a fallback
	// in case we miss the UpkeepPerformed log or the tx errors. It does not need to be exact.
	err = executer.orm.SetLastRunHeightForUpkeepOnJob(ctxCombined, executer.job.ID, upkeep.UpkeepID, headNumber)
	if err != nil {
		logger.Errorw("UpkeepExecuter: unable to setLastRunHeightForUpkeep for upkeep", logArgs...)
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
