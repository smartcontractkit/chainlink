package sessions

import (
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/sqlx"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	FindUser() (User, error)
	AuthorizedUserWithSession(sessionID string) (User, error)
	DeleteUser() error
	DeleteUserSession(sessionID string) error
	CreateSession(sr SessionRequest) (string, error)
	ClearNonCurrentSessions(sessionID string) error
	CreateUser(user *User) error
	SetAuthToken(user *User, token *auth.Token) error
	CreateAndSetAuthToken(user *User) (*auth.Token, error)
	DeleteAuthToken(user *User) error
	SetPassword(user *User, newPassword string) error
	Sessions(offset, limit int) ([]Session, error)
	GetUserWebAuthn(email string) ([]WebAuthn, error)
	SaveWebAuthn(token *WebAuthn) error

	FindExternalInitiator(eia *auth.Token) (initiator *bridges.ExternalInitiator, err error)
}

type orm struct {
	db              *sqlx.DB
	sessionDuration time.Duration
	lggr            logger.Logger
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB, sessionDuration time.Duration, lggr logger.Logger) ORM {
	return &orm{db, sessionDuration, lggr.Named("SessionsORM")}
}

// FindUser will return the one API user, or an error.
func (o *orm) FindUser() (User, error) {
	return o.findUser()
}

func (o *orm) findUser() (user User, err error) {
	sql := "SELECT * FROM users ORDER BY created_at desc LIMIT 1"
	err = o.db.Get(&user, sql)
	return
}

// AuthorizedUserWithSession will return the one API user if the Session ID exists
// and hasn't expired, and update session's LastUsed field.
func (o *orm) AuthorizedUserWithSession(sessionID string) (User, error) {
	if len(sessionID) == 0 {
		return User{}, errors.New("Session ID cannot be empty")
	}

	result, err := o.db.Exec("UPDATE sessions SET last_used = now() WHERE id = $1 AND last_used + $2 >= now()", sessionID, o.sessionDuration)
	if err != nil {
		return User{}, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return User{}, err
	}
	if rowsAffected == 0 {
		return User{}, sql.ErrNoRows
	}
	return o.FindUser()
}

// DeleteUser will delete the API User in the db.
func (o *orm) DeleteUser() error {
	ctx, cancel := pg.DefaultQueryCtx()
	defer cancel()
	return pg.SqlxTransaction(ctx, o.db, o.lggr, func(tx pg.Queryer) error {
		if _, err := tx.Exec("DELETE FROM users"); err != nil {
			return err
		}

		_, err := tx.Exec("DELETE FROM sessions")
		return err
	})
}

// DeleteUserSession will erase the session ID for the sole API User.
func (o *orm) DeleteUserSession(sessionID string) error {
	_, err := o.db.Exec("DELETE FROM sessions WHERE id = $1", sessionID)
	return err
}

// GetUserWebAuthn will return a list of structures representing all enrolled WebAuthn
// tokens for the user. This list must be used when logging in (for obvious reasons) but
// must also be used for registration to prevent the user from enrolling the same hardware
// token multiple times.
func (o *orm) GetUserWebAuthn(email string) ([]WebAuthn, error) {
	var uwas []WebAuthn
	err := o.db.Select(&uwas, "SELECT email, public_key_data FROM web_authns WHERE LOWER(email) = $1", strings.ToLower(email))
	if err != nil {
		return uwas, err
	}
	// In the event of not found, there is no MFA on this account and it is not an error
	// so this returns either an empty list or list of WebAuthn rows
	return uwas, nil
}

