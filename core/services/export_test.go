package services

import (
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func ExportedExecuteRun(
	run *models.JobRun,
	store *store.Store,
) error {
	return executeRun(run, store)
}

func ExportedChannelForRun(jr JobRunner, runID *models.ID) chan<- struct{} {
	return jr.channelForRun(runID)
}

func ExportedResumeRunsSinceLastShutdown(jr JobRunner) error {
	return jr.resumeRunsSinceLastShutdown()
}

func ExportedWorkerCount(jr JobRunner) int {
	return jr.workerCount()
}

func ExportedNewPendingConnectionResumer(
	store *store.Store,
	resumer func(*models.JobRun, *store.Store) error,
) store.HeadTrackable {
	return &pendingConnectionResumer{
		store:   store,
		resumer: resumer,
	}
}
