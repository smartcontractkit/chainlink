package sessions

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"github.com/smartcontractkit/chainlink/core/logger"
	sqlxTypes "github.com/smartcontractkit/sqlx/types"
)

// User holds the credentials for API user.
type WebAuthn struct {
	Email         string
	PublicKeyData sqlxTypes.JSONText
}

// This struct implements the required duo-labs/webauthn/ 'User' interface
// kept seperate from our internal 'User' struct
type WebAuthnUser struct {
	Email         string
	WACredentials []webauthn.Credential
}

type WebAuthnConfiguration struct {
	RPID     string
	RPOrigin string
}

func BeginWebAuthnRegistration(user User, uwas []WebAuthn, sessionStore *WebAuthnSessionStore, ctx *gin.Context, config WebAuthnConfiguration) (*protocol.CredentialCreation, error) {
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Chainlink Operator", // Display Name
		RPID:          config.RPID,          // Generally the domain name
		RPOrigin:      config.RPOrigin,      // The origin URL for WebAuthn requests
	})

	if err != nil {
		return nil, err
	}

	if sessionStore == nil {
		return nil, errors.New("SessionStore must not be nil")
	}

	waUser, err := duoWebAuthUserFromUser(user, uwas)
	if err != nil {
		return nil, err
	}

	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = waUser.CredentialExcludeList()
	}

	// generate PublicKeyCredentialCreationOptions, session data
	options, sessionData, err := webAuthn.BeginRegistration(
		waUser,
		registerOptions,
	)

	if err != nil {
		return nil, err
	}

	userRegistrationIndexKey := fmt.Sprintf("%s-registration", user.Email)
	err = sessionStore.SaveWebauthnSession(userRegistrationIndexKey, sessionData, ctx.Request, ctx.Writer)
	if err != nil {
		return nil, err
	}

	return options, nil
}

func FinishWebAuthnRegistration(user User, uwas []WebAuthn, sessionStore *WebAuthnSessionStore, ctx *gin.Context, config WebAuthnConfiguration) (*webauthn.Credential, error) {
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Chainlink Operator", // Display Name
		RPID:          config.RPID,          // Generally the domain name
		RPOrigin:      config.RPOrigin,      // The origin URL for WebAuthn requests
	})
	if err != nil {
		return nil, err
	}

	if sessionStore == nil {
		return nil, errors.New("SessionStore must not be nil")
	}

	userRegistrationIndexKey := fmt.Sprintf("%s-registration", user.Email)
	sessionData, err := sessionStore.GetWebauthnSession(userRegistrationIndexKey, ctx.Request)
	if err != nil {
		return nil, err
	}

	waUser, err := duoWebAuthUserFromUser(user, uwas)
	if err != nil {
		return nil, err
	}

	credential, err := webAuthn.FinishRegistration(waUser, sessionData, ctx.Request)
	if err != nil {
		logger.Errorf("Finish registration failed %v", err)
		return nil, err
	}

	return credential, nil
}

func BeginWebAuthnLogin(user User, uwas []WebAuthn, sr SessionRequest) (*protocol.CredentialAssertion, error) {
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Chainlink Operator",       // Display Name
		RPID:          sr.WebAuthnConfig.RPID,     // Generally the domain name
		RPOrigin:      sr.WebAuthnConfig.RPOrigin, // The origin URL for WebAuthn requests
	})

	if err != nil {
		return nil, err
	}

	waUser, err := duoWebAuthUserFromUser(user, uwas)
	if err != nil {
		return nil, err
	}

	options, sessionData, err := webAuthn.BeginLogin(waUser)
	if err != nil {
		return nil, err
	}

	userLoginIndexKey := fmt.Sprintf("%s-authentication", user.Email)
	err = sr.SessionStore.SaveWebauthnSession(userLoginIndexKey, sessionData, sr.RequestContext.Request, sr.RequestContext.Writer)
	if err != nil {
		return nil, err
	}

	return options, nil
}

func FinishWebAuthnLogin(user User, uwas []WebAuthn, sr SessionRequest) error {
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Chainlink Operator",       // Display Name
		RPID:          sr.WebAuthnConfig.RPID,     // Generally the domain name
		RPOrigin:      sr.WebAuthnConfig.RPOrigin, // The origin URL for WebAuthn requests
	})

	if err != nil {
		logger.Errorf("Could not create webAuthn structure with RPID: %s and RPOrigin: %s", sr.WebAuthnConfig.RPID, sr.WebAuthnConfig.RPOrigin)
		return err
	}

	credential, err := protocol.ParseCredentialRequestResponseBody(strings.NewReader(sr.WebAuthnData))
	if err != nil {
		return err
	}

	userLoginIndexKey := fmt.Sprintf("%s-authentication", user.Email)
	sessionData, err := sr.SessionStore.GetWebauthnSession(userLoginIndexKey, sr.RequestContext.Request)
	if err != nil {
		return err
	}

	waUser, err := duoWebAuthUserFromUser(user, uwas)
	if err != nil {
		return err
	}

	_, err = webAuthn.ValidateLogin(waUser, sessionData, credential)
	return err
}

