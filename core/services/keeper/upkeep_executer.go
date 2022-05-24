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

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
	bigmath "github.com/smartcontractkit/chainlink/core/utils/big_math"
)

const (
	executionQueueSize = 10
)

// UpkeepExecuter fulfills Service and HeadTrackable interfaces
var (
	_ job.ServiceCtx        = (*UpkeepExecuter)(nil)
	_ httypes.HeadTrackable = (*UpkeepExecuter)(nil)
)

var (
	promCheckUpkeepExecutionTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "keeper_check_upkeep_execution_time",
		Help: "Time taken to fully execute the check upkeep logic",
	},
		[]string{"upkeepID"},
	)
)

// UpkeepExecuter implements the logic to communicate with KeeperRegistry
type UpkeepExecuter struct {
	chStop          chan struct{}
	ethClient       evmclient.Client
	config          Config
	executionQueue  chan struct{}
	headBroadcaster httypes.HeadBroadcasterRegistry
	gasEstimator    gas.Estimator
	job             job.Job
	mailbox         *utils.Mailbox[*evmtypes.Head]
	orm             ORM
	pr              pipeline.Runner
	logger          logger.Logger
	wgDone          sync.WaitGroup
	utils.StartStopOnce
}

// NewUpkeepExecuter is the constructor of UpkeepExecuter
func NewUpkeepExecuter(
	job job.Job,
	orm ORM,
	pr pipeline.Runner,
	ethClient evmclient.Client,
	headBroadcaster httypes.HeadBroadcaster,
	gasEstimator gas.Estimator,
	logger logger.Logger,
	config Config,
) *UpkeepExecuter {
	return &UpkeepExecuter{
		chStop:          make(chan struct{}),
		ethClient:       ethClient,
		executionQueue:  make(chan struct{}, executionQueueSize),
		headBroadcaster: headBroadcaster,
		gasEstimator:    gasEstimator,
		job:             job,
		mailbox:         utils.NewMailbox[*evmtypes.Head](1),
		config:          config,
		orm:             orm,
		pr:              pr,
		logger:          logger.Named("UpkeepExecuter"),
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
	head, exists := ex.mailbox.Retrieve()
	if !exists {
		ex.logger.Info("no head to retrieve. It might have been skipped")
		return
	}

	ex.logger.Debugw("checking active upkeeps", "blockheight", head.Number)

	registry, err := ex.orm.RegistryByContractAddress(ex.job.KeeperSpec.ContractAddress)
	if err != nil {
		ex.logger.Error(errors.Wrap(err, "unable to load registry"))
		return
	}

	var activeUpkeeps []UpkeepRegistration
	if ex.config.KeeperTurnFlagEnabled() {
		turnBinary, err2 := ex.turnBlockHashBinary(registry, head, ex.config.KeeperTurnLookBack())
		if err2 != nil {
			ex.logger.Error(errors.Wrap(err2, "unable to get turn block number hash"))
			return
		}
		activeUpkeeps, err2 = ex.orm.NewEligibleUpkeepsForRegistry(
			ex.job.KeeperSpec.ContractAddress,
			head.Number,
			ex.config.KeeperMaximumGracePeriod(),
			turnBinary)
		if err2 != nil {
			ex.logger.Error(errors.Wrap(err2, "unable to load active registrations"))
			return
		}
	} else {
		activeUpkeeps, err = ex.orm.EligibleUpkeepsForRegistry(
			ex.job.KeeperSpec.ContractAddress,
			head.Number,
			ex.config.KeeperMaximumGracePeriod(),
		)
		if err != nil {
			ex.logger.Error(errors.Wrap(err, "unable to load active registrations"))
			return
		}
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
}

// execute triggers the pipeline run
func (ex *UpkeepExecuter) execute(upkeep UpkeepRegistration, head *evmtypes.Head, done func()) {
	defer done()

	start := time.Now()
	svcLogger := ex.logger.With("jobID", ex.job.ID, "blockNum", head.Number, "upkeepID", upkeep.UpkeepID)
	svcLogger.Debugw("checking upkeep", "lastRunBlockHeight", upkeep.LastRunBlockHeight, "lastKeeperIndex", upkeep.LastKeeperIndex)

	ctxService, cancel := utils.ContextFromChanWithDeadline(ex.chStop, time.Minute)
	defer cancel()

	evmChainID := ""
	if ex.job.KeeperSpec.EVMChainID != nil {
		evmChainID = ex.job.KeeperSpec.EVMChainID.String()
	}

	var gasPrice, gasTipCap, gasFeeCap *big.Int
	if ex.config.KeeperCheckUpkeepGasPriceFeatureEnabled() {
		price, fee, err := ex.estimateGasPrice(upkeep)
		if err != nil {
			svcLogger.Error(errors.Wrap(err, "estimating gas price"))
			return
		}
		gasPrice, gasTipCap, gasFeeCap = price, fee.TipCap, fee.FeeCap

		// Make sure the gas price is at least as large as the basefee to avoid ErrFeeCapTooLow error from geth during eth call.
		// If head.BaseFeePerGas, we assume it is a EIP-1559 chain.
		// Note: gasPrice will be nil if EvmEIP1559DynamicFees is enabled.
		if head.BaseFeePerGas != nil && head.BaseFeePerGas.ToInt().BitLen() > 0 {
			baseFee := addBuffer(head.BaseFeePerGas.ToInt(), ex.config.KeeperBaseFeeBufferPercent())
			if gasPrice == nil || gasPrice.Cmp(baseFee) < 0 {
				gasPrice = baseFee
			}
		}
	}

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"jobID":                 ex.job.ID,
			"fromAddress":           upkeep.Registry.FromAddress.String(),
			"contractAddress":       upkeep.Registry.ContractAddress.String(),
			"upkeepID":              upkeep.UpkeepID,
			"prettyID":              upkeep.PrettyID(),
			"performUpkeepGasLimit": upkeep.ExecuteGas + ex.orm.config.KeeperRegistryPerformGasOverhead(),
			"checkUpkeepGasLimit": ex.config.KeeperRegistryCheckGasOverhead() + uint64(upkeep.Registry.CheckGas) +
				ex.config.KeeperRegistryPerformGasOverhead() + upkeep.ExecuteGas,
			"gasPrice":   gasPrice,
			"gasTipCap":  gasTipCap,
			"gasFeeCap":  gasFeeCap,
			"evmChainID": evmChainID,
		},
	})

	run := pipeline.NewRun(*ex.job.PipelineSpec, vars)
	if _, err := ex.pr.Run(ctxService, &run, svcLogger, true, nil); err != nil {
		svcLogger.Error(errors.Wrap(err, "failed executing run"))
		return
	}

	// Only after task runs where a tx was broadcast
	if run.State == pipeline.RunStatusCompleted {
		err := ex.orm.SetLastRunInfoForUpkeepOnJob(ex.job.ID, upkeep.UpkeepID, head.Number, upkeep.Registry.FromAddress, pg.WithParentCtx(ctxService))
		if err != nil {
			svcLogger.Error(errors.Wrap(err, "failed to set last run height for upkeep"))
		}
		svcLogger.Debugw("execute pipeline status completed", "fromAddr", upkeep.Registry.FromAddress)

		elapsed := time.Since(start)
		promCheckUpkeepExecutionTime.
			WithLabelValues(upkeep.PrettyID()).
			Set(float64(elapsed))
	}
}

