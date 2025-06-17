package keeper

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

const (
	executionQueueSize  = 10
	maxUpkeepPerformGas = 5_000_000 // Max perform gas for upkeep is 5M on all chains for v1.x
)

// UpkeepExecuter fulfills Service and HeadTrackable interfaces
var (
	_ job.ServiceCtx  = (*UpkeepExecuter)(nil)
	_ heads.Trackable = (*UpkeepExecuter)(nil)
)

var (
	promCheckUpkeepExecutionTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "keeper_check_upkeep_execution_time",
		Help: "Time taken to fully execute the check upkeep logic",
	},
		[]string{"upkeepID"},
	)
)

type UpkeepExecuterConfig interface {
	MaxGracePeriod() int64
	TurnLookBack() int64
	Registry() config.Registry
}

// UpkeepExecuter implements the logic to communicate with KeeperRegistry
type UpkeepExecuter struct {
	services.StateMachine
	chStop                 services.StopChan
	ethClient              evmclient.Client
	config                 UpkeepExecuterConfig
	executionQueue         chan struct{}
	headBroadcaster        heads.Broadcaster
	gasEstimator           gas.EvmFeeEstimator
	job                    job.Job
	mailbox                *mailbox.Mailbox[*evmtypes.Head]
	orm                    *ORM
	pr                     pipeline.Runner
	logger                 logger.Logger
	wgDone                 sync.WaitGroup
	effectiveKeeperAddress common.Address
}

// NewUpkeepExecuter is the constructor of UpkeepExecuter
func NewUpkeepExecuter(
	job job.Job,
	orm *ORM,
	pr pipeline.Runner,
	ethClient evmclient.Client,
	headBroadcaster heads.Broadcaster,
	gasEstimator gas.EvmFeeEstimator,
	logger logger.Logger,
	config UpkeepExecuterConfig,
	effectiveKeeperAddress common.Address,
) *UpkeepExecuter {
	return &UpkeepExecuter{
		chStop:                 make(services.StopChan),
		ethClient:              ethClient,
		executionQueue:         make(chan struct{}, executionQueueSize),
		headBroadcaster:        headBroadcaster,
		gasEstimator:           gasEstimator,
		job:                    job,
		mailbox:                mailbox.NewSingle[*evmtypes.Head](),
		config:                 config,
		orm:                    orm,
		pr:                     pr,
		effectiveKeeperAddress: effectiveKeeperAddress,
		logger:                 logger.Named("UpkeepExecuter"),
	}
}

