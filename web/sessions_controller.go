package web

import (
	"errors"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
)

// SnapshotsController manages Snapshot requests.
type SessionsController struct {
	App *services.ChainlinkApplication
}

// CreateSnapshot begins the job run for the given Assignment ID
// Example:
//  "/assignments/:AID/snapshots"
func (sc *SessionsController) Create(c *gin.Context) {
	session := sessions.Default(c)
	var sr models.SessionRequest
	if err := c.ShouldBindJSON(&sr); err != nil {
		publicError(c, 400, err)
	} else if sid, err := sc.App.GetStore().CheckPasswordForSession(sr); err != nil {
		publicError(c, 400, err) // TODO: I never differentiate between the errors
	} else if err := saveSessionId(session, sid); err != nil {
		c.JSON(200, gin.H{})
	}
}

func (sc *SessionsController) Destroy(c *gin.Context) {
	publicError(c, 404, errors.New("Job not found"))
}

func saveSessionId(session sessions.Session, sessionId string) error {
	session.Set(sessionIdKey, sessionId)
	return session.Save()
}
