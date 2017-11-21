package web

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type JobsController struct{}

type Job struct {
	Schedule string    `json:"schedule"`
	Subtasks []Subtask `json:"subtasks"`
}

type Subtask struct {
	Type   string                 `json:"adapterType"`
	Params map[string]interface{} `json:"adapterParams"`
}

func (jc *JobsController) Create(c *gin.Context) {
	var j Job
	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if _, err = j.valid(); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, gin.H{"id": 1})
	}
}

func (j *Job) valid() (bool, error) {
	for _, s := range j.Subtasks {
		if s.Type != "httpJSON" {
			return false, errors.New(`"` + s.Type + `" is not a supported adapter type.`)
		}
	}
	return true, nil
}
