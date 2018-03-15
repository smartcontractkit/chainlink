package web

import (
	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
)

// JobSpecsController manages JobSpec requests.
type JobSpecsController struct {
	App *services.ChainlinkApplication
}

// Index lists all of the existing JobSpecs.
// Example:
//  "<application>/specs"
func (jsc *JobSpecsController) Index(c *gin.Context) {
	var jobs []models.JobSpec
	if err := jsc.App.Store.AllByIndex("CreatedAt", &jobs); err != nil {
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

// Create adds validates, saves, and starts a new JobSpec.
// Example:
//  "<application>/specs"
func (jsc *JobSpecsController) Create(c *gin.Context) {
	j := models.NewJob()

	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(400, gin.H{
			"errors": []string{err.Error()},
		})
	} else if err = services.ValidateJob(j, jsc.App.Store); err != nil {
		c.JSON(400, gin.H{
			"errors": []string{err.Error()},
		})
	} else if err = jsc.App.AddJob(j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, presenters.JobSpec{JobSpec: j})
	}
}

// Show returns the details of a JobSpec.
// Example:
//  "<application>/specs/:SpecID"
func (jsc *JobSpecsController) Show(c *gin.Context) {
	id := c.Param("SpecID")
	if j, err := jsc.App.Store.FindJob(id); err == storm.ErrNotFound {
		c.JSON(404, gin.H{
			"errors": []string{"JobSpec not found."},
		})
	} else if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if runs, err := jsc.App.Store.JobRunsFor(j.ID); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, presenters.JobSpec{j, runs})
	}
}
