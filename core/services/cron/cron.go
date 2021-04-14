package cron

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/robfig/cron"
	cronParser "github.com/robfig/cron/v3"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/operator_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store"
)

// CronJob runs a cron jobSpec from a CronJobSpec
type CronJob struct {
	jobID int32

	config   Config
	logger   *logger.Logger
	orm      *orm
	runner   pipeline.Runner
	Schedule string
	store    *store.Store

	jobSpec      job.CronJobSpec
	pipelineSpec pipeline.Spec
}

func NewCronJob(jobID int32, store store.Store, config Config, logger *logger.Logger, schedule string, runner pipeline.Runner, spec pipeline.Spec, orm *orm) (*CronJob, error) {
	return &CronJob{
		jobID:        jobID,
		config:       config,
		logger:       logger,
		Schedule:     schedule,
		runner:       runner,
		pipelineSpec: spec,
		orm:          orm,
		store:        &store,
	}, nil
}

// NewFromJobSpec() - instantiates a CronJob singleton to execute the given job pipelineSpec
func NewFromJobSpec(
	jobSpec job.Job,
	store store.Store,
	config Config,
	runner pipeline.Runner,
	orm *orm,
) (*CronJob, error) {
	cronSpec := jobSpec.CronRequestSpec
	spec := jobSpec.PipelineSpec

	var gasLimit uint64
	if cronSpec.EthGasLimit == 0 {
		gasLimit = store.Config.EthGasLimitDefault()
	} else {
		gasLimit = cronSpec.EthGasLimit
	}
	config.EthGasLimit = gasLimit

	if err := validateCronJobSpec(*cronSpec); err != nil {
		return nil, err
	}

	cronLogger := logger.CreateLogger(
		logger.Default.With(
			"jobID", jobSpec.ID,
			"schedule", cronSpec.CronSchedule,
		),
	)

	return NewCronJob(jobSpec.ID, store, config, cronLogger, cronSpec.CronSchedule, runner, *spec, orm)
}

// validateCronJobSpec() - validates the cron job spec and included fields
func validateCronJobSpec(cronSpec job.CronJobSpec) error {
	if cronSpec.CronSchedule == "" {
		return fmt.Errorf("schedule must not be empty")
	}
	_, err := cron.Parse(cronSpec.CronSchedule)
	if err != nil {
		return err
	}

	return nil
}

// Start implements the job.Service interface.
func (cron *CronJob) Start() error {
	cron.logger.Debug("Starting cron jobSpec")

	go gracefulpanic.WrapRecover(func() {
		cron.run()
	})

	return nil
}

// Close implements the job.Service interface. It stops this instance from
// polling, cleaning up resources.
func (cron *CronJob) Close() error {
	cron.logger.Debug("Closing cron jobSpec")
	return nil
}

// TODO: Do we use native cron.Schedule() to run or AddFunc with string schedule?
// run() runs the cron jobSpec in the pipeline runner
func (cron *CronJob) run() {
	defer cron.runner.Close()

	c := cronParser.New()
	_, err := c.AddFunc(cron.Schedule, func() {
		cron.runPipeline()
	})
	if err != nil {
		cron.logger.Errorf("Error running cron job(id: %d): %v", cron.jobID, err)
	}
	c.Start()
}

var OperatorABI = eth.MustGetABI(operator_wrapper.OperatorABI)

func (cron *CronJob) runPipeline() {
	ctx := context.Background()
	_, results, err := cron.runner.ExecuteAndInsertNewRun(ctx, cron.pipelineSpec, *cron.logger)
	if err != nil {
		cron.logger.Errorf("error executing new run for jobSpec ID %v", cron.jobID)
	}

	result, err := results.SingularResult()
	if err != nil {
		cron.logger.Errorf("error getting singular result for jobSpec ID %v", cron.jobID)
	}
	if result.Error != nil {
		cron.logger.Errorf("error getting singular result: %v", result.Error)
	}

	// TODO: callbackAddress (address) - fromAddress? ? - empty string? nil?
	fromAddress, err := cron.store.GetRoundRobinAddress()
	if err != nil {
		cron.logger.Errorf("error getting fom address for cron job %d: %v", cron.jobID, err)
	}

	// TODO: payment - payment amount that is released for oracle? (1?) (configurable?) (uint256)
	payment, err := utils.EVMWordBigInt(cron.jobSpec.OraclePayment.ToInt())
	if err != nil {
		cron.logger.Errorf("error parsing payment: %v", err)
	}

	// TODO: requestId? -> runId? (bytes32) ?
	// TODO: callbackFunctionId (bytes4) ?
	callBackFunctionId := [4]byte{0, 1, 2} // TODO: what should this be? (bytes4)

	// TODO: expiration experiation the node should respond by before the request can cancel (uint256)
	expiration, err := utils.EVMWordBigInt(big.NewInt(time.Now().Unix()))
	if err != nil {
		cron.logger.Errorf("error parsing payment: %v", err)
	}

	// TODO: Verify OperatorABI by creating transaction w/ operator_wrapper + Verifying transaction creation
	payload, err := OperatorABI.Pack("fulfillOracleRequest", payment, fromAddress, callBackFunctionId, expiration, result.Value)
	if err != nil {
		cron.logger.Errorf("abi.Pack failed: %v", err)
	}

	// TODO: make acceptance & unit test list...

	// Send Eth Transaction with Payload
	err = cron.orm.CreateEthTransaction(ctx, fromAddress, cron.jobSpec.ToAddress.Address(), payload, cron.config.EthGasLimit, cron.config.MaxUnconfirmedTransactions)
	if err != nil {
		cron.logger.Errorf("error creating eth tx for cron job %d: %v", cron.jobID, err)
	}
}
