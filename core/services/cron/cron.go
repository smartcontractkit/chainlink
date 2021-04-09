package cron

import (
	"github.com/robfig/cron/v3"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var CronParser cron.Parser

// CronJob runs a cron job from a CronJobSpec
type CronJob struct {
	jobID  int32
	logger *logger.Logger
	Clock  utils.Nower
}

func NewCronJob(jobID int32, logger *logger.Logger, clock utils.Nower) (*CronJob, error) {
	cron := &CronJob{
		jobID:  jobID,
		logger: logger,
		Clock:  clock,
	}

	return cron, nil
}

func NewJobFromSpec(
	jobSpec job.Job,
) (*CronJob, error) {
	cronSpec := jobSpec.CronRequestSpec

	//TODO: do some validation on the spec

	// create logger from spec
	cronLogger := logger.CreateLogger(
		logger.Default.With(
			"jobID", jobSpec.ID,
			"schedule", cronSpec.CronSchedule,
		),
	)

	// TODO: pull cron schedule from spec (clock) and create Schedules / Jobs (Recurring or OneTime jobs?)
	// TODO: complete logic for parsing job from job spec...run job in context / on schedule
	// TODO: pass clock to NewCronJob(...)

	return NewCronJob(jobSpec.ID, cronLogger, nil)
}
