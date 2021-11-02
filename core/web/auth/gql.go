package auth

import (
	"context"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	clsessions "github.com/smartcontractkit/chainlink/core/sessions"
)

type sessionUserKey struct{}

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

		ctx := SetGQLAuthenticatedUser(c.Request.Context(), user)

		c.Request = c.Request.WithContext(ctx)
	}
}

// SetGQLAuthenticatedUser sets the authenticated user in the context
//
// There shouldn't be a need to do this outside of testing
func SetGQLAuthenticatedUser(ctx context.Context, user clsessions.User) context.Context {
	return context.WithValue(ctx, sessionUserKey{}, &user)
}

// GetGQLAuthenticatedUser extracts the authentication user from a context.
func GetGQLAuthenticatedUser(ctx context.Context) (*clsessions.User, bool) {
	obj := ctx.Value(sessionUserKey{})

	user, ok := obj.(*clsessions.User)

	return user, ok
}
