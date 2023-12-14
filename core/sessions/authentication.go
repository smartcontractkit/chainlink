package sessions

import (
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

//go:generate mockery --quiet --name BasicAdminUsersORM --output ./mocks/ --case=underscore

// BasicAdminUsersORM is the interface that defines the functionality required for supporting basic admin functionality
// adjacent to the identity provider authentication provider implementation. It is currently implemented by the local
// users/sessions ORM containing local admin CLI actions. This is separate from the AuthenticationProvider,
// as local admin management (ie initial core node setup, initial admin user creation), is always
// required no matter what the pluggable AuthenticationProvider implementation is.
type BasicAdminUsersORM interface {
	ListUsers() ([]User, error)
	CreateUser(user *User) error
	FindUser(email string) (User, error)
}

//go:generate mockery --quiet --name AuthenticationProvider --output ./mocks/ --case=underscore

// AuthenticationProvider is an interface that abstracts the required application calls to a user management backend
// Currently localauth (users table DB) or LDAP server (readonly)
type AuthenticationProvider interface {
	FindUser(email string) (User, error)
	FindUserByAPIToken(apiToken string) (User, error)
	ListUsers() ([]User, error)
	AuthorizedUserWithSession(sessionID string) (User, error)
	DeleteUser(email string) error
	DeleteUserSession(sessionID string) error
	CreateSession(sr SessionRequest) (string, error)
	ClearNonCurrentSessions(sessionID string) error
	CreateUser(user *User) error
	UpdateRole(email, newRole string) (User, error)
	SetAuthToken(user *User, token *auth.Token) error
	CreateAndSetAuthToken(user *User) (*auth.Token, error)
	DeleteAuthToken(user *User) error
	SetPassword(user *User, newPassword string) error
	TestPassword(email, password string) error
	Sessions(offset, limit int) ([]Session, error)
	GetUserWebAuthn(email string) ([]WebAuthn, error)
	SaveWebAuthn(token *WebAuthn) error

	FindExternalInitiator(eia *auth.Token) (initiator *bridges.ExternalInitiator, err error)
}
