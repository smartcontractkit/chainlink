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

func (self JobRun) ForLogger(kvs ...interface{}) []interface{} {
	output := []interface{}{
		"job", self.JobID,
		"run", self.ID,
		"status", self.Status,
	}

	if self.Result.HasError() {
		output = append(output, "error", self.Result.Error())
	}

	return append(kvs, output...)
}
