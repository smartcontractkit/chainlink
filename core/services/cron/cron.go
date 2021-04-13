package cron

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/robfig/cron"
	cronParser "github.com/robfig/cron/v3"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store"
)

// CronJob runs a cron jobSpec from a CronJobSpec
type CronJob struct {
	jobID int32

	config      Config
	contractABI abi.ABI
	logger      *logger.Logger
	orm         *orm
	runner      pipeline.Runner
	Schedule    string
	store       *store.Store

	jobSpec      job.CronJobSpec
	pipelineSpec pipeline.Spec
}

func NewCronJob(jobID int32, store store.Store, config Config, contractABI abi.ABI, logger *logger.Logger, schedule string, runner pipeline.Runner, spec pipeline.Spec, orm *orm) (*CronJob, error) {
	return &CronJob{
		jobID:        jobID,
		config:       config,
		contractABI:  contractABI,
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
	contractABI abi.ABI,
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

	return NewCronJob(jobSpec.ID, store, config, contractABI, cronLogger, cronSpec.CronSchedule, runner, *spec, orm)
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

// run() runs the cron jobSpec in the pipeline runner
func (cron *CronJob) run() {
	defer cron.runner.Close()

	c := cronParser.New()
	// TODO: Do we use native cron.Schedule() to run or AddFunc with string schedule?
	_, err := c.AddFunc(cron.Schedule, func() {
		cron.runPipeline()
	})
	if err != nil {
		cron.logger.Errorf("Error running cron job(id: %d): %v", cron.jobID, err)
	}
	c.Start()
}

func (cron *CronJob) runPipeline() {
	ctx := context.Background()
	runID, results, err := cron.runner.ExecuteAndInsertNewRun(ctx, cron.pipelineSpec, *cron.logger)
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

	fromAddress, err := cron.store.GetRoundRobinAddress()
	if err != nil {
		cron.logger.Errorf("error getting fom address for cron job %d: %v", cron.jobID, err)
	}

	// TODO: results.Value -> array of one answer
	// TODO: how do we get values from result.Value?
	payload, err := cron.contractABI.Pack("fulfillOracleRequest", runID, "payment...", fromAddress, "callBackFunctionId?", "expiration?", "data?")
	if err != nil {
		cron.logger.Errorf("abi.Pack failed: %v", err)
	}

	// Send Eth Transaction with Payload
	err = cron.orm.CreateEthTransaction(ctx, fromAddress, cron.jobSpec.ToAddress.Address(), payload, cron.config.EthGasLimit, cron.config.MaxUnconfirmedTransactions)
	if err != nil {
		cron.logger.Errorf("error creating eth tx for cron job %d: %v", cron.jobID, err)
	}
}
