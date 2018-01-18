package models

import (
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	uuid "github.com/satori/go.uuid"
	null "gopkg.in/guregu/null.v3"
)

const (
	StatusInProgress = "in progress"
	StatusPending    = "pending"
	StatusErrored    = "errored"
	StatusCompleted  = "completed"
)

type Job struct {
	ID         string      `storm:"id,index,unique"`
	Initiators []Initiator `json:"initiators"`
	Tasks      []Task      `json:"tasks" storm:"inline"`
	StartAt    null.Time   `storm:"index"`
	EndAt      null.Time   `storm:"index"`
	CreatedAt  Time        `storm:"index"`
}

func NewJob() *Job {
	return &Job{ID: uuid.NewV4().String(), CreatedAt: Time{Time: time.Now()}}
}

func (j *Job) NewRun() *JobRun {
	taskRuns := make([]TaskRun, len(j.Tasks))
	for i, task := range j.Tasks {
		taskRuns[i] = TaskRun{
			ID:   uuid.NewV4().String(),
			Task: task,
		}
	}

	return &JobRun{
		ID:        uuid.NewV4().String(),
		JobID:     j.ID,
		CreatedAt: time.Now(),
		TaskRuns:  taskRuns,
	}
}

func (j *Job) InitiatorsFor(t string) []Initiator {
	list := []Initiator{}
	for _, initr := range j.Initiators {
		if initr.Type == t {
			list = append(list, initr)
		}
	}
	return list
}

func (j *Job) WebAuthorized() bool {
	for _, initr := range j.Initiators {
		if initr.Type == "web" {
			return true
		}
	}
	return false
}

func (j *Job) Ended(t time.Time) bool {
	if !j.EndAt.Valid {
		return false
	}
	return t.After(j.EndAt.Time)
}

func (j *Job) Started(t time.Time) bool {
	if !j.StartAt.Valid {
		return true
	}
	return t.After(j.StartAt.Time) || t.Equal(j.StartAt.Time)
}

type Initiator struct {
	ID       int            `storm:"id,increment"`
	JobID    string         `storm:"index"`
	Type     string         `json:"type" storm:"index"`
	Schedule Cron           `json:"schedule,omitempty"`
	Time     Time           `json:"time,omitempty"`
	Ran      bool           `json:"ranAt,omitempty"`
	Address  common.Address `json:"address,omitempty" storm:"index"`
}

type Task struct {
	Type   string          `json:"type" storm:"index"`
	Params json.RawMessage `json:"params,omitempty"`
}
