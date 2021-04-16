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

// Cron runs a cron jobSpec from a CronSpec
type Cron struct {
	jobID int32

	logger   *logger.Logger
	runner   pipeline.Runner
	Schedule string

	pipelineSpec pipeline.Spec
}

func NewCron(jobID int32, logger *logger.Logger, schedule string, runner pipeline.Runner, spec pipeline.Spec) (*Cron, error) {
	return &Cron{
		jobID:        jobID,
		logger:       logger,
		Schedule:     schedule,
		runner:       runner,
		pipelineSpec: spec,
	}, nil
}

// NewFromJobSpec() - instantiates a Cron singleton to execute the given job pipelineSpec
func NewFromJobSpec(
	jobSpec job.Job,
	runner pipeline.Runner,
) (*Cron, error) {
	cronSpec := jobSpec.CronSpec
	spec := jobSpec.PipelineSpec

	if err := validateCronSpec(*cronSpec); err != nil {
		return nil, err
	}

	cronLogger := logger.CreateLogger(
		logger.Default.With(
			"jobID", jobSpec.ID,
			"schedule", cronSpec.CronSchedule,
		),
	)

	return NewCron(jobSpec.ID, cronLogger, cronSpec.CronSchedule, runner, *spec)
}

// validateCronSpec() - validates the cron job spec and included fields
func validateCronSpec(cronSpec job.CronSpec) error {
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
func (cron *Cron) Start() error {
	cron.logger.Debug("Cron: Starting")

	go gracefulpanic.WrapRecover(func() {
		cron.run()
	})

	return nil
}

// Close implements the job.Service interface. It stops this instance from
// polling, cleaning up resources.
func (cron *Cron) Close() error {
	cron.logger.Debug("Cron: Closign")
	return nil
}

// run() runs the cron jobSpec in the pipeline runner
func (cron *Cron) run() {

	c := cronParser.New()
	_, err := c.AddFunc(cron.Schedule, func() {
		cron.runPipeline()
	})
	if err != nil {
		cron.logger.Errorf("Error running cron job(id: %d): %v", cron.jobID, err)
	}
	c.Start()
}

func (cron *Cron) runPipeline() {
	ctx := context.Background()
	_, _, err := cron.runner.ExecuteAndInsertNewRun(ctx, cron.pipelineSpec, pipeline.JSONSerializable{}, *cron.logger, true)
	if err != nil {
		cron.logger.Errorf("Error executing new run for jobSpec ID %v", cron.jobID)
	}
}
