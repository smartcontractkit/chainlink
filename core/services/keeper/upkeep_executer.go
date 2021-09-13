package keeper

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const (
	executionQueueSize = 10
)

// UpkeepExecuter fulfills Service and HeadTrackable interfaces
var (
	_ job.Service           = (*UpkeepExecuter)(nil)
	_ httypes.HeadTrackable = (*UpkeepExecuter)(nil)
)

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
func (ex *UpkeepExecuter) OnNewLongestChain(_ context.Context, head models.Head) {
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

	svcLogger := ex.logger.With("blockNum", headNumber, "upkeepID", upkeep.UpkeepID)
	svcLogger.Debug("checking upkeep")

	ctxService, cancel := utils.ContextFromChanWithDeadline(ex.chStop, time.Minute)
	defer cancel()

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"jobID":                 ex.job.ID,
			"fromAddress":           upkeep.Registry.FromAddress.String(),
			"contractAddress":       upkeep.Registry.ContractAddress.String(),
			"upkeepID":              upkeep.UpkeepID,
			"performUpkeepGasLimit": upkeep.ExecuteGas + ex.orm.config.KeeperRegistryPerformGasOverhead(),
			"checkUpkeepGasLimit": ex.config.KeeperRegistryCheckGasOverhead() + uint64(upkeep.Registry.CheckGas) +
				ex.config.KeeperRegistryPerformGasOverhead() + upkeep.ExecuteGas,
		},
	})

	run := pipeline.NewRun(*ex.job.PipelineSpec, vars)
	if _, err := ex.pr.Run(ctxService, &run, *ex.logger, true, nil); err != nil {
		ex.logger.WithError(err).Errorw("failed executing run")
		return
	}

	// Only after task runs where a tx was broadcast
	if run.State == pipeline.RunStatusCompleted {
		err := ex.orm.SetLastRunHeightForUpkeepOnJob(ctxService, ex.job.ID, upkeep.UpkeepID, headNumber)
		if err != nil {
			ex.logger.WithError(err).Errorw("failed to set last run height for upkeep")
		}
	}
}
