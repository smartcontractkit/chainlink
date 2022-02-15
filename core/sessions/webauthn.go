package sessions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/pkg/errors"
	sqlxTypes "github.com/smartcontractkit/sqlx/types"
)

// User holds the credentials for API user.
type WebAuthn struct {
	Email         string
	PublicKeyData sqlxTypes.JSONText
}

// This struct implements the required duo-labs/webauthn/ 'User' interface
// kept separate from our internal 'User' struct
type WebAuthnUser struct {
	Email         string
	WACredentials []webauthn.Credential
}

type WebAuthnConfiguration struct {
	RPID     string
	RPOrigin string
}

func (store *WebAuthnSessionStore) BeginWebAuthnRegistration(user User, uwas []WebAuthn, config WebAuthnConfiguration) (*protocol.CredentialCreation, error) {
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Chainlink Operator", // Display Name
		RPID:          config.RPID,          // Generally the domain name
		RPOrigin:      config.RPOrigin,      // The origin URL for WebAuthn requests
	})

	if err != nil {
		return nil, err
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
	err = store.SaveWebauthnSession(userRegistrationIndexKey, sessionData)
	if err != nil {
		return nil, err
	}

	return options, nil
}

func (store *WebAuthnSessionStore) FinishWebAuthnRegistration(user User, uwas []WebAuthn, response *http.Request, config WebAuthnConfiguration) (*webauthn.Credential, error) {
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Chainlink Operator", // Display Name
		RPID:          config.RPID,          // Generally the domain name
		RPOrigin:      config.RPOrigin,      // The origin URL for WebAuthn requests
	})
	if err != nil {
		return nil, err
	}

	userRegistrationIndexKey := fmt.Sprintf("%s-registration", user.Email)
	sessionData, err := store.GetWebauthnSession(userRegistrationIndexKey)
	if err != nil {
		return nil, err
	}

	waUser, err := duoWebAuthUserFromUser(user, uwas)
	if err != nil {
		return nil, err
	}

	credential, err := webAuthn.FinishRegistration(waUser, sessionData, response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to FinishRegistration")
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
	err = sr.SessionStore.SaveWebauthnSession(userLoginIndexKey, sessionData)
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
		return errors.Wrapf(err, "failed to create webAuthn structure with RPID: %s and RPOrigin: %s", sr.WebAuthnConfig.RPID, sr.WebAuthnConfig.RPOrigin)
	}

	credential, err := protocol.ParseCredentialRequestResponseBody(strings.NewReader(sr.WebAuthnData))
	if err != nil {
		return err
	}

	userLoginIndexKey := fmt.Sprintf("%s-authentication", user.Email)
	sessionData, err := sr.SessionStore.GetWebauthnSession(userLoginIndexKey)
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
			return fmt.Errorf("error unmarshalling provided PublicKeyData: %s", err)
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

// WebAuthnSessionStore is a wrapper around an in memory key value store which provides some helper
// methods related to webauthn operations.
type WebAuthnSessionStore struct {
	inProgressRegistrations map[string]string
	mu                      sync.Mutex
}

// NewWebAuthnSessionStore returns a new session store.
func NewWebAuthnSessionStore() *WebAuthnSessionStore {
	return &WebAuthnSessionStore{
		inProgressRegistrations: map[string]string{},
	}
}

// SaveWebauthnSession marshals and saves the webauthn data to the provided
// key given the request and responsewriter
func (store *WebAuthnSessionStore) SaveWebauthnSession(key string, data *webauthn.SessionData) error {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	store.put(key, string(marshaledData))
	return nil
}

func (store *WebAuthnSessionStore) put(key, val string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.inProgressRegistrations[key] = val
}

// take returns the val for key, as well as removing it.
func (store *WebAuthnSessionStore) take(key string) (val string, ok bool) {
	store.mu.Lock()
	defer store.mu.Unlock()
	val, ok = store.inProgressRegistrations[key]
	if ok {
		delete(store.inProgressRegistrations, key)
	}
	return
}

// GetWebauthnSession unmarshals and returns the webauthn session information
// from the session cookie.
func (store *WebAuthnSessionStore) GetWebauthnSession(key string) (webauthn.SessionData, error) {
	sessionData := webauthn.SessionData{}

	assertion, ok := store.take(key)
	if !ok {
		return sessionData, fmt.Errorf("assertion not in challenge store")
	}
	err := json.Unmarshal([]byte(assertion), &sessionData)
	if err != nil {
		return sessionData, err
	}
	return sessionData, nil
}
