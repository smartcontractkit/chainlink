package web

import (
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
func (jsc *AssignmentsController) Create(c *gin.Context) {
	var a models.AssignmentSpec

	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(400, gin.H{
			"errors": []string{err.Error()},
		})
	} else if j, err := a.ConvertToJobSpec(); err != nil {
		c.JSON(500, gin.H{
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
