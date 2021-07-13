package cron

import (
	"fmt"

	"github.com/robfig/cron/v3"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Cron runs a cron jobSpec from a CronSpec
type Cron struct {
	cronRunner     *cron.Cron
	logger         *logger.Logger
	jobSpec        job.Job
	pipelineRunner pipeline.Runner
	chStop         chan struct{}
}

// NewCronFromJobSpec instantiates a job that executes on a predefined schedule.
func NewCronFromJobSpec(
	jobSpec job.Job,
	pipelineRunner pipeline.Runner,
) (*Cron, error) {
	cronLogger := logger.CreateLogger(
		logger.Default.With(
			"jobID", jobSpec.ID,
			"schedule", jobSpec.CronSpec.CronSchedule,
		),
	)

	return &Cron{
		cronRunner:     cronRunner(),
		logger:         cronLogger,
		jobSpec:        jobSpec,
		pipelineRunner: pipelineRunner,
		chStop:         make(chan struct{}),
	}, nil
}

// Start implements the job.Service interface.
func (cr *Cron) Start() error {
	cr.logger.Debug("Cron: Starting")

	_, err := cr.cronRunner.AddFunc(cr.jobSpec.CronSpec.CronSchedule, cr.runPipeline)
	if err != nil {
		cr.logger.Errorw(fmt.Sprintf("Error running cron job %d", cr.jobSpec.ID), "error", err, "schedule", cr.jobSpec.CronSpec.CronSchedule, "jobID", cr.jobSpec.ID)
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

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"databaseID":    cr.jobSpec.ID,
			"externalJobID": cr.jobSpec.ExternalJobID,
			"name":          cr.jobSpec.Name.ValueOrZero(),
		},
		"jobRun": map[string]interface{}{
			"meta": map[string]interface{}{},
		},
	})

	run := pipeline.NewRun(*cr.jobSpec.PipelineSpec, vars)

	_, err := cr.pipelineRunner.Run(ctx, &run, *cr.logger, false)
	if err != nil {
		cr.logger.Errorf("Error executing new run for jobSpec ID %v", cr.jobSpec.ID)
	}
}

func cronRunner() *cron.Cron {
	return cron.New(cron.WithSeconds())
}
