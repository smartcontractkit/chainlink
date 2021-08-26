package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const (
	checkUpkeep        = "checkUpkeep"
	performUpkeep      = "performUpkeep"
	executionQueueSize = 10
)

// Revert reasons
const (
	UpkeepNotNeededReason     RevertReason = "upkeep not needed"
	OutOfTurnReason           RevertReason = "keepers must take turns"
	PerformUpkeepFailedReason RevertReason = "call to perform upkeep failed"
	CheckTargetFailedReason   RevertReason = "call to check target failed"
	InsufficientFundsReason   RevertReason = "insufficient funds"
)

var (
	// debugRevertReasons contains revert reasons that should be logged with the debug log level
	debugRevertReasons = []RevertReason{
		UpkeepNotNeededReason,
		OutOfTurnReason,
		PerformUpkeepFailedReason,
		CheckTargetFailedReason,
		InsufficientFundsReason,
	}
)

// UpkeepExecuter fulfills Service and HeadTrackable interfaces
var (
	_ job.Service           = (*UpkeepExecuter)(nil)
	_ httypes.HeadTrackable = (*UpkeepExecuter)(nil)
)

// RevertReason represents the revert reason message
type RevertReason string

// IsOneOf returns true if the "rr" is one of the provided revert reasons.
func (rr RevertReason) IsOneOf(revertReasons ...RevertReason) bool {
	for _, revertReason := range revertReasons {
		if strings.Contains(string(rr), string(revertReason)) {
			return true
		}
	}
	return false
}

// UpkeepExecuter implements the logic to communicate with KeeperRegistry
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
	logger          *logger.Logger
	wgDone          sync.WaitGroup
	utils.StartStopOnce
}

// NewUpkeepExecuter is the constructor of UpkeepExecuter
func NewUpkeepExecuter(
	job job.Job,
	orm ORM,
	pr pipeline.Runner,
	ethClient eth.Client,
	headBroadcaster httypes.HeadBroadcaster,
	logger *logger.Logger,
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
		logger:          logger,
	}
}

// Start starts the upkeep executer logic
func (ex *UpkeepExecuter) Start() error {
	return ex.StartOnce("UpkeepExecuter", func() error {
		ex.wgDone.Add(2)
		go ex.run()
		latestHead, unsubscribeHeads := ex.headBroadcaster.Subscribe(ex)
		if latestHead != nil {
			ex.mailbox.Deliver(*latestHead)
		}
		go func() {
			defer unsubscribeHeads()
			defer ex.wgDone.Done()
			<-ex.chStop
		}()
		return nil
	})
}

// Close stops and closes upkeep executer
func (ex *UpkeepExecuter) Close() error {
	return ex.StopOnce("UpkeepExecuter", func() error {
		close(ex.chStop)
		ex.wgDone.Wait()
		return nil
	})
}

// OnNewLongestChain handles the given head of a new longest chain
func (ex *UpkeepExecuter) OnNewLongestChain(ctx context.Context, head models.Head) {
	ex.mailbox.Deliver(head)
}

func (ex *UpkeepExecuter) run() {
	defer ex.wgDone.Done()
	for {
		select {
		case <-ex.chStop:
			return
		case <-ex.mailbox.Notify():
			ex.processActiveUpkeeps()
		}
	}
}

func (ex *UpkeepExecuter) processActiveUpkeeps() {
	// Keepers could miss their turn in the turn taking algo if they are too overloaded
	// with work because processActiveUpkeeps() blocks
	item, exists := ex.mailbox.Retrieve()
	if !exists {
		ex.logger.Info("no head to retrieve. It might have been skipped")
		return
	}

	head, ok := item.(models.Head)
	if !ok {
		ex.logger.Errorf("expected `models.Head`, got %T", head)
		return
	}

	ex.logger.Debugw("checking active upkeeps", "blockheight", head.Number)

	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()

	activeUpkeeps, err := ex.orm.EligibleUpkeepsForRegistry(
		ctx,
		ex.job.KeeperSpec.ContractAddress,
		head.Number,
		ex.config.KeeperMaximumGracePeriod(),
	)
	if err != nil {
		ex.logger.WithError(err).Error("unable to load active registrations")
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(activeUpkeeps))
	done := func() {
		<-ex.executionQueue
		wg.Done()
	}
	for _, reg := range activeUpkeeps {
		ex.executionQueue <- struct{}{}
		go ex.execute(reg, head.Number, done)
	}

	wg.Wait()
}

