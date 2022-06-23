package resolver

import (
	"context"
	"fmt"

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

// Authenticates the user from the session cookie and asserts at least 'edit_minimal' role.
func authenticateUserCanEditMinimal(ctx context.Context) error {
	session, ok := auth.GetGQLAuthenticatedSession(ctx)
	if !ok {
		return fmt.Errorf("Not permitted with current role")
	}
	if session.User.Role == sessions.UserRoleView {
		return fmt.Errorf("Not permitted with current role")
	}
	return nil
}

// Authenticates the user from the session cookie and asserts at least 'edit' role.
func authenticateUserCanEdit(ctx context.Context) error {
	session, ok := auth.GetGQLAuthenticatedSession(ctx)
	if !ok {
		return fmt.Errorf("Not permitted with current role")
	}
	if session.User.Role == sessions.UserRoleView || session.User.Role == sessions.UserRoleEditMinimal {
		return fmt.Errorf("Not permitted with current role")
	}
	return nil
}

// Authenticates the user from the session cookie and asserts has 'admin' role
func authenticateUserIsAdmin(ctx context.Context) error {
	session, ok := auth.GetGQLAuthenticatedSession(ctx)
	if !ok {
		return fmt.Errorf("Not permitted with current role")
	}
	if session.User.Role != sessions.UserRoleAdmin {
		return fmt.Errorf("Not permitted with current role")
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
