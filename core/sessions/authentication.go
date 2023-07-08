package sessions

import (
	"github.com/pkg/errors"
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

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore

// UserManager interface abstracts the required application calls to a user management backend
// Currently localauth (users table DB) or LDAP server (readonly)
type UserManager interface {
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

	// Multiple auth providers require a local DB user always, for other authentication methods
	// provide interface level access to continue allowing local queries
	LocalAdminListUsers() ([]User, error)
	LocalAdminCreateUser(user *User) error
	LocalAdminFindUser(email string) (User, error)
}
