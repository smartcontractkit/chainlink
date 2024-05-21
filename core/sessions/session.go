package sessions

import (
	"crypto/subtle"
	"time"

	pkgerrors "github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/v2/core/auth"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

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
		return pkgerrors.Wrap(err, "user")
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
