package web

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
)

// SnapshotsController manages Snapshot requests.
type SessionsController struct {
	App *services.ChainlinkApplication
}

// CreateSnapshot begins the job run for the given Assignment ID
// Example:
//  "/assignments/:AID/snapshots"
func (sc *SessionsController) Create(c *gin.Context) {
	publicError(c, 404, errors.New("Job not found"))
}

func (sc *SessionsController) Destroy(c *gin.Context) {
	publicError(c, 404, errors.New("Job not found"))
}
