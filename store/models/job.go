package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mrwonko/cron"
	uuid "github.com/satori/go.uuid"
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
	CreatedAt  Time        `storm:"index"`
}

func NewJob() Job {
	return Job{ID: uuid.NewV4().String(), CreatedAt: Time{Time: time.Now()}}
}

func (j Job) NewRun() *JobRun {
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

func (j Job) InitiatorsFor(t string) []Initiator {
	list := []Initiator{}
	for _, initr := range j.Initiators {
		if initr.Type == t {
			list = append(list, initr)
		}
	}
	return list
}

func (j Job) WebAuthorized() bool {
	for _, initr := range j.Initiators {
		if initr.Type == "web" {
			return true
		}
	}
	return false
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

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	newTime, err := dateparse.ParseAny(s)
	t.Time = newTime
	return err
}

func (t *Time) ISO8601() string {
	return t.UTC().Format("2006-01-02T15:04:05Z07:00")
}

func (t *Time) DurationFromNow() time.Duration {
	return t.Time.Sub(time.Now())
}

type Cron string

func (c *Cron) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return fmt.Errorf("Cron: %v", err)
	}
	if s == "" {
		return nil
	}

	_, err = cron.Parse(s)
	if err != nil {
		return fmt.Errorf("Cron: %v", err)
	}
	*c = Cron(s)
	return nil
}
