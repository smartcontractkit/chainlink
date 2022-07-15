package resolver

import (
	"context"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/web/auth"
)

// Authenticates the user from the session cookie, presence of user inherently provides 'view' access.
func authenticateUser(ctx context.Context) error {
	if _, ok := auth.GetGQLAuthenticatedSession(ctx); !ok {
		return unauthorizedError{}
	}
	return nil
}

// Authenticates the user from the session cookie and asserts at least 'run' role.
func authenticateUserCanRun(ctx context.Context) error {
	session, ok := auth.GetGQLAuthenticatedSession(ctx)
	if !ok {
		return unauthorizedError{}
	}
	if session.User.Role == sessions.UserRoleView {
		return errors.Errorf("Not permitted with current role %s\n", session.User.Role)
	}
	return nil
}

// Authenticates the user from the session cookie and asserts at least 'edit' role.
func authenticateUserCanEdit(ctx context.Context) error {
	session, ok := auth.GetGQLAuthenticatedSession(ctx)
	if !ok {
		return unauthorizedError{}
	}
	switch session.User.Role {
	case sessions.UserRoleView, sessions.UserRoleRun:
		return errors.Errorf("Not permitted with current role %s\n", session.User.Role)
	}
	return nil
}

// Authenticates the user from the session cookie and asserts has 'admin' role
func authenticateUserIsAdmin(ctx context.Context) error {
	session, ok := auth.GetGQLAuthenticatedSession(ctx)
	if !ok {
		return unauthorizedError{}
	}
	if session.User.Role != sessions.UserRoleAdmin {
		return errors.Errorf("Not permitted with current role %s\n", session.User.Role)
	}
	return nil
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
