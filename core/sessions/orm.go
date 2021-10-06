package sessions

import (
	"crypto/subtle"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/sqlx"
)

type ORM interface {
	FindUser() (User, error)
	AuthorizedUserWithSession(sessionID string) (User, error)
	DeleteUser() error
	DeleteUserSession(sessionID string) error
	CreateSession(sr SessionRequest) (string, error)
	ClearNonCurrentSessions(sessionID string) error
	CreateUser(user *User) error
	SetAuthToken(user *User, token *auth.Token) error
	DeleteAuthToken(user *User) error
	SetPassword(user *User, newPassword string) error
	Sessions(offset, limit int) ([]Session, error)

	FindExternalInitiator(eia *auth.Token) (initiator *bridges.ExternalInitiator, err error)
}

type orm struct {
	db              *sqlx.DB
	sessionDuration time.Duration
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB, sessionDuration time.Duration) ORM {
	return &orm{db, sessionDuration}
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
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	return postgres.SqlxTransaction(ctx, o.db, func(tx *sqlx.Tx) error {
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

// CreateSession will check the password in the SessionRequest against
// the hashed API User password in the db.
func (o *orm) CreateSession(sr SessionRequest) (string, error) {
	user, err := o.FindUser()
	if err != nil {
		return "", err
	}
	logger.Debugw("Found user", "user", user)

	if !constantTimeEmailCompare(sr.Email, user.Email) {
		return "", errors.New("Invalid email")
	}

	if utils.CheckPasswordHash(sr.Password, user.HashedPassword) {
		session := NewSession()
		_, err := o.db.Exec("INSERT INTO sessions (id, last_used, created_at) VALUES ($1, now(), now())", session.ID)
		return session.ID, err
	}
	return "", errors.New("Invalid password")
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
