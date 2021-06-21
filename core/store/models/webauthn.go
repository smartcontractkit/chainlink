package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"

	"github.com/duo-labs/webauthn.io/session"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

// User holds the credentials for API user.
type WebAuthn struct {
	Email         string
	PublicKeyData string
	Settings      string
}

type WebAuthnUser struct {
	Email         string
	WACredentials []webauthn.Credential
}

type WebAuthnConfiguration struct {
	RPID     string
	RPOrigin string
}

func BeginWebAuthnRegistration(user User, uwas []WebAuthn, sessionStore *session.Store, ctx *gin.Context, config WebAuthnConfiguration) (*protocol.CredentialCreation, error) {
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

func FinishWebAuthnRegistration(user User, uwas []WebAuthn, sessionStore *session.Store, ctx *gin.Context, config WebAuthnConfiguration) (*webauthn.Credential, error) {
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
		json.Unmarshal([]byte(v.PublicKeyData), &credential)
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
