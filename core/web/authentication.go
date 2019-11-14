package web

import (
	"chainlink/core/auth"
	"chainlink/core/store"
	"chainlink/core/store/models"
	"chainlink/core/store/orm"
	"net/http"

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

type authType func(store *store.Store, ctx *gin.Context) error

func authenticatedUser(c *gin.Context) (*models.User, bool) {
	obj, ok := c.Get(SessionUserKey)
	if !ok {
		return nil, false
	}
	return obj.(*models.User), ok
}

func AuthenticateExternalInitiator(store *store.Store, c *gin.Context) error {
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

func AuthenticateBySession(store *store.Store, c *gin.Context) error {
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

func RequireAuth(store *store.Store, methods ...authType) gin.HandlerFunc {
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
