package web

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type AssignmentsController struct{}

type Assignment struct {
	Schedule string    `json:"schedule"`
	Subtasks []Subtask `json:"subtasks"`
}

type Subtask struct {
	Type   string                 `json:"adapterType"`
	Params map[string]interface{} `json:"adapterParams"`
}

func (ac *AssignmentsController) Create(c *gin.Context) {
	var a Assignment
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if _, err = a.valid(); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, gin.H{"id": 1})
	}
}

func (a *Assignment) valid() (bool, error) {
	for _, s := range a.Subtasks {
		if s.Type != "httpJSON" {
			return false, errors.New(`"` + s.Type + `" is not a supported adapter type.`)
		}
	}
	return true, nil
}
