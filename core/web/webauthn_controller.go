package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/gin-gonic/gin"

	"github.com/duo-labs/webauthn.io/session"
	"github.com/duo-labs/webauthn/webauthn"
)

// WebAuthnController manages registers new keys as well as authentication
// with those keys
type WebAuthnController struct {
	App      chainlink.Application
	Sessions *session.Store
}

func (c *WebAuthnController) BeginRegistration(ctx *gin.Context) {
	if c.Sessions == nil {
		sessionStore, err := session.NewStore()
		if err != nil {
			jsonAPIError(ctx, http.StatusInternalServerError, errors.New("Internal Server Error"))
			return
		}
		c.Sessions = sessionStore
	}

	user, err := c.App.GetStore().FindUser()
	if err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to obtain current user record: %+v", err))
		return
	}

	uwas, err := c.App.GetStore().GetUserWebAuthn(&user)
	if err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to obtain current user MFA tokens: %+v", err))
		return
	}

	webAuthnConfig := c.App.GetWebAuthnConfiguration()

	options, err := models.BeginWebAuthnRegistration(user, uwas, c.Sessions, ctx, webAuthnConfig)
	if err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, errors.New("Internal Server Error"))
		return
	}

	optionsp := presenters.NewRegistrationSettings(*options)

	jsonAPIResponse(ctx, optionsp, "settings")
	return
}

func (c *WebAuthnController) FinishRegistration(ctx *gin.Context) {
	// This should never be nil at this stage (if it registration will surely fail)
	if c.Sessions == nil {
		jsonAPIError(ctx, http.StatusBadRequest, errors.New("Registration was unsuccessful"))
		return
	}

	user, err := c.App.GetStore().FindUser()
	if err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to obtain current user record: %+v", err))
		return
	}

	uwas, err := c.App.GetStore().GetUserWebAuthn(&user)
	if err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to obtain current user MFA tokens: %+v", err))
		return
	}

	webAuthnConfig := c.App.GetWebAuthnConfiguration()

	credential, err := models.FinishWebAuthnRegistration(user, uwas, c.Sessions, ctx, webAuthnConfig)
	if err != nil {
		jsonAPIError(ctx, http.StatusBadRequest, errors.New("Registration was unsuccessful"))
		return
	}

	if c.addCredentialToUser(user, credential) != nil {
		logger.Errorf("Could not save WebAuthn credential to DB for user: %s", user.Email)
		jsonAPIError(ctx, http.StatusInternalServerError, errors.New("Internal Server Error"))
		return
	}

	ctx.String(http.StatusOK, "{}")
	return
}

func (c *WebAuthnController) addCredentialToUser(user models.User, credential *webauthn.Credential) error {
	credj, err := json.Marshal(credential)
	if err != nil {
		return err
	}

	token := models.WebAuthn{
		Email:         user.Email,
		PublicKeyData: string(credj),
		Settings:      "",
	}
	err = c.App.GetStore().SaveWebAuthn(&token)
	if err != nil {
		logger.Errorf("Database error: %v", err)
	}
	return err
}
