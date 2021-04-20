package cron

import (
	"context"

	cronParser "github.com/robfig/cron/v3"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

// Cron runs a cron jobSpec from a CronSpec
type Cron struct {
	jobID int32

	ctx            context.Context
	cancel         context.CancelFunc
	cronRunner     *cronParser.Cron
	done           chan struct{}
	logger         *logger.Logger
	pipelineSpec   pipeline.Spec
	pipelineRunner pipeline.Runner
	Schedule       string
}

// NewCronFromJobSpec() - instantiates a Cron singleton to execute the given job pipelineSpec
func NewCronFromJobSpec(
	jobSpec job.Job,
	runner pipeline.Runner,
) (*Cron, error) {

	cronSpec := jobSpec.CronSpec
	spec := jobSpec.PipelineSpec

	cronLogger := logger.CreateLogger(
		logger.Default.With(
			"jobID", jobSpec.ID,
			"schedule", cronSpec.CronSchedule,
		),
	)

	ctx, cancel := context.WithCancel(context.Background())

	return &Cron{
		jobID:          jobSpec.ID,
		logger:         cronLogger,
		Schedule:       cronSpec.CronSchedule,
		pipelineRunner: runner,
		pipelineSpec:   *spec,
		cronRunner:     cronParser.New(),
		ctx:            ctx,
		cancel:         cancel,
	}, nil
}

// Start implements the job.Service interface.
func (cron *Cron) Start() error {
	cron.logger.Debug("Cron: Starting")
	go gracefulpanic.WrapRecover(func() {
		defer close(cron.done)
		cron.cronRunner.Run()
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

// run() runs the cron jobSpec in the pipeline pipelineRunner
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
	_, _, err := cron.pipelineRunner.ExecuteAndInsertNewRun(cron.ctx, cron.pipelineSpec, pipeline.JSONSerializable{}, *cron.logger, true)
	if err != nil {
		cron.logger.Errorf("Error executing new run for jobSpec ID %v", cron.jobID)
	}
}
