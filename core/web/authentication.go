package web

import (
	"net/http"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	// APIKey is the header name for the API token identifier for user authentication.
	APIKey = "X-API-KEY"
	// APISecret is the header name for the API token secret for user authentication.
	APISecret = "X-API-SECRET"
	// ExternalInitiatorAccessKeyHeader is the header name for the access key
	// used by external initiators to authenticate
	ExternalInitiatorAccessKeyHeader = "X-Chainlink-EA-AccessKey"
	// ExternalInitiatorSecretHeader is the header name for the secret used by
	// external initiators to authenticate
	ExternalInitiatorSecretHeader = "X-Chainlink-EA-Secret"
)

type AuthStorer interface {
	AuthorizedUserWithSession(sessionID string) (models.User, error)
	FindExternalInitiator(eia *auth.Token) (*models.ExternalInitiator, error)
	FindUser() (models.User, error)
}

type authType func(store AuthStorer, ctx *gin.Context) error

func authenticatedUser(c *gin.Context) (*models.User, bool) {
	obj, ok := c.Get(SessionUserKey)
	if !ok {
		return nil, false
	}
	return obj.(*models.User), ok
}

func AuthenticateExternalInitiator(store AuthStorer, c *gin.Context) error {
	eia := &auth.Token{
		AccessKey: c.GetHeader(ExternalInitiatorAccessKeyHeader),
		Secret:    c.GetHeader(ExternalInitiatorSecretHeader),
	}

	ei, err := store.FindExternalInitiator(eia)
	if errors.Cause(err) == orm.ErrorNotFound {
		return auth.ErrorAuthFailed
	} else if err != nil {
		return errors.Wrap(err, "finding external intiator")
	}

	ok, err := models.AuthenticateExternalInitiator(eia, ei)
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

func authenticatedEI(c *gin.Context) (*models.ExternalInitiator, bool) {
	obj, ok := c.Get(SessionExternalInitiatorKey)
	if !ok {
		return nil, false
	}
	return obj.(*models.ExternalInitiator), ok
}

// AuthenticateByToken authenticates a User by their API token.
func AuthenticateByToken(store AuthStorer, c *gin.Context) error {
	token := &auth.Token{
		AccessKey: c.GetHeader(APIKey),
		Secret:    c.GetHeader(APISecret),
	}

	user, err := store.FindUser()
	if errors.Cause(err) == orm.ErrorNotFound {
		return auth.ErrorAuthFailed
	} else if err != nil {
		return err
	}

	ok, err := models.AuthenticateUserByToken(token, &user)
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
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			c.Next()
		}
	}
}
