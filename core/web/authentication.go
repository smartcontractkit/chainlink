package web

import (
	"database/sql"
	"net/http"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/bridges"
	clsessions "github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/static"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	// APIKey is the header name for the API token identifier for user authentication.
	APIKey = "X-API-KEY"
	// APISecret is the header name for the API token secret for user authentication.
	APISecret = "X-API-SECRET"
)

type AuthStorer interface {
	AuthorizedUserWithSession(sessionID string) (clsessions.User, error)
	FindExternalInitiator(eia *auth.Token) (*bridges.ExternalInitiator, error)
	FindUser() (clsessions.User, error)
}

type authType func(store AuthStorer, ctx *gin.Context) error

func authenticatedUser(c *gin.Context) (*clsessions.User, bool) {
	obj, ok := c.Get(SessionUserKey)
	if !ok {
		return nil, false
	}
	return obj.(*clsessions.User), ok
}

func AuthenticateExternalInitiator(store AuthStorer, c *gin.Context) error {
	eia := &auth.Token{
		AccessKey: c.GetHeader(static.ExternalInitiatorAccessKeyHeader),
		Secret:    c.GetHeader(static.ExternalInitiatorSecretHeader),
	}

	ei, err := store.FindExternalInitiator(eia)
	if errors.Is(err, sql.ErrNoRows) {
		return auth.ErrorAuthFailed
	} else if err != nil {
		return errors.Wrap(err, "finding external intiator")
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

var _ authType = AuthenticateExternalInitiator

func authenticatedEI(c *gin.Context) (*bridges.ExternalInitiator, bool) {
	obj, ok := c.Get(SessionExternalInitiatorKey)
	if !ok {
		return nil, false
	}
	return obj.(*bridges.ExternalInitiator), ok
}

// AuthenticateByToken authenticates a User by their API token.
func AuthenticateByToken(store AuthStorer, c *gin.Context) error {
	token := &auth.Token{
		AccessKey: c.GetHeader(APIKey),
		Secret:    c.GetHeader(APISecret),
	}

	user, err := store.FindUser()
	if errors.Is(err, sql.ErrNoRows) {
		return auth.ErrorAuthFailed
	} else if err != nil {
		return err
	}

	ok, err := clsessions.AuthenticateUserByToken(token, &user)
	if err != nil {
		return err
	} else if !ok {
		return auth.ErrorAuthFailed
	}
	c.Set(SessionUserKey, &user)
	return nil
}

var _ authType = AuthenticateByToken

func AuthenticateBySession(store AuthStorer, c *gin.Context) error {
	session := sessions.Default(c)
	sessionID, ok := session.Get(SessionIDKey).(string)
	if !ok {
		return auth.ErrorAuthFailed
	}

	user, err := store.AuthorizedUserWithSession(sessionID)
	if err != nil {
		return err
	}
	c.Set(SessionUserKey, &user)
	return nil
}

var _ authType = AuthenticateBySession

func RequireAuth(store AuthStorer, methods ...authType) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		for _, method := range methods {
			err = method(store, c)
			if err != auth.ErrorAuthFailed {
				break
			}
		}
		if err != nil {
			c.Abort()
			jsonAPIError(c, http.StatusUnauthorized, err)
		} else {
			c.Next()
		}
	}
}
