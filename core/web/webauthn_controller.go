package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/web/auth"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// WebAuthnController manages registers new keys as well as authentication
// with those keys
type WebAuthnController struct {
	App                          chainlink.Application
	inProgressRegistrationsStore *sessions.WebAuthnSessionStore
}

func NewWebAuthnController(app chainlink.Application) WebAuthnController {
	return WebAuthnController{
		App:                          app,
		inProgressRegistrationsStore: sessions.NewWebAuthnSessionStore(),
	}
}

func (c *WebAuthnController) BeginRegistration(ctx *gin.Context) {
	user, ok := auth.GetAuthenticatedUser(ctx)
	if !ok {
		jsonAPIError(ctx, http.StatusInternalServerError, errors.New("failed to obtain current user from context"))
		return
	}

	orm := c.App.SessionORM()
	uwas, err := orm.GetUserWebAuthn(user.Email)
	if err != nil {
		c.App.GetLogger().Errorf("failed to obtain current user MFA tokens: error in GetUserWebAuthn: %+v", err)
		jsonAPIError(ctx, http.StatusInternalServerError, errors.New("Unable to register key"))
		return
	}

	webAuthnConfig := c.App.GetWebAuthnConfiguration()

	options, err := c.inProgressRegistrationsStore.BeginWebAuthnRegistration(*user, uwas, webAuthnConfig)
	if err != nil {
		c.App.GetLogger().Errorf("error in BeginWebAuthnRegistration: %s", err)
		jsonAPIError(ctx, http.StatusInternalServerError, errors.New("internal Server Error"))
		return
	}

	optionsp := presenters.NewRegistrationSettings(*options)

	jsonAPIResponse(ctx, optionsp, "settings")
}

func (c *WebAuthnController) FinishRegistration(ctx *gin.Context) {
	user, ok := auth.GetAuthenticatedUser(ctx)
	if !ok {
		logger.Sugared(c.App.GetLogger()).AssumptionViolationf("failed to obtain current user from context")
		jsonAPIError(ctx, http.StatusInternalServerError, errors.New("Unable to register key"))
		return
	}

	orm := c.App.SessionORM()
	uwas, err := orm.GetUserWebAuthn(user.Email)
	if err != nil {
		c.App.GetLogger().Errorf("failed to obtain current user MFA tokens: error in GetUserWebAuthn: %s", err)
		jsonAPIError(ctx, http.StatusInternalServerError, errors.New("Unable to register key"))
		return
	}

	webAuthnConfig := c.App.GetWebAuthnConfiguration()

	credential, err := c.inProgressRegistrationsStore.FinishWebAuthnRegistration(*user, uwas, ctx.Request, webAuthnConfig)
	if err != nil {
		c.App.GetLogger().Errorf("error in FinishWebAuthnRegistration: %s", err)
		jsonAPIError(ctx, http.StatusBadRequest, errors.New("registration was unsuccessful"))
		return
	}

	if sessions.AddCredentialToUser(c.App.SessionORM(), user.Email, credential) != nil {
		c.App.GetLogger().Errorf("Could not save WebAuthn credential to DB for user: %s", user.Email)
		jsonAPIError(ctx, http.StatusInternalServerError, errors.New("internal Server Error"))
		return
	}

	// Forward registered credentials for audit logs
	credj, err := json.Marshal(credential)
	if err != nil {
		c.App.GetLogger().Errorf("error in Marshal credentials: %s", err)
		jsonAPIError(ctx, http.StatusBadRequest, errors.New("registration was unsuccessful"))
		return
	}
	c.App.GetAuditLogger().Audit(audit.Auth2FAEnrolled, map[string]interface{}{"email": user.Email, "credential": string(credj)})

	ctx.String(http.StatusOK, "{}")
}
