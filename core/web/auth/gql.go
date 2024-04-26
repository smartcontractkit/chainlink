package auth

import (
	"context"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	clsessions "github.com/smartcontractkit/chainlink/v2/core/sessions"
)

type sessionUserKey struct{}
type GQLSession struct {
	SessionID string
	User      *clsessions.User
}

// AuthenticateGQL middleware checks the session cookie for a user and sets it
// on the request context if it exists. It is the responsibility of each resolver
// to validate whether it requires an authenticated user.
//
// We currently only support GQL authentication by session cookie.
func AuthenticateGQL(authenticator Authenticator, lggr logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		session := sessions.Default(c)
		sessionID, ok := session.Get(SessionIDKey).(string)
		if !ok {
			return
		}

		user, err := authenticator.AuthorizedUserWithSession(ctx, sessionID)
		if err != nil {
			if errors.Is(err, clsessions.ErrUserSessionExpired) {
				lggr.Warnw("Failed to authenticate session", "err", err)
			} else {
				lggr.Errorw("Failed call to AuthorizedUserWithSession, unable to get user", "err", err)
			}
			return
		}

		ctx = WithGQLAuthenticatedSession(c.Request.Context(), user, sessionID)

		c.Request = c.Request.WithContext(ctx)
	}
}

// WithGQLAuthenticatedSession sets the authenticated session in the context
//
// There shouldn't be a need to do this outside of testing
func WithGQLAuthenticatedSession(ctx context.Context, user clsessions.User, sessionID string) context.Context {
	return context.WithValue(
		ctx,
		sessionUserKey{},
		&GQLSession{sessionID, &user},
	)
}

// GetGQLAuthenticatedSession extracts the authentication session from a context.
func GetGQLAuthenticatedSession(ctx context.Context) (*GQLSession, bool) {
	obj := ctx.Value(sessionUserKey{})
	if obj == nil {
		return nil, false
	}

	session, ok := obj.(*GQLSession)

	return session, ok
}
