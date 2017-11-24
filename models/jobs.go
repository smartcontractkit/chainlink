package models

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Job struct {
	ID        string    `storm:"id,index,unique"`
	Schedule  string    `json:"schedule"`
	Tasks     []Task    `json:"tasks" storm:"inline"`
	CreatedAt time.Time `storm:"index"`
}

type Task struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params"`
}

func (j *Job) Valid() (bool, error) {
	for _, s := range j.Tasks {
		if !isValidTask(s) {
			return false, errors.New(`"` + s.Type + `" is not a supported adapter type.`)
		}
	}
	return true, nil
}

func NewJob() Job {
	return Job{ID: uuid.NewV4().String(), CreatedAt: time.Now()}
}

func isValidTask(t Task) bool {
	switch t.Type {
	case "HttpGet":
		return true
	}
	return false
}
