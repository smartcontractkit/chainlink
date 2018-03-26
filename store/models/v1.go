package models

import (
	"fmt"
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
	EndAt       Time        `json:"endAt"`
	Hour        null.String `json:"hour"`
	Minute      null.String `json:"minute"`
	DayOfMonth  null.String `json:"dayOfMonth"`
	MonthOfYear null.String `json:"monthOfYear"`
	DayOfWeek   null.String `json:"dayOfWeek"`
	RunAt       []Time      `json:"runAt"`
}

func (s Schedule) hasCron() bool {
	return s.Minute.Valid || s.Hour.Valid || s.DayOfMonth.Valid ||
		s.MonthOfYear.Valid || s.DayOfWeek.Valid
}

func (s Schedule) toCron() Cron {
	return Cron(fmt.Sprintf("0 %v %v %v %v %v",
		cronUnitOrDefault(s.Minute),
		cronUnitOrDefault(s.Hour),
		cronUnitOrDefault(s.DayOfMonth),
		cronUnitOrDefault(s.MonthOfYear),
		cronUnitOrDefault(s.DayOfWeek),
	))
}

func cronUnitOrDefault(s null.String) string {
	if s.Valid {
		return s.String
	}
	return "*"
}

func appendCronInitiator(initiators []Initiator, s AssignmentSpec) []Initiator {
	if s.Schedule.hasCron() {
		initiators = append(initiators, Initiator{
			Type:     "cron",
			Schedule: s.Schedule.toCron(),
		})
	}

	return initiators
}

func (s AssignmentSpec) ConvertToJobSpec() (JobSpec, error) {
	var merr error
	tasks := []TaskSpec{}
	for _, st := range s.Assignment.Subtasks {
		params, err := st.Params.Add("type", st.Type)
		multierr.Append(merr, err)
		tasks = append(tasks, TaskSpec{
			Type:   st.Type,
			Params: params,
		})
	}
	initiators := []Initiator{{Type: "web"}}
	for _, r := range s.Schedule.RunAt {
		initiators = append(initiators, Initiator{
			Type: "runAt",
			Time: r,
		})
	}
	initiators = appendCronInitiator(initiators, s)

	j := JobSpec{
		ID:         utils.NewBytes32ID(),
		CreatedAt:  Time{Time: time.Now()},
		Tasks:      tasks,
		EndAt:      null.TimeFrom(s.Schedule.EndAt.Time),
		Initiators: initiators,
	}

	return j, merr
}
