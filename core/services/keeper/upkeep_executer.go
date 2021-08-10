package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum"
	"github.com/smartcontractkit/chainlink/core/logger"
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

// UpkeepExecuter fulfills Service and HeadBroadcastable interfaces
var _ job.Service = (*UpkeepExecuter)(nil)
var _ httypes.HeadTrackable = (*UpkeepExecuter)(nil)

type UpkeepExecuter struct {
	chStop          chan struct{}
	ethClient       eth.Client
	config          Config
	executionQueue  chan struct{}
	headBroadcaster httypes.HeadBroadcasterRegistry
	job             job.Job
	mailbox         *utils.Mailbox
	orm             ORM
	pr              pipeline.Runner
	wgDone          sync.WaitGroup
	utils.StartStopOnce
}

func NewUpkeepExecuter(
	job job.Job,
	orm ORM,
	pr pipeline.Runner,
	ethClient eth.Client,
	headBroadcaster httypes.HeadBroadcaster,
	config Config,
) *UpkeepExecuter {
	return &UpkeepExecuter{
		chStop:          make(chan struct{}),
		ethClient:       ethClient,
		executionQueue:  make(chan struct{}, executionQueueSize),
		headBroadcaster: headBroadcaster,
		job:             job,
		mailbox:         utils.NewMailbox(1),
		config:          config,
		orm:             orm,
		pr:              pr,
		wgDone:          sync.WaitGroup{},
		StartStopOnce:   utils.StartStopOnce{},
	}
}

func (executer *UpkeepExecuter) Start() error {
	return executer.StartOnce("UpkeepExecuter", func() error {
		executer.wgDone.Add(2)
		go executer.run()
		latestHead, unsubscribeHeads := executer.headBroadcaster.Subscribe(executer)
		if latestHead != nil {
			executer.mailbox.Deliver(*latestHead)
		}
		go func() {
			defer unsubscribeHeads()
			defer executer.wgDone.Done()
			<-executer.chStop
		}()
		return nil
	})
}

func (executer *UpkeepExecuter) Close() error {
	return executer.StopOnce("UpkeepExecuter", func() error {
		close(executer.chStop)
		executer.wgDone.Wait()
		return nil
	})
}

