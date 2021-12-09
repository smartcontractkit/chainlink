package auth

import (
	"context"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/bridges"
	clsessions "github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/static"
)

type sessionUserKey struct{}
type GQLSession struct {
	SessionID string
	User      *clsessions.User
}

type sessionExternalInitiatorKey struct{}
type GQLExternalInitiatorSession struct {
	SessionID         string
	ExternalInitiator *bridges.ExternalInitiator
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

		ctx := SetGQLAuthenticatedSession(c.Request.Context(), user, sessionID)

		c.Request = c.Request.WithContext(ctx)
	}
}

// SetGQLAuthenticatedSession sets the authenticated session in the context
//
// There shouldn't be a need to do this outside of testing
func SetGQLAuthenticatedSession(ctx context.Context, user clsessions.User, sessionID string) context.Context {
	return context.WithValue(
		ctx,
		sessionUserKey{},
		&GQLSession{sessionID, &user},
	)
}

// GetGQLAuthenticatedSession extracts the authentication session from a context.
func GetGQLAuthenticatedSession(ctx context.Context) (*GQLSession, bool) {
	obj := ctx.Value(sessionUserKey{})

	session, ok := obj.(*GQLSession)

	return session, ok
}

func AuthenticateExternalInitiatorGQL(authenticator Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionID, ok := session.Get(SessionExternalInitiatorKey).(string)
		if !ok {
			return
		}

		eia := &auth.Token{
			AccessKey: c.GetHeader(static.ExternalInitiatorAccessKeyHeader),
			Secret:    c.GetHeader(static.ExternalInitiatorSecretHeader),
		}

		ei, err := authenticator.FindExternalInitiator(eia)
		if err != nil {
			return
		}

		ok, err = bridges.AuthenticateExternalInitiator(eia, ei)
		if err != nil || !ok {
			return
		}

		ctx := SetGQLAuthenticatedExternalInitiatorSession(c.Request.Context(), *ei, sessionID)

		c.Request = c.Request.WithContext(ctx)
	}
}

// SetGQLAuthenticatedSession sets the authenticated external initiator session in the context
//
// There shouldn't be a need to do this outside of testing
func SetGQLAuthenticatedExternalInitiatorSession(ctx context.Context, ei bridges.ExternalInitiator, sessionID string) context.Context {
	return context.WithValue(
		ctx,
		sessionExternalInitiatorKey{},
		&GQLExternalInitiatorSession{sessionID, &ei},
	)
}

// GetGQLAuthenticatedExternalInitiatorSession extracts the authenticated external initiator session from a context.
func GetGQLAuthenticatedExternalInitiatorSession(ctx context.Context) (*GQLExternalInitiatorSession, bool) {
	obj := ctx.Value(sessionExternalInitiatorKey{})

	session, ok := obj.(*GQLExternalInitiatorSession)

	return session, ok
}
