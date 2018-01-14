package web

import (
	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
)

type JobsController struct {
	App *services.Application
}

func (jrc *JobsController) Index(c *gin.Context) {
	var jobs []models.Job
	if err := jrc.App.Store.AllByIndex("CreatedAt", &jobs); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, jobs)
	}
}

func (jc *JobsController) Create(c *gin.Context) {
	j := models.NewJob()

	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if err = adapters.Validate(j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if err = jc.App.AddJob(j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, gin.H{"id": j.ID})
	}
}

type JobPresenter struct {
	models.Job
	Runs []models.JobRun `json:"runs,omitempty"`
}

func (jc *JobsController) Show(c *gin.Context) {
	id := c.Param("ID")
	var j models.Job

	if err := jc.App.Store.One("ID", id, &j); err == storm.ErrNotFound {
		c.JSON(404, gin.H{
			"errors": []string{"Job not found."},
		})
	} else if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if runs, err := jc.App.Store.JobRunsFor(j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, JobPresenter{j, runs})
	}
}
