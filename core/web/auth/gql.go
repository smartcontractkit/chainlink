package auth

import (
	"context"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	clsessions "github.com/smartcontractkit/chainlink/core/sessions"
)

type sessionUserKey struct{}
type GQLQUserSession struct {
	SessionID string
	User      *clsessions.User
}

// AuthenticateGQL middleware checks the session cookie for a user and sets it
// on the request context if it exists. It is the responsiblity of each resolver
// to validate whether it requires an authenticated user.
//
// We currently only support GQL authentication by session cookie.
func AuthenticateGQL(authenticator Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionID, ok := session.Get(SessionIDKey).(string)
		if !ok {
			return
		}

		user, err := authenticator.AuthorizedUserWithSession(sessionID)
		if err != nil {
			return
		}

		ctx := SetGQLAuthenticatedUser(c.Request.Context(), user, sessionID)

		c.Request = c.Request.WithContext(ctx)
	}
}

// SetGQLAuthenticatedUser sets the authenticated user in the context
//
// There shouldn't be a need to do this outside of testing
func SetGQLAuthenticatedUser(ctx context.Context, user clsessions.User, sessionID string) context.Context {
	return context.WithValue(
		ctx,
		sessionUserKey{},
		&GQLQUserSession{sessionID, &user},
	)
}

// GetGQLAuthenticatedUser extracts the authentication user from a context.
func GetGQLAuthenticatedUser(ctx context.Context) (*GQLQUserSession, bool) {
	obj := ctx.Value(sessionUserKey{})

	session, ok := obj.(*GQLQUserSession)

	return session, ok
}
