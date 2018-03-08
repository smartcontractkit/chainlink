package web

import (
	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
)

// JobsController manages Job requests in the node.
type JobsController struct {
	App *services.ChainlinkApplication
}

// Index adds the root of the Jobs to the given context.
// Example:
//  "<application>/jobs"
func (jc *JobsController) Index(c *gin.Context) {
	var jobs []models.JobSpec
	if err := jc.App.Store.AllByIndex("CreatedAt", &jobs); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		pjs := make([]presenters.JobSpec, len(jobs))
		for i, j := range jobs {
			pjs[i] = presenters.JobSpec{JobSpec: j}
		}
		c.JSON(200, pjs)
	}
}

// Create adds the Jobs to the given context.
// Example:
//  "<application>/jobs"
func (jc *JobsController) Create(c *gin.Context) {
	j := models.NewJob()

	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(400, gin.H{
			"errors": []string{err.Error()},
		})
	} else if err = services.ValidateJob(j, jc.App.Store); err != nil {
		c.JSON(400, gin.H{
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

// Show returns the details of a job if it exists.
// Example:
//  "<application>/jobs/:JobID"
func (jc *JobsController) Show(c *gin.Context) {
	id := c.Param("JobID")
	if j, err := jc.App.Store.FindJob(id); err == storm.ErrNotFound {
		c.JSON(404, gin.H{
			"errors": []string{"Job not found."},
		})
	} else if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if runs, err := jc.App.Store.JobRunsFor(j.ID); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, presenters.JobSpec{j, runs})
	}
}
