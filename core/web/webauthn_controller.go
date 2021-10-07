package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	sqlxTypes "github.com/smartcontractkit/sqlx/types"
)

// WebAuthnController manages registers new keys as well as authentication
// with those keys
type WebAuthnController struct {
	App                          chainlink.Application
	InProgressRegistrationsStore *sessions.WebAuthnSessionStore
}

func (c *WebAuthnController) BeginRegistration(ctx *gin.Context) {
	if c.InProgressRegistrationsStore == nil {
		c.InProgressRegistrationsStore = sessions.NewWebAuthnSessionStore()
	}

	orm := c.App.SessionORM()
	user, err := orm.FindUser()
	if err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to obtain current user record: %+v", err))
		return
	}

	uwas, err := orm.GetUserWebAuthn(user.Email)
	if err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to obtain current user MFA tokens: %+v", err))
		return
	}

	webAuthnConfig := c.App.GetWebAuthnConfiguration()

	options, err := sessions.BeginWebAuthnRegistration(user, uwas, c.InProgressRegistrationsStore, ctx, webAuthnConfig)
	if err != nil {
		logger.Errorf("error in BeginWebAuthnRegistration: %s", err)
		jsonAPIError(ctx, http.StatusInternalServerError, errors.New("Internal Server Error"))
		return
	}

	optionsp := presenters.NewRegistrationSettings(*options)

	jsonAPIResponse(ctx, optionsp, "settings")
}

func (c *WebAuthnController) FinishRegistration(ctx *gin.Context) {
	// This should never be nil at this stage (if it registration will surely fail)
	if c.InProgressRegistrationsStore == nil {
		jsonAPIError(ctx, http.StatusBadRequest, errors.New("Registration was unsuccessful"))
		return
	}

	orm := c.App.SessionORM()
	user, err := orm.FindUser()
	if err != nil {
		logger.Errorf("error finding user: %s", err)
		jsonAPIError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to obtain current user record: %+v", err))
		return
	}

	uwas, err := orm.GetUserWebAuthn(user.Email)
	if err != nil {
		logger.Errorf("error in GetUserWebAuthn: %s", err)
		jsonAPIError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to obtain current user MFA tokens: %+v", err))
		return
	}

	webAuthnConfig := c.App.GetWebAuthnConfiguration()

	credential, err := sessions.FinishWebAuthnRegistration(user, uwas, c.InProgressRegistrationsStore, ctx, webAuthnConfig)
	if err != nil {
		logger.Errorf("error in FinishWebAuthnRegistration: %s", err)
		jsonAPIError(ctx, http.StatusBadRequest, errors.New("Registration was unsuccessful"))
		return
	}

	if c.addCredentialToUser(user, credential) != nil {
		logger.Errorf("Could not save WebAuthn credential to DB for user: %s", user.Email)
		jsonAPIError(ctx, http.StatusInternalServerError, errors.New("Internal Server Error"))
		return
	}

	ctx.String(http.StatusOK, "{}")
}

func (c *WebAuthnController) addCredentialToUser(user sessions.User, credential *webauthn.Credential) error {
	credj, err := json.Marshal(credential)
	if err != nil {
		return err
	}

	token := sessions.WebAuthn{
		Email:         user.Email,
		PublicKeyData: sqlxTypes.JSONText(credj),
	}
	err = c.App.SessionORM().SaveWebAuthn(&token)
	if err != nil {
		logger.Errorf("Database error: %v", err)
	}
	return err
}
