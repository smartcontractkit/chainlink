package web

import (
	"errors"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
)

// SnapshotsController manages Snapshot requests.
type SnapshotsController struct {
	App *services.ChainlinkApplication
}

// CreateSnapshot begins the job run for the given Assignment ID
// Example:
//  "/assignments/:AID/snapshots"
func (sc *SnapshotsController) CreateSnapshot(c *gin.Context) {
	id := c.Param("AID")

	if j, err := sc.App.Store.FindJob(id); err == storm.ErrNotFound {
		publicError(c, 404, errors.New("Job not found"))
	} else if err != nil {
		c.AbortWithError(500, err)
	} else if jr, err := startJob(j, sc.App.Store, models.JSON{}); err != nil {
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

	if jr, err := sc.App.Store.FindJobRun(id); err == storm.ErrNotFound {
		publicError(c, 404, errors.New("Job not found"))
	} else if err != nil {
		c.AbortWithError(500, err)
	} else {
		c.JSON(200, models.ConvertToSnapshot(jr.Result))
	}
}
