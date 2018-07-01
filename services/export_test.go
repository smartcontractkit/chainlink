package services

import "github.com/smartcontractkit/chainlink/store"

func ExportedChannelForRun(jr JobRunner, runID string) chan<- store.RunRequest {
	return jr.channelForRun(runID)
}

func ExportedResumeSleepingRuns(jr JobRunner) error {
	return jr.resumeSleepingRuns()
}

func ExportedWorkerCount(jr JobRunner) int {
	return jr.workerCount()
}
