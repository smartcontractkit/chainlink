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

func (w *WebAuthnController) BeginRegistration(c *gin.Context) {
	ctx := c.Request.Context()
	user, ok := auth.GetAuthenticatedUser(c)
	if !ok {
		jsonAPIError(c, http.StatusInternalServerError, errors.New("failed to obtain current user from context"))
		return
	}

	orm := w.App.AuthenticationProvider()
	uwas, err := orm.GetUserWebAuthn(ctx, user.Email)
	if err != nil {
		w.App.GetLogger().Errorf("failed to obtain current user MFA tokens: error in GetUserWebAuthn: %+v", err)
		jsonAPIError(c, http.StatusInternalServerError, errors.New("Unable to register key"))
		return
	}

	webAuthnConfig := w.App.GetWebAuthnConfiguration()

	options, err := w.inProgressRegistrationsStore.BeginWebAuthnRegistration(*user, uwas, webAuthnConfig)
	if err != nil {
		w.App.GetLogger().Errorf("error in BeginWebAuthnRegistration: %s", err)
		jsonAPIError(c, http.StatusInternalServerError, errors.New("internal Server Error"))
		return
	}

	optionsp := presenters.NewRegistrationSettings(*options)

	jsonAPIResponse(c, optionsp, "settings")
}

func (w *WebAuthnController) FinishRegistration(c *gin.Context) {
	ctx := c.Request.Context()
	user, ok := auth.GetAuthenticatedUser(c)
	if !ok {
		logger.Sugared(w.App.GetLogger()).AssumptionViolationf("failed to obtain current user from context")
		jsonAPIError(c, http.StatusInternalServerError, errors.New("Unable to register key"))
		return
	}

	orm := w.App.AuthenticationProvider()
	uwas, err := orm.GetUserWebAuthn(ctx, user.Email)
	if err != nil {
		w.App.GetLogger().Errorf("failed to obtain current user MFA tokens: error in GetUserWebAuthn: %s", err)
		jsonAPIError(c, http.StatusInternalServerError, errors.New("Unable to register key"))
		return
	}

	webAuthnConfig := w.App.GetWebAuthnConfiguration()

	credential, err := w.inProgressRegistrationsStore.FinishWebAuthnRegistration(*user, uwas, c.Request, webAuthnConfig)
	if err != nil {
		w.App.GetLogger().Errorf("error in FinishWebAuthnRegistration: %s", err)
		jsonAPIError(c, http.StatusBadRequest, errors.New("registration was unsuccessful"))
		return
	}

	if sessions.AddCredentialToUser(ctx, w.App.AuthenticationProvider(), user.Email, credential) != nil {
		w.App.GetLogger().Errorf("Could not save WebAuthn credential to DB for user: %s", user.Email)
		jsonAPIError(c, http.StatusInternalServerError, errors.New("internal Server Error"))
		return
	}

	// Forward registered credentials for audit logs
	credj, err := json.Marshal(credential)
	if err != nil {
		w.App.GetLogger().Errorf("error in Marshal credentials: %s", err)
		jsonAPIError(c, http.StatusBadRequest, errors.New("registration was unsuccessful"))
		return
	}
	w.App.GetAuditLogger().Audit(audit.Auth2FAEnrolled, map[string]interface{}{"email": user.Email, "credential": string(credj)})

	c.String(http.StatusOK, "{}")
}
