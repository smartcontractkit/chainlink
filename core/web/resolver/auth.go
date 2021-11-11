package resolver

import (
	"context"

	"github.com/smartcontractkit/chainlink/core/web/auth"
)

// Authenticates the user from the session cookie.
func authenticateUser(ctx context.Context) (*auth.GQLQUserSession, error) {
	session, ok := auth.GetGQLAuthenticatedUser(ctx)
	if !ok {
		return nil, unauthorizedError{}
	}

	return session, nil
}

type unauthorizedError struct{}

func (e unauthorizedError) Error() string {
	return "Unauthorized"
}

func (e unauthorizedError) Extensions() map[string]interface{} {
	return map[string]interface{}{
		"code": "UNAUTHORIZED",
	}
}
