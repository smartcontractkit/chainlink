package web

import (
	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/store/models"
)

type JobRunsController struct {
	App *services.Application
}

func (self *JobRunsController) Index(c *gin.Context) {
	id := c.Param("ID")
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
	id := c.Param("JobID")
	j := models.Job{}

	if err := self.App.Store.One("ID", id, &j); err == storm.ErrNotFound {
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
	} else if jr, err := startJob(j, self.App.Store); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, gin.H{"id": jr.ID})
	}
}

func startJob(j models.Job, s *store.Store) (models.JobRun, error) {
	jr := j.NewRun()
	if err := services.StartJob(jr, s); err != nil {
		return jr, err
	}
	return jr, nil
}
