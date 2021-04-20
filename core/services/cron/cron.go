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

	ctx          context.Context
	cancel       context.CancelFunc
	cronRunner   *cronParser.Cron
	done         chan struct{}
	logger       *logger.Logger
	pipelineSpec pipeline.Spec
	runner       pipeline.Runner
	Schedule     string
}

// NewCronFromJobSpec() - instantiates a Cron singleton to execute the given job pipelineSpec
func NewCronFromJobSpec(
	jobSpec job.Job,
	runner pipeline.Runner,
) (*Cron, error) {
	ctx, cancel := context.WithCancel(context.Background())

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

	return &Cron{
		jobID:        jobSpec.ID,
		logger:       cronLogger,
		Schedule:     cronSpec.CronSchedule,
		runner:       runner,
		pipelineSpec: *spec,
		cronRunner:   cronParser.New(),
		ctx:          ctx,
		cancel:       cancel,
	}, nil
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
		defer close(cron.done)
	})
	return nil
}

// Close implements the job.Service interface. It stops this instance from
// polling, cleaning up resources.
func (cron *Cron) Close() error {
	cron.logger.Debug("Cron: Closing")
	cron.cancel() // Cancel any inflight cron runs.
	cron.cronRunner.Stop()
	<-cron.done

	return nil
}

// run() runs the cron jobSpec in the pipeline runner
func (cron *Cron) run() {
	_, err := cron.cronRunner.AddFunc(cron.Schedule, func() {
		cron.runPipeline()
	})
	if err != nil {
		cron.logger.Errorf("Error running cron job(id: %d): %v", cron.jobID, err)
	}
	cron.cronRunner.Start()
}

func (cron *Cron) runPipeline() {
	_, _, err := cron.runner.ExecuteAndInsertNewRun(cron.ctx, cron.pipelineSpec, pipeline.JSONSerializable{}, *cron.logger, true)
	if err != nil {
		cron.logger.Errorf("Error executing new run for jobSpec ID %v", cron.jobID)
	}
}
