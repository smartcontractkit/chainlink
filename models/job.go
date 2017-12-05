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
	ID        string    `storm:"id,index,unique"`
	Schedule  Schedule  `json:"schedule" storm:"inline"`
	Tasks     []Task    `json:"tasks" storm:"inline"`
	CreatedAt time.Time `storm:"index"`
}

type Schedule struct {
	Cron    Cron   `json:"cron" storm:"index"`
	StartAt *Time  `json:"startAt"`
	EndAt   *Time  `json:"endAt"`
	RunAt   []Time `json:"runAt"`
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

func (self Job) Validate() error {
	var err error
	for _, t := range self.Tasks {
		err = t.Validate()
		if err != nil {
			break
		}
	}

	return err
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
	_, err = cron.Parse(s)
	if err != nil {
		return fmt.Errorf("Cron: %v", err)
	}
	*self = Cron(s)
	return nil
}
