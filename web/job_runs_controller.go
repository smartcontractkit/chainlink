package web

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/smartcontractkit/chainlink-go/store/models"
)

type JobRunsController struct {
	App *services.Application
}

func (self *JobRunsController) Index(c *gin.Context) {
	id := c.Param("id")
	jobRuns := []models.JobRun{}

	if err := self.App.Store.Where("JobID", id, &jobRuns); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, gin.H{"runs": jobRuns})
	}
}

func (self *JobRunsController) Create(c *gin.Context) {
	jID := c.Param("jobID")
	j := models.Job{}
	if err := self.App.Store.One("ID", jID, &j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if j.ID == "" {
		c.JSON(500, gin.H{
			"errors": []string{"Job not found"},
		})
	} else {
		jr := j.NewRun()
		defer services.StartJob(jr, self.App.Store)
		c.JSON(200, gin.H{"id": jr.ID})
	}
}
