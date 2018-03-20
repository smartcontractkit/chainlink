package models

import (
	"time"

	"github.com/smartcontractkit/chainlink/utils"
	"go.uber.org/multierr"
)

type AssignmentSpec struct {
	Assignment Assignment `json:"assignment"`
}

type Assignment struct {
	Subtasks []Subtask `json:"subtasks"`
}

type Subtask struct {
	Type   string `json:"adapterType"`
	Params JSON   `json:"adapterParams"`
}

func (s AssignmentSpec) ConvertToJobSpec() (JobSpec, error) {
	tasks := []TaskSpec{}
	var merr error
	for _, st := range s.Assignment.Subtasks {
		params, err := st.Params.Add("type", st.Type)
		multierr.Append(merr, err)
		tasks = append(tasks, TaskSpec{
			Type:   st.Type,
			Params: params,
		})
	}
	j := JobSpec{
		ID:        utils.NewBytes32ID(),
		CreatedAt: Time{Time: time.Now()},
		Tasks:     tasks,
	}

	return j, merr
}
