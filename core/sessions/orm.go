package sessions

import (
	"crypto/subtle"
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/logger/audit"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/utils/mathutil"
)

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore

type ORM interface {
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
	Sessions(offset, limit int) ([]Session, error)
	GetUserWebAuthn(email string) ([]WebAuthn, error)
	SaveWebAuthn(token *WebAuthn) error

	FindExternalInitiator(eia *auth.Token) (initiator *bridges.ExternalInitiator, err error)
}

type orm struct {
	q               pg.Q
	sessionDuration time.Duration
	lggr            logger.Logger
	auditLogger     audit.AuditLogger
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB, sd time.Duration, lggr logger.Logger, cfg pg.QConfig, auditLogger audit.AuditLogger) ORM {
	namedLogger := lggr.Named("SessionsORM")
	return &orm{
		q:               pg.NewQ(db, namedLogger, cfg),
		sessionDuration: sd,
		lggr:            lggr.Named("SessionsORM"),
		auditLogger:     auditLogger,
	}
}

// FindUser will attempt to return an API user by email.
func (o *orm) FindUser(email string) (User, error) {
	return o.findUser(email)
}

// FindUserByAPIToken will attempt to return an API user via the user's table token_key column.
func (o *orm) FindUserByAPIToken(apiToken string) (user User, err error) {
	sql := "SELECT * FROM users WHERE token_key = $1"
	err = o.q.Get(&user, sql, apiToken)
	return
}

func (o *orm) findUser(email string) (user User, err error) {
	sql := "SELECT * FROM users WHERE lower(email) = lower($1)"
	err = o.q.Get(&user, sql, email)
	return
}

// ListUsers will load and return all user rows from the db.
func (o *orm) ListUsers() (users []User, err error) {
	sql := "SELECT * FROM users ORDER BY email ASC;"
	err = o.q.Select(&users, sql)
	return
}

// ErrUserSessionExpired defines the error triggered when the user session has expired
var ErrUserSessionExpired = errors.New("session missing or expired, please login again")

// AuthorizedUserWithSession will return the API user associated with the Session ID if it
// exists and hasn't expired, and update session's LastUsed field.
func (o *orm) AuthorizedUserWithSession(sessionID string) (User, error) {
	if len(sessionID) == 0 {
		return User{}, errors.New("Session ID cannot be empty")
	}

	var user User
	err := o.q.Transaction(func(tx pg.Queryer) error {
		// First find user based on session token
		var foundSession struct {
			Email string
			Valid bool
		}
		if err := tx.Get(&foundSession, "SELECT email, last_used + $2 >= now() as valid FROM sessions WHERE id = $1 FOR UPDATE", sessionID, o.sessionDuration); err != nil {
			return errors.Wrap(err, "no matching user for provided session token")
		}

		if !foundSession.Valid {
			return ErrUserSessionExpired
		}

		if err := tx.Get(&user, "SELECT * FROM users WHERE lower(email) = lower($1)", foundSession.Email); err != nil {
			return errors.Wrap(err, "no matching user for provided session email")
		}
		// Session valid and tied to user, update last_used
		_, err := tx.Exec("UPDATE sessions SET last_used = now() WHERE id = $1 AND last_used + $2 >= now()", sessionID, o.sessionDuration)
		if err != nil {
			return errors.Wrap(err, "unable to update sessions table")
		}
		return nil
	})

	if err != nil {
		return User{}, err
	}

	return user, nil
}

// DeleteUser will delete an API User and sessions by email.
func (o *orm) DeleteUser(email string) error {
	return o.q.Transaction(func(tx pg.Queryer) error {
		// session table rows are deleted on cascade through the user email constraint
		if _, err := tx.Exec("DELETE FROM users WHERE email = $1", email); err != nil {
			return err
		}
		return nil
	})
}

// DeleteUserSession will delete a session by ID.
func (o *orm) DeleteUserSession(sessionID string) error {
	_, err := o.q.Exec("DELETE FROM sessions WHERE id = $1", sessionID)
	return err
}

