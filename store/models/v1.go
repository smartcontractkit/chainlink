package models

import (
	"time"

	"github.com/smartcontractkit/chainlink/utils"
	"go.uber.org/multierr"
	null "gopkg.in/guregu/null.v3"
)

type AssignmentSpec struct {
	Assignment Assignment `json:"assignment"`
	Schedule   Schedule   `json:"schedule"`
}

type Assignment struct {
	Subtasks []Subtask `json:"subtasks"`
}

type Subtask struct {
	Type   string `json:"adapterType"`
	Params JSON   `json:"adapterParams"`
}

type Schedule struct {
	EndAt null.Time `json:"endAt"`
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
		EndAt:     s.Schedule.EndAt,
	}

	return j, merr
}