func (executer *UpkeepExecuter) Connect(head *models.Head) error { return nil }

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
	item, exists := executer.mailbox.Retrieve()
	if !exists {
		logger.Info("UpkeepExecuter: no head to retrieve. It might have been skipped")
		return
	}

	head, ok := item.(models.Head)
	if !ok {
		logger.Errorf("expected `models.Head`, got %T", head)
		return
	}

	logger.Debugw("UpkeepExecuter: checking active upkeeps", "blockheight", head.Number, "jobID", executer.job.ID)

	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()

	activeUpkeeps, err := executer.orm.EligibleUpkeepsForRegistry(
		ctx,
		executer.job.KeeperSpec.ContractAddress,
		head.Number,
		executer.config.KeeperMaximumGracePeriod(),
	)
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

	msg, err := executer.constructCheckUpkeepCallMsg(upkeep)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debugw("UpkeepExecuter: checking upkeep", logArgs...)

	ctxService, cancel := utils.ContextFromChan(executer.chStop)
	defer cancel()

	checkUpkeepResult, err := executer.ethClient.CallContract(ctxService, msg, nil)
	if err != nil {
		revertReason, err2 := eth.ExtractRevertReasonFromRPCError(err)
		if err2 != nil {
			revertReason = fmt.Sprintf("unknown revert reason: error during extraction: %v", err2)
		}
		logArgs = append(logArgs, "revertReason", revertReason)
		logger.Debugw(fmt.Sprintf("UpkeepExecuter: checkUpkeep failed: %v", err), logArgs...)
		return
	}

	performTxData, err := constructPerformUpkeepTxData(checkUpkeepResult, upkeep.UpkeepID)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debugw("UpkeepExecuter: performing upkeep", logArgs...)

	// Save a run indicating we performed an upkeep.
	f := time.Now()
	var runErrors pipeline.RunErrors
	if err == nil {
		runErrors = pipeline.RunErrors{null.String{}}
	} else {
		runErrors = pipeline.RunErrors{null.StringFrom(errors.Wrap(err, "UpkeepExecuter: failed to construct upkeep txdata").Error())}
	}

	var etx bulletprooftxmanager.EthTx
	ctxQuery, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	err = postgres.GormTransaction(ctxQuery, executer.orm.DB, func(dbtx *gorm.DB) (err error) {
		etx, err = executer.orm.CreateEthTransactionForUpkeep(dbtx, upkeep, performTxData)
		if err != nil {
			return errors.Wrap(err, "failed to create eth_tx for upkeep")
		}

		// NOTE: this is the block that initiated the run, not the block height when broadcast nor the block
		// that the tx gets confirmed in. This is fine because this grace period is just used as a fallback
		// in case we miss the UpkeepPerformed log or the tx errors. It does not need to be exact.
		err = executer.orm.SetLastRunHeightForUpkeepOnJob(dbtx, executer.job.ID, upkeep.UpkeepID, headNumber)
		if err != nil {
			return errors.Wrap(err, "failed to set last run height for upkeep")
		}

		_, err = executer.pr.InsertFinishedRun(dbtx, pipeline.Run{
			State:          pipeline.RunStatusCompleted,
			PipelineSpecID: executer.job.PipelineSpecID,
			Meta: pipeline.JSONSerializable{
				Val: map[string]interface{}{"eth_tx_id": etx.ID},
			},
			Errors: runErrors,
			Outputs: pipeline.JSONSerializable{Val: []interface{}{
				fmt.Sprintf("queued tx from %v to %v txdata %v", etx.FromAddress, etx.ToAddress, hex.EncodeToString(etx.EncodedPayload)),
			}},
			CreatedAt:  start,
			FinishedAt: null.TimeFrom(f),
		}, nil, false)
		if err != nil {
			return errors.Wrap(err, "UpkeepExecuter: failed to insert finished run")
		}
		return nil
	})
	if err != nil {
		logger.Errorw("UpkeepExecuter: failed to update database state", "err", err)
	}

	// TODO: Remove in
	// https://app.clubhouse.io/chainlinklabs/story/6065/hook-keeper-up-to-use-tasks-in-the-pipeline
	elapsed := time.Since(start)
	pipeline.PromPipelineTaskExecutionTime.WithLabelValues(fmt.Sprintf("%d", executer.job.ID), executer.job.Name.String, "", job.Keeper.String()).Set(float64(elapsed))
	var status string
	if runErrors.HasError() || err != nil {
		status = "error"
		pipeline.PromPipelineRunErrors.WithLabelValues(fmt.Sprintf("%d", executer.job.ID), executer.job.Name.String).Inc()
	} else {
		status = "completed"
	}
	pipeline.PromPipelineRunTotalTimeToCompletion.WithLabelValues(fmt.Sprintf("%d", executer.job.ID), executer.job.Name.String).Set(float64(elapsed))
	pipeline.PromPipelineTasksTotalFinished.WithLabelValues(fmt.Sprintf("%d", executer.job.ID), executer.job.Name.String, "", job.Keeper.String(), status).Inc()
}

func (executer *UpkeepExecuter) constructCheckUpkeepCallMsg(upkeep UpkeepRegistration) (ethereum.CallMsg, error) {
	checkPayload, err := RegistryABI.Pack(
		checkUpkeep,
		big.NewInt(int64(upkeep.UpkeepID)),
		upkeep.Registry.FromAddress.Address(),
	)
	if err != nil {
		return ethereum.CallMsg{}, err
	}

	to := upkeep.Registry.ContractAddress.Address()
	gasLimit := executer.config.KeeperRegistryCheckGasOverhead() + uint64(upkeep.Registry.CheckGas) +
		executer.config.KeeperRegistryPerformGasOverhead() + upkeep.ExecuteGas
	msg := ethereum.CallMsg{
		From: utils.ZeroAddress,
		To:   &to,
		Gas:  gasLimit,
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
