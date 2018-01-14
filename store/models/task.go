package models

import (
	"encoding/json"
)

type Task struct {
	Type   string          `json:"type" storm:"index"`
	Params json.RawMessage `json:"params,omitempty"`
}

type TaskRun struct {
	Task
	ID     string `storm:"id,index,unique"`
	Status string
	Result RunResult
}

func (tr TaskRun) Completed() bool {
	return tr.Status == "completed"
}

func (tr TaskRun) Errored() bool {
	return tr.Status == "errored"
}