// CreateSession will check the password in the SessionRequest against
// the hashed API User password in the db. Also will check WebAuthn if it's
// enabled for that user.
func (o *orm) CreateSession(sr SessionRequest) (string, error) {
	user, err := o.FindUser()
	if err != nil {
		return "", err
	}
	lggr := o.lggr.With("user", user.Email)
	lggr.Debugw("Found user")

	// Do email and password check first to prevent extra database look up
	// for MFA tokens leaking if an account has MFA tokens or not.
	if !constantTimeEmailCompare(sr.Email, user.Email) {
		return "", errors.New("Invalid email")
	}

	if !utils.CheckPasswordHash(sr.Password, user.HashedPassword) {
		return "", errors.New("Invalid password")
	}

	// Load all valid MFA tokens associated with user's email
	uwas, err := o.GetUserWebAuthn(user.Email)
	if err != nil {
		// There was an error with the database query
		lggr.Errorf("Could not fetch user's MFA data: %v", err)
		return "", errors.New("MFA Error")
	}

	// No webauthn tokens registered for the current user, so normal authentication is now complete
	if len(uwas) == 0 {
		lggr.Infof("No MFA for user. Creating Session")
		session := NewSession()
		_, err = o.db.Exec("INSERT INTO sessions (id, last_used, created_at) VALUES ($1, now(), now())", session.ID)
		return session.ID, err
	}

	// Next check if this session request includes the required WebAuthn challenge data
	// if not, return a 401 error for the frontend to prompt the user to provide this
	// data in the next round trip request (tap key to include webauthn data on the login page)
	if sr.WebAuthnData == "" {
		lggr.Warnf("Attempted login to MFA user. Generating challenge for user.")
		options, webauthnError := BeginWebAuthnLogin(user, uwas, sr)
		if webauthnError != nil {
			lggr.Errorf("Could not begin WebAuthn verification: %v", err)
			return "", errors.New("MFA Error")
		}

		j, jsonError := json.Marshal(options)
		if jsonError != nil {
			lggr.Errorf("Could not serialize WebAuthn challenge: %v", err)
			return "", errors.New("MFA Error")
		}

		return "", errors.New(string(j))
	}

	// The user is at the final stage of logging in with MFA. We have an
	// attestation back from the user, we now need to verify that it is
	// correct.
	err = FinishWebAuthnLogin(user, uwas, sr)

	if err != nil {
		// The user does have WebAuthn enabled but failed the check
		lggr.Errorf("User sent an invalid attestation: %v", err)
		return "", errors.New("MFA Error")
	}

	lggr.Infof("User passed MFA authentication and login will proceed")
	// This is a success so we can create the sessions
	session := NewSession()
	_, err = o.db.Exec("INSERT INTO sessions (id, last_used, created_at) VALUES ($1, now(), now())", session.ID)
	return session.ID, err
}

const constantTimeEmailLength = 256

func constantTimeEmailCompare(left, right string) bool {
	length := utils.MaxInt(constantTimeEmailLength, len(left), len(right))
	leftBytes := make([]byte, length)
	rightBytes := make([]byte, length)
	copy(leftBytes, left)
	copy(rightBytes, right)
	return subtle.ConstantTimeCompare(leftBytes, rightBytes) == 1
}

// ClearNonCurrentSessions removes all sessions but the id passed in.
func (o *orm) ClearNonCurrentSessions(sessionID string) error {
	_, err := o.db.Exec("DELETE FROM sessions where id != $1", sessionID)
	return err
}

// Creates creates the user.
func (o *orm) CreateUser(user *User) error {
	sql := "INSERT INTO users (email, hashed_password, created_at, updated_at) VALUES ($1, $2, now(), now()) RETURNING *"
	return o.db.Get(user, sql, user.Email, user.HashedPassword)
}

// SetAuthToken updates the user to use the given Authentication Token.
func (o *orm) SetPassword(user *User, newPassword string) error {
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	sql := "UPDATE users SET hashed_password = $1, updated_at = now() WHERE email = $2 RETURNING *"
	return o.db.Get(user, sql, hashedPassword, user.Email)
}

func (o *orm) CreateAndSetAuthToken(user *User) (*auth.Token, error) {
	newToken := auth.NewToken()

	err := o.SetAuthToken(user, newToken)
	if err != nil {
		return nil, err
	}

	return newToken, nil
}

// SetAuthToken updates the user to use the given Authentication Token.
func (o *orm) SetAuthToken(user *User, token *auth.Token) error {
	salt := utils.NewSecret(utils.DefaultSecretSize)
	hashedSecret, err := auth.HashedSecret(token, salt)
	if err != nil {
		return errors.Wrap(err, "user")
	}
	sql := "UPDATE users SET token_salt = $1, token_key = $2, token_hashed_secret = $3, updated_at = now() WHERE email = $4 RETURNING *"
	return o.db.Get(user, sql, salt, token.AccessKey, hashedSecret, user.Email)
}

// DeleteAuthToken clears and disables the users Authentication Token.
func (o *orm) DeleteAuthToken(user *User) error {
	sql := "UPDATE users SET token_salt = '', token_key = '', token_hashed_secret = '', updated_at = now() WHERE email = $1 RETURNING *"
	return o.db.Get(user, sql, user.Email)
}

// SaveWebAuthn saves new WebAuthn token information.
func (o *orm) SaveWebAuthn(token *WebAuthn) error {
	sql := "INSERT INTO web_authns (email, public_key_data) VALUES ($1, $2)"
	_, err := o.db.Exec(sql, token.Email, token.PublicKeyData)
	return err
}

// Sessions returns all sessions limited by the parameters.
func (o *orm) Sessions(offset, limit int) (sessions []Session, err error) {
	sql := `SELECT * FROM sessions ORDER BY created_at, id LIMIT $1 OFFSET $2;`
	if err = o.db.Select(&sessions, sql, limit, offset); err != nil {
		return
	}
	return
}

// NOTE: this is duplicated from the bridges ORM to appease the AuthStorer interface
func (o *orm) FindExternalInitiator(
	eia *auth.Token,
) (*bridges.ExternalInitiator, error) {
	exi := &bridges.ExternalInitiator{}
	err := o.db.Get(exi, `SELECT * FROM external_initiators WHERE access_key = $1`, eia.AccessKey)
	return exi, err
}
