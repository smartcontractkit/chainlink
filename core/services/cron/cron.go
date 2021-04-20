package cron

import (
	cronParser "github.com/robfig/cron/v3"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Cron runs a cron jobSpec from a CronSpec
type Cron struct {
	cronRunner     *cronParser.Cron
	chDone         chan struct{}
	chStop         chan struct{}
	jobID          int32
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

	return &Cron{
		chDone:         make(chan struct{}),
		chStop:         make(chan struct{}),
		cronRunner:     cronParser.New(),
		jobID:          jobSpec.ID,
		logger:         cronLogger,
		pipelineRunner: pipelineRunner,
		pipelineSpec:   *spec,
		Schedule:       cronSpec.CronSchedule,
	}, nil
}

// Start implements the job.Service interface.
func (cr *Cron) Start() error {
	cr.logger.Debug("Cron: Starting")
	go gracefulpanic.WrapRecover(func() {
		defer close(cr.chDone)
		cr.run()
	})
	return nil
}

// Close implements the job.Service interface. It stops this instance from
// polling, cleaning up resources.
func (cr *Cron) Close() error {
	cr.logger.Debug("Cron: Closing")
	cr.cronRunner.Stop()
	close(cr.chStop)
	<-cr.chDone

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
	ctx, cancel := utils.ContextFromChan(cr.chStop)
	defer cancel()
	_, _, err := cr.pipelineRunner.ExecuteAndInsertFinishedRun(ctx, cr.pipelineSpec, pipeline.JSONSerializable{}, *cr.logger, true)
	if err != nil {
		cr.logger.Errorf("Error executing new run for jobSpec ID %v", cr.jobID)
	}
}
