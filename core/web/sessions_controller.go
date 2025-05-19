package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	clsessions "github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/web/auth"
)

// SessionsController manages session requests.
type SessionsController struct {
	App      chainlink.Application
	sessions *clsessions.WebAuthnSessionStore
}

func NewSessionsController(app chainlink.Application) *SessionsController {
	return &SessionsController{app, clsessions.NewWebAuthnSessionStore()}
}

// Create creates a session ID for the given user credentials, and returns it
// in a cookie.
func (sc *SessionsController) Create(c *gin.Context) {
	defer sc.App.WakeSessionReaper()
	ctx := c.Request.Context()
	sc.App.GetLogger().Debugf("TRACE: Starting Session Creation")

	session := sessions.Default(c)
	var sr clsessions.SessionRequest
	if err := c.ShouldBindJSON(&sr); err != nil {
		jsonAPIError(c, http.StatusBadRequest, fmt.Errorf("error binding json %w", err))
		return
	}

	// Does this user have 2FA enabled?
	userWebAuthnTokens, err := sc.App.AuthenticationProvider().GetUserWebAuthn(ctx, sr.Email)
	if err != nil {
		sc.App.GetLogger().Errorf("Error loading user WebAuthn data: %s", err)
		jsonAPIError(c, http.StatusInternalServerError, errors.New("internal Server Error"))
		return
	}

	// If the user has registered MFA tokens, then populate our session store and context
	// required for successful WebAuthn authentication
	if len(userWebAuthnTokens) > 0 {
		sr.SessionStore = sc.sessions
		sr.WebAuthnConfig = sc.App.GetWebAuthnConfiguration()
	}

	sid, err := sc.App.AuthenticationProvider().CreateSession(ctx, sr)
	if err != nil {
		jsonAPIError(c, http.StatusUnauthorized, err)
		return
	}

	if err := saveSessionID(session, sid); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, multierr.Append(errors.New("unable to save session id"), err))
		return
	}

	jsonAPIResponse(c, Session{Authenticated: true}, "session")
}

// Destroy removes the specified session ID from the database.
func (sc *SessionsController) Destroy(c *gin.Context) {
	defer sc.App.WakeSessionReaper()
	ctx := c.Request.Context()

	session := sessions.Default(c)
	defer session.Clear()
	sessionID, ok := session.Get(auth.SessionIDKey).(string)
	if !ok {
		jsonAPIResponse(c, Session{Authenticated: false}, "session")
		return
	}
	if err := sc.App.AuthenticationProvider().DeleteUserSession(ctx, sessionID); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	sc.App.GetAuditLogger().Audit(audit.AuthSessionDeleted, map[string]interface{}{"sessionID": sessionID})
	jsonAPIResponse(c, Session{Authenticated: false}, "session")
}

func saveSessionID(session sessions.Session, sessionID string) error {
	session.Set(auth.SessionIDKey, sessionID)
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
