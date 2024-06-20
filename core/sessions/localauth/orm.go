package localauth

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mathutil"
	"github.com/smartcontractkit/chainlink/v2/core/auth"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type orm struct {
	ds              sqlutil.DataSource
	sessionDuration time.Duration
	lggr            logger.Logger
	auditLogger     audit.AuditLogger
}

// orm implements sessions.AuthenticationProvider and sessions.BasicAdminUsersORM interfaces
var _ sessions.AuthenticationProvider = (*orm)(nil)
var _ sessions.BasicAdminUsersORM = (*orm)(nil)

func NewORM(ds sqlutil.DataSource, sd time.Duration, lggr logger.Logger, auditLogger audit.AuditLogger) sessions.AuthenticationProvider {
	return &orm{
		ds:              ds,
		sessionDuration: sd,
		lggr:            lggr.Named("LocalAuthAuthenticationProviderORM"),
		auditLogger:     auditLogger,
	}
}

// FindUser will attempt to return an API user by email.
func (o *orm) FindUser(ctx context.Context, email string) (sessions.User, error) {
	return o.findUser(ctx, email)
}

// FindUserByAPIToken will attempt to return an API user via the user's table token_key column.
func (o *orm) FindUserByAPIToken(ctx context.Context, apiToken string) (user sessions.User, err error) {
	sql := "SELECT * FROM users WHERE token_key = $1"
	err = o.ds.GetContext(ctx, &user, sql, apiToken)
	return
}

func (o *orm) findUser(ctx context.Context, email string) (user sessions.User, err error) {
	sql := "SELECT * FROM users WHERE lower(email) = lower($1)"
	err = o.ds.GetContext(ctx, &user, sql, email)
	return
}

// ListUsers will load and return all user rows from the db.
func (o *orm) ListUsers(ctx context.Context) (users []sessions.User, err error) {
	sql := "SELECT * FROM users ORDER BY email ASC;"
	err = o.ds.SelectContext(ctx, &users, sql)
	return
}

// findValidSession finds an unexpired session by its ID and returns the associated email.
func (o *orm) findValidSession(ctx context.Context, sessionID string) (email string, err error) {
	if err := o.ds.GetContext(ctx, &email, "SELECT email FROM sessions WHERE id = $1 AND last_used + $2 >= now() FOR UPDATE", sessionID, o.sessionDuration); err != nil {
		o.lggr.Infof("query result: %v", email)
		return email, pkgerrors.Wrap(err, "no matching user for provided session token")
	}
	return email, nil
}

// updateSessionLastUsed updates a session by its ID and sets the LastUsed field to now().
func (o *orm) updateSessionLastUsed(ctx context.Context, sessionID string) error {
	_, err := o.ds.ExecContext(ctx, "UPDATE sessions SET last_used = now() WHERE id = $1", sessionID)
	return err
}

// AuthorizedUserWithSession will return the API user associated with the Session ID if it
// exists and hasn't expired, and update session's LastUsed field.
// AuthorizedUserWithSession will return the API user associated with the Session ID if it
// exists and hasn't expired, and update session's LastUsed field.
func (o *orm) AuthorizedUserWithSession(ctx context.Context, sessionID string) (user sessions.User, err error) {
	if len(sessionID) == 0 {
		return sessions.User{}, sessions.ErrEmptySessionID
	}

	email, err := o.findValidSession(ctx, sessionID)
	if err != nil {
		return sessions.User{}, sessions.ErrUserSessionExpired
	}

	user, err = o.findUser(ctx, email)
	if err != nil {
		return sessions.User{}, sessions.ErrUserSessionExpired
	}

	if err := o.updateSessionLastUsed(ctx, sessionID); err != nil {
		return sessions.User{}, err
	}

	return user, nil
}

// DeleteUser will delete an API User and sessions by email.
func (o *orm) DeleteUser(ctx context.Context, email string) error {
	return sqlutil.TransactDataSource(ctx, o.ds, nil, func(tx sqlutil.DataSource) error {
		// session table rows are deleted on cascade through the user email constraint
		if _, err := tx.ExecContext(ctx, "DELETE FROM users WHERE email = $1", email); err != nil {
			return err
		}
		return nil
	})
}

// DeleteUserSession will delete a session by ID.
func (o *orm) DeleteUserSession(ctx context.Context, sessionID string) error {
	_, err := o.ds.ExecContext(ctx, "DELETE FROM sessions WHERE id = $1", sessionID)
	return err
}

