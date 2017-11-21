package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type Job struct {
	gorm.Model
	Schedule string    `json:"schedule"`
	Subtasks []Subtask `json:"subtasks"`
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
