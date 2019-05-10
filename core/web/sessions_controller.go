package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
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
		jsonAPIError(c, http.StatusBadRequest, fmt.Errorf("error binding json %v", err))
	} else if sid, err := sc.App.GetStore().CreateSession(sr); err != nil {
		jsonAPIError(c, http.StatusUnauthorized, err)
	} else if err := saveSessionID(session, sid); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, multierr.Append(errors.New("Unable to save session id"), err))
	} else {
		jsonAPIResponse(c, Session{Authenticated: true}, "session")
	}
}

// Destroy erases the session ID for the sole API user.
func (sc *SessionsController) Destroy(c *gin.Context) {
	defer sc.App.WakeSessionReaper()

	session := sessions.Default(c)
	defer session.Clear()
	sessionID, ok := session.Get(SessionIDKey).(string)
	if !ok {
		jsonAPIResponse(c, Session{Authenticated: false}, "session")
	} else if err := sc.App.GetStore().DeleteUserSession(sessionID); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponse(c, Session{Authenticated: false}, "session")
	}
}

func saveSessionID(session sessions.Session, sessionID string) error {
	session.Set(SessionIDKey, sessionID)
	return session.Save()
}

type Session struct {
	Authenticated bool `json:"authenticated"`
}

// GetID returns the jsonapi ID.
func (s Session) GetID() string {
	return "sessionID"
}

// GetName returns the collection name for jsonapi.
func (Session) GetName() string {
	return "session"
}

// SetID is used to conform to the UnmarshallIdentifier interface for
// deserializing from jsonapi documents.
func (*Session) SetID(string) error {
	return nil
}
