package sessions

import (
	"crypto/subtle"
	"fmt"
	"net/mail"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// User holds the credentials for API user.
type User struct {
	Email             string
	HashedPassword    string
	Role              UserRole
	CreatedAt         time.Time
	TokenKey          null.String
	TokenSalt         null.String
	TokenHashedSecret null.String
	UpdatedAt         time.Time
}

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleEdit  UserRole = "edit"
	UserRoleRun   UserRole = "run"
	UserRoleView  UserRole = "view"
)

// https://security.stackexchange.com/questions/39849/does-bcrypt-have-a-maximum-password-length
const (
	MaxBcryptPasswordLength = 50
)

// NewUser creates a new user by hashing the passed plainPwd with bcrypt.
func NewUser(email string, plainPwd string, role UserRole) (User, error) {
	if err := ValidateEmail(email); err != nil {
		return User{}, err
	}

	pwd, err := ValidateAndHashPassword(plainPwd)
	if err != nil {
		return User{}, err
	}

	return User{
		Email:          email,
		HashedPassword: pwd,
		Role:           role,
	}, nil
}

// ValidateEmail is the single point of logic for user email validations
func ValidateEmail(email string) error {
	if len(email) == 0 {
		return errors.New("Must enter an email")
	}
	_, err := mail.ParseAddress(email)
	return err
}

// ValidateAndHashPassword is the single point of logic for user password validations
func ValidateAndHashPassword(plainPwd string) (string, error) {
	if err := utils.VerifyPasswordComplexity(plainPwd); err != nil {
		return "", errors.Wrapf(err, "password insufficiently complex:\n%s", utils.PasswordComplexityRequirements)
	}
	if len(plainPwd) > MaxBcryptPasswordLength {
		return "", errors.Errorf("must enter a password less than %v characters", MaxBcryptPasswordLength)
	}

	pwd, err := utils.HashPassword(plainPwd)
	if err != nil {
		return "", err
	}

	return pwd, nil
}

// GetUserRole is the single point of logic for mapping role string to UserRole
func GetUserRole(role string) (UserRole, error) {
	if role == string(UserRoleAdmin) {
		return UserRoleAdmin, nil
	}
	if role == string(UserRoleEdit) {
		return UserRoleEdit, nil
	}
	if role == string(UserRoleRun) {
		return UserRoleRun, nil
	}
	if role == string(UserRoleView) {
		return UserRoleView, nil
	}

	errStr := fmt.Sprintf(
		"Invalid role: %s. Allowed roles: '%s', '%s', '%s', '%s'.",
		role,
		UserRoleAdmin,
		UserRoleEdit,
		UserRoleRun,
		UserRoleView,
	)
	return UserRole(""), errors.New(errStr)
}

// SessionRequest encapsulates the fields needed to generate a new SessionID,
// including the hashed password.
type SessionRequest struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	WebAuthnData   string `json:"webauthndata"`
	WebAuthnConfig WebAuthnConfiguration
	SessionStore   *WebAuthnSessionStore
}

// Session holds the unique id for the authenticated session.
type Session struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	LastUsed  time.Time `json:"lastUsed"`
	CreatedAt time.Time `json:"createdAt"`
}

// NewSession returns a session instance with ID set to a random ID and
// LastUsed to now.
func NewSession() Session {
	return Session{
		ID:       utils.NewBytes32ID(),
		LastUsed: time.Now(),
	}
}

// Changeauth.TokenRequest is sent when updating a User's authentication token.
type ChangeAuthTokenRequest struct {
	Password string `json:"password"`
}

// GenerateAuthToken randomly generates and sets the users Authentication
// Token.
func (u *User) GenerateAuthToken() (*auth.Token, error) {
	token := auth.NewToken()
	return token, u.SetAuthToken(token)
}

// SetAuthToken updates the user to use the given Authentication Token.
func (u *User) SetAuthToken(token *auth.Token) error {
	salt := utils.NewSecret(utils.DefaultSecretSize)
	hashedSecret, err := auth.HashedSecret(token, salt)
	if err != nil {
		return errors.Wrap(err, "user")
	}
	u.TokenSalt = null.StringFrom(salt)
	u.TokenKey = null.StringFrom(token.AccessKey)
	u.TokenHashedSecret = null.StringFrom(hashedSecret)
	return nil
}

// AuthenticateUserByToken returns true on successful authentication of the
// user against the given Authentication Token.
func AuthenticateUserByToken(token *auth.Token, user *User) (bool, error) {
	hashedSecret, err := auth.HashedSecret(token, user.TokenSalt.ValueOrZero())
	if err != nil {
		return false, err
	}
	return subtle.ConstantTimeCompare([]byte(hashedSecret), []byte(user.TokenHashedSecret.ValueOrZero())) == 1, nil
}
