package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/services"
	"go.uber.org/multierr"
)

// SessionsController manages session requests.
type SessionsController struct {
	App services.Application
}

// Create creates a session ID for the given user credentials, and returns it
// in a cookie.
func (sc *SessionsController) Create(c *gin.Context) {
	defer sc.App.WakeSessionReaper()

	session := sessions.Default(c)
	var sr models.SessionRequest
	if err := c.ShouldBindJSON(&sr); err != nil {
		publicError(c, http.StatusBadRequest, fmt.Errorf("error binding json %v", err))
	} else if sid, err := sc.App.GetStore().CreateSession(sr); err != nil {
		publicError(c, http.StatusUnauthorized, err)
	} else if err := saveSessionID(session, sid); err != nil {
		c.AbortWithError(http.StatusInternalServerError, multierr.Append(errors.New("Unable to save session id"), err))
	} else {
		c.JSON(http.StatusOK, gin.H{"authenticated": true})
	}
}

// Destroy erases the session ID for the sole API user.
func (sc *SessionsController) Destroy(c *gin.Context) {
	defer sc.App.WakeSessionReaper()

	session := sessions.Default(c)
	defer session.Clear()
	sessionID, ok := session.Get(SessionIDKey).(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
	} else if err := sc.App.GetStore().DeleteUserSession(sessionID); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
	}
}

func saveSessionID(session sessions.Session, sessionID string) error {
	session.Set(SessionIDKey, sessionID)
	return session.Save()
}
