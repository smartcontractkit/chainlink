package web

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/store/models"
)

// StatusCodeForError returns an http status code for an error type.
func StatusCodeForError(err interface{}) int {
	switch err.(type) {
	case *models.ValidationError:
		return 400
	default:
		return 500
	}
}

// publicError adds an error to the gin context and sets
// the JSON value of errors.
func publicError(c *gin.Context, statusCode int, err error) {
	c.Error(err).SetType(gin.ErrorTypePublic)
	switch v := err.(type) {
	case *models.JSONAPIErrors:
		c.JSON(statusCode, v)
	default:
		c.JSON(statusCode, models.NewJSONAPIErrorsWith(err.Error()))
	}
}
