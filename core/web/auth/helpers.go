package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/core/store/models"
)

// jsonAPIError adds an error to the gin context and sets
// the JSON value of errors.
//
// This is duplicated code, but we plan to deprecate and remove the JSONAPI
// so this is ok for now
func jsonAPIError(c *gin.Context, statusCode int, err error) {
	_ = c.Error(err).SetType(gin.ErrorTypePublic)
	switch v := err.(type) {
	case *models.JSONAPIErrors:
		c.JSON(statusCode, v)
	default:
		c.JSON(statusCode, models.NewJSONAPIErrorsWith(err.Error()))
	}
}
