package models

import (
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/robfig/cron"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-go/models/tasks"
	"time"
)

type Job struct {
	ID        string       `storm:"id,index,unique"`
	Schedule  Schedule     `json:"schedule"`
	Tasks     []tasks.Task `json:"tasks" storm:"inline"`
	CreatedAt time.Time    `storm:"index"`
}

type Schedule struct {
	Cron    Cron   `json:"cron"`
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
