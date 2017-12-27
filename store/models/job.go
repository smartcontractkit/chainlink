package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/mrwonko/cron"
	uuid "github.com/satori/go.uuid"
)

type Job struct {
	ID         string      `storm:"id,index,unique"`
	Initiators []Initiator `json:"initiators"`
	Tasks      []Task      `json:"tasks" storm:"inline"`
	CreatedAt  time.Time   `storm:"index"`
}

type Initiator struct {
	ID       int    `storm:"id,increment"`
	Type     string `json:"type" storm:"index"`
	Schedule Cron   `json:"schedule"`
	JobID    string `storm:"index"`
}

type Cron string

type Time struct {
	time.Time
}

func NewJob() Job {
	return Job{ID: uuid.NewV4().String(), CreatedAt: time.Now()}
}

func (self Job) NewRun() JobRun {
	taskRuns := make([]TaskRun, len(self.Tasks))
	for i, task := range self.Tasks {
		taskRuns[i] = TaskRun{
			ID:   uuid.NewV4().String(),
			Task: task,
		}
	}

	return JobRun{
		ID:        uuid.NewV4().String(),
		JobID:     self.ID,
		CreatedAt: time.Now(),
		TaskRuns:  taskRuns,
	}
}

func (self Job) Schedules() []Initiator {
	list := []Initiator{}
	for _, initr := range self.Initiators {
		if initr.Type == "cron" {
			list = append(list, initr)
		}
	}
	return list
}

func (self Job) WebAuthorized() bool {
	for _, initr := range self.Initiators {
		if initr.Type == "web" {
			return true
		}
	}
	return false
}

func (self *Time) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	t, err := dateparse.ParseAny(s)
	self.Time = t
	return err
}

func (self *Cron) UnmarshalJSON(b []byte) error {
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
	*self = Cron(s)
	return nil
}
