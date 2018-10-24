package services

import (
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

func ExportedExecuteRunAtBlock(
	run *models.JobRun,
	store *store.Store,
	input models.RunResult,
) (*models.JobRun, error) {
	return executeRun(run, store)
}

func ExportedChannelForRun(jr JobRunner, runID string) chan<- struct{} {
	return jr.channelForRun(runID)
}

func ExportedResumeRuns(jr JobRunner) error {
	return jr.resumeRuns()
}

func ExportedWorkerCount(jr JobRunner) int {
	return jr.workerCount()
}
