package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	sqlxTypes "github.com/jmoiron/sqlx/types"
	pkgerrors "github.com/pkg/errors"
)

// WebAuthn holds the credentials for API user.
type WebAuthn struct {
	Email         string
	PublicKeyData sqlxTypes.JSONText
}

// WebAuthnUser implements the required duo-labs/webauthn/ 'User' interface
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

	userRegistrationIndexKey := user.Email + "-registration"
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

	userRegistrationIndexKey := user.Email + "-registration"
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
		return nil, pkgerrors.Wrap(err, "failed to FinishRegistration")
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

	userLoginIndexKey := user.Email + "-authentication"
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
		return pkgerrors.Wrapf(err, "failed to create webAuthn structure with RPID: %s and RPOrigin: %s", sr.WebAuthnConfig.RPID, sr.WebAuthnConfig.RPOrigin)
	}

	credential, err := protocol.ParseCredentialRequestResponseBody(strings.NewReader(sr.WebAuthnData))
	if err != nil {
		return err
	}

	userLoginIndexKey := user.Email + "-authentication"
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
			return fmt.Errorf("error unmarshalling provided PublicKeyData: %w", err)
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
// from the session cookie, which is removed.
func (store *WebAuthnSessionStore) GetWebauthnSession(key string) (data webauthn.SessionData, err error) {
	assertion, ok := store.take(key)
	if !ok {
		err = pkgerrors.New("assertion not in challenge store")
		return
	}
	err = json.Unmarshal([]byte(assertion), &data)
	return
}

func AddCredentialToUser(ctx context.Context, ap AuthenticationProvider, email string, credential *webauthn.Credential) error {
	credj, err := json.Marshal(credential)
	if err != nil {
		return err
	}

	token := WebAuthn{
		Email:         email,
		PublicKeyData: sqlxTypes.JSONText(credj),
	}
	return ap.SaveWebAuthn(ctx, &token)
}
