package auth

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/bridges"
	clsessions "github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/static"
)

const (
	// APIKey is the header name for the API token identifier for user authentication.
	APIKey = "X-API-KEY"

	// APISecret is the header name for the API token secret for user authentication.
	APISecret = "X-API-SECRET"

	// SessionName is the session name
	SessionName = "clsession"

	// SessionIDKey is the session ID key in the session map
	SessionIDKey = "clsession_id"

	// SessionUserKey is the User key in the session map
	SessionUserKey = "user"

	// SessionExternalInitiatorKey is the External Initiator key in the session map
	SessionExternalInitiatorKey = "external_initiator"
)

// Authenticator defines the interface to authenticate requests against a
// datastore.
type Authenticator interface {
	AuthorizedUserWithSession(sessionID string) (clsessions.User, error)
	FindExternalInitiator(eia *auth.Token) (*bridges.ExternalInitiator, error)
	FindUser() (clsessions.User, error)
}

// authMethod defines a method which can be used to authenticate a request. This
// can be implemented according to your authentication method (i.e by session,
// token, etc)
type authMethod func(ctx *gin.Context, store Authenticator) error

// AuthenticateBySession authenticates the request by the session cookie.
//
// Implements authMethod
func AuthenticateBySession(c *gin.Context, authr Authenticator) error {
	session := sessions.Default(c)
	sessionID, ok := session.Get(SessionIDKey).(string)
	if !ok {
		return auth.ErrorAuthFailed
	}

	user, err := authr.AuthorizedUserWithSession(sessionID)
	if err != nil {
		return err
	}

	c.Set(SessionUserKey, &user)

	return nil
}

var _ authMethod = AuthenticateBySession

// AuthenticateByToken authenticates a User by their API token.
//
// Implements authMethod
func AuthenticateByToken(c *gin.Context, authr Authenticator) error {
	token := &auth.Token{
		AccessKey: c.GetHeader(APIKey),
		Secret:    c.GetHeader(APISecret),
	}

	user, err := authr.FindUser()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.ErrorAuthFailed
		}

		return err
	}

	ok, err := clsessions.AuthenticateUserByToken(token, &user)
	if err != nil {
		return err
	}
	if !ok {
		return auth.ErrorAuthFailed
	}

	c.Set(SessionUserKey, &user)

	return nil
}

var _ authMethod = AuthenticateByToken

// AuthenticateExternalInitiator authenticates an external initiator request.
//
// Implements authMethod
func AuthenticateExternalInitiator(c *gin.Context, store Authenticator) error {
	eia := &auth.Token{
		AccessKey: c.GetHeader(static.ExternalInitiatorAccessKeyHeader),
		Secret:    c.GetHeader(static.ExternalInitiatorSecretHeader),
	}

	ei, err := store.FindExternalInitiator(eia)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.ErrorAuthFailed
		}

		return errors.Wrap(err, "finding external initiator")
	}

	ok, err := bridges.AuthenticateExternalInitiator(eia, ei)
	if err != nil {
		return err
	}
	if !ok {
		return auth.ErrorAuthFailed
	}

	c.Set(SessionExternalInitiatorKey, ei)

	return nil
}

var _ authMethod = AuthenticateExternalInitiator

// Authenticate is middleware which authenticates the request by attempting to
// authenticate using all the provided methods.
func Authenticate(store Authenticator, methods ...authMethod) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		for _, method := range methods {
			err = method(c, store)
			if !errors.Is(err, auth.ErrorAuthFailed) {
				break
			}
		}
		if err != nil {
			c.Abort()
			jsonAPIError(c, http.StatusUnauthorized, err)

			return
		}

		c.Next()
	}
}

// GetAuthenticatedUser extracts the authentication user from the context.
func GetAuthenticatedUser(c *gin.Context) (*clsessions.User, bool) {
	obj, ok := c.Get(SessionUserKey)
	if !ok {
		return nil, false
	}

	user, ok := obj.(*clsessions.User)

	return user, ok
}

// GetAuthenticatedExternalInitiator extracts the external initiator from the
// context.
func GetAuthenticatedExternalInitiator(c *gin.Context) (*bridges.ExternalInitiator, bool) {
	obj, ok := c.Get(SessionExternalInitiatorKey)
	if !ok {
		return nil, false
	}

	return obj.(*bridges.ExternalInitiator), ok
}
