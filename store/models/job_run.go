package models

import (
	"time"
)

type JobRun struct {
	ID        string    `storm:"id,index,unique"`
	JobID     string    `storm:"index"`
	Status    string    `storm:"index"`
	CreatedAt time.Time `storm:"index"`
	Result    RunResult `storm:"inline"`
	TaskRuns  []TaskRun `storm:"inline"`
}

func (jr JobRun) ForLogger(kvs ...interface{}) []interface{} {
	output := []interface{}{
		"job", jr.JobID,
		"run", jr.ID,
		"status", jr.Status,
	}

	if jr.Result.HasError() {
		output = append(output, "error", jr.Result.Error())
	}

	return append(kvs, output...)
}

func (jr JobRun) UnfinishedTaskRuns() []TaskRun {
	unfinished := jr.TaskRuns
	for _, tr := range jr.TaskRuns {
		if tr.Completed() {
			unfinished = unfinished[1:]
		} else if tr.Errored() {
			return []TaskRun{}
		} else {
			return unfinished
		}
	}
	return unfinished
}

func (jr JobRun) NextTaskRun() TaskRun {
	return jr.UnfinishedTaskRuns()[0]
}