// execute calls checkForUpkeep and, if it succeeds, trigger a job on the CL node
// DEV: must perform contract call "manually" because abigen wrapper can only send tx
func (ex *UpkeepExecuter) execute(upkeep UpkeepRegistration, headNumber int64, done func()) {
	defer done()
	start := time.Now()

	svcLogger := ex.logger.With("blockNum", headNumber, "upkeepID", upkeep.UpkeepID)

	msg, err := ex.constructCheckUpkeepCallMsg(upkeep)
	if err != nil {
		svcLogger.WithError(err).Error("failed to construct check upkeep call message")
		return
	}

	svcLogger.Debug("checking upkeep")

	ctxService, cancel := utils.ContextFromChan(ex.chStop)
	defer cancel()

	checkUpkeepResult, err := ex.ethClient.CallContract(ctxService, msg, nil)
	if err != nil {
		logRevertReason(svcLogger, err)
		return
	}

	performTxData, err := constructPerformUpkeepTxData(checkUpkeepResult, upkeep.UpkeepID)
	if err != nil {
		svcLogger.WithError(err).Error("failed to construct check upkeep call message")
		return
	}

	svcLogger.Debug("performing upkeep")

	// Save a run indicating we performed an upkeep.
	f := time.Now()
	var runErrors pipeline.RunErrors
	if err == nil {
		runErrors = pipeline.RunErrors{null.String{}}
	} else {
		runErrors = pipeline.RunErrors{null.StringFrom(errors.Wrap(err, "failed to construct upkeep txdata").Error())}
	}

	var etx bulletprooftxmanager.EthTx
	err = ex.orm.WithTransaction(func(ctx context.Context) error {
		etx, err = ex.orm.CreateEthTransactionForUpkeep(ctx, upkeep, performTxData)
		if err != nil {
			return errors.Wrap(err, "failed to create eth_tx for upkeep")
		}

		// NOTE: this is the block that initiated the run, not the block height when broadcast nor the block
		// that the tx gets confirmed in. This is fine because this grace period is just used as a fallback
		// in case we miss the UpkeepPerformed log or the tx errors. It does not need to be exact.
		err = ex.orm.SetLastRunHeightForUpkeepOnJob(ctx, ex.job.ID, upkeep.UpkeepID, headNumber)
		if err != nil {
			return errors.Wrap(err, "failed to set last run height for upkeep")
		}

		_, err = ex.pr.InsertFinishedRun(postgres.UnwrapGorm(postgres.TxFromContext(ctx, ex.orm.DB)), pipeline.Run{
			State:          pipeline.RunStatusCompleted,
			PipelineSpecID: ex.job.PipelineSpecID,
			Meta: pipeline.JSONSerializable{
				Val: map[string]interface{}{"eth_tx_id": etx.ID},
			},
			Errors: runErrors,
			Outputs: pipeline.JSONSerializable{Val: []interface{}{
				fmt.Sprintf("queued tx from %v to %v txdata %v", etx.FromAddress, etx.ToAddress, hex.EncodeToString(etx.EncodedPayload)),
			}},
			CreatedAt:  start,
			FinishedAt: null.TimeFrom(f),
		}, false)
		if err != nil {
			return errors.Wrap(err, "failed to insert finished run")
		}
		return nil
	})
	if err != nil {
		svcLogger.WithError(err).Error("failed to update database state")
	}

	// TODO: Remove in
	// https://app.clubhouse.io/chainlinklabs/story/6065/hook-keeper-up-to-use-tasks-in-the-pipeline
	elapsed := time.Since(start)
	pipeline.PromPipelineTaskExecutionTime.WithLabelValues(fmt.Sprintf("%d", ex.job.ID), ex.job.Name.String, "", job.Keeper.String()).Set(float64(elapsed))
	var status string
	if runErrors.HasError() || err != nil {
		status = "error"
		pipeline.PromPipelineRunErrors.WithLabelValues(fmt.Sprintf("%d", ex.job.ID), ex.job.Name.String).Inc()
	} else {
		status = "completed"
	}
	pipeline.PromPipelineRunTotalTimeToCompletion.WithLabelValues(fmt.Sprintf("%d", ex.job.ID), ex.job.Name.String).Set(float64(elapsed))
	pipeline.PromPipelineTasksTotalFinished.WithLabelValues(fmt.Sprintf("%d", ex.job.ID), ex.job.Name.String, "", job.Keeper.String(), status).Inc()
}

func (ex *UpkeepExecuter) constructCheckUpkeepCallMsg(upkeep UpkeepRegistration) (ethereum.CallMsg, error) {
	checkPayload, err := RegistryABI.Pack(
		checkUpkeep,
		big.NewInt(int64(upkeep.UpkeepID)),
		upkeep.Registry.FromAddress.Address(),
	)
	if err != nil {
		return ethereum.CallMsg{}, errors.Wrap(err, "failed to pack payload of check upkeep")
	}

	to := upkeep.Registry.ContractAddress.Address()
	gasLimit := ex.config.KeeperRegistryCheckGasOverhead() + uint64(upkeep.Registry.CheckGas) +
		ex.config.KeeperRegistryPerformGasOverhead() + upkeep.ExecuteGas

	return ethereum.CallMsg{
		From: utils.ZeroAddress,
		To:   &to,
		Gas:  gasLimit,
		Data: checkPayload,
	}, nil
}

// logRevertReason logs the given error with a log level depends on this given error context
// Default log level is error. Mapping between a reason and its log level:
//  - "upkeep not needed": Debug
func logRevertReason(logger *logger.Logger, err error) {
	revertReason, err2 := eth.ExtractRevertReasonFromRPCError(err)
	if err2 != nil {
		logger.WithError(err).Errorf("call failed and failed to extract revert reason, err2: %v", err2)
		return
	}

	logger = logger.With("revertReason", revertReason)
	if RevertReason(revertReason).IsOneOf(debugRevertReasons...) {
		logger.Debug("checkUpkeep call failed with known reason")
	} else {
		logger.Error("checkUpkeep call failed with some reason")
	}
}

func constructPerformUpkeepTxData(checkUpkeepResult []byte, upkeepID int64) ([]byte, error) {
	unpackedResult, err := RegistryABI.Unpack(checkUpkeep, checkUpkeepResult)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unpack check upkeep result")
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
		return nil, errors.Wrap(err, "failed to pack a payload of perform upkeep")
	}

	return performTxData, nil
}
