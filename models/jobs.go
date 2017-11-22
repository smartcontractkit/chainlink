package models

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Job struct {
	ID        string    `storm:"id,index,unique"`
	Schedule  string    `json:"schedule"`
	Subtasks  []Subtask `json:"subtasks" storm:"inline"`
	CreatedAt time.Time `storm:"index"`
}

type Subtask struct {
	Type   string                 `json:"adapterType"`
	Params map[string]interface{} `json:"adapterParams"`
}

func (j *Job) Valid() (bool, error) {
	for _, s := range j.Subtasks {
		if s.Type != "httpJSON" {
			return false, errors.New(`"` + s.Type + `" is not a supported adapter type.`)
		}
	}
	return true, nil
}

func NewJob() Job {
	return Job{ID: uuid.NewV4().String(), CreatedAt: time.Now()}
}
