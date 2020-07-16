package models

import (
	"crypto/subtle"
	"fmt"
	"regexp"
	"time"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/pkg/errors"
)

// User holds the credentials for API user.
type User struct {
	Email             string    `json:"email" gorm:"primary_key"`
	HashedPassword    string    `json:"hashedPassword"`
	CreatedAt         time.Time `json:"createdAt" gorm:"index"`
	TokenKey          string    `json:"tokenKey"`
	TokenSalt         string    `json:"-"`
	TokenHashedSecret string    `json:"-"`
	UpdatedAt         time.Time `json:"-"`
}

// https://davidcel.is/posts/stop-validating-email-addresses-with-regex/
var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// https://security.stackexchange.com/questions/39849/does-bcrypt-have-a-maximum-password-length
const (
	MaxBcryptPasswordLength = 50
)

// NewUser creates a new user by hashing the passed plainPwd with bcrypt.
func NewUser(email, plainPwd string) (User, error) {
	if len(email) == 0 {
		return User{}, errors.New("Must enter an email")
	}

	if !emailRegexp.MatchString(email) {
		return User{}, errors.New("Invalid email format")
	}

	if len(plainPwd) < 8 || len(plainPwd) > MaxBcryptPasswordLength {
		return User{}, fmt.Errorf("must enter a password with 8 - %v characters", MaxBcryptPasswordLength)
	}

	pwd, err := utils.HashPassword(plainPwd)
	if err != nil {
		return User{}, err
	}

	return User{
		Email:          email,
		HashedPassword: pwd,
	}, nil
}

// SessionRequest encapsulates the fields needed to generate a new SessionID,
// including the hashed password.
type SessionRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Session holds the unique id for the authenticated session.
type Session struct {
	ID        string    `json:"id" gorm:"primary_key"`
	LastUsed  time.Time `json:"lastUsed" gorm:"index"`
	CreatedAt time.Time `json:"createdAt" gorm:"index"`
}

// NewSession returns a session instance with ID set to a random ID and
// LastUsed to to now.
func NewSession() Session {
	return Session{
		ID:       utils.NewBytes32ID(),
		LastUsed: time.Now(),
	}
}

// ChangePasswordRequest sets a new password for the current Session's User.
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
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

// DeleteAuthToken clears and disables the users Authentication Token.
func (u *User) DeleteAuthToken() {
	u.TokenKey = ""
	u.TokenSalt = ""
	u.TokenHashedSecret = ""
}

// SetAuthToken updates the user to use the given Authentication Token.
func (u *User) SetAuthToken(token *auth.Token) error {
	salt := utils.NewSecret(utils.DefaultSecretSize)
	hashedSecret, err := auth.HashedSecret(token, salt)
	if err != nil {
		return errors.Wrap(err, "user")
	}
	u.TokenSalt = salt
	u.TokenKey = token.AccessKey
	u.TokenHashedSecret = hashedSecret
	return nil
}

// AuthenticateUserByToken returns true on successful authentication of the
// user against the given Authentication Token.
func AuthenticateUserByToken(token *auth.Token, user *User) (bool, error) {
	hashedSecret, err := auth.HashedSecret(token, user.TokenSalt)
	if err != nil {
		return false, err
	}
	return subtle.ConstantTimeCompare([]byte(hashedSecret), []byte(user.TokenHashedSecret)) == 1, nil
}