// WebAuthnID returns the user's ID
func (u WebAuthnUser) WebAuthnID() []byte {
	return []byte(u.Email)
}

// WebAuthnName returns the user's email
func (u WebAuthnUser) WebAuthnName() string {
	return u.Email
}

// WebAuthnDisplayName returns the user's display name.
// In this case we just return the email
func (u WebAuthnUser) WebAuthnDisplayName() string {
	return u.Email
}

// WebAuthnIcon should be the logo in some form. How it should
// be is currently unclear to me.
func (u WebAuthnUser) WebAuthnIcon() string {
	return ""
}

// WebAuthnCredentials returns credentials owned by the user
func (u WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return u.WACredentials
}

// CredentialExcludeList returns a CredentialDescriptor array filled
// with all the user's credentials to prevent them from re-registering
// keys
func (u WebAuthnUser) CredentialExcludeList() []protocol.CredentialDescriptor {
	credentialExcludeList := []protocol.CredentialDescriptor{}

	for _, cred := range u.WACredentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}

func (u *WebAuthnUser) LoadWebAuthnCredentials(uwas []WebAuthn) error {
	for _, v := range uwas {
		var credential webauthn.Credential
		err := v.PublicKeyData.Unmarshal(&credential)
		if err != nil {
			return fmt.Errorf("Error unmarshalling provided PublicKeyData: %s", err)
		}
		u.WACredentials = append(u.WACredentials, credential)
	}
	return nil
}

func duoWebAuthUserFromUser(user User, uwas []WebAuthn) (WebAuthnUser, error) {
	waUser := WebAuthnUser{
		Email: user.Email,
	}
	err := waUser.LoadWebAuthnCredentials(uwas)

	return waUser, err
}

// DefaultEncryptionKeyLength is the length of the generated encryption keys
// used for session management.
const DefaultEncryptionKeyLength = 32

// WebauthnSession is the name of the session cookie used to manage session-
// related information.
const WebauthnSession = "webauthn-session"

// ErrInsufficientBytesRead is returned in the rare case that an unexpected
// number of bytes are returned from the crypto/rand reader when creating
// session cookie encryption keys.
var ErrInsufficientBytesRead = errors.New("insufficient bytes read")

// ErrMarshal is returned if unexpected data is present in a webauthn session.
var ErrMarshal = errors.New("error unmarshaling data")

// WebAuthnSessionStore is a wrapper around sessions.CookieStore which provides some helper
// methods related to webauthn operations.
type WebAuthnSessionStore struct {
	*sessions.CookieStore
}

// GenerateSecureKey reads and returns n bytes from the crypto/rand reader
func GenerateSecureKey(n int) ([]byte, error) {
	buf := make([]byte, n)
	read, err := rand.Read(buf)
	if err != nil {
		return buf, err
	}
	if read != n {
		return buf, ErrInsufficientBytesRead
	}
	return buf, nil
}

// NewWebAuthnSessionStore returns a new session store.
func NewWebAuthnSessionStore(keyPairs ...[]byte) (*WebAuthnSessionStore, error) {
	// Generate a default encryption key if one isn't provided
	if len(keyPairs) == 0 {
		key, err := GenerateSecureKey(DefaultEncryptionKeyLength)
		if err != nil {
			return nil, err
		}
		keyPairs = append(keyPairs, key)
	}
	store := &WebAuthnSessionStore{
		sessions.NewCookieStore(keyPairs...),
	}
	return store, nil
}

// SaveWebauthnSession marhsals and saves the webauthn data to the provided
// key given the request and responsewriter
func (store *WebAuthnSessionStore) SaveWebauthnSession(key string, data *webauthn.SessionData, r *http.Request, w http.ResponseWriter) error {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return store.Set(key, marshaledData, r, w)
}

// GetWebauthnSession unmarshals and returns the webauthn session information
// from the session cookie.
func (store *WebAuthnSessionStore) GetWebauthnSession(key string, r *http.Request) (webauthn.SessionData, error) {
	sessionData := webauthn.SessionData{}
	session, err := store.Get(r, WebauthnSession)
	if err != nil {
		return sessionData, err
	}
	assertion, ok := session.Values[key].([]byte)
	if !ok {
		return sessionData, ErrMarshal
	}
	err = json.Unmarshal(assertion, &sessionData)
	if err != nil {
		return sessionData, err
	}
	// Delete the value from the session now that it's been read
	delete(session.Values, key)
	return sessionData, nil
}

// Set stores a value to the session with the provided key.
func (store *WebAuthnSessionStore) Set(key string, value interface{}, r *http.Request, w http.ResponseWriter) error {
	session, err := store.Get(r, WebauthnSession)
	if err != nil {
		return fmt.Errorf("error getting session %s", err)
	}
	session.Values[key] = value
	return session.Save(r, w)
}
