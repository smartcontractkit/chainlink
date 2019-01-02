package web

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
	"github.com/smartcontractkit/chainlink/store/presenters"
)

// AssignmentsController manages Assignment requests.
type AssignmentsController struct {
	App services.Application
}

// Create adds validates, saves, and starts a new JobSpec from the v1
// assignment spec format.
// Example:
//  "<application>/assignments"
func (ac *AssignmentsController) Create(c *gin.Context) {
	var a models.AssignmentSpec

	if err := c.ShouldBindJSON(&a); err != nil {
		publicError(c, 400, err)
	} else if j, err := a.ConvertToJobSpec(); err != nil {
		c.AbortWithError(500, err)
	} else if err = ac.App.AddJob(j); err != nil {
		c.AbortWithError(500, err)
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

	if j, err := ac.App.GetStore().FindJob(id); err == orm.ErrorNotFound {
		publicError(c, 404, errors.New("ID not found"))
	} else if err != nil {
		c.AbortWithError(500, err)
	} else if as, err := models.ConvertToAssignment(j); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.JSON(200, as)
	}
}
