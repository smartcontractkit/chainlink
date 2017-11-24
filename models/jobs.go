package models

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-go/models/tasks"
	"time"
)

type Job struct {
	ID        string       `storm:"id,index,unique"`
	Schedule  string       `json:"schedule"`
	Tasks     []tasks.Task `json:"tasks" storm:"inline"`
	CreatedAt time.Time    `storm:"index"`
}

func NewJob() Job {
	return Job{ID: uuid.NewV4().String(), CreatedAt: time.Now()}
}

func (j *Job) Valid() (bool, error) {
	for _, t := range j.Tasks {
		if !t.Valid() {
			return false, errors.New(`"` + t.Type + `" is not a supported adapter type.`)
		}
	}
	return true, nil
}
