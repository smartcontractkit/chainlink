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
