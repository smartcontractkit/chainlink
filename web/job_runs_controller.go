package web

import (
	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

type JobRunsController struct {
	App *services.Application
}

func (jrc *JobRunsController) Index(c *gin.Context) {
	id := c.Param("ID")
	jobRuns := []models.JobRun{}

	if err := jrc.App.Store.Where("JobID", id, &jobRuns); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, gin.H{"runs": jobRuns})
	}
}

func (jrc *JobRunsController) Create(c *gin.Context) {
	id := c.Param("JobID")
	if j, err := jrc.App.Store.FindJob(id); err == storm.ErrNotFound {
		c.JSON(404, gin.H{
			"errors": []string{"Job not found"},
		})
	} else if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if !j.WebAuthorized() {
		c.JSON(403, gin.H{
			"errors": []string{"Job not available on web API. Recreate with web initiator."},
		})
	} else {
		jr := startJob(j, jrc.App.Store)
		c.JSON(200, gin.H{"id": jr.ID})
	}
}

func startJob(j models.Job, s *store.Store) models.JobRun {
	jr := j.NewRun()
	go func() {
		if _, err := services.ResumeRun(jr, s); err != nil {
			logger.Panic(err)
		}
	}()
	return jr
}