// GetUserWebAuthn will return a list of structures representing all enrolled WebAuthn
// tokens for the user. This list must be used when logging in (for obvious reasons) but
// must also be used for registration to prevent the user from enrolling the same hardware
// token multiple times.
func (o *orm) GetUserWebAuthn(ctx context.Context, email string) ([]sessions.WebAuthn, error) {
	var uwas []sessions.WebAuthn
	err := o.ds.SelectContext(ctx, &uwas, "SELECT email, public_key_data FROM web_authns WHERE LOWER(email) = $1", strings.ToLower(email))
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
func (o *orm) CreateSession(ctx context.Context, sr sessions.SessionRequest) (string, error) {
	user, err := o.FindUser(ctx, sr.Email)
	if err != nil {
		return "", err
	}
	lggr := o.lggr.With("user", user.Email)
	lggr.Debugw("Found user")

	// Do email and password check first to prevent extra database look up
	// for MFA tokens leaking if an account has MFA tokens or not.
	if !constantTimeEmailCompare(strings.ToLower(sr.Email), strings.ToLower(user.Email)) {
		o.auditLogger.Audit(audit.AuthLoginFailedEmail, map[string]interface{}{"email": sr.Email})
		return "", pkgerrors.New("Invalid email")
	}

	if !utils.CheckPasswordHash(sr.Password, user.HashedPassword) {
		o.auditLogger.Audit(audit.AuthLoginFailedPassword, map[string]interface{}{"email": sr.Email})
		return "", pkgerrors.New("Invalid password")
	}

	// Load all valid MFA tokens associated with user's email
	uwas, err := o.GetUserWebAuthn(ctx, user.Email)
	if err != nil {
		// There was an error with the database query
		lggr.Errorf("Could not fetch user's MFA data: %v", err)
		return "", pkgerrors.New("MFA Error")
	}

	// No webauthn tokens registered for the current user, so normal authentication is now complete
	if len(uwas) == 0 {
		lggr.Infof("No MFA for user. Creating Session")
		session := sessions.NewSession()
		_, err = o.ds.ExecContext(ctx, "INSERT INTO sessions (id, email, last_used, created_at) VALUES ($1, $2, now(), now())", session.ID, user.Email)
		o.auditLogger.Audit(audit.AuthLoginSuccessNo2FA, map[string]interface{}{"email": sr.Email})
		return session.ID, err
	}

	// Next check if this session request includes the required WebAuthn challenge data
	// if not, return a 401 error for the frontend to prompt the user to provide this
	// data in the next round trip request (tap key to include webauthn data on the login page)
	if sr.WebAuthnData == "" {
		lggr.Warnf("Attempted login to MFA user. Generating challenge for user.")
		options, webauthnError := sessions.BeginWebAuthnLogin(user, uwas, sr)
		if webauthnError != nil {
			lggr.Errorf("Could not begin WebAuthn verification: %v", webauthnError)
			return "", pkgerrors.New("MFA Error")
		}

		j, jsonError := json.Marshal(options)
		if jsonError != nil {
			lggr.Errorf("Could not serialize WebAuthn challenge: %v", jsonError)
			return "", pkgerrors.New("MFA Error")
		}

		return "", pkgerrors.New(string(j))
	}

	// The user is at the final stage of logging in with MFA. We have an
	// attestation back from the user, we now need to verify that it is
	// correct.
	err = sessions.FinishWebAuthnLogin(user, uwas, sr)

	if err != nil {
		// The user does have WebAuthn enabled but failed the check
		o.auditLogger.Audit(audit.AuthLoginFailed2FA, map[string]interface{}{"email": sr.Email, "error": err})
		lggr.Errorf("User sent an invalid attestation: %v", err)
		return "", pkgerrors.New("MFA Error")
	}

	lggr.Infof("User passed MFA authentication and login will proceed")
	// This is a success so we can create the sessions
	session := sessions.NewSession()
	_, err = o.ds.ExecContext(ctx, "INSERT INTO sessions (id, email, last_used, created_at) VALUES ($1, $2, now(), now())", session.ID, user.Email)
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
func (o *orm) ClearNonCurrentSessions(ctx context.Context, sessionID string) error {
	_, err := o.ds.ExecContext(ctx, "DELETE FROM sessions where id != $1", sessionID)
	return err
}

// CreateUser creates a new API user
func (o *orm) CreateUser(ctx context.Context, user *sessions.User) error {
	sql := "INSERT INTO users (email, hashed_password, role, created_at, updated_at) VALUES ($1, $2, $3, now(), now()) RETURNING *"
	return o.ds.GetContext(ctx, user, sql, strings.ToLower(user.Email), user.HashedPassword, user.Role)
}

// UpdateRole overwrites role field of the user specified by email.
func (o *orm) UpdateRole(ctx context.Context, email, newRole string) (sessions.User, error) {
	var userToEdit sessions.User

	if newRole == "" {
		return userToEdit, pkgerrors.New("user role must be specified")
	}

	err := sqlutil.TransactDataSource(ctx, o.ds, nil, func(tx sqlutil.DataSource) error {
		// First, attempt to load specified user by email
		if err := tx.GetContext(ctx, &userToEdit, "SELECT * FROM users WHERE lower(email) = lower($1)", email); err != nil {
			return pkgerrors.New("no matching user for provided email")
		}

		// Patch validated role
		userRole, err := sessions.GetUserRole(newRole)
		if err != nil {
			return err
		}
		userToEdit.Role = userRole

		_, err = tx.ExecContext(ctx, "DELETE FROM sessions WHERE email = lower($1)", email)
		if err != nil {
			o.lggr.Errorw("Failed to purge user sessions for UpdateRole", "err", err)
			return pkgerrors.New("error updating API user")
		}

		sql := "UPDATE users SET role = $1, updated_at = now() WHERE lower(email) = lower($2) RETURNING *"
		if err := tx.GetContext(ctx, &userToEdit, sql, userToEdit.Role, email); err != nil {
			o.lggr.Errorw("Error updating API user", "err", err)
			return pkgerrors.New("error updating API user")
		}

		return nil
	})

	return userToEdit, err
}

// SetAuthToken updates the user to use the given Authentication Token.
func (o *orm) SetPassword(ctx context.Context, user *sessions.User, newPassword string) error {
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	sql := "UPDATE users SET hashed_password = $1, updated_at = now() WHERE email = $2 RETURNING *"
	return o.ds.GetContext(ctx, user, sql, hashedPassword, user.Email)
}

// TestPassword checks plaintext user provided password with hashed database password, returns nil if matched
func (o *orm) TestPassword(ctx context.Context, email string, password string) error {
	var hashedPassword string
	if err := o.ds.GetContext(ctx, &hashedPassword, "SELECT hashed_password FROM users WHERE lower(email) = lower($1)", email); err != nil {
		return pkgerrors.New("no matching user for provided email")
	}
	if !utils.CheckPasswordHash(password, hashedPassword) {
		return pkgerrors.New("passwords don't match")
	}
	return nil
}

func (o *orm) CreateAndSetAuthToken(ctx context.Context, user *sessions.User) (*auth.Token, error) {
	newToken := auth.NewToken()

	err := o.SetAuthToken(ctx, user, newToken)
	if err != nil {
		return nil, err
	}

	return newToken, nil
}

// SetAuthToken updates the user to use the given Authentication Token.
func (o *orm) SetAuthToken(ctx context.Context, user *sessions.User, token *auth.Token) error {
	salt := utils.NewSecret(utils.DefaultSecretSize)
	hashedSecret, err := auth.HashedSecret(token, salt)
	if err != nil {
		return pkgerrors.Wrap(err, "user")
	}
	sql := "UPDATE users SET token_salt = $1, token_key = $2, token_hashed_secret = $3, updated_at = now() WHERE email = $4 RETURNING *"
	return o.ds.GetContext(ctx, user, sql, salt, token.AccessKey, hashedSecret, user.Email)
}

// DeleteAuthToken clears and disables the users Authentication Token.
func (o *orm) DeleteAuthToken(ctx context.Context, user *sessions.User) error {
	sql := "UPDATE users SET token_salt = '', token_key = '', token_hashed_secret = '', updated_at = now() WHERE email = $1 RETURNING *"
	return o.ds.GetContext(ctx, user, sql, user.Email)
}

// SaveWebAuthn saves new WebAuthn token information.
func (o *orm) SaveWebAuthn(ctx context.Context, token *sessions.WebAuthn) error {
	sql := "INSERT INTO web_authns (email, public_key_data) VALUES ($1, $2)"
	_, err := o.ds.ExecContext(ctx, sql, token.Email, token.PublicKeyData)
	return err
}

// Sessions returns all sessions limited by the parameters.
func (o *orm) Sessions(ctx context.Context, offset, limit int) (sessions []sessions.Session, err error) {
	sql := `SELECT * FROM sessions ORDER BY created_at, id LIMIT $1 OFFSET $2;`
	if err = o.ds.SelectContext(ctx, &sessions, sql, limit, offset); err != nil {
		return
	}
	return
}

// NOTE: this is duplicated from the bridges ORM to appease the AuthStorer interface
func (o *orm) FindExternalInitiator(
	ctx context.Context,
	eia *auth.Token,
) (*bridges.ExternalInitiator, error) {
	exi := &bridges.ExternalInitiator{}
	err := o.ds.GetContext(ctx, exi, `SELECT * FROM external_initiators WHERE access_key = $1`, eia.AccessKey)
	return exi, err
}