// Start starts the upkeep executer logic
func (ex *UpkeepExecuter) Start(context.Context) error {
	return ex.StartOnce("UpkeepExecuter", func() error {
		ex.wgDone.Add(2)
		go ex.run()
		latestHead, unsubscribeHeads := ex.headBroadcaster.Subscribe(ex)
		if latestHead != nil {
			ex.mailbox.Deliver(latestHead)
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
func (ex *UpkeepExecuter) OnNewLongestChain(_ context.Context, head *evmtypes.Head) {
	ex.mailbox.Deliver(head)
}

func (ex *UpkeepExecuter) run() {
	defer ex.wgDone.Done()
	ctx, cancel := ex.chStop.NewCtx()
	defer cancel()
	for {
		select {
		case <-ex.chStop:
			return
		case <-ex.mailbox.Notify():
			ex.processActiveUpkeeps(ctx)
		}
	}
}

func (ex *UpkeepExecuter) processActiveUpkeeps(ctx context.Context) {
	// Keepers could miss their turn in the turn taking algo if they are too overloaded
	// with work because processActiveUpkeeps() blocks
	head, exists := ex.mailbox.Retrieve()
	if !exists {
		ex.logger.Info("no head to retrieve. It might have been skipped")
		return
	}

	ex.logger.Debugw("checking active upkeeps", "blockheight", head.Number)

	registry, err := ex.orm.RegistryByContractAddress(ctx, ex.job.KeeperSpec.ContractAddress)
	if err != nil {
		ex.logger.Error(errors.Wrap(err, "unable to load registry"))
		return
	}

	var activeUpkeeps []UpkeepRegistration
	turnBinary, err2 := ex.turnBlockHashBinary(ctx, registry, head, ex.config.TurnLookBack())
	if err2 != nil {
		ex.logger.Error(errors.Wrap(err2, "unable to get turn block number hash"))
		return
	}
	activeUpkeeps, err2 = ex.orm.EligibleUpkeepsForRegistry(
		ctx,
		ex.job.KeeperSpec.ContractAddress,
		head.Number,
		ex.config.MaxGracePeriod(),
		turnBinary)
	if err2 != nil {
		ex.logger.Error(errors.Wrap(err2, "unable to load active registrations"))
		return
	}

	if head.Number%10 == 0 {
		// Log this once every 10 blocks
		fetchedUpkeepIDs := make([]string, len(activeUpkeeps))
		for i, activeUpkeep := range activeUpkeeps {
			fetchedUpkeepIDs[i] = NewUpkeepIdentifier(activeUpkeep.UpkeepID).String()
		}
		ex.logger.Debugw("Fetched list of active upkeeps", "blockNum", head.Number, "active upkeeps list", fetchedUpkeepIDs)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(activeUpkeeps))
	done := func() {
		<-ex.executionQueue
		wg.Done()
	}
	for _, reg := range activeUpkeeps {
		ex.executionQueue <- struct{}{}
		go ex.execute(reg, head, done)
	}

	wg.Wait()
	ex.logger.Debugw("Finished checking upkeeps", "blockNum", head.Number)
}

// execute triggers the pipeline run
func (ex *UpkeepExecuter) execute(upkeep UpkeepRegistration, head *evmtypes.Head, done func()) {
	defer done()

	start := time.Now()
	svcLogger := ex.logger.With("jobID", ex.job.ID, "blockNum", head.Number, "upkeepID", upkeep.UpkeepID)
	svcLogger.Debugw("checking upkeep", "lastRunBlockHeight", upkeep.LastRunBlockHeight, "lastKeeperIndex", upkeep.LastKeeperIndex)

	ctxService, cancel := ex.chStop.CtxWithTimeout(time.Minute)
	defer cancel()

	evmChainID := ""
	if ex.job.KeeperSpec.EVMChainID != nil {
		evmChainID = ex.job.KeeperSpec.EVMChainID.String()
	}

	var gasPrice, gasTipCap, gasFeeCap *assets.Wei
	// effectiveKeeperAddress is always fromAddress when forwarding is not enabled.
	// when forwarding is enabled, effectiveKeeperAddress is on-chain forwarder.
	vars := pipeline.NewVarsFrom(buildJobSpec(ex.job, ex.effectiveKeeperAddress, upkeep, ex.config.Registry(), gasPrice, gasTipCap, gasFeeCap, evmChainID))

	// DotDagSource in database is empty because all the Keeper pipeline runs make use of the same observation source
	ex.job.PipelineSpec.DotDagSource = pipeline.KeepersObservationSource
	run := pipeline.NewRun(*ex.job.PipelineSpec, vars)

	if _, err := ex.pr.Run(ctxService, run, true, nil); err != nil {
		svcLogger.Error(errors.Wrap(err, "failed executing run"))
		return
	}

	// Only after task runs where a tx was broadcast
	if run.State == pipeline.RunStatusCompleted {
		rowsAffected, err := ex.orm.SetLastRunInfoForUpkeepOnJob(ctxService, ex.job.ID, upkeep.UpkeepID, head.Number, upkeep.Registry.FromAddress)
		if err != nil {
			svcLogger.Error(errors.Wrap(err, "failed to set last run height for upkeep"))
		}
		svcLogger.Debugw("execute pipeline status completed", "fromAddr", upkeep.Registry.FromAddress, "rowsAffected", rowsAffected)

		elapsed := time.Since(start)
		promCheckUpkeepExecutionTime.
			WithLabelValues(upkeep.PrettyID()).
			Set(float64(elapsed))
	}
}

func (ex *UpkeepExecuter) turnBlockHashBinary(ctx context.Context, registry Registry, head *evmtypes.Head, lookback int64) (string, error) {
	turnBlock := head.Number - (head.Number % int64(registry.BlockCountPerTurn)) - lookback
	block, err := ex.ethClient.HeadByNumber(ctx, big.NewInt(turnBlock))
	if err != nil {
		return "", err
	}
	hashAtHeight := block.Hash
	binaryString := fmt.Sprintf("%b", hashAtHeight.Big())
	return binaryString, nil
}

func buildJobSpec(
	jb job.Job,
	effectiveKeeperAddress common.Address,
	upkeep UpkeepRegistration,
	ormConfig RegistryGasChecker,
	gasPrice *assets.Wei,
	gasTipCap *assets.Wei,
	gasFeeCap *assets.Wei,
	chainID string,
) map[string]interface{} {
	return map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"jobID":                  jb.ID,
			"fromAddress":            upkeep.Registry.FromAddress.String(),
			"effectiveKeeperAddress": effectiveKeeperAddress.String(),
			"contractAddress":        upkeep.Registry.ContractAddress.String(),
			"upkeepID":               upkeep.UpkeepID.String(),
			"prettyID":               upkeep.PrettyID(),
			"pipelineSpec": &pipeline.Spec{
				ForwardingAllowed: jb.ForwardingAllowed,
			},
			"performUpkeepGasLimit": maxUpkeepPerformGas + ormConfig.PerformGasOverhead(),
			"maxPerformDataSize":    ormConfig.MaxPerformDataSize(),
			"gasPrice":              gasPrice.ToInt(),
			"gasTipCap":             gasTipCap.ToInt(),
			"gasFeeCap":             gasFeeCap.ToInt(),
			"evmChainID":            chainID,
		},
	}
}
