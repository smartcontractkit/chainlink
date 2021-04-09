package cron

import (
	"context"
	"fmt"

	"github.com/robfig/cron"
	cronParser "github.com/robfig/cron/v3"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

// CronJob runs a cron jobSpec from a CronSpec
type CronJob struct {
	jobID int32

	logger   *logger.Logger
	runner   pipeline.Runner
	Schedule string

	jobSpec      job.CronSpec
	pipelineSpec pipeline.Spec
}

func NewCronJob(jobID int32, logger *logger.Logger, schedule string, runner pipeline.Runner, spec pipeline.Spec) (*CronJob, error) {
	return &CronJob{
		jobID:        jobID,
		logger:       logger,
		Schedule:     schedule,
		runner:       runner,
		pipelineSpec: spec,
	}, nil
}

// NewFromJobSpec() - instantiates a CronJob singleton to execute the given job pipelineSpec
func NewFromJobSpec(
	jobSpec job.Job,
	runner pipeline.Runner,
) (*CronJob, error) {
	cronSpec := jobSpec.CronSpec
	spec := jobSpec.PipelineSpec

	if err := validateCronJobSpec(*cronSpec); err != nil {
		return nil, err
	}

	cronLogger := logger.CreateLogger(
		logger.Default.With(
			"jobID", jobSpec.ID,
			"schedule", cronSpec.CronSchedule,
		),
	)

	return NewCronJob(jobSpec.ID, cronLogger, cronSpec.CronSchedule, runner, *spec)
}

// validateCronJobSpec() - validates the cron job spec and included fields
func validateCronJobSpec(cronSpec job.CronSpec) error {
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
	_, _, err := cron.runner.ExecuteAndInsertNewRun(ctx, cron.pipelineSpec, *cron.logger)
	if err != nil {
		cron.logger.Errorf("Error executing new run for jobSpec ID %v", cron.jobID)
	}
}