// GetUserWebAuthn will return a list of structures representing all enrolled WebAuthn
// tokens for the user. This list must be used when logging in (for obvious reasons) but
// must also be used for registration to prevent the user from enrolling the same hardware
// token multiple times.
func (o *orm) GetUserWebAuthn(email string) ([]WebAuthn, error) {
	var uwas []WebAuthn
	err := o.q.Select(&uwas, "SELECT email, public_key_data FROM web_authns WHERE LOWER(email) = $1", strings.ToLower(email))
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
	user, err := o.FindUser(sr.Email)
	if err != nil {
		return "", err
	}
	lggr := o.lggr.With("user", user.Email)
	lggr.Debugw("Found user")

	// Do email and password check first to prevent extra database look up
	// for MFA tokens leaking if an account has MFA tokens or not.
	if !constantTimeEmailCompare(strings.ToLower(sr.Email), strings.ToLower(user.Email)) {
		o.auditLogger.Audit(audit.AuthLoginFailedEmail, map[string]interface{}{"email": sr.Email})
		return "", errors.New("Invalid email")
	}

	if !utils.CheckPasswordHash(sr.Password, user.HashedPassword) {
		o.auditLogger.Audit(audit.AuthLoginFailedPassword, map[string]interface{}{"email": sr.Email})
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
		_, err = o.q.Exec("INSERT INTO sessions (id, email, last_used, created_at) VALUES ($1, $2, now(), now())", session.ID, user.Email)
		o.auditLogger.Audit(audit.AuthLoginSuccessNo2FA, map[string]interface{}{"email": sr.Email})
		return session.ID, err
	}

	// Next check if this session request includes the required WebAuthn challenge data
	// if not, return a 401 error for the frontend to prompt the user to provide this
	// data in the next round trip request (tap key to include webauthn data on the login page)
	if sr.WebAuthnData == "" {
		lggr.Warnf("Attempted login to MFA user. Generating challenge for user.")
		options, webauthnError := BeginWebAuthnLogin(user, uwas, sr)
		if webauthnError != nil {
			lggr.Errorf("Could not begin WebAuthn verification: %v", webauthnError)
			return "", errors.New("MFA Error")
		}

		j, jsonError := json.Marshal(options)
		if jsonError != nil {
			lggr.Errorf("Could not serialize WebAuthn challenge: %v", jsonError)
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
		o.auditLogger.Audit(audit.AuthLoginFailed2FA, map[string]interface{}{"email": sr.Email, "error": err})
		lggr.Errorf("User sent an invalid attestation: %v", err)
		return "", errors.New("MFA Error")
	}

	lggr.Infof("User passed MFA authentication and login will proceed")
	// This is a success so we can create the sessions
	session := NewSession()
	_, err = o.q.Exec("INSERT INTO sessions (id, email, last_used, created_at) VALUES ($1, $2, now(), now())", session.ID, user.Email)
	if err != nil {
		return "", err
	}

	// Forward registered credentials for audit logs
	uwasj, err := json.Marshal(uwas)
	if err != nil {
		lggr.Errorf("error in Marshal credentials: %s", err)
	} else {
		o.auditLogger.Audit(audit.AuthLoginSuccessWith2FA, map[string]interface{}{"email": sr.Email, "credential": string(uwasj)})
	}

	return session.ID, nil
}

const constantTimeEmailLength = 256

func constantTimeEmailCompare(left, right string) bool {
	length := mathutil.Max(constantTimeEmailLength, len(left), len(right))
	leftBytes := make([]byte, length)
	rightBytes := make([]byte, length)
	copy(leftBytes, left)
	copy(rightBytes, right)
	return subtle.ConstantTimeCompare(leftBytes, rightBytes) == 1
}

// ClearNonCurrentSessions removes all sessions but the id passed in.
func (o *orm) ClearNonCurrentSessions(sessionID string) error {
	_, err := o.q.Exec("DELETE FROM sessions where id != $1", sessionID)
	return err
}

// CreateUser creates a new API user
func (o *orm) CreateUser(user *User) error {
	sql := "INSERT INTO users (email, hashed_password, role, created_at, updated_at) VALUES ($1, $2, $3, now(), now()) RETURNING *"
	return o.q.Get(user, sql, strings.ToLower(user.Email), user.HashedPassword, user.Role)
}

// UpdateRole overwrites role field of the user specified by email.
func (o *orm) UpdateRole(email, newRole string) (User, error) {
	var userToEdit User

	if newRole == "" {
		return userToEdit, errors.New("user role must be specified")
	}

	err := o.q.Transaction(func(tx pg.Queryer) error {
		// First, attempt to load specified user by email
		if err := tx.Get(&userToEdit, "SELECT * FROM users WHERE lower(email) = lower($1)", email); err != nil {
			return errors.New("no matching user for provided email")
		}

		// Patch validated role
		userRole, err := GetUserRole(newRole)
		if err != nil {
			return err
		}
		userToEdit.Role = userRole

		_, err = tx.Exec("DELETE FROM sessions WHERE email = lower($1)", email)
		if err != nil {
			o.lggr.Errorf("Failed to purge user sessions for UpdateRole", "err", err)
			return errors.New("error updating API user")
		}

		sql := "UPDATE users SET role = $1, updated_at = now() WHERE lower(email) = lower($2) RETURNING *"
		if err := tx.Get(&userToEdit, sql, userToEdit.Role, email); err != nil {
			o.lggr.Errorf("Error updating API user", "err", err)
			return errors.New("error updating API user")
		}

		return nil
	})

	return userToEdit, err
}

// SetAuthToken updates the user to use the given Authentication Token.
func (o *orm) SetPassword(user *User, newPassword string) error {
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	sql := "UPDATE users SET hashed_password = $1, updated_at = now() WHERE email = $2 RETURNING *"
	return o.q.Get(user, sql, hashedPassword, user.Email)
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
	return o.q.Get(user, sql, salt, token.AccessKey, hashedSecret, user.Email)
}

// DeleteAuthToken clears and disables the users Authentication Token.
func (o *orm) DeleteAuthToken(user *User) error {
	sql := "UPDATE users SET token_salt = '', token_key = '', token_hashed_secret = '', updated_at = now() WHERE email = $1 RETURNING *"
	return o.q.Get(user, sql, user.Email)
}

// SaveWebAuthn saves new WebAuthn token information.
func (o *orm) SaveWebAuthn(token *WebAuthn) error {
	sql := "INSERT INTO web_authns (email, public_key_data) VALUES ($1, $2)"
	_, err := o.q.Exec(sql, token.Email, token.PublicKeyData)
	return err
}

// Sessions returns all sessions limited by the parameters.
func (o *orm) Sessions(offset, limit int) (sessions []Session, err error) {
	sql := `SELECT * FROM sessions ORDER BY created_at, id LIMIT $1 OFFSET $2;`
	if err = o.q.Select(&sessions, sql, limit, offset); err != nil {
		return
	}
	return
}

// NOTE: this is duplicated from the bridges ORM to appease the AuthStorer interface
func (o *orm) FindExternalInitiator(
	eia *auth.Token,
) (*bridges.ExternalInitiator, error) {
	exi := &bridges.ExternalInitiator{}
	err := o.q.Get(exi, `SELECT * FROM external_initiators WHERE access_key = $1`, eia.AccessKey)
	return exi, err
}
