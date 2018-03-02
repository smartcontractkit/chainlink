package web

import (
	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// JobRunsController manages JobRun requests in the node.
type JobRunsController struct {
	App *services.ChainlinkApplication
}

// Index adds the root of the JobRuns to the given context.
// Example:
//  "<application>/jobs/:JobID/runs"
func (jrc *JobRunsController) Index(c *gin.Context) {
	id := c.Param("JobID")

	if jobRuns, err := jrc.App.Store.JobRunsFor(id); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, gin.H{"runs": jobRuns})
	}
}

// Create starts a new JobRun for the Job specified.
// Example:
//  "<application>/jobs/:JobID/runs"
func (jrc *JobRunsController) Create(c *gin.Context) {
	id := c.Param("JobID")

	var body models.JSON
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
	} else if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if jr, err := startJob(j, jrc.App.Store, body); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, gin.H{"id": jr.ID})
	}
}

// Update marks the JobRun no longer pending, and resumes the Job's pipeline.
// Example:
//  "<application>/runs/:RunID"
func (jrc *JobRunsController) Update(c *gin.Context) {
	id := c.Param("RunID")
	var rr models.RunResult
	if jr, err := jrc.App.Store.FindJobRun(id); err == storm.ErrNotFound {
		c.JSON(404, gin.H{
			"errors": []string{"Job Run not found"},
		})
	} else if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if !jr.Result.Pending {
		c.JSON(405, gin.H{
			"errors": []string{"Cannot resume a job run that isn't pending"},
		})
	} else if err := c.ShouldBindJSON(&rr); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		executeRun(jr, jrc.App.Store, rr)
		c.JSON(200, gin.H{"id": jr.ID})
	}
}

func startJob(j models.Job, s *store.Store, body models.JSON) (models.JobRun, error) {
	jr, err := services.BuildRun(j, s)
	if err != nil {
		return jr, err
	}
	executeRun(jr, s, models.RunResult{Data: body})
	return jr, nil
}

func executeRun(jr models.JobRun, s *store.Store, rr models.RunResult) {
	go func() {
		if _, err := services.ExecuteRun(jr, s, rr); err != nil {
			logger.Error("Web initiator: ", err.Error())
		}
	}()
}
