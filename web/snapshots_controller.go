package web

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
)

// SnapshotsController manages Snapshot requests.
type SnapshotsController struct {
	App services.Application
}

// CreateSnapshot begins the job run for the given Assignment ID
// Example:
//  "/assignments/:AID/snapshots"
func (sc *SnapshotsController) CreateSnapshot(c *gin.Context) {
	id := c.Param("AID")

	if j, err := sc.App.GetStore().FindJob(id); err == orm.ErrorNotFound {
		publicError(c, 404, errors.New("Job not found"))
	} else if err != nil {
		c.AbortWithError(500, err)
	} else if jr, err := services.ExecuteJob(j, j.InitiatorsFor(models.InitiatorWeb)[0], models.RunResult{}, nil, sc.App.GetStore()); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.JSON(200, gin.H{"id": jr.ID})
	}
}

// ShowSnapshot returns snapshot for given ID
// Example:
//  "/snapshots/:ID"
func (sc *SnapshotsController) ShowSnapshot(c *gin.Context) {
	id := c.Param("ID")

	if jr, err := sc.App.GetStore().FindJobRun(id); err == orm.ErrorNotFound {
		publicError(c, 404, errors.New("Job not found"))
	} else if err != nil {
		c.AbortWithError(500, err)
	} else {
		c.JSON(200, models.ConvertToSnapshot(jr.Result))
	}
}
