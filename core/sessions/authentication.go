package sessions

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/auth"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
)

// Application config constant options
type AuthenticationProviderName string

const (
	LocalAuth AuthenticationProviderName = "local"
	LDAPAuth  AuthenticationProviderName = "ldap"
)

// ErrUserSessionExpired defines the error triggered when the user session has expired
var ErrUserSessionExpired = errors.New("session missing or expired, please login again")

// ErrNotSupported defines the error where interface functionality doesn't align with the underlying Auth Provider
var ErrNotSupported = fmt.Errorf("functionality not supported with current authentication provider: %w", errors.ErrUnsupported)

// ErrEmptySessionID captures the empty case error message
var ErrEmptySessionID = errors.New("session ID cannot be empty")

// BasicAdminUsersORM is the interface that defines the functionality required for supporting basic admin functionality
// adjacent to the identity provider authentication provider implementation. It is currently implemented by the local
// users/sessions ORM containing local admin CLI actions. This is separate from the AuthenticationProvider,
// as local admin management (ie initial core node setup, initial admin user creation), is always
// required no matter what the pluggable AuthenticationProvider implementation is.
type BasicAdminUsersORM interface {
	ListUsers(ctx context.Context) ([]User, error)
	CreateUser(ctx context.Context, user *User) error
	FindUser(ctx context.Context, email string) (User, error)
}

// AuthenticationProvider is an interface that abstracts the required application calls to a user management backend
// Currently localauth (users table DB) or LDAP server (readonly)
type AuthenticationProvider interface {
	FindUser(ctx context.Context, email string) (User, error)
	FindUserByAPIToken(ctx context.Context, apiToken string) (User, error)
	ListUsers(ctx context.Context) ([]User, error)
	AuthorizedUserWithSession(ctx context.Context, sessionID string) (User, error)
	DeleteUser(ctx context.Context, email string) error
	DeleteUserSession(ctx context.Context, sessionID string) error
	CreateSession(ctx context.Context, sr SessionRequest) (string, error)
	ClearNonCurrentSessions(ctx context.Context, sessionID string) error
	CreateUser(ctx context.Context, user *User) error
	UpdateRole(ctx context.Context, email, newRole string) (User, error)
	SetAuthToken(ctx context.Context, user *User, token *auth.Token) error
	CreateAndSetAuthToken(ctx context.Context, user *User) (*auth.Token, error)
	DeleteAuthToken(ctx context.Context, user *User) error
	SetPassword(ctx context.Context, user *User, newPassword string) error
	TestPassword(ctx context.Context, email, password string) error
	Sessions(ctx context.Context, offset, limit int) ([]Session, error)
	GetUserWebAuthn(ctx context.Context, email string) ([]WebAuthn, error)
	SaveWebAuthn(ctx context.Context, token *WebAuthn) error

	FindExternalInitiator(ctx context.Context, eia *auth.Token) (initiator *bridges.ExternalInitiator, err error)
}
