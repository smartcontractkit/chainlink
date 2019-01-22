package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/utils"
	"go.uber.org/multierr"
	null "gopkg.in/guregu/null.v3"
)

// AssignmentSpec represents a specification of work to be given to an oracle
// this consists of an Assignment and a Schedule
type AssignmentSpec struct {
	Assignment Assignment `json:"assignment"`
	Schedule   Schedule   `json:"schedule"`
}

// Assignment contains all the subtasks to perform
type Assignment struct {
	Subtasks []Subtask `json:"subtasks"`
}

// Subtask is a step taken by the oracle to complete an assignment
type Subtask struct {
	Type   string `json:"adapterType"`
	Params JSON   `json:"adapterParams"`
}

// Schedule defines the frequency to run the Assignment
// Schedule uses standard cron syntax
type Schedule struct {
	EndAt       Time        `json:"endAt"`
	Hour        null.String `json:"hour"`
	Minute      null.String `json:"minute"`
	DayOfMonth  null.String `json:"dayOfMonth"`
	MonthOfYear null.String `json:"monthOfYear"`
	DayOfWeek   null.String `json:"dayOfWeek"`
	RunAt       []Time      `json:"runAt"`
}

// Snapshot captures the result of an individual subtask
type Snapshot struct {
	Details JSON        `json:"details"`
	ID      string      `json:"xid"`
	Error   null.String `json:"error"`
	Pending bool        `json:"pending"`
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
			Type: "cron",
			InitiatorParams: InitiatorParams{
				Schedule: s.Schedule.toCron(),
			},
		})
	}

	return initiators
}

// ConvertToJobSpec converts an AssignmentSpec to a JobSpec
func (s AssignmentSpec) ConvertToJobSpec() (JobSpec, error) {
	var merr error
	tasks := []TaskSpec{}
	for _, st := range s.Assignment.Subtasks {
		tt, err := NewTaskType(st.Type)
		multierr.Append(merr, err)

		tasks = append(tasks, TaskSpec{
			Type:   tt,
			Params: st.Params,
		})
	}
	initiators := []Initiator{{Type: "web"}}
	for _, r := range s.Schedule.RunAt {
		initiators = append(initiators, Initiator{
			Type: "runAt",
			InitiatorParams: InitiatorParams{
				Time: r,
			},
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
	if j.EndAt.Time.IsZero() {
		j.EndAt.Valid = false
	}
	return j, merr
}

func addCronToSchedule(s Schedule, it Initiator) Schedule {
	t := strings.Split(it.Schedule.String(), " ")

	tk := make([]null.String, len(t))
	for i, v := range t {
		if v != "*" {
			tk[i] = null.StringFrom(v)
		}
	}

	s.Minute = tk[1]
	s.Hour = tk[2]
	s.DayOfMonth = tk[3]
	s.MonthOfYear = tk[4]
	s.DayOfWeek = tk[5]

	return s
}

func removeTypeFromParams(s string) (JSON, error) {
	var m map[string]interface{}

	json.Unmarshal([]byte(s), &m)
	if _, ok := m["type"]; ok {
		delete(m, "type")
	}

	var err error
	if b, err := json.Marshal(m); err == nil {
		return ParseJSON(b)
	}

	return JSON{}, err
}

func buildAssignment(ts []TaskSpec) (Assignment, error) {
	var merr error
	st := []Subtask{}

	for _, t := range ts {
		var err error
		t.Params, err = removeTypeFromParams(t.Params.String())
		if err != nil {
			multierr.Append(merr, err)
		}

		st = append(st, Subtask{
			Type:   t.Type.String(),
			Params: t.Params,
		})
	}

	a := Assignment{
		Subtasks: st,
	}

	return a, merr
}

func buildScheduleFromJobSpec(j JobSpec) Schedule {
	var s Schedule
	for _, r := range j.Initiators {
		switch r.Type {
		case InitiatorCron:
			s = addCronToSchedule(s, r)
		case InitiatorRunAt:
			s.RunAt = append(s.RunAt, r.Time)
		}
	}
	s.EndAt.Time = j.EndAt.Time

	return s
}

// ConvertToAssignment converts JobSpec to AssignmentSpec
func ConvertToAssignment(j JobSpec) (AssignmentSpec, error) {
	var merr error

	a, err := buildAssignment(j.Tasks)
	merr = multierr.Append(merr, err)

	s := buildScheduleFromJobSpec(j)

	as := AssignmentSpec{
		Assignment: a,
		Schedule:   s,
	}

	return as, merr
}

// ConvertToSnapshot convert given RunResult to a Snapshot
func ConvertToSnapshot(rr RunResult) Snapshot {
	return Snapshot{
		Details: rr.Data,
		ID:      rr.CachedJobRunID,
		Error:   rr.ErrorMessage,
		Pending: rr.Status.PendingBridge(),
	}
}