func (ex *UpkeepExecuter) estimateGasPrice(upkeep UpkeepRegistration) (gasPrice *big.Int, fee gas.DynamicFee, err error) {
	var performTxData []byte
	performTxData, err = Registry1_1ABI.Pack(
		"performUpkeep", // performUpkeep is same across registry ABI versions
		upkeep.UpkeepID.ToInt(),
		common.Hex2Bytes("1234"), // placeholder
	)
	if err != nil {
		return nil, fee, errors.Wrap(err, "unable to construct performUpkeep data")
	}

	if ex.config.EvmEIP1559DynamicFees() {
		fee, _, err = ex.gasEstimator.GetDynamicFee(upkeep.ExecuteGas)
		fee.TipCap = addBuffer(fee.TipCap, ex.config.KeeperGasTipCapBufferPercent())
	} else {
		gasPrice, _, err = ex.gasEstimator.GetLegacyGas(performTxData, upkeep.ExecuteGas)
		gasPrice = addBuffer(gasPrice, ex.config.KeeperGasPriceBufferPercent())
	}
	if err != nil {
		return nil, fee, errors.Wrap(err, "unable to estimate gas")
	}

	return gasPrice, fee, nil
}

func addBuffer(val *big.Int, prct uint32) *big.Int {
	return bigmath.Div(
		bigmath.Mul(val, 100+prct),
		100,
	)
}

func (ex *UpkeepExecuter) turnBlockHashBinary(registry Registry, head *evmtypes.Head, lookback int64) (string, error) {
	turnBlock := head.Number - (head.Number % int64(registry.BlockCountPerTurn)) - lookback
	block, err := ex.ethClient.BlockByNumber(context.Background(), big.NewInt(turnBlock))
	if err != nil {
		return "", err
	}
	hashAtHeight := block.Hash()
	binaryString := fmt.Sprintf("%b", hashAtHeight.Big())
	return binaryString, nil
}
