package cron

import (
	"github.com/robfig/cron/v3"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Cron runs a cron jobSpec from a CronSpec
type Cron struct {
	cronRunner     *cron.Cron
	jobID          int32
	logger         *logger.Logger
	pipelineSpec   pipeline.Spec
	pipelineRunner pipeline.Runner
	Schedule       string
	chStop         chan struct{}
}

// NewCronFromJobSpec instantiates a job that executes on a predefined schedule.
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
		cronRunner:     cron.New(cron.WithSeconds()),
		jobID:          jobSpec.ID,
		logger:         cronLogger,
		pipelineRunner: pipelineRunner,
		pipelineSpec:   *spec,
		Schedule:       cronSpec.CronSchedule,
		chStop:         make(chan struct{}),
	}, nil
}

// Start implements the job.Service interface.
func (cr *Cron) Start() error {
	cr.logger.Debug("Cron: Starting")

	_, err := cr.cronRunner.AddFunc(cr.Schedule, cr.runPipeline)
	if err != nil {
		cr.logger.Errorf("Error running cr job(id: %d): %v", cr.jobID, err)
		return err
	}
	cr.cronRunner.Start()
	return nil
}

// Close implements the job.Service interface. It stops this job from
// running and cleans up resources.
func (cr *Cron) Close() error {
	cr.logger.Debug("Cron: Closing")
	cr.cronRunner.Stop()
	return nil
}

func (cr *Cron) runPipeline() {
	ctx, cancel := utils.ContextFromChan(cr.chStop)
	defer cancel()
	_, _, err := cr.pipelineRunner.ExecuteAndInsertFinishedRun(ctx, cr.pipelineSpec, nil, pipeline.JSONSerializable{}, *cr.logger, false)
	if err != nil {
		cr.logger.Errorf("Error executing new run for jobSpec ID %v", cr.jobID)
	}
}
