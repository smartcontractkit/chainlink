package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// jsonAPIError adds an error to the gin context and sets
// the JSON value of errors.
//
// This is duplicated code, but we plan to deprecate and remove the JSONAPI
// so this is ok for now
func jsonAPIError(c *gin.Context, statusCode int, err error) {
	_ = c.Error(err).SetType(gin.ErrorTypePublic)
	var jsonErr *models.JSONAPIErrors
	if errors.As(err, &jsonErr) {
		c.JSON(statusCode, jsonErr)
		return
	}
	c.JSON(statusCode, models.NewJSONAPIErrorsWith(err.Error()))
}

// addForbiddenErrorHeaders adds custom headers to the 403 (Forbidden) response
// so that they can be parsed by the remote client for friendly/actionable error messages.
//
// The fields are specific because Forbidden error is caused by the user not having the correct role
// for the required action
func addForbiddenErrorHeaders(c *gin.Context, requiredRole string, providedRole string, providedEmail string) {
	c.Header("forbidden-required-role", requiredRole)
	c.Header("forbidden-provided-role", providedRole)
	c.Header("forbidden-provided-email", providedEmail)
}
