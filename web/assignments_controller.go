package web

import (
	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
)

// AssignmentsController manages Assignment requests.
type AssignmentsController struct {
	App *services.ChainlinkApplication
}

// Create adds validates, saves, and starts a new JobSpec from the v1
// assignment spec format.
// Example:
//  "<application>/assignments"
func (ac *AssignmentsController) Create(c *gin.Context) {
	var a models.AssignmentSpec

	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(400, gin.H{
			"errors": []string{err.Error()},
		})
	} else if j, err := a.ConvertToJobSpec(); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if err = ac.App.AddJob(j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, presenters.JobSpec{JobSpec: j})
	}
}

// Show returns specified assignment ID from the v1
// assignment spec format.
// Example:
//  "<application>/assignments/:ID"
func (ac *AssignmentsController) Show(c *gin.Context) {
	id := c.Param("ID")

	if j, err := ac.App.Store.FindJob(id); err == storm.ErrNotFound {
		c.JSON(404, gin.H{
			"errors": []string{"ID not found."},
		})
	} else if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if as, err := models.ConvertToAssignment(j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, as)
	}
}

// CreateSnapshot will always return the Assignment ID of a snapshot. It may return all of the snapshot immediately,
// or if the details take time to compute the results of the snapshot will be retrievable with the assignment’s XID.
// Example:
//  "/assignments/:AID/snapshots"
func (ac *AssignmentsController) CreateSnapshot(c *gin.Context) {
	id := c.Param("AID")

	if j, err := ac.App.Store.FindJob(id); err == storm.ErrNotFound {
		c.JSON(404, gin.H{
			"errors": []string{"Job not found"},
		})
	} else if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if jr, err := startJob(j, ac.App.Store, models.JSON{}); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, gin.H{"id": jr.ID})
	}
}

// ShowSnapshot returns snapshot for given ID
// Example:
//  "/snapshots/:ID"
func (ac *AssignmentsController) ShowSnapshot(c *gin.Context) {
	id := c.Param("ID")

	if jr, err := ac.App.Store.FindJobRun(id); err == storm.ErrNotFound {
		c.JSON(404, gin.H{
			"errors": []string{"Job not found"},
		})
	} else if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, models.ConvertToSnapshot(jr.Result))
	}
}
