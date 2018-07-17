package web

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
)

// SessionsController manages session requests.
type SessionsController struct {
	App *services.ChainlinkApplication
}

// Create creates a session ID for the given user credentials, and returns it
// in a cookie.
func (sc *SessionsController) Create(c *gin.Context) {
	session := sessions.Default(c)
	var sr models.SessionRequest
	if err := c.ShouldBindJSON(&sr); err != nil {
		publicError(c, 400, err)
	} else if sid, err := sc.App.GetStore().CheckPasswordForSession(sr); err != nil {
		publicError(c, 400, err) // TODO: I never differentiate between the errors
	} else if err := saveSessionID(session, sid); err != nil {
		c.JSON(200, gin.H{})
	}
}

// Destroy erases the session ID for the sole API user.
func (sc *SessionsController) Destroy(c *gin.Context) {
	err := sc.App.GetStore().DeleteUserSession()
	if err != nil {
		c.AbortWithError(500, err)
	} else {
		c.JSON(200, gin.H{})
	}
}

func saveSessionID(session sessions.Session, sessionID string) error {
	session.Set(SessionIDKey, sessionID)
	return session.Save()
}
