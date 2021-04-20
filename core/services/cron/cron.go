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
	pipelineRunner pipeline.Runner,
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
		pipelineRunner: pipelineRunner,
		pipelineSpec:   *spec,
		cronRunner:     cronParser.New(),
		ctx:            ctx,
		cancel:         cancel,
	}, nil
}

// Start implements the job.Service interface.
func (cr *Cron) Start() error {
	cr.logger.Debug("Cron: Starting")
	go gracefulpanic.WrapRecover(func() {
		cr.run()
		defer close(cr.done)
	})
	return nil
}

// Close implements the job.Service interface. It stops this instance from
// polling, cleaning up resources.
func (cr *Cron) Close() error {
	cr.logger.Debug("Cron: Closing")
	cr.cancel() // Cancel any inflight cr runs.
	cr.cronRunner.Stop()
	<-cr.done

	return nil
}

// run() runs the cron jobSpec in the pipeline pipelineRunner
func (cr *Cron) run() {
	_, err := cr.cronRunner.AddFunc(cr.Schedule, func() {
		cr.runPipeline()
	})
	if err != nil {
		cr.logger.Errorf("Error running cr job(id: %d): %v", cr.jobID, err)
	}
	cr.cronRunner.Start()
}

func (cr *Cron) runPipeline() {
	_, _, err := cr.pipelineRunner.ExecuteAndInsertNewRun(cr.ctx, cr.pipelineSpec, pipeline.JSONSerializable{}, *cr.logger, true)
	if err != nil {
		cr.logger.Errorf("Error executing new run for jobSpec ID %v", cr.jobID)
	}
}
